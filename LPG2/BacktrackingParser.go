package lpg2

import "math"
type BacktrackingParser struct {
	*Stacks
	monitor         Monitor
	START_STATE     int
	NUM_RULES       int
	NT_OFFSET       int
	LA_STATE_OFFSET int
	EOFT_SYMBOL     int
	ERROR_SYMBOL    int
	ACCEPT_ACTION   int
	ERROR_ACTION    int

	lastToken     int
	currentAction int

	tokStream TokenStream
	prs       ParseTable
	ra        RuleAction

	action      *IntSegmentedTuple
	tokens      *IntTuple
	actionStack []int
	skipTokens  bool

	//
	// A starting marker indicates that we are dealing with an entry point
	// for a given nonterminal. We need to execute a shift action on the
	// marker in order to parse the entry point in question.
	//
	markerTokenIndex int
}

func NewBacktrackingParser(tokStream TokenStream, prs ParseTable,
	ra RuleAction, monitor Monitor) *BacktrackingParser {
	a := new(BacktrackingParser)
	a.skipTokens = false
	a.Stacks = NewStacks()
	a.action = NewIntSegmentedTuple(10, 1024)
	err := a.reset(tokStream, prs, ra, monitor)
	if err != nil {
		return nil
	}
	return a
}
//
// A starting marker indicates that we are dealing with an entry point
// for a given nonterminal. We need to execute a shift action on the
// marker in order to parse the entry point in question.
//
func (self *BacktrackingParser) getMarkerToken(marker_kind int, start_token_index int) (int, error) {
	if marker_kind == 0 {
		return 0, nil
	} else {
		_ipsream, ok := self.tokStream.(IPrsStream)
		if self.markerTokenIndex == 0 {

			if !ok {
				return -1, NewTokenStreamNotIPrsStreamException("")
			}
			self.markerTokenIndex = _ipsream.makeErrorToken(self.tokStream.getPrevious(start_token_index),
                                                            self.tokStream.getPrevious(start_token_index),
                                                            self.tokStream.getPrevious(start_token_index),
                                                            marker_kind)
		} else {
			_ipsream.getIToken(self.markerTokenIndex).setKind(marker_kind)
		}
	}
	return self.markerTokenIndex, nil
}

//
// Override the getToken function in Stacks.
//
func (self *BacktrackingParser) getToken(i int) int {
	return self.tokens.get(self.locationStack[self.stateStackTop+(i-1)])
}

func (self *BacktrackingParser) getCurrentRule() int {
	return self.currentAction
}
func (self *BacktrackingParser) getFirstToken() int {
	return self.tokStream.getFirstRealToken(self.getToken(1))
}
func (self *BacktrackingParser) getFirstTokenAt(i int) int {

	return self.tokStream.getFirstRealToken(self.getToken(i))
}
func (self *BacktrackingParser) getLastToken() int {
	return self.tokStream.getLastRealToken(self.lastToken)
}
func (self *BacktrackingParser) getLastTokenAt(i int) int {
	var l int
	if i >= self.prs.rhs(self.currentAction) {
		l = self.lastToken
	} else {
		l = self.tokens.get(self.locationStack[self.stateStackTop+i] - 1)
	}
	return self.tokStream.getLastRealToken(l)
}
func (self *BacktrackingParser) setMonitor(monitor Monitor) {
	self.monitor = monitor
}
func (self *BacktrackingParser) reset1() error {
	self.action.reset()
	self.skipTokens = false
	self.markerTokenIndex = 0
	return nil
}
func (self *BacktrackingParser) reset2(tokStream TokenStream, monitor Monitor) error {
	self.monitor = monitor
	self.tokStream = tokStream
	return self.reset1()
}

