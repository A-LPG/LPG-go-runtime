package lpg2

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GLRParser is a generalized LR driver for LPG GLR conflict tables. It keeps
// shared stack prefixes in a GSS, records reductions in an SPPF, and projects
// compatible AST alternatives through GetNextAst.
type GLRParser struct {
	*Stacks

	monitor         Monitor
	START_STATE     int
	NUM_RULES       int
	NT_OFFSET       int
	LA_STATE_OFFSET int
	ACCEPT_ACTION   int
	ERROR_ACTION    int

	tokStream TokenStream
	prs       ParseTable
	ra        RuleAction

	takingActions  bool
	currentAction  int
	lastToken      int
	parseStackRoot int
	frameTop       int
	frameLocation  []int
	frameParse     []interface{}

	familyCache map[string]IAst
	forestCache map[string]IAst
	gssNodes    map[string]*GssNode
	sppfNodes   map[string]*SppfNode
	sppfRoot    *SppfNode
	sppfSymbols int
}

var glrNullResult = &struct{}{}

type glrRecoverAction interface {
	SetRecoverParser(*BacktrackingParser)
}

func NewGLRParser(tokStream TokenStream, prs ParseTable, ra RuleAction,
	monitor Monitor) (*GLRParser, error) {
	my := &GLRParser{Stacks: NewStacks()}
	if err := my.Reset(tokStream, prs, ra, monitor); err != nil {
		return nil, err
	}
	return my, nil
}

func (my *GLRParser) lookahead(act int, token int) int {
	act = my.prs.LookAhead(act-my.LA_STATE_OFFSET, my.tokStream.GetKind(token))
	if act > my.LA_STATE_OFFSET {
		return my.lookahead(act, my.tokStream.GetNext(token))
	}
	return act
}

// tAction acts on sym in state, with lookahead beginning after curtok.
func (my *GLRParser) tAction(state int, sym int, curtok int) int {
	act := my.prs.TAction(state, sym)
	if act > my.LA_STATE_OFFSET {
		return my.lookahead(act, my.tokStream.GetNext(curtok))
	}
	return act
}

func (my *GLRParser) expandConflict(act int) []int {
	out := make([]int, 0, 2)
	for i := act; ; i++ {
		candidate := my.prs.BaseAction(i)
		if candidate == 0 {
			break
		}
		out = append(out, candidate)
	}
	return out
}

func (my *GLRParser) GetCurrentRule() int {
	return my.currentAction
}

func (my *GLRParser) GetToken(i int) int {
	if my.takingActions {
		return my.frameLocation[my.frameTop+(i-1)]
	}
	return my.Stacks.GetToken(i)
}

func (my *GLRParser) GetSym(i int) interface{} {
	if my.takingActions {
		return my.frameParse[my.frameTop+(i-1)]
	}
	return my.Stacks.GetSym(i)
}

func (my *GLRParser) SetSym1(ast interface{}) {
	if my.takingActions {
		my.frameParse[my.frameTop] = ast
		return
	}
	my.Stacks.SetSym1(ast)
}

func (my *GLRParser) GetFirstToken() int {
	return my.GetToken(1)
}

func (my *GLRParser) GetFirstTokenAt(i int) int {
	return my.GetToken(i)
}

func (my *GLRParser) GetLastToken() int {
	return my.lastToken
}

func (my *GLRParser) GetLastTokenAt(i int) int {
	if i >= my.prs.Rhs(my.currentAction) {
		return my.lastToken
	}
	return my.tokStream.GetPrevious(my.GetToken(i + 1))
}

// GetSppfRoot returns the root from the last successful error-free parse.
func (my *GLRParser) GetSppfRoot() *SppfNode {
	return my.sppfRoot
}

func (my *GLRParser) GetSppfSymbolCount() int {
	return my.sppfSymbols
}

func (my *GLRParser) SetMonitor(monitor Monitor) {
	my.monitor = monitor
}

func (my *GLRParser) Reset1() {
	my.takingActions = false
	my.sppfRoot = nil
	my.sppfSymbols = 0
}

func (my *GLRParser) Reset2(tokStream TokenStream, monitor Monitor) {
	my.monitor = monitor
	my.tokStream = tokStream
	my.Reset1()
}

