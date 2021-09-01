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
	err := a.Reset(tokStream, prs, ra, monitor)
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
func (my *BacktrackingParser) GetMarkerToken(marker_kind int, start_token_index int) (int, error) {
	if marker_kind == 0 {
		return 0, nil
	} else {
		_ipsream, ok := my.tokStream.(IPrsStream)
		if my.markerTokenIndex == 0 {

			if !ok {
				return -1, NewTokenStreamNotIPrsStreamException("")
			}
			my.markerTokenIndex = _ipsream.MakeErrorToken(my.tokStream.GetPrevious(start_token_index),
                                                            my.tokStream.GetPrevious(start_token_index),
                                                            my.tokStream.GetPrevious(start_token_index),
                                                            marker_kind)
		} else {
			_ipsream.GetIToken(my.markerTokenIndex).SetKind(marker_kind)
		}
	}
	return my.markerTokenIndex, nil
}

//
// Override the GetToken function in Stacks.
//
func (my *BacktrackingParser) GetToken(i int) int {
	return my.tokens.Get(my.locationStack[my.stateStackTop+(i-1)])
}

func (my *BacktrackingParser) GetCurrentRule() int {
	return my.currentAction
}
func (my *BacktrackingParser) GetFirstToken() int {
	return my.tokStream.GetFirstRealToken(my.GetToken(1))
}
func (my *BacktrackingParser) GetFirstTokenAt(i int) int {

	return my.tokStream.GetFirstRealToken(my.GetToken(i))
}
func (my *BacktrackingParser) GetLastToken() int {
	return my.tokStream.GetLastRealToken(my.lastToken)
}
func (my *BacktrackingParser) GetLastTokenAt(i int) int {
	var l int
	if i >= my.prs.Rhs(my.currentAction) {
		l = my.lastToken
	} else {
		l = my.tokens.Get(my.locationStack[my.stateStackTop+i] - 1)
	}
	return my.tokStream.GetLastRealToken(l)
}
func (my *BacktrackingParser) SetMonitor(monitor Monitor) {
	my.monitor = monitor
}
func (my *BacktrackingParser) Reset1() error {
	my.action.ReSet()
	my.skipTokens = false
	my.markerTokenIndex = 0
	return nil
}
func (my *BacktrackingParser) Reset2(tokStream TokenStream, monitor Monitor) error {
	my.monitor = monitor
	my.tokStream = tokStream
	return my.Reset1()
}

func (my *BacktrackingParser) Reset(tokStream TokenStream, prs ParseTable, ra RuleAction, monitor Monitor) error {
	if prs != nil {
		my.prs = prs
		my.START_STATE = prs.GetStartState()
		my.NUM_RULES = prs.GetNumRules()
		my.NT_OFFSET = prs.GetNtOffSet()
		my.LA_STATE_OFFSET = prs.GetLaStateOffSet()
		my.EOFT_SYMBOL = prs.GetEoftSymbol()
		my.ERROR_SYMBOL = prs.GetErrorSymbol()
		my.ACCEPT_ACTION = prs.GetAcceptAction()
		my.ERROR_ACTION = prs.GetErrorAction()
		if !prs.IsValidForParser() {
			return NewBadParseSymFileException("")
		}
		if !prs.GetBacktrack() {
			return NewNotBacktrackParseTableException("")
		}

	}
	if nil != ra {
		my.ra = ra
	}

	if nil == tokStream {
		err := my.Reset1()
		if err != nil {
			return err
		}
		return nil
	}
	return my.Reset2(tokStream, monitor)
}
func (my *BacktrackingParser) Reset3(tokStream TokenStream, prs ParseTable, ra RuleAction) error {
	return my.Reset(tokStream, prs, ra, nil)
}