func (self *BacktrackingParser) reset(tokStream TokenStream, prs ParseTable, ra RuleAction, monitor Monitor) error {
	if prs != nil {
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
		if !prs.getBacktrack() {
			return NewNotBacktrackParseTableException("")
		}

	}
	if nil != ra {
		self.ra = ra
	}

	if nil == tokStream {
		err := self.reset1()
		if err != nil {
			return err
		}
		return nil
	}
	return self.reset2(tokStream, monitor)
}
func (self *BacktrackingParser) reset3(tokStream TokenStream, prs ParseTable, ra RuleAction) error {
	return self.reset(tokStream, prs, ra, nil)
}

//
// Allocate or reallocate all the stacks. Their sizes should always be the same.
//
func (self *BacktrackingParser) reallocateOtherStacks(startTokenIndex int) {
	if len(self.actionStack) == 0 {
		self.actionStack = make([]int, len(self.stateStack))
		self.locationStack = make([]int, len(self.stateStack))
		self.parseStack = make([]interface{}, len(self.stateStack))
		self.actionStack[0] = 0
		self.locationStack[0] = startTokenIndex
	} else {
		if len(self.actionStack) < len(self.stateStack) {
			var old_length int = len(self.actionStack)
			self.actionStack = arraycopy(self.actionStack, 0, make([]int, len(self.stateStack)), 0, old_length)
			self.locationStack = arraycopy(self.locationStack, 0, make([]int, len(self.stateStack)), 0, old_length)
			self.parseStack = ObjectArraycopy(self.parseStack, 0, make([]interface{}, len(self.stateStack)), 0, old_length)
		}
	}
	return
}
//
// Always attempt to recover
//
func (self *BacktrackingParser) fuzzyParse() (interface{}, error) {
    return self.fuzzyParseEntry(0, math.MaxInt32)
}
//
// Recover up to max_error_count times and then quit
//
func (self *BacktrackingParser) fuzzyParseWithErrorCount(max_error_count int) (interface{}, error) {
	return self.fuzzyParseEntry(0, max_error_count)
}

func (self *BacktrackingParser) fuzzyParseEntry(marker_kind int, max_error_count int) (interface{}, error) {

	self.action.reset()
	self.tokStream.reset() // Position at first token.
	self.reallocateStateStack()
	self.stateStackTop = 0
	self.stateStack[0] = self.START_STATE

	//
	// The tuple tokens will eventually contain the sequence
	// of tokens that resulted in a successful parse. We leave
	// it up to the "Stream" implementer to define the predecessor
	// of the first token as he sees fit.
	//
	var first_token int = self.tokStream.peek()
	var start_token int = first_token
	marker_token, _ := self.getMarkerToken(marker_kind, first_token)

	self.tokens = NewIntTupleWithEstimate(self.tokStream.getStreamLength())
	self.tokens.add(self.tokStream.getPrevious(first_token))

	var error_token int = self.backtrackParseInternal(self.action, marker_token)
	if error_token != 0 { // an error was detected?
		_stream, ok := self.tokStream.(IPrsStream)
		if !ok {
			return nil, NewTokenStreamNotIPrsStreamException("")
		}
		var rp = NewRecoveryParser(self, self.action, self.tokens, _stream, self.prs, max_error_count, 0, self.monitor)
		start_token, _ = rp.recover(marker_token, error_token)
	}
	if marker_token != 0 && start_token == first_token {
		self.tokens.add(marker_token)
	}
	var t int = start_token
	for ; self.tokStream.getKind(t) != self.EOFT_SYMBOL; t = self.tokStream.getNext(t) {
		self.tokens.add(t)
	}
	self.tokens.add(t)
	return self.parseActions(marker_kind)
}

//
// Parse input allowing up to max_error_count Error token recoveries.
// When max_error_count is 0, no Error token recoveries occur.
// When max_error is > 0, it limits the number of Error token recoveries.
// When max_error is < 0, the number of error token recoveries is unlimited.
// Also, such recoveries only require one token to be parsed beyond the recovery point.
// (normally two tokens beyond the recovery point must be parsed)
// Thus, a negative max_error_count should be used when error productions are used to
// skip tokens.
//
func (self *BacktrackingParser) parse(max_error_count int) (interface{}, error) {
	return self.parseEntry(0, max_error_count)
}

