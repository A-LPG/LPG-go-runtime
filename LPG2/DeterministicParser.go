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

func NewDeterministicParser(tokStream TokenStream, prs ParseTable, ra RuleAction, monitor Monitor) *DeterministicParser {

	a := new(DeterministicParser)
	a.Stacks = NewStacks()

	err := a.reset(tokStream, prs, ra, monitor)
	if err != nil {
		return nil
	}

	return a
}

//
// keep looking ahead until we compute a valid action
//
func (self *DeterministicParser) lookahead(act int, token int) int {
	act = self.prs.lookAhead(act-self.LA_STATE_OFFSET, self.tokStream.getKind(token))
	if act > self.LA_STATE_OFFSET {
		return self.lookahead(act, self.tokStream.getNext(token))
	} else {
		return act
	}

}

//
// Compute the next action defined on act and sym. If self
// action requires more lookahead, these lookahead symbols
// are in the token stream beginning at the next token that
// is yielded by peek().
//
func (self *DeterministicParser) tAction1(act int, sym int) int {
	act = self.prs.tAction(act, sym)
	if act > self.LA_STATE_OFFSET {
		return self.lookahead(act, self.tokStream.peek())
	} else {
		return act
	}

}

//
// Compute the next action defined on act and the next k tokens
// whose types are stored in the array sym starting at location
// index. The array sym is a circular buffer. If we reach the last
// element of sym and we need more lookahead, we proceed to the
// first element.
//
// assert(sym.length == prs.getMaxLa())
//
func (self *DeterministicParser) tAction(act int, sym []int, index int) int {

	act = self.prs.tAction(act, sym[index])
	for act > self.LA_STATE_OFFSET {
		index = (index + 1) % len(sym)
		act = self.prs.lookAhead(act-self.LA_STATE_OFFSET, sym[index])
	}
	return act
}

