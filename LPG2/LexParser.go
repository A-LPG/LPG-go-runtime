package lpg2

type LexParser struct {
	taking_actions  bool
	STACK_INCREMENT int
	START_STATE     int
	LA_STATE_OFFSET int
	EOFT_SYMBOL     int
	ACCEPT_ACTION   int
	ERROR_ACTION    int
	START_SYMBOL    int
	NUM_RULES       int

	tokStream     ILexStream
	prs           ParseTable
	ra            RuleAction
	action        *IntTuple
	stateStackTop int
	stackLength   int
	stack         []int
	locationStack []int
	tempStack     []int
	lastToken     int
	currentAction int
	curtok        int
	starttok      int
	current_kind  int
}

func NewLexParser(tokStream ILexStream, prs ParseTable, ra RuleAction) *LexParser {
	self := new(LexParser)
	self.STACK_INCREMENT = 1024
    self.stackLength=0
	if tokStream != nil && prs != nil && ra != nil {
		self.reset(tokStream, prs, ra)
	}

	return self
}

func (self *LexParser) reset(tokStream ILexStream, prs ParseTable, ra RuleAction) {
	self.tokStream = tokStream
	self.prs = prs
	self.ra = ra
	self.START_STATE = prs.getStartState()
	self.LA_STATE_OFFSET = prs.getLaStateOffset()
	self.EOFT_SYMBOL = prs.getEoftSymbol()
	self.ACCEPT_ACTION = prs.getAcceptAction()
	self.ERROR_ACTION = prs.getErrorAction()
	self.START_SYMBOL = prs.getStartSymbol()
	self.NUM_RULES = prs.getNumRules()
}

//
// Stacks portion
//

func (self *LexParser) reallocateStacks() {
	var old_stack_length int
	if len(self.stack) == 0 {
		old_stack_length = 0
	} else {
		old_stack_length = self.stackLength
	}
	self.stackLength += self.STACK_INCREMENT
	if old_stack_length == 0 {
		self.stack = make([]int, self.stackLength)
		self.locationStack = make([]int, self.stackLength)
		self.tempStack = make([]int, self.stackLength)
	} else {
		self.stack = arraycopy(self.stack, 0, make([]int, self.stackLength), 0, old_stack_length)
		self.locationStack = arraycopy(self.locationStack, 0, make([]int, self.stackLength), 0, old_stack_length)
		self.tempStack = arraycopy(self.tempStack, 0, make([]int, self.stackLength), 0, old_stack_length)
	}
	return
}

//
// The following functions can be invoked only when the parser is
// processing actions. Thus, they can be invoked when the parser
// was entered via the main entry point (parseCharacters()). When using
// the incremental parser (via the entry point scanNextToken(int [], int)),
// they always return 0 when invoked. // TODO Should we throw an Exception instead?
// However, note that when parseActions() is invoked after successfully
// parsing an input with the incremental parser, then they can be invoked.
//
func (self *LexParser) getFirstTokenAt(i int) (int, error) {
	return self.getToken(i)
}
func (self *LexParser) getFirstToken() int {
	return self.starttok
}
func (self *LexParser) getLastToken() int {
	return self.lastToken
}
func (self *LexParser) getLastTokenAt(i int) (int, error) {

	if self.taking_actions {
		if i >= self.prs.rhs(self.currentAction) {
			return self.lastToken, nil
		} else {
			var index, e = self.getToken(i + 1)
			if e != nil {
				return -1, e
			}
			return self.tokStream.getPrevious(index), nil
		}

	}
	return -1, NewUnavailableParserInformationException("")
}
func (self *LexParser) getCurrentRule() (int, error) {
	if self.taking_actions {
		return self.currentAction, nil
	}
	return -1, NewUnavailableParserInformationException("")
}