//
// Parse input allowing up to max_error_count Error token recoveries.
// When max_error_count is 0, no Error token recoveries occur.
// When max_error is > 0, it limits the int of Error token recoveries.
// When max_error is < 0, the int of error token recoveries is unlimited.
// Also, such recoveries only require one token to be parsed beyond the recovery point.
// (normally two tokens beyond the recovery point must be parsed)
// Thus, a negative max_error_count should be used when error productions are used to
// skip tokens.
//
func (self *BacktrackingParser) parseEntry(marker_kind int, max_error_count int) (interface{}, error) {
	self.action.reset()
	self.tokStream.reset() // Position at first token.

	self.reallocateStateStack()
	self.stateStackTop = 0
	self.stateStack[0] = self.START_STATE

	self.skipTokens = max_error_count < 0
	_stream, ok := self.tokStream.(IPrsStream)
	if max_error_count > 0 && ok {
		max_error_count = 0
	}
	//
	// The tuple tokens will eventually contain the sequence
	// of tokens that resulted in a successful parse. We leave
	// it up to the "Stream" implementer to define the predecessor
	// of the first token as he sees fit.
	//
	self.tokens = NewIntTupleWithEstimate(self.tokStream.getStreamLength())
	self.tokens.add(self.tokStream.getPrevious(self.tokStream.peek()))

	var start_token_index int = self.tokStream.peek()
	var repair_token, _ = self.getMarkerToken(marker_kind, start_token_index)
	var start_action_index int = self.action.size() // obviously 0
	var temp_stack []int = make([]int, self.stateStackTop+1)
	arraycopy(self.stateStack, 0, temp_stack, 0, len(temp_stack))

	var initial_error_token = self.backtrackParseInternal(self.action, repair_token)
	var error_token int = initial_error_token
	var count int = 0
	for error_token != 0 {
		if count == max_error_count {
			return nil, NewBadParseException(initial_error_token)
		}
		self.action.resetTo(start_action_index)
		self.tokStream.resetTo(start_token_index)
		self.stateStackTop = len(temp_stack) - 1
		arraycopy(temp_stack, 0, self.stateStack, 0, len(temp_stack))
		self.reallocateOtherStacks(start_token_index)

		self.backtrackParseUpToError(repair_token, error_token)

		for self.stateStackTop = self.findRecoveryStateIndex(self.stateStackTop);
            self.stateStackTop >= 0;
            self.stateStackTop = self.findRecoveryStateIndex(self.stateStackTop - 1) {
			var recovery_token = self.tokens.get(self.locationStack[self.stateStackTop] - 1)
			var temp int
			if recovery_token >= start_token_index {
				temp = recovery_token
			} else {
				temp = error_token
			}
			repair_token = self.errorRepair(_stream, temp, error_token)
			if repair_token != 0 {
				break
			}
		}
		if self.stateStackTop < 0 {
			return nil, NewBadParseException(initial_error_token)
		}
		temp_stack = make([]int, self.stateStackTop+1)
		arraycopy(self.stateStack, 0, temp_stack, 0, len(temp_stack))

		start_action_index = self.action.size()
		start_token_index = self.tokStream.peek()

		error_token = self.backtrackParseInternal(self.action, repair_token)
		count++
	}
	if repair_token != 0 {
		self.tokens.add(repair_token)
	}
	var t int = start_token_index
	for ; self.tokStream.getKind(t) != self.EOFT_SYMBOL; t = self.tokStream.getNext(t) {
		self.tokens.add(t)
	}
	self.tokens.add(t)
	return self.parseActions(marker_kind)
}

