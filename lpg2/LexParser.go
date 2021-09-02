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

func NewLexParser() *LexParser {
	return NewLexParserAndInit(nil,nil,nil)
}
func NewLexParserAndInit(tokStream ILexStream, prs ParseTable, ra RuleAction) *LexParser {
	my := new(LexParser)
	my.STACK_INCREMENT = 1024
    my.stackLength=0
	if tokStream != nil && prs != nil && ra != nil {
		my.Reset(tokStream, prs, ra)
	}

	return my
}

func (my *LexParser) Reset(tokStream ILexStream, prs ParseTable, ra RuleAction) {
	my.tokStream = tokStream
	my.prs = prs
	my.ra = ra
	my.START_STATE = prs.GetStartState()
	my.LA_STATE_OFFSET = prs.GetLaStateOffset()
	my.EOFT_SYMBOL = prs.GetEoftSymbol()
	my.ACCEPT_ACTION = prs.GetAcceptAction()
	my.ERROR_ACTION = prs.GetErrorAction()
	my.START_SYMBOL = prs.GetStartSymbol()
	my.NUM_RULES = prs.GetNumRules()
}

//
// Stacks portion
//

func (my *LexParser) ReallocateStacks() {
	var old_stack_length int
	if len(my.stack) == 0 {
		old_stack_length = 0
	} else {
		old_stack_length = my.stackLength
	}
	my.stackLength += my.STACK_INCREMENT
	if old_stack_length == 0 {
		my.stack = make([]int, my.stackLength)
		my.locationStack = make([]int, my.stackLength)
		my.tempStack = make([]int, my.stackLength)
	} else {
		my.stack = Arraycopy(my.stack, 0, make([]int, my.stackLength), 0, old_stack_length)
		my.locationStack = Arraycopy(my.locationStack, 0, make([]int, my.stackLength), 0, old_stack_length)
		my.tempStack = Arraycopy(my.tempStack, 0, make([]int, my.stackLength), 0, old_stack_length)
	}
	return
}

//
// The following functions can be invoked only when the parser is
// processing actions. Thus, they can be invoked when the parser
// was entered via the main entry point (ParseCharacters()). When using
// the incremental parser (via the entry point ScanNextToken(int [], int)),
// they always return 0 when invoked. // TODO Should we throw an Exception instead?
// However, note that when ParseActions() is invoked after successfully
// parsing an input with the incremental parser, then they can be invoked.
//
func (my *LexParser) GetFirstTokenAt(i int) int {
	return my.GetToken(i)
}
func (my *LexParser) GetFirstToken() int {
	return my.starttok
}
func (my *LexParser) GetLastToken() int {
	return my.lastToken
}
func (my *LexParser) GetLastTokenAt(i int) int {

	if my.taking_actions {
		if i >= my.prs.Rhs(my.currentAction) {
			return my.lastToken
		} else {
			var index = my.GetToken(i + 1)
			return my.tokStream.GetPrevious(index)
		}

	}
	return -1
}
func (my *LexParser) GetCurrentRule() (int, error) {
	if my.taking_actions {
		return my.currentAction, nil
	}
	return -1, NewUnavailableParserInformationException("")
}

//
// Given a rule of the form     A = x1 x2 ... xn     n > 0
//
// the function GetToken(i) yields the symbol xi, if xi is a terminal
// or ti, if xi is a nonterminal that produced a string of the form
// xi => ti w. If xi is a nullable nonterminal, then ti is the first
//  symbol that immediately follows xi in the input (the Lookahead).
//
func (my *LexParser) GetToken(i int) int {
	if my.taking_actions {
		return my.locationStack[my.stateStackTop+(i-1)]
	}
	return -1
}
func (my *LexParser) SetSym1(i int) {}
func (my *LexParser) GetSym(i int) int {
	i = my.GetLastTokenAt(i)
	return i
}

func (my *LexParser) ResetTokenStream(i int) {
	//
	// if i exceeds the upper bound, Reset it to point to the last element.
	//
	var temp int
	if i > my.tokStream.GetStreamLength() {
		temp = my.tokStream.GetStreamLength()
	} else {
		temp = i
	}
	my.tokStream.ResetTo(temp)
	my.curtok = my.tokStream.GetToken()
	my.current_kind = my.tokStream.GetKind(my.curtok)
	if len(my.stack) == 0 {
		my.ReallocateStacks()
	}
	if my.action == nil {
		my.action = NewIntTupleWithEstimate(1 << 10)
	}
}

