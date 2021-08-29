package lpg2


type RepairCandidate  struct{
    symbol int 
    location int 
}
func NewRepairCandidate() *RepairCandidate {
    t := new(RepairCandidate)
    t.symbol = 0
    t.location = 0
    return t
}  

type PrimaryRepairInfo struct{
    distance int 
    misspellIndex int 
    code int 
    bufferPosition int 
    symbol int
}
func NewPrimaryRepairInfo() *PrimaryRepairInfo {
      t := new(PrimaryRepairInfo)
      return t
}
func (self *PrimaryRepairInfo)  copy(clone PrimaryRepairInfo) {
        self.distance = clone.distance
        self.misspellIndex = clone.misspellIndex
        self.code = clone.code
        self.bufferPosition = clone.bufferPosition
        self.symbol = clone.symbol
        return
}

type SecondaryRepairInfo struct{
    code int 
    distance int 
    bufferPosition int 
    stackPosition int 
    numDeletions int 
    symbol int 
    recoveryOnNextStack  bool
}
func NewSecondaryRepairInfo() *SecondaryRepairInfo {
    t := new(SecondaryRepairInfo)
    return t
}

type StateInfo struct{
    state int 
    next int 
}
func NewStateInfo(state int , next int ) *StateInfo{
    t := new(StateInfo)
    t.state = state
    t.next = next
    return t
}

const  STACK_INCREMENT int  = 256
const  BUFF_UBOUND int  = 31
const  BUFF_SIZE int  = 32
const  MAX_DISTANCE int  = 30
const  MIN_DISTANCE int  = 3

type DiagnoseParser  struct {

    monitor Monitor
    tokStream TokenStream

    prs ParseTable

    ERROR_SYMBOL int 
    SCOPE_SIZE int 
    MAX_NAME_LENGTH int 
    NT_OFFSET int 
    LA_STATE_OFFSET int 
    NUM_RULES int 
    NUM_SYMBOLS int 
    START_STATE int 
    EOFT_SYMBOL int 
    EOLT_SYMBOL int 
    ACCEPT_ACTION int 
    ERROR_ACTION int 

    list []int



    maxErrors int 

    maxTime int 

    stateStackTop int 
    stateStack []int 

    locationStack []int 

    tempStackTop int 
    tempStack []int 

    prevStackTop int 
    prevStack []int 

    nextStackTop int 
    nextStack []int 

    scopeStackTop int 
    scopeIndex []int 
    scopePosition []int 

    buffer []int 

    stateSeen []int 

    statePoolTop int 
    statePool []int
    main_configuration_stack *ConfigurationStack
}
func NewDiagnoseParser(tokStream TokenStream, prs ParseTable, maxErrors int , maxTime int,monitor Monitor ) *DiagnoseParser {
    self := new(DiagnoseParser)
    self.monitor = monitor
    self.maxErrors = maxErrors
    self.maxTime = maxTime
    self.tokStream = tokStream
    self.prs = prs
    self.main_configuration_stack = NewConfigurationStack(prs)
    self.ERROR_SYMBOL = prs.getErrorSymbol()
    self.SCOPE_SIZE = prs.getScopeSize()
    self.MAX_NAME_LENGTH = prs.getMaxNameLength()
    self.NT_OFFSET = prs.getNtOffset()
    self.LA_STATE_OFFSET = prs.getLaStateOffset()
    self.NUM_RULES = prs.getNumRules()
    self.NUM_SYMBOLS = prs.getNumSymbols()
    self.START_STATE = prs.getStartState()
    self.EOFT_SYMBOL = prs.getEoftSymbol()
    self.EOLT_SYMBOL = prs.getEoltSymbol()
    self.ACCEPT_ACTION = prs.getAcceptAction()
    self.ERROR_ACTION = prs.getErrorAction()
    self.list = make([]int,self.NUM_SYMBOLS + 1,self.NUM_SYMBOLS + 1)
    return self
}


