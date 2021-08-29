package lpg2


type RecoveryParser(DiagnoseParser):

    //
    // maxErrors is the maximum int of errors to be repaired
    // maxTime is the maximum amount of time allowed for diagnosing
    // but at least one error must be diagnosed 
    //
    def __init__( parser: BacktrackingParser, action: IntSegmentedTuple, tokens: IntTuple, tokStream: IPrsStream,
                 prs: ParseTable, maxErrors int = 0, maxTime int = 0, monitor: Monitor = nil):
        super().__init__(tokStream, prs, maxErrors, maxTime, monitor)
        a.parser: BacktrackingParser = parser
        a.action: IntSegmentedTuple = action
        a.tokens: IntTuple = tokens
        a.actionStack: list = []
        a.scope_repair: PrimaryRepairInfo = PrimaryRepairInfo()

    def reallocateStacks()
        super().reallocateStacks()
        if a.actionStack is nil or a.actionStack.__len__() == 0:
            a.actionStack = [0] * (len(a.stateStack))
        else:
            old_stack_length int = a.actionStack.__len__()
            a.actionStack = arraycopy(a.actionStack, 0, [0] * (len(a.stateStack)), 0, old_stack_length)

        return

    def reportError( scope_index int, error_token int):
        text: string = "\""

        i int = a.scopeSuffix(scope_index)
        while a.scopeRhs(i) != 0:
            if not a.isNullable(a.scopeRhs(i)):
                symbol_index int = (a.nonterminalIndex(a.scopeRhs(i) - a.NT_OFFSET)
                                     if a.scopeRhs(i) > a.NT_OFFSET
                                     else a.terminalIndex(a.scopeRhs(i)))
                if a.name(symbol_index).__len__() > 0:
                    if text.__len__() > 1:  // Not just starting quote?
                        text += " "  // add a space separator

                    text += a.name(symbol_index)
            i += 1

        text += "\""
        a.tokStream.reportError(ParseErrorCodes.SCOPE_CODE, error_token, error_token, [text])
        return

    def recover( marker_token int, error_token int) int {
        if not a.stateStack or len(a.stateStack) == 0:
            a.reallocateStacks()

        a.tokens.reset()
        a.tokStream.reset()
        a.tokens.add(a.tokStream.getPrevious(a.tokStream.peek()))
        restart_token int = (marker_token if marker_token != 0 else a.tokStream.getToken())
        old_action_size int = 0
        a.stateStackTop = 0
        a.stateStack[a.stateStackTop] = a.START_STATE
        while true:
            a.action.reset(old_action_size)
            if not a.fixError(restart_token, error_token):
                raise BadParseException(error_token)

            //
            // if the parser needs to stop processing,
            // it may do so here.
            //
            if a.monitor and a.monitor.isCancelled()
                break

            //
            // At a stage, we have a recovery configuration. See how
            // far we can go with it.
            //
            restart_token = error_token
            a.tokStream.reset(error_token)
            old_action_size = a.action.size()  // save the old size in case we encounter a  error
            error_token = a.parser.backtrackParse(a.stateStack, a.stateStackTop, a.action, 0)
            a.tokStream.reset(a.tokStream.getNext(restart_token))
            if not error_token != 0:  // no error found
                break
        return restart_token

    //
    // Given the configuration consisting of the states in stateStack
    // and the sequence of tokens (current_kind, followed by the tokens
    // in tokStream):, fixError parses up to error_token in the tokStream
    // recovers, if possible, from that error and returns the result.
    // While doing  it also computes the location_stack information
    // and the sequence of actions that matches up with the result that
    // it returns.
    //
    def fixError( start_token int, error_token int)  bool:
        //
        // Save information about the current configuration.
        //
        curtok int = start_token
        current_kind int = a.tokStream.getKind(curtok)
        first_stream_token int = a.tokStream.peek()

        a.buffer[1] = error_token
        a.buffer[0] = a.tokStream.getPrevious(a.buffer[1])
        for k in range(2, a.BUFF_SIZE):
            a.buffer[k] = a.tokStream.getNext(a.buffer[k - 1])

        a.scope_repair.distance = 0
        a.scope_repair.misspellIndex = 0
        a.scope_repair.bufferPosition = 1

        //
        // Clear the configuration stack.
        //
        a.main_configuration_stack = ConfigurationStack(a.prs)

        //
        // Keep parsing until we reach the end of file and succeed or
        // an error is encountered. The list of actions executed will
        // be stored in the "action" tuple.
        //
        a.locationStack[a.stateStackTop] = curtok
        a.actionStack[a.stateStackTop] = a.action.size()
        act int = a.tAction(a.stateStack[a.stateStackTop], current_kind)
        while true:
            //
            // if the parser needs to stop processing,
            // it may do so here.
            //
            if a.monitor and a.monitor.isCancelled()
                return true

            if act <= a.NUM_RULES:
                a.action.add(act)  // save this reduce action
                a.stateStackTop -= 1

                while true:
                    a.stateStackTop -= (a.rhs(act) - 1)
                    act = a.ntAction(a.stateStack[a.stateStackTop], a.lhs(act))
                    if not act <= a.NUM_RULES:
                        break

                a.stateStackTop += 1
                if a.stateStackTop >= len(a.stateStack):
                    a.reallocateStacks()
                a.stateStack[a.stateStackTop] = act

                a.locationStack[a.stateStackTop] = curtok
                a.actionStack[a.stateStackTop] = a.action.size()
                act = a.tAction(act, current_kind)
                continue

            elif act == a.ERROR_ACTION:

                if curtok != error_token or a.main_configuration_stack.size() > 0:
                    configuration = a.main_configuration_stack.pop()
                    if configuration is nil:
                        act = a.ERROR_ACTION
                    else:
                        a.stateStackTop = configuration.stack_top
                        configuration.retrieveStack(a.stateStack)
                        act = configuration.act
                        curtok = configuration.curtok
                        a.action.reset(configuration.action_length)
                        current_kind = a.tokStream.getKind(curtok)
                        a.tokStream.reset(a.tokStream.getNext(curtok))
                        continue

                break

            elif act > a.ACCEPT_ACTION and act < a.ERROR_ACTION:

                if a.main_configuration_stack.findConfiguration(a.stateStack, a.stateStackTop, curtok):
                    act = a.ERROR_ACTION
                else:
                    a.main_configuration_stack.push(a.stateStack, a.stateStackTop, act + 1, curtok,
                                                       a.action.size())
                    act = a.baseAction(act)

                continue

            else:

                if act < a.ACCEPT_ACTION:

                    a.action.add(act)  // save a shift action
                    curtok = a.tokStream.getToken()
                    current_kind = a.tokStream.getKind(curtok)

                elif act > a.ERROR_ACTION:

                    a.action.add(act)  // save a shift-reduce action
                    curtok = a.tokStream.getToken()
                    current_kind = a.tokStream.getKind(curtok)
                    act -= a.ERROR_ACTION
                    while true:
                        a.stateStackTop -= (a.rhs(act) - 1)
                        act = a.ntAction(a.stateStack[a.stateStackTop], a.lhs(act))
                        if not act <= a.NUM_RULES:
                            break

                else:
                    break  // assert(act == ACCEPT_ACTION):  THIS IS NOT SUPPOSED TO HAPPEN!!!

                a.stateStackTop += 1
                if a.stateStackTop >= len(a.stateStack):
                    a.reallocateStacks()
                a.stateStack[a.stateStackTop] = act

                if curtok == error_token:
                    a.scopeTrial(a.scope_repair, a.stateStack, a.stateStackTop)
                    if a.scope_repair.distance >= a.MIN_DISTANCE:

                        a.tokens.add(start_token)
                        token int = first_stream_token
                        while token != error_token:
                            a.tokens.add(token)
                            token = a.tokStream.getNext(token)

                        a.acceptRecovery(error_token)
                        break  // equivalent to: return true;

                a.locationStack[a.stateStackTop] = curtok
                a.actionStack[a.stateStackTop] = a.action.size()
                act = a.tAction(act, current_kind)

        return act != a.ERROR_ACTION

    def acceptRecovery( error_token int):
        //
        //
        //
        // int action_size = action.size()

        //
        // Simulate parsing actions required for a sequence of scope
        // recoveries.
        // TODO: need to add action and fix the location_stack?
        //
        recovery_action: IntTuple = IntTuple()
        for k in range(0, a.scopeStackTop + 1):
            scope_index int = a.scopeIndex[k]
            la int = a.scopeLa(scope_index)

            //
            // Compute the action (or set of actions in case of conflicts): that
            // can be executed on the scope lookahead symbol. Save the action(s):
            // in the action tuple.
            //
            recovery_action.reset()
            act int = a.tAction(a.stateStack[a.stateStackTop], la)
            if act > a.ACCEPT_ACTION and act < a.ERROR_ACTION:  // conflicting actions?
                while true:

                    recovery_action.add(a.baseAction(act))
                    act += 1

                    if not a.baseAction(act) != 0:
                        break
            else:
                recovery_action.add(act)

            //
            // For each action defined on the scope lookahead symbol,
            // try scope recovery. At least one action should succeed!
            //
            start_action_size int = a.action.size()
            index int
            for index in range(0, recovery_action.size()):

                //
                // Reset the action tuple each time through a loop
                // to clear previous actions that may have been added
                // because of a failed call to completeScope.
                //
                a.action.reset(start_action_size)
                a.tokStream.reset(error_token)
                a.tempStackTop = a.stateStackTop - 1
                max_pos int = a.stateStackTop

                act = recovery_action.get(index)
                while act <= a.NUM_RULES:
                    a.action.add(act)  // save a reduce action
                    //
                    // ... Process all goto-reduce actions following
                    // reduction, until a goto action is computed ...
                    //
                    while true:
                        lhs_symbol int = a.lhs(act)
                        a.tempStackTop -= (a.rhs(act) - 1)
                        act = (a.tempStack[a.tempStackTop]
                               if a.tempStackTop > max_pos else a.stateStack[a.tempStackTop])
                        act = a.ntAction(act, lhs_symbol)
                        if not act <= a.NUM_RULES:
                            break
                    if a.tempStackTop + 1 >= len(a.stateStack):
                        a.reallocateStacks()

                    max_pos = max_pos if max_pos < a.tempStackTop else a.tempStackTop
                    a.tempStack[a.tempStackTop + 1] = act
                    act = a.tAction(act, la)

                //
                // If the lookahead symbol is parsable, then we check
                // whether or not we have a match between the scope
                // prefix and the transition symbols corresponding to
                // the states on top of the stack.
                //
                if act != a.ERROR_ACTION:
                    a.tempStackTop += 1
                    a.nextStackTop = a.tempStackTop

                    for i in range(0, max_pos + 1):
                        a.nextStack[i] = a.stateStack[i]

                    //
                    // NOTE that we do not need to update location_stack and
                    // actionStack here because, once the rules associated with
                    // these scopes are reduced, all these states will be popped
                    // from the stack.
                    //
                    for i in range(max_pos + 1, a.tempStackTop + 1):
                        a.nextStack[i] = a.tempStack[i]

                    if a.completeScope(a.action, a.scopeSuffix(scope_index)):
                        i int = a.scopeSuffix(a.scopeIndex[k])
                        while a.scopeRhs(i) != 0:
                            a.tokens.add(a.tokStream.makeErrorToken
                                            (error_token,
                                             a.tokStream.getPrevious(error_token),
                                             error_token, a.scopeRhs(i)))
                            i += 1

                        a.reportError(a.scopeIndex[k], a.tokStream.getPrevious(error_token))
                        break

            // assert (index < recovery_action.size()): // sanity check!
            a.stateStackTop = a.nextStackTop
            arraycopy(a.nextStack, 0, a.stateStack, 0, a.stateStackTop + 1)

        return

    def completeScope( action: IntSegmentedTuple, scope_rhs_index int)  bool:
        kind int = a.scopeRhs(scope_rhs_index)
        if kind == 0:
            return true

        act int = a.nextStack[a.nextStackTop]

        if kind > a.NT_OFFSET:
            lhs_symbol int = kind - a.NT_OFFSET
            if a.baseCheck(act + lhs_symbol) != lhs_symbol:
                // is there a valid
                // action defined on
                // lhs_symbol?
                return False

            act = a.ntAction(act, lhs_symbol)

            //
            // if action is a goto-reduce action, save it as a shift-reduce
            // action.
            //
            action.add(act + a.ERROR_ACTION if act <= a.NUM_RULES else act)
            while act <= a.NUM_RULES:
                a.nextStackTop -= (a.rhs(act) - 1)
                act = a.ntAction(a.nextStack[a.nextStackTop], a.lhs(act))

            a.nextStackTop += 1
            a.nextStack[a.nextStackTop] = act
            return a.completeScope(action, scope_rhs_index + 1)

        //
        // Processing a terminal
        //
        act = a.tAction(act, kind)
        action.add(act)  // save a terminal action
        if act < a.ACCEPT_ACTION:
            a.nextStackTop += 1
            a.nextStack[a.nextStackTop] = act
            return a.completeScope(action, scope_rhs_index + 1)

        elif act > a.ERROR_ACTION:
            act -= a.ERROR_ACTION
            while true:
                a.nextStackTop -= (a.rhs(act) - 1)
                act = a.ntAction(a.nextStack[a.nextStackTop], a.lhs(act))
                if not act <= a.NUM_RULES:
                    break

            a.nextStackTop += 1
            a.nextStack[a.nextStackTop] = act
            return true

        elif act > a.ACCEPT_ACTION and act < a.ERROR_ACTION:  // conflicting actions?

            save_action_size int = action.size()
            i int = act
            while a.baseAction(i) != 0:  // consider only shift and shift-reduce actions

                action.reset(save_action_size)
                act = a.baseAction(i)
                action.add(act)  // save a terminal action
                if act <= a.NUM_RULES:  // Ignore reduce actions
                    pass
                elif act < a.ACCEPT_ACTION:

                    a.nextStackTop += 1
                    a.nextStack[a.nextStackTop] = act
                    if a.completeScope(action, scope_rhs_index + 1):
                        return true

                elif act > a.ERROR_ACTION:

                    act -= a.ERROR_ACTION
                    while true:
                        a.nextStackTop -= (a.rhs(act) - 1)
                        act = a.ntAction(a.nextStack[a.nextStackTop], a.lhs(act))
                        if not act <= a.NUM_RULES:
                            break
                    a.nextStackTop += 1
                    a.nextStack[a.nextStackTop] = act
                    return true
                i += 1

        return False