//
// Allocate or reallocate all the stacks. Their sizes should always be the same.
//
func (my *BacktrackingParser) ReallocateOtherStacks(startTokenIndex int) {
	if len(my.actionStack) == 0 {
		my.actionStack = make([]int, len(my.stateStack))
		my.locationStack = make([]int, len(my.stateStack))
		my.parseStack = make([]interface{}, len(my.stateStack))
		my.actionStack[0] = 0
		my.locationStack[0] = startTokenIndex
	} else {
		if len(my.actionStack) < len(my.stateStack) {
			var old_length int = len(my.actionStack)
			my.actionStack = Arraycopy(my.actionStack, 0, make([]int, len(my.stateStack)), 0, old_length)
			my.locationStack = Arraycopy(my.locationStack, 0, make([]int, len(my.stateStack)), 0, old_length)
			my.parseStack = ObjectArraycopy(my.parseStack, 0, make([]interface{}, len(my.stateStack)), 0, old_length)
		}
	}
	return
}
//
// Always attempt to recover
//
func (my *BacktrackingParser) FuzzyParse() (interface{}, error) {
    return my.FuzzyParseEntry(0, math.MaxInt32)
}
//
// Recover up to max_error_count times and then quit
//
func (my *BacktrackingParser) FuzzyParseWithErrorCount(max_error_count int) (interface{}, error) {
	return my.FuzzyParseEntry(0, max_error_count)
}

