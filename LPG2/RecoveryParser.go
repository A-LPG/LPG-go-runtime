package lpg2


type RecoveryParser struct   {
	 *DiagnoseParser
     parser *BacktrackingParser
     action *IntSegmentedTuple
     tokens *IntTuple
     actionStack []int
     scope_repair *PrimaryRepairInfo
 }
//
// maxErrors is the maximum int of errors to be repaired
// maxTime is the maximum amount of time allowed for diagnosing
// but at least one error must be diagnosed 
//
func  NewRecoveryParser(parser *BacktrackingParser, action *IntSegmentedTuple, tokens *IntTuple, tokStream IPrsStream,
		prs ParseTable, maxErrors int, maxTime int,monitor Monitor ) *RecoveryParser{
	    t := new(RecoveryParser)
	    t.DiagnoseParser=NewDiagnoseParser(tokStream,prs,maxErrors,maxTime,monitor)
		t.parser = parser
		t.action = action
		t.tokens = tokens
		return t
}

func (this *RecoveryParser) reallocateStacks() {
	this.DiagnoseParser.reallocateStacks()
	if len(this.actionStack) == 0 {
		this.actionStack = make([]int,len(this.stateStack))
	} else {
		var old_stack_length int = len(this.actionStack)
		this.actionStack =arraycopy(this.actionStack, 0,  make([]int,len(this.stateStack)), 0, old_stack_length)
	}
	return
}
func (this *RecoveryParser) reportError(scope_index int, error_token int) {
	var text string = "\""
	var i int = this.scopeSuffix(scope_index)
	for ; this.scopeRhs(i) != 0 ;i++{
		if !this.isNullable(this.scopeRhs(i)) {

			var symbol_index int
			if this.scopeRhs(i) > this.NT_OFFSET{
				symbol_index = this.nonterminalIndex(this.scopeRhs(i) - this.NT_OFFSET)
			}else{
				symbol_index =this.terminalIndex(this.scopeRhs(i))
			}

			if len(this.name(symbol_index)) > 0 {
				if len(text) > 1 { // Not just starting quote?
					text += " "// add a space separator
				}
				text += this.name(symbol_index)
			}
		}
	}
	text += "\""
	this.tokStream.reportError(SCOPE_CODE, error_token, error_token,[]string{text},0)
	return
}
func (this *RecoveryParser) recover(marker_token int, error_token int) (int,error) {
	if len(this.stateStack) == 0 {
		this.reallocateStacks()
	}

	this.tokens.reset()
	this.tokStream.reset()
	this.tokens.add(this.tokStream.getPrevious(this.tokStream.peek()))
	var restart_token int
	if marker_token != 0 {
		restart_token=marker_token
	}else{
		restart_token= this.tokStream.getToken()
	}

	var	old_action_size int = 0
	this.stateStackTop = 0
	this.stateStack[this.stateStackTop] = this.START_STATE
	for;; {
		this.action.resetTo(old_action_size)
		if !this.fixError(restart_token, error_token) {
			return -1, NewBadParseException(error_token)
		}
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if nil == this.monitor  && this.monitor.isCancelled() {
			break
		}
		//
		// At this stage, we have a recovery configuration. See how
		// far we can go with it.
		//
		restart_token = error_token
		this.tokStream.resetTo(error_token)
		old_action_size = this.action.size()// save the old size in case we encounter a new error
		error_token = this.parser.backtrackParse(this.stateStack, this.stateStackTop, this.action, 0)
		this.tokStream.resetTo(this.tokStream.getNext(restart_token))
		if error_token != 0 {
			continue
		}else{
			break
		}
	} // no error found
	return restart_token,nil
}
//
// Given the configuration consisting of the states in stateStack
// and the sequence of tokens (current_kind, followed by the tokens
// in tokStream), fixError parses up to error_token in the tokStream
// recovers, if possible, from that error and returns the result.
// While doing this, it also computes the location_stack information
// and the sequence of actions that matches up with the result that
// it returns.
//
func (this *RecoveryParser) fixError(start_token int, error_token int) bool {
	//
	// Save information about the current configuration.
	//
	var curtok int = start_token
	var	current_kind int = this.tokStream.getKind(curtok)
	var	first_stream_token int = this.tokStream.peek()

	this.buffer[1] = error_token
	this.buffer[0] = this.tokStream.getPrevious(this.buffer[1])
	var k int = 2
	for ;k < BUFF_SIZE; k++ {
		this.buffer[k] = this.tokStream.getNext(this.buffer[k - 1])
	}

	this.scope_repair.distance = 0
	this.scope_repair.misspellIndex = 0
	this.scope_repair.bufferPosition = 1

	//
	// Clear the configuration stack.
	//
	this.main_configuration_stack = NewConfigurationStack(this.prs)

	//
	// Keep parsing until we reach the end of file and succeed or
	// an error is encountered. The list of actions executed will
	// be stored in the "action" tuple.
	//
	this.locationStack[this.stateStackTop] = curtok
	this.actionStack[this.stateStackTop] = this.action.size()
	var act int = this.tAction(this.stateStack[this.stateStackTop], current_kind)
	for ;; {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if nil == this.monitor  && this.monitor.isCancelled() {
			return true
		}
		if act <= this.NUM_RULES {
			this.action.add(act)// save this reduce action
			this.stateStackTop--

			for;; {
				this.stateStackTop -= (this.rhs(act) - 1)
				act = this.ntAction(this.stateStack[this.stateStackTop], this.lhs(act))
				if act <= this.NUM_RULES {
					continue
				}else {
					break
				}
			}
			this.stateStackTop+=1
			if this.stateStackTop >= len(this.stateStack) {
				this.reallocateStacks()
			}
			this.stateStack[this.stateStackTop] = act
		
			this.locationStack[this.stateStackTop] = curtok
			this.actionStack[this.stateStackTop] = this.action.size()
			act = this.tAction(act, current_kind)
			continue
		}else{
			if act == this.ERROR_ACTION {
				if curtok != error_token || this.main_configuration_stack.size() > 0 {
					var configuration = this.main_configuration_stack.pop()
					if configuration == nil {
						act = this.ERROR_ACTION
					} else {
						this.stateStackTop = configuration.stack_top
						configuration.retrieveStack(this.stateStack)
						act = configuration.act
						curtok = configuration.curtok
						this.action.resetTo(configuration.action_length)
						current_kind = this.tokStream.getKind(curtok)
						this.tokStream.resetTo(this.tokStream.getNext(curtok))
						continue
					}
				}
				break
			}else {
				if act > this.ACCEPT_ACTION && act < this.ERROR_ACTION {
					if this.main_configuration_stack.findConfiguration(this.stateStack, this.stateStackTop, curtok) {
						act = this.ERROR_ACTION
					} else {
						this.main_configuration_stack.push(this.stateStack, this.stateStackTop, act + 1, curtok, this.action.size())
						act = this.baseAction(act)
					}
					continue
				} else{
					if act < this.ACCEPT_ACTION {
						this.action.add(act)// save this shift action
						curtok = this.tokStream.getToken()
						current_kind = this.tokStream.getKind(curtok)
					}else{ 
						if act > this.ERROR_ACTION {
							this.action.add(act)// save this shift-reduce action
							curtok = this.tokStream.getToken()
							current_kind = this.tokStream.getKind(curtok)
							act -= this.ERROR_ACTION
							for;; {
								this.stateStackTop -= (this.rhs(act) - 1)
								act = this.ntAction(this.stateStack[this.stateStackTop], this.lhs(act))
								if act <= this.NUM_RULES {
									continue
								}else {
									break
								}
							}
						} else{
							break// assert(act == ACCEPT_ACTION)  THIS IS NOT SUPPOSED TO HAPPEN!!!
						}
					}
					this.stateStackTop +=1
					if this.stateStackTop >= len(this.stateStack) {
						this.reallocateStacks()
					}
					this.stateStack[this.stateStackTop] = act
			
					if curtok == error_token {
						this.scopeTrial(this.scope_repair, this.stateStack, this.stateStackTop)
						if this.scope_repair.distance >= MIN_DISTANCE {

							this.tokens.add(start_token)
							var token int = first_stream_token
							for ; token != error_token; token = this.tokStream.getNext(token) {
								this.tokens.add(token)
							}
							this.acceptRecovery(error_token)
							break
						}
					}
					this.locationStack[this.stateStackTop] = curtok
					this.actionStack[this.stateStackTop] = this.action.size()
					act = this.tAction(act, current_kind)
				}
			}		
		}
	}
	return act != this.ERROR_ACTION
}
func (this *RecoveryParser) cast() IPrsStream{
	t, _ := this.tokStream.(IPrsStream)
	return  t
}
func (this *RecoveryParser) acceptRecovery(error_token int) {
	//
	//
	//
	// int action_size = action.size()

	//
	// Simulate parsing actions required for this sequence of scope
	// recoveries.
	// TODO need to add action and fix the location_stack?
	//
	var recovery_action  = NewIntTuple()
	var k int = 0
	for  ;k <= this.scopeStackTop ;k++ {
		var scope_index int = this.scopeIndex[k]
		var la int = this.scopeLa(scope_index)

		//
		// Compute the action (or set of actions in case of conflicts) that
		// can be executed on the scope lookahead symbol. Save the action(s)
		// in the action tuple.
		//
		recovery_action.reset()
		var act int = this.tAction(this.stateStack[this.stateStackTop], la)
		if act > this.ACCEPT_ACTION && act < this.ERROR_ACTION { // conflicting actions?
			for;; {
				recovery_action.add(this.baseAction(act))
				act++
				if this.baseAction(act) != 0 {
					continue
				}else {
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
		var start_action_size int = this.action.size()
		var index int
		for index = 0; index < recovery_action.size(); index++{
			//
			// Reset the action tuple each time through this loop
			// to clear previous actions that may have been added
			// because of a failed call to completeScope.
			//
			this.action.resetTo(start_action_size)
			this.tokStream.resetTo(error_token)
			this.tempStackTop = this.stateStackTop - 1
			var max_pos int = this.stateStackTop

			act = recovery_action.get(index)
			for; act <= this.NUM_RULES; {
				this.action.add(act)// save this reduce action
				//
				// ... Process all goto-reduce actions following
				// reduction, until a goto action is computed ...
				//
				for;; {
					var lhs_symbol int = this.lhs(act)
					this.tempStackTop -= (this.rhs(act) - 1)

						if this.tempStackTop > max_pos{
							act =this.tempStack[this.tempStackTop]
						}else{
							act =this.stateStack[this.tempStackTop]
						}

					act = this.ntAction(act, lhs_symbol)
					if act <= this.NUM_RULES {
						continue
					}else {
						break
					}
				}
				if this.tempStackTop + 1 >= len(this.stateStack) {
					this.reallocateStacks()
				}
				if !(max_pos < this.tempStackTop){
					max_pos= this.tempStackTop
				}

				this.tempStack[this.tempStackTop + 1] = act
				act = this.tAction(act, la)
			}


			//
			// If the lookahead symbol is parsable, then we check
			// whether or not we have a match between the scope
			// prefix and the transition symbols corresponding to
			// the states on top of the stack.
			//
			if (act != this.ERROR_ACTION) {
				this.tempStackTop+=1
				this.nextStackTop = this.tempStackTop
				var i int = 0
				for  ;i <= max_pos ;i++ {
					this.nextStack[i] = this.stateStack[i]
				}

				//
				// NOTE that we do not need to update location_stack and
				// actionStack here because, once the rules associated with
				// these scopes are reduced, all these states will be popped
				// from the stack.
				//
				i  = max_pos + 1
				for ; i <= this.tempStackTop; i++ {
					this.nextStack[i] = this.tempStack[i]
				}
				if this.completeScope(this.action, this.scopeSuffix(scope_index)) {
					var i int = this.scopeSuffix(this.scopeIndex[k])
					for ; this.scopeRhs(i) != 0 ;i++{

						this.tokens.add( this.cast().makeErrorToken(error_token,
								this.tokStream.getPrevious(error_token),
								error_token, this.scopeRhs(i)))
					}
					this.reportError(this.scopeIndex[k], this.tokStream.getPrevious(error_token))
					break
				}
			}
		}
		// assert (index < recovery_action.size()) // sanity check!
		this.stateStackTop = this.nextStackTop
		arraycopy(this.nextStack, 0, this.stateStack, 0, this.stateStackTop + 1)
	}
	return
}
func (this *RecoveryParser) completeScope(action *IntSegmentedTuple, scope_rhs_index int) bool {
	var kind int = this.scopeRhs(scope_rhs_index)
	if (kind == 0) {
		return true
	}

	var act int = this.nextStack[this.nextStackTop]

	if kind > this.NT_OFFSET {
		var lhs_symbol int = kind - this.NT_OFFSET
		if this.baseCheck(act + lhs_symbol) != lhs_symbol {
			// is there a valid
			// action defined on
			// lhs_symbol?
			return false
		}
		act = this.ntAction(act, lhs_symbol)

		//
		// if action is a goto-reduce action, save it as a shift-reduce
		// action.
		//
		var temp int
		if act <= this.NUM_RULES {
			temp = act + this.ERROR_ACTION
		}else {
			temp = act
		}
		action.add(temp)
		for;act <= this.NUM_RULES; {
			this.nextStackTop -= (this.rhs(act) - 1)
			act = this.ntAction(this.nextStack[this.nextStackTop], this.lhs(act))
		}
		this.nextStackTop++
		this.nextStack[this.nextStackTop] = act
		return this.completeScope(action, scope_rhs_index + 1)
	}

	//
	// Processing a terminal
	//
	act = this.tAction(act, kind)
	action.add(act)// save this terminal action
	if act < this.ACCEPT_ACTION {
		this.nextStackTop++
		this.nextStack[this.nextStackTop] = act
		return this.completeScope(action, scope_rhs_index + 1)
	}else {
		if act > this.ERROR_ACTION{
				act -= this.ERROR_ACTION
				for;; {
					this.nextStackTop -= (this.rhs(act) - 1)
					act = this.ntAction(this.nextStack[this.nextStackTop], this.lhs(act))
					if act <= this.NUM_RULES{
						continue
					}else {
						break
					}
				}
				this.nextStackTop++
				this.nextStack[this.nextStackTop] = act
				return true
		}else {
			if act > this.ACCEPT_ACTION && act < this.ERROR_ACTION { // conflicting actions?

				var save_action_size int = action.size()
				var i int = act
				for ; this.baseAction(i) != 0 ;i++{// consider only shift and shift-reduce actions
				
					action.resetTo(save_action_size)
					act = this.baseAction(i)
					action.add(act)// save this terminal action
					if act <= this.NUM_RULES {
					} else {
						if act < this.ACCEPT_ACTION {
							this.nextStackTop++
							this.nextStack[this.nextStackTop] = act
							if this.completeScope(action, scope_rhs_index+1) {
								return true
							}
						} else {
							if act > this.ERROR_ACTION {
								act -= this.ERROR_ACTION
								for ; ; {
									this.nextStackTop -= (this.rhs(act) - 1)
									act = this.ntAction(this.nextStack[this.nextStackTop], this.lhs(act))
									if act <= this.NUM_RULES {
										continue
									} else {
										break
									}
								}
								this.nextStackTop++
								this.nextStack[this.nextStackTop] = act
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