//
// Parse the input and create a stream of tokens.
//
func (my *LexParser) ParseCharacters(start_offSet int, end_offSet int, monitor Monitor) {
	my.ResetTokenStream(start_offSet)
	for ; my.curtok <= end_offSet; {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if monitor != nil && monitor.IsCancelled() {
			return
		}
		my.LexNextToken(end_offSet)
	}
}

//
// Parse the input and create a stream of tokens.
//
func (my *LexParser) ParseCharactersWhitMonitor(monitor Monitor) {
	//
	// Indicate that we are running the regular parser and that it's
	// ok to use the utility functions to query the parser.
	//
	my.taking_actions = true
	my.ResetTokenStream(0)
	my.lastToken = my.tokStream.GetPrevious(my.curtok)

	//
	// Until it reaches the end-of-file token, my outer loop
	// reSets the parser and processes the next token.
	//
ProcessTokens:
	for ; my.current_kind != my.EOFT_SYMBOL; {
		//
		// if the parser needs to stop processing,
		// it may do so here.
		//
		if monitor != nil && monitor.IsCancelled() {
			break ProcessTokens
		}

		my.stateStackTop = -1
		my.currentAction = my.START_STATE
		my.starttok = my.curtok

	ScanToken:
		for ; ; {
			my.stateStackTop += 1
			if my.stateStackTop >= len(my.stack) {
				my.ReallocateStacks()
			}
			my.stack[my.stateStackTop] = my.currentAction

			my.locationStack[my.stateStackTop] = my.curtok

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
			my.ParseNextCharacter(my.curtok, my.current_kind)
			if my.currentAction == my.ERROR_ACTION && my.current_kind != my.EOFT_SYMBOL { // if not successful try EOF

				var save_next_token = my.tokStream.Peek()                  // save position after curtok
				my.tokStream.ResetTo(my.tokStream.GetStreamLength() - 1) // point to the end of the input
				my.ParseNextCharacter(my.curtok, my.EOFT_SYMBOL)
				// assert (currentAction == ACCEPT_ACTION || currentAction == ERROR_ACTION)
				my.tokStream.ResetTo(save_next_token) // Reset the stream for the next token after curtok.
			}

			//
			// At my point, currentAction is either a Shift, Shift-Reduce, Accept or Error action.
			//
			if my.currentAction > my.ERROR_ACTION { // Shift-reduce

				my.lastToken = my.curtok
				my.curtok = my.tokStream.GetToken()
				my.current_kind = my.tokStream.GetKind(my.curtok)
				my.currentAction -= my.ERROR_ACTION
				for ; ; {
					my.stateStackTop -= my.prs.Rhs(my.currentAction) - 1
					my.ra.RuleAction(my.currentAction)
					var lhs_symbol = my.prs.Lhs(my.currentAction)
					if lhs_symbol == my.START_SYMBOL {
						continue ProcessTokens
					}

					my.currentAction = my.prs.NtAction(my.stack[my.stateStackTop], lhs_symbol)
					if my.currentAction <= my.NUM_RULES {
						continue
					} else {
						break
					}
				}
			} else {
				if my.currentAction < my.ACCEPT_ACTION { // Shift

					my.lastToken = my.curtok
					my.curtok = my.tokStream.GetToken()
					my.current_kind = my.tokStream.GetKind(my.curtok)
				} else {
					if my.currentAction == my.ACCEPT_ACTION {
						continue ProcessTokens
					} else {
						break ScanToken // ERROR_ACTION
					}
				}
			}
		}

		//
		// Whenever we reach my point, an error has been detected.
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
		if my.starttok == my.curtok {
			if my.current_kind == my.EOFT_SYMBOL {
				break ProcessTokens
			}

			my.tokStream.ReportLexicalErrorPosition(my.starttok, my.curtok)
			my.lastToken = my.curtok
			my.curtok = my.tokStream.GetToken()
			my.current_kind = my.tokStream.GetKind(my.curtok)
		} else {
			my.tokStream.ReportLexicalErrorPosition(my.starttok, my.lastToken)
		}

	}

	my.taking_actions = false // indicate that we are done

	return

}

