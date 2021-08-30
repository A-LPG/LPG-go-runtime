package lpg2


type BacktrackingParser struct{
         *Stacks
         monitor Monitor
         START_STATE int
         NUM_RULES int
         NT_OFFSET int
         LA_STATE_OFFSET int
         EOFT_SYMBOL int
         ERROR_SYMBOL int
         ACCEPT_ACTION int
         ERROR_ACTION int 
    
         lastToken int
         currentAction int 
    
         tokStream TokenStream
         prs ParseTable 
         ra RuleAction 
    
         action *IntSegmentedTuple
         tokens *IntTuple
         actionStack []int
         skipTokens bool 
    
        //
        // A starting marker indicates that we are dealing with an entry point
        // for a given nonterminal. We need to execute a shift action on the
        // marker in order to parse the entry point in question.
        //
        markerTokenIndex int
}
func NewBacktrackingParser(tokStream TokenStream , prs ParseTable ,
    ra RuleAction , monitor Monitor)*BacktrackingParser {
    a := new(BacktrackingParser)
    a.Stacks = NewStacks()
    a.action = NewIntSegmentedTuple(10, 1024)
    a.reset(tokStream, prs, ra, monitor)
    a.skipTokens = false
    return a
}

func (this *BacktrackingParser) getMarkerToken(marker_kind int, start_token_index int) (int,error) {
    if marker_kind == 0 {
        return 0,nil
    } else {
        _ipsream,ok :=this.tokStream.(IPrsStream)
        if this.markerTokenIndex == 0 {

            if !ok {
                return -1, NewTokenStreamNotIPrsStreamException("")
            }
            this.markerTokenIndex = _ipsream.makeErrorToken(  this.tokStream.getPrevious(start_token_index),
                                                                    this.tokStream.getPrevious(start_token_index),
                                                                    this.tokStream.getPrevious(start_token_index),
                                                                    marker_kind)
        } else {
            _ipsream.getIToken(this.markerTokenIndex).setKind(marker_kind)
        }
    }
    return this.markerTokenIndex, nil
}

//
// Override the getToken function in Stacks.
//
func (this *BacktrackingParser) getToken(i int) int {
    return this.tokens.get(this.locationStack[this.stateStackTop + (i - 1)])
}

func (this *BacktrackingParser) getCurrentRule() int {
    return this.currentAction
}
func (this *BacktrackingParser) getFirstToken2() int {
    return this.tokStream.getFirstRealToken(this.getToken(1))
}
func (this *BacktrackingParser) getFirstToken(i int) int {

    return this.tokStream.getFirstRealToken(this.getToken(i))
}
func (this *BacktrackingParser) getLastToken2() int {
    return this.tokStream.getLastRealToken(this.lastToken)
}
func (this *BacktrackingParser) getLastToken(i  int) int {

    var l int 
    if i >= this.prs.rhs(this.currentAction){
        l = this.lastToken
    }else{
        l = this.tokens.get(this.locationStack[this.stateStackTop + i] - 1)
    }

    return this.tokStream.getLastRealToken(l)
}
func (this *BacktrackingParser) setMonitor(monitor Monitor)  {
    this.monitor = monitor
}
func (this *BacktrackingParser) reset1()  error{
    this.action.reset()
    this.skipTokens = false
    this.markerTokenIndex = 0
    return nil
}
func (this *BacktrackingParser) reset2(tokStream TokenStream, monitor Monitor) error {
    this.monitor = monitor
    this.tokStream = tokStream
    return this.reset1()
}

