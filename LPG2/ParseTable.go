package lpg2

type ParseTable interface {
	baseCheck(index int) int

	rhs(index int) int

	baseAction(index int) int

	lhs(index int) int

	termCheck(index int) int

	termAction(index int) int

	asb(index int) int

	asr(index int) int

	nasb(index int) int

	nasr(index int) int

	terminalIndex(index int) int

	nonterminalIndex(index int) int

	scopePrefix(index int) int

	scopeSuffix(index int) int

	scopeLhs(index int) int

	scopeLa(index int) int

	scopeStateSet(index int) int

	scopeRhs(index int) int

	scopeState(index int) int

	inSymb(index int) int

	name(index int) string

	originalState(state int) int

	asi(state int) int

	nasi(state int) int

	inSymbol(state int) int

	ntAction(state int, sym int) int

	tAction(act int, sym int) int

	lookAhead(act int, sym int) int

	getErrorSymbol() int

	getScopeUbound() int

	getScopeSize() int

	getMaxNameLength() int

	getNumStates() int

	getNtOffset() int

	getLaStateOffset() int

	getMaxLa() int

	getNumRules() int

	getNumNonterminals() int

	getNumSymbols() int

	getStartState() int

	getStartSymbol() int

	getEoftSymbol() int

	getEoltSymbol() int

	getAcceptAction() int

	getErrorAction() int

	isNullable(symbol int) bool

	isValidForParser() bool

	getBacktrack() bool
}
