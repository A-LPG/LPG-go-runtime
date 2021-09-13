package lpg2

type RecoveryParser struct {
	*DiagnoseParser
	parser       *BacktrackingParser
	action       *IntSegmentedTuple
	tokens       *IntTuple
	actionStack  []int
	scope_repair *PrimaryRepairInfo
}

// NewRecoveryParser
// maxErrors is the maximum int of errors to be repaired
// maxTime is the maximum amount of time allowed for diagnosing
// but at least one error must be diagnosed
//
func NewRecoveryParser(parser *BacktrackingParser, action *IntSegmentedTuple, tokens *IntTuple, tokStream IPrsStream,
	prs ParseTable, maxErrors int, maxTime int, monitor Monitor) *RecoveryParser {
	t := new(RecoveryParser)
	t.scope_repair = NewPrimaryRepairInfo()
	t.DiagnoseParser = NewDiagnoseParserExt(t,tokStream, prs, maxErrors, maxTime, monitor)
	t.parser = parser
	t.action = action
	t.tokens = tokens
	return t
}

func (my *RecoveryParser) ReallocateStacks() {
	my.DiagnoseParser.ReallocateStacks()
	if len(my.actionStack) == 0 {
		my.actionStack = make([]int, len(my.stateStack))
	} else {
		var old_stack_length int = len(my.actionStack)
		my.actionStack = Arraycopy(my.actionStack, 0, make([]int, len(my.stateStack)), 0, old_stack_length)
	}
	return
}
func (my *RecoveryParser) ReportError(scope_index int, error_token int) {
	var text string = "\""
	var i int = my.ScopeSuffix(scope_index)
	for ; my.ScopeRhs(i) != 0; i++ {
		if !my.IsNullable(my.ScopeRhs(i)) {

			var symbol_index int
			if my.ScopeRhs(i) > my.NT_OFFSET {
				symbol_index = my.NonterminalIndex(my.ScopeRhs(i) - my.NT_OFFSET)
			} else {
				symbol_index = my.TerminalIndex(my.ScopeRhs(i))
			}

			if len(my.Name(symbol_index)) > 0 {
				if len(text) > 1 { // Not just starting quote?
					text += " " // add a space separator
				}
				text += my.Name(symbol_index)
			}
		}
	}
	text += "\""
	my.tokStream.ReportError(SCOPE_CODE, error_token, error_token, []string{text}, 0)
	return
}
func (my *RecoveryParser) Recover(marker_token int, error_token int) (int, error) {
	if len(my.stateStack) == 0 {
		my.ReallocateStacks()
	}

	my.tokens.Reset()
	my.tokStream.Reset()
	my.tokens.Add(my.tokStream.GetPrevious(my.tokStream.Peek()))
	var restart_token int
	if marker_token != 0 {
		restart_token = marker_token
	} else {
		restart_token = my.tokStream.GetToken()
	}

	var old_action_size int = 0
	my.stateStackTop = 0
	my.stateStack[my.stateStackTop] = my.START_STATE
	for {
		my.action.ResetTo(old_action_size)
		if !my.FixError(restart_token, error_token) {
			return -1, NewBadParseException(error_token)
		}
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if nil != my.monitor && my.monitor.IsCancelled() {
			break
		}
		//
		// At my stage, we have a recovery configuration. See how
		// far we can go with it.
		//
		restart_token = error_token
		my.tokStream.ResetTo(error_token)
		old_action_size = my.action.Size() // save the old Size in case we encounter a new error
		error_token = my.parser.BacktrackParse(my.stateStack, my.stateStackTop, my.action, 0)
		my.tokStream.ResetTo(my.tokStream.GetNext(restart_token))
		if error_token != 0 {
			continue
		} else {
			break
		}
	} // no error found
	return restart_token, nil
}

