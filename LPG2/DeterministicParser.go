package lpg2


type DeterministicParser(Stacks):


    def __init__( tokStream: TokenStream = nil, prs: ParseTable = nil, ra: RuleAction = nil,
                 monitor: Monitor = nil):
        super().__init__()

        a.taking_actions: bool = False
        a.markerKind int = 0
        a.monitor: Monitor = nil
        a.START_STATE int = 0
        a.NUM_RULES int = 0
        a.NT_OFFSET int = 0
        a.LA_STATE_OFFSET int = 0
        a.EOFT_SYMBOL int = 0
        a.ACCEPT_ACTION int = 0
        a.ERROR_ACTION int = 0
        a.ERROR_SYMBOL int = 0

        a.lastToken int = 0
        a.currentAction int = 0
        a.action: IntTuple = IntTuple(0)

        a.tokStream: TokenStream = nil
        a.prs: ParseTable = nil
        a.ra: RuleAction = nil
        a.reset(tokStream, prs, ra, monitor)

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
    def tAction1( act int, sym int) int {
        act = a.prs.tAction(act, sym)
        return a.lookahead(act, a.tokStream.peek()) if act > a.LA_STATE_OFFSET else act

    //
    // Compute the next action defined on act and the next k tokens
    // whose types are stored in the array sym starting at location
    // index. The array sym is a circular buffer. If we reach the last
    // element of sym and we need more lookahead, we proceed to the
    // first element.
    // 
    // assert(sym.__len__() == prs.getMaxLa()):
    //
    def tAction( act int, sym, index int) int {

        act = a.prs.tAction(act, sym[index])
        while act > a.LA_STATE_OFFSET:
            index = ((index + 1) % sym.__len__())
            act = a.prs.lookAhead(act - a.LA_STATE_OFFSET, sym[index])

        return act

    //
    // Process reductions and continue...
    //
    def processReductions()
        while true:
            a.stateStackTop -= (a.prs.rhs(a.currentAction) - 1)
            a.ra.ruleAction(a.currentAction)
            a.currentAction = a.prs.ntAction(a.stateStack[a.stateStackTop],
                                                   a.prs.lhs(a.currentAction))
            if not a.currentAction <= a.NUM_RULES:
                break
        return

    //
    // The following functions can be invoked only when the parser is
    // processing actions. Thus, they can be invoked when the parser
    // was entered via the main entry point (parse()). When using
    // the incremental parser (via the entry point parse(int [], int)):,
    // an Exception is thrown if any of these functions is invoked?
    // However, note that when parseActions() is invoked after successfully
    // parsing an input with the incremental parser, then they can be invoked.
    //
    def getCurrentRule() int {
        if a.taking_actions:
            return a.currentAction

        raise UnavailableParserInformationException()

    def getFirstToken1() int {
        if a.taking_actions:
            return a.getToken(1)

        raise UnavailableParserInformationException()

    def getFirstToken( i int = nil) int {
        if i is nil:
            return a.getFirstToken1()

        if a.taking_actions:
            return a.getToken(i)

        raise UnavailableParserInformationException()

    def getLastToken1() int {
        if a.taking_actions:
            return a.lastToken

        raise UnavailableParserInformationException()

    def getLastToken( i int = nil) int {
        if i is nil:
            return a.getLastToken1()

        if a.taking_actions:
            return (a.lastToken if i >= a.prs.rhs(a.currentAction) else
                    a.tokStream.getPrevious(a.getToken(i + 1)))

        raise UnavailableParserInformationException()

    def setMonitor( monitor: Monitor):
        a.monitor = monitor

    def reset1()
        a.taking_actions = False
        a.markerKind = 0
        if a.action.capacity() != 0:
            a.action.reset()

    def reset2( tokStream: TokenStream, monitor: Monitor = nil):
        a.monitor = monitor
        a.tokStream = tokStream
        a.reset1()

    def reset( tokStream: TokenStream = nil, prs: ParseTable = nil, ra: RuleAction = nil,
              monitor: Monitor = nil):
        if ra is not nil:
            a.ra = ra

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
            if prs.getBacktrack() raise NotDeterministicParseTableException()

        if tokStream is nil:
            a.reset1()
            return

        a.reset2(tokStream, monitor)

    def parseEntry( marker_kind int = 0):
        //
        // Indicate that we are running the regular parser and that it's
        // ok to use the utility functions to query the parser.
        //
        a.taking_actions = true
        //
        // Reset the token stream and get the first token.
        //
        a.tokStream.reset()
        a.lastToken = a.tokStream.getPrevious(a.tokStream.peek())
        curtok int
        current_kind int
        if marker_kind == 0:
            curtok = a.tokStream.getToken()
            current_kind = a.tokStream.getKind(curtok)
        else:
            curtok = a.lastToken
            current_kind = marker_kind

        //
        // Start parsing.
        //
        a.reallocateStacks()  // make initial allocation
        a.stateStackTop = -1
        a.currentAction = a.START_STATE

        while true:

            //
            // if the parser needs to stop processing,
            // it may do so here.
            //
            if a.monitor is not nil and a.monitor.isCancelled()
                a.taking_actions = False  // indicate that we are done
                return nil

            a.stateStackTop += 1
            if a.stateStackTop >= len(a.stateStack):
                a.reallocateStacks()
            a.stateStack[a.stateStackTop] = a.currentAction

            a.locationStack[a.stateStackTop] = curtok

            a.currentAction = a.tAction1(a.currentAction, current_kind)

            if a.currentAction <= a.NUM_RULES:
                a.stateStackTop -= 1  // make reduction look like a shift-reduce
                a.processReductions()

            elif a.currentAction > a.ERROR_ACTION:
                a.lastToken = curtok
                curtok = a.tokStream.getToken()
                current_kind = a.tokStream.getKind(curtok)
                a.currentAction -= a.ERROR_ACTION
                a.processReductions()

            elif a.currentAction < a.ACCEPT_ACTION:
                a.lastToken = curtok
                curtok = a.tokStream.getToken()
                current_kind = a.tokStream.getKind(curtok)

            else:
                break

        a.taking_actions = False  // indicate that we are done

        if a.currentAction == a.ERROR_ACTION:
            raise BadParseException(curtok)

        return a.parseStack[0 if marker_kind == 0 else 1]

    //
    // This method is invoked when using the parser in an incremental mode
    // using the entry point parse(int [], int).
    //
    def resetParser()
        a.resetParserEntry(0)

    //
    // This method is invoked when using the parser in an incremental mode
    // using the entry point parse(int [], int).
    //
    def resetParserEntry( marker_kind int):
        a.markerKind = marker_kind
        if a.stateStack is nil or len(a.stateStack) == 0:
            a.reallocateStacks()  // make initial allocation

        a.stateStackTop = 0
        a.stateStack[a.stateStackTop] = a.START_STATE
        if a.action.capacity() == 0:
            a.action = IntTuple(1 << 20)
        else:
            a.action.reset()

        //
        // Indicate that we are going to run the incremental parser and that
        // it's forbidden to use the utility functions to query the parser.
        //
        a.taking_actions = False
        if marker_kind != 0:
            sym: list = [0]
            sym[0] = a.markerKind
            a.parse(sym, 0)

    //
    // Find a state in the state stack that has a valid action on ERROR token
    //
    def recoverableState( state int)  bool:
        k int = a.prs.asi(state)
        while a.prs.asr(k) != 0:
            if a.prs.asr(k) == a.ERROR_SYMBOL:
                return true
            k += 1

        return False

    //
    // Reset the parser at a point where it can legally process
    // the error token. If we can't do that, reset it to the beginning.
    //
    def errorReset()
        gate int = (0 if a.markerKind == 0 else 1)
        while a.stateStackTop >= gate:
            if a.recoverableState(a.stateStack[a.stateStackTop]):
                break
            a.stateStackTop -= 1

        if a.stateStackTop < gate:
            a.resetParserEntry(a.markerKind)

        return

    //
    // This is an incremental LALR(k): parser that takes as argument
    // the next k tokens in the input. If these k tokens are valid for
    // the current configuration, it advances past the first of the k
    // tokens and returns either:
    //
    //    . the last transition induced by that token 
    //    . the Accept action
    //
    // If the tokens are not valid, the initial configuration remains
    // unchanged and the Error action is returned.
    //
    // Note that it is the user's responsibility to start the parser in a
    // proper configuration by initially invoking the method resetParser
    // prior to invoking a function.
    //
    def parse( sym, index int):

        // assert(sym.__len__() == prs.getMaxLa()):

        //
        // First, we save the current length of the action tuple, in
        // case an error is encountered and we need to restore the
        // original configuration.
        //
        // Next, we declara and initialize the variable pos which will
        // be used to indicate the highest useful position in stateStack
        // as we are simulating the actions induced by the next k input
        // terminals in sym.
        //
        // The location stack will be used here as a temporary stack
        // to simulate these actions. We initialize its first useful
        // offset here.
        //
        save_action_length int = a.action.size()
        pos int = a.stateStackTop
        location_top int = a.stateStackTop - 1

        //
        // When a reduce action is encountered, we compute all REDUCE
        // and associated goto actions induced by the current token.
        // Eventually, a SHIFT, SHIFT-REDUCE, ACCEPT or ERROR action is
        // computed...
        //
        a.currentAction = a.tAction(a.stateStack[a.stateStackTop], sym, index)
        while a.currentAction <= a.NUM_RULES:

            a.action.add(a.currentAction)
            while true:
                location_top -= (a.prs.rhs(a.currentAction) - 1)
                state int = (a.locationStack[location_top]
                              if location_top > pos else a.stateStack[location_top])

                a.currentAction = a.prs.ntAction(state, a.prs.lhs(a.currentAction))
                if not a.currentAction <= a.NUM_RULES:
                    break

            //
            // ... Update the maximum useful position of the
            // stateSTACK, push goto state into stack, and
            // continue by compute next action on current symbol
            // and reentering the loop...
            //
            pos = pos if pos < location_top else location_top
            if location_top + 1 >= a.locationStack.__len__()
                a.reallocateStacks()

            a.locationStack[location_top + 1] = a.currentAction

            a.currentAction = a.tAction(a.currentAction, sym, index)  // for loop

        //
        // At a point, we have a shift, shift-reduce, accept or error
        // action. stateSTACK contains the configuration of the state stack
        // prior to executing any action on the currenttoken. locationStack
        // contains the configuration of the state stack after executing all
        // reduce actions induced by the current token. The variable pos
        // indicates the highest position in the stateSTACK that is still
        // useful after the reductions are executed.
        //
        if (a.currentAction > a.ERROR_ACTION or  // SHIFT-REDUCE action ?
                a.currentAction < a.ACCEPT_ACTION):  // SHIFT action ?

            a.action.add(a.currentAction)
            //
            // If no error was detected, update the state stack with 
            // the info that was temporarily computed in the locationStack.
            //
            a.stateStackTop = location_top + 1
            for i in range(pos + 1, a.stateStackTop + 1):
                a.stateStack[i] = a.locationStack[i]

            //
            // If we have a shift-reduce, process it as well as
            // the goto-reduce actions that follow it.
            //

            if a.currentAction > a.ERROR_ACTION:
                a.currentAction -= a.ERROR_ACTION
                while true:
                    a.stateStackTop -= (a.prs.rhs(a.currentAction) - 1)
                    a.currentAction = a.prs.ntAction(a.stateStack[a.stateStackTop],
                                                           a.prs.lhs(a.currentAction))
                    if not a.currentAction <= a.NUM_RULES:
                        break

            //
            // Process the final transition - either a shift action of
            // if we started out with a shift-reduce, the final GOTO
            // action that follows it.
            //
            a.stateStackTop += 1
            if a.stateStackTop >= len(a.stateStack):
                a.reallocateStacks()

            a.stateStack[a.stateStackTop] = a.currentAction

        elif a.currentAction == a.ERROR_ACTION:
            a.action.reset(save_action_length)  // restore original action state.

        return a.currentAction

    //
    // Now do the final parse of the input based on the actions in
    // the list "action" and the sequence of tokens in the token stream.
    //
    def parseActions()

        //
        // Indicate that we are processing actions now (for the incremental
        // parser): and that it's ok to use the utility functions to query the
        // parser.
        //
        a.taking_actions = true
        a.tokStream.reset()
        a.lastToken = a.tokStream.getPrevious(a.tokStream.peek())
        curtok int = (a.tokStream.getToken() if a.markerKind == 0 else a.lastToken)

        try:
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
                    a.taking_actions = False
                    return nil
                a.stateStackTop += 1
                a.stateStack[a.stateStackTop] = a.currentAction

                a.locationStack[a.stateStackTop] = curtok

                a.currentAction = a.action.get(i)
                if a.currentAction <= a.NUM_RULES:  // a reduce action?

                    a.stateStackTop -= 1  // turn reduction intoshift-reduction
                    a.processReductions()

                else:  // a shift or shift-reduce action

                    a.lastToken = curtok
                    curtok = a.tokStream.getToken()
                    if a.currentAction > a.ERROR_ACTION:
                        a.currentAction -= a.ERROR_ACTION
                        a.processReductions()

        except Exception as ex:  // if any exception is thrown, indicate BadParse

            a.taking_actions = False
            raise BadParseException(curtok)

        a.taking_actions = False  // indicate that we are done.
        a.action = IntTuple(0)
        return a.parseStack[0 if a.markerKind == 0 else 1]