//
// Process reductions and continue...
//
func (self *BacktrackingParser) process_reductions() {
	for {
		self.stateStackTop -= (self.prs.rhs(self.currentAction) - 1)
		self.ra.ruleAction(self.currentAction)
		self.currentAction = self.prs.ntAction(self.stateStack[self.stateStackTop], self.prs.lhs(self.currentAction))
		if self.currentAction <= self.NUM_RULES {
			continue
		} else {
			break
		}
	}
	return
}

//
// Now do the final parse of the input based on the actions in
// the list "action" and the sequence of tokens in list "tokens".
//
func (self *BacktrackingParser) parseActions(marker_kind int) (interface{}, error) {
	var ti int = -1
	ti += 1
	self.lastToken = self.tokens.get(ti)

	ti += 1
	var curtok = self.tokens.get(ti)
	self.allocateOtherStacks()
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
			return nil, nil
		}
		self.stateStackTop += 1
		self.stateStack[self.stateStackTop] = self.currentAction

		self.locationStack[self.stateStackTop] = ti

		self.currentAction = self.action.get(i)
		if self.currentAction <= self.NUM_RULES { // a reduce action?
			self.stateStackTop-- // make reduction look like shift-reduction
			self.process_reductions()
		} else { // a shift or shift-reduce action
			if self.tokStream.getKind(curtok) > self.NT_OFFSET {
				_stream, _ := self.tokStream.(IPrsStream)
				badtok, _ := _stream.getIToken(curtok).(*ErrorToken)
				return nil, NewBadParseException(badtok.getErrorToken().getTokenIndex())
			}
			self.lastToken = curtok
			ti += 1
			curtok = self.tokens.get(ti)
			if self.currentAction > self.ERROR_ACTION {
				self.currentAction -= self.ERROR_ACTION
				self.process_reductions()
			}
		}
	}

	if marker_kind == 0 {
		return self.parseStack[0], nil
	} else {
		return self.parseStack[1], nil
	}

}

//
// Process reductions and continue...
//
func (self *BacktrackingParser) process_backtrack_reductions(act int) int {
	for {
		self.stateStackTop -= (self.prs.rhs(act) - 1)
		act = self.prs.ntAction(self.stateStack[self.stateStackTop], self.prs.lhs(act))
		if act <= self.NUM_RULES {
			continue
		} else {
			break
		}
	}
	return act
}

//
// This method is intended to be used by the type RecoveryParser.
// Note that the action tuple passed here must be the same action
// tuple that was passed down to RecoveryParser. It is passed back
// to self method as documention.
//
func (self *BacktrackingParser) backtrackParse(stack []int, stack_top int, action *IntSegmentedTuple, initial_token int) int {
	self.stateStackTop = stack_top
	arraycopy(stack, 0, self.stateStack, 0, self.stateStackTop+1)
	return self.backtrackParseInternal(action, initial_token)
}

