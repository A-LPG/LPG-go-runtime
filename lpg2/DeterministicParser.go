package lpg2

type DeterministicParser struct {
	*Stacks
	taking_actions bool
	markerKind     int

	monitor         Monitor
	START_STATE     int
	NUM_RULES       int
	NT_OFFSET       int
	LA_STATE_OFFSET int
	EOFT_SYMBOL     int
	ACCEPT_ACTION   int
	ERROR_ACTION    int
	ERROR_SYMBOL    int

	lastToken     int
	currentAction int
	action        *IntTuple

	tokStream TokenStream
	prs       ParseTable
	ra        RuleAction
}

func NewDeterministicParser(tokStream TokenStream, prs ParseTable, ra RuleAction, monitor Monitor) (*DeterministicParser,error) {

	a := new(DeterministicParser)
	a.Stacks = NewStacks()

	err := a.Reset(tokStream, prs, ra, monitor)
	if err != nil {
		return nil,err
	}

	return a,nil
}

//
// keep looking ahead until we compute a valid action
//
func (my *DeterministicParser) Lookahead(act int, token int) int {
	act = my.prs.LookAhead(act-my.LA_STATE_OFFSET, my.tokStream.GetKind(token))
	if act > my.LA_STATE_OFFSET {
		return my.Lookahead(act, my.tokStream.GetNext(token))
	} else {
		return act
	}

}

//
// Compute the next action defined on act and sym. If my
// action requires more Lookahead, these Lookahead symbols
// are in the token stream beginning at the next token that
// is yielded by Peek().
//
func (my *DeterministicParser) TAction1(act int, sym int) int {
	act = my.prs.TAction(act, sym)
	if act > my.LA_STATE_OFFSET {
		return my.Lookahead(act, my.tokStream.Peek())
	} else {
		return act
	}

}

//
// Compute the next action defined on act and the next k tokens
// whose types are stored in the array sym starting at location
// index. The array sym is a circular buffer. If we reach the last
// element of sym and we need more Lookahead, we proceed to the
// first element.
//
// assert(sym.length == prs.GetMaxLa())
//
func (my *DeterministicParser) TAction(act int, sym []int, index int) int {

	act = my.prs.TAction(act, sym[index])
	for act > my.LA_STATE_OFFSET {
		index = (index + 1) % len(sym)
		act = my.prs.LookAhead(act-my.LA_STATE_OFFSET, sym[index])
	}
	return act
}

//
// Process reductions and continue...
//
func (my *DeterministicParser) ProcessReductions() {
	for {
		my.stateStackTop -= (my.prs.Rhs(my.currentAction) - 1)
		my.ra.RuleAction(my.currentAction)
		my.currentAction = my.prs.NtAction( my.stateStack[my.stateStackTop],
												my.prs.Lhs(my.currentAction))
		if my.currentAction <= my.NUM_RULES {
			continue
		} else {
			break
		}
	}
	return
}

//
// The following functions can be invoked only when the parser is
// processing actions. Thus, they can be invoked when the parser
// was entered via the main entry point (parse()). When using
// the incremental parser (via the entry point parse(int [], int)),
// an Exception is thrown if any of these functions is invoked?
// However, note that when ParseActions() is invoked after successfully
// parsing an input with the incremental parser, then they can be invoked.
//
func (my *DeterministicParser) GetCurrentRule() int{
	if my.taking_actions {
		return my.currentAction
	}
	return -1
}
func (my *DeterministicParser) GetFirstToken() int {
	if my.taking_actions {
		return my.GetToken(1)
	}
	return -1
}
func (my *DeterministicParser) GetFirstTokenAt(i int) int {

	if my.taking_actions {
		return my.GetToken(i)
	}
	return -1
}
func (my *DeterministicParser) GetLastToken() int{
	if my.taking_actions {
		return my.lastToken
	}
	return -1
}
func (my *DeterministicParser) GetLastTokenAt(i int) int {

	if my.taking_actions {
		if i >= my.prs.Rhs(my.currentAction) {
			return my.lastToken
		} else {
			return my.tokStream.GetPrevious(my.GetToken(i + 1))
		}
	}
	return -1
}
func (my *DeterministicParser) SetMonitor(monitor Monitor) {
	my.monitor = monitor
}
func (my *DeterministicParser) Reset1() {
	my.taking_actions = false
	my.markerKind = 0
	if my.action != nil {
		my.action.Reset()
	}
}
func (my *DeterministicParser) Reset2(tokStream TokenStream, monitor Monitor) {
	my.monitor = monitor
	my.tokStream = tokStream
	my.Reset1()
}