func (my *GLRParser) Reset(tokStream TokenStream, prs ParseTable,
	ra RuleAction, monitor Monitor) error {
	if prs != nil {
		my.prs = prs
		my.START_STATE = prs.GetStartState()
		my.NUM_RULES = prs.GetNumRules()
		my.NT_OFFSET = prs.GetNtOffset()
		my.LA_STATE_OFFSET = prs.GetLaStateOffset()
		my.ACCEPT_ACTION = prs.GetAcceptAction()
		my.ERROR_ACTION = prs.GetErrorAction()
		if !prs.IsValidForParser() {
			return NewBadParseSymFileException("")
		}
		if !prs.IsGLR() {
			return NewNotGLRParseTableException("")
		}
	}
	if ra != nil {
		my.ra = ra
	}
	if tokStream == nil {
		my.Reset1()
		return nil
	}
	my.Reset2(tokStream, monitor)
	return nil
}

// Parse parses from the grammar start symbol. A positive error count falls
// back to the backtracking recovery driver if the GLR parse fails.
func (my *GLRParser) Parse(maxErrorCount int) (interface{}, error) {
	return my.ParseEntry(0, maxErrorCount)
}

func (my *GLRParser) ParseEntry(markerKind int,
	maxErrorCount int) (interface{}, error) {
	result, err := my.parseEntryNoRepair(markerKind)
	if err == nil || maxErrorCount <= 0 {
		return result, err
	}
	if _, ok := err.(*BadParseException); !ok {
		return nil, err
	}

	bt, btErr := NewBacktrackingParser(my.tokStream, my.prs, my.ra, my.monitor)
	if btErr != nil {
		return nil, btErr
	}
	if action, ok := my.ra.(glrRecoverAction); ok {
		action.SetRecoverParser(bt)
		defer action.SetRecoverParser(nil)
	}
	return bt.FuzzyParseEntry(markerKind, maxErrorCount)
}

func (my *GLRParser) parseEntryNoRepair(markerKind int) (interface{}, error) {
	my.tokStream.Reset()
	my.familyCache = make(map[string]IAst)
	my.forestCache = make(map[string]IAst)
	my.gssNodes = make(map[string]*GssNode)
	my.sppfNodes = make(map[string]*SppfNode)
	my.sppfRoot = nil
	my.sppfSymbols = 0

	firstTok := my.tokStream.GetToken()
	prev := my.tokStream.GetPrevious(firstTok)
	startTok := firstTok
	startKind := my.tokStream.GetKind(firstTok)
	my.parseStackRoot = 0
	if markerKind != 0 {
		startTok = prev
		startKind = markerKind
		my.parseStackRoot = 1
	}

	start := &glrConfig{
		stateStackTop: -1,
		currentAction: my.START_STATE,
		curtok:        startTok,
		lastToken:     prev,
		currentKind:   startKind,
	}
	my.ensureCapacity(start, 16)

	live := []*glrConfig{start}
	accepts := make([]*glrAccept, 0)
	errorTok := startTok
	outerGuard := my.prs.GetNumStates()*64 +
		my.tokStream.GetStreamLength()*8 + 256

	for len(live) > 0 {
		if my.monitor != nil && my.monitor.IsCancelled() {
			return nil, nil
		}
		outerGuard--
		if outerGuard < 0 {
			return nil, fmt.Errorf(
				"cyclic/epsilon-loop grammar not supported by GLR v2")
		}

		next := make([]*glrConfig, 0)
		packed := make(map[string][]*glrConfig)
		for _, cfg := range live {
			if cfg.curtok > errorTok {
				errorTok = cfg.curtok
			}
			stepResults := make([]*glrConfig, 0)
			stepAccepts := make([]*glrAccept, 0)
			if err := my.stepConfig(cfg, &stepResults, &stepAccepts); err != nil {
				return nil, err
			}
			for _, candidate := range stepAccepts {
				if err := my.packAccept(&accepts, candidate); err != nil {
					return nil, err
				}
			}
			for _, result := range stepResults {
				key := result.key()
				bucket := packed[key]
				merged := false
				for _, existing := range bucket {
					if my.canPackParseStacks(existing, result) {
						if err := my.packParseStacks(existing, result); err != nil {
							return nil, err
						}
						merged = true
						break
					}
				}
				if !merged {
					packed[key] = append(bucket, result)
					next = append(next, result)
				}
			}
		}

		if len(accepts) > 0 && len(next) == 0 {
			break
		}
		live = next
		if len(live) == 0 && len(accepts) == 0 {
			return nil, NewBadParseException(errorTok)
		}
	}

	if len(accepts) == 0 {
		return nil, NewBadParseException(errorTok)
	}

	root := accepts[0].ast
	rootSymbol := accepts[0].grammarSymbol
	my.sppfRoot = accepts[0].sppf
	for i := 1; i < len(accepts); i++ {
		other := accepts[i]
		if other.grammarSymbol != rootSymbol {
			return nil, fmt.Errorf("GLR accepted distinct start symbols")
		}
		if my.sppfRoot == nil {
			my.sppfRoot = other.sppf
		}
		if !appendNextAst(root, other.ast, true) {
			return nil, fmt.Errorf("overlapping GLR accept forests")
		}
	}
	my.sppfSymbols = len(my.sppfNodes)
	if root == glrNullResult {
		return nil, nil
	}
	return root, nil
}