func (my *BacktrackingParser) FuzzyParseEntry(marker_kind int, max_error_count int) (interface{}, error) {

	my.action.ReSet()
	my.tokStream.Reset() // Position at first token.
	my.ReallocateStateStack()
	my.stateStackTop = 0
	my.stateStack[0] = my.START_STATE

	//
	// The tuple tokens will eventually contain the sequence
	// of tokens that resulted in a successful parse. We leave
	// it up to the "Stream" implementer to define the predecessor
	// of the first token as he sees fit.
	//
	var first_token int = my.tokStream.Peek()
	var start_token int = first_token
	marker_token, _ := my.GetMarkerToken(marker_kind, first_token)

	my.tokens = NewIntTupleWithEstimate(my.tokStream.GetStreamLength())
	my.tokens.Add(my.tokStream.GetPrevious(first_token))

	var error_token int = my.backtrackParseInternal(my.action, marker_token)
	if error_token != 0 { // an error was detected?
		_stream, ok := my.tokStream.(IPrsStream)
		if !ok {
			return nil, NewTokenStreamNotIPrsStreamException("")
		}
		var rp = NewRecoveryParser(my, my.action, my.tokens, _stream, my.prs, max_error_count, 0, my.monitor)
		start_token, _ = rp.Recover(marker_token, error_token)
	}
	if marker_token != 0 && start_token == first_token {
		my.tokens.Add(marker_token)
	}
	var t int = start_token
	for ; my.tokStream.GetKind(t) != my.EOFT_SYMBOL; t = my.tokStream.GetNext(t) {
		my.tokens.Add(t)
	}
	my.tokens.Add(t)
	return my.ParseActions(marker_kind)
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
func (my *BacktrackingParser) Parse(max_error_count int) (interface{}, error) {
	return my.ParseEntry(0, max_error_count)
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
func (my *BacktrackingParser) ParseEntry(marker_kind int, max_error_count int) (interface{}, error) {
	my.action.ReSet()
	my.tokStream.Reset() // Position at first token.

	my.ReallocateStateStack()
	my.stateStackTop = 0
	my.stateStack[0] = my.START_STATE

	my.skipTokens = max_error_count < 0
	_stream, ok := my.tokStream.(IPrsStream)
	if max_error_count > 0 && ok {
		max_error_count = 0
	}
	//
	// The tuple tokens will eventually contain the sequence
	// of tokens that resulted in a successful parse. We leave
	// it up to the "Stream" implementer to define the predecessor
	// of the first token as he sees fit.
	//
	my.tokens = NewIntTupleWithEstimate(my.tokStream.GetStreamLength())
	my.tokens.Add(my.tokStream.GetPrevious(my.tokStream.Peek()))

	var start_token_index int = my.tokStream.Peek()
	var repair_token, _ = my.GetMarkerToken(marker_kind, start_token_index)
	var start_action_index int = my.action.Size() // obviously 0
	var temp_stack []int = make([]int, my.stateStackTop+1)
	Arraycopy(my.stateStack, 0, temp_stack, 0, len(temp_stack))

	var initial_error_token = my.backtrackParseInternal(my.action, repair_token)
	var error_token int = initial_error_token
	var count int = 0
	for error_token != 0 {
		if count == max_error_count {
			return nil, NewBadParseException(initial_error_token)
		}
		my.action.ReSetTo(start_action_index)
		my.tokStream.ResetTo(start_token_index)
		my.stateStackTop = len(temp_stack) - 1
		Arraycopy(temp_stack, 0, my.stateStack, 0, len(temp_stack))
		my.ReallocateOtherStacks(start_token_index)

		my.BacktrackParseUpToError(repair_token, error_token)

		for my.stateStackTop = my.FindRecoveryStateIndex(my.stateStackTop);
            my.stateStackTop >= 0;
            my.stateStackTop = my.FindRecoveryStateIndex(my.stateStackTop - 1) {
			var recovery_token = my.tokens.Get(my.locationStack[my.stateStackTop] - 1)
			var temp int
			if recovery_token >= start_token_index {
				temp = recovery_token
			} else {
				temp = error_token
			}
			repair_token = my.ErrorRepair(_stream, temp, error_token)
			if repair_token != 0 {
				break
			}
		}
		if my.stateStackTop < 0 {
			return nil, NewBadParseException(initial_error_token)
		}
		temp_stack = make([]int, my.stateStackTop+1)
		Arraycopy(my.stateStack, 0, temp_stack, 0, len(temp_stack))

		start_action_index = my.action.Size()
		start_token_index = my.tokStream.Peek()

		error_token = my.backtrackParseInternal(my.action, repair_token)
		count++
	}
	if repair_token != 0 {
		my.tokens.Add(repair_token)
	}
	var t int = start_token_index
	for ; my.tokStream.GetKind(t) != my.EOFT_SYMBOL; t = my.tokStream.GetNext(t) {
		my.tokens.Add(t)
	}
	my.tokens.Add(t)
	return my.ParseActions(marker_kind)
}

//
// Process reductions and continue...
//
func (my *BacktrackingParser) Process_reductions() {
	for {
		my.stateStackTop -= (my.prs.Rhs(my.currentAction) - 1)
		my.ra.RuleAction(my.currentAction)
		my.currentAction = my.prs.NtAction(my.stateStack[my.stateStackTop], my.prs.Lhs(my.currentAction))
		if my.currentAction <= my.NUM_RULES {
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
func (my *BacktrackingParser) ParseActions(marker_kind int) (interface{}, error) {
	var ti int = -1
	ti += 1
	my.lastToken = my.tokens.Get(ti)

	ti += 1
	var curtok = my.tokens.Get(ti)
	my.AllocateOtherStacks()
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
			return nil, nil
		}
		my.stateStackTop += 1
		my.stateStack[my.stateStackTop] = my.currentAction

		my.locationStack[my.stateStackTop] = ti

		my.currentAction = my.action.Get(i)
		if my.currentAction <= my.NUM_RULES { // a reduce action?
			my.stateStackTop-- // make reduction look like shift-reduction
			my.Process_reductions()
		} else { // a shift or shift-reduce action
			if my.tokStream.GetKind(curtok) > my.NT_OFFSET {
				_stream, _ := my.tokStream.(IPrsStream)
				badtok, _ := _stream.GetIToken(curtok).(*ErrorToken)
				return nil, NewBadParseException(badtok.GetErrorToken().GetTokenIndex())
			}
			my.lastToken = curtok
			ti += 1
			curtok = my.tokens.Get(ti)
			if my.currentAction > my.ERROR_ACTION {
				my.currentAction -= my.ERROR_ACTION
				my.Process_reductions()
			}
		}
	}

	if marker_kind == 0 {
		return my.parseStack[0], nil
	} else {
		return my.parseStack[1], nil
	}

}

//
// Process reductions and continue...
//
func (my *BacktrackingParser) Process_backtrack_reductions(act int) int {
	for {
		my.stateStackTop -= (my.prs.Rhs(act) - 1)
		act = my.prs.NtAction(my.stateStack[my.stateStackTop], my.prs.Lhs(act))
		if act <= my.NUM_RULES {
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
// to my method as documention.
//
func (my *BacktrackingParser) BacktrackParse(stack []int, stack_top int, action *IntSegmentedTuple, initial_token int) int {
	my.stateStackTop = stack_top
	Arraycopy(stack, 0, my.stateStack, 0, my.stateStackTop+1)
	return my.backtrackParseInternal(action, initial_token)
}

//
// Parse the input until either the parse completes successfully or
// an error is encountered. This function returns an integer that
// represents the last action that was executed by the parser. If
// the parse was succesful, then the tuple "action" contains the
// successful sequence of actions that was executed.
//
func (my *BacktrackingParser) backtrackParseInternal(action *IntSegmentedTuple, initial_token int) int {
	//
	// Allocate configuration stack.
	//
	var configuration_stack = NewConfigurationStack(my.prs)

	//
	// Keep parsing until we successfully reach the end of file or
	// an error is encountered. The list of actions executed will
	// be stored in the "action" tuple.
	//
	var error_token int = 0
	var maxStackTop int = my.stateStackTop
	var start_token int = my.tokStream.Peek()

	var curtok int
	if initial_token > 0 {
		curtok = initial_token
	} else {
		curtok = my.tokStream.GetToken()
	}

	var current_kind int = my.tokStream.GetKind(curtok)
	var act int = my.TAction(my.stateStack[my.stateStackTop], current_kind)
	//
	// The main driver loop
	//
	for {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if my.monitor != nil && my.monitor.IsCancelled() {
			return 0
		}
		if act <= my.NUM_RULES {
			action.Add(act) // save my reduce action
			my.stateStackTop--
			act = my.Process_backtrack_reductions(act)
		} else {
			if act > my.ERROR_ACTION {
				action.Add(act) // save my shift-reduce action
				curtok = my.tokStream.GetToken()
				current_kind = my.tokStream.GetKind(curtok)
				act = my.Process_backtrack_reductions(act - my.ERROR_ACTION)
			} else {
				if act < my.ACCEPT_ACTION {
					action.Add(act) // save my shift action
					curtok = my.tokStream.GetToken()
					current_kind = my.tokStream.GetKind(curtok)
				} else {
					if act == my.ERROR_ACTION {

						if !(error_token > curtok) {
							error_token = curtok
						}
						var configuration = configuration_stack.Pop()
						if configuration == nil {
							act = my.ERROR_ACTION
						} else {
							action.ReSetTo(configuration.action_length)
							act = configuration.act
							curtok = configuration.curtok
							current_kind = my.tokStream.GetKind(curtok)
							var index int
							if curtok == initial_token {
								index = start_token
							} else {
								index = my.tokStream.GetNext(curtok)
							}
							my.tokStream.ResetTo(index)
							my.stateStackTop = configuration.stack_top
							configuration.RetrieveStack(my.stateStack)
							continue
						}
						break
					} else {
						if act > my.ACCEPT_ACTION {
							if configuration_stack.FindConfiguration(my.stateStack, my.stateStackTop, curtok) {
								act = my.ERROR_ACTION
							} else {
								configuration_stack.Push(my.stateStack, my.stateStackTop, act+1, curtok, action.Size())
								act = my.prs.BaseAction(act)
								if my.stateStackTop > maxStackTop {
									maxStackTop = my.stateStackTop
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
		my.stateStackTop += 1
		if my.stateStackTop >= len(my.stateStack) {
			my.ReallocateStateStack()
		}
		my.stateStack[my.stateStackTop] = act

		act = my.TAction(act, current_kind)
	}
	if act == my.ERROR_ACTION {
		return error_token
	} else {
		return 0
	}
}
func (my *BacktrackingParser) BacktrackParseUpToError(initial_token int, error_token int) {
	//
	// Allocate configuration stack.
	//
	var configuration_stack = NewConfigurationStack(my.prs)

	//
	// Keep parsing until we successfully reach the end of file or
	// an error is encountered. The list of actions executed will
	// be stored in the "action" tuple.
	//
	var start_token int = my.tokStream.Peek()
	var curtok int
	if initial_token > 0 {
		curtok = initial_token
	} else {
		curtok = my.tokStream.GetToken()
	}
	var current_kind int = my.tokStream.GetKind(curtok)
	var act int = my.TAction(my.stateStack[my.stateStackTop], current_kind)

	my.tokens.Add(curtok)
	my.locationStack[my.stateStackTop] = my.tokens.Size()
	my.actionStack[my.stateStackTop] = my.action.Size()

	for {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if my.monitor != nil && my.monitor.IsCancelled() {
			return
		}

		if act <= my.NUM_RULES {
			my.action.Add(act) // save my reduce action
			my.stateStackTop--
			act = my.Process_backtrack_reductions(act)
		} else {
			if act > my.ERROR_ACTION {
				my.action.Add(act) // save my shift-reduce action
				curtok = my.tokStream.GetToken()
				current_kind = my.tokStream.GetKind(curtok)
				my.tokens.Add(curtok)
				act = my.Process_backtrack_reductions(act - my.ERROR_ACTION)
			} else {
				if act < my.ACCEPT_ACTION {
					my.action.Add(act) // save my shift action
					curtok = my.tokStream.GetToken()
					current_kind = my.tokStream.GetKind(curtok)
					my.tokens.Add(curtok)
				} else {
					if act == my.ERROR_ACTION {
						if curtok != error_token {
							var configuration = configuration_stack.Pop()
							if configuration == nil {
								act = my.ERROR_ACTION
							} else {
								my.action.ReSetTo(configuration.action_length)
								act = configuration.act
								var next_token_index int = configuration.curtok
								my.tokens.ReSetTo(next_token_index)
								curtok = my.tokens.Get(next_token_index - 1)
								current_kind = my.tokStream.GetKind(curtok)
								var index int
								if curtok == initial_token {
									index = start_token
								} else {
									index = my.tokStream.GetNext(curtok)
								}
								my.tokStream.ResetTo(index)

								my.stateStackTop = configuration.stack_top
								configuration.RetrieveStack(my.stateStack)
								my.locationStack[my.stateStackTop] = my.tokens.Size()
								my.actionStack[my.stateStackTop] = my.action.Size()
								continue
							}
						}
						break
					} else {
						if act > my.ACCEPT_ACTION {
							if configuration_stack.FindConfiguration(my.stateStack, my.stateStackTop, my.tokens.Size()) {
								act = my.ERROR_ACTION
							} else {
								configuration_stack.Push(my.stateStack, my.stateStackTop, act+1, my.tokens.Size(), my.action.Size())
								act = my.prs.BaseAction(act)
							}
							continue
						} else {
							break // assert(act == ACCEPT_ACTION)
						}
					}
				}
			}
		}

		my.stateStackTop += 1
		my.stateStack[my.stateStackTop] = act // no need to check if out of bounds

		my.locationStack[my.stateStackTop] = my.tokens.Size()
		my.actionStack[my.stateStackTop] = my.action.Size()
		act = my.TAction(act, current_kind)
	}
	return
}
func (my *BacktrackingParser) Repairable(error_token int) bool {
	//
	// Allocate configuration stack.
	//
	var configuration_stack = NewConfigurationStack(my.prs)

	//
	// Keep parsing until we successfully reach the end of file or
	// an error is encountered. The list of actions executed will
	// be stored in the "action" tuple.
	//
	var start_token int = my.tokStream.Peek()
	var final_token int = my.tokStream.GetStreamLength() // unreachable
	var curtok int = 0
	var current_kind int = my.ERROR_SYMBOL
	var act int = my.TAction(my.stateStack[my.stateStackTop], current_kind)
	for {
		if act <= my.NUM_RULES {
			my.stateStackTop--
			act = my.Process_backtrack_reductions(act)
		} else {
			if act > my.ERROR_ACTION {
				curtok = my.tokStream.GetToken()
				if curtok > final_token {
					return true
				}
				current_kind = my.tokStream.GetKind(curtok)
				act = my.Process_backtrack_reductions(act - my.ERROR_ACTION)
			} else {
				if act < my.ACCEPT_ACTION {
					curtok = my.tokStream.GetToken()
					if curtok > final_token {
						return true
					}
					current_kind = my.tokStream.GetKind(curtok)
				} else {
					if act == my.ERROR_ACTION {
						var configuration = configuration_stack.Pop()
						if configuration == nil {
							act = my.ERROR_ACTION
						} else {
							my.stateStackTop = configuration.stack_top
							configuration.RetrieveStack(my.stateStack)
							act = configuration.act
							curtok = configuration.curtok
							if curtok == 0 {
								current_kind = my.ERROR_SYMBOL
								my.tokStream.ResetTo(start_token)
							} else {
								current_kind = my.tokStream.GetKind(curtok)
								my.tokStream.ResetTo(my.tokStream.GetNext(curtok))
							}
							continue
						}
						break
					} else {
						if act > my.ACCEPT_ACTION {
							if configuration_stack.FindConfiguration(my.stateStack, my.stateStackTop, curtok) {
								act = my.ERROR_ACTION
							} else {
								configuration_stack.Push(my.stateStack, my.stateStackTop, act+1, curtok, 0)
								act = my.prs.BaseAction(act)
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
		if (curtok > error_token) && (final_token == my.tokStream.GetStreamLength()) {
			//
			// If the ERROR_SYMBOL is a valid Action Adjunct in the state
			// "act" then we Set the terminating token as the successor of
			// the current token. I.e., we have to be able to parse at least
			// two tokens past the resynch point before we claim victory.
			//
			if my.recoverableState(act) {
				if my.skipTokens {
					final_token = curtok
				} else {
					final_token = my.tokStream.GetNext(curtok)
				}

			}
		}
		my.stateStackTop += 1
		if my.stateStackTop >= len(my.stateStack) {
			my.ReallocateStateStack()
		}
		my.stateStack[my.stateStackTop] = act

		act = my.TAction(act, current_kind)
	}
	//
	// If we can reach the end of the input successfully, we claim victory.
	//
	return act == my.ACCEPT_ACTION
}
func (my *BacktrackingParser) recoverableState(state int) bool {
	var k int = my.prs.Asi(state)
	for ; my.prs.Asr(k) != 0; k++ {
		if my.prs.Asr(k) == my.ERROR_SYMBOL {
			return true
		}
	}
	return false
}
func (my *BacktrackingParser) FindRecoveryStateIndex(start_index int) int {
	var i int = start_index
	for ; i >= 0; i-- {
		//
		// If the ERROR_SYMBOL is an Action Adjunct in state stateStack[i]
		// then chose i as the index of the state to recover on.
		//
		if my.recoverableState(my.stateStack[i]) {
			break
		}
	}
	if i >= 0 { // if a recoverable state, remove null reductions, if any.
		var k = i - 1
		for ; k >= 0; k-- {
			if my.locationStack[k] != my.locationStack[i] {
				break
			}
		}
		i = k + 1
	}
	return i
}

func (my *BacktrackingParser) ErrorRepair(stream IPrsStream, recovery_token int, error_token int) int {
	var temp_stack []int = make([]int, my.stateStackTop+1)
	Arraycopy(my.stateStack, 0, temp_stack, 0, len(temp_stack))
	for ; stream.GetKind(recovery_token) != my.EOFT_SYMBOL;
	      recovery_token = stream.GetNext(recovery_token) {
		stream.ResetTo(recovery_token)
		if my.Repairable(error_token) {
			break
		}
		my.stateStackTop = len(temp_stack) - 1
		Arraycopy(temp_stack, 0, my.stateStack, 0, len(temp_stack))
	}

	if stream.GetKind(recovery_token) == my.EOFT_SYMBOL {
		stream.ResetTo(recovery_token)
		if !my.Repairable(error_token) {
			my.stateStackTop = len(temp_stack) - 1
			Arraycopy(temp_stack, 0, my.stateStack, 0, len(temp_stack))
			return 0
		}
	}

	my.stateStackTop = len(temp_stack) - 1
	Arraycopy(temp_stack, 0, my.stateStack, 0, len(temp_stack))
	stream.ResetTo(recovery_token)
	my.tokens.ReSetTo(my.locationStack[my.stateStackTop] - 1)
	my.action.ReSetTo(my.actionStack[my.stateStackTop])

	return stream.MakeErrorToken(   my.tokens.Get(my.locationStack[my.stateStackTop]-1),
                                    stream.GetPrevious(recovery_token),
                                    error_token,
                                    my.ERROR_SYMBOL)
}

//
// keep looking ahead until we compute a valid action
//
func (my *BacktrackingParser) Lookahead(act int, token int) int {
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
func (my *BacktrackingParser) TAction(act int, sym int) int {
	act = my.prs.TAction(act, sym)
	if act > my.LA_STATE_OFFSET {
		return my.Lookahead(act, my.tokStream.Peek())
	} else {
		return act
	}
}