//
// Parse the input until either the parse completes successfully or
// an error is encountered. This function returns an integer that
// represents the last action that was executed by the parser. If
// the parse was succesful, then the tuple "action" contains the
// successful sequence of actions that was executed.
//
func (self *BacktrackingParser) backtrackParseInternal(action *IntSegmentedTuple, initial_token int) int {
	//
	// Allocate configuration stack.
	//
	var configuration_stack = NewConfigurationStack(self.prs)

	//
	// Keep parsing until we successfully reach the end of file or
	// an error is encountered. The list of actions executed will
	// be stored in the "action" tuple.
	//
	var error_token int = 0
	var maxStackTop int = self.stateStackTop
	var start_token int = self.tokStream.peek()

	var curtok int
	if initial_token > 0 {
		curtok = initial_token
	} else {
		curtok = self.tokStream.getToken()
	}

	var current_kind int = self.tokStream.getKind(curtok)
	var act int = self.tAction(self.stateStack[self.stateStackTop], current_kind)
	//
	// The main driver loop
	//
	for {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if self.monitor != nil && self.monitor.isCancelled() {
			return 0
		}
		if act <= self.NUM_RULES {
			action.add(act) // save self reduce action
			self.stateStackTop--
			act = self.process_backtrack_reductions(act)
		} else {
			if act > self.ERROR_ACTION {
				action.add(act) // save self shift-reduce action
				curtok = self.tokStream.getToken()
				current_kind = self.tokStream.getKind(curtok)
				act = self.process_backtrack_reductions(act - self.ERROR_ACTION)
			} else {
				if act < self.ACCEPT_ACTION {
					action.add(act) // save self shift action
					curtok = self.tokStream.getToken()
					current_kind = self.tokStream.getKind(curtok)
				} else {
					if act == self.ERROR_ACTION {

						if !(error_token > curtok) {
							error_token = curtok
						}
						var configuration = configuration_stack.pop()
						if configuration == nil {
							act = self.ERROR_ACTION
						} else {
							action.resetTo(configuration.action_length)
							act = configuration.act
							curtok = configuration.curtok
							current_kind = self.tokStream.getKind(curtok)
							var index int
							if curtok == initial_token {
								index = start_token
							} else {
								index = self.tokStream.getNext(curtok)
							}
							self.tokStream.resetTo(index)
							self.stateStackTop = configuration.stack_top
							configuration.retrieveStack(self.stateStack)
							continue
						}
						break
					} else {
						if act > self.ACCEPT_ACTION {
							if configuration_stack.findConfiguration(self.stateStack, self.stateStackTop, curtok) {
								act = self.ERROR_ACTION
							} else {
								configuration_stack.push(self.stateStack, self.stateStackTop, act+1, curtok, action.size())
								act = self.prs.baseAction(act)
								if self.stateStackTop > maxStackTop {
									maxStackTop = self.stateStackTop
								}
							}
							continue
						} else {
							break // assert(act == ACCEPT_ACTION)
						}
					}
				}
			}
		}
		self.stateStackTop += 1
		if self.stateStackTop >= len(self.stateStack) {
			self.reallocateStateStack()
		}
		self.stateStack[self.stateStackTop] = act

		act = self.tAction(act, current_kind)
	}
	if act == self.ERROR_ACTION {
		return error_token
	} else {
		return 0
	}
}
func (self *BacktrackingParser) backtrackParseUpToError(initial_token int, error_token int) {
	//
	// Allocate configuration stack.
	//
	var configuration_stack = NewConfigurationStack(self.prs)

	//
	// Keep parsing until we successfully reach the end of file or
	// an error is encountered. The list of actions executed will
	// be stored in the "action" tuple.
	//
	var start_token int = self.tokStream.peek()
	var curtok int
	if initial_token > 0 {
		curtok = initial_token
	} else {
		curtok = self.tokStream.getToken()
	}
	var current_kind int = self.tokStream.getKind(curtok)
	var act int = self.tAction(self.stateStack[self.stateStackTop], current_kind)

	self.tokens.add(curtok)
	self.locationStack[self.stateStackTop] = self.tokens.size()
	self.actionStack[self.stateStackTop] = self.action.size()

	for {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if self.monitor != nil && self.monitor.isCancelled() {
			return
		}

		if act <= self.NUM_RULES {
			self.action.add(act) // save self reduce action
			self.stateStackTop--
			act = self.process_backtrack_reductions(act)
		} else {
			if act > self.ERROR_ACTION {
				self.action.add(act) // save self shift-reduce action
				curtok = self.tokStream.getToken()
				current_kind = self.tokStream.getKind(curtok)
				self.tokens.add(curtok)
				act = self.process_backtrack_reductions(act - self.ERROR_ACTION)
			} else {
				if act < self.ACCEPT_ACTION {
					self.action.add(act) // save self shift action
					curtok = self.tokStream.getToken()
					current_kind = self.tokStream.getKind(curtok)
					self.tokens.add(curtok)
				} else {
					if act == self.ERROR_ACTION {
						if curtok != error_token {
							var configuration = configuration_stack.pop()
							if configuration == nil {
								act = self.ERROR_ACTION
							} else {
								self.action.resetTo(configuration.action_length)
								act = configuration.act
								var next_token_index int = configuration.curtok
								self.tokens.resetTo(next_token_index)
								curtok = self.tokens.get(next_token_index - 1)
								current_kind = self.tokStream.getKind(curtok)
								var index int
								if curtok == initial_token {
									index = start_token
								} else {
									index = self.tokStream.getNext(curtok)
								}
								self.tokStream.resetTo(index)

								self.stateStackTop = configuration.stack_top
								configuration.retrieveStack(self.stateStack)
								self.locationStack[self.stateStackTop] = self.tokens.size()
								self.actionStack[self.stateStackTop] = self.action.size()
								continue
							}
						}
						break
					} else {
						if act > self.ACCEPT_ACTION {
							if configuration_stack.findConfiguration(self.stateStack, self.stateStackTop, self.tokens.size()) {
								act = self.ERROR_ACTION
							} else {
								configuration_stack.push(self.stateStack, self.stateStackTop, act+1, self.tokens.size(), self.action.size())
								act = self.prs.baseAction(act)
							}
							continue
						} else {
							break // assert(act == ACCEPT_ACTION)
						}
					}
				}
			}
		}

		self.stateStackTop += 1
		self.stateStack[self.stateStackTop] = act // no need to check if out of bounds

		self.locationStack[self.stateStackTop] = self.tokens.size()
		self.actionStack[self.stateStackTop] = self.action.size()
		act = self.tAction(act, current_kind)
	}
	return
}
func (self *BacktrackingParser) repairable(error_token int) bool {
	//
	// Allocate configuration stack.
	//
	var configuration_stack = NewConfigurationStack(self.prs)

	//
	// Keep parsing until we successfully reach the end of file or
	// an error is encountered. The list of actions executed will
	// be stored in the "action" tuple.
	//
	var start_token int = self.tokStream.peek()
	var final_token int = self.tokStream.getStreamLength() // unreachable
	var curtok int = 0
	var current_kind int = self.ERROR_SYMBOL
	var act int = self.tAction(self.stateStack[self.stateStackTop], current_kind)
	for {
		if act <= self.NUM_RULES {
			self.stateStackTop--
			act = self.process_backtrack_reductions(act)
		} else {
			if act > self.ERROR_ACTION {
				curtok = self.tokStream.getToken()
				if curtok > final_token {
					return true
				}
				current_kind = self.tokStream.getKind(curtok)
				act = self.process_backtrack_reductions(act - self.ERROR_ACTION)
			} else {
				if act < self.ACCEPT_ACTION {
					curtok = self.tokStream.getToken()
					if curtok > final_token {
						return true
					}
					current_kind = self.tokStream.getKind(curtok)
				} else {
					if act == self.ERROR_ACTION {
						var configuration = configuration_stack.pop()
						if configuration == nil {
							act = self.ERROR_ACTION
						} else {
							self.stateStackTop = configuration.stack_top
							configuration.retrieveStack(self.stateStack)
							act = configuration.act
							curtok = configuration.curtok
							if curtok == 0 {
								current_kind = self.ERROR_SYMBOL
								self.tokStream.resetTo(start_token)
							} else {
								current_kind = self.tokStream.getKind(curtok)
								self.tokStream.resetTo(self.tokStream.getNext(curtok))
							}
							continue
						}
						break
					} else {
						if act > self.ACCEPT_ACTION {
							if configuration_stack.findConfiguration(self.stateStack, self.stateStackTop, curtok) {
								act = self.ERROR_ACTION
							} else {
								configuration_stack.push(self.stateStack, self.stateStackTop, act+1, curtok, 0)
								act = self.prs.baseAction(act)
							}
							continue
						} else {
							break // assert(act == ACCEPT_ACTION)
						}
					}
				}
			}
		}

		//
		// We consider a configuration to be acceptable for recovery
		// if we are able to consume enough symbols in the remainining
		// tokens to reach another potential recovery point past the
		// original error token.
		//
		if (curtok > error_token) && (final_token == self.tokStream.getStreamLength()) {
			//
			// If the ERROR_SYMBOL is a valid Action Adjunct in the state
			// "act" then we set the terminating token as the successor of
			// the current token. I.e., we have to be able to parse at least
			// two tokens past the resynch point before we claim victory.
			//
			if self.recoverableState(act) {
				if self.skipTokens {
					final_token = curtok
				} else {
					final_token = self.tokStream.getNext(curtok)
				}

			}
		}
		self.stateStackTop += 1
		if self.stateStackTop >= len(self.stateStack) {
			self.reallocateStateStack()
		}
		self.stateStack[self.stateStackTop] = act

		act = self.tAction(act, current_kind)
	}
	//
	// If we can reach the end of the input successfully, we claim victory.
	//
	return act == self.ACCEPT_ACTION
}
func (self *BacktrackingParser) recoverableState(state int) bool {
	var k int = self.prs.asi(state)
	for ; self.prs.asr(k) != 0; k++ {
		if self.prs.asr(k) == self.ERROR_SYMBOL {
			return true
		}
	}
	return false
}
func (self *BacktrackingParser) findRecoveryStateIndex(start_index int) int {
	var i int = start_index
	for ; i >= 0; i-- {
		//
		// If the ERROR_SYMBOL is an Action Adjunct in state stateStack[i]
		// then chose i as the index of the state to recover on.
		//
		if self.recoverableState(self.stateStack[i]) {
			break
		}
	}
	if i >= 0 { // if a recoverable state, remove null reductions, if any.
		var k = i - 1
		for ; k >= 0; k-- {
			if self.locationStack[k] != self.locationStack[i] {
				break
			}
		}
		i = k + 1
	}
	return i
}