//
// Process reductions and continue...
//
func (self *DeterministicParser) processReductions() {
	for {
		self.stateStackTop -= (self.prs.rhs(self.currentAction) - 1)
		self.ra.ruleAction(self.currentAction)
		self.currentAction = self.prs.ntAction( self.stateStack[self.stateStackTop],
												self.prs.lhs(self.currentAction))
		if self.currentAction <= self.NUM_RULES {
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
// However, note that when parseActions() is invoked after successfully
// parsing an input with the incremental parser, then they can be invoked.
//
func (self *DeterministicParser) getCurrentRule() (int, error) {
	if self.taking_actions {
		return self.currentAction, nil
	}
	return -1, NewUnavailableParserInformationException("")
}
func (self *DeterministicParser) getFirstToken() (int, error) {
	if self.taking_actions {
		return self.getToken(1), nil
	}
	return -1, NewUnavailableParserInformationException("")
}
func (self *DeterministicParser) getFirstTokenAt(i int) (int, error) {

	if self.taking_actions {
		return self.getToken(i), nil
	}
	return -1, NewUnavailableParserInformationException("")
}
func (self *DeterministicParser) getLastToken() (int, error) {
	if self.taking_actions {
		return self.lastToken, nil
	}
	return -1, NewUnavailableParserInformationException("")
}
func (self *DeterministicParser) getLastTokenAt(i int) (int, error) {

	if self.taking_actions {
		if i >= self.prs.rhs(self.currentAction) {
			return self.lastToken, nil
		} else {
			return self.tokStream.getPrevious(self.getToken(i + 1)), nil
		}
	}
	return -1, NewUnavailableParserInformationException("")
}
func (self *DeterministicParser) setMonitor(monitor Monitor) {
	self.monitor = monitor
}
func (self *DeterministicParser) reset1() {
	self.taking_actions = false
	self.markerKind = 0
	if self.action == nil {
		self.action.reset()
	}
}
func (self *DeterministicParser) reset2(tokStream TokenStream, monitor Monitor) {
	self.monitor = monitor
	self.tokStream = tokStream
	self.reset1()
}

func (self *DeterministicParser) reset(tokStream TokenStream, prs ParseTable, ra RuleAction, monitor Monitor) error {
	if nil != ra {
		self.ra = ra
	}
	if nil != prs {
		self.prs = prs

		self.START_STATE = prs.getStartState()
		self.NUM_RULES = prs.getNumRules()
		self.NT_OFFSET = prs.getNtOffset()
		self.LA_STATE_OFFSET = prs.getLaStateOffset()
		self.EOFT_SYMBOL = prs.getEoftSymbol()
		self.ERROR_SYMBOL = prs.getErrorSymbol()
		self.ACCEPT_ACTION = prs.getAcceptAction()
		self.ERROR_ACTION = prs.getErrorAction()
		if !prs.isValidForParser() {
			return NewBadParseSymFileException("")
		}
		if prs.getBacktrack() {
			return NewNotDeterministicParseTableException("")
		}
	}
	if nil == tokStream {
		self.reset1()
		return nil
	}
	self.reset2(tokStream, monitor)
	return nil
}

func (self *DeterministicParser) parseEntry(marker_kind int) (interface{}, error) {
	//
	// Indicate that we are running the regular parser and that it's
	// ok to use the utility functions to query the parser.
	//
	self.taking_actions = true
	//
	// Reset the token stream and get the first token.
	//
	self.tokStream.reset()
	self.lastToken = self.tokStream.getPrevious(self.tokStream.peek())
	var curtok int
	var current_kind int
	if marker_kind == 0 {
		curtok = self.tokStream.getToken()
		current_kind = self.tokStream.getKind(curtok)
	} else {
		curtok = self.lastToken
		current_kind = marker_kind
	}
	//
	// Start parsing.
	//
	self.reallocateStacks() // make initial allocation
	self.stateStackTop = -1
	self.currentAction = self.START_STATE

	processTerminals:
	for {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if self.monitor != nil && self.monitor.isCancelled() {
			self.taking_actions = false // indicate that we are done
			return nil, nil
		}

		self.stateStackTop += 1
		if self.stateStackTop >= len(self.stateStack) {
			self.reallocateStacks()
		}

		self.stateStack[self.stateStackTop] = self.currentAction

		self.locationStack[self.stateStackTop] = curtok

		self.currentAction = self.tAction1(self.currentAction, current_kind)

		if self.currentAction <= self.NUM_RULES {
			self.stateStackTop-- // make reduction look like a shift-reduce
			self.processReductions()
		} else {
			if self.currentAction > self.ERROR_ACTION {
				self.lastToken = curtok
				curtok = self.tokStream.getToken()
				current_kind = self.tokStream.getKind(curtok)
				self.currentAction -= self.ERROR_ACTION
				self.processReductions()
			} else {
				if self.currentAction < self.ACCEPT_ACTION {
					self.lastToken = curtok
					curtok = self.tokStream.getToken()
					current_kind = self.tokStream.getKind(curtok)
				} else {
					break processTerminals
				}
			}
		}
	}

	self.taking_actions = false // indicate that we are done

	if self.currentAction == self.ERROR_ACTION {
		return nil, NewBadParseException(curtok)
	}

	if marker_kind == 0 {
		return self.parseStack[0], nil
	} else {
		return self.parseStack[1], nil
	}
}

//
// This method is invoked when using the parser in an incremental mode
// using the entry point parse(int [], int).
//
func (self *DeterministicParser) resetParser() {
	self.resetParserEntry(0)
}

//
// This method is invoked when using the parser in an incremental mode
// using the entry point parse(int [], int).
//
func (self *DeterministicParser) resetParserEntry(marker_kind int) {
	self.markerKind = marker_kind
	if self.stateStack == nil || len(self.stateStack) == 0 {
		self.reallocateStacks() // make initial allocation
	}
	self.stateStackTop = 0
	self.stateStack[self.stateStackTop] = self.START_STATE
	if self.action== nil {
		self.action = NewIntTupleWithEstimate(1 << 20)
	} else {
		self.action.reset()
	}
	//
	// Indicate that we are going to run the incremental parser and that
	// it's forbidden to use the utility functions to query the parser.
	//
	self.taking_actions = false
	if marker_kind != 0 {
		var sym = []int{self.markerKind}
		self.parse(sym, 0)
	}
}

//
// Find a state in the state stack that has a valid action on ERROR token
//
func (self *DeterministicParser) recoverableState(state int) bool {
	var k = self.prs.asi(state)
	for ; self.prs.asr(k) != 0; k++ {
		if self.prs.asr(k) == self.ERROR_SYMBOL {
			return true
		}
	}
	return false
}

//
// Reset the parser at a point where it can legally process
// the error token. If we can't do that, reset it to the beginning.
//
func (self *DeterministicParser) errorReset() {
	var gate int
	if self.markerKind == 0 {
		gate = 0
	} else {
		gate = 1
	}
	for ; self.stateStackTop >= gate; self.stateStackTop-- {
		if self.recoverableState(self.stateStack[self.stateStackTop]) {
			break
		}
	}
	if self.stateStackTop < gate {
		self.resetParserEntry(self.markerKind)
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
// proper configuration by initially invoking the method resetParser
// prior to invoking self function.
//
func (self *DeterministicParser) parse(sym []int, index int) int {

	// assert(sym.length == prs.getMaxLa())

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
	// offset here.
	//
	var save_action_length int = self.action.size()
	var pos int = self.stateStackTop
	var location_top int = self.stateStackTop - 1

	//
	// When a reduce action is encountered, we compute all REDUCE
	// and associated goto actions induced by the current token.
	// Eventually, a SHIFT, SHIFT-REDUCE, ACCEPT or ERROR action is
	// computed...
	//
	for self.currentAction = self.tAction(self.stateStack[self.stateStackTop], sym, index);
		self.currentAction <= self.NUM_RULES;
		self.currentAction = self.tAction(self.currentAction, sym, index) {
		self.action.add(self.currentAction)
		for {
			location_top -= (self.prs.rhs(self.currentAction) - 1)

			var state int
			if location_top > pos {
				state = self.locationStack[location_top]
			} else {
				state = self.stateStack[location_top]
			}

			self.currentAction = self.prs.ntAction(state, self.prs.lhs(self.currentAction))
			if self.currentAction <= self.NUM_RULES {
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
		if location_top+1 >= len(self.locationStack) {
			self.reallocateStacks()
		}
		self.locationStack[location_top+1] = self.currentAction

	}
	//
	// At self point, we have a shift, shift-reduce, accept or error
	// action. stateSTACK contains the configuration of the state stack
	// prior to executing any action on the currenttoken. locationStack
	// contains the configuration of the state stack after executing all
	// reduce actions induced by the current token. The variable pos
	// indicates the highest position in the stateSTACK that is still
	// useful after the reductions are executed.
	//
	if self.currentAction > self.ERROR_ACTION || // SHIFT-REDUCE action ?
		self.currentAction < self.ACCEPT_ACTION { // SHIFT action ?

		self.action.add(self.currentAction)
		//
		// If no error was detected, update the state stack with
		// the info that was temporarily computed in the locationStack.
		//
		self.stateStackTop = location_top + 1
		var i int = pos + 1
		for ; i <= self.stateStackTop; i++ {
			self.stateStack[i] = self.locationStack[i]
		}

		//
		// If we have a shift-reduce, process it as well as
		// the goto-reduce actions that follow it.
		//

		if self.currentAction > self.ERROR_ACTION {
			self.currentAction -= self.ERROR_ACTION
			for {
				self.stateStackTop -= self.prs.rhs(self.currentAction) - 1
				self.currentAction = self.prs.ntAction( self.stateStack[self.stateStackTop],
														self.prs.lhs(self.currentAction))
				if self.currentAction <= self.NUM_RULES {
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
		self.stateStackTop += 1
		if self.stateStackTop >= len(self.stateStack) {
			self.reallocateStacks()
		}
		self.stateStack[self.stateStackTop] = self.currentAction
	} else {
		if self.currentAction == self.ERROR_ACTION {
			self.action.resetTo(save_action_length) // restore original action state.
		}
	}
	return self.currentAction
}

//
// Now do the final parse of the input based on the actions in
// the list "action" and the sequence of tokens in the token stream.
//
func (self *DeterministicParser) parseActions() interface{} {
	//
	// Indicate that we are processing actions now (for the incremental
	// parser) and that it's ok to use the utility functions to query the
	// parser.
	//
	self.taking_actions = true
	self.tokStream.reset()
	self.lastToken = self.tokStream.getPrevious(self.tokStream.peek())
	var curtok int
	if self.markerKind == 0 {
		curtok = self.tokStream.getToken()
	} else {
		curtok = self.lastToken
	}

	//
	// Reparse the input...
	//
	self.stateStackTop = -1
	self.currentAction = self.START_STATE
	var i int = 0
	for ; i < self.action.size(); i++ {
		//
		// if the parser needs to stop processing, it may do so here.
		//
		if self.monitor != nil && self.monitor.isCancelled() {
			self.taking_actions = false
			return nil
		}

		self.stateStackTop += 1
		self.stateStack[self.stateStackTop] = self.currentAction
		self.locationStack[self.stateStackTop] = curtok

		self.currentAction = self.action.get(i)
		if self.currentAction <= self.NUM_RULES { // a reduce action?

			self.stateStackTop-- // turn reduction intoshift-reduction
			self.processReductions()
		} else { // a shift or shift-reduce action

			self.lastToken = curtok
			curtok = self.tokStream.getToken()
			if self.currentAction > self.ERROR_ACTION {
				self.currentAction -= self.ERROR_ACTION
				self.processReductions()
			}
		}
	}

	self.taking_actions = false // indicate that we are done.
	self.action = nil
	if self.markerKind == 0 {
		return self.parseStack[0]
	} else {
		return self.parseStack[1]
	}
}