func (this *BacktrackingParser) reset(tokStream TokenStream , prs ParseTable , ra RuleAction , monitor Monitor) error {
    if prs!= nil {
        this.prs = prs
        this.START_STATE = prs.getStartState()
        this.NUM_RULES = prs.getNumRules()
        this.NT_OFFSET = prs.getNtOffset()
        this.LA_STATE_OFFSET = prs.getLaStateOffset()
        this.EOFT_SYMBOL = prs.getEoftSymbol()
        this.ERROR_SYMBOL = prs.getErrorSymbol()
        this.ACCEPT_ACTION = prs.getAcceptAction()
        this.ERROR_ACTION = prs.getErrorAction()
        if !prs.isValidForParser() {
            return NewBadParseSymFileException("")
        }
        if !prs.getBacktrack() {
            return NewNotBacktrackParseTableException("")
        }

    }
    if nil !=ra {
        this.ra = ra
    }

        
    if nil == tokStream {
        this.reset1()
        return nil
    }
    return this.reset2(tokStream, monitor)
}
func (this *BacktrackingParser) reset3(tokStream TokenStream, prs ParseTable, ra RuleAction) error {
   return this.reset(tokStream, prs, ra,nil)
}


//
// Allocate or reallocate all the stacks. Their sizes should always be the same.
//
func (this *BacktrackingParser) reallocateOtherStacks(start_token_index int)  {
    if len(this.actionStack) == 0 {
        this.actionStack = make([]int,len(this.stateStack))
        this.locationStack = make([]int,len(this.stateStack))
        this.parseStack =make([]interface{},len(this.stateStack))
        this.actionStack[0] = 0
        this.locationStack[0] = start_token_index
    } else {
        if len(this.actionStack) < len(this.stateStack) {
            var old_length int = len(this.actionStack)
            this.actionStack =arraycopy(this.actionStack, 0,  make([]int,len(this.stateStack)), 0, old_length)
            this.locationStack =arraycopy(this.locationStack, 0, make([]int,len(this.stateStack)), 0, old_length)
            this.parseStack =ObjectArraycopy(this.parseStack, 0, make([]interface{},len(this.stateStack)), 0, old_length)
        }
    }
    return
}
//func (this *BacktrackingParser) fuzzyParse() any {
//    return this.fuzzyParseEntry(0, lpg.lang.Integer.MAX_VALUE)
//}
func (this *BacktrackingParser) fuzzyParse(max_error_count int) (interface{},error) {
    return this.fuzzyParseEntry(0, max_error_count)
}
//func (this *BacktrackingParser) fuzzyParseEntry(marker_kind int) any {
//    return this.fuzzyParseEntry(marker_kind, lpg.lang.Integer.MAX_VALUE)
//}
func (this *BacktrackingParser) fuzzyParseEntry(marker_kind int, max_error_count int) (interface{},error) {

    this.action.reset()
    this.tokStream.reset()// Position at first token.
    this.reallocateStateStack()
    this.stateStackTop = 0
    this.stateStack[0] = this.START_STATE

    //
    // The tuple tokens will eventually contain the sequence 
    // of tokens that resulted in a successful parse. We leave
    // it up to the "Stream" implementer to define the predecessor
    // of the first token as he sees fit.
    //
    var first_token int = this.tokStream.peek()
    var start_token int = first_token
    marker_token ,_ := this.getMarkerToken(marker_kind, first_token)

    this.tokens = NewIntTupleWithEstimate(this.tokStream.getStreamLength())
    this.tokens.add(this.tokStream.getPrevious(first_token))

    var error_token int = this.backtrackParseInternal(this.action, marker_token)
    if error_token != 0 { // an error was detected?
        _ipsream,ok :=this.tokStream.(IPrsStream)
        if !ok {
            return nil,NewTokenStreamNotIPrsStreamException("")
        }
        var rp = NewRecoveryParser(this, this.action, this.tokens, _ipsream, this.prs, max_error_count, 0, this.monitor)
        start_token,_ = rp.recover(marker_token, error_token)
    }
    if marker_token != 0 && start_token == first_token {
        this.tokens.add(marker_token)
    }
    var t int = start_token
    for ; this.tokStream.getKind(t) != this.EOFT_SYMBOL; t = this.tokStream.getNext(t) {
        this.tokens.add(t)
    }
    this.tokens.add(t)
    return this.parseActions(marker_kind)
}