//
// Given a rule of the form     A = x1 x2 ... xn     n > 0
//
// the function getToken(i) yields the symbol xi, if xi is a terminal
// or ti, if xi is a nonterminal that produced a string of the form
// xi => ti w. If xi is a nullable nonterminal, then ti is the first
//  symbol that immediately follows xi in the input (the lookahead).
//
func (self *LexParser) getToken(i int) (int, error) {
	if self.taking_actions {
		return self.locationStack[self.stateStackTop+(i-1)], nil
	}
	return -1, NewUnavailableParserInformationException("")
}
func (self *LexParser) setSym1(i int) {}
func (self *LexParser) getSym(i int) int {
	i, _ = self.getLastTokenAt(i)
	return i
}

func (self *LexParser) resetTokenStream(i int) {
	//
	// if i exceeds the upper bound, reset it to point to the last element.
	//
	var temp int
	if i > self.tokStream.getStreamLength() {
		temp = self.tokStream.getStreamLength()
	} else {
		temp = i
	}
	self.tokStream.resetTo(temp)
	self.curtok = self.tokStream.getToken()
	self.current_kind = self.tokStream.getKind(self.curtok)
	if len(self.stack) == 0 {
		self.reallocateStacks()
	}
	if self.action == nil {
		self.action = NewIntTupleWithEstimate(1 << 10)
	}
}

//
// Parse the input and create a stream of tokens.
//
func (self *LexParser) parseCharacters(start_offset int, end_offset int, monitor Monitor) {
	self.resetTokenStream(start_offset)
	for ; self.curtok <= end_offset; {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if monitor != nil && monitor.isCancelled() {
			return
		}
		self.lexNextToken(end_offset)
	}
}