//
// This function takes as argument a configuration ([stack, stackTop], [tokStream, curtok])
// and determines whether or not curtok can be validly parsed in my configuration. If so,
// it parses curtok and returns the final shift or shift-reduce action on it. Otherwise, it
// leaves the configuration unchanged and returns ERROR_ACTION.
//
func (my *LexParser) ParseNextCharacter(token int, kind int) {
	var start_action = my.stack[my.stateStackTop]
	var pos = my.stateStackTop
	var tempStackTop = my.stateStackTop - 1

Scan:
	for my.currentAction = my.TAction(start_action, kind);
		my.currentAction <= my.NUM_RULES;
		my.currentAction = my.TAction(my.currentAction, kind) {
		for ; ; {
			var lhs_symbol = my.prs.Lhs(my.currentAction)
			if lhs_symbol == my.START_SYMBOL {
				break Scan
			}
			tempStackTop -= my.prs.Rhs(my.currentAction) - 1
			var state int
			if tempStackTop > pos {
				state = my.tempStack[tempStackTop]
			} else {
				state = my.stack[tempStackTop]
			}

			my.currentAction = my.prs.NtAction(state, lhs_symbol)
			if my.currentAction <= my.NUM_RULES {
				continue
			} else {
				break
			}
		}
		if tempStackTop+1 >= len(my.stack) {
			my.ReallocateStacks()
		}
		//
		// ... Update the maximum useful position of the stack,
		// push goto state into (temporary) stack, and compute
		// the next action on the current symbol ...
		//
		if !(pos < tempStackTop) {
			pos = tempStackTop
		}
		my.tempStack[tempStackTop+1] = my.currentAction
	}

	//
	// If no error was detected, we update the configuration up to the point prior to the
	// shift or shift-reduce on the token by processing all reduce and goto actions associated
	// with the current token.
	//
	if my.currentAction != my.ERROR_ACTION {
		//
		// Note that it is important that the global variable currentAction be used here when
		// we are actually processing the rules. The reason being that the user-defined function
		// ra.RuleAction() may call func (my *LexParser) functions defined in my class (such as GetLastToken())
		// which require that currentAction be properly initialized.
		//
	Replay:
		for my.currentAction = my.TAction(start_action, kind);
			my.currentAction <= my.NUM_RULES;
			my.currentAction = my.TAction(my.currentAction, kind) {
			my.stateStackTop--
			for ; ; {
				my.stateStackTop -= my.prs.Rhs(my.currentAction) - 1
				my.ra.RuleAction(my.currentAction)
				var lhs_symbol = my.prs.Lhs(my.currentAction)
				if lhs_symbol == my.START_SYMBOL {

					if my.starttok == token { // null string reduction to START_SYMBOL is illegal
						my.currentAction = my.ERROR_ACTION
					} else {
						my.currentAction = my.ACCEPT_ACTION
					}
					break Replay
				}
				my.currentAction = my.prs.NtAction(my.stack[my.stateStackTop], lhs_symbol)
				if my.currentAction <= my.NUM_RULES {
					continue
				} else {
					break
				}
			}
			my.stateStackTop += 1
			if my.stateStackTop >= len(my.stack) {
				my.ReallocateStacks()
			}
			my.stack[my.stateStackTop] = my.currentAction

			my.locationStack[my.stateStackTop] = token
		}
	}

	return
}