func (my *DeterministicParser) Reset(tokStream TokenStream, prs ParseTable, ra RuleAction, monitor Monitor) error {
	if nil != ra {
		my.ra = ra
	}
	if nil != prs {
		my.prs = prs

		my.START_STATE = prs.GetStartState()
		my.NUM_RULES = prs.GetNumRules()
		my.NT_OFFSET = prs.GetNtOffset()
		my.LA_STATE_OFFSET = prs.GetLaStateOffset()
		my.EOFT_SYMBOL = prs.GetEoftSymbol()
		my.ERROR_SYMBOL = prs.GetErrorSymbol()
		my.ACCEPT_ACTION = prs.GetAcceptAction()
		my.ERROR_ACTION = prs.GetErrorAction()
		if !prs.IsValidForParser() {
			return NewBadParseSymFileException("")
		}
		if prs.GetBacktrack() {
			return NewNotDeterministicParseTableException("")
		}
	}
	if nil == tokStream {
		my.Reset1()
		return nil
	}
	my.Reset2(tokStream, monitor)
	return nil
}

func (my *DeterministicParser) ParseEntry(marker_kind int) (interface{}, error) {
	//
	// Indicate that we are running the regular parser and that it's
	// ok to use the utility functions to query the parser.
	//
	my.taking_actions = true
	//
	// Reset the token stream and Get the first token.
	//
	my.tokStream.Reset()
	my.lastToken = my.tokStream.GetPrevious(my.tokStream.Peek())
	var curtok int
	var current_kind int
	if marker_kind == 0 {
		curtok = my.tokStream.GetToken()
		current_kind = my.tokStream.GetKind(curtok)
	} else {
		curtok = my.lastToken
		current_kind = marker_kind
	}
	//
	// Start parsing.
	//
	my.ReallocateStacks() // make initial allocation
	my.stateStackTop = -1
	my.currentAction = my.START_STATE

	processTerminals:
	for {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if my.monitor != nil && my.monitor.IsCancelled() {
			my.taking_actions = false // indicate that we are done
			return nil, nil
		}

		my.stateStackTop += 1
		if my.stateStackTop >= len(my.stateStack) {
			my.ReallocateStacks()
		}

		my.stateStack[my.stateStackTop] = my.currentAction

		my.locationStack[my.stateStackTop] = curtok

		my.currentAction = my.TAction1(my.currentAction, current_kind)

		if my.currentAction <= my.NUM_RULES {
			my.stateStackTop-- // make reduction look like a shift-reduce
			my.ProcessReductions()
		} else {
			if my.currentAction > my.ERROR_ACTION {
				my.lastToken = curtok
				curtok = my.tokStream.GetToken()
				current_kind = my.tokStream.GetKind(curtok)
				my.currentAction -= my.ERROR_ACTION
				my.ProcessReductions()
			} else {
				if my.currentAction < my.ACCEPT_ACTION {
					my.lastToken = curtok
					curtok = my.tokStream.GetToken()
					current_kind = my.tokStream.GetKind(curtok)
				} else {
					break processTerminals
				}
			}
		}
	}

	my.taking_actions = false // indicate that we are done

	if my.currentAction == my.ERROR_ACTION {
		return nil, NewBadParseException(curtok)
	}

	if marker_kind == 0 {
		return my.parseStack[0], nil
	} else {
		return my.parseStack[1], nil
	}
}