//
// Parse the input and create a stream of tokens.
//
func (self *LexParser) parseCharactersWhitMonitor(monitor Monitor) {
	//
	// Indicate that we are running the regular parser and that it's
	// ok to use the utility functions to query the parser.
	//
	self.taking_actions = true
	self.resetTokenStream(0)
	self.lastToken = self.tokStream.getPrevious(self.curtok)

	//
	// Until it reaches the end-of-file token, self outer loop
	// resets the parser and processes the next token.
	//
ProcessTokens:
	for ; self.current_kind != self.EOFT_SYMBOL; {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if monitor != nil && monitor.isCancelled() {
			break ProcessTokens
		}

		self.stateStackTop = -1
		self.currentAction = self.START_STATE
		self.starttok = self.curtok

	ScanToken:
		for ; ; {
			self.stateStackTop += 1
			if self.stateStackTop >= len(self.stack) {
				self.reallocateStacks()
			}
			self.stack[self.stateStackTop] = self.currentAction

			self.locationStack[self.stateStackTop] = self.curtok

			//
			// Compute the action on the next character. If it is a reduce action, we do not
			// want to accept it until we are sure that the character in question is can be parsed.
			// What we are trying to avoid is a situation where Curtok is not the EOF token
			// but it yields a default reduce action in the current configuration even though
			// it cannot ultimately be shifted However, the state on top of the configuration also
			// contains a valid reduce action on EOF which, if taken, would lead to the successful
			// scanning of the token.
			//
			// Thus, if the character can be parsed, we proceed normally. Otherwise, we proceed
			// as if we had reached the end of the file (end of the token, since we are really
			// scanning).
			//
			self.parseNextCharacter(self.curtok, self.current_kind)
			if self.currentAction == self.ERROR_ACTION && self.current_kind != self.EOFT_SYMBOL { // if not successful try EOF

				var save_next_token = self.tokStream.peek()                  // save position after curtok
				self.tokStream.resetTo(self.tokStream.getStreamLength() - 1) // point to the end of the input
				self.parseNextCharacter(self.curtok, self.EOFT_SYMBOL)
				// assert (currentAction == ACCEPT_ACTION || currentAction == ERROR_ACTION)
				self.tokStream.resetTo(save_next_token) // reset the stream for the next token after curtok.
			}

			//
			// At self point, currentAction is either a Shift, Shift-Reduce, Accept or Error action.
			//
			if self.currentAction > self.ERROR_ACTION { // Shift-reduce

				self.lastToken = self.curtok
				self.curtok = self.tokStream.getToken()
				self.current_kind = self.tokStream.getKind(self.curtok)
				self.currentAction -= self.ERROR_ACTION
				for ; ; {
					self.stateStackTop -= self.prs.rhs(self.currentAction) - 1
					self.ra.ruleAction(self.currentAction)
					var lhs_symbol = self.prs.lhs(self.currentAction)
					if lhs_symbol == self.START_SYMBOL {
						continue ProcessTokens
					}

					self.currentAction = self.prs.ntAction(self.stack[self.stateStackTop], lhs_symbol)
					if self.currentAction <= self.NUM_RULES {
						continue
					} else {
						break
					}
				}
			} else {
				if self.currentAction < self.ACCEPT_ACTION { // Shift

					self.lastToken = self.curtok
					self.curtok = self.tokStream.getToken()
					self.current_kind = self.tokStream.getKind(self.curtok)
				} else {
					if self.currentAction == self.ACCEPT_ACTION {
						continue ProcessTokens
					} else {
						break ScanToken // ERROR_ACTION
					}
				}
			}
		}

		//
		// Whenever we reach self point, an error has been detected.
		// Note that the parser loop above can never reach the ACCEPT
		// point as it is short-circuited each time it reduces a phrase
		// to the START_SYMBOL.
		//
		// If an error is detected on a single bad character,
		// we advance to the next character before resuming the
		// scan. However, if an error is detected after we start
		// scanning a construct, we form a bad token out of the
		// characters that have already been scanned and resume
		// scanning on the character on which the problem was
		// detected. In other words, in that case, we do not advance.
		//
		if self.starttok == self.curtok {
			if self.current_kind == self.EOFT_SYMBOL {
				break ProcessTokens
			}

			self.tokStream.reportLexicalErrorPosition(self.starttok, self.curtok)
			self.lastToken = self.curtok
			self.curtok = self.tokStream.getToken()
			self.current_kind = self.tokStream.getKind(self.curtok)
		} else {
			self.tokStream.reportLexicalErrorPosition(self.starttok, self.lastToken)
		}

	}

	self.taking_actions = false // indicate that we are done

	return

}

