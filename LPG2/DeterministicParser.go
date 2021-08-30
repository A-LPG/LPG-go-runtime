package lpg2


type DeterministicParser struct{
     *Stacks
     taking_actions bool
     markerKind int 

     monitor Monitor
     START_STATE int
     NUM_RULES int
     NT_OFFSET int
     LA_STATE_OFFSET int
     EOFT_SYMBOL int
     ACCEPT_ACTION int
     ERROR_ACTION int
     ERROR_SYMBOL int 

     lastToken int
     currentAction int
     action *IntTuple 

     tokStream TokenStream 
     prs ParseTable 
     ra RuleAction 
}

func NewDeterministicParser(tokStream TokenStream , prs ParseTable ,ra RuleAction , monitor Monitor) *DeterministicParser{
       
        a := new(DeterministicParser)
        a.Stacks = NewStacks()
        
        a.reset(tokStream, prs, ra, monitor)

        return a
}
//
// keep looking ahead until we compute a valid action
//
func (this *DeterministicParser) lookahead(act int, token int) int {
    act = this.prs.lookAhead(act - this.LA_STATE_OFFSET, this.tokStream.getKind(token))
    if act > this.LA_STATE_OFFSET{
        return this.lookahead(act, this.tokStream.getNext(token))
    }else{
        return act
    }

}
//
// Compute the next action defined on act and sym. If this
// action requires more lookahead, these lookahead symbols
// are in the token stream beginning at the next token that
// is yielded by peek().
//
func (this *DeterministicParser) tAction1(act int, sym int) int {
    act = this.prs.tAction(act, sym)
    if act > this.LA_STATE_OFFSET{
        return this.lookahead(act, this.tokStream.peek())
    }else{
        return  act
    }

}
//
// Compute the next action defined on act and the next k tokens
// whose types are stored in the array sym starting at location
// index. The array sym is a circular buffer. If we reach the last
// element of sym and we need more lookahead, we proceed to the
// first element.
// 
// assert(sym.length == prs.getMaxLa())
//
func (this *DeterministicParser) tAction(act int, sym []int , index int) int {
    
    act = this.prs.tAction(act, sym[index])
    for;act > this.LA_STATE_OFFSET; {
        index = (index + 1) % len(sym)
        act = this.prs.lookAhead(act - this.LA_STATE_OFFSET, sym[index])
    }
    return act
}
//
// Process reductions and continue...
//
func (this *DeterministicParser)   processReductions()  {
    for;; {
        this.stateStackTop -= (this.prs.rhs(this.currentAction) - 1)
        this.ra.ruleAction(this.currentAction)
        this.currentAction = this.prs.ntAction( this.stateStack[this.stateStackTop],
                                                this.prs.lhs(this.currentAction))
        if this.currentAction <= this.NUM_RULES {
          continue
        }else {
          break
        }
    }
    return
}
//
// The following functions can be invoked only when the parser is
// processing actions. Thus, they can be invoked when the parser
// was entered via the main entry point (parse()). When using
// the incremental parser (via the entry point parse(int [], int)),
// an Exception is thrown if any of these functions is invoked?
// However, note that when parseActions() is invoked after successfully
// parsing an input with the incremental parser, then they can be invoked.
//
func (this *DeterministicParser)  getCurrentRule() (int ,error){
    if (this.taking_actions) {
        return this.currentAction,nil
    }
    return -1,NewUnavailableParserInformationException("")
}
func (this *DeterministicParser)  getFirstToken1() (int ,error) {
    if (this.taking_actions) {
        return this.getToken(1),nil
    }
    return -1,NewUnavailableParserInformationException("")
}
func (this *DeterministicParser)  getFirstToken(i int) (int ,error) {

    if (this.taking_actions) {
        return this.getToken(i),nil
    }
    return -1, NewUnavailableParserInformationException("")
}
func (this *DeterministicParser)  getLastToken1() (int ,error) {
    if (this.taking_actions) {
        return this.lastToken,nil
    }
    return -1,NewUnavailableParserInformationException("")
}
func (this *DeterministicParser)  getLastToken(i int) (int ,error)  {

    if this.taking_actions {
        if i >= this.prs.rhs(this.currentAction) {
            return this.lastToken,nil
        }else{
            return   this.tokStream.getPrevious(this.getToken(i + 1)),nil
        }
    }
    return -1 ,NewUnavailableParserInformationException("")
}
func (this *DeterministicParser)  setMonitor(monitor Monitor)  {
    this.monitor = monitor
}
func (this *DeterministicParser)  reset1()  {
    this.taking_actions = false
    this.markerKind = 0
    if this.action == nil {
        this.action.reset()
    }
}
func (this *DeterministicParser)  reset2(tokStream TokenStream,monitor Monitor)  {
    this.monitor = monitor
    this.tokStream = tokStream
    this.reset1()
}

func (this *DeterministicParser)  reset(tokStream TokenStream  , prs ParseTable , ra RuleAction  , monitor Monitor) error  {
    if nil != ra {
        this.ra = ra
    }
    if nil !=prs {
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
                return  NewBadParseSymFileException("")
            }
            if prs.getBacktrack() {
                return NewNotDeterministicParseTableException("")
            }
    }
    if nil == tokStream {
        this.reset1()
        return nil
    }
    this.reset2(tokStream, monitor)
    return nil
}