func (this *BacktrackingParser) parse(max_error_count int)  (interface{},error)  {
    return this.parseEntry(0, max_error_count)
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
func (this *BacktrackingParser) parseEntry(marker_kind int, max_error_count int)  (interface{},error)  {
    this.action.reset()
    this.tokStream.reset()// Position at first token.

    this.reallocateStateStack()
    this.stateStackTop = 0
    this.stateStack[0] = this.START_STATE

    this.skipTokens = max_error_count < 0
    _ipstream,ok := this.tokStream.(IPrsStream)
    if max_error_count > 0 &&  ok {
        max_error_count = 0
    }
    //
    // The tuple tokens will eventually contain the sequence 
    // of tokens that resulted in a successful parse. We leave
    // it up to the "Stream" implementer to define the predecessor
    // of the first token as he sees fit.
    //
    this.tokens = NewIntTupleWithEstimate(this.tokStream.getStreamLength())
    this.tokens.add(this.tokStream.getPrevious(this.tokStream.peek()))

    var start_token_index int = this.tokStream.peek()
    var repair_token,_ = this.getMarkerToken(marker_kind, start_token_index)
    var start_action_index int = this.action.size()// obviously 0
    var temp_stack []int = make([]int,this.stateStackTop + 1)
    arraycopy(this.stateStack, 0, temp_stack, 0, len(temp_stack))

    var initial_error_token = this.backtrackParseInternal(this.action, repair_token)
    var error_token int = initial_error_token
    var  count int = 0
    for ;error_token != 0;{
        if count == max_error_count {
            return nil,NewBadParseException(initial_error_token)
        }
        this.action.reset(start_action_index)
        this.tokStream.reset(start_token_index)
        this.stateStackTop = len(temp_stack) - 1
        arraycopy(temp_stack, 0, this.stateStack, 0, len(temp_stack))
        this.reallocateOtherStacks(start_token_index)


        this.backtrackParseUpToError(repair_token, error_token)

        for this.stateStackTop = this.findRecoveryStateIndex(this.stateStackTop);
            this.stateStackTop >= 0;
            this.stateStackTop = this.findRecoveryStateIndex(this.stateStackTop - 1) {
            var recovery_token = this.tokens.get(this.locationStack[this.stateStackTop] - 1)
            var temp int
            if recovery_token >= start_token_index {
                temp = recovery_token
            }else {
                temp = error_token
            }
            repair_token = this.errorRepair(_ipstream, temp, error_token)
            if repair_token != 0 {
                break
            }
        }
        if this.stateStackTop < 0 {
            return nil,NewBadParseException(initial_error_token)
        }
        temp_stack = make([]int,this.stateStackTop + 1)
        arraycopy(this.stateStack, 0, temp_stack, 0, len(temp_stack))

        start_action_index = this.action.size()
        start_token_index = this.tokStream.peek()

        error_token = this.backtrackParseInternal(this.action, repair_token)
        count++ 
    }
    if (repair_token != 0) {
        this.tokens.add(repair_token)
    }
    var t int= start_token_index
    for ; this.tokStream.getKind(t) != this.EOFT_SYMBOL ;t = this.tokStream.getNext(t) {
        this.tokens.add(t)
    }
    this.tokens.add(t)
    return this.parseActions(marker_kind)
}
//
// Process reductions and continue...
//
func (this *BacktrackingParser) process_reductions()  {
    for;; {
        this.stateStackTop -= (this.prs.rhs(this.currentAction) - 1)
        this.ra.ruleAction(this.currentAction)
        this.currentAction = this.prs.ntAction(this.stateStack[this.stateStackTop], this.prs.lhs(this.currentAction))
        if this.currentAction <= this.NUM_RULES{
            continue
        }else{
            break
        }
    } 
    return
}

//
// Now do the final parse of the input based on the actions in
// the list "action" and the sequence of tokens in list "tokens".
//
func(this *BacktrackingParser) parseActions(marker_kind int) (interface{}, error) {
    var ti int = -1
    ti+=1
    this.lastToken = this.tokens.get(ti)

    ti+=1
    var curtok = this.tokens.get(ti)
    this.allocateOtherStacks()
    //
    // Reparse the input...
    //
    this.stateStackTop = -1
    this.currentAction = this.START_STATE
    var i int = 0
    for ; i < this.action.size(); i++ {
        //
        // if the parser needs to stop processing, it may do so here.
        //
        if this.monitor!= nil  && this.monitor.isCancelled() {
            return nil,nil
        }
        this.stateStackTop+=1
        this.stateStack[this.stateStackTop] = this.currentAction

        this.locationStack[this.stateStackTop] = ti

        this.currentAction = this.action.get(i)
        if this.currentAction <= this.NUM_RULES { // a reduce action?
            this.stateStackTop--// make reduction look like shift-reduction
            this.process_reductions()
        } else { // a shift or shift-reduce action
            if this.tokStream.getKind(curtok) > this.NT_OFFSET {
                _ipstream,_ := this.tokStream.(IPrsStream)
                var badtok1  = _ipstream.getIToken(curtok)
                badtok,_ := badtok1.(*ErrorToken)
                return nil,NewBadParseException(badtok.getErrorToken().getTokenIndex())
            }
            this.lastToken = curtok
            ti+=1
            curtok = this.tokens.get(ti)
            if this.currentAction > this.ERROR_ACTION {
                this.currentAction -= this.ERROR_ACTION
                this.process_reductions()
            }
        }
    }

    if marker_kind == 0 {
        return this.parseStack[0 ],nil
    }else{
        return this.parseStack[1 ],nil
    }

}

//
// Process reductions and continue...
//
func (this *BacktrackingParser) process_backtrack_reductions(act int) int {
    for;; {
        this.stateStackTop -= (this.prs.rhs(act) - 1)
        act = this.prs.ntAction(this.stateStack[this.stateStackTop], this.prs.lhs(act))
        if act <= this.NUM_RULES {
            continue
        }else {
            break
        }
    }
    return act
}

//
// This method is intended to be used by the type RecoveryParser.
// Note that the action tuple passed here must be the same action
// tuple that was passed down to RecoveryParser. It is passed back
// to this method as documention.
//
func (this *BacktrackingParser) backtrackParse(stack []int, stack_top int, action* IntSegmentedTuple, initial_token int) int {
    this.stateStackTop = stack_top
    arraycopy(stack, 0, this.stateStack, 0, this.stateStackTop + 1)
    return this.backtrackParseInternal(action, initial_token)
}


//
// Parse the input until either the parse completes successfully or
// an error is encountered. This function returns an integer that
// represents the last action that was executed by the parser. If
// the parse was succesful, then the tuple "action" contains the
// successful sequence of actions that was executed.
//
func (this *BacktrackingParser) backtrackParseInternal(action *IntSegmentedTuple, initial_token int) int {
    //
    // Allocate configuration stack.
    //
    var configuration_stack  = NewConfigurationStack(this.prs)

    //
    // Keep parsing until we successfully reach the end of file or
    // an error is encountered. The list of actions executed will
    // be stored in the "action" tuple.
    //
    var error_token int = 0
    var maxStackTop int = this.stateStackTop
    var start_token int = this.tokStream.peek()

    var   curtok int
     if initial_token > 0 {
         curtok = initial_token
     } else{
         curtok =this.tokStream.getToken()
     }

    var    current_kind int = this.tokStream.getKind(curtok)
    var    act int = this.tAction(this.stateStack[this.stateStackTop], current_kind)
    //
    // The main driver loop
    //
    for;; {
        //
        // if the parser needs to stop processing,
        // it may do so here.
        //
        if this.monitor!=nil  && this.monitor.isCancelled() {
            return 0
        }
        if act <= this.NUM_RULES {
            action.add(act)// save this reduce action
            this.stateStackTop--
            act = this.process_backtrack_reductions(act)
        }else {
            if act > this.ERROR_ACTION {
                action.add(act) // save this shift-reduce action
                curtok = this.tokStream.getToken()
                current_kind = this.tokStream.getKind(curtok)
                act = this.process_backtrack_reductions(act - this.ERROR_ACTION)
            } else {
                if act < this.ACCEPT_ACTION {
                    action.add(act) // save this shift action
                    curtok = this.tokStream.getToken()
                    current_kind = this.tokStream.getKind(curtok)
                } else {
                    if act == this.ERROR_ACTION {

                        if !(error_token > curtok ){
                            error_token = curtok
                        }
                        var configuration = configuration_stack.pop()
                        if configuration == nil {
                            act = this.ERROR_ACTION
                        } else {
                            action.reset(configuration.action_length)
                            act = configuration.act
                            curtok = configuration.curtok
                            current_kind = this.tokStream.getKind(curtok)
                            var index int
                            if curtok == initial_token {
                                index = start_token
                            }else {
                                index=  this.tokStream.getNext(curtok)
                            }
                            this.tokStream.reset(index)
                            this.stateStackTop = configuration.stack_top
                            configuration.retrieveStack(this.stateStack)
                            continue
                        }
                        break
                    } else {
                        if act > this.ACCEPT_ACTION {
                            if configuration_stack.findConfiguration(this.stateStack, this.stateStackTop, curtok) {
                                act = this.ERROR_ACTION
                            } else {
                                configuration_stack.push(this.stateStack, this.stateStackTop, act+1, curtok, action.size())
                                act = this.prs.baseAction(act)
                                if this.stateStackTop > maxStackTop{
                                    maxStackTop = this.stateStackTop
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
        this.stateStackTop+=1
        if this.stateStackTop >= len(this.stateStack) {
            this.reallocateStateStack()
        }
        this.stateStack[this.stateStackTop] = act

        act = this.tAction(act, current_kind)
    }
    if  act == this.ERROR_ACTION {
        return  error_token
    } else {
        return 0
    }
}
func (this *BacktrackingParser) backtrackParseUpToError(initial_token int, error_token int)  {
    //
    // Allocate configuration stack.
    //
    var configuration_stack = NewConfigurationStack(this.prs)

    //
    // Keep parsing until we successfully reach the end of file or
    // an error is encountered. The list of actions executed will
    // be stored in the "action" tuple.
    //
    var start_token int = this.tokStream.peek()
    var    curtok int
    if initial_token > 0 {
        curtok = initial_token
    } else {
        curtok=this.tokStream.getToken()
    }
    var    current_kind int = this.tokStream.getKind(curtok)
    var   act int = this.tAction(this.stateStack[this.stateStackTop], current_kind)

    this.tokens.add(curtok)
    this.locationStack[this.stateStackTop] = this.tokens.size()
    this.actionStack[this.stateStackTop] = this.action.size()

    for ;; {
        //
        // if the parser needs to stop processing,
        // it may do so here.
        //
        if this.monitor!= nil  && this.monitor.isCancelled() {
            return
        }

        if act <= this.NUM_RULES {
            this.action.add(act)// save this reduce action
            this.stateStackTop--
            act = this.process_backtrack_reductions(act)
        }else{
             if (act > this.ERROR_ACTION){
                this.action.add(act) // save this shift-reduce action
                curtok = this.tokStream.getToken()
                current_kind = this.tokStream.getKind(curtok)
                this.tokens.add(curtok)
                act = this.process_backtrack_reductions(act - this.ERROR_ACTION)
            }else{
                 if (act < this.ACCEPT_ACTION){
                    this.action.add(act) // save this shift action
                    curtok = this.tokStream.getToken()
                    current_kind = this.tokStream.getKind(curtok)
                    this.tokens.add(curtok)
                }else{
                     if (act == this.ERROR_ACTION){
                        if (curtok != error_token) {
                            var configuration = configuration_stack.pop()
                            if configuration == nil {
                                act = this.ERROR_ACTION
                            } else {
                                this.action.reset(configuration.action_length)
                                act = configuration.act
                                var next_token_index int = configuration.curtok
                                this.tokens.reset(next_token_index)
                                curtok = this.tokens.get(next_token_index - 1)
                                current_kind = this.tokStream.getKind(curtok)
                                var index int
                                if curtok == initial_token {
                                    index = start_token
                                }else {
                                    index = this.tokStream.getNext(curtok)
                                }
                                this.tokStream.reset(index)

                                this.stateStackTop = configuration.stack_top
                                configuration.retrieveStack(this.stateStack)
                                this.locationStack[this.stateStackTop] = this.tokens.size()
                                this.actionStack[this.stateStackTop] = this.action.size()
                                continue
                            }
                        }
                        break
                    }else{
                         if (act > this.ACCEPT_ACTION){
                            if (configuration_stack.findConfiguration(this.stateStack, this.stateStackTop, this.tokens.size())) {
                                act = this.ERROR_ACTION
                            } else {
                                configuration_stack.push(this.stateStack, this.stateStackTop, act + 1, this.tokens.size(), this.action.size())
                                act = this.prs.baseAction(act)
                            }
                            continue
                        } else {
                            break// assert(act == ACCEPT_ACTION)
                        }
                    }
                }
            }
        }

        this.stateStackTop+=1
        this.stateStack[this.stateStackTop] = act// no need to check if out of bounds
        this.locationStack[this.stateStackTop] = this.tokens.size()
        this.actionStack[this.stateStackTop] = this.action.size()
        act = this.tAction(act, current_kind)
    }
    return
}
func (this *BacktrackingParser) repairable(error_token int) bool {
    //
    // Allocate configuration stack.
    //
    var configuration_stack  = NewConfigurationStack(this.prs)

    //
    // Keep parsing until we successfully reach the end of file or
    // an error is encountered. The list of actions executed will
    // be stored in the "action" tuple.
    //
    var start_token int = this.tokStream.peek()
    var final_token int = this.tokStream.getStreamLength() // unreachable
    var curtok int = 0
    var current_kind int = this.ERROR_SYMBOL
    var act int = this.tAction(this.stateStack[this.stateStackTop], current_kind)
    for ;;{
        if (act <= this.NUM_RULES) {
            this.stateStackTop--
            act = this.process_backtrack_reductions(act)
        }else{
             if (act > this.ERROR_ACTION){
                curtok = this.tokStream.getToken()
                if (curtok > final_token) {
                    return true
                }
                current_kind = this.tokStream.getKind(curtok)
                act = this.process_backtrack_reductions(act - this.ERROR_ACTION)
            }else{
                 if (act < this.ACCEPT_ACTION){
                    curtok = this.tokStream.getToken()
                    if (curtok > final_token) {
                        return true
                    }
                    current_kind = this.tokStream.getKind(curtok)
                }else{
                     if (act == this.ERROR_ACTION){
                        var configuration = configuration_stack.pop()
                        if (configuration == nil) {
                            act = this.ERROR_ACTION
                        } else {
                            this.stateStackTop = configuration.stack_top
                            configuration.retrieveStack(this.stateStack)
                            act = configuration.act
                            curtok = configuration.curtok
                            if (curtok == 0) {
                                current_kind = this.ERROR_SYMBOL
                                this.tokStream.reset(start_token)
                            } else {
                                current_kind = this.tokStream.getKind(curtok)
                                this.tokStream.reset(this.tokStream.getNext(curtok))
                            }
                            continue
                        }
                        break
                    }else{
                         if (act > this.ACCEPT_ACTION){
                            if (configuration_stack.findConfiguration(this.stateStack, this.stateStackTop, curtok)) {
                                act = this.ERROR_ACTION
                            } else {
                                configuration_stack.push(this.stateStack, this.stateStackTop, act + 1, curtok, 0)
                                act = this.prs.baseAction(act)
                            }
                            continue
                        } else {
                            break// assert(act == ACCEPT_ACTION)
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
        if ((curtok > error_token) && (final_token == this.tokStream.getStreamLength())) {
            //
            // If the ERROR_SYMBOL is a valid Action Adjunct in the state
            // "act" then we set the terminating token as the successor of
            // the current token. I.e., we have to be able to parse at least
            // two tokens past the resynch point before we claim victory.
            //
            if this.recoverableState(act) {
                if this.skipTokens {
                    final_token = curtok
                }else{
                    final_token = this.tokStream.getNext(curtok)
                }

            }
        }
        this.stateStackTop+=1
        if this.stateStackTop >= len(this.stateStack) {
            this.reallocateStateStack()
        }
        this.stateStack[this.stateStackTop] = act

        act = this.tAction(act, current_kind)
    }
    //
    // If we can reach the end of the input successfully, we claim victory.
    //
    return (act == this.ACCEPT_ACTION)
}
func (this *BacktrackingParser) recoverableState(state int) bool {
        var k int = this.prs.asi(state)
    for ; this.prs.asr(k) != 0; k++ {
        if (this.prs.asr(k) == this.ERROR_SYMBOL) {
            return true
        }
    }
    return false
}
func (this *BacktrackingParser) findRecoveryStateIndex(start_index int) int {
    var i int = start_index
    for ;i >= 0; i-- {
        //
        // If the ERROR_SYMBOL is an Action Adjunct in state stateStack[i]
        // then chose i as the index of the state to recover on.
        //
        if (this.recoverableState(this.stateStack[i])) {
            break
        }
    }
    if (i >= 0) {// if a recoverable state, remove null reductions, if any.
        var k = i - 1
        for ; k >= 0; k-- {
            if (this.locationStack[k] != this.locationStack[i]) {
                break
            }
        }
        i = k + 1
    }
    return i
}

func (this *BacktrackingParser) errorRepair(stream IPrsStream, recovery_token int, error_token int) int {
    var temp_stack []int = make([]int,this.stateStackTop + 1)
    arraycopy(this.stateStack, 0, temp_stack, 0, len(temp_stack))
    for ;
        stream.getKind(recovery_token) != this.EOFT_SYMBOL;
        recovery_token = stream.getNext(recovery_token){
        stream.reset(recovery_token)
        if (this.repairable(error_token)) {
            break
        }
        this.stateStackTop = len(temp_stack) - 1
        arraycopy(temp_stack, 0, this.stateStack, 0, len(temp_stack))
    }

    if (stream.getKind(recovery_token) == this.EOFT_SYMBOL) {
        stream.reset(recovery_token)
        if (!this.repairable(error_token)) {
            this.stateStackTop = len(temp_stack) - 1
            arraycopy(temp_stack, 0, this.stateStack, 0, len(temp_stack))
            return 0
        }
    }

    this.stateStackTop = len(temp_stack) - 1
    arraycopy(temp_stack, 0, this.stateStack, 0, len(temp_stack))
    stream.reset(recovery_token)
    this.tokens.reset(this.locationStack[this.stateStackTop] - 1)
    this.action.reset(this.actionStack[this.stateStackTop])

    return stream.makeErrorToken(this.tokens.get(this.locationStack[this.stateStackTop] - 1),
                                                stream.getPrevious(recovery_token),
                                                error_token,
                                                this.ERROR_SYMBOL)
}

//
// keep looking ahead until we compute a valid action
//
func (this *BacktrackingParser) lookahead(act int, token int) int {
    act = this.prs.lookAhead(act - this.LA_STATE_OFFSET, this.tokStream.getKind(token))
    if act > this.LA_STATE_OFFSET{
        return this.lookahead(act, this.tokStream.getNext(token))
    }else{
        return  act
    }

}

//
// Compute the next action defined on act and sym. If this
// action requires more lookahead, these lookahead symbols
// are in the token stream beginning at the next token that
// is yielded by peek().
//
func (this *BacktrackingParser) tAction(act int, sym int) int {
    act = this.prs.tAction(act, sym)
    if act > this.LA_STATE_OFFSET{
        return this.lookahead(act, this.tokStream.peek())
    }else{
        return  act
    }
}
    