//
// This function takes as argument a configuration ([stack, stackTop], [tokStream, curtok])
// and determines whether or not curtok can be validly parsed in self configuration. If so,
// it parses curtok and returns the final shift or shift-reduce action on it. Otherwise, it
// leaves the configuration unchanged and returns ERROR_ACTION.
//
func (self *LexParser) parseNextCharacter(token int, kind int) {
	var start_action = self.stack[self.stateStackTop]
	var pos = self.stateStackTop
	var tempStackTop = self.stateStackTop - 1

Scan:
	for self.currentAction = self.tAction(start_action, kind);
		self.currentAction <= self.NUM_RULES;
		self.currentAction = self.tAction(self.currentAction, kind) {
		for ; ; {
			var lhs_symbol = self.prs.lhs(self.currentAction)
			if lhs_symbol == self.START_SYMBOL {
				break Scan
			}
			tempStackTop -= self.prs.rhs(self.currentAction) - 1
			var state int
			if tempStackTop > pos {
				state = self.tempStack[tempStackTop]
			} else {
				state = self.stack[tempStackTop]
			}

			self.currentAction = self.prs.ntAction(state, lhs_symbol)
			if self.currentAction <= self.NUM_RULES {
				continue
			} else {
				break
			}
		}
		if tempStackTop+1 >= len(self.stack) {
			self.reallocateStacks()
		}
		//
		// ... Update the maximum useful position of the stack,
		// push goto state into (temporary) stack, and compute
		// the next action on the current symbol ...
		//
		if !(pos < tempStackTop) {
			pos = tempStackTop
		}
		self.tempStack[tempStackTop+1] = self.currentAction
	}

	//
	// If no error was detected, we update the configuration up to the point prior to the
	// shift or shift-reduce on the token by processing all reduce and goto actions associated
	// with the current token.
	//
	if self.currentAction != self.ERROR_ACTION {
		//
		// Note that it is important that the global variable currentAction be used here when
		// we are actually processing the rules. The reason being that the user-defined function
		// ra.ruleAction() may call func (self *LexParser) functions defined in self class (such as getLastToken())
		// which require that currentAction be properly initialized.
		//
	Replay:
		for self.currentAction = self.tAction(start_action, kind);
			self.currentAction <= self.NUM_RULES;
			self.currentAction = self.tAction(self.currentAction, kind) {
			self.stateStackTop--
			for ; ; {
				self.stateStackTop -= self.prs.rhs(self.currentAction) - 1
				self.ra.ruleAction(self.currentAction)
				var lhs_symbol = self.prs.lhs(self.currentAction)
				if lhs_symbol == self.START_SYMBOL {

					if self.starttok == token { // null string reduction to START_SYMBOL is illegal
						self.currentAction = self.ERROR_ACTION
					} else {
						self.currentAction = self.ACCEPT_ACTION
					}
					break Replay
				}
				self.currentAction = self.prs.ntAction(self.stack[self.stateStackTop], lhs_symbol)
				if self.currentAction <= self.NUM_RULES {
					continue
				} else {
					break
				}
			}
			self.stateStackTop += 1
			if self.stateStackTop >= len(self.stack) {
				self.reallocateStacks()
			}
			self.stack[self.stateStackTop] = self.currentAction

			self.locationStack[self.stateStackTop] = token
		}
	}

	return
}