func (this *DeterministicParser)  parseEntry(marker_kind int) (interface{},error) {
    //
    // Indicate that we are running the regular parser and that it's
    // ok to use the utility functions to query the parser.
    //
    this.taking_actions = true
    //
    // Reset the token stream and get the first token.
    //
    this.tokStream.reset()
    this.lastToken = this.tokStream.getPrevious(this.tokStream.peek())
    var curtok int
    var current_kind int
    if marker_kind == 0 {
        curtok = this.tokStream.getToken()
        current_kind = this.tokStream.getKind(curtok)
    } else {
        curtok = this.lastToken
        current_kind = marker_kind
    }
    //
    // Start parsing.
    //
    this.reallocateStacks()// make initial allocation
    this.stateStackTop = -1
    this.currentAction = this.START_STATE

   // processTerminals:
        for ;; {
        //
        // if the parser needs to stop processing,
        // it may do so here.
        //
        if this.monitor != nil && this. monitor.isCancelled() {
            this.taking_actions = false // indicate that we are done
            return nil,nil
        }

        this.stateStackTop+=1
        if (this.stateStackTop >= len(this.stateStack)) {
            this.reallocateStacks()
        }

        this.stateStack[this.stateStackTop] = this.currentAction

        this.locationStack[this.stateStackTop] = curtok

        this.currentAction = this.tAction1(this.currentAction, current_kind)

        if this.currentAction <= this. NUM_RULES {
            this. stateStackTop-- // make reduction look like a shift-reduce
            this.  processReductions()
        }else{
             if this.currentAction > this.ERROR_ACTION {
                this.lastToken = curtok
                curtok = this.tokStream.getToken()
                current_kind = this.tokStream.getKind(curtok)
                this.currentAction -= this.ERROR_ACTION
                this.processReductions()
            }else{
                 if this.currentAction < this.ACCEPT_ACTION {
                    this. lastToken = curtok
                    curtok = this.tokStream.getToken()
                    current_kind = this.tokStream.getKind(curtok)
                } else {
                    break //processTerminals
                }
            }
        }
    }

    this. taking_actions = false // indicate that we are done

    if this.currentAction == this.ERROR_ACTION {
        return nil,NewBadParseException(curtok)
    }


    if marker_kind == 0 {
        return   this.parseStack[0],nil
    }else{
        return   this.parseStack[1],nil
    }
}
//
// This method is invoked when using the parser in an incremental mode
// using the entry point parse(int [], int).
//
func (this *DeterministicParser)  resetParser()  {
    this.resetParserEntry(0)
}
//
// This method is invoked when using the parser in an incremental mode
// using the entry point parse(int [], int).
//
func (this *DeterministicParser)  resetParserEntry(marker_kind int)  {
    this.markerKind = marker_kind
    if (this.stateStack == nil || len(this.stateStack) == 0) {
        this.reallocateStacks()// make initial allocation
    }
    this.stateStackTop = 0
    this.stateStack[this.stateStackTop] = this.START_STATE
    if (this.action.capacity() == 0) {
        this.action = NewIntTupleWithEstimate(1 << 20)
    } else {
        this.action.reset()
    }
    //
    // Indicate that we are going to run the incremental parser and that
    // it's forbidden to use the utility functions to query the parser.
    //
    this.taking_actions = false
    if marker_kind != 0 {
        var sym []int = []int{this.markerKind}
        this.parse(sym, 0)
    }
}
//
// Find a state in the state stack that has a valid action on ERROR token
//
func (this *DeterministicParser) recoverableState(state int) bool {
    var k int = this.prs.asi(state)
    for  ;this.prs.asr(k) != 0; k++ {
        if this.prs.asr(k) == this.ERROR_SYMBOL {
            return true
        }
    }
    return false
}

//
// Reset the parser at a point where it can legally process
// the error token. If we can't do that, reset it to the beginning.
//
func (this *DeterministicParser)  errorReset()  {
    var gate int
        if this.markerKind == 0 {
            gate = 0
        } else {
            gate = 1
        }
    for ;this.stateStackTop >= gate; this.stateStackTop-- {
        if this.recoverableState(this.stateStack[this.stateStackTop]) {
            break
        }
    }
    if this.stateStackTop < gate {
        this.resetParserEntry(this.markerKind)
    }
    return
}

