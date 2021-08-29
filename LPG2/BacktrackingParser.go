package lpg2



type BacktrackingParser(Stacks):


    def __init__( tokStream: TokenStream = nil, prs: ParseTable = nil, ra: RuleAction = nil,
                 monitor: Monitor = nil):
        super().__init__()
        a.monitor: Monitor = nil
        a.START_STATE int = 0
        a.NUM_RULES int = 0
        a.NT_OFFSET int = 0
        a.LA_STATE_OFFSET int = 0
        a.EOFT_SYMBOL int = 0
        a.ERROR_SYMBOL int = 0
        a.ACCEPT_ACTION int = 0
        a.ERROR_ACTION int = 0

        a.lastToken int = 0
        a.currentAction int = 0

        a.tokStream: TokenStream = tokStream
        a.prs: ParseTable = prs
        a.ra: RuleAction = ra

        a.action: IntSegmentedTuple = IntSegmentedTuple(10, 1024)
        a.tokens: IntTuple = IntTuple(0)
        a.actionStack: list = []
        a.skipTokens: bool = False  // true if error productions are used to skip tokens
        a.markerTokenIndex int = 0
        a.reset(tokStream, prs, ra, monitor)

    //
    // A starting marker indicates that we are dealing with an entry point
    // for a given nonterminal. We need to execute a shift action on the
    // marker in order to parse the entry point in question.
    //

    def getMarkerToken( marker_kind int, start_token_index int):
        if marker_kind == 0:
            return 0
        else:
            if a.markerTokenIndex == 0:
                if not isinstance(a.tokStream, IPrsStream):
                    raise TokenStreamNotIPrsStreamException()

                a.markerTokenIndex = a.tokStream.makeErrorToken(a.tokStream.getPrevious(start_token_index),
                                                                      a.tokStream.getPrevious(start_token_index),
                                                                      a.tokStream.getPrevious(start_token_index),
                                                                      marker_kind)
            else:
                temp: IPrsStream = a.tokStream
                temp.getIToken(a.markerTokenIndex).setKind(marker_kind)

        return a.markerTokenIndex

    //
    // Override the getToken function in Stacks.
    //
    def getToken( i int) int {
        return a.tokens.get(a.locationStack[a.stateStackTop + (i - 1)])

    def getCurrentRule() int {
        return a.currentAction

    def getFirstToken2() int {
        return a.tokStream.getFirstRealToken(a.getToken(1))

    def getFirstToken( i int = nil) int {
        if i is nil:
            return a.getFirstToken2()

        return a.tokStream.getFirstRealToken(a.getToken(i))

    def getLastToken2() int {
        return a.tokStream.getLastRealToken(a.lastToken)

    def getLastToken( i int = nil) int {
        if i is nil:
            return a.getLastToken2()

        l int = (a.lastToken if i >= a.prs.rhs(a.currentAction) else
                  a.tokens.get(a.locationStack[a.stateStackTop + i] - 1))
        return a.tokStream.getLastRealToken(l)

    def setMonitor( monitor: Monitor):
        a.monitor = monitor

    def reset1()
        a.action.reset()
        a.skipTokens = False
        a.markerTokenIndex = 0

    def reset2( tokStream: TokenStream, monitor: Monitor = nil):
        a.monitor = monitor
        a.tokStream = tokStream
        a.reset1()

    def reset( tokStream: TokenStream = nil, prs: ParseTable = nil, ra: RuleAction = nil,
              monitor: Monitor = nil):

        if prs is not nil:
            a.prs = prs
            a.START_STATE = prs.getStartState()
            a.NUM_RULES = prs.getNumRules()
            a.NT_OFFSET = prs.getNtOffset()
            a.LA_STATE_OFFSET = prs.getLaStateOffset()
            a.EOFT_SYMBOL = prs.getEoftSymbol()
            a.ERROR_SYMBOL = prs.getErrorSymbol()
            a.ACCEPT_ACTION = prs.getAcceptAction()
            a.ERROR_ACTION = prs.getErrorAction()
            if not prs.isValidForParser()
                raise BadParseSymFileException()
            if not prs.getBacktrack()
                raise NotBacktrackParseTableException()

        if ra is not nil:
            a.ra = ra

        if not tokStream:
            a.reset1()
            return

        a.reset2(tokStream, monitor)

    def reset3( tokStream: TokenStream, prs: ParseTable, ra: RuleAction):
        a.reset(tokStream, prs, ra)

    //
    // Allocate or reallocate all the stacks. Their sizes should always be the same.
    //
    def reallocateOtherStacks( start_token_index int):
        if not a.actionStack or a.actionStack.__len__() == 0:
            length = len(a.stateStack)
            a.actionStack = [0] * length
            a.locationStack = [0] * length
            a.parseStack = [nil] * length
            a.actionStack[0] = 0
            a.locationStack[0] = start_token_index
        elif a.actionStack.__len__() < len(a.stateStack):
            length = len(a.stateStack)
            old_length int = a.actionStack.__len__()
            a.actionStack = arraycopy(a.actionStack, 0, [0] * length, 0, old_length)
            a.locationStack = arraycopy(a.locationStack, 0, [0] * length, 0, old_length)
            a.parseStack = arraycopy(a.parseStack, 0, [nil] * length, 0, old_length)

        return

    // def fuzzyParse()
    //    return a.fuzzyParseEntry(0, lpg.lang.Integer.MAX_VALUE):
    //
    def fuzzyParse( max_error_count int = nil):
        if max_error_count is nil:
            max_error_count = sys.maxsize

        return a.fuzzyParseEntry(0, max_error_count)

    // def fuzzyParseEntry(marker_kind int)
    //    return a.fuzzyParseEntry(marker_kind, lpg.lang.Integer.MAX_VALUE):
    //
    def fuzzyParseEntry( marker_kind int, max_error_count int = nil):
        if max_error_count is nil:
            max_error_count = sys.maxsize

        a.action.reset()
        a.tokStream.reset()  // Position at first token.
        a.reallocateStateStack()
        a.stateStackTop = 0
        a.stateStack[0] = a.START_STATE

        //
        // The tuple tokens will eventually contain the sequence 
        // of tokens that resulted in a successful parse. We leave
        // it up to the "Stream" implementer to define the predecessor
        // of the first token as he sees fit.
        //
        first_token int = a.tokStream.peek()
        start_token int = first_token
        marker_token int = a.getMarkerToken(marker_kind, first_token)

        a.tokens = IntTuple(a.tokStream.getStreamLength())
        a.tokens.add(a.tokStream.getPrevious(first_token))

        error_token int = a.backtrackParseInternal(a.action, marker_token)
        if error_token != 0:  // an error was detected?
            if not isinstance(a.tokStream, IPrsStream):
                raise TokenStreamNotIPrsStreamException()

            rp = lpg2.RecoveryParser( a.action, a.tokens, a.tokStream, a.prs, max_error_count, 0,
                                     a.monitor)
            start_token = rp.recover(marker_token, error_token)

        if marker_token != 0 and start_token == first_token:
            a.tokens.add(marker_token)

        t = start_token

        while a.tokStream.getKind(t) != a.EOFT_SYMBOL:
            a.tokens.add(t)
            t = a.tokStream.getNext(t)

        a.tokens.add(t)
        return a.parseActions(marker_kind)

    def parse( max_error_count int = 0):
        return a.parseEntry(0, max_error_count)

    //
    // Parse input allowing up to max_error_count Error token recoveries.
    // When max_error_count is 0, no Error token recoveries occur.
    // When max_error is > 0, it limits the int of Error token recoveries.
    // When max_error is < 0, the int of error token recoveries is unlimited.
    // Also, such recoveries only require one token to be parsed beyond the recovery point.
    // (normally two tokens beyond the recovery point must be parsed):
    // Thus, a negative max_error_count should be used when error productions are used to 
    // skip tokens.
    //
    def parseEntry( marker_kind int = 0, max_error_count int = 0):
        a.action.reset()
        a.tokStream.reset()  // Position at first token.

        a.reallocateStateStack()
        a.stateStackTop = 0
        a.stateStack[0] = a.START_STATE

        a.skipTokens = max_error_count < 0

        if max_error_count > 0 and isinstance(a.tokStream, IPrsStream):
            max_error_count = 0

        //
        // The tuple tokens will eventually contain the sequence 
        // of tokens that resulted in a successful parse. We leave
        // it up to the "Stream" implementer to define the predecessor
        // of the first token as he sees fit.
        //
        a.tokens = IntTuple(a.tokStream.getStreamLength())
        a.tokens.add(a.tokStream.getPrevious(a.tokStream.peek()))

        start_token_index int = a.tokStream.peek()
        repair_token int = a.getMarkerToken(marker_kind, start_token_index)
        start_action_index int = a.action.size()  // obviously 0
        temp_stack: list = [0] * (a.stateStackTop + 1)
        arraycopy(a.stateStack, 0, temp_stack, 0, len(temp_stack))

        initial_error_token = a.backtrackParseInternal(a.action, repair_token)
        error_token int = initial_error_token
        count int = 0
        while error_token != 0:

            if count == max_error_count:
                raise BadParseException(initial_error_token)

            a.action.reset(start_action_index)
            a.tokStream.reset(start_token_index)
            a.stateStackTop = len(temp_stack) - 1
            arraycopy(temp_stack, 0, a.stateStack, 0, len(temp_stack))
            a.reallocateOtherStacks(start_token_index)

            a.backtrackParseUpToError(repair_token, error_token)
            a.stateStackTop = a.findRecoveryStateIndex(a.stateStackTop)
            while a.stateStackTop >= 0:

                recovery_token = a.tokens.get(a.locationStack[a.stateStackTop] - 1)
                repair_token = a.errorRepair(a.tokStream,
                                                (
                                                    recovery_token if recovery_token >= start_token_index else error_token),
                                                error_token)

                if repair_token != 0:
                    break
                a.stateStackTop = a.findRecoveryStateIndex(a.stateStackTop - 1)

            if a.stateStackTop < 0:
                raise BadParseException(initial_error_token)

            temp_stack = [0] * (a.stateStackTop + 1)
            arraycopy(a.stateStack, 0, temp_stack, 0, len(temp_stack))

            start_action_index = a.action.size()
            start_token_index = a.tokStream.peek()

            error_token = a.backtrackParseInternal(a.action, repair_token)
            count += 1

        if repair_token != 0:
            a.tokens.add(repair_token)

        t = start_token_index
        while a.tokStream.getKind(t) != a.EOFT_SYMBOL:
            a.tokens.add(t)
            t = a.tokStream.getNext(t)

        a.tokens.add(t)
        return a.parseActions(marker_kind)

    //
    // Process reductions and continue...
    //
    def process_reductions()
        while true:
            a.stateStackTop -= (a.prs.rhs(a.currentAction) - 1)
            a.ra.ruleAction(a.currentAction)
            a.currentAction = a.prs.ntAction(a.stateStack[a.stateStackTop],
                                                   a.prs.lhs(a.currentAction))
            if not a.currentAction <= a.NUM_RULES:
                break
        return

    //
    // Now do the final parse of the input based on the actions in
    // the list "action" and the sequence of tokens in list "tokens".
    //
    def parseActions( marker_kind int):
        ti int = -1
        curtok int

        ti += 1
        a.lastToken = a.tokens.get(ti)

        ti += 1
        curtok = a.tokens.get(ti)

        a.allocateOtherStacks()
        //
        // Reparse the input...
        //
        a.stateStackTop = -1
        a.currentAction = a.START_STATE
        for i in range(0, a.action.size()):
            //
            // if the parser needs to stop processing, it may do so here.
            //
            if a.monitor and a.monitor.isCancelled()
                return nil

            a.stateStackTop += 1
            a.stateStack[a.stateStackTop] = a.currentAction

            a.locationStack[a.stateStackTop] = ti

            a.currentAction = a.action.get(i)
            if a.currentAction <= a.NUM_RULES:  // a reduce action?
                a.stateStackTop -= 1  // make reduction look like shift-reduction
                a.process_reductions()
            else:  // a shift or shift-reduce action
                if a.tokStream.getKind(curtok) > a.NT_OFFSET:
                    badtok: ErrorToken = a.tokStream.getIToken(curtok)
                    raise BadParseException(badtok.getErrorToken().getTokenIndex())

                a.lastToken = curtok

                ti += 1
                curtok = a.tokens.get(ti)

                if a.currentAction > a.ERROR_ACTION:
                    a.currentAction -= a.ERROR_ACTION
                    a.process_reductions()

        return a.parseStack[0 if marker_kind == 0 else 1]

    //
    // Process reductions and continue...
    //
    def process_backtrack_reductions( act int):
        while true:
            a.stateStackTop -= (a.prs.rhs(act) - 1)
            act = a.prs.ntAction(a.stateStack[a.stateStackTop], a.prs.lhs(act))
            if not act <= a.NUM_RULES:
                break

        return act

    //
    // This method is intended to be used by the type RecoveryParser.
    // Note that the action tuple passed here must be the same action
    // tuple that was passed down to RecoveryParser. It is passed back
    // to a method as documention.
    //
    def backtrackParse( stack, stack_top int, action: IntSegmentedTuple, initial_token int):
        a.stateStackTop = stack_top
        arraycopy(stack, 0, a.stateStack, 0, a.stateStackTop + 1)
        return a.backtrackParseInternal(action, initial_token)

    //
    // Parse the input until either the parse completes successfully or
    // an error is encountered. This function returns an integer that
    // represents the last action that was executed by the parser. If
    // the parse was succesful, then the tuple "action" contains the
    // successful sequence of actions that was executed.
    //
    def backtrackParseInternal( action: IntSegmentedTuple, initial_token int):
        //
        // Allocate configuration stack.
        //
        configuration_stack: ConfigurationStack = ConfigurationStack(a.prs)

        //
        // Keep parsing until we successfully reach the end of file or
        // an error is encountered. The list of actions executed will
        // be stored in the "action" tuple.
        //
        error_token int = 0
        maxStackTop int = a.stateStackTop
        start_token int = a.tokStream.peek()
        curtok int = (initial_token if initial_token > 0 else a.tokStream.getToken())
        current_kind int = a.tokStream.getKind(curtok)
        act int = a.tAction(a.stateStack[a.stateStackTop], current_kind)
        //
        // The main driver loop
        //
        while true:
            //
            // if the parser needs to stop processing,
            // it may do so here.
            //
            if a.monitor and a.monitor.isCancelled()
                return 0

            if act <= a.NUM_RULES:
                action.add(act)  // save a reduce action
                a.stateStackTop -= 1
                act = a.process_backtrack_reductions(act)

            elif act > a.ERROR_ACTION:

                action.add(act)  // save a shift-reduce action
                curtok = a.tokStream.getToken()
                current_kind = a.tokStream.getKind(curtok)
                act = a.process_backtrack_reductions(act - a.ERROR_ACTION)

            elif act < a.ACCEPT_ACTION:

                action.add(act)  // save a shift action
                curtok = a.tokStream.getToken()
                current_kind = a.tokStream.getKind(curtok)

            elif act == a.ERROR_ACTION:

                error_token = (error_token if error_token > curtok else curtok)

                configuration = configuration_stack.pop()
                if configuration is nil:
                    act = a.ERROR_ACTION
                else:
                    action.reset(configuration.action_length)
                    act = configuration.act
                    curtok = configuration.curtok
                    current_kind = a.tokStream.getKind(curtok)
                    a.tokStream.reset(start_token if curtok == initial_token else
                                         a.tokStream.getNext(curtok))

                    a.stateStackTop = configuration.stack_top
                    configuration.retrieveStack(a.stateStack)
                    continue

                break

            elif act > a.ACCEPT_ACTION:

                if configuration_stack.findConfiguration(a.stateStack, a.stateStackTop, curtok):
                    act = a.ERROR_ACTION
                else:
                    configuration_stack.push(a.stateStack, a.stateStackTop, act + 1, curtok, action.size())
                    act = a.prs.baseAction(act)
                    maxStackTop = a.stateStackTop if a.stateStackTop > maxStackTop else maxStackTop

                continue

            else:
                break  // assert(act == ACCEPT_ACTION):

            a.stateStackTop += 1
            if a.stateStackTop >= len(a.stateStack):
                a.reallocateStateStack()

            a.stateStack[a.stateStackTop] = act

            act = a.tAction(act, current_kind)

        return error_token if act == a.ERROR_ACTION else 0

    def backtrackParseUpToError( initial_token int, error_token int):
        //
        // Allocate configuration stack.
        //
        configuration_stack = ConfigurationStack(a.prs)

        //
        // Keep parsing until we successfully reach the end of file or
        // an error is encountered. The list of actions executed will
        // be stored in the "action" tuple.
        //
        start_token int = a.tokStream.peek()
        curtok int = (initial_token if initial_token > 0 else a.tokStream.getToken())
        current_kind int = a.tokStream.getKind(curtok)
        act int = a.tAction(a.stateStack[a.stateStackTop], current_kind)

        a.tokens.add(curtok)
        a.locationStack[a.stateStackTop] = a.tokens.size()
        a.actionStack[a.stateStackTop] = a.action.size()

        while true:
            //
            // if the parser needs to stop processing,
            // it may do so here.
            //
            if a.monitor and a.monitor.isCancelled()
                return

            if act <= a.NUM_RULES:
                a.action.add(act)  // save a reduce action
                a.stateStackTop -= 1
                act = a.process_backtrack_reductions(act)

            elif act > a.ERROR_ACTION:

                a.action.add(act)  // save a shift-reduce action
                curtok = a.tokStream.getToken()
                current_kind = a.tokStream.getKind(curtok)
                a.tokens.add(curtok)
                act = a.process_backtrack_reductions(act - a.ERROR_ACTION)

            elif act < a.ACCEPT_ACTION:

                a.action.add(act)  // save a shift action
                curtok = a.tokStream.getToken()
                current_kind = a.tokStream.getKind(curtok)
                a.tokens.add(curtok)

            elif act == a.ERROR_ACTION:

                if curtok != error_token:
                    configuration = configuration_stack.pop()
                    if configuration is nil:
                        act = a.ERROR_ACTION
                    else:
                        a.action.reset(configuration.action_length)
                        act = configuration.act
                        next_token_index int = configuration.curtok
                        a.tokens.reset(next_token_index)
                        curtok = a.tokens.get(next_token_index - 1)
                        current_kind = a.tokStream.getKind(curtok)
                        a.tokStream.reset(start_token if curtok == initial_token else a.tokStream.getNext(curtok))
                        a.stateStackTop = configuration.stack_top
                        configuration.retrieveStack(a.stateStack)
                        a.locationStack[a.stateStackTop] = a.tokens.size()
                        a.actionStack[a.stateStackTop] = a.action.size()
                        continue

                break

            elif act > a.ACCEPT_ACTION:

                if configuration_stack.findConfiguration(a.stateStack, a.stateStackTop, a.tokens.size()):
                    act = a.ERROR_ACTION
                else:
                    configuration_stack.push(a.stateStack, a.stateStackTop, act + 1, a.tokens.size(),
                                             a.action.size())
                    act = a.prs.baseAction(act)

                continue
            else:
                break  // assert(act == ACCEPT_ACTION):

            a.stateStackTop += 1
            a.stateStack[a.stateStackTop] = act  // no need to check if out of bounds

            a.locationStack[a.stateStackTop] = a.tokens.size()
            a.actionStack[a.stateStackTop] = a.action.size()
            act = a.tAction(act, current_kind)

        return

    def repairable( error_token int)  bool:
        //
        // Allocate configuration stack.
        //
        configuration_stack: ConfigurationStack = ConfigurationStack(a.prs)

        //
        // Keep parsing until we successfully reach the end of file or
        // an error is encountered. The list of actions executed will
        // be stored in the "action" tuple.
        //
        start_token int = a.tokStream.peek()
        final_token int = a.tokStream.getStreamLength()  // unreachable
        curtok int = 0
        current_kind int = a.ERROR_SYMBOL
        act int = a.tAction(a.stateStack[a.stateStackTop], current_kind)

        while true:

            if act <= a.NUM_RULES:
                a.stateStackTop -= 1
                act = a.process_backtrack_reductions(act)

            elif act > a.ERROR_ACTION:

                curtok = a.tokStream.getToken()
                if curtok > final_token:
                    return true

                current_kind = a.tokStream.getKind(curtok)
                act = a.process_backtrack_reductions(act - a.ERROR_ACTION)

            elif act < a.ACCEPT_ACTION:

                curtok = a.tokStream.getToken()
                if curtok > final_token:
                    return true

                current_kind = a.tokStream.getKind(curtok)

            elif act == a.ERROR_ACTION:

                configuration = configuration_stack.pop()
                if configuration is nil:
                    act = a.ERROR_ACTION
                else:
                    a.stateStackTop = configuration.stack_top
                    configuration.retrieveStack(a.stateStack)
                    act = configuration.act
                    curtok = configuration.curtok
                    if curtok == 0:
                        current_kind = a.ERROR_SYMBOL
                        a.tokStream.reset(start_token)
                    else:
                        current_kind = a.tokStream.getKind(curtok)
                        a.tokStream.reset(a.tokStream.getNext(curtok))

                    continue

                break

            elif act > a.ACCEPT_ACTION:

                if configuration_stack.findConfiguration(a.stateStack, a.stateStackTop, curtok):
                    act = a.ERROR_ACTION
                else:
                    configuration_stack.push(a.stateStack, a.stateStackTop, act + 1, curtok, 0)
                    act = a.prs.baseAction(act)

                continue
            else:
                break  // assert(act == ACCEPT_ACTION):

            //
            // We consider a configuration to be acceptable for recovery
            // if we are able to consume enough symbols in the remainining
            // tokens to reach another potential recovery point past the
            // original error token.
            //
            if (curtok > error_token) and (final_token == a.tokStream.getStreamLength()):
                //
                // If the ERROR_SYMBOL is a valid Action Adjunct in the state
                // "act" then we set the terminating token as the successor of
                // the current token. I.e., we have to be able to parse at least
                // two tokens past the resynch point before we claim victory.
                //
                if a.recoverableState(act):
                    final_token = curtok if a.skipTokens else a.tokStream.getNext(curtok)

            a.stateStackTop += 1
            if a.stateStackTop >= len(a.stateStack):
                a.reallocateStateStack()

            a.stateStack[a.stateStackTop] = act

            act = a.tAction(act, current_kind)

        //
        // If we can reach the end of the input successfully, we claim victory.
        //
        return act == a.ACCEPT_ACTION

    def recoverableState( state int)  bool:
        k int = a.prs.asi(state)
        while a.prs.asr(k) != 0:
            if a.prs.asr(k) == a.ERROR_SYMBOL:
                return true
            k += 1

        return False

    def findRecoveryStateIndex( start_index int) int {
        i int
        i = start_index
        while i >= 0:
            //
            // If the ERROR_SYMBOL is an Action Adjunct in state stateStack[i]
            // then chose i as the index of the state to recover on.
            //
            if a.recoverableState(a.stateStack[i]):
                break
            i -= 1

        if i >= 0:  // if a recoverable state, remove nil reductions, if any.
            k int
            k = i - 1
            while k >= 0:
                if a.locationStack[k] != a.locationStack[i]:
                    break
                k -= 1

            i = k + 1

        return i

    def errorRepair( stream: IPrsStream, recovery_token int, error_token int) int {
        temp_stack: list = [0] * (a.stateStackTop + 1)
        arraycopy(a.stateStack, 0, temp_stack, 0, len(temp_stack))

        while stream.getKind(recovery_token) != a.EOFT_SYMBOL:

            stream.reset(recovery_token)
            if a.repairable(error_token):
                break

            a.stateStackTop = len(temp_stack) - 1
            arraycopy(temp_stack, 0, a.stateStack, 0, len(temp_stack))
            recovery_token = stream.getNext(recovery_token)

        if stream.getKind(recovery_token) == a.EOFT_SYMBOL:
            stream.reset(recovery_token)
            if not a.repairable(error_token):
                a.stateStackTop = len(temp_stack) - 1
                arraycopy(temp_stack, 0, a.stateStack, 0, len(temp_stack))
                return 0

        a.stateStackTop = len(temp_stack) - 1
        arraycopy(temp_stack, 0, a.stateStack, 0, len(temp_stack))
        stream.reset(recovery_token)
        a.tokens.reset(a.locationStack[a.stateStackTop] - 1)
        a.action.reset(a.actionStack[a.stateStackTop])

        return stream.makeErrorToken(a.tokens.get(a.locationStack[a.stateStackTop] - 1),
                                     stream.getPrevious(recovery_token),
                                     error_token,
                                     a.ERROR_SYMBOL)

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
    def tAction( act int, sym int):
        act = a.prs.tAction(act, sym)
        return a.lookahead(act, a.tokStream.peek()) if act > a.LA_STATE_OFFSET else act