func (self *DiagnoseParser) setMonitor(monitor Monitor) {
    self.monitor = monitor
}
   




    func (self *DiagnoseParser) rhs(index int ) int  {
        return self.prs.rhs(index)
    }
    func (self *DiagnoseParser) baseAction(index int ) int  {
        return self.prs.baseAction(index)
    }
    func (self *DiagnoseParser) baseCheck(index int ) int  {
        return self.prs.baseCheck(index)
    }
    func (self *DiagnoseParser) lhs(index int ) int  {
        return self.prs.lhs(index)
    }
    func (self *DiagnoseParser) termCheck(index int ) int  {
        return self.prs.termCheck(index)
    }
    func (self *DiagnoseParser) termAction(index int ) int  {
        return self.prs.termAction(index)
    }
    func (self *DiagnoseParser) asb(index int ) int  {
        return self.prs.asb(index)
    }
    func (self *DiagnoseParser) asr(index int ) int  {
        return self.prs.asr(index)
    }
    func (self *DiagnoseParser) nasb(index int ) int  {
        return self.prs.nasb(index)
    }
    func (self *DiagnoseParser) nasr(index int ) int  {
        return self.prs.nasr(index)
    }
    func (self *DiagnoseParser) terminalIndex(index int ) int  {
        return self.prs.terminalIndex(index)
    }
    func (self *DiagnoseParser) nonterminalIndex(index int ) int  {
        return self.prs.nonterminalIndex(index)
    }
    func (self *DiagnoseParser) symbolIndex(index int ) int  {
        return index > self.NT_OFFSET
            ? self.nonterminalIndex(index - self.NT_OFFSET)
            : self.terminalIndex(index)
    }
    func (self *DiagnoseParser) scopePrefix(index int ) int  {
        return self.prs.scopePrefix(index)
    }
    func (self *DiagnoseParser) scopeSuffix(index int ) int  {
        return self.prs.scopeSuffix(index)
    }
    func (self *DiagnoseParser) scopeLhs(index int ) int  {
        return self.prs.scopeLhs(index)
    }
    func (self *DiagnoseParser) scopeLa(index int ) int  {
        return self.prs.scopeLa(index)
    }
    func (self *DiagnoseParser) scopeStateSet(index int ) int  {
        return self.prs.scopeStateSet(index)
    }
    func (self *DiagnoseParser) scopeRhs(index int ) int  {
        return self.prs.scopeRhs(index)
    }
    func (self *DiagnoseParser) scopeState(index int ) int  {
        return self.prs.scopeState(index)
    }
    func (self *DiagnoseParser) inSymb(index int ) int  {
        return self.prs.inSymb(index)
    }
    func (self *DiagnoseParser) name(index int ): string {
        return self.prs.name(index)
    }
    func (self *DiagnoseParser) originalState(state int ) int  {
        return self.prs.originalState(state)
    }
    func (self *DiagnoseParser) asi(state int ) int  {
        return self.prs.asi(state)
    }
    func (self *DiagnoseParser) nasi(state int ) int  {
        return self.prs.nasi(state)
    }
    func (self *DiagnoseParser) inSymbol(state int ) int  {
        return self.prs.inSymbol(state)
    }
    func (self *DiagnoseParser) ntAction(state int , sym int ) int  {
        return self.prs.ntAction(state, sym)
    }
    func (self *DiagnoseParser) isNullable(symbol int )  bool {
        return self.prs.isNullable(symbol)
    }

  
    func (self *DiagnoseParser) reallocateStacks() {
        var old_stack_length int  =  len(self.stateStack)
        var stack_length int  = old_stack_length + STACK_INCREMENT

        if len(self.stateStack) == 0 {
            self.stateStack = make([]int,stack_length,stack_length)
            self.locationStack =  make([]int,stack_length,stack_length)
            self.tempStack =  make([]int,stack_length,stack_length)
            self.prevStack = make([]int,stack_length,stack_length)
            self.nextStack = make([]int,stack_length,stack_length)
            self.scopeIndex = make([]int,stack_length,stack_length)
            self.scopePosition = make([]int,stack_length,stack_length)
        } else {
            self.stateStack = arraycopy(self.stateStack, 0, make([]int,stack_length,stack_length), 0, old_stack_length)
            self.locationStack = arraycopy(self.locationStack, 0, make([]int,stack_length,stack_length), 0, old_stack_length)
            self.tempStack = arraycopy(self.tempStack, 0, make([]int,stack_length,stack_length), 0, old_stack_length)
            self.prevStack = arraycopy(self.prevStack, 0, make([]int,stack_length,stack_length), 0, old_stack_length)
            self.nextStack = arraycopy(self.nextStack, 0, make([]int,stack_length,stack_length), 0, old_stack_length)
            self.scopeIndex = arraycopy(self.scopeIndex, 0, make([]int,stack_length,stack_length), 0, old_stack_length)
            self.scopePosition = arraycopy(self.scopePosition, 0, make([]int,stack_length,stack_length), 0, old_stack_length)
        }
        return
    }
   
    func (self *DiagnoseParser) diagnose(error_token int ) {
        self.diagnoseEntry2(0, error_token)
    }

    func (self *DiagnoseParser) diagnoseEntry1(marker_kind int ) {
        self.reallocateStacks()
        self.tempStackTop = 0
        self.tempStack[self.tempStackTop] = self.START_STATE
        self.tokStream.resetDefault()
        var current_token int 
        var current_kind int 
        if marker_kind == 0 {
            current_token = self.tokStream.getToken()
            current_kind = self.tokStream.getKind(current_token)
        } else {
            current_token = self.tokStream.peek()
            current_kind = marker_kind
        }

        //
        // If an error was found, start the diagnosis and recovery.
        //
        var error_token int  = self.parseForError(current_kind)
        if error_token != 0 {
            self.diagnoseEntry2(marker_kind, error_token)
        }
        return
    }
    func (self *DiagnoseParser) diagnoseEntry2(marker_kind int , error_token int ) {
        var action  = NewIntTupleWithEstimate(1 << 18)
        var startTime int  = Date.now()
        var errorCount int  = 0

        //
        // Compute sequence of actions that leads us to the
        // error_token.
        //
        if len(self.stateStack)== 0 {
            self.reallocateStacks()
        }

        self.tempStackTop = 0
        self.tempStack[self.tempStackTop] = self.START_STATE
        self.tokStream.resetDefault()
        var current_token int
        var current_kind int
        if (marker_kind == 0) {
            current_token = self.tokStream.getToken()
            current_kind = self.tokStream.getKind(current_token)
        } else {
            current_token = self.tokStream.peek()
            current_kind = marker_kind
        }
        self.parseUpToError(action, current_kind, error_token)

        //
        // Start parsing
        //
        self.stateStackTop = 0
        self.stateStack[self.stateStackTop] = self.START_STATE

        self.tempStackTop = self.stateStackTop
        arraycopy(self.tempStack, 0, self.stateStack, 0, self.tempStackTop + 1)

        self.tokStream.resetDefault()
        if (marker_kind == 0) {
            current_token = self.tokStream.getToken()
            current_kind = self.tokStream.getKind(current_token)
        } else {
            current_token = self.tokStream.peek()
            current_kind = marker_kind
        }
        self.locationStack[self.stateStackTop] = current_token

        //
        // Process a terminal
        //
        var act int 
        for ;; {
            //
            // Synchronize state stacks and update the location stack
            //
            var prev_pos int  = -1
            self.prevStackTop = -1

            var next_pos int  = -1
            self.nextStackTop = -1

            var pos int  = self.stateStackTop
            self.tempStackTop = self.stateStackTop - 1
            arraycopy(self.stateStack, 0, self.tempStack, 0, self.stateStackTop + 1)

            var action_index int  = 0
            act = action.get(action_index++)// tAction(act, current_kind)

            //
            // When a reduce action is encountered, we compute all REDUCE
            // and associated goto actions induced by the current token.
            // Eventually, a SHIFT, SHIFT-REDUCE, ACCEPT or ERROR action is
            // computed...
            //
            for ;act <= self.NUM_RULES; {
                for;; {
                    self.tempStackTop -= (self.rhs(act) - 1)
                    act = self.ntAction(self.tempStack[self.tempStackTop], self.lhs(act))
                    if  act <= self.NUM_RULES {
                        continue
                    }else {
                        break
                    }
                }

                //
                // ... Update the maximum useful position of the
                // (STATE_)STACK, push goto state into stack, and
                // compute next action on current symbol ...
                //
                if self.tempStackTop + 1 >= len(self.stateStack) {
                    self.reallocateStacks()
                }
                pos = pos < self.tempStackTop ? pos : self.tempStackTop
                self.tempStack[self.tempStackTop + 1] = act
                act = action.get(action_index++) // tAction(act, current_kind)
            }
            //
            // At self point, we have a shift, shift-reduce, accept or error
            // action.  STACK contains the configuration of the state stack
            // prior to executing any action on current_token. next_stack contains
            // the configuration of the state stack after executing all
            // reduce actions induced by current_token.  The variable pos indicates
            // the highest position in STACK that is still useful after the
            // reductions are executed.
            //
            for;act > self.ERROR_ACTION || act < self.ACCEPT_ACTION; {

                //
                // if the parser needs to stop processing,
                // it may do so here.
                //
                if self.monitor != nil  && self.monitor.isCancelled() {
                    return
                }

                self.nextStackTop = self.tempStackTop + 1
                var i int
                for i   = next_pos + 1 ;i <= self.nextStackTop; i++ {
                    self.nextStack[i] = self.tempStack[i]
                }
                var k int
                for  k   = pos + 1 ;k <= self.nextStackTop ;k++ {
                    self.locationStack[k] = self.locationStack[self.stateStackTop]
                }

                //
                // If we have a shift-reduce, process it as well as
                // the goto-reduce actions that follow it.
                //
                if (act > self.ERROR_ACTION) {
                    act -= self.ERROR_ACTION
                    for;; {
                        self.nextStackTop -= (self.rhs(act) - 1)
                        act = self.ntAction(self.nextStack[self.nextStackTop], self.lhs(act))
                        if act <= self.NUM_RULES {
                            continue
                        }else {
                            break
                        }
                    }
                    pos = pos < self.nextStackTop ? pos : self.nextStackTop
                }

                if self.nextStackTop + 1 >=  len(self.stateStack) {
                    self.reallocateStacks()
                }

                self.tempStackTop = self.nextStackTop
                self.nextStack[++self.nextStackTop] = act
                next_pos = self.nextStackTop
                //
                // Simulate the parser through the next token without
                // destroying STACK or next_stack.
                //
                current_token = self.tokStream.getToken()
                current_kind = self.tokStream.getKind(current_token)
                act = action.get(action_index++)// tAction(act, current_kind)
                while (act <= self.NUM_RULES) {
                    //
                    // ... Process all goto-reduce actions following
                    // reduction, until a goto action is computed ...
                    //

                    for ;; {
                        var lhs_symbol int  = self.lhs(act)
                        self.tempStackTop -= (self.rhs(act) - 1)
                        act = (self.tempStackTop > next_pos
                                                 ? self.tempStack[self.tempStackTop]
                                                 : self.nextStack[self.tempStackTop])
                        act = self.ntAction(act, lhs_symbol)
                    } while (act <= self.NUM_RULES)
                    //
                    // ... Update the maximum useful position of the
                    // (STATE_)STACK, push GOTO state into stack, and
                    // compute next action on current symbol ...
                    //
                    if (self.tempStackTop + 1 >= self.stateStack.length) {
                        self.reallocateStacks()
                    }
                    next_pos = next_pos < self.tempStackTop ? next_pos : self.tempStackTop
                    self.tempStack[self.tempStackTop + 1] = act
                    act = action.get(action_index++)// tAction(act, current_kind)
                }
                //
                // No error was detected, Read next token into
                // PREVTOK element, advance CURRENT_TOKEN pointer and
                // update stacks.
                //
                if (act != self.ERROR_ACTION) {
                    self.prevStackTop = self.stateStackTop
                    for (var i int  = prev_pos + 1 i <= self.prevStackTop i++) {
                        self.prevStack[i] = self.stateStack[i]
                    }
                    prev_pos = pos

                    self.stateStackTop = self.nextStackTop
                    for (var k int  = pos + 1 k <= self.stateStackTop k++) {
                        self.stateStack[k] = self.nextStack[k]
                    }
                    self.locationStack[self.stateStackTop] = current_token
                    pos = next_pos
                }
            }

            //
            // At self stage, either we have an ACCEPT or an ERROR
            // action.
            //
            if (act == self.ERROR_ACTION) {
                //
                // An error was detected.
                //
                errorCount += 1
                //
                // Check time and error limits after the first error encountered
                // Exit if number of errors exceeds maxError or if maxTime reached
                //
                if (errorCount > 1) {
                    if (self.maxErrors > 0 && errorCount > self.maxErrors) {
                        break
                    }
                    if (self.maxTime > 0 && Date.now() - startTime > self.maxTime) {
                        break
                    }
                }
                var candidate: RepairCandidate = self.errorRecovery(current_token)
                //
                // if the parser needs to stop processing,
                // it may do so here.
                //
                if (self.monitor  && self.monitor.isCancelled()) {
                    return
                }
                act = self.stateStack[self.stateStackTop]

                //
                // If the recovery was successful on a nonterminal candidate,
                // parse through that candidate and "read" the next token.
                //
                if (candidate.symbol == 0) {
                    break
                } else {
                    if (candidate.symbol > self.NT_OFFSET) {
                        var lhs_symbol int  = candidate.symbol - self.NT_OFFSET
                        act = self.ntAction(act, lhs_symbol)
                        while (act <= self.NUM_RULES) {
                            self.stateStackTop -= (self.rhs(act) - 1)
                            act = self.ntAction(self.stateStack[self.stateStackTop], self.lhs(act))
                        }
                        self.stateStack[++self.stateStackTop] = act
                        current_token = self.tokStream.getToken()
                        current_kind = self.tokStream.getKind(current_token)
                        self.locationStack[self.stateStackTop] = current_token
                    } else {
                        current_kind = candidate.symbol
                        self.locationStack[self.stateStackTop] = candidate.location
                    }
                }
                //
                // At self stage, we have a recovery configuration. See how
                // far we can go with it.
                //
                var next_token int  = self.tokStream.peek()
                self.tempStackTop = self.stateStackTop
                arraycopy(self.stateStack, 0, self.tempStack, 0, self.stateStackTop + 1)
                error_token = self.parseForError(current_kind)

                if (error_token != 0) {
                    self.tokStream.reset(next_token)
                    self.tempStackTop = self.stateStackTop
                    arraycopy(self.stateStack, 0, self.tempStack, 0, self.stateStackTop + 1)
                    self.parseUpToError(action, current_kind, error_token)
                    self.tokStream.reset(next_token)
                } else {
                    act = self.ACCEPT_ACTION
                }
            }
        } while (act != self.ACCEPT_ACTION)
        return
    }
    //
    // Given the configuration consisting of the states in tempStack
    // and the sequence of tokens (current_kind, followed by the tokens
    // in tokStream), keep parsing until either the parse completes
    // successfully or it encounters an error. If the parse is not
    // succesful, we return the farthest token on which an error was
    // encountered. Otherwise, we return 0.
    //
    func (self *DiagnoseParser) parseForError(current_kind int ) int  {
        var error_token int  = 0
        //
        // Get next token in stream and compute initial action
        //
        var curtok int  = self.tokStream.getPrevious(self.tokStream.peek())
        var act int  = self.tAction(self.tempStack[self.tempStackTop], current_kind)
        //
        // Allocate configuration stack.
        //
        var configuration_stack = NewConfigurationStack(self.prs)

        //
        // Keep parsing until we reach the end of file and succeed or
        // an error is encountered. The list of actions executed will
        // be store in the "action" tuple.
        //
        for ;; {
            if (act <= self.NUM_RULES) {

                self.tempStackTop--

                for;; {
                    self.tempStackTop -= self.rhs(act) - 1
                    act = self.ntAction(self.tempStack[self.tempStackTop], self.lhs(act))
                    if act <= self.NUM_RULES{
                        continue
                    }else {
                        break
                    }
                }

            }
            else if (act > self.ERROR_ACTION)
            {
                curtok = self.tokStream.getToken()
                current_kind = self.tokStream.getKind(curtok)
                act -= self.ERROR_ACTION

                for ;; {
                    self.tempStackTop -= (self.rhs(act) - 1)
                    act = self.ntAction(self.tempStack[self.tempStackTop], self.lhs(act))
                } while (act <= self.NUM_RULES)
            }
            else if (act < self.ACCEPT_ACTION) {
                curtok = self.tokStream.getToken()
                current_kind = self.tokStream.getKind(curtok)
            }
            else if (act == self.ERROR_ACTION)
            {
                error_token = (error_token > curtok ? error_token : curtok)

                var configuration = configuration_stack.pop()
                if (configuration == nil) {
                    act = self.ERROR_ACTION
                } else {
                    self.tempStackTop = configuration.stack_top
                    configuration.retrieveStack(self.tempStack)
                    act = configuration.act
                    curtok = configuration.curtok
                    // no need to execute: action.reset(configuration.action_length)
                    current_kind = self.tokStream.getKind(curtok)
                    self.tokStream.reset(self.tokStream.getNext(curtok))
                    continue
                }
                break
            }
            else if (act > self.ACCEPT_ACTION)
            {
                if (configuration_stack.findConfiguration(self.tempStack, self.tempStackTop, curtok)) {
                    act = self.ERROR_ACTION
                } else {
                    configuration_stack.push(self.tempStack, self.tempStackTop, act + 1, curtok, 0)
                    act = self.baseAction(act)
                }
                continue
            }
            else {
                break// assert(act == ACCEPT_ACTION)
            }

            if (++self.tempStackTop >= len(self.tempStack)) {
                self.reallocateStacks()
            }
            self.tempStack[self.tempStackTop] = act
            act = self.tAction(act, current_kind)
        }
        return (act == self.ERROR_ACTION ? error_token : 0)
    }
    //
    // Given the configuration consisting of the states in tempStack
    // and the sequence of tokens (current_kind, followed by the tokens
    // in tokStream), parse up to error_token in the tokStream and store
    // all the parsing actions executed in the "action" tuple.
    //
    func (self *DiagnoseParser) parseUpToError(action *IntTuple, current_kind int , error_token int ) {
        //
        // Assume predecessor of next token and compute initial action
        //
        var curtok int  = self.tokStream.getPrevious(self.tokStream.peek())
        var act int  = self.tAction(self.tempStack[self.tempStackTop], current_kind)
        //
        // Allocate configuration stack.
        //
        var configuration_stack = NewConfigurationStack(self.prs)
        //
        // Keep parsing until we reach the end of file and succeed or
        // an error is encountered. The list of actions executed will
        // be store in the "action" tuple.
        //
        action.resetFromZero()
        for ;;  {
            if act <= self.NUM_RULES {
                action.add(act)// save self reduce action
                self.tempStackTop--

                for;; {
                    self.tempStackTop -= (self.rhs(act) - 1)
                    act = self.ntAction(self.tempStack[self.tempStackTop], self.lhs(act))
                    if act <= self.NUM_RULES{
                        continue
                    }else {
                        break
                    }
                }

            }
            else if (act > self.ERROR_ACTION)
            {
                action.add(act) // save self shift-reduce action
                curtok = self.tokStream.getToken()
                current_kind = self.tokStream.getKind(curtok)
                act -= self.ERROR_ACTION

                for ;; {
                    self.tempStackTop -= (self.rhs(act) - 1)
                    act = self.ntAction(self.tempStack[self.tempStackTop], self.lhs(act))
                } while (act <= self.NUM_RULES)
            }
            else if (act < self.ACCEPT_ACTION)
            {
                action.add(act)// save self shift action
                curtok = self.tokStream.getToken()
                current_kind = self.tokStream.getKind(curtok)
            }
            else if (act == self.ERROR_ACTION)
            {
                if (curtok != error_token) {
                    var configuration = configuration_stack.pop()
                    if (configuration == nil) {
                        act = self.ERROR_ACTION
                    } else {
                        self.tempStackTop = configuration.stack_top
                        configuration.retrieveStack(self.tempStack)
                        act = configuration.act
                        curtok = configuration.curtok
                        action.reset(configuration.action_length)
                        current_kind = self.tokStream.getKind(curtok)
                        self.tokStream.reset(self.tokStream.getNext(curtok))
                        continue
                    }
                }
                break
            }
            else if (act > self.ACCEPT_ACTION) {
                if (configuration_stack.findConfiguration(self.tempStack, self.tempStackTop, curtok)) {
                    act = self.ERROR_ACTION
                } else {
                    configuration_stack.push(self.tempStack, self.tempStackTop, act + 1, curtok, action.size())
                    act = self.baseAction(act)
                }
                continue
            } else {
                break// assert(act == ACCEPT_ACTION)
            }


            if (++self.tempStackTop >= len(self.tempStack)) {
                self.reallocateStacks()
            }
            self.tempStack[self.tempStackTop] = act
            act = self.tAction(act, current_kind)
        }
        action.add(self.ERROR_ACTION)
        return
    }
    //
    // Try to parse until first_symbol and all tokens in BUFFER have
    // been consumed, or an error is encountered. Return the number
    // of tokens that were expended before the parse blocked.
    //
    func (self *DiagnoseParser) parseCheck(stack []int, stack_top int , first_symbol int , buffer_position int ) int  {
        var buffer_index int
        var current_kind int

        var local_stack []int = new []int(stack.length)
        var local_stack_top int  = stack_top
        for (var i int  = 0; i <= stack_top ;i++) {
            local_stack[i] = stack[i]
        }
        var configuration_stack: ConfigurationStack = new ConfigurationStack(self.prs)

        //
        // If the first symbol is a nonterminal, process it here.
        //
        var act int  = local_stack[local_stack_top]
        if (first_symbol > self.NT_OFFSET) {
            var lhs_symbol int  = first_symbol - self.NT_OFFSET
            buffer_index = buffer_position
            current_kind = self.tokStream.getKind(self.buffer[buffer_index])
            self.tokStream.reset(self.tokStream.getNext(self.buffer[buffer_index]))
            act = self.ntAction(act, lhs_symbol)
            for ;act <= self.NUM_RULES; {
                local_stack_top -= (self.rhs(act) - 1)
                act = self.ntAction(local_stack[local_stack_top], self.lhs(act))
            }
        } else {
            local_stack_top--
            buffer_index = buffer_position - 1
            current_kind = first_symbol
            self.tokStream.reset(self.buffer[buffer_position])
        }

        //
        // Start parsing the remaining symbols in the buffer
        //
        if (++local_stack_top >= local_stack.length) {// Stack overflow!!!
            return buffer_index
        }
        local_stack[local_stack_top] = act

        act = self.tAction(act, current_kind)

        for ;;
        {
            if (act <= self.NUM_RULES) // reduce action
            {  
                local_stack_top -= self.rhs(act)
                act = self.ntAction(local_stack[local_stack_top], self.lhs(act))
                for;act <= self.NUM_RULES; {
                    local_stack_top -= (self.rhs(act) - 1)
                    act = self.ntAction(local_stack[local_stack_top], self.lhs(act))
                }
            }
            else if (act > self.ERROR_ACTION)     // shift-reduce action
            {
                if (buffer_index++ == MAX_DISTANCE) {
                    break
                }
                current_kind = self.tokStream.getKind(self.buffer[buffer_index])
                self.tokStream.reset(self.tokStream.getNext(self.buffer[buffer_index]))
                act -= self.ERROR_ACTION

                for ;; {
                    local_stack_top -= (self.rhs(act) - 1)
                    act = self.ntAction(local_stack[local_stack_top], self.lhs(act))
                } while (act <= self.NUM_RULES)
            }
            else if (act < self.ACCEPT_ACTION)    // shift action
            {
                if (buffer_index++ == MAX_DISTANCE) {
                    break
                }
                current_kind = self.tokStream.getKind(self.buffer[buffer_index])
                self.tokStream.reset(self.tokStream.getNext(self.buffer[buffer_index]))
            }
            else if (act == self.ERROR_ACTION)
            {
                var configuration = configuration_stack.pop()
                if (configuration == nil) {
                    act = self.ERROR_ACTION
                } else {
                    local_stack_top = configuration.stack_top
                    configuration.retrieveStack(local_stack)
                    act = configuration.act
                    buffer_index = configuration.curtok
                    // no need to execute: action.reset(configuration.action_length)
                    current_kind = self.tokStream.getKind(self.buffer[buffer_index])
                    self.tokStream.reset(self.tokStream.getNext(self.buffer[buffer_index]))
                    continue
                }
                break
            }
            else if (act > self.ACCEPT_ACTION)
            {
                if (configuration_stack.findConfiguration(local_stack, local_stack_top, buffer_index)) {
                    act = self.ERROR_ACTION
                } else {
                    configuration_stack.push(local_stack, local_stack_top, act + 1, buffer_index, 0)
                    act = self.baseAction(act)
                }
                continue
            }
            else {
                break
            }
                        

            if (++local_stack_top >= local_stack.length) {
                break
            }
            local_stack[local_stack_top] = act
            act = self.tAction(act, current_kind)
        }
        return (act == self.ACCEPT_ACTION ? MAX_DISTANCE : buffer_index)
    }

    //
    //  This routine is invoked when an error is encountered.  It
    // tries to diagnose the error and recover from it.  If it is
    // successful, the state stack, the current token and the buffer
    // are readjusted i.e., after a successful recovery,
    // state_stack_top points to the location in the state stack
    // that contains the state on which to recover current_token
    // identifies the symbol on which to recover.
    //
    // Up to three configurations may be available when self routine
    // is invoked. PREV_STACK may contain the sequence of states
    // preceding any action on prevtok, STACK always contains the
    // sequence of states preceding any action on current_token, and
    // NEXT_STACK may contain the sequence of states preceding any
    // action on the successor of current_token.
    //
    func (self *DiagnoseParser) errorRecovery(error_token int ): RepairCandidate {
        var prevtok int  = self.tokStream.getPrevious(error_token)

        //
        // Try primary phase recoveries. If not successful, try secondary
        // phase recoveries.  If not successful and we are at end of the
        // file, we issue the end-of-file error and quit. Otherwise, ...
        //
        var candidate: RepairCandidate = self.primaryPhase(error_token)
        if (candidate.symbol != 0) {
            return candidate
        }
        candidate = self.secondaryPhase(error_token)
        if (candidate.symbol != 0) {
            return candidate
        }
        //
        // At self point, primary and (initial attempt at) secondary
        // recovery did not work.  We will now get into "panic mode" and
        // keep trying secondary phase recoveries until we either find
        // a successful recovery or have consumed the remaining input
        // tokens.
        //
        if (self.tokStream.getKind(error_token) != self.EOFT_SYMBOL) {
            while (self.tokStream.getKind(self.buffer[BUFF_UBOUND]) != self.EOFT_SYMBOL) {
                candidate = self.secondaryPhase(self.buffer[MAX_DISTANCE - MIN_DISTANCE + 2])
                if (candidate.symbol != 0) {
                    return candidate
                }
            }
        }
        //
        // If no successful recovery is found and we have reached the
        // end of the file, check whether or not any scope recovery is
        // applicable at the end of the file after discarding some
        // states.
        //
        var scope_repair = new PrimaryRepairInfo()
        scope_repair.bufferPosition = BUFF_UBOUND
        for (var top int  = self.stateStackTop top >= 0 top--) {
            self.scopeTrial(scope_repair, self.stateStack, top)
            if (scope_repair.distance > 0) {
                break
            }
        }
        //
        // If any scope repair was successful, emit the message now
        //
        for (var i int  = 0 i < self.scopeStackTop i++) {
            self.emitError(SCOPE_CODE,
                            -self.scopeIndex[i],
                            self.locationStack[self.scopePosition[i]],
                            self.buffer[1],
                            self.nonterminalIndex(self.scopeLhs(self.scopeIndex[i])))
        }

        //
        // If the original error_token was already pointing to the EOF, issue the EOF-reached message.
        //
        if (self.tokStream.getKind(error_token) == self.EOFT_SYMBOL) {
            self.emitError(EOF_CODE,
                            self.terminalIndex(self.EOFT_SYMBOL),
                            prevtok,
                            prevtok)
        } else {
            //
            // We reached the end of the file while panicking. Delete all
            // remaining tokens in the input.
            //
            var i int 
            for (i = BUFF_UBOUND ;self.tokStream.getKind(self.buffer[i]) == self.EOFT_SYMBOL; i--)
            {
            }

            self.emitError(DELETION_CODE,
                            self.terminalIndex(self.tokStream.getKind(error_token)),
                            error_token,
                            self.buffer[i])
        }
        //
        // Create the "failed" candidate and return it.
        //
        candidate.symbol = 0
        candidate.location = self.buffer[BUFF_UBOUND]// point to EOF
        return candidate
    }

    //
    // This function tries primary and scope recovery on each
    // available configuration.  If a successful recovery is found
    // and no secondary phase recovery can do better, a diagnosis is
    // issued, the configuration is updated and the function returns
    // "true".  Otherwise, it returns "false".
    //
    func (self *DiagnoseParser) primaryPhase(error_token int ): RepairCandidate {
        //
        // Initialize the buffer.
        //
        var i int  = (self.nextStackTop >= 0 ? 3 : 2)
        self.buffer[i] = error_token

        for (var j int  = i j > 0 j--) {
            self.buffer[j - 1] = self.tokStream.getPrevious(self.buffer[j])
        }
        for (var k int  = i + 1 k < BUFF_SIZE k++) {
            self.buffer[k] = self.tokStream.getNext(self.buffer[k - 1])
        }

        //
        // If NEXT_STACK_TOP > 0 then the parse was successful on CURRENT_TOKEN
        // and the error was detected on the successor of CURRENT_TOKEN. In
        // that case, first check whether or not primary recovery is
        // possible on next_stack ...
        //
        var repair PrimaryRepairInfo = new PrimaryRepairInfo()
        if (self.nextStackTop >= 0) {
            repair.bufferPosition = 3
            self.checkPrimaryDistance(repair, self.nextStack, self.nextStackTop)
        }

        //
        // ... Try primary recovery on the current token and compare
        // the quality of self recovery to the one on the next token...
        //
        var base_repair PrimaryRepairInfo = new PrimaryRepairInfo(repair)
        base_repair.bufferPosition = 2
        self.checkPrimaryDistance(base_repair, self.stateStack, self.stateStackTop)
        if (base_repair.distance > repair.distance || base_repair.misspellIndex > repair.misspellIndex) {
            repair = base_repair
        }

        //
        // Finally, if prev_stack_top >= 0 try primary recovery on
        // the prev_stack configuration and compare it to the best
        // recovery computed thus far.
        //
        if (self.prevStackTop >= 0) {
            var prev_repair PrimaryRepairInfo = new PrimaryRepairInfo(repair)
            prev_repair.bufferPosition = 1
            self.checkPrimaryDistance(prev_repair, self.prevStack, self.prevStackTop)
            if (prev_repair.distance > repair.distance || prev_repair.misspellIndex > repair.misspellIndex) {
                repair = prev_repair
            }
        }


        //
        // Before accepting the best primary phase recovery obtained,
        // ensure that we cannot do better with a similar secondary
        // phase recovery.
        //
        var candidate: RepairCandidate = new RepairCandidate()
        if (self.nextStackTop >= 0)// next_stack available
        {
            if (self.secondaryCheck(self.nextStack, self.nextStackTop, 3, repair.distance)) {
                return candidate
            }
        } else {
            if (self.secondaryCheck(self.stateStack, self.stateStackTop, 2, repair.distance)) {
                return candidate
            }
        }

        //
        // First, adjust distance if the recovery is on the error token
        // it is important that the adjustment be made here and not at
        // each primary trial to prevent the distance tests from being
        // biased in favor of deferred recoveries which have access to
        // more input tokens...
        //
        repair.distance = repair.distance - repair.bufferPosition + 1

        //
        // ...Next, adjust the distance if the recovery is a deletion or
        // (some form of) substitution...
        //
        if (repair.code == INVALID_CODE ||
            repair.code == DELETION_CODE ||
            repair.code == SUBSTITUTION_CODE ||
            repair.code == MERGE_CODE) {
            repair.distance--
        }

        //
        // ... After adjustment, check if the most successful primary
        // recovery can be applied.  If not, continue with more radical
        // recoveries...
        //
        if (repair.distance < MIN_DISTANCE) {
            return candidate
        }

        //
        // When processing an insertion error, if the token preceeding
        // the error token is not available, we change the repair code
        // into a BEFORE_CODE to instruct the reporting routine that it
        // indicates that the repair symbol should be inserted before
        // the error token.
        //
        if (repair.code == INSERTION_CODE) {
            if (self.tokStream.getKind(self.buffer[repair.bufferPosition - 1]) == 0) {
                repair.code = BEFORE_CODE
            }
        }


        //
        // Select the proper sequence of states on which to recover,
        // update stack accordingly and call diagnostic routine.
        //
        if (repair.bufferPosition == 1) {
            self.stateStackTop = self.prevStackTop
            arraycopy(self.prevStack, 0, self.stateStack, 0, self.stateStackTop + 1)
        } else {
            if (self.nextStackTop >= 0 && repair.bufferPosition >= 3) {
                self.stateStackTop = self.nextStackTop
                arraycopy(self.nextStack, 0, self.stateStack, 0, self.stateStackTop + 1)
                self.locationStack[self.stateStackTop] = self.buffer[3]
            }
        }
        return self.primaryDiagnosis(repair)
    }



    //
    //  This function checks whether or not a given state has a
    // candidate, whose string representaion is a merging of the two
    // tokens at positions buffer_position and buffer_position+1 in
    // the buffer.  If so, it returns the candidate in question
    // otherwise it returns 0.
    //
    func (self *DiagnoseParser) mergeCandidate(state int , buffer_position int ) int  {
        var str: string = self.tokStream.getName(self.buffer[buffer_position]) + self.tokStream.getName(self.buffer[buffer_position + 1])
        for (var k int  = self.asi(state) self.asr(k) != 0 k++) {
            var i int  = self.terminalIndex(self.asr(k))
            if (str.length == self.name(i).length) {
                if (str.toLowerCase()==(self.name(i).toLowerCase())) {
                    return self.asr(k)
                }
            }
        }
        return 0
    }




    //
    // This procedure takes as arguments a parsing configuration
    // consisting of a state stack (stack and stack_top) and a fixed
    // number of input tokens (starting at buffer_position) in the
    // input BUFFER and some reference arguments: repair_code,
    // distance, misspell_index, candidate, and stack_position
    // which it sets based on the best possible recovery that it
    // finds in the given configuration.  The effectiveness of a
    // a repair is judged based on two criteria:
    //
    //       1) the number of tokens that can be parsed after the repair
    //              is applied: distance.
    //       2) how close to perfection is the candidate that is chosen:
    //              misspell_index.
    //
    // When self procedure is entered, distance, misspell_index and
    // repair_code are assumed to be initialized.
    //

    func (self *DiagnoseParser) checkPrimaryDistance(repair PrimaryRepairInfo, stck []int, stack_top int ) {
        //
        //  First, try scope recovery.
        //
        var scope_repair PrimaryRepairInfo = new PrimaryRepairInfo(repair)
        self.scopeTrial(scope_repair, stck, stack_top)
        if (scope_repair.distance > repair.distance) {
            repair.copy(scope_repair)
        }

        //
        //  Next, try merging the error token with its successor.
        //
        var symbol int  = self.mergeCandidate(stck[stack_top], repair.bufferPosition)
        if (symbol != 0) {
            var j int  = self.parseCheck(stck, stack_top, symbol, repair.bufferPosition + 2)
            if ((j > repair.distance) || (j == repair.distance && repair.misspellIndex < 10)) {
                repair.misspellIndex = 10
                repair.symbol = symbol
                repair.distance = j
                repair.code = MERGE_CODE
            }
        }

        //
        // Next, try deletion of the error token.
        //
        var j int  = self.parseCheck(stck,
                                        stack_top,
                                        self.tokStream.getKind(self.buffer[repair.bufferPosition + 1]),
                                        repair.bufferPosition + 2)

        var k int  = (self.tokStream.getKind(self.buffer[repair.bufferPosition]) == self.EOLT_SYMBOL &&
                         self.tokStream.afterEol(self.buffer[repair.bufferPosition + 1])
                         ? 10
                         : 0)

        if (j > repair.distance || (j == repair.distance && k > repair.misspellIndex)) {
            repair.misspellIndex = k
            repair.code = DELETION_CODE
            repair.distance = j
        }

        //
        // Update the error configuration by simulating all reduce and
        // goto actions induced by the error token. Then assign the top
        // most state of the new configuration to next_state.
        //
        var next_state int  = stck[stack_top],
            max_pos int  = stack_top
        self.tempStackTop = stack_top - 1

        self.tokStream.reset(self.buffer[repair.bufferPosition + 1])
        var tok int  = self.tokStream.getKind(self.buffer[repair.bufferPosition]),
            act int  = self.tAction(next_state, tok)
        for ;act <= self.NUM_RULES; {
            for ;;{
                var lhs_symbol int  = self.lhs(act)
                self.tempStackTop -= (self.rhs(act) - 1)
                act = (self.tempStackTop > max_pos
                                        ? self.tempStack[self.tempStackTop]
                                        : stck[self.tempStackTop])

                act = self.ntAction(act, lhs_symbol)
            } while (act <= self.NUM_RULES)
            max_pos = max_pos < self.tempStackTop ? max_pos : self.tempStackTop
            self.tempStack[self.tempStackTop + 1] = act
            next_state = act
            act = self.tAction(next_state, tok)
        }

        //
        //  Next, place the list of candidates in proper order.
        //
        var root int  = 0
        for (var i int  = self.asi(next_state) ;self.asr(i) != 0; i++) {
            symbol = self.asr(i)
            if (symbol != self.EOFT_SYMBOL && symbol != self.ERROR_SYMBOL) {
                if (root == 0) {
                    self.list[symbol] = symbol
                } else {
                    self.list[symbol] = self.list[root]
                    self.list[root] = symbol
                }
                root = symbol
            }
        }
        if (stck[stack_top] != next_state) {
            for (var i int  = self.asi(stck[stack_top]); self.asr(i) != 0 ;i++) {
                symbol = self.asr(i)
                if (symbol != self.EOFT_SYMBOL && symbol != self.ERROR_SYMBOL && self.list[symbol] == 0) {
                    if (root == 0) {
                        self.list[symbol] = symbol
                    } else {
                        self.list[symbol] = self.list[root]
                        self.list[root] = symbol
                    }
                    root = symbol
                }
            }
        }

        var head int  = self.list[root]
        self.list[root] = 0
        root = head

        //
        //  Next, try insertion for each possible candidate available in
        // the current state, except EOFT and ERROR_SYMBOL.
        //

        symbol = root
        while (symbol != 0) {
            var m int  = self.parseCheck(stck, stack_top, symbol, repair.bufferPosition),
                n int  = (symbol == self.EOLT_SYMBOL && self.tokStream.afterEol(self.buffer[repair.bufferPosition])
                        ? 10
                        : 0)
            if (m > repair.distance ||
                (m == repair.distance && n > repair.misspellIndex)) {
                repair.misspellIndex = n
                repair.distance = m
                repair.symbol = symbol
                repair.code = INSERTION_CODE
            }
            symbol = self.list[symbol]
        }

        //
        //  Next, Try substitution for each possible candidate available
        // in the current state, except EOFT and ERROR_SYMBOL.
        //
        symbol = root
        for ;symbol != 0; {

            var m int  = self.parseCheck(stck, stack_top, symbol, repair.bufferPosition + 1),
                n int  = (symbol == self.EOLT_SYMBOL && self.tokStream.afterEol(self.buffer[repair.bufferPosition + 1])
                    ? 10
                    : self.misspell(symbol, self.buffer[repair.bufferPosition]))
            if (m > repair.distance ||
                (m == repair.distance && n > repair.misspellIndex))
            {
                repair.misspellIndex = n
                repair.distance = m
                repair.symbol = symbol
                repair.code = SUBSTITUTION_CODE
            }
            var s int  = symbol
            symbol = self.list[symbol]
            self.list[s] = 0// reset element
        }

        //
        // Next, we try to insert a nonterminal candidate in front of the
        // error token, or substituting a nonterminal candidate for the
        // error token. Precedence is given to insertion.
        //
        var nt_index int
        for nt_index = self.nasi(stck[stack_top]); self.nasr(nt_index) != 0 ;nt_index++{
            symbol = self.nasr(nt_index) + self.NT_OFFSET
            var n int  = self.parseCheck(stck, stack_top, symbol, repair.bufferPosition + 1)
            if (n > repair.distance) {
                repair.misspellIndex = 0
                repair.distance = n
                repair.symbol = symbol
                repair.code = INVALID_CODE
            }

            n = self.parseCheck(stck, stack_top, symbol, repair.bufferPosition)
            if (n > repair.distance || (n == repair.distance && repair.code == INVALID_CODE))
            {
                repair.misspellIndex = 0
                repair.distance = n
                repair.symbol = symbol
                repair.code = INSERTION_CODE
            }
        }
        return
    }

    //
    // This procedure is invoked to issue a diagnostic message and
    // adjust the input buffer.  The recovery in question is either
    // the insertion of one or more scopes, the merging of the error
    // token with its successor, the deletion of the error token,
    // the insertion of a single token in front of the error token
    // or the substitution of another token for the error token.
    //
    func (self *DiagnoseParser) primaryDiagnosis(repair PrimaryRepairInfo): RepairCandidate {
        //
        //  Issue diagnostic.
        //
        var prevtok int  = self.buffer[repair.bufferPosition - 1],
            current_token int  = self.buffer[repair.bufferPosition]

        switch (repair.code) {
            case INSERTION_CODE:
            case BEFORE_CODE:
                {
                    var name_index int  = (repair.symbol > self.NT_OFFSET
                                                    ? self.getNtermIndex(self.stateStack[self.stateStackTop],
                                                                        repair.symbol,
                                                                        repair.bufferPosition)

                                                    : self.getTermIndex(self.stateStack,
                                                                        self.stateStackTop,
                                                                        repair.symbol,
                                                                        repair.bufferPosition))

                    var tok int  = (repair.code == INSERTION_CODE ? prevtok : current_token)
                    self.emitError(repair.code, name_index, tok, tok)
                }
                break
            case INVALID_CODE:
                {
                    var name_index int  = self.getNtermIndex(self.stateStack[self.stateStackTop],
                                                                repair.symbol,
                                                                repair.bufferPosition + 1)
                    self.emitError(repair.code, name_index, current_token, current_token)
                }
                break
            case SUBSTITUTION_CODE:
                {
                    var name_index int 
                    if (repair.misspellIndex >= 6) {
                        name_index = self.terminalIndex(repair.symbol)
                    } else {
                        name_index = self.getTermIndex(self.stateStack,
                            self.stateStackTop,
                            repair.symbol,
                            repair.bufferPosition + 1)
                        if (name_index != self.terminalIndex(repair.symbol)) {
                            repair.code = INVALID_CODE
                        }
                    }
                    self.emitError(repair.code, name_index, current_token, current_token)
                }
                break
            case MERGE_CODE:
                 self.emitError(repair.code,
                                self.terminalIndex(repair.symbol),
                                current_token,
                                self.tokStream.getNext(current_token))
                break
            case SCOPE_CODE:
            {
                for (var i int  = 0 i < self.scopeStackTop i++) {
                    self.emitError(repair.code,
                                    -self.scopeIndex[i],
                                    self.locationStack[self.scopePosition[i]],
                                    prevtok,
                                    self.nonterminalIndex(self.scopeLhs(self.scopeIndex[i])))
                }
                repair.symbol = self.scopeLhs(self.scopeIndex[self.scopeStackTop]) + self.NT_OFFSET
                self.stateStackTop = self.scopePosition[self.scopeStackTop]
                self.emitError(repair.code,
                                -self.scopeIndex[self.scopeStackTop],
                                self.locationStack[self.scopePosition[self.scopeStackTop]],
                                prevtok,
                                self.getNtermIndex(self.stateStack[self.stateStackTop], repair.symbol, repair.bufferPosition))
                break
            }
            default:// deletion
                self.emitError(repair.code, self.terminalIndex(self.ERROR_SYMBOL), current_token, current_token)
        }


        //
        //  Update buffer.
        //
        var candidate: RepairCandidate = new RepairCandidate()
        switch (repair.code) {
            case INSERTION_CODE:
            case BEFORE_CODE:
            case SCOPE_CODE:
                candidate.symbol = repair.symbol
                candidate.location = self.buffer[repair.bufferPosition]
                self.tokStream.reset(self.buffer[repair.bufferPosition])
                break
            case INVALID_CODE:
            case SUBSTITUTION_CODE:
                candidate.symbol = repair.symbol
                candidate.location = self.buffer[repair.bufferPosition]
                self.tokStream.reset(self.buffer[repair.bufferPosition + 1])
                break
            case MERGE_CODE:
                candidate.symbol = repair.symbol
                candidate.location = self.buffer[repair.bufferPosition]
                self.tokStream.reset(self.buffer[repair.bufferPosition + 2])
                break
            default:// deletion
                candidate.location = self.buffer[repair.bufferPosition + 1]
                candidate.symbol = self.tokStream.getKind(self.buffer[repair.bufferPosition + 1])
                self.tokStream.reset(self.buffer[repair.bufferPosition + 2])
                break
        }
        return candidate
    }

    //
    // This function takes as parameter an integer STACK_TOP that
    // points to a STACK element containing the state on which a
    // primary recovery will be made the terminal candidate on which
    // to recover and an integer: buffer_position, which points to
    // the position of the next input token in the BUFFER.  The
    // parser is simulated until a shift (or shift-reduce) action
    // is computed on the candidate.  Then we proceed to compute the
    // the name index of the highest level nonterminal that can
    // directly or indirectly produce the candidate.
    //
    func (self *DiagnoseParser) getTermIndex(stck []int, stack_top int , tok int , buffer_position int ) int  {
        //
        // Initialize stack index of temp_stack and initialize maximum
        // position of state stack that is still useful.
        //
        var act int  = stck[stack_top]
         var   max_pos int  = stack_top
         var   highest_symbol int  = tok

        self.tempStackTop = stack_top - 1

        //
        // Compute all reduce and associated actions induced by the
        // candidate until a SHIFT or SHIFT-REDUCE is computed. ERROR
        // and ACCEPT actions cannot be computed on the candidate in
        // self context, since we know that it is suitable for recovery.
        //
        self.tokStream.reset(self.buffer[buffer_position])
        act = self.tAction(act, tok)
        for ;act <= self.NUM_RULES; {
            //
            // Process all goto-reduce actions following reduction,
            // until a goto action is computed ...
            //
            for;; {
                var lhs_symbol int  = self.lhs(act)
                self.tempStackTop -= (self.rhs(act) - 1)
                act = (self.tempStackTop > max_pos
                                        ? self.tempStack[self.tempStackTop]
                                        : stck[self.tempStackTop])
                act = self.ntAction(act, lhs_symbol)
                if act <= self.NUM_RULES {
                    continue
                }
                else{
                    break
                }
            } 
            //
            // Compute new maximum useful position of (STATE_)stack,
            // push goto state into the stack, and compute next
            // action on candidate ...
            //
            max_pos = max_pos < self.tempStackTop ? max_pos : self.tempStackTop
            self.tempStack[self.tempStackTop + 1] = act
            act = self.tAction(act, tok)
        }

        //
        // At self stage, we have simulated all actions induced by the
        // candidate and we are ready to shift or shift-reduce it. First,
        // set tok and next_ptr appropriately and identify the candidate
        // as the initial highest_symbol. If a shift action was computed
        // on the candidate, update the stack and compute the next
        // action. Next, simulate all actions possible on the next input
        // token until we either have to shift it or are about to reduce
        // below the initial starting point in the stack (indicated by
        // max_pos as computed in the previous loop).  At that point,
        // return the highest_symbol computed.
        //
        self.tempStackTop++// adjust top of stack to reflect last goto
                                             // next move is shift or shift-reduce.

        var threshold int  = self.tempStackTop

        tok = self.tokStream.getKind(self.buffer[buffer_position])
        self.tokStream.reset(self.buffer[buffer_position + 1])

        if (act > self.ERROR_ACTION) {// shift-reduce on candidate?
            act -= self.ERROR_ACTION
        } else {
            if (act < self.ACCEPT_ACTION) {// shift on candidate
                self.tempStack[self.tempStackTop + 1] = act
                act = self.tAction(act, tok)
            }
        }
        for ;act <= self.NUM_RULES; {
            //
            // Process all goto-reduce actions following reduction,
            // until a goto action is computed ...
            //
            for;; {
                var lhs_symbol int  = self.lhs(act)
                self.tempStackTop -= (self.rhs(act) - 1)

                if (self.tempStackTop < threshold) {
                    return (highest_symbol  > self.NT_OFFSET
                                            ? self.nonterminalIndex(highest_symbol - self.NT_OFFSET)
                                            : self.terminalIndex(highest_symbol))
                }
                if (self.tempStackTop == threshold) {
                    highest_symbol = lhs_symbol + self.NT_OFFSET
                }
                act = (self.tempStackTop > max_pos
                                        ? self.tempStack[self.tempStackTop]
                                        : stck[self.tempStackTop])
                act = self.ntAction(act, lhs_symbol)
                if act <= self.NUM_RULES{
                    continue
                }
                else{
                    break
                }
            } 

            self.tempStack[self.tempStackTop + 1] = act
            act = self.tAction(act, tok)

        }
        return (highest_symbol  > self.NT_OFFSET
                                ? self.nonterminalIndex(highest_symbol - self.NT_OFFSET)
                                : self.terminalIndex(highest_symbol))
    }
    //
    // This function takes as parameter a starting state number:
    // start, a nonterminal symbol, A (candidate), and an integer,
    // buffer_position,  which points to the position of the next
    // input token in the BUFFER.
    // It returns the highest level non-terminal B such that
    // B =>*rm A.  I.e., there does not exists a nonterminal C such
    // that C =>+rm B. (Recall that for an LALR(k) grammar if
    // C =>+rm B, it cannot be the case that B =>+rm C)
    //
    func (self *DiagnoseParser) getNtermIndex(start int , sym int , buffer_position int ) int  {
        var highest_symbol int  = sym - self.NT_OFFSET,
            tok int  = self.tokStream.getKind(self.buffer[buffer_position])
        self.tokStream.reset(self.buffer[buffer_position + 1])

        //
        // Initialize stack index of temp_stack and initialize maximum
        // position of state stack that is still useful.
        //
        self.tempStackTop = 0
        self.tempStack[self.tempStackTop] = start

        var act int  = self.ntAction(start, highest_symbol)
        if (act > self.NUM_RULES) {// goto action?
            self.tempStack[self.tempStackTop + 1] = act
            act = self.tAction(act, tok)
        }

        for; act <= self.NUM_RULES; {
            //
            // Process all goto-reduce actions following reduction,
            // until a goto action is computed ...
            //
            for ;; {
                self.tempStackTop -= (self.rhs(act) - 1)
                if (self.tempStackTop < 0) {
                    return self.nonterminalIndex(highest_symbol)
                }
                if (self.tempStackTop == 0) {
                    highest_symbol = self.lhs(act)
                }
                act = self.ntAction(self.tempStack[self.tempStackTop], self.lhs(act))
                if  act <= self.NUM_RULES {
                    continue
                }else{
                    break
                }
            }
            self.tempStack[self.tempStackTop + 1] = act
            act = self.tAction(act, tok)
        }
        return self.nonterminalIndex(highest_symbol)
    }
    //
    //  Check whether or not there is a high probability that a
    // given string is a misspelling of another.
    // Certain singleton symbols (such as ":" and "") are also
    // considered to be misspellings of each other.
    //
    func (self *DiagnoseParser) misspell(sym int , tok int ) int  {
        //
        // Set up the two strings in question. Note that there is a "0"
        // gate added at the end of each string. This is important as
        // the algorithm assumes that it can "peek" at the symbol immediately
        // following the one that is being analysed.
        //
        var s1: string = (self.name(self.terminalIndex(sym))).toLowerCase()
        var n int  = s1.length
        s1 += '\u0000'
        var s2: string = (self.tokStream.getName(tok)).toLowerCase()
        var m int  = (s2.length < self.MAX_NAME_LENGTH ? s2.length : self.MAX_NAME_LENGTH)
        s2 = s2.substring(0, m) + '\u0000'

        //
        //  Singleton misspellings:
        //
        //        <---->     ,
        //
        //        <---->     :
        //
        //  .      <---->     ,
        //
        //  '      <---->     "
        //
        //
        if (n == 1 && m == 1) {
            if ((s1.charAt(0) == '' && s2.charAt(0) == ',') ||
                (s1.charAt(0) == ','  && s2.charAt(0) == '') ||
                (s1.charAt(0) == '' && s2.charAt(0) == ':') ||
                (s1.charAt(0) == ':'  && s2.charAt(0) == '') ||
                (s1.charAt(0) == '.'  && s2.charAt(0) == ',') ||
                (s1.charAt(0) == ','  && s2.charAt(0) == '.') ||
                (s1.charAt(0) == '\'' && s2.charAt(0) == '\"') ||
                (s1.charAt(0) == '\"' && s2.charAt(0) == '\'')) {
                return 3
            }
        }

        //
        // Scan the two strings. Increment "match" count for each match.
        // When a transposition is encountered, increase "match" count
        // by two but count it as one error. When a typo is found, skip
        // it and count it as one error. Otherwise we have a mismatch if
        // one of the strings is longer, increment its index, otherwise,
        // increment both indices and continue.
        //
        // This algorithm is an adaptation of a bool misspelling
        // algorithm proposed by Juergen Uhl.
        //
        var count int  = 0,
            prefix_length int  = 0,
            num_errors int  = 0

        var i int  = 0,
            j int  = 0

        for;(i < n) && (j < m) ; {
            if (s1.charAt(i) == s2.charAt(j)) {
                count++
                i++
                j++
                if (num_errors == 0) {
                    prefix_length++
                }
            }
            else if (s1.charAt(i + 1) == s2.charAt(j) && s1.charAt(i) == s2.charAt(j + 1))//transposition
            {
                count += 2
                i += 2
                j += 2
                num_errors++
            }
            else if (s1.charAt(i + 1) == s2.charAt(j + 1)) // mismatch
            {
                i += 2
                j += 2
                num_errors++
            }
            else
            {
                if ((n - i) > (m - j)) {
                    i++
                }
                else if ((m - j) > (n - i)) {
                    j++
                }
                else {
                    i++
                    j++
                }
                num_errors++
            }

        }

        if (i < n || j < m) {
            num_errors++
        }

        if (num_errors > ((n < m ? n : m) / 6 + 1)) {
            count = prefix_length
        }

        return Math.floor(count * 10 / ((n < s1.length ? s1.length : n) + num_errors))
    }

    func (self *DiagnoseParser) scopeTrial(repair PrimaryRepairInfo, stack []int, stack_top int ) {
        if (!self.stateSeen  || self.stateSeen.length == 0 || self.stateSeen.length < self.stateStack.length) {
            self.stateSeen = new []int(self.stateStack.length)
        }
        for (var i int  = 0 i < self.stateStack.length i++) {
            self.stateSeen[i] = DiagnoseParser.NIL
        }

        self.statePoolTop = 0
        if (!self.statePool ||self.statePool.length == 0 || self.statePool.length < self.stateStack.length) {
            self.statePool = new Array<StateInfo>(self.stateStack.length)
        }
        self.scopeTrialCheck(repair, stack, stack_top, 0)
        repair.code = SCOPE_CODE
        repair.misspellIndex = 10
        return
    }

    func (self *DiagnoseParser) scopeTrialCheck(repair PrimaryRepairInfo, stack []int, stack_top int , indx int ) {

        for (var i int  = self.stateSeen[stack_top] i != DiagnoseParser.NIL i = self.statePool[i].next) {
            if (self.statePool[i].state == stack[stack_top]) {
                return
            }
        }
        var old_state_pool_top int  = self.statePoolTop++
        if (self.statePoolTop >= self.statePool.length) {
            arraycopy(self.statePool, 0, self.statePool = new Array<StateInfo>(self.statePoolTop * 2), 0, self.statePoolTop)
        }

        self.statePool[old_state_pool_top] = new StateInfo(stack[stack_top], self.stateSeen[stack_top])
        self.stateSeen[stack_top] = old_state_pool_top

        var action: IntTuple = new IntTuple(1 << 3)
        for var i int  = 0; i < self.SCOPE_SIZE ;i++{
            //
            // Compute the action (or set of actions in case of conflicts) that
            // can be executed on the scope lookahead symbol. Save the action(s)
            // in the action tuple.
            //
            action.reset()
            var act int  = self.tAction(stack[stack_top], self.scopeLa(i))
            if (act > self.ACCEPT_ACTION && act < self.ERROR_ACTION)// conflicting actions?
            {
                for;; {
                    action.add(self.baseAction(act++))
                    if self.baseAction(act) != 0{
                        continue
                    }
                    else{
                        break
                    }
                } 
            } else {
                action.add(act)
            }

            //
            // For each action defined on the scope lookahead symbol,
            // try scope recovery.
            //
            var action_index int
            for (action_index  = 0; action_index < action.size(); action_index++{
                self.tokStream.reset(self.buffer[repair.bufferPosition])
                self.tempStackTop = stack_top - 1
                var max_pos int  = stack_top

                act = action.get(action_index)
                for;act <= self.NUM_RULES; {
                    //
                    // ... Process all goto-reduce actions following
                    // reduction, until a goto action is computed ...
                    //
                    for;; {
                        var lhs_symbol int  = self.lhs(act)
                        self.tempStackTop -= (self.rhs(act) - 1)
                        act = (self.tempStackTop > max_pos
                                                ? self.tempStack[self.tempStackTop]
                                                : stack[self.tempStackTop])

                        act = self.ntAction(act, lhs_symbol)
                        if act <= self.NUM_RULES {
                            continue
                        }else{
                            break
                        }
                    } 
                    if (self.tempStackTop + 1 >= self.stateStack.length) {
                        return
                    }
                    max_pos = max_pos < self.tempStackTop ? max_pos : self.tempStackTop
                    self.tempStack[self.tempStackTop + 1] = act
                    act = self.tAction(act, self.scopeLa(i))
                }
                //
                // If the lookahead symbol is parsable, then we check
                // whether or not we have a match between the scope
                // prefix and the transition symbols corresponding to
                // the states on top of the stack.
                //
                if (act != self.ERROR_ACTION) {
                    var j int 
                    var k int  = self.scopePrefix(i)
                    for (j = self.tempStackTop + 1
                        j >= (max_pos + 1) &&
                            self.inSymbol(self.tempStack[j]) == self.scopeRhs(k) j--) {
                        k++
                    }
                    if (j == max_pos) {
                        for (j = max_pos
                            j >= 1 && self.inSymbol(stack[j]) == self.scopeRhs(k)
                            j--) {
                            k++
                        }
                    }
                    //
                    // If the prefix matches, check whether the state
                    // newly exposed on top of the stack, (after the
                    // corresponding prefix states are popped from the
                    // stack), is in the set of "source states" for the
                    // scope in question and that it is at a position
                    // below the threshold indicated by MARKED_POS.
                    //
                    var marked_pos int  = (max_pos < stack_top ? max_pos + 1 : stack_top)
                    if (self.scopeRhs(k) == 0 && j < marked_pos) {// match?
                        var stack_position int  = j
                        for j = self.scopeStateSet(i);
                             stack[stack_position] != self.scopeState(j) &&
                             self.scopeState(j) != 0;
                             j++ {
                        }
                        //
                        // If the top state is valid for scope recovery,
                        // the left-hand side of the scope is used as
                        // starting symbol and we calculate how far the
                        // parser can advance within the forward context
                        // after parsing the left-hand symbol.
                        //
                        if (self.scopeState(j) != 0) {// state was found
                            var previous_distance int  = repair.distance,
                                distance int  = self.parseCheck(stack,
                                                                    stack_position,
                                                                    self.scopeLhs(i) + self.NT_OFFSET,
                                                                    repair.bufferPosition)

                            //
                            // if the recovery is not successful, we
                            // update the stack with all actions induced
                            // by the left-hand symbol, and recursively
                            // call SCOPE_TRIAL_CHECK to try again.
                            // Otherwise, the recovery is successful. If
                            // the new distance is greater than the
                            // initial SCOPE_DISTANCE, we update
                            // SCOPE_DISTANCE and set scope_stack_top to INDX
                            // to indicate the number of scopes that are
                            // to be applied for a succesful  recovery.
                            // NOTE that self procedure cannot get into
                            // an infinite loop, since each prefix match
                            // is guaranteed to take us to a lower point
                            // within the stack.
                            //
                            if ((distance - repair.bufferPosition + 1) < MIN_DISTANCE) {
                                var top int  = stack_position
                                act = self.ntAction(stack[top], self.scopeLhs(i))
                                for ;act <= self.NUM_RULES; {
                                    top -= (self.rhs(act) - 1)
                                    act = self.ntAction(stack[top], self.lhs(act))
                                }
                                top++
                                j = act
                                act = stack[top]// save
                                stack[top] = j// swap
                                self.scopeTrialCheck(repair, stack, top, indx + 1)
                                stack[top] = act// restore
                            }
                            else if (distance > repair.distance)
                            {
                                self.scopeStackTop = indx
                                repair.distance = distance
                            }

                            //
                            // If no other recovery possibility is left (due to
                            // backtracking and we are at the end of the input,
                            // then we favor a scope recovery over all other kinds
                            // of recovery.
                            //
                            if ( // TODO: main_configuration_stack.size() == 0 && // no other bactracking possibilities left
                                self.tokStream.getKind(self.buffer[repair.bufferPosition]) == self.EOFT_SYMBOL &&
                                repair.distance == previous_distance)
                            {
                                self.scopeStackTop = indx
                                repair.distance = MAX_DISTANCE
                            }
                            //
                            // If self scope recovery has beaten the
                            // previous distance, then we have found a
                            // better recovery (or self recovery is one
                            // of a list of scope recoveries). Record
                            // its information at the proper location
                            // (INDX) in SCOPE_INDEX and SCOPE_STACK.
                            //
                            if (repair.distance > previous_distance)
                            {
                                self.scopeIndex[indx] = i
                                self.scopePosition[indx] = stack_position
                                return
                            }
                        }
                    }
                }
            }
        }
    }
    //
    // This function computes the ParseCheck distance for the best
    // possible secondary recovery for a given configuration that
    // either deletes none or only one symbol in the forward context.
    // If the recovery found is more effective than the best primary
    // recovery previously computed, then the function returns true.
    // Only misplacement, scope and manual recoveries are attempted
    // simple insertion or substitution of a nonterminal are tried
    // in CHECK_PRIMARY_DISTANCE as part of primary recovery.
    //
    func (self *DiagnoseParser) secondaryCheck(stack []int, stack_top int , buffer_position int , distance int )  bool {
        var top int
        for top int  = stack_top - 1; top >= 0; top-- {
            var j int  = self.parseCheck(stack,
                                            top,
                                            self.tokStream.getKind(self.buffer[buffer_position]),
                                            buffer_position + 1)
            if (((j - buffer_position + 1) > MIN_DISTANCE) && (j > distance)) {
                return true
            }
        }

        var scope_repair PrimaryRepairInfo = new PrimaryRepairInfo()
        scope_repair.bufferPosition = buffer_position + 1
        scope_repair.distance = distance
        self.scopeTrial(scope_repair, stack, stack_top)
        return ((scope_repair.distance - buffer_position) > MIN_DISTANCE && scope_repair.distance > distance)
    }



    //
    // Secondary_phase is a bool function that checks whether or
    // not some form of secondary recovery is applicable to one of
    // the error configurations. First, if "next_stack" is available,
    // misplacement and secondary recoveries are attempted on it.
    // Then, in any case, these recoveries are attempted on "stack".
    // If a successful recovery is found, a diagnosis is issued, the
    // configuration is updated and the function returns "true".
    // Otherwise, the function returns false.
    //
    func (self *DiagnoseParser) secondaryPhase(error_token int ): RepairCandidate {
        var repair: SecondaryRepairInfo = new SecondaryRepairInfo(),
            misplaced_repair: SecondaryRepairInfo = new SecondaryRepairInfo()

        //
        // If the next_stack is available, try misplaced and secondary
        // recovery on it first.
        //
        var next_last_index int  = 0
        if (self.nextStackTop >= 0) {

            var save_location int 

            self.buffer[2] = error_token
            self.buffer[1] = self.tokStream.getPrevious(self.buffer[2])
            self.buffer[0] = self.tokStream.getPrevious(self.buffer[1])

            for (var k int  = 3; k < BUFF_UBOUND ;k++) {
                self.buffer[k] = self.tokStream.getNext(self.buffer[k - 1])
            }

            self.buffer[BUFF_UBOUND] = self.tokStream.badToken()// elmt not available
            //
            // If we are at the end of the input stream, compute the
            // index position of the first EOFT symbol (last useful
            // index).
            //
            for (next_last_index = MAX_DISTANCE - 1
                next_last_index >= 1 &&
                self.tokStream.getKind(self.buffer[next_last_index]) == self.EOFT_SYMBOL
                next_last_index--) {
            }

            next_last_index = next_last_index + 1

            save_location = self.locationStack[self.nextStackTop]
            self.locationStack[self.nextStackTop] = self.buffer[2]
            misplaced_repair.numDeletions = self.nextStackTop
            self.misplacementRecovery(misplaced_repair, self.nextStack, self.nextStackTop, next_last_index, true)
            if (misplaced_repair.recoveryOnNextStack) {
                misplaced_repair.distance++
            }
            repair.numDeletions = self.nextStackTop + BUFF_UBOUND
            self.secondaryRecovery( repair,
                                    self.nextStack,
                                    self.nextStackTop,
                                    next_last_index, true)

            if (repair.recoveryOnNextStack) {
                repair.distance++
            }
            self.locationStack[self.nextStackTop] = save_location
        } else {// next_stack not available, initialize ...
            misplaced_repair.numDeletions = self.stateStackTop
            repair.numDeletions = self.stateStackTop + BUFF_UBOUND
        }

        //
        // Try secondary recovery on the "stack" configuration.
        //
        self.buffer[3] = error_token

        self.buffer[2] = self.tokStream.getPrevious(self.buffer[3])
        self.buffer[1] = self.tokStream.getPrevious(self.buffer[2])
        self.buffer[0] = self.tokStream.getPrevious(self.buffer[1])

        for (var k int  = 4 k < BUFF_SIZE k++) {
            self.buffer[k] = self.tokStream.getNext(self.buffer[k - 1])
        }

        var last_index int 
        for (last_index = MAX_DISTANCE - 1
            last_index >= 1 &&
            self.tokStream.getKind(self.buffer[last_index]) == self.EOFT_SYMBOL
            last_index--) {
        }
        last_index++

        self.misplacementRecovery(misplaced_repair, self.stateStack, self.stateStackTop, last_index, false)

        self.secondaryRecovery(repair, self.stateStack, self.stateStackTop, last_index, false)

        //
        // If a successful misplaced recovery was found, compare it with
        // the most successful secondary recovery.  If the misplaced
        // recovery either deletes fewer symbols or parse-checks further
        // then it is chosen.
        //
        if (misplaced_repair.distance > MIN_DISTANCE)
        {
            if (misplaced_repair.numDeletions <= repair.numDeletions ||
                (misplaced_repair.distance - misplaced_repair.numDeletions) >=
                (repair.distance - repair.numDeletions))
            {
                repair.code = MISPLACED_CODE
                repair.stackPosition = misplaced_repair.stackPosition
                repair.bufferPosition = 2
                repair.numDeletions = misplaced_repair.numDeletions
                repair.distance = misplaced_repair.distance
                repair.recoveryOnNextStack = misplaced_repair.recoveryOnNextStack
            }
        }


        //
        // If the successful recovery was on next_stack, update: stack,
        // buffer, location_stack and last_index.
        //
        if (repair.recoveryOnNextStack)
        {
            self.stateStackTop = self.nextStackTop
            arraycopy(self.nextStack, 0, self.stateStack, 0, self.stateStackTop + 1)

            self.buffer[2] = error_token
            self.buffer[1] = self.tokStream.getPrevious(self.buffer[2])
            self.buffer[0] = self.tokStream.getPrevious(self.buffer[1])

            for (var k int  = 3 ;k < BUFF_UBOUND ;k++) {
                self.buffer[k] = self.tokStream.getNext(self.buffer[k - 1])
            }

            self.buffer[BUFF_UBOUND] = self.tokStream.badToken()// elmt not available

            self.locationStack[self.nextStackTop] = self.buffer[2]
            last_index = next_last_index
        }

        //
        // Next, try scope recoveries after deletion of one, two, three,
        // four ... buffer_position tokens from the input stream.
        //
        if (repair.code == SECONDARY_CODE || repair.code == DELETION_CODE) {
            var scope_repair PrimaryRepairInfo = new PrimaryRepairInfo()
            for scope_repair.bufferPosition = 2;
                scope_repair.bufferPosition <= repair.bufferPosition &&
                repair.code != SCOPE_CODE ;scope_repair.bufferPosition++{
                self.scopeTrial(scope_repair, self.stateStack, self.stateStackTop)
                var j int  = (scope_repair.distance == MAX_DISTANCE
                                                        ? last_index
                                                        : scope_repair.distance),
                    k int  = scope_repair.bufferPosition - 1
                if ((scope_repair.distance - k) > MIN_DISTANCE && (j - k) > (repair.distance - repair.numDeletions)) {
                    var i int  = self.scopeIndex[self.scopeStackTop]// upper bound
                    repair.code = SCOPE_CODE
                    repair.symbol = self.scopeLhs(i) + self.NT_OFFSET
                    repair.stackPosition = self.stateStackTop
                    repair.bufferPosition = scope_repair.bufferPosition
                }
            }
        }
        //
        // If a successful repair was not found, quit!  Otherwise, issue
        // diagnosis and adjust configuration...
        //
        var candidate: RepairCandidate = new RepairCandidate()
        if (repair.code == 0) {
            return candidate
        }
        self.secondaryDiagnosis(repair)

        //
        // Update buffer based on number of elements that are deleted.
        //
        switch (repair.code) {
            case MISPLACED_CODE:
                candidate.location = self.buffer[2]
                candidate.symbol = self.tokStream.getKind(self.buffer[2])
                self.tokStream.reset(self.tokStream.getNext(self.buffer[2]))
                break
            case DELETION_CODE:
                candidate.location = self.buffer[repair.bufferPosition]
                candidate.symbol = self.tokStream.getKind(self.buffer[repair.bufferPosition])
                self.tokStream.reset(self.tokStream.getNext(self.buffer[repair.bufferPosition]))
                break
            default:// SCOPE_CODE || SECONDARY_CODE
                candidate.symbol = repair.symbol
                candidate.location = self.buffer[repair.bufferPosition]
                self.tokStream.reset(self.buffer[repair.bufferPosition])
                break
        }
        return candidate
    }

    //
    // This bool function checks whether or not a given
    // configuration yields a better misplacement recovery than
    // the best misplacement recovery computed previously.
    //
    func (self *DiagnoseParser) misplacementRecovery(repair: SecondaryRepairInfo, stack []int, stack_top int , last_index int , stack_flag  bool) {
        var previous_loc int  = self.buffer[2]
        var   stack_deletions int  = 0
        var top int
        for top  = stack_top - 1; top >= 0; top-- {
            if (self.locationStack[top] < previous_loc) {
                stack_deletions++
            }
            previous_loc = self.locationStack[top]

            var parse_distance int  = self.parseCheck(stack, top, self.tokStream.getKind(self.buffer[2]), 3),
                j int  = (parse_distance == MAX_DISTANCE ? last_index : parse_distance)
            if ((parse_distance > MIN_DISTANCE) && (j - stack_deletions) > (repair.distance - repair.numDeletions)) {
                repair.stackPosition = top
                repair.distance = j
                repair.numDeletions = stack_deletions
                repair.recoveryOnNextStack = stack_flag
            }
        }
        return
    }

    //
    // This function checks whether or not a given
    // configuration yields a better secondary recovery than the
    // best misplacement recovery computed previously.
    //
    func (self *DiagnoseParser) secondaryRecovery(repair: SecondaryRepairInfo, stack []int, stack_top int , last_index int , stack_flag  bool) {
        var previous_loc int  = self.buffer[2],
            stack_deletions int  = 0

        for (var top int  = stack_top top >= 0 && repair.numDeletions >= stack_deletions top--) {
            if (self.locationStack[top] < previous_loc) {
                stack_deletions++
            }
            previous_loc = self.locationStack[top]
            var i int  
            for  i   = 2;
                i <= (last_index - MIN_DISTANCE + 1) &&
                (repair.numDeletions >= (stack_deletions + i - 1)); i++
            {
                var parse_distance int  = self.parseCheck(stack, top, self.tokStream.getKind(self.buffer[i]), i + 1),
                    j int  = (parse_distance == MAX_DISTANCE ? last_index : parse_distance)

                if ((parse_distance - i + 1) > MIN_DISTANCE) {

                    var k int  = stack_deletions + i - 1
                    if ((k < repair.numDeletions) ||
                        (j - k) > (repair.distance - repair.numDeletions) ||
                        ((repair.code == SECONDARY_CODE) && (j - k) == (repair.distance - repair.numDeletions)))
                    {
                        repair.code = DELETION_CODE
                        repair.distance = j
                        repair.stackPosition = top
                        repair.bufferPosition = i
                        repair.numDeletions = k
                        repair.recoveryOnNextStack = stack_flag
                    }
                }
                var l int
                for (l = self.nasi(stack[top]) ;l >= 0 && self.nasr(l) != 0 ;l++) {
                    var symbol int  = self.nasr(l) + self.NT_OFFSET
                    parse_distance = self.parseCheck(stack, top, symbol, i)
                    j = (parse_distance == MAX_DISTANCE ? last_index : parse_distance)

                    if ((parse_distance - i + 1) > MIN_DISTANCE)
                    {
                        var k int  = stack_deletions + i - 1
                        if (k < repair.numDeletions || (j - k) > (repair.distance - repair.numDeletions)) {
                            repair.code = SECONDARY_CODE
                            repair.symbol = symbol
                            repair.distance = j
                            repair.stackPosition = top
                            repair.bufferPosition = i
                            repair.numDeletions = k
                            repair.recoveryOnNextStack = stack_flag
                        }
                    }
                }
            }
        }
        return
    }
    //
    // This procedure is invoked to issue a secondary diagnosis and
    // adjust the input buffer.  The recovery in question is either
    // an automatic scope recovery, a manual scope recovery, a
    // secondary substitution or a secondary deletion.
    //
    func (self *DiagnoseParser) secondaryDiagnosis(repair: SecondaryRepairInfo) {
        switch (repair.code) {
            case SCOPE_CODE:
                if (repair.stackPosition < self.stateStackTop) {
                    self.emitError( DELETION_CODE,
                                    self.terminalIndex(self.ERROR_SYMBOL),
                                    self.locationStack[repair.stackPosition],
                                    self.buffer[1])
                }
                for (var i int  = 0 i < self.scopeStackTop i++) {
                    self.emitError(SCOPE_CODE,
                        -self.scopeIndex[i],
                        self.locationStack[self.scopePosition[i]],
                        self.buffer[1],
                        self.nonterminalIndex(self.scopeLhs(self.scopeIndex[i])))
                }

                repair.symbol = self.scopeLhs(self.scopeIndex[self.scopeStackTop]) + self.NT_OFFSET
                self.stateStackTop = self.scopePosition[self.scopeStackTop]
                self.emitError(SCOPE_CODE,
                    -self.scopeIndex[self.scopeStackTop],
                    self.locationStack[self.scopePosition[self.scopeStackTop]],
                    self.buffer[1],
                    self.getNtermIndex(self.stateStack[self.stateStackTop],
                                        repair.symbol,
                                        repair.bufferPosition))
                break
            default:
                self.emitError(repair.code,
                    (repair.code == SECONDARY_CODE
                                    ? self.getNtermIndex(self.stateStack[repair.stackPosition],
                                                        repair.symbol,
                                                        repair.bufferPosition)
                                    : self.terminalIndex(self.ERROR_SYMBOL)),
                    self.locationStack[repair.stackPosition],
                    self.buffer[repair.bufferPosition - 1])
                self.stateStackTop = repair.stackPosition
        }
        return
    }

  
    //
    // This method is invoked by an LPG PARSER or a semantic
    // routine to process an error message.
    //

    func (self *DiagnoseParser) emitError(msg_code int , name_index int , left_token int , right_token int , scope_name_index int =0) {
        var left_token_loc int  = (left_token > right_token ? right_token : left_token),
            right_token_loc int  = right_token
        
        var   token_name = (name_index >= 0 &&
            !(self.name(name_index).toUpperCase() == "ERROR")
                ? "\"" + self.name(name_index) + "\""
                : "")
        
        if (msg_code == INVALID_CODE) {
            msg_code = token_name.length == 0 ? INVALID_CODE : INVALID_TOKEN_CODE
        }
        if (msg_code == SCOPE_CODE) {
            token_name = "\""
            for (var i int  = self.scopeSuffix(-<number>name_index) self.scopeRhs(i) != 0 i++) {

                if (!self.isNullable(self.scopeRhs(i))) {
                    var symbol_index int  = (self.scopeRhs(i) > self.NT_OFFSET
                                                                ? self.nonterminalIndex(self.scopeRhs(i) - self.NT_OFFSET)
                                                                : self.terminalIndex(self.scopeRhs(i)))

                    if (self.name(symbol_index).length > 0) {
                        if (token_name.length > 1) {// Not just starting quote?
                            token_name += " "// add a space separator
                        }
                        token_name += self.name(symbol_index)
                    }
                }
            }
            token_name += "\""
        }
        self.tokStream.reportError(msg_code, left_token, right_token, token_name)
        return
    }
    //
    // keep looking ahead until we compute a valid action
    //
    private lookahead(act int , token int ) int  {
        act = self.prs.lookAhead(act - self.LA_STATE_OFFSET, self.tokStream.getKind(token))
        return (act > self.LA_STATE_OFFSET
                    ? self.lookahead(act, self.tokStream.getNext(token))
                    : act)
    }
    //
    // Compute the next action defined on act and sym. If self
    // action requires more lookahead, these lookahead symbols
    // are in the token stream beginning at the next token that
    // is yielded by peek().
    //
    func (self *DiagnoseParser) tAction(act int , sym int ) int 
    {
        act = self.prs.tAction(act, sym)
        return (act > self.LA_STATE_OFFSET
                    ? self.lookahead(act, self.tokStream.peek())
                    : act)
    }
}