func (my *GLRParser) stepConfig(cfg *glrConfig, out *[]*glrConfig,
	accepts *[]*glrAccept) error {
	work := []*glrConfig{cfg.copy()}
	guard := my.prs.GetNumStates()*4 + 8

	for len(work) > 0 {
		guard--
		if guard < 0 {
			return fmt.Errorf(
				"cyclic/epsilon-loop grammar not supported by GLR v2")
		}
		last := len(work) - 1
		current := work[last]
		work = work[:last]

		my.ensureCapacity(current, current.stateStackTop+2)
		current.stateStackTop++
		top := current.stateStackTop
		current.stateStack[top] = current.currentAction
		current.locationStack[top] = current.curtok
		current.symbolStack[top] = 0
		current.sppfStack[top] = nil
		if top != my.parseStackRoot {
			current.parseStack[top] = nil
		}
		current.gssTip = my.gssPush(current.gssTip, current.currentAction,
			current.curtok, 0, nil, nil)

		act := my.tAction(current.currentAction, current.currentKind,
			current.curtok)
		candidates := []int{act}
		if act > my.ACCEPT_ACTION && act < my.ERROR_ACTION {
			candidates = my.expandConflict(act)
		}
		for _, candidate := range candidates {
			fork := current
			if len(candidates) != 1 {
				fork = current.copy()
			}
			if err := my.applyConcreteAction(fork, candidate,
				&work, out, accepts); err != nil {
				return err
			}
		}
	}
	return nil
}

func (my *GLRParser) applyConcreteAction(fork *glrConfig, candidate int,
	work *[]*glrConfig, out *[]*glrConfig, accepts *[]*glrAccept) error {
	switch {
	case candidate <= my.NUM_RULES:
		fork.stateStackTop--
		fork.gssTip = gssPop(fork.gssTip)
		return my.applyReduceClosure(fork, candidate, work)

	case candidate > my.ERROR_ACTION:
		top := fork.stateStackTop
		fork.symbolStack[top] = fork.currentKind
		term := my.terminalSppf(fork.currentKind, fork.curtok)
		fork.sppfStack[top] = term
		fork.gssTip = gssRelabel(fork.gssTip, fork.currentKind,
			fork.curtok, nil, term)
		fork.lastToken = fork.curtok
		fork.curtok = my.tokStream.GetNext(fork.curtok)
		fork.currentKind = my.tokStream.GetKind(fork.curtok)
		return my.applyReduceClosure(fork,
			candidate-my.ERROR_ACTION, work)

	case candidate < my.ACCEPT_ACTION:
		top := fork.stateStackTop
		fork.symbolStack[top] = fork.currentKind
		term := my.terminalSppf(fork.currentKind, fork.curtok)
		fork.sppfStack[top] = term
		fork.gssTip = gssRelabel(fork.gssTip, fork.currentKind,
			fork.curtok, nil, term)
		fork.lastToken = fork.curtok
		fork.curtok = my.tokStream.GetNext(fork.curtok)
		fork.currentKind = my.tokStream.GetKind(fork.curtok)
		fork.currentAction = candidate
		*out = append(*out, fork)

	case candidate == my.ACCEPT_ACTION:
		var root interface{}
		rootSymbol := 0
		var rootSppf *SppfNode
		if my.parseStackRoot < len(fork.parseStack) {
			root = fork.parseStack[my.parseStackRoot]
		}
		if my.parseStackRoot <= fork.stateStackTop {
			rootSymbol = fork.symbolStack[my.parseStackRoot]
		}
		if my.parseStackRoot < len(fork.sppfStack) {
			rootSppf = fork.sppfStack[my.parseStackRoot]
		}
		if root == nil {
			root = glrNullResult
		}
		*accepts = append(*accepts, &glrAccept{
			ast: root, grammarSymbol: rootSymbol, sppf: rootSppf,
		})
	}
	// candidate == ERROR_ACTION drops this fork.
	return nil
}