func (self *BacktrackingParser) errorRepair(stream IPrsStream, recovery_token int, error_token int) int {
	var temp_stack []int = make([]int, self.stateStackTop+1)
	arraycopy(self.stateStack, 0, temp_stack, 0, len(temp_stack))
	for ; stream.getKind(recovery_token) != self.EOFT_SYMBOL;
	      recovery_token = stream.getNext(recovery_token) {
		stream.resetTo(recovery_token)
		if self.repairable(error_token) {
			break
		}
		self.stateStackTop = len(temp_stack) - 1
		arraycopy(temp_stack, 0, self.stateStack, 0, len(temp_stack))
	}

	if stream.getKind(recovery_token) == self.EOFT_SYMBOL {
		stream.resetTo(recovery_token)
		if !self.repairable(error_token) {
			self.stateStackTop = len(temp_stack) - 1
			arraycopy(temp_stack, 0, self.stateStack, 0, len(temp_stack))
			return 0
		}
	}

	self.stateStackTop = len(temp_stack) - 1
	arraycopy(temp_stack, 0, self.stateStack, 0, len(temp_stack))
	stream.resetTo(recovery_token)
	self.tokens.resetTo(self.locationStack[self.stateStackTop] - 1)
	self.action.resetTo(self.actionStack[self.stateStackTop])

	return stream.makeErrorToken(   self.tokens.get(self.locationStack[self.stateStackTop]-1),
                                    stream.getPrevious(recovery_token),
                                    error_token,
                                    self.ERROR_SYMBOL)
}

//
// keep looking ahead until we compute a valid action
//
func (self *BacktrackingParser) lookahead(act int, token int) int {
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
func (self *BacktrackingParser) tAction(act int, sym int) int {
	act = self.prs.tAction(act, sym)
	if act > self.LA_STATE_OFFSET {
		return self.lookahead(act, self.tokStream.peek())
	} else {
		return act
	}
}
