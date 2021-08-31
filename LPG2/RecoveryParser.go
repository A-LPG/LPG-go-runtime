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
	t.DiagnoseParser = NewDiagnoseParser(tokStream, prs, maxErrors, maxTime, monitor)
	t.parser = parser
	t.action = action
	t.tokens = tokens
	return t
}

func (self *RecoveryParser) reallocateStacks() {
	self.DiagnoseParser.reallocateStacks()
	if len(self.actionStack) == 0 {
		self.actionStack = make([]int, len(self.stateStack))
	} else {
		var old_stack_length int = len(self.actionStack)
		self.actionStack = arraycopy(self.actionStack, 0, make([]int, len(self.stateStack)), 0, old_stack_length)
	}
	return
}
func (self *RecoveryParser) reportError(scope_index int, error_token int) {
	var text string = "\""
	var i int = self.scopeSuffix(scope_index)
	for ; self.scopeRhs(i) != 0; i++ {
		if !self.isNullable(self.scopeRhs(i)) {

			var symbol_index int
			if self.scopeRhs(i) > self.NT_OFFSET {
				symbol_index = self.nonterminalIndex(self.scopeRhs(i) - self.NT_OFFSET)
			} else {
				symbol_index = self.terminalIndex(self.scopeRhs(i))
			}

			if len(self.name(symbol_index)) > 0 {
				if len(text) > 1 { // Not just starting quote?
					text += " " // add a space separator
				}
				text += self.name(symbol_index)
			}
		}
	}
	text += "\""
	self.tokStream.reportError(SCOPE_CODE, error_token, error_token, []string{text}, 0)
	return
}
func (self *RecoveryParser) recover(marker_token int, error_token int) (int, error) {
	if len(self.stateStack) == 0 {
		self.reallocateStacks()
	}

	self.tokens.reset()
	self.tokStream.reset()
	self.tokens.add(self.tokStream.getPrevious(self.tokStream.peek()))
	var restart_token int
	if marker_token != 0 {
		restart_token = marker_token
	} else {
		restart_token = self.tokStream.getToken()
	}

	var old_action_size int = 0
	self.stateStackTop = 0
	self.stateStack[self.stateStackTop] = self.START_STATE
	for {
		self.action.resetTo(old_action_size)
		if !self.fixError(restart_token, error_token) {
			return -1, NewBadParseException(error_token)
		}
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if nil != self.monitor && self.monitor.isCancelled() {
			break
		}
		//
		// At self stage, we have a recovery configuration. See how
		// far we can go with it.
		//
		restart_token = error_token
		self.tokStream.resetTo(error_token)
		old_action_size = self.action.size() // save the old size in case we encounter a new error
		error_token = self.parser.backtrackParse(self.stateStack, self.stateStackTop, self.action, 0)
		self.tokStream.resetTo(self.tokStream.getNext(restart_token))
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
// in tokStream), fixError parses up to error_token in the tokStream
// recovers, if possible, from that error and returns the result.
// While doing self, it also computes the location_stack information
// and the sequence of actions that matches up with the result that
// it returns.
//
func (self *RecoveryParser) fixError(start_token int, error_token int) bool {
	//
	// Save information about the current configuration.
	//
	var curtok int = start_token
	var current_kind int = self.tokStream.getKind(curtok)
	var first_stream_token int = self.tokStream.peek()

	self.buffer[1] = error_token
	self.buffer[0] = self.tokStream.getPrevious(self.buffer[1])
	var k int = 2
	for ; k < BUFF_SIZE; k++ {
		self.buffer[k] = self.tokStream.getNext(self.buffer[k-1])
	}

	self.scope_repair.distance = 0
	self.scope_repair.misspellIndex = 0
	self.scope_repair.bufferPosition = 1

	//
	// Clear the configuration stack.
	//
	self.main_configuration_stack = NewConfigurationStack(self.prs)

	//
	// Keep parsing until we reach the end of file and succeed or
	// an error is encountered. The list of actions executed will
	// be stored in the "action" tuple.
	//
	self.locationStack[self.stateStackTop] = curtok
	self.actionStack[self.stateStackTop] = self.action.size()
	var act int = self.tAction(self.stateStack[self.stateStackTop], current_kind)
	for {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if nil != self.monitor && self.monitor.isCancelled() {
			return true
		}
		if act <= self.NUM_RULES {
			self.action.add(act) // save self reduce action
			self.stateStackTop--

			for {
				self.stateStackTop -= (self.rhs(act) - 1)
				act = self.ntAction(self.stateStack[self.stateStackTop], self.lhs(act))
				if act <= self.NUM_RULES {
					continue
				} else {
					break
				}
			}
			self.stateStackTop += 1
			if self.stateStackTop >= len(self.stateStack) {
				self.reallocateStacks()
			}
			self.stateStack[self.stateStackTop] = act

			self.locationStack[self.stateStackTop] = curtok
			self.actionStack[self.stateStackTop] = self.action.size()
			act = self.tAction(act, current_kind)
			continue
		} else {
			if act == self.ERROR_ACTION {
				if curtok != error_token || self.main_configuration_stack.size() > 0 {
					var configuration = self.main_configuration_stack.pop()
					if configuration == nil {
						act = self.ERROR_ACTION
					} else {
						self.stateStackTop = configuration.stack_top
						configuration.retrieveStack(self.stateStack)
						act = configuration.act
						curtok = configuration.curtok
						self.action.resetTo(configuration.action_length)
						current_kind = self.tokStream.getKind(curtok)
						self.tokStream.resetTo(self.tokStream.getNext(curtok))
						continue
					}
				}
				break
			} else {
				if act > self.ACCEPT_ACTION && act < self.ERROR_ACTION {
					if self.main_configuration_stack.findConfiguration(self.stateStack, self.stateStackTop, curtok) {
						act = self.ERROR_ACTION
					} else {
						self.main_configuration_stack.push(self.stateStack, self.stateStackTop, act+1, curtok, self.action.size())
						act = self.baseAction(act)
					}
					continue
				} else {
					if act < self.ACCEPT_ACTION {
						self.action.add(act) // save self shift action
						curtok = self.tokStream.getToken()
						current_kind = self.tokStream.getKind(curtok)
					} else {
						if act > self.ERROR_ACTION {
							self.action.add(act) // save self shift-reduce action
							curtok = self.tokStream.getToken()
							current_kind = self.tokStream.getKind(curtok)
							act -= self.ERROR_ACTION
							for {
								self.stateStackTop -= (self.rhs(act) - 1)
								act = self.ntAction(self.stateStack[self.stateStackTop], self.lhs(act))
								if act <= self.NUM_RULES {
									continue
								} else {
									break
								}
							}
						} else {
							break // assert(act == ACCEPT_ACTION)  THIS IS NOT SUPPOSED TO HAPPEN!!!
						}
					}
					self.stateStackTop += 1
					if self.stateStackTop >= len(self.stateStack) {
						self.reallocateStacks()
					}
					self.stateStack[self.stateStackTop] = act

					if curtok == error_token {
						self.scopeTrial(self.scope_repair, self.stateStack, self.stateStackTop)
						if self.scope_repair.distance >= MIN_DISTANCE {

							self.tokens.add(start_token)
							var token int = first_stream_token
							for ; token != error_token; token = self.tokStream.getNext(token) {
								self.tokens.add(token)
							}
							self.acceptRecovery(error_token)
							break
						}
					}
					self.locationStack[self.stateStackTop] = curtok
					self.actionStack[self.stateStackTop] = self.action.size()
					act = self.tAction(act, current_kind)
				}
			}
		}
	}
	return act != self.ERROR_ACTION
}
func (self *RecoveryParser) cast() IPrsStream {
	t, _ := self.tokStream.(IPrsStream)
	return t
}
func (self *RecoveryParser) acceptRecovery(error_token int) {
	//
	//
	//
	// int action_size = action.size()

	//
	// Simulate parsing actions required for self sequence of scope
	// recoveries.
	// TODO need to add action and fix the location_stack?
	//
	var recovery_action = NewIntTuple()
	var k int = 0
	for ; k <= self.scopeStackTop; k++ {
		var scope_index int = self.scopeIndex[k]
		var la int = self.scopeLa(scope_index)

		//
		// Compute the action (or set of actions in case of conflicts) that
		// can be executed on the scope lookahead symbol. Save the action(s)
		// in the action tuple.
		//
		recovery_action.reset()
		var act int = self.tAction(self.stateStack[self.stateStackTop], la)
		if act > self.ACCEPT_ACTION && act < self.ERROR_ACTION { // conflicting actions?
			for {
				recovery_action.add(self.baseAction(act))
				act++
				if self.baseAction(act) != 0 {
					continue
				} else {
					break
				}
			}
		} else {
			recovery_action.add(act)
		}

		//
		// For each action defined on the scope lookahead symbol,
		// try scope recovery. At least one action should succeed!
		//
		var start_action_size int = self.action.size()
		var index int
		for index = 0; index < recovery_action.size(); index++ {
			//
			// Reset the action tuple each time through self loop
			// to clear previous actions that may have been added
			// because of a failed call to completeScope.
			//
			self.action.resetTo(start_action_size)
			self.tokStream.resetTo(error_token)
			self.tempStackTop = self.stateStackTop - 1
			var max_pos int = self.stateStackTop

			act = recovery_action.get(index)
			for act <= self.NUM_RULES {
				self.action.add(act) // save self reduce action
				//
				// ... Process all goto-reduce actions following
				// reduction, until a goto action is computed ...
				//
				for {
					var lhs_symbol int = self.lhs(act)
					self.tempStackTop -= (self.rhs(act) - 1)

					if self.tempStackTop > max_pos {
						act = self.tempStack[self.tempStackTop]
					} else {
						act = self.stateStack[self.tempStackTop]
					}

					act = self.ntAction(act, lhs_symbol)
					if act <= self.NUM_RULES {
						continue
					} else {
						break
					}
				}
				if self.tempStackTop+1 >= len(self.stateStack) {
					self.reallocateStacks()
				}
				if !(max_pos < self.tempStackTop) {
					max_pos = self.tempStackTop
				}

				self.tempStack[self.tempStackTop+1] = act
				act = self.tAction(act, la)
			}

			//
			// If the lookahead symbol is parsable, then we check
			// whether or not we have a match between the scope
			// prefix and the transition symbols corresponding to
			// the states on top of the stack.
			//
			if act != self.ERROR_ACTION {
				self.tempStackTop += 1
				self.nextStackTop = self.tempStackTop
				var i int = 0
				for ; i <= max_pos; i++ {
					self.nextStack[i] = self.stateStack[i]
				}

				//
				// NOTE that we do not need to update location_stack and
				// actionStack here because, once the rules associated with
				// these scopes are reduced, all these states will be popped
				// from the stack.
				//
				i = max_pos + 1
				for ; i <= self.tempStackTop; i++ {
					self.nextStack[i] = self.tempStack[i]
				}
				if self.completeScope(self.action, self.scopeSuffix(scope_index)) {
					var i int = self.scopeSuffix(self.scopeIndex[k])
					for ; self.scopeRhs(i) != 0; i++ {

						self.tokens.add(self.cast().makeErrorToken(error_token,
							self.tokStream.getPrevious(error_token),
							error_token, self.scopeRhs(i)))
					}
					self.reportError(self.scopeIndex[k], self.tokStream.getPrevious(error_token))
					break
				}
			}
		}
		// assert (index < recovery_action.size()) // sanity check!
		self.stateStackTop = self.nextStackTop
		arraycopy(self.nextStack, 0, self.stateStack, 0, self.stateStackTop+1)
	}
	return
}
func (self *RecoveryParser) completeScope(action *IntSegmentedTuple, scope_rhs_index int) bool {
	var kind int = self.scopeRhs(scope_rhs_index)
	if kind == 0 {
		return true
	}

	var act int = self.nextStack[self.nextStackTop]

	if kind > self.NT_OFFSET {
		var lhs_symbol int = kind - self.NT_OFFSET
		if self.baseCheck(act+lhs_symbol) != lhs_symbol {
			// is there a valid
			// action defined on
			// lhs_symbol?
			return false
		}
		act = self.ntAction(act, lhs_symbol)

		//
		// if action is a goto-reduce action, save it as a shift-reduce
		// action.
		//
		var temp int
		if act <= self.NUM_RULES {
			temp = act + self.ERROR_ACTION
		} else {
			temp = act
		}
		action.add(temp)
		for act <= self.NUM_RULES {
			self.nextStackTop -= (self.rhs(act) - 1)
			act = self.ntAction(self.nextStack[self.nextStackTop], self.lhs(act))
		}
		self.nextStackTop++
		self.nextStack[self.nextStackTop] = act
		return self.completeScope(action, scope_rhs_index+1)
	}

	//
	// Processing a terminal
	//
	act = self.tAction(act, kind)
	action.add(act) // save self terminal action
	if act < self.ACCEPT_ACTION {
		self.nextStackTop++
		self.nextStack[self.nextStackTop] = act
		return self.completeScope(action, scope_rhs_index+1)
	} else {
		if act > self.ERROR_ACTION {
			act -= self.ERROR_ACTION
			for {
				self.nextStackTop -= (self.rhs(act) - 1)
				act = self.ntAction(self.nextStack[self.nextStackTop], self.lhs(act))
				if act <= self.NUM_RULES {
					continue
				} else {
					break
				}
			}
			self.nextStackTop++
			self.nextStack[self.nextStackTop] = act
			return true
		} else {
			if act > self.ACCEPT_ACTION && act < self.ERROR_ACTION { // conflicting actions?

				var save_action_size int = action.size()
				var i int = act
				for ; self.baseAction(i) != 0; i++ { // consider only shift and shift-reduce actions

					action.resetTo(save_action_size)
					act = self.baseAction(i)
					action.add(act) // save self terminal action
					if act <= self.NUM_RULES {
					} else {
						if act < self.ACCEPT_ACTION {
							self.nextStackTop++
							self.nextStack[self.nextStackTop] = act
							if self.completeScope(action, scope_rhs_index+1) {
								return true
							}
						} else {
							if act > self.ERROR_ACTION {
								act -= self.ERROR_ACTION
								for {
									self.nextStackTop -= (self.rhs(act) - 1)
									act = self.ntAction(self.nextStack[self.nextStackTop], self.lhs(act))
									if act <= self.NUM_RULES {
										continue
									} else {
										break
									}
								}
								self.nextStackTop++
								self.nextStack[self.nextStackTop] = act
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