func (my *GLRParser) applyReduceClosure(fork *glrConfig, rule int,
	work *[]*glrConfig) error {
	action := rule
	for {
		rhs := my.prs.Rhs(action)
		if fork.stateStackTop-(rhs-1) < 0 {
			return fmt.Errorf("GLR reduce stack underflow")
		}

		children := make([]*SppfNode, rhs)
		for i := 0; i < rhs; i++ {
			children[i] =
				fork.sppfStack[fork.stateStackTop-rhs+1+i]
		}
		fork.stateStackTop -= rhs - 1
		if rhs > 0 {
			for i := 0; i < rhs-1; i++ {
				fork.gssTip = gssPop(fork.gssTip)
			}
		} else {
			my.ensureCapacity(fork, fork.stateStackTop+1)
			fork.gssTip = my.gssPush(fork.gssTip,
				fork.stateStack[fork.stateStackTop],
				fork.locationStack[fork.stateStackTop], 0, nil, nil)
		}

		reductionKey := my.reductionKey(action, fork.lastToken, rhs,
			fork.stateStackTop, fork.locationStack, fork.symbolStack,
			fork.parseStack)
		my.currentAction = action
		my.lastToken = fork.lastToken
		my.frameTop = fork.stateStackTop
		my.frameLocation = fork.locationStack
		my.frameParse = fork.parseStack

		my.takingActions = true
		my.ra.RuleAction(action)
		my.takingActions = false

		lhs := my.prs.Lhs(action)
		lhsSymbol := my.NT_OFFSET + lhs
		result := fork.parseStack[fork.stateStackTop]
		if ast, ok := result.(IAst); ok && ast != nil {
			canonical := my.familyCache[reductionKey]
			if canonical == nil {
				forestKey, packable := astForestKey(lhsSymbol, ast)
				if packable {
					canonical = my.forestCache[forestKey]
				}
				if canonical == nil {
					canonical = ast
					if packable {
						my.forestCache[forestKey] = canonical
					}
				} else if !sameIdentity(canonical, ast) &&
					!appendNextAst(canonical, ast, true) {
					return fmt.Errorf("cannot merge GLR production family")
				}
				my.familyCache[reductionKey] = canonical
			}
			fork.parseStack[fork.stateStackTop] = canonical
			result = canonical
		}

		leftExtent := fork.locationStack[fork.stateStackTop]
		rightExtent := fork.lastToken
		if ast, ok := result.(IAst); ok && ast != nil {
			left := ast.GetLeftIToken()
			right := ast.GetRightIToken()
			if left != nil && right != nil {
				leftExtent = left.GetTokenIndex()
				rightExtent = right.GetTokenIndex()
			}
		}
		symbolNode := my.sppfSymbol(lhsSymbol, leftExtent, rightExtent)
		my.addPacked(symbolNode, action, children, result)
		if _, ok := result.(IAst); ok {
			symbolNode.astForest = result
		}
		top := fork.stateStackTop
		fork.sppfStack[top] = symbolNode
		fork.symbolStack[top] = lhsSymbol
		fork.gssTip = gssRelabel(fork.gssTip, lhsSymbol,
			leftExtent, result, symbolNode)
		action = my.prs.NtAction(fork.stateStack[top], lhs)
		if action > my.NUM_RULES {
			break
		}
	}

	fork.currentAction = action
	*work = append(*work, fork)
	return nil
}