//
// This is an incremental LALR(k) parser that takes as argument
// the next k tokens in the input. If these k tokens are valid for
// the current configuration, it advances past the first of the k
// tokens and returns either
//
//    . the last transition induced by that token 
//    . the Accept action
//
// If the tokens are not valid, the initial configuration remains
// unchanged and the Error action is returned.
//
// Note that it is the user's responsibility to start the parser in a
// proper configuration by initially invoking the method resetParser
// prior to invoking this function.
//
func (this *DeterministicParser)  parse(sym []int, index int) int {

    // assert(sym.length == prs.getMaxLa())

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
    var save_action_length int = this.action.size()
    var pos int = this.stateStackTop
    var  location_top int = this.stateStackTop - 1

    //
    // When a reduce action is encountered, we compute all REDUCE
    // and associated goto actions induced by the current token.
    // Eventually, a SHIFT, SHIFT-REDUCE, ACCEPT or ERROR action is
    // computed...
    //
    for this.currentAction = this.tAction(this.stateStack[this.stateStackTop], sym, index);
        this.currentAction <= this.NUM_RULES;
        this.currentAction = this.tAction(this.currentAction, sym, index){
        this.action.add(this.currentAction)
        for;; {
            location_top -= (this.prs.rhs(this.currentAction) - 1)

            var state int
            if location_top > pos{
                state=this.locationStack[location_top]
            }else{
                state=   this.stateStack[location_top]
            }

            this.currentAction = this.prs.ntAction(state, this.prs.lhs(this.currentAction))
            if this.currentAction <= this.NUM_RULES {
                continue
            }else{
                break
            }
        }

        //
        // ... Update the maximum useful position of the
        // stateSTACK, push goto state into stack, and
        // continue by compute next action on current symbol
        // and reentering the loop...
        //

        if !(pos < location_top){
            pos = location_top
        }
        if location_top + 1 >= len(this.locationStack) {
            this.reallocateStacks()
        }
        this.locationStack[location_top + 1] = this.currentAction

    }
    //
    // At this point, we have a shift, shift-reduce, accept or error
    // action. stateSTACK contains the configuration of the state stack
    // prior to executing any action on the currenttoken. locationStack
    // contains the configuration of the state stack after executing all
    // reduce actions induced by the current token. The variable pos
    // indicates the highest position in the stateSTACK that is still
    // useful after the reductions are executed.
    //
    if this.currentAction > this.ERROR_ACTION || // SHIFT-REDUCE action ?
            this.currentAction < this.ACCEPT_ACTION { // SHIFT action ?

        this.action.add(this.currentAction)
        //
        // If no error was detected, update the state stack with 
        // the info that was temporarily computed in the locationStack.
        //
        this.stateStackTop = location_top + 1
        var i int = pos + 1
        for  ;i <= this.stateStackTop; i++ {
            this.stateStack[i] = this.locationStack[i]
        }

        //
        // If we have a shift-reduce, process it as well as
        // the goto-reduce actions that follow it.
        //

        if this.currentAction > this.ERROR_ACTION {
            this.currentAction -= this.ERROR_ACTION
            for;; {
                this.stateStackTop -= this.prs.rhs(this.currentAction) - 1
                this.currentAction = this.prs.ntAction(this.stateStack[this.stateStackTop],
                                                        this.prs.lhs(this.currentAction))
                if this.currentAction <= this.NUM_RULES{
                    continue
                }else{
                    break
                }
            }
        }
        //
        // Process the final transition - either a shift action of
        // if we started out with a shift-reduce, the final GOTO
        // action that follows it.
        //
        this.stateStackTop+=1
        if this.stateStackTop >= len(this.stateStack) {
            this.reallocateStacks()
        }
        this.stateStack[this.stateStackTop] = this.currentAction
    } else{
        if this.currentAction == this.ERROR_ACTION {
        this.action.reset(save_action_length)// restore original action state.
        }
    }
    return this.currentAction
}
//
// Now do the final parse of the input based on the actions in
// the list "action" and the sequence of tokens in the token stream.
//
func (this *DeterministicParser)  parseActions() interface{} {
    //
    // Indicate that we are processing actions now (for the incremental
    // parser) and that it's ok to use the utility functions to query the
    // parser.
    //
    this.taking_actions = true
    this.tokStream.reset()
    this.lastToken = this.tokStream.getPrevious(this.tokStream.peek())
    var curtok int
    if this.markerKind == 0 {
        curtok = this.tokStream.getToken()
    }else{
        curtok=this.lastToken
    }


    //
    // Reparse the input...
    //
    this.stateStackTop = -1
    this.currentAction = this.START_STATE
    var i int = 0
    for ; i < this.action.size(); i++{
        //
        // if the parser needs to stop processing, it may do so here.
        //
        if this.monitor != nil && this.monitor.isCancelled() {
            this.taking_actions = false
            return nil
        }

        this.stateStackTop+=1
        this.stateStack[this.stateStackTop] = this.currentAction
        this.locationStack[this.stateStackTop] = curtok

        this.currentAction = this.action.get(i)
        if this.currentAction <= this.NUM_RULES { // a reduce action?

            this.stateStackTop--// turn reduction intoshift-reduction
            this.processReductions()
        } else{ // a shift or shift-reduce action

            this.lastToken = curtok
            curtok = this.tokStream.getToken()
            if this.currentAction > this.ERROR_ACTION {
                this.currentAction -= this.ERROR_ACTION
                this.processReductions()
            }
        }
    }


    this.taking_actions = false// indicate that we are done.
    this.action = nil
    if this.markerKind == 0 {
        return this.parseStack[0]
    }else{
        return this.parseStack[1]
    }
}
        
        
        