//
// Given the configuration consisting of the states in stateStack
// and the sequence of tokens (current_kind, followed by the tokens
// in tokStream), FixError parses up to error_token in the tokStream
// recovers, if possible, from that error and returns the result.
// While doing my, it also computes the location_stack information
// and the sequence of actions that matches up with the result that
// it returns.
//
func (my *RecoveryParser) FixError(start_token int, error_token int) bool {
	//
	// Save information about the current configuration.
	//
	var curtok int = start_token
	var current_kind int = my.tokStream.GetKind(curtok)
	var first_stream_token int = my.tokStream.Peek()

	my.buffer[1] = error_token
	my.buffer[0] = my.tokStream.GetPrevious(my.buffer[1])
	var k int = 2
	for ; k < BUFF_SIZE; k++ {
		my.buffer[k] = my.tokStream.GetNext(my.buffer[k-1])
	}

	my.scope_repair.distance = 0
	my.scope_repair.misspellIndex = 0
	my.scope_repair.bufferPosition = 1

	//
	// Clear the configuration stack.
	//
	my.main_configuration_stack = NewConfigurationStack(my.prs)

	//
	// Keep parsing until we reach the end of file and succeed or
	// an error is encountered. The list of actions executed will
	// be stored in the "action" tuple.
	//
	my.locationStack[my.stateStackTop] = curtok
	my.actionStack[my.stateStackTop] = my.action.Size()
	var act int = my.TAction(my.stateStack[my.stateStackTop], current_kind)
	for {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if nil != my.monitor && my.monitor.IsCancelled() {
			return true
		}
		if act <= my.NUM_RULES {
			my.action.Add(act) // save my reduce action
			my.stateStackTop--

			for {
				my.stateStackTop -= (my.Rhs(act) - 1)
				act = my.NtAction(my.stateStack[my.stateStackTop], my.Lhs(act))
				if act <= my.NUM_RULES {
					continue
				} else {
					break
				}
			}
			my.stateStackTop += 1
			if my.stateStackTop >= len(my.stateStack) {
				my.ReallocateStacks()
			}
			my.stateStack[my.stateStackTop] = act

			my.locationStack[my.stateStackTop] = curtok
			my.actionStack[my.stateStackTop] = my.action.Size()
			act = my.TAction(act, current_kind)
			continue
		} else {
			if act == my.ERROR_ACTION {
				if curtok != error_token || my.main_configuration_stack.size() > 0 {
					var configuration = my.main_configuration_stack.Pop()
					if configuration == nil {
						act = my.ERROR_ACTION
					} else {
						my.stateStackTop = configuration.stack_top
						configuration.RetrieveStack(my.stateStack)
						act = configuration.act
						curtok = configuration.curtok
						my.action.ResetTo(configuration.action_length)
						current_kind = my.tokStream.GetKind(curtok)
						my.tokStream.ResetTo(my.tokStream.GetNext(curtok))
						continue
					}
				}
				break
			} else {
				if act > my.ACCEPT_ACTION && act < my.ERROR_ACTION {
					if my.main_configuration_stack.FindConfiguration(my.stateStack, my.stateStackTop, curtok) {
						act = my.ERROR_ACTION
					} else {
						my.main_configuration_stack.Push(my.stateStack, my.stateStackTop, act+1, curtok, my.action.Size())
						act = my.BaseAction(act)
					}
					continue
				} else {
					if act < my.ACCEPT_ACTION {
						my.action.Add(act) // save my shift action
						curtok = my.tokStream.GetToken()
						current_kind = my.tokStream.GetKind(curtok)
					} else {
						if act > my.ERROR_ACTION {
							my.action.Add(act) // save my shift-reduce action
							curtok = my.tokStream.GetToken()
							current_kind = my.tokStream.GetKind(curtok)
							act -= my.ERROR_ACTION
							for {
								my.stateStackTop -= (my.Rhs(act) - 1)
								act = my.NtAction(my.stateStack[my.stateStackTop], my.Lhs(act))
								if act <= my.NUM_RULES {
									continue
								} else {
									break
								}
							}
						} else {
							break // assert(act == ACCEPT_ACTION)  THIS IS NOT SUPPOSED TO HAPPEN!!!
						}
					}
					my.stateStackTop += 1
					if my.stateStackTop >= len(my.stateStack) {
						my.ReallocateStacks()
					}
					my.stateStack[my.stateStackTop] = act

					if curtok == error_token {
						my.scopeTrial(my.scope_repair, my.stateStack, my.stateStackTop)
						if my.scope_repair.distance >= MIN_DISTANCE {

							my.tokens.Add(start_token)
							var token int = first_stream_token
							for ; token != error_token; token = my.tokStream.GetNext(token) {
								my.tokens.Add(token)
							}
							my.AcceptRecovery(error_token)
							break
						}
					}
					my.locationStack[my.stateStackTop] = curtok
					my.actionStack[my.stateStackTop] = my.action.Size()
					act = my.TAction(act, current_kind)
				}
			}
		}
	}
	return act != my.ERROR_ACTION
}
func (my *RecoveryParser) cast() IPrsStream {
	t, _ := my.tokStream.(IPrsStream)
	return t
}
func (my *RecoveryParser) AcceptRecovery(error_token int) {
	//
	//
	//
	// int action_size = action.Size()

	//
	// Simulate parsing actions required for my sequence of scope
	// recoveries.
	// TODO need to add action and fix the location_stack?
	//
	var recovery_action = NewIntTuple()
	var k int = 0
	for ; k <= my.scopeStackTop; k++ {
		var scope_index int = my.scopeIndex[k]
		var la int = my.ScopeLa(scope_index)

		//
		// Compute the action (or Set of actions in case of conflicts) that
		// can be executed on the scope Lookahead symbol. Save the action(s)
		// in the action tuple.
		//
		recovery_action.Reset()
		var act int = my.TAction(my.stateStack[my.stateStackTop], la)
		if act > my.ACCEPT_ACTION && act < my.ERROR_ACTION { // conflicting actions?
			for {
				recovery_action.Add(my.BaseAction(act))
				act++
				if my.BaseAction(act) != 0 {
					continue
				} else {
					break
				}
			}
		} else {
			recovery_action.Add(act)
		}

		//
		// For each action defined on the scope Lookahead symbol,
		// try scope recovery. At least one action should succeed!
		//
		var start_action_size int = my.action.Size()
		var index int
		for index = 0; index < recovery_action.Size(); index++ {
			//
			// Reset the action tuple each time through my loop
			// to Clear previous actions that may have been added
			// because of a failed call to completeScope.
			//
			my.action.ResetTo(start_action_size)
			my.tokStream.ResetTo(error_token)
			my.tempStackTop = my.stateStackTop - 1
			var max_pos int = my.stateStackTop

			act = recovery_action.Get(index)
			for act <= my.NUM_RULES {
				my.action.Add(act) // save my reduce action
				//
				// ... Process all goto-reduce actions following
				// reduction, until a goto action is computed ...
				//
				for {
					var lhs_symbol int = my.Lhs(act)
					my.tempStackTop -= (my.Rhs(act) - 1)

					if my.tempStackTop > max_pos {
						act = my.tempStack[my.tempStackTop]
					} else {
						act = my.stateStack[my.tempStackTop]
					}

					act = my.NtAction(act, lhs_symbol)
					if act <= my.NUM_RULES {
						continue
					} else {
						break
					}
				}
				if my.tempStackTop+1 >= len(my.stateStack) {
					my.ReallocateStacks()
				}
				if !(max_pos < my.tempStackTop) {
					max_pos = my.tempStackTop
				}

				my.tempStack[my.tempStackTop+1] = act
				act = my.TAction(act, la)
			}

			//
			// If the Lookahead symbol is parsable, then we check
			// whether or not we have a match between the scope
			// prefix and the transition symbols corresponding to
			// the states on top of the stack.
			//
			if act != my.ERROR_ACTION {
				my.tempStackTop += 1
				my.nextStackTop = my.tempStackTop
				var i int = 0
				for ; i <= max_pos; i++ {
					my.nextStack[i] = my.stateStack[i]
				}

				//
				// NOTE that we do not need to update location_stack and
				// actionStack here because, once the rules associated with
				// these scopes are reduced, all these states will be popped
				// from the stack.
				//
				i = max_pos + 1
				for ; i <= my.tempStackTop; i++ {
					my.nextStack[i] = my.tempStack[i]
				}
				if my.CompleteScope(my.action, my.ScopeSuffix(scope_index)) {
					var i int = my.ScopeSuffix(my.scopeIndex[k])
					for ; my.ScopeRhs(i) != 0; i++ {

						my.tokens.Add(my.cast().MakeErrorToken(error_token,
							my.tokStream.GetPrevious(error_token),
							error_token, my.ScopeRhs(i)))
					}
					my.ReportError(my.scopeIndex[k], my.tokStream.GetPrevious(error_token))
					break
				}
			}
		}
		// assert (index < recovery_action.Size()) // sanity check!
		my.stateStackTop = my.nextStackTop
		Arraycopy(my.nextStack, 0, my.stateStack, 0, my.stateStackTop+1)
	}
	return
}
func (my *RecoveryParser) CompleteScope(action *IntSegmentedTuple, scope_rhs_index int) bool {
	var kind int = my.ScopeRhs(scope_rhs_index)
	if kind == 0 {
		return true
	}

	var act int = my.nextStack[my.nextStackTop]

	if kind > my.NT_OFFSET {
		var lhs_symbol int = kind - my.NT_OFFSET
		if my.BaseCheck(act+lhs_symbol) != lhs_symbol {
			// is there a valid
			// action defined on
			// lhs_symbol?
			return false
		}
		act = my.NtAction(act, lhs_symbol)

		//
		// if action is a goto-reduce action, save it as a shift-reduce
		// action.
		//
		var temp int
		if act <= my.NUM_RULES {
			temp = act + my.ERROR_ACTION
		} else {
			temp = act
		}
		action.Add(temp)
		for act <= my.NUM_RULES {
			my.nextStackTop -= (my.Rhs(act) - 1)
			act = my.NtAction(my.nextStack[my.nextStackTop], my.Lhs(act))
		}
		my.nextStackTop++
		my.nextStack[my.nextStackTop] = act
		return my.CompleteScope(action, scope_rhs_index+1)
	}

	//
	// Processing a terminal
	//
	act = my.TAction(act, kind)
	action.Add(act) // save my terminal action
	if act < my.ACCEPT_ACTION {
		my.nextStackTop++
		my.nextStack[my.nextStackTop] = act
		return my.CompleteScope(action, scope_rhs_index+1)
	} else {
		if act > my.ERROR_ACTION {
			act -= my.ERROR_ACTION
			for {
				my.nextStackTop -= (my.Rhs(act) - 1)
				act = my.NtAction(my.nextStack[my.nextStackTop], my.Lhs(act))
				if act <= my.NUM_RULES {
					continue
				} else {
					break
				}
			}
			my.nextStackTop++
			my.nextStack[my.nextStackTop] = act
			return true
		} else {
			if act > my.ACCEPT_ACTION && act < my.ERROR_ACTION { // conflicting actions?

				var save_action_size int = action.Size()
				var i int = act
				for ; my.BaseAction(i) != 0; i++ { // consider only shift and shift-reduce actions

					action.ResetTo(save_action_size)
					act = my.BaseAction(i)
					action.Add(act) // save my terminal action
					if act <= my.NUM_RULES {
					} else {
						if act < my.ACCEPT_ACTION {
							my.nextStackTop++
							my.nextStack[my.nextStackTop] = act
							if my.CompleteScope(action, scope_rhs_index+1) {
								return true
							}
						} else {
							if act > my.ERROR_ACTION {
								act -= my.ERROR_ACTION
								for {
									my.nextStackTop -= (my.Rhs(act) - 1)
									act = my.NtAction(my.nextStack[my.nextStackTop], my.Lhs(act))
									if act <= my.NUM_RULES {
										continue
									} else {
										break
									}
								}
								my.nextStackTop++
								my.nextStack[my.nextStackTop] = act
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}