func (my *GLRParser) ensureCapacity(config *glrConfig, need int) {
	oldLength := len(config.stateStack)
	if need < oldLength {
		return
	}
	newLength := need + 8
	if oldLength+my.STACK_INCREMENT > newLength {
		newLength = oldLength + my.STACK_INCREMENT
	}
	config.stateStack = growInts(config.stateStack, newLength)
	config.symbolStack = growInts(config.symbolStack, newLength)
	config.locationStack = growInts(config.locationStack, newLength)
	config.parseStack = growInterfaces(config.parseStack, newLength)
	config.sppfStack = growSppf(config.sppfStack, newLength)
}

func (my *GLRParser) sppfSymbol(grammarSymbol int, leftExtent int,
	rightExtent int) *SppfNode {
	key := fmt.Sprintf("%d:%d:%d", grammarSymbol, leftExtent, rightExtent)
	node := my.sppfNodes[key]
	if node == nil {
		node = NewSppfNode(grammarSymbol, leftExtent, rightExtent)
		my.sppfNodes[key] = node
	}
	return node
}

func (my *GLRParser) terminalSppf(kind int, token int) *SppfNode {
	term := my.sppfSymbol(kind, token, token)
	if len(term.packs) == 0 {
		term.packs = append(term.packs,
			NewSppfPackedNode(-kind, nil, nil))
	}
	return term
}

func (my *GLRParser) addPacked(symbolNode *SppfNode, rule int,
	children []*SppfNode, semantic interface{}) {
	for _, packed := range symbolNode.packs {
		if packed.rule != rule || len(packed.children) != len(children) {
			continue
		}
		same := true
		for i := range children {
			if packed.children[i] != children[i] {
				same = false
				break
			}
		}
		if same {
			return
		}
	}
	copyChildren := append([]*SppfNode(nil), children...)
	symbolNode.packs = append(symbolNode.packs,
		NewSppfPackedNode(rule, copyChildren, semantic))
}

func (my *GLRParser) gssPush(tip *GssNode, state int, index int,
	symbol int, semantic interface{}, sppf *SppfNode) *GssNode {
	node := NewGssNode(state, index)
	predecessor := tip
	if predecessor == nil {
		predecessor = NewGssNode(-int(^uint(0)>>1)-1, -1)
	}
	node.edges = append(node.edges,
		NewGssEdge(predecessor, symbol, index, semantic, sppf))

	key := fmt.Sprintf("%d:%d", state, index)
	canonical := my.gssNodes[key]
	if canonical == nil {
		canonical = NewGssNode(state, index)
		my.gssNodes[key] = canonical
	}
	canonical.edges = append(canonical.edges,
		NewGssEdge(predecessor, symbol, index, semantic, sppf))
	return node
}

func gssPop(tip *GssNode) *GssNode {
	if tip == nil || len(tip.edges) == 0 {
		return nil
	}
	predecessor := tip.edges[0].predecessor
	if predecessor.state == -int(^uint(0)>>1)-1 {
		return nil
	}
	return predecessor
}

func gssRelabel(tip *GssNode, symbol int, location int,
	semantic interface{}, sppf *SppfNode) *GssNode {
	if tip == nil || len(tip.edges) == 0 {
		return tip
	}
	node := NewGssNode(tip.state, tip.index)
	node.edges = append(node.edges,
		NewGssEdge(tip.edges[0].predecessor, symbol, location,
			semantic, sppf))
	return node
}

func (my *GLRParser) packAccept(accepts *[]*glrAccept,
	candidate *glrAccept) error {
	if candidate.ast == glrNullResult {
		for _, existing := range *accepts {
			if existing.ast == glrNullResult {
				return nil
			}
		}
		*accepts = append(*accepts, candidate)
		return nil
	}
	ast, ok := candidate.ast.(IAst)
	if !ok || ast == nil {
		return nil
	}
	for _, existing := range *accepts {
		other, ok := existing.ast.(IAst)
		if !ok || other == nil {
			continue
		}
		if existing.grammarSymbol == candidate.grammarSymbol &&
			sameSpan(other, ast) && appendNextAst(other, ast, true) {
			return nil
		}
	}
	*accepts = append(*accepts, candidate)
	return nil
}