//
// keep looking ahead until we compute a valid action
//
func (my *LexParser) Lookahead(act int, token int) int {
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
func (my *LexParser) TAction(act int, sym int) int {
	act = my.prs.TAction(act, sym)
	if act > my.LA_STATE_OFFSET {
		return my.Lookahead(act, my.tokStream.Peek())
	} else {
		return act
	}

}

func (my *LexParser) ScanNextToken() bool {
	return my.LexNextToken(my.tokStream.GetStreamLength())
}
func (my *LexParser) ScanNextTokenFromStartOffSet(start_offSet int) bool {
	my.ResetTokenStream(start_offSet)
	return my.LexNextToken(my.tokStream.GetStreamLength())
}
func (my *LexParser) LexNextToken(end_offSet int) bool {
	//
	// Indicate that we are going to run the incremental parser and that
	// it's forbidden to use the utility functions to query the parser.
	//
	my.taking_actions = false

	my.stateStackTop = -1
	my.currentAction = my.START_STATE
	my.starttok = my.curtok
	my.action.ReSet()

ScanToken:
	for ; ; {
		my.stateStackTop += 1
		if my.stateStackTop >= len(my.stack) {
			my.ReallocateStacks()
		}
		my.stack[my.stateStackTop] = my.currentAction

		//
		// Compute the my.action on the next character. If it is a reduce my.action, we do not
		// want to accept it until we are sure that the character in question is parsable.
		// What we are trying to avoid is a situation where my.curtok is not the EOF token
		// but it yields a default reduce my.action in the current configuration even though
		// it cannot ultimately be shifted However, the state on top of the configuration also
		// contains a valid reduce my.action on EOF which, if taken, would lead to the succesful
		// scanning of the token.
		//
		// Thus, if the character is parsable, we proceed normally. Otherwise, we proceed
		// as if we had reached the end of the file (end of the token, since we are really
		// scanning).
		//
		my.currentAction = my.LexNextCharacter(my.currentAction, my.current_kind)
		if my.currentAction == my.ERROR_ACTION && my.current_kind != my.EOFT_SYMBOL { // if not successful try EOF

			var save_next_token = my.tokStream.Peek()                  // save position after my.curtok
			my.tokStream.ResetTo(my.tokStream.GetStreamLength() - 1) // point to the end of the input
			my.currentAction = my.LexNextCharacter(my.stack[my.stateStackTop], my.EOFT_SYMBOL)
			// assert (my.currentAction == my.ACCEPT_ACTION || my.currentAction == my.ERROR_ACTION)
			my.tokStream.ResetTo(save_next_token) // Reset the stream for the next token after my.curtok.
		}

		my.action.Add(my.currentAction) // save the my.action

		//
		// At my point, my.currentAction is either a Shift, Shift-Reduce, Accept or Error my.action.
		//
		if my.currentAction > my.ERROR_ACTION { //Shift-reduce

			my.curtok = my.tokStream.GetToken()
			if my.curtok > end_offSet {
				my.curtok = my.tokStream.GetStreamLength()
			}

			my.current_kind = my.tokStream.GetKind(my.curtok)
			my.currentAction -= my.ERROR_ACTION
			for ; ; {
				var lhs_symbol = my.prs.Lhs(my.currentAction)
				if lhs_symbol == my.START_SYMBOL {
					my.ParseActions()
					return true
				}
				my.stateStackTop -= my.prs.Rhs(my.currentAction) - 1
				my.currentAction = my.prs.NtAction(my.stack[my.stateStackTop], lhs_symbol)
				if my.currentAction <= my.NUM_RULES {
					continue
				} else {
					break
				}
			}
		} else {
			if my.currentAction < my.ACCEPT_ACTION { // Shift

				my.curtok = my.tokStream.GetToken()
				if my.curtok > end_offSet {
					my.curtok = my.tokStream.GetStreamLength()
				}
				my.current_kind = my.tokStream.GetKind(my.curtok)
			} else {
				if my.currentAction == my.ACCEPT_ACTION {
					return true
				} else {
					break ScanToken // my.ERROR_ACTION
				}
			}

		}
	}

	//
	// Whenever we reach my point, an error has been detected.
	// Note that the parser loop above can never reach the ACCEPT
	// point as it is short-circuited each time it reduces a phrase
	// to the my.START_SYMBOL.
	//
	// If an error is detected on a single bad character,
	// we advance to the next character before resuming the
	// scan. However, if an error is detected after we start
	// scanning a construct, we form a bad token out of the
	// characters that have already been scanned and resume
	// scanning on the character on which the problem was
	// detected. In other words, in that case, we do not advance.
	//
	if my.starttok == my.curtok {
		if my.current_kind == my.EOFT_SYMBOL {
			my.action = nil // turn into garbage!
			return false
		}
		my.lastToken = my.curtok
		my.tokStream.ReportLexicalErrorPosition(my.starttok, my.curtok)
		my.curtok = my.tokStream.GetToken()
		if my.curtok > end_offSet {
			my.curtok = my.tokStream.GetStreamLength()
		}
		my.current_kind = my.tokStream.GetKind(my.curtok)
	} else {
		my.lastToken = my.tokStream.GetPrevious(my.curtok)
		my.tokStream.ReportLexicalErrorPosition(my.starttok, my.lastToken)
	}

	return true
}

//
// This function takes as argument a configuration ([my.stack, stackTop], [my.tokStream, my.curtok])
// and determines whether or not the reduce my.action the my.curtok can be validly parsed in my
// configuration.
//
func (my *LexParser) LexNextCharacter(act int, kind int) int {
	var action_save = my.action.Size()
	var pos = my.stateStackTop
	var tempStackTop = my.stateStackTop - 1
	act = my.TAction(act, kind)
Scan:
	for ; act <= my.NUM_RULES; {
		my.action.Add(act)

		for ; ; {
			var lhs_symbol = my.prs.Lhs(act)
			if lhs_symbol == my.START_SYMBOL {
				if my.starttok == my.curtok { // null string reduction to my.START_SYMBOL is illegal
					act = my.ERROR_ACTION
					break Scan
				} else {
					my.ParseActions()
					return my.ACCEPT_ACTION
				}
			}
			tempStackTop -= my.prs.Rhs(act) - 1
			var state int
			if tempStackTop > pos {
				state = my.tempStack[tempStackTop]
			} else {
				state = my.stack[tempStackTop]
			}
			act = my.prs.NtAction(state, lhs_symbol)
			if act <= my.NUM_RULES {
				continue
			} else {
				break
			}
		}
		if tempStackTop+1 >= len(my.stack) {
			my.ReallocateStacks()
		}
		//
		// ... Update the maximum useful position of the my.stack,
		// push goto state into (temporary) my.stack, and compute
		// the next my.action on the current symbol ...
		//

		if !(pos < tempStackTop) {
			pos = tempStackTop
		}
		my.tempStack[tempStackTop+1] = act
		act = my.TAction(act, kind)
	}

	//
	// If an error was detected, we restore the original configuration.
	// Otherwise, we update configuration up to the point prior to the
	// shift or shift-reduce on the token.
	//
	if act == my.ERROR_ACTION {
		my.action.ReSetTo(action_save)
	} else {
		my.stateStackTop = tempStackTop + 1
		var i = pos + 1
		for ; i <= my.stateStackTop; i++ { // update my.stack
			my.stack[i] = my.tempStack[i]
		}
	}

	return act
}

//
// Now do the final parse of the input based on the actions in
// the list "my.action" and the sequence of tokens in the token stream.
//
func (my *LexParser) ParseActions() {
	//
	// Indicate that we are running the regular parser and that it's
	// ok to use the utility functions to query the parser.
	//
	my.taking_actions = true

	my.curtok = my.starttok
	my.lastToken = my.tokStream.GetPrevious(my.curtok)

	//
	// Reparse the input...
	//
	my.stateStackTop = -1
	my.currentAction = my.START_STATE
	var i = 0

process_actions:
	for ; i < my.action.Size(); i++ {
		my.stateStackTop += 1
		my.stack[my.stateStackTop] = my.currentAction
		my.locationStack[my.stateStackTop] = my.curtok

		my.currentAction = my.action.Get(i)
		if my.currentAction <= my.NUM_RULES { // a reduce my.action?

			my.stateStackTop-- // turn reduction intoshift-reduction
			for ; ; {
				my.stateStackTop -= my.prs.Rhs(my.currentAction) - 1
				my.ra.RuleAction(my.currentAction)
				var lhs_symbol = my.prs.Lhs(my.currentAction)
				if lhs_symbol == my.START_SYMBOL {
					// assert(starttok != my.curtok)  // null string reduction to my.START_SYMBOL is illegal
					break process_actions
				}
				my.currentAction = my.prs.NtAction(my.stack[my.stateStackTop], lhs_symbol)
				if my.currentAction <= my.NUM_RULES {
					continue
				} else {
					break
				}
			}
		} else { // a shift or shift-reduce my.action

			my.lastToken = my.curtok
			my.curtok = my.tokStream.GetNext(my.curtok)
			if my.currentAction > my.ERROR_ACTION { // a shift-reduce my.action?

				my.current_kind = my.tokStream.GetKind(my.curtok)
				my.currentAction -= my.ERROR_ACTION
				for ; ; {
					my.stateStackTop -= my.prs.Rhs(my.currentAction) - 1
					my.ra.RuleAction(my.currentAction)
					var lhs_symbol = my.prs.Lhs(my.currentAction)
					if lhs_symbol == my.START_SYMBOL {
						break process_actions
					}
					my.currentAction = my.prs.NtAction(my.stack[my.stateStackTop], lhs_symbol)
					if my.currentAction <= my.NUM_RULES {
						continue
					} else {
						break
					}
				}
			}
		}
	}

	my.taking_actions = false // indicate that we are done

	return
}