//
// This method is invoked when using the parser in an incremental mode
// using the entry point parse(int [], int).
//
func (my *DeterministicParser) ResetParser() {
	my.ResetParserEntry(0)
}

//
// This method is invoked when using the parser in an incremental mode
// using the entry point parse(int [], int).
//
func (my *DeterministicParser) ResetParserEntry(marker_kind int) {
	my.markerKind = marker_kind
	if my.stateStack == nil || len(my.stateStack) == 0 {
		my.ReallocateStacks() // make initial allocation
	}
	my.stateStackTop = 0
	my.stateStack[my.stateStackTop] = my.START_STATE
	if my.action== nil {
		my.action = NewIntTupleWithEstimate(1 << 20)
	} else {
		my.action.Reset()
	}
	//
	// Indicate that we are going to run the incremental parser and that
	// it's forbidden to use the utility functions to query the parser.
	//
	my.taking_actions = false
	if marker_kind != 0 {
		var sym = []int{my.markerKind}
		my.Parse(sym, 0)
	}
}

//
// Find a state in the state stack that has a valid action on ERROR token
//
func (my *DeterministicParser) RecoverableState(state int) bool {
	var k = my.prs.Asi(state)
	for ; my.prs.Asr(k) != 0; k++ {
		if my.prs.Asr(k) == my.ERROR_SYMBOL {
			return true
		}
	}
	return false
}

//
// Reset the parser at a point where it can legally process
// the error token. If we can't do that, Reset it to the beginning.
//
func (my *DeterministicParser) ErrorReset() {
	var gate int
	if my.markerKind == 0 {
		gate = 0
	} else {
		gate = 1
	}
	for ; my.stateStackTop >= gate; my.stateStackTop-- {
		if my.RecoverableState(my.stateStack[my.stateStackTop]) {
			break
		}
	}
	if my.stateStackTop < gate {
		my.ResetParserEntry(my.markerKind)
	}
	return
}