func (my *GLRParser) canPackParseStacks(existing *glrConfig,
	incoming *glrConfig) bool {
	if existing.stateStackTop != incoming.stateStackTop {
		return false
	}
	for i := 0; i <= existing.stateStackTop; i++ {
		a, b := existing.parseStack[i], incoming.parseStack[i]
		if sameIdentity(a, b) {
			continue
		}
		astA, okA := a.(IAst)
		astB, okB := b.(IAst)
		if !okA || !okB || astA == nil || astB == nil ||
			!sameSpan(astA, astB) ||
			!appendNextAst(astA, astB, false) {
			return false
		}
	}
	return true
}

func (my *GLRParser) packParseStacks(existing *glrConfig,
	incoming *glrConfig) error {
	for i := 0; i <= existing.stateStackTop; i++ {
		a, b := existing.parseStack[i], incoming.parseStack[i]
		if sameIdentity(a, b) || a == nil || b == nil {
			continue
		}
		astA, okA := a.(IAst)
		astB, okB := b.(IAst)
		if !okA || !okB || !appendNextAst(astA, astB, false) {
			return fmt.Errorf("overlapping GLR semantic forests")
		}
	}

	for i := 0; i <= existing.stateStackTop; i++ {
		a, b := existing.parseStack[i], incoming.parseStack[i]
		switch {
		case a == nil:
			existing.parseStack[i] = b
		case b == nil || sameIdentity(a, b):
		default:
			if !appendNextAst(a, b, true) {
				return fmt.Errorf("overlapping GLR semantic forests")
			}
		}

		left, right := existing.sppfStack[i], incoming.sppfStack[i]
		if left == nil {
			existing.sppfStack[i] = right
		} else if right != nil && left != right &&
			left.grammarSymbol == right.grammarSymbol &&
			left.leftExtent == right.leftExtent &&
			left.rightExtent == right.rightExtent {
			for _, packed := range right.packs {
				my.addPacked(left, packed.rule, packed.children,
					packed.semantic)
			}
			if _, ok := existing.parseStack[i].(IAst); ok {
				left.astForest = existing.parseStack[i]
			}
		}
	}
	if incoming.gssTip != nil {
		existing.gssTip = incoming.gssTip
	}
	return nil
}

func (my *GLRParser) reductionKey(rule int, lastToken int, rhs int,
	frameTop int, locations []int, symbols []int,
	semantics []interface{}) string {
	var out strings.Builder
	out.WriteString(strconv.Itoa(rule))
	out.WriteByte(':')
	out.WriteString(strconv.Itoa(lastToken))
	for i := 0; i < rhs; i++ {
		index := frameTop + i
		out.WriteByte(':')
		out.WriteString(strconv.Itoa(locations[index]))
		out.WriteByte(':')
		out.WriteString(strconv.Itoa(symbols[index]))
		out.WriteByte(':')
		out.WriteString(strconv.FormatUint(
			uint64(identityOf(semantics[index])), 10))
	}
	return out.String()
}

type glrAccept struct {
	ast           interface{}
	grammarSymbol int
	sppf          *SppfNode
}

type glrConfig struct {
	stateStack    []int
	symbolStack   []int
	parseStack    []interface{}
	locationStack []int
	sppfStack     []*SppfNode
	gssTip        *GssNode
	stateStackTop int
	currentAction int
	curtok        int
	lastToken     int
	currentKind   int
}

func (config *glrConfig) copy() *glrConfig {
	return &glrConfig{
		stateStack:    append([]int(nil), config.stateStack...),
		symbolStack:   append([]int(nil), config.symbolStack...),
		parseStack:    append([]interface{}(nil), config.parseStack...),
		locationStack: append([]int(nil), config.locationStack...),
		sppfStack:     append([]*SppfNode(nil), config.sppfStack...),
		gssTip:        config.gssTip,
		stateStackTop: config.stateStackTop,
		currentAction: config.currentAction,
		curtok:        config.curtok,
		lastToken:     config.lastToken,
		currentKind:   config.currentKind,
	}
}

