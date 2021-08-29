package lpg2


type LexParser struct {

    def reset( tokStream: ILexStream, prs: ParseTable, ra: RuleAction):
        a.tokStream = tokStream
        a.prs = prs
        a.ra = ra
        a.START_STATE = prs.getStartState()
        a.LA_STATE_OFFSET = prs.getLaStateOffset()
        a.EOFT_SYMBOL = prs.getEoftSymbol()
        a.ACCEPT_ACTION = prs.getAcceptAction()
        a.ERROR_ACTION = prs.getErrorAction()
        a.START_SYMBOL = prs.getStartSymbol()
        a.NUM_RULES = prs.getNumRules()

    def __init__( tokStream: ILexStream = nil, prs: ParseTable = nil, ra: RuleAction = nil):
        a.taking_actions: bool = False

        a.START_STATE int = 0
        a.LA_STATE_OFFSET int = 0
        a.EOFT_SYMBOL int = 0
        a.ACCEPT_ACTION int = 0
        a.ERROR_ACTION int = 0
        a.START_SYMBOL int = 0
        a.NUM_RULES int = 0

        a.tokStream: ILexStream = nil
        a.prs: ParseTable = nil
        a.ra: RuleAction = nil
        a.action: IntTuple = IntTuple(0)

        a.stateStackTop int = 0
        a.stackLength int = 0
        a.stack: list = []
        a.locationStack: list = []
        a.tempStack: list = []

        a.lastToken int = 0
        a.currentAction int = 0
        a.curtok int = 0
        a.starttok int = 0
        a.current_kind int = 0
        if tokStream and prs and ra:
            a.reset(tokStream, prs, ra)

    //
    // Stacks portion
    //
    STACK_INCREMENT int = 1024

    def reallocateStacks()
        old_stack_length int = (0 if a.stack.__len__() == 0 else a.stackLength)
        a.stackLength += a.STACK_INCREMENT
        if old_stack_length == 0:
            a.stack = [0] * a.stackLength
            a.locationStack = [0] * a.stackLength
            a.tempStack = [0] * a.stackLength
        else:
            a.stack = arraycopy(a.stack, 0, [0] * a.stackLength, 0, old_stack_length)
            a.locationStack = arraycopy(a.locationStack, 0, [0] * a.stackLength, 0, old_stack_length)
            a.tempStack = arraycopy(a.tempStack, 0, [0] * a.stackLength, 0, old_stack_length)

        return

    //
    // The following functions can be invoked only when the parser is
    // processing actions. Thus, they can be invoked when the parser
    // was entered via the main entry point (parseCharacters()). When using
    // the incremental parser (via the entry point scanNextToken(int [], int)):,
    // they always return 0 when invoked. // TODO: Should we raise an Exception instead?
    // However, note that when parseActions() is invoked after successfully
    // parsing an input with the incremental parser, then they can be invoked.
    //
    def getFirstToken( i int = nil) int {
        if i is nil:
            return a.starttok

        return a.getToken(i)

    def getLastToken( i int = nil) int {
        if i is nil:
            return a.lastToken

        if a.taking_actions:
            return (a.lastToken if i >= a.prs.rhs(a.currentAction) else a.tokStream.getPrevious(
                a.getToken(i + 1)))

        raise UnavailableParserInformationException()

    def getCurrentRule() int {
        if a.taking_actions:
            return a.currentAction

        raise UnavailableParserInformationException()

    //
    // Given a rule of the form     A ::= x1 x2 ... xn     n > 0
    //
    // the function getToken(i): yields the symbol xi, if xi is a terminal
    // or ti, if xi is a nonterminal that produced a string of the form
    // xi => ti w. If xi is a nullable nonterminal, then ti is the first
    //  symbol that immediately follows xi in the input (the lookahead).
    //
    def getToken( i int) int {
        if a.taking_actions:
            return a.locationStack[a.stateStackTop + (i - 1)]

        raise UnavailableParserInformationException()

    def setSym1( i int):
        pass

    def getSym( i int) int {
        return a.getLastToken(i)

    def resetTokenStream( i int):
        //
        // if i exceeds the upper bound, reset it to point to the last element.
        //
        a.tokStream.reset(a.tokStream.getStreamLength() if i > a.tokStream.getStreamLength() else i)
        a.curtok = a.tokStream.getToken()
        a.current_kind = a.tokStream.getKind(a.curtok)
        if not a.stack or a.stack.__len__() == 0:
            a.reallocateStacks()

        if a.action.capacity() == 0:
            a.action = IntTuple(1 << 10)

    //
    // Parse the input and create a stream of tokens.
    //
    def parseCharacters( start_offset int, end_offset int, monitor: Monitor):
        a.resetTokenStream(start_offset)
        while a.curtok <= end_offset:
            //
            // if the parser needs to stop processing,
            // it may do so here.
            //
            if monitor and monitor.isCancelled()
                return

            a.lexNextToken(end_offset)

    //
    // Parse the input and create a stream of tokens.
    //
    def parseCharactersWhitMonitor( monitor: Monitor = nil):
        //
        // Indicate that we are running the regular parser and that it's
        // ok to use the utility functions to query the parser.
        //
        a.taking_actions = true

        a.resetTokenStream(0)
        a.lastToken = a.tokStream.getPrevious(a.curtok)
        //
        // Until it reaches the end-of-file token, a outer loop
        // resets the parser and processes the next token.
        //

        // ProcessTokens:
        while a.current_kind != a.EOFT_SYMBOL:
            //
            // if the parser needs to stop processing,
            // it may do so here.
            //
            if monitor is not nil and monitor.isCancelled()
                break  // ProcessTokens

            a.stateStackTop = -1
            a.currentAction = a.START_STATE
            a.starttok = a.curtok
            b_continue_process_tokens: bool = False
            // ScanToken:
            while true:

                a.stateStackTop += 1
                if a.stateStackTop >= a.stack.__len__()
                    a.reallocateStacks()
                a.stack[a.stateStackTop] = a.currentAction

                a.locationStack[a.stateStackTop] = a.curtok

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
                if a.curtok == 275:
                    a.curtok = 275
                a.parseNextCharacter(a.curtok, a.current_kind)
                if a.curtok == 275:
                    a.curtok = 275

                if (a.currentAction == a.ERROR_ACTION and
                        a.current_kind != a.EOFT_SYMBOL):  // if not successful try EOF

                    save_next_token = a.tokStream.peek()  // save position after curtok
                    a.tokStream.reset(a.tokStream.getStreamLength() - 1)  // point to the end of the input
                    a.parseNextCharacter(a.curtok, a.EOFT_SYMBOL)
                    // assert (currentAction == ACCEPT_ACTION or currentAction == ERROR_ACTION):
                    a.tokStream.reset(save_next_token)  // reset the stream for the next token after curtok.

                //
                // At a point, currentAction is either a Shift, Shift-Reduce, Accept or Error action.
                //
                if a.currentAction > a.ERROR_ACTION:  // Shift-reduce

                    a.lastToken = a.curtok
                    a.curtok = a.tokStream.getToken()
                    a.current_kind = a.tokStream.getKind(a.curtok)
                    a.currentAction -= a.ERROR_ACTION
                    while true:
                        a.stateStackTop -= (a.prs.rhs(a.currentAction) - 1)
                        a.ra.ruleAction(a.currentAction)
                        lhs_symbol = a.prs.lhs(a.currentAction)
                        if lhs_symbol == a.START_SYMBOL:
                            b_continue_process_tokens = true
                            break
                        a.currentAction = a.prs.ntAction(a.stack[a.stateStackTop], lhs_symbol)
                        if not a.currentAction <= a.NUM_RULES:
                            break
                    if b_continue_process_tokens:
                        break
                elif a.currentAction < a.ACCEPT_ACTION:  // Shift

                    a.lastToken = a.curtok
                    a.curtok = a.tokStream.getToken()
                    a.current_kind = a.tokStream.getKind(a.curtok)

                elif a.currentAction == a.ACCEPT_ACTION:
                    b_continue_process_tokens = true
                    break
                else:
                    break  // ScanToken ERROR_ACTION

            if b_continue_process_tokens:
                continue
            //
            // Whenever we reach a point, an error has been detected.
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
            if a.starttok == a.curtok:
                if a.current_kind == a.EOFT_SYMBOL:
                    break  // ProcessTokens
                a.tokStream.reportLexicalError(a.starttok, a.curtok)
                a.lastToken = a.curtok
                a.curtok = a.tokStream.getToken()
                a.current_kind = a.tokStream.getKind(a.curtok)

            else:
                a.tokStream.reportLexicalError(a.starttok, a.lastToken)

        a.taking_actions = False  // indicate that we are done

        return

    //
    // This function takes as argument a configuration ([stack, stackTop], [tokStream, curtok]):
    // and determines whether or not curtok can be validly parsed in a configuration. If so,
    // it parses curtok and returns the final shift or shift-reduce action on it. Otherwise, it
    // leaves the configuration unchanged and returns ERROR_ACTION.
    //
    def parseNextCharacter( token int, kind int):
        start_action int = a.stack[a.stateStackTop]
        pos int = a.stateStackTop
        tempStackTop int = a.stateStackTop - 1

        a.currentAction = a.tAction(start_action, kind)

        b_break_scan: bool = False
        // Scan:
        while a.currentAction <= a.NUM_RULES:
            while true:
                lhs_symbol = a.prs.lhs(a.currentAction)
                if lhs_symbol == a.START_SYMBOL:
                    b_break_scan = true
                    break

                tempStackTop -= (a.prs.rhs(a.currentAction) - 1)

                state = (a.tempStack[tempStackTop] if tempStackTop > pos else a.stack[tempStackTop])

                a.currentAction = a.prs.ntAction(state, lhs_symbol)
                if not a.currentAction <= a.NUM_RULES:
                    break

            if b_break_scan:
                break

            if tempStackTop + 1 >= a.stack.__len__()
                a.reallocateStacks()
            //
            // ... Update the maximum useful position of the stack,
            // push goto state into (temporary): stack, and compute
            // the next action on the current symbol ...
            //
            pos = pos if pos < tempStackTop else tempStackTop
            a.tempStack[tempStackTop + 1] = a.currentAction

            a.currentAction = a.tAction(a.currentAction, kind)
        //
        // If no error was detected, we update the configuration up to the point prior to the
        // shift or shift-reduce on the token by processing all reduce and goto actions associated
        // with the current token.
        //
        if a.currentAction != a.ERROR_ACTION:
            //
            // Note that it is important that the global variable currentAction be used here when
            // we are actually processing the rules. The reason being that the user-defined function
            // ra.ruleAction() may call def functions defined in a type (such as getLastToken()):
            // which require that currentAction be properly initialized.
            //

            a.currentAction = a.tAction(start_action, kind)
            // Replay:
            bBreakReplay: bool = False
            while a.currentAction <= a.NUM_RULES:
                a.stateStackTop -= 1
                while true:
                    a.stateStackTop -= (a.prs.rhs(a.currentAction) - 1)
                    a.ra.ruleAction(a.currentAction)
                    lhs_symbol = a.prs.lhs(a.currentAction)
                    if lhs_symbol == a.START_SYMBOL:
                        a.currentAction = (a.ERROR_ACTION
                                              if a.starttok == token  // nil string reduction to START_SYMBOL is illegal
                                              else a.ACCEPT_ACTION)
                        bBreakReplay = true
                        break  // Replay

                    a.currentAction = a.prs.ntAction(a.stack[a.stateStackTop], lhs_symbol)
                    if not a.currentAction <= a.NUM_RULES:
                        break
                if bBreakReplay:
                    break

                a.stateStackTop += 1
                if a.stateStackTop >= a.stack.__len__()
                    a.reallocateStacks()
                a.stack[a.stateStackTop] = a.currentAction

                a.locationStack[a.stateStackTop] = token

                a.currentAction = a.tAction(a.currentAction, kind)

        return

    //
    // keep looking ahead until we compute a valid action
    //
    def lookahead( act int, token int) int {
        act = a.prs.lookAhead(act - a.LA_STATE_OFFSET, a.tokStream.getKind(token))
        return a.lookahead(act, a.tokStream.getNext(token)) if act > a.LA_STATE_OFFSET else act

    //
    // Compute the next action defined on act and sym. If a
    // action requires more lookahead, these lookahead symbols
    // are in the token stream beginning at the next token that
    // is yielded by peek().
    //
    def tAction( act int, sym int) int {
        act = a.prs.tAction(act, sym)
        return a.lookahead(act, a.tokStream.peek()) if act > a.LA_STATE_OFFSET else act

    def scanNextToken2()  bool:
        return a.lexNextToken(a.tokStream.getStreamLength())

    def scanNextToken( start_offset int = nil)  bool:
        if start_offset is nil:
            return a.scanNextToken2()

        a.resetTokenStream(start_offset)
        return a.lexNextToken(a.tokStream.getStreamLength())

    def lexNextToken( end_offset int)  bool:
        //
        // Indicate that we are going to run the incremental parser and that
        // it's forbidden to use the utility functions to query the parser.
        //
        a.taking_actions = False

        a.stateStackTop = -1
        a.currentAction = a.START_STATE
        a.starttok = a.curtok
        a.action.reset()

        // ScanToken:
        while true:

            a.stateStackTop += 1
            if a.stateStackTop >= a.stack.__len__()
                a.reallocateStacks()
            a.stack[a.stateStackTop] = a.currentAction

            //
            // Compute the the action on the next character. If it is a reduce action, we do not
            // want to accept it until we are sure that the character in question is parsable.
            // What we are trying to avoid is a situation where a.curtok is not the EOF token
            // but it yields a default reduce a.action in the current configuration even though
            // it cannot ultimately be shifted However, the state on top of the configuration also
            // contains a valid reduce a.action on EOF which, if taken, would lead to the succesful
            // scanning of the token.
            //
            // Thus, if the character is parsable, we proceed normally. Otherwise, we proceed
            // as if we had reached the end of the file (end of the token, since we are really
            // scanning).
            //
            a.currentAction = a.lexNextCharacter(a.currentAction, a.current_kind)
            if (a.currentAction == a.ERROR_ACTION and
                    a.current_kind != a.EOFT_SYMBOL):  // if not successful try EOF

                save_next_token = a.tokStream.peek()  // save position after a.curtok
                a.tokStream.reset(a.tokStream.getStreamLength() - 1)  // point to the end of the input
                a.currentAction = a.lexNextCharacter(a.stack[a.stateStackTop], a.EOFT_SYMBOL)
                // assert (a.currentAction == a.ACCEPT_ACTION or a.currentAction == a.ERROR_ACTION):
                a.tokStream.reset(save_next_token)  // reset the stream for the next token after a.curtok.

            a.action.add(a.currentAction)  // save the a.action

            //
            // At a point, a.currentAction is either a Shift, Shift-Reduce, Accept or Error a.action.
            //
            if a.currentAction > a.ERROR_ACTION:  // Shift-reduce

                a.curtok = a.tokStream.getToken()
                if a.curtok > end_offset:
                    a.curtok = a.tokStream.getStreamLength()
                a.current_kind = a.tokStream.getKind(a.curtok)
                a.currentAction -= a.ERROR_ACTION
                while true:
                    lhs_symbol = a.prs.lhs(a.currentAction)
                    if lhs_symbol == a.START_SYMBOL:
                        a.parseActions()
                        return true

                    a.stateStackTop -= (a.prs.rhs(a.currentAction) - 1)
                    a.currentAction = a.prs.ntAction(a.stack[a.stateStackTop], lhs_symbol)
                    if not a.currentAction <= a.NUM_RULES:
                        break

            elif a.currentAction < a.ACCEPT_ACTION:  // Shift

                a.curtok = a.tokStream.getToken()
                if a.curtok > end_offset:
                    a.curtok = a.tokStream.getStreamLength()
                a.current_kind = a.tokStream.getKind(a.curtok)

            elif a.currentAction == a.ACCEPT_ACTION:
                return true
            else:
                break  // ScanToken // a.ERROR_ACTION

        //
        // Whenever we reach a point, an error has been detected.
        // Note that the parser loop above can never reach the ACCEPT
        // point as it is short-circuited each time it reduces a phrase
        // to the a.START_SYMBOL.
        //
        // If an error is detected on a single bad character,
        // we advance to the next character before resuming the
        // scan. However, if an error is detected after we start
        // scanning a construct, we form a bad token out of the
        // characters that have already been scanned and resume
        // scanning on the character on which the problem was
        // detected. In other words, in that case, we do not advance.
        //
        if a.starttok == a.curtok:
            if a.current_kind == a.EOFT_SYMBOL:
                a.action = IntTuple(0)  // turn into garbage!
                return False

            a.lastToken = a.curtok
            a.tokStream.reportLexicalError(a.starttok, a.curtok)
            a.curtok = a.tokStream.getToken()
            if a.curtok > end_offset:
                a.curtok = a.tokStream.getStreamLength()
            a.current_kind = a.tokStream.getKind(a.curtok)

        else:
            a.lastToken = a.tokStream.getPrevious(a.curtok)
            a.tokStream.reportLexicalError(a.starttok, a.lastToken)

        return true

    //
    // This function takes as argument a configuration ([a.stack, stackTop], [a.tokStream, a.curtok]):
    // and determines whether or not the reduce a.action the a.curtok can be validly parsed in a
    // configuration.
    //
    def lexNextCharacter( act int, kind int):
        action_save = a.action.size()
        pos = a.stateStackTop,
        tempStackTop = a.stateStackTop - 1
        act = a.tAction(act, kind)
        // Scan:
        b_break_scan = False
        while act <= a.NUM_RULES:
            a.action.add(act)

            while true:
                lhs_symbol = a.prs.lhs(act)
                if lhs_symbol == a.START_SYMBOL:
                    if a.starttok == a.curtok:  // nil string reduction to a.START_SYMBOL is illegal

                        act = a.ERROR_ACTION
                        b_break_scan = true
                        break  // Scan

                    else:
                        a.parseActions()
                        return a.ACCEPT_ACTION

                tempStackTop -= (a.prs.rhs(act) - 1)
                state = (a.tempStack[tempStackTop] if tempStackTop > pos else a.stack[tempStackTop])
                act = a.prs.ntAction(state, lhs_symbol)

                if not act <= a.NUM_RULES:
                    break
            if b_break_scan:
                break

            if tempStackTop + 1 >= a.stack.__len__()
                a.reallocateStacks()
            //
            // ... Update the maximum useful position of the a.stack,
            // push goto state into (temporary): a.stack, and compute
            // the next a.action on the current symbol ...
            //
            pos = pos if pos < tempStackTop else tempStackTop
            a.tempStack[tempStackTop + 1] = act
            act = a.tAction(act, kind)

        //
        // If an error was detected, we restore the original configuration.
        // Otherwise, we update configuration up to the point prior to the
        // shift or shift-reduce on the token.
        //
        if act == a.ERROR_ACTION:
            a.action.reset(action_save)
        else:
            a.stateStackTop = tempStackTop + 1
            for i in range(pos + 1, a.stateStackTop + 1):  // update stack
                a.stack[i] = a.tempStack[i]

        return act

    //
    // Now do the final parse of the input based on the actions in
    // the list "a.action" and the sequence of tokens in the token stream.
    //
    def parseActions()
        //
        // Indicate that we are running the regular parser and that it's
        // ok to use the utility functions to query the parser.
        //
        a.taking_actions = true

        a.curtok = a.starttok
        a.lastToken = a.tokStream.getPrevious(a.curtok)

        //
        // Reparse the input...
        //
        a.stateStackTop = -1
        a.currentAction = a.START_STATE
        // process_actions:
        b_break_process_actions = False
        for i in range(0, a.action.size()):
            a.stateStackTop += 1
            a.stack[a.stateStackTop] = a.currentAction
            a.locationStack[a.stateStackTop] = a.curtok

            a.currentAction = a.action.get(i)
            if a.currentAction <= a.NUM_RULES:  // a reduce a.action?

                a.stateStackTop -= 1  // turn reduction into shift-reduction
                while true:
                    a.stateStackTop -= (a.prs.rhs(a.currentAction) - 1)
                    a.ra.ruleAction(a.currentAction)
                    lhs_symbol = a.prs.lhs(a.currentAction)
                    if lhs_symbol == a.START_SYMBOL:
                        // assert(starttok != a.curtok):  // nil string reduction to a.START_SYMBOL is illegal
                        b_break_process_actions = true
                        break  // process_actions

                    a.currentAction = a.prs.ntAction(a.stack[a.stateStackTop], lhs_symbol)
                    if not a.currentAction <= a.NUM_RULES:
                        break
            else:  // a shift or shift-reduce a.action

                a.lastToken = a.curtok
                a.curtok = a.tokStream.getNext(a.curtok)
                if a.currentAction > a.ERROR_ACTION:  // a shift-reduce a.action?

                    a.current_kind = a.tokStream.getKind(a.curtok)
                    a.currentAction -= a.ERROR_ACTION
                    while true:
                        a.stateStackTop -= (a.prs.rhs(a.currentAction) - 1)
                        a.ra.ruleAction(a.currentAction)
                        lhs_symbol = a.prs.lhs(a.currentAction)
                        if lhs_symbol == a.START_SYMBOL:
                            b_break_process_actions = true
                            break  // process_actions
                        a.currentAction = a.prs.ntAction(a.stack[a.stateStackTop], lhs_symbol)
                        if not a.currentAction <= a.NUM_RULES:
                            break

            if b_break_process_actions:
                break

        a.taking_actions = False  // indicate that we are done

        return