//
// keep looking ahead until we compute a valid action
//
func (self *LexParser) lookahead(act int, token int) int {
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
func (self *LexParser) tAction(act int, sym int) int {
	act = self.prs.tAction(act, sym)
	if act > self.LA_STATE_OFFSET {
		return self.lookahead(act, self.tokStream.peek())
	} else {
		return act
	}

}

func (self *LexParser) scanNextToken() bool {
	return self.lexNextToken(self.tokStream.getStreamLength())
}
func (self *LexParser) scanNextTokenFromStartOffset(start_offset int) bool {
	self.resetTokenStream(start_offset)
	return self.lexNextToken(self.tokStream.getStreamLength())
}
func (self *LexParser) lexNextToken(end_offset int) bool {
	//
	// Indicate that we are going to run the incremental parser and that
	// it's forbidden to use the utility functions to query the parser.
	//
	self.taking_actions = false

	self.stateStackTop = -1
	self.currentAction = self.START_STATE
	self.starttok = self.curtok
	self.action.reset()

ScanToken:
	for ; ; {
		self.stateStackTop += 1
		if self.stateStackTop >= len(self.stack) {
			self.reallocateStacks()
		}
		self.stack[self.stateStackTop] = self.currentAction

		//
		// Compute the self.action on the next character. If it is a reduce self.action, we do not
		// want to accept it until we are sure that the character in question is parsable.
		// What we are trying to avoid is a situation where self.curtok is not the EOF token
		// but it yields a default reduce self.action in the current configuration even though
		// it cannot ultimately be shifted However, the state on top of the configuration also
		// contains a valid reduce self.action on EOF which, if taken, would lead to the succesful
		// scanning of the token.
		//
		// Thus, if the character is parsable, we proceed normally. Otherwise, we proceed
		// as if we had reached the end of the file (end of the token, since we are really
		// scanning).
		//
		self.currentAction = self.lexNextCharacter(self.currentAction, self.current_kind)
		if self.currentAction == self.ERROR_ACTION && self.current_kind != self.EOFT_SYMBOL { // if not successful try EOF

			var save_next_token = self.tokStream.peek()                  // save position after self.curtok
			self.tokStream.resetTo(self.tokStream.getStreamLength() - 1) // point to the end of the input
			self.currentAction = self.lexNextCharacter(self.stack[self.stateStackTop], self.EOFT_SYMBOL)
			// assert (self.currentAction == self.ACCEPT_ACTION || self.currentAction == self.ERROR_ACTION)
			self.tokStream.resetTo(save_next_token) // reset the stream for the next token after self.curtok.
		}

		self.action.add(self.currentAction) // save the self.action

		//
		// At self point, self.currentAction is either a Shift, Shift-Reduce, Accept or Error self.action.
		//
		if self.currentAction > self.ERROR_ACTION { //Shift-reduce

			self.curtok = self.tokStream.getToken()
			if self.curtok > end_offset {
				self.curtok = self.tokStream.getStreamLength()
			}

			self.current_kind = self.tokStream.getKind(self.curtok)
			self.currentAction -= self.ERROR_ACTION
			for ; ; {
				var lhs_symbol = self.prs.lhs(self.currentAction)
				if lhs_symbol == self.START_SYMBOL {
					self.parseActions()
					return true
				}
				self.stateStackTop -= self.prs.rhs(self.currentAction) - 1
				self.currentAction = self.prs.ntAction(self.stack[self.stateStackTop], lhs_symbol)
				if self.currentAction <= self.NUM_RULES {
					continue
				} else {
					break
				}
			}
		} else {
			if self.currentAction < self.ACCEPT_ACTION { // Shift

				self.curtok = self.tokStream.getToken()
				if self.curtok > end_offset {
					self.curtok = self.tokStream.getStreamLength()
				}
				self.current_kind = self.tokStream.getKind(self.curtok)
			} else {
				if self.currentAction == self.ACCEPT_ACTION {
					return true
				} else {
					break ScanToken // self.ERROR_ACTION
				}
			}

		}
	}

	//
	// Whenever we reach self point, an error has been detected.
	// Note that the parser loop above can never reach the ACCEPT
	// point as it is short-circuited each time it reduces a phrase
	// to the self.START_SYMBOL.
	//
	// If an error is detected on a single bad character,
	// we advance to the next character before resuming the
	// scan. However, if an error is detected after we start
	// scanning a construct, we form a bad token out of the
	// characters that have already been scanned and resume
	// scanning on the character on which the problem was
	// detected. In other words, in that case, we do not advance.
	//
	if self.starttok == self.curtok {
		if self.current_kind == self.EOFT_SYMBOL {
			self.action = nil // turn into garbage!
			return false
		}
		self.lastToken = self.curtok
		self.tokStream.reportLexicalErrorPosition(self.starttok, self.curtok)
		self.curtok = self.tokStream.getToken()
		if self.curtok > end_offset {
			self.curtok = self.tokStream.getStreamLength()
		}
		self.current_kind = self.tokStream.getKind(self.curtok)
	} else {
		self.lastToken = self.tokStream.getPrevious(self.curtok)
		self.tokStream.reportLexicalErrorPosition(self.starttok, self.lastToken)
	}

	return true
}

//
// This function takes as argument a configuration ([self.stack, stackTop], [self.tokStream, self.curtok])
// and determines whether or not the reduce self.action the self.curtok can be validly parsed in self
// configuration.
//
func (self *LexParser) lexNextCharacter(act int, kind int) int {
	var action_save = self.action.size()
	var pos = self.stateStackTop
	var tempStackTop = self.stateStackTop - 1
	act = self.tAction(act, kind)
Scan:
	for ; act <= self.NUM_RULES; {
		self.action.add(act)

		for ; ; {
			var lhs_symbol = self.prs.lhs(act)
			if lhs_symbol == self.START_SYMBOL {
				if self.starttok == self.curtok { // null string reduction to self.START_SYMBOL is illegal
					act = self.ERROR_ACTION
					break Scan
				} else {
					self.parseActions()
					return self.ACCEPT_ACTION
				}
			}
			tempStackTop -= self.prs.rhs(act) - 1
			var state int
			if tempStackTop > pos {
				state = self.tempStack[tempStackTop]
			} else {
				state = self.stack[tempStackTop]
			}
			act = self.prs.ntAction(state, lhs_symbol)
			if act <= self.NUM_RULES {
				continue
			} else {
				break
			}
		}
		if tempStackTop+1 >= len(self.stack) {
			self.reallocateStacks()
		}
		//
		// ... Update the maximum useful position of the self.stack,
		// push goto state into (temporary) self.stack, and compute
		// the next self.action on the current symbol ...
		//

		if !(pos < tempStackTop) {
			pos = tempStackTop
		}
		self.tempStack[tempStackTop+1] = act
		act = self.tAction(act, kind)
	}

	//
	// If an error was detected, we restore the original configuration.
	// Otherwise, we update configuration up to the point prior to the
	// shift or shift-reduce on the token.
	//
	if act == self.ERROR_ACTION {
		self.action.resetTo(action_save)
	} else {
		self.stateStackTop = tempStackTop + 1
		var i = pos + 1
		for ; i <= self.stateStackTop; i++ { // update self.stack
			self.stack[i] = self.tempStack[i]
		}
	}

	return act
}

//
// Now do the final parse of the input based on the actions in
// the list "self.action" and the sequence of tokens in the token stream.
//
func (self *LexParser) parseActions() {
	//
	// Indicate that we are running the regular parser and that it's
	// ok to use the utility functions to query the parser.
	//
	self.taking_actions = true

	self.curtok = self.starttok
	self.lastToken = self.tokStream.getPrevious(self.curtok)

	//
	// Reparse the input...
	//
	self.stateStackTop = -1
	self.currentAction = self.START_STATE
	var i = 0

process_actions:
	for ; i < self.action.size(); i++ {
		self.stateStackTop += 1
		self.stack[self.stateStackTop] = self.currentAction
		self.locationStack[self.stateStackTop] = self.curtok

		self.currentAction = self.action.get(i)
		if self.currentAction <= self.NUM_RULES { // a reduce self.action?

			self.stateStackTop-- // turn reduction intoshift-reduction
			for ; ; {
				self.stateStackTop -= self.prs.rhs(self.currentAction) - 1
				self.ra.ruleAction(self.currentAction)
				var lhs_symbol = self.prs.lhs(self.currentAction)
				if lhs_symbol == self.START_SYMBOL {
					// assert(starttok != self.curtok)  // null string reduction to self.START_SYMBOL is illegal
					break process_actions
				}
				self.currentAction = self.prs.ntAction(self.stack[self.stateStackTop], lhs_symbol)
				if self.currentAction <= self.NUM_RULES {
					continue
				} else {
					break
				}
			}
		} else { // a shift or shift-reduce self.action

			self.lastToken = self.curtok
			self.curtok = self.tokStream.getNext(self.curtok)
			if self.currentAction > self.ERROR_ACTION { // a shift-reduce self.action?

				self.current_kind = self.tokStream.getKind(self.curtok)
				self.currentAction -= self.ERROR_ACTION
				for ; ; {
					self.stateStackTop -= self.prs.rhs(self.currentAction) - 1
					self.ra.ruleAction(self.currentAction)
					var lhs_symbol = self.prs.lhs(self.currentAction)
					if lhs_symbol == self.START_SYMBOL {
						break process_actions
					}
					self.currentAction = self.prs.ntAction(self.stack[self.stateStackTop], lhs_symbol)
					if self.currentAction <= self.NUM_RULES {
						continue
					} else {
						break
					}
				}
			}
		}
	}

	self.taking_actions = false // indicate that we are done

	return
}