func (config *glrConfig) key() string {
	var out strings.Builder
	fmt.Fprintf(&out, "%d:%d:%d:%d:%d", config.curtok,
		config.currentKind, config.lastToken, config.currentAction,
		config.stateStackTop)
	for i := 0; i <= config.stateStackTop; i++ {
		fmt.Fprintf(&out, ":%d:%d:%d", config.stateStack[i],
			config.locationStack[i], config.symbolStack[i])
	}
	return out.String()
}

func sameSpan(a IAst, b IAst) bool {
	leftA, rightA := a.GetLeftIToken(), a.GetRightIToken()
	leftB, rightB := b.GetLeftIToken(), b.GetRightIToken()
	if leftA == nil || rightA == nil || leftB == nil || rightB == nil {
		return false
	}
	return sameIdentity(leftA.GetILexStream(), leftB.GetILexStream()) &&
		sameIdentity(rightA.GetILexStream(), rightB.GetILexStream()) &&
		leftA.GetTokenIndex() == leftB.GetTokenIndex() &&
		rightA.GetTokenIndex() == rightB.GetTokenIndex()
}

func appendNextAst(root interface{}, alternative interface{}, commit bool) bool {
	current, okCurrent := root.(IAst)
	incoming, okIncoming := alternative.(IAst)
	if !okCurrent || !okIncoming || current == nil || incoming == nil {
		return false
	}
	if sameIdentity(current, incoming) {
		return true
	}

	seen := make(map[uintptr]bool)
	var tail IAst
	for node := current; node != nil; node = node.GetNextAst() {
		id := identityOf(node)
		if id == 0 || seen[id] {
			return false
		}
		seen[id] = true
		tail = node
	}

	incomingSeen := make(map[uintptr]bool)
	for node := incoming; node != nil; {
		id := identityOf(node)
		if id == 0 || incomingSeen[id] {
			return false
		}
		incomingSeen[id] = true
		if seen[id] {
			node = node.GetNextAst()
			continue
		}
		for next := node.GetNextAst(); next != nil; next = next.GetNextAst() {
			nextID := identityOf(next)
			if nextID == 0 || incomingSeen[nextID] || seen[nextID] {
				return false
			}
			incomingSeen[nextID] = true
		}
		if commit {
			setter, ok := tail.(interface{ SetNextAst(IAst) })
			if !ok {
				return false
			}
			setter.SetNextAst(node)
		}
		return true
	}
	return true
}

func astForestKey(grammarSymbol int, ast IAst) (string, bool) {
	left, right := ast.GetLeftIToken(), ast.GetRightIToken()
	if left == nil || right == nil {
		return "", false
	}
	return fmt.Sprintf("%d:%d:%d:%d", grammarSymbol,
		identityOf(left.GetILexStream()), left.GetTokenIndex(),
		right.GetTokenIndex()), true
}

func sameIdentity(a interface{}, b interface{}) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return reflect.TypeOf(a) == reflect.TypeOf(b) &&
		identityOf(a) == identityOf(b)
}

func identityOf(value interface{}) uintptr {
	if value == nil {
		return 0
	}
	ref := reflect.ValueOf(value)
	switch ref.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr,
		reflect.Slice, reflect.UnsafePointer:
		if ref.IsNil() {
			return 0
		}
		return ref.Pointer()
	default:
		// GLR semantic values are normally pointers. This stable-enough fallback
		// distinguishes scalar values without making them map keys directly.
		text := fmt.Sprintf("%T:%v", value, value)
		var hash uintptr = 1469598103934665603
		for i := 0; i < len(text); i++ {
			hash ^= uintptr(text[i])
			hash *= 1099511628211
		}
		return hash
	}
}

func growInts(source []int, length int) []int {
	out := make([]int, length)
	copy(out, source)
	return out
}

func growInterfaces(source []interface{}, length int) []interface{} {
	out := make([]interface{}, length)
	copy(out, source)
	return out
}

func growSppf(source []*SppfNode, length int) []*SppfNode {
	out := make([]*SppfNode, length)
	copy(out, source)
	return out
}