//
// This is an incremental LALR(k) parser that takes as argument
// the next k tokens in the input. If these k tokens are valid for
// the current configuration, it advances past the first of the k
// tokens and returns either
//
//    . the last transition induced by that token
//    . the Accept action
//
// If the tokens are not valid, the initial configuration remains
// unchanged and the Error action is returned.
//
// Note that it is the user's responsibility to start the parser in a
// proper configuration by initially invoking the method ResetParser
// prior to invoking my function.
//
func (my *DeterministicParser) Parse(sym []int, index int) (int,error) {

	// assert(sym.length == prs.GetMaxLa())

	//
	// First, we save the current length of the action tuple, in
	// case an error is encountered and we need to restore the
	// original configuration.
	//
	// Next, we declara and initialize the variable pos which will
	// be used to indicate the highest useful position in stateStack
	// as we are simulating the actions induced by the next k input
	// terminals in sym.
	//
	// The location stack will be used here as a temporary stack
	// to simulate these actions. We initialize its first useful
	// offSet here.
	//
	var save_action_length int = my.action.Size()
	var pos int = my.stateStackTop
	var location_top int = my.stateStackTop - 1

	//
	// When a reduce action is encountered, we compute all REDUCE
	// and associated goto actions induced by the current token.
	// Eventually, a SHIFT, SHIFT-REDUCE, ACCEPT or ERROR action is
	// computed...
	//
	for my.currentAction = my.TAction(my.stateStack[my.stateStackTop], sym, index);
		my.currentAction <= my.NUM_RULES;
		my.currentAction = my.TAction(my.currentAction, sym, index) {
		my.action.Add(my.currentAction)
		for {
			location_top -= (my.prs.Rhs(my.currentAction) - 1)

			var state int
			if location_top > pos {
				state = my.locationStack[location_top]
			} else {
				state = my.stateStack[location_top]
			}

			my.currentAction = my.prs.NtAction(state, my.prs.Lhs(my.currentAction))
			if my.currentAction <= my.NUM_RULES {
				continue
			} else {
				break
			}
		}

		//
		// ... Update the maximum useful position of the
		// stateSTACK, push goto state into stack, and
		// continue by compute next action on current symbol
		// and reentering the loop...
		//

		if !(pos < location_top) {
			pos = location_top
		}
		if location_top+1 >= len(my.locationStack) {
			my.ReallocateStacks()
		}
		my.locationStack[location_top+1] = my.currentAction

	}
	//
	// At my point, we have a shift, shift-reduce, accept or error
	// action. stateSTACK contains the configuration of the state stack
	// prior to executing any action on the currenttoken. locationStack
	// contains the configuration of the state stack after executing all
	// reduce actions induced by the current token. The variable pos
	// indicates the highest position in the stateSTACK that is still
	// useful after the reductions are executed.
	//
	if my.currentAction > my.ERROR_ACTION || // SHIFT-REDUCE action ?
		my.currentAction < my.ACCEPT_ACTION { // SHIFT action ?

		my.action.Add(my.currentAction)
		//
		// If no error was detected, update the state stack with
		// the info that was temporarily computed in the locationStack.
		//
		my.stateStackTop = location_top + 1
		var i int = pos + 1
		for ; i <= my.stateStackTop; i++ {
			my.stateStack[i] = my.locationStack[i]
		}

		//
		// If we have a shift-reduce, process it as well as
		// the goto-reduce actions that follow it.
		//

		if my.currentAction > my.ERROR_ACTION {
			my.currentAction -= my.ERROR_ACTION
			for {
				my.stateStackTop -= my.prs.Rhs(my.currentAction) - 1
				my.currentAction = my.prs.NtAction( my.stateStack[my.stateStackTop],
														my.prs.Lhs(my.currentAction))
				if my.currentAction <= my.NUM_RULES {
					continue
				} else {
					break
				}
			}
		}
		//
		// Process the final transition - either a shift action of
		// if we started out with a shift-reduce, the final GOTO
		// action that follows it.
		//
		my.stateStackTop += 1
		if my.stateStackTop >= len(my.stateStack) {
			my.ReallocateStacks()
		}
		my.stateStack[my.stateStackTop] = my.currentAction
	} else {
		if my.currentAction == my.ERROR_ACTION {
			my.action.ResetTo(save_action_length) // restore original action state.
		}
	}
	return my.currentAction,nil
}

//
// Now do the final parse of the input based on the actions in
// the list "action" and the sequence of tokens in the token stream.
//
func (my *DeterministicParser) ParseActions() interface{} {
	//
	// Indicate that we are processing actions now (for the incremental
	// parser) and that it's ok to use the utility functions to query the
	// parser.
	//
	my.taking_actions = true
	my.tokStream.Reset()
	my.lastToken = my.tokStream.GetPrevious(my.tokStream.Peek())
	var curtok int
	if my.markerKind == 0 {
		curtok = my.tokStream.GetToken()
	} else {
		curtok = my.lastToken
	}

	//
	// Reparse the input...
	//
	my.stateStackTop = -1
	my.currentAction = my.START_STATE
	var i int = 0
	for ; i < my.action.Size(); i++ {
		//
		// if the parser needs to stop processing, it may do so here.
		//
		if my.monitor != nil && my.monitor.IsCancelled() {
			my.taking_actions = false
			return nil
		}

		my.stateStackTop += 1
		my.stateStack[my.stateStackTop] = my.currentAction
		my.locationStack[my.stateStackTop] = curtok

		my.currentAction = my.action.Get(i)
		if my.currentAction <= my.NUM_RULES { // a reduce action?

			my.stateStackTop-- // turn reduction intoshift-reduction
			my.ProcessReductions()
		} else { // a shift or shift-reduce action

			my.lastToken = curtok
			curtok = my.tokStream.GetToken()
			if my.currentAction > my.ERROR_ACTION {
				my.currentAction -= my.ERROR_ACTION
				my.ProcessReductions()
			}
		}
	}

	my.taking_actions = false // indicate that we are done.
	my.action = nil
	if my.markerKind == 0 {
		return my.parseStack[0]
	} else {
		return my.parseStack[1]
	}
}
