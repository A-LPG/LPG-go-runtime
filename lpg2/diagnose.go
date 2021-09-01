package lpg2

type RepairCandidate struct {
	symbol   int
	location int
}

func NewRepairCandidate() *RepairCandidate {
	t := new(RepairCandidate)
	t.symbol = 0
	t.location = 0
	return t
}

type PrimaryRepairInfo struct {
	distance       int
	misspellIndex  int
	code           int
	bufferPosition int
	symbol         int
}

func NewPrimaryRepairInfo() *PrimaryRepairInfo {
	t := new(PrimaryRepairInfo)
	return t
}
func NewPrimaryRepairInfoAndClone(clone *PrimaryRepairInfo) *PrimaryRepairInfo {
	t := new(PrimaryRepairInfo)
	t.copy(clone)
	return t
}
func (my *PrimaryRepairInfo) copy(clone *PrimaryRepairInfo) {
	my.distance = clone.distance
	my.misspellIndex = clone.misspellIndex
	my.code = clone.code
	my.bufferPosition = clone.bufferPosition
	my.symbol = clone.symbol
	return
}

type SecondaryRepairInfo struct {
	code                int
	distance            int
	bufferPosition      int
	stackPosition       int
	numDeletions        int
	symbol              int
	recoveryOnNextStack bool
}

func NewSecondaryRepairInfo() *SecondaryRepairInfo {
	t := new(SecondaryRepairInfo)
	return t
}

type StateInfo struct {
	state int
	next  int
}

func StateInfoArraycopy(src []*StateInfo, srcPos int,
	dest []*StateInfo, destPos int, length int) []*StateInfo {
	for i := 0; i < length; i++ {
		dest[destPos+i] = src[srcPos+i]
	}
	return dest
}
func NewStateInfo(state int, next int) *StateInfo {
	t := new(StateInfo)
	t.state = state
	t.next = next
	return t
}


const BUFF_UBOUND int = 31
const BUFF_SIZE int = 32
const MAX_DISTANCE int = 30
const MIN_DISTANCE int = 3
const NIL int = -1

type DiagnoseParser struct {
	monitor   Monitor
	tokStream TokenStream

	prs ParseTable

	ERROR_SYMBOL    int
	SCOPE_SIZE      int
	MAX_NAME_LENGTH int
	NT_OFFSET       int
	LA_STATE_OFFSET int
	NUM_RULES       int
	NUM_SYMBOLS     int
	START_STATE     int
	EOFT_SYMBOL     int
	EOLT_SYMBOL     int
	ACCEPT_ACTION   int
	ERROR_ACTION    int

	list []int

	maxErrors int

	maxTime int

	stateStackTop int
	stateStack    []int

	locationStack []int

	tempStackTop int
	tempStack    []int

	prevStackTop int
	prevStack    []int

	nextStackTop int
	nextStack    []int

	scopeStackTop int
	scopeIndex    []int
	scopePosition []int

	buffer []int

	stateSeen []int

	statePoolTop             int
	statePool                []*StateInfo
	main_configuration_stack *ConfigurationStack
	STACK_INCREMENT int
}

func NewDiagnoseParser(tokStream TokenStream, prs ParseTable, maxErrors int, maxTime int, monitor Monitor) *DiagnoseParser {
	my := new(DiagnoseParser)
	my.STACK_INCREMENT  = 256
	my.monitor = monitor
	my.maxErrors = maxErrors
	my.maxTime = maxTime
	my.tokStream = tokStream
	my.prs = prs
	my.main_configuration_stack = NewConfigurationStack(prs)
	my.ERROR_SYMBOL = prs.GetErrorSymbol()
	my.SCOPE_SIZE = prs.GetScopeSize()
	my.MAX_NAME_LENGTH = prs.GetMaxNameLength()
	my.NT_OFFSET = prs.GetNtOffSet()
	my.LA_STATE_OFFSET = prs.GetLaStateOffSet()
	my.NUM_RULES = prs.GetNumRules()
	my.NUM_SYMBOLS = prs.GetNumSymbols()
	my.START_STATE = prs.GetStartState()
	my.EOFT_SYMBOL = prs.GetEoftSymbol()
	my.EOLT_SYMBOL = prs.GetEoltSymbol()
	my.ACCEPT_ACTION = prs.GetAcceptAction()
	my.ERROR_ACTION = prs.GetErrorAction()
	my.list = make([]int, my.NUM_SYMBOLS+1)
	return my
}

func (my *DiagnoseParser) SetMonitor(monitor Monitor) {
	my.monitor = monitor
}

func (my *DiagnoseParser) Rhs(index int) int {
	return my.prs.Rhs(index)
}
func (my *DiagnoseParser) BaseAction(index int) int {
	return my.prs.BaseAction(index)
}
func (my *DiagnoseParser) BaseCheck(index int) int {
	return my.prs.BaseCheck(index)
}
func (my *DiagnoseParser) Lhs(index int) int {
	return my.prs.Lhs(index)
}
func (my *DiagnoseParser) TermCheck(index int) int {
	return my.prs.TermCheck(index)
}
func (my *DiagnoseParser) TermAction(index int) int {
	return my.prs.TermAction(index)
}
func (my *DiagnoseParser) Asb(index int) int {
	return my.prs.Asb(index)
}
func (my *DiagnoseParser) asr(index int) int {
	return my.prs.Asr(index)
}
func (my *DiagnoseParser) Nasb(index int) int {
	return my.prs.Nasb(index)
}
func (my *DiagnoseParser) Nasr(index int) int {
	return my.prs.Nasr(index)
}
func (my *DiagnoseParser) TerminalIndex(index int) int {
	return my.prs.TerminalIndex(index)
}
func (my *DiagnoseParser) NonterminalIndex(index int) int {
	return my.prs.NonterminalIndex(index)
}

func (my *DiagnoseParser) SymbolIndex(index int) int {

	if index > my.NT_OFFSET {
		return my.NonterminalIndex(index - my.NT_OFFSET)
	} else {
		return my.TerminalIndex(index)
	}
}
func (my *DiagnoseParser) ScopePrefix(index int) int {
	return my.prs.ScopePrefix(index)
}
func (my *DiagnoseParser) ScopeSuffix(index int) int {
	return my.prs.ScopeSuffix(index)
}
func (my *DiagnoseParser) ScopeLhs(index int) int {
	return my.prs.ScopeLhs(index)
}
func (my *DiagnoseParser) ScopeLa(index int) int {
	return my.prs.ScopeLa(index)
}
func (my *DiagnoseParser) ScopeStateSet(index int) int {
	return my.prs.ScopeStateSet(index)
}
func (my *DiagnoseParser) ScopeRhs(index int) int {
	return my.prs.ScopeRhs(index)
}
func (my *DiagnoseParser) ScopeState(index int) int {
	return my.prs.ScopeState(index)
}
func (my *DiagnoseParser) InSymb(index int) int {
	return my.prs.InSymb(index)
}
func (my *DiagnoseParser) Name(index int) string {
	return my.prs.Name(index)
}
func (my *DiagnoseParser) OriginalState(state int) int {
	return my.prs.OriginalState(state)
}
func (my *DiagnoseParser) Asi(state int) int {
	return my.prs.Asi(state)
}
func (my *DiagnoseParser) Nasi(state int) int {
	return my.prs.Nasi(state)
}
func (my *DiagnoseParser) InSymbol(state int) int {
	return my.prs.InSymbol(state)
}
func (my *DiagnoseParser) NtAction(state int, sym int) int {
	return my.prs.NtAction(state, sym)
}
func (my *DiagnoseParser) IsNullable(symbol int) bool {
	return my.prs.IsNullable(symbol)
}

func (my *DiagnoseParser) ReallocateStacks() {
	var old_stack_length int = len(my.stateStack)
	var stack_length int = old_stack_length + my.STACK_INCREMENT

	if len(my.stateStack) == 0 {
		my.stateStack = make([]int, stack_length)
		my.locationStack = make([]int, stack_length)
		my.tempStack = make([]int, stack_length)
		my.prevStack = make([]int, stack_length)
		my.nextStack = make([]int, stack_length)
		my.scopeIndex = make([]int, stack_length)
		my.scopePosition = make([]int, stack_length)
	} else {
		my.stateStack = Arraycopy(my.stateStack, 0, make([]int, stack_length), 0, old_stack_length)
		my.locationStack = Arraycopy(my.locationStack, 0, make([]int, stack_length), 0, old_stack_length)
		my.tempStack = Arraycopy(my.tempStack, 0, make([]int, stack_length), 0, old_stack_length)
		my.prevStack = Arraycopy(my.prevStack, 0, make([]int, stack_length), 0, old_stack_length)
		my.nextStack = Arraycopy(my.nextStack, 0, make([]int, stack_length), 0, old_stack_length)
		my.scopeIndex = Arraycopy(my.scopeIndex, 0, make([]int, stack_length), 0, old_stack_length)
		my.scopePosition = Arraycopy(my.scopePosition, 0, make([]int, stack_length), 0, old_stack_length)
	}

}

func (my *DiagnoseParser) Diagnose(error_token int) {
	my.DiagnoseEntry2(0, error_token)
}

func (my *DiagnoseParser) DiagnoseEntry1(marker_kind int) {
	my.ReallocateStacks()
	my.tempStackTop = 0
	my.tempStack[my.tempStackTop] = my.START_STATE
	my.tokStream.Reset()
	var current_token int
	var current_kind int
	if marker_kind == 0 {
		current_token = my.tokStream.GetToken()
		current_kind = my.tokStream.GetKind(current_token)
	} else {
		current_token = my.tokStream.Peek()
		current_kind = marker_kind
	}

	//
	// If an error was found, start the diagnosis and recovery.
	//
	var error_token int = my.ParseForError(current_kind)
	if error_token != 0 {
		my.DiagnoseEntry2(marker_kind, error_token)
	}
	return
}
func (my *DiagnoseParser) DiagnoseEntry2(marker_kind int, error_token int) {
	var action = NewIntTupleWithEstimate(1 << 18)
	var startTime int = Now()
	var errorCount int = 0

	//
	// Compute sequence of actions that leads us to the
	// error_token.
	//
	if len(my.stateStack) == 0 {
		my.ReallocateStacks()
	}

	my.tempStackTop = 0
	my.tempStack[my.tempStackTop] = my.START_STATE
	my.tokStream.Reset()
	var current_token int
	var current_kind int
	if marker_kind == 0 {
		current_token = my.tokStream.GetToken()
		current_kind = my.tokStream.GetKind(current_token)
	} else {
		current_token = my.tokStream.Peek()
		current_kind = marker_kind
	}
	my.ParseUpToError(action, current_kind, error_token)

	//
	// Start parsing
	//
	my.stateStackTop = 0
	my.stateStack[my.stateStackTop] = my.START_STATE

	my.tempStackTop = my.stateStackTop
	Arraycopy(my.tempStack, 0, my.stateStack, 0, my.tempStackTop+1)

	my.tokStream.Reset()
	if marker_kind == 0 {
		current_token = my.tokStream.GetToken()
		current_kind = my.tokStream.GetKind(current_token)
	} else {
		current_token = my.tokStream.Peek()
		current_kind = marker_kind
	}
	my.locationStack[my.stateStackTop] = current_token

	//
	// Process a terminal
	//
	var act int = -1// make it different
	for act != my.ACCEPT_ACTION {
		//
		// Synchronize state stacks and update the location stack
		//
		var prev_pos int = -1
		my.prevStackTop = -1

		var next_pos int = -1
		my.nextStackTop = -1

		var pos int = my.stateStackTop
		my.tempStackTop = my.stateStackTop - 1
		Arraycopy(my.stateStack, 0, my.tempStack, 0, my.stateStackTop+1)

		var action_index int = 0
		act = action.Get(action_index) // TAction(act, current_kind)
		action_index++
		//
		// When a reduce action is encountered, we compute all REDUCE
		// and associated goto actions induced by the current token.
		// Eventually, a SHIFT, SHIFT-REDUCE, ACCEPT or ERROR action is
		// computed...
		//
		for act <= my.NUM_RULES {
			for {
				my.tempStackTop -= my.Rhs(act) - 1
				act = my.NtAction(my.tempStack[my.tempStackTop], my.Lhs(act))
				if act <= my.NUM_RULES {
					continue
				} else {
					break
				}
			}

			//
			// ... Update the maximum useful position of the
			// (STATE_)STACK, push goto state into stack, and
			// compute next action on current symbol ...
			//
			if my.tempStackTop+1 >= len(my.stateStack) {
				my.ReallocateStacks()
			}
			if !(pos < my.tempStackTop) {
				pos = my.tempStackTop
			}
			my.tempStack[my.tempStackTop+1] = act
			act = action.Get(action_index) // TAction(act, current_kind)
			action_index++
		}
		//
		// At my point, we have a shift, shift-reduce, accept or error
		// action.  STACK contains the configuration of the state stack
		// prior to executing any action on current_token. next_stack contains
		// the configuration of the state stack after executing all
		// reduce actions induced by current_token.  The variable pos indicates
		// the highest position in STACK that is still useful after the
		// reductions are executed.
		//
		for act > my.ERROR_ACTION || act < my.ACCEPT_ACTION {

			//
			// if the parser needs to stop processing,
			// it may do so here.
			//
			if my.monitor != nil && my.monitor.IsCancelled() {
				return
			}

			my.nextStackTop = my.tempStackTop + 1
			var i int
			for i = next_pos + 1; i <= my.nextStackTop; i++ {
				my.nextStack[i] = my.tempStack[i]
			}
			var k int
			for k = pos + 1; k <= my.nextStackTop; k++ {
				my.locationStack[k] = my.locationStack[my.stateStackTop]
			}

			//
			// If we have a shift-reduce, process it as well as
			// the goto-reduce actions that follow it.
			//
			if act > my.ERROR_ACTION {
				act -= my.ERROR_ACTION
				for {
					my.nextStackTop -= (my.Rhs(act) - 1)
					act = my.NtAction(my.nextStack[my.nextStackTop], my.Lhs(act))
					if act <= my.NUM_RULES {
						continue
					} else {
						break
					}
				}
				if !(pos < my.nextStackTop) {
					pos = my.nextStackTop
				}
			}

			if my.nextStackTop+1 >= len(my.stateStack) {
				my.ReallocateStacks()
			}

			my.tempStackTop = my.nextStackTop

			my.nextStackTop += 1
			my.nextStack[my.nextStackTop] = act

			next_pos = my.nextStackTop
			//
			// Simulate the parser through the next token without
			// destroying STACK or next_stack.
			//
			current_token = my.tokStream.GetToken()
			current_kind = my.tokStream.GetKind(current_token)
			act = action.Get(action_index) // TAction(act, current_kind)
			action_index++
			for act <= my.NUM_RULES {
				//
				// ... Process all goto-reduce actions following
				// reduction, until a goto action is computed ...
				//

				for {
					var lhs_symbol int = my.Lhs(act)
					my.tempStackTop -= (my.Rhs(act) - 1)
					if my.tempStackTop > next_pos {
						act = my.tempStack[my.tempStackTop]
					} else {
						act = my.nextStack[my.tempStackTop]
					}

					act = my.NtAction(act, lhs_symbol)
					if act <= my.NUM_RULES {
						continue
					} else {
						break
					}
				}
				//
				// ... Update the maximum useful position of the
				// (STATE_)STACK, push GOTO state into stack, and
				// compute next action on current symbol ...
				//
				if my.tempStackTop+1 >= len(my.stateStack) {
					my.ReallocateStacks()
				}
				if !(next_pos < my.tempStackTop) {
					next_pos = my.tempStackTop
				}
				my.tempStack[my.tempStackTop+1] = act

				act = action.Get(action_index) // TAction(act, current_kind)
				action_index++

			}
			//
			// No error was detected, Read next token into
			// PREVTOK element, advance CURRENT_TOKEN pointer and
			// update stacks.
			//
			if act != my.ERROR_ACTION {
				my.prevStackTop = my.stateStackTop
				var i int = prev_pos + 1
				for ; i <= my.prevStackTop; i++ {
					my.prevStack[i] = my.stateStack[i]
				}
				prev_pos = pos

				my.stateStackTop = my.nextStackTop
				var k int = pos + 1
				for ; k <= my.stateStackTop; k++ {
					my.stateStack[k] = my.nextStack[k]
				}
				my.locationStack[my.stateStackTop] = current_token
				pos = next_pos
			}
		}

		//
		// At my stage, either we have an ACCEPT or an ERROR
		// action.
		//
		if act == my.ERROR_ACTION {
			//
			// An error was detected.
			//
			errorCount += 1
			//
			// Check time and error limits after the first error encountered
			// Exit if number of errors exceeds maxError or if maxTime reached
			//
			if errorCount > 1 {
				if my.maxErrors > 0 && errorCount > my.maxErrors {
					break
				}
				if my.maxTime > 0 && Now()-startTime > my.maxTime {
					break
				}
			}
			var candidate = my.ErrorRecovery(current_token)
			//
			// if the parser needs to stop processing,
			// it may do so here.
			//
			if my.monitor != nil && my.monitor.IsCancelled() {
				return
			}
			act = my.stateStack[my.stateStackTop]

			//
			// If the recovery was successful on a nonterminal candidate,
			// parse through that candidate and "read" the next token.
			//
			if candidate.symbol == 0 {
				break
			} else {
				if candidate.symbol > my.NT_OFFSET {
					var lhs_symbol int = candidate.symbol - my.NT_OFFSET
					act = my.NtAction(act, lhs_symbol)
					for act <= my.NUM_RULES {
						my.stateStackTop -= (my.Rhs(act) - 1)
						act = my.NtAction(my.stateStack[my.stateStackTop], my.Lhs(act))
					}

					my.stateStackTop += 1
					my.stateStack[my.stateStackTop] = act

					current_token = my.tokStream.GetToken()
					current_kind = my.tokStream.GetKind(current_token)
					my.locationStack[my.stateStackTop] = current_token
				} else {
					current_kind = candidate.symbol
					my.locationStack[my.stateStackTop] = candidate.location
				}
			}
			//
			// At my stage, we have a recovery configuration. See how
			// far we can go with it.
			//
			var next_token int = my.tokStream.Peek()
			my.tempStackTop = my.stateStackTop
			Arraycopy(my.stateStack, 0, my.tempStack, 0, my.stateStackTop+1)
			error_token = my.ParseForError(current_kind)

			if error_token != 0 {
				my.tokStream.ResetTo(next_token)
				my.tempStackTop = my.stateStackTop
				Arraycopy(my.stateStack, 0, my.tempStack, 0, my.stateStackTop+1)
				my.ParseUpToError(action, current_kind, error_token)
				my.tokStream.ResetTo(next_token)
			} else {
				act = my.ACCEPT_ACTION
			}
		}
	}
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
func (my *DiagnoseParser) ParseForError(current_kind int) int {
	var error_token int = 0
	//
	// Get next token in stream and compute initial action
	//
	var curtok int = my.tokStream.GetPrevious(my.tokStream.Peek())
	var act int = my.TAction(my.tempStack[my.tempStackTop], current_kind)
	//
	// Allocate configuration stack.
	//
	var configuration_stack = NewConfigurationStack(my.prs)

	//
	// Keep parsing until we reach the end of file and succeed or
	// an error is encountered. The list of actions executed will
	// be store in the "action" tuple.
	//
	for {
		if act <= my.NUM_RULES {

			my.tempStackTop--

			for {
				my.tempStackTop -= my.Rhs(act) - 1
				act = my.NtAction(my.tempStack[my.tempStackTop], my.Lhs(act))
				if act <= my.NUM_RULES {
					continue
				} else {
					break
				}
			}

		} else {
			if act > my.ERROR_ACTION {
				curtok = my.tokStream.GetToken()
				current_kind = my.tokStream.GetKind(curtok)
				act -= my.ERROR_ACTION

				for {
					my.tempStackTop -= (my.Rhs(act) - 1)
					act = my.NtAction(my.tempStack[my.tempStackTop], my.Lhs(act))
					if act <= my.NUM_RULES {
						continue
					} else {
						break
					}
				}

			} else {
				if act < my.ACCEPT_ACTION {
					curtok = my.tokStream.GetToken()
					current_kind = my.tokStream.GetKind(curtok)
				} else {
					if act == my.ERROR_ACTION {

						if !(error_token > curtok) {
							error_token = curtok
						}

						var configuration = configuration_stack.Pop()
						if configuration == nil {
							act = my.ERROR_ACTION
						} else {
							my.tempStackTop = configuration.stack_top
							configuration.RetrieveStack(my.tempStack)
							act = configuration.act
							curtok = configuration.curtok
							// no need to execute: action.Reset(configuration.action_length)
							current_kind = my.tokStream.GetKind(curtok)
							my.tokStream.ResetTo(my.tokStream.GetNext(curtok))
							continue
						}
						break
					} else {
						if act > my.ACCEPT_ACTION {
							if configuration_stack.FindConfiguration(my.tempStack, my.tempStackTop, curtok) {
								act = my.ERROR_ACTION
							} else {
								configuration_stack.Push(my.tempStack, my.tempStackTop, act+1, curtok, 0)
								act = my.BaseAction(act)
							}
							continue
						} else {
							break // assert(act == ACCEPT_ACTION)
						}
					}
				}
			}
		}

		my.tempStackTop += 1
		if my.tempStackTop >= len(my.tempStack) {
			my.ReallocateStacks()
		}
		my.tempStack[my.tempStackTop] = act
		act = my.TAction(act, current_kind)
	}
	if act == my.ERROR_ACTION {
		return error_token
	} else {
		return 0
	}
}

//
// Given the configuration consisting of the states in tempStack
// and the sequence of tokens (current_kind, followed by the tokens
// in tokStream), parse up to error_token in the tokStream and store
// all the parsing actions executed in the "action" tuple.
//
func (my *DiagnoseParser) ParseUpToError(action *IntTuple, current_kind int, error_token int) {
	//
	// Assume predecessor of next token and compute initial action
	//
	var curtok int = my.tokStream.GetPrevious(my.tokStream.Peek())
	var act int = my.TAction(my.tempStack[my.tempStackTop], current_kind)
	//
	// Allocate configuration stack.
	//
	var configuration_stack = NewConfigurationStack(my.prs)
	//
	// Keep parsing until we reach the end of file and succeed or
	// an error is encountered. The list of actions executed will
	// be store in the "action" tuple.
	//
	action.ReSet()
	for {
		if act <= my.NUM_RULES {
			action.Add(act) // save my reduce action
			my.tempStackTop--

			for {
				my.tempStackTop -= (my.Rhs(act) - 1)
				act = my.NtAction(my.tempStack[my.tempStackTop], my.Lhs(act))
				if act <= my.NUM_RULES {
					continue
				} else {
					break
				}
			}

		} else {
			if act > my.ERROR_ACTION {
				action.Add(act) // save my shift-reduce action
				curtok = my.tokStream.GetToken()
				current_kind = my.tokStream.GetKind(curtok)
				act -= my.ERROR_ACTION

				for {
					my.tempStackTop -= (my.Rhs(act) - 1)
					act = my.NtAction(my.tempStack[my.tempStackTop], my.Lhs(act))
					if act <= my.NUM_RULES {
						continue
					} else {
						break
					}
				}
			} else {
				if act < my.ACCEPT_ACTION {
					action.Add(act) // save my shift action
					curtok = my.tokStream.GetToken()
					current_kind = my.tokStream.GetKind(curtok)
				} else {
					if act == my.ERROR_ACTION {
						if curtok != error_token {
							var configuration = configuration_stack.Pop()
							if configuration == nil {
								act = my.ERROR_ACTION
							} else {
								my.tempStackTop = configuration.stack_top
								configuration.RetrieveStack(my.tempStack)
								act = configuration.act
								curtok = configuration.curtok
								action.ReSetTo(configuration.action_length)
								current_kind = my.tokStream.GetKind(curtok)
								my.tokStream.ResetTo(my.tokStream.GetNext(curtok))
								continue
							}
						}
						break
					} else {
						if act > my.ACCEPT_ACTION {
							if configuration_stack.FindConfiguration(my.tempStack, my.tempStackTop, curtok) {
								act = my.ERROR_ACTION
							} else {
								configuration_stack.Push(my.tempStack, my.tempStackTop, act+1, curtok, action.Size())
								act = my.BaseAction(act)
							}
							continue
						} else {
							break // assert(act == ACCEPT_ACTION)
						}
					}
				}
			}
		}

		my.tempStackTop += 1
		if my.tempStackTop >= len(my.tempStack) {
			my.ReallocateStacks()
		}
		my.tempStack[my.tempStackTop] = act
		act = my.TAction(act, current_kind)
	}
	action.Add(my.ERROR_ACTION)
	return
}

//
// Try to parse until first_symbol and all tokens in BUFFER have
// been consumed, or an error is encountered. Return the number
// of tokens that were expended before the parse blocked.
//
func (my *DiagnoseParser) ParseCheck(stack []int, stack_top int, first_symbol int, buffer_position int) int {
	var buffer_index int
	var current_kind int

	var local_stack []int = make([]int, len(stack))
	var local_stack_top int = stack_top
	var i int = 0
	for ; i <= stack_top; i++ {
		local_stack[i] = stack[i]
	}
	var configuration_stack = NewConfigurationStack(my.prs)

	//
	// If the first symbol is a nonterminal, process it here.
	//
	var act int = local_stack[local_stack_top]
	if first_symbol > my.NT_OFFSET {
		var lhs_symbol int = first_symbol - my.NT_OFFSET
		buffer_index = buffer_position
		current_kind = my.tokStream.GetKind(my.buffer[buffer_index])
		my.tokStream.ResetTo(my.tokStream.GetNext(my.buffer[buffer_index]))
		act = my.NtAction(act, lhs_symbol)
		for act <= my.NUM_RULES {
			local_stack_top -= (my.Rhs(act) - 1)
			act = my.NtAction(local_stack[local_stack_top], my.Lhs(act))
		}
	} else {
		local_stack_top--
		buffer_index = buffer_position - 1
		current_kind = first_symbol
		my.tokStream.ResetTo(my.buffer[buffer_position])
	}

	//
	// Start parsing the remaining symbols in the buffer
	//
	local_stack_top += 1
	if local_stack_top >= len(local_stack) { // Stack overflow!!!
		return buffer_index
	}
	local_stack[local_stack_top] = act

	act = my.TAction(act, current_kind)

	for {
		if act <= my.NUM_RULES { // reduce action

			local_stack_top -= my.Rhs(act)
			act = my.NtAction(local_stack[local_stack_top], my.Lhs(act))
			for act <= my.NUM_RULES {
				local_stack_top -= (my.Rhs(act) - 1)
				act = my.NtAction(local_stack[local_stack_top], my.Lhs(act))
			}
		} else {
			if act > my.ERROR_ACTION { // shift-reduce action

				if buffer_index == MAX_DISTANCE {
					buffer_index++
					break
				}
				buffer_index++

				current_kind = my.tokStream.GetKind(my.buffer[buffer_index])
				my.tokStream.ResetTo(my.tokStream.GetNext(my.buffer[buffer_index]))
				act -= my.ERROR_ACTION

				for {
					local_stack_top -= (my.Rhs(act) - 1)
					act = my.NtAction(local_stack[local_stack_top], my.Lhs(act))
					if act <= my.NUM_RULES {
						continue
					} else {
						break
					}
				}
			} else {
				if act < my.ACCEPT_ACTION { // shift action

					if buffer_index == MAX_DISTANCE {
						buffer_index++
						break
					}
					buffer_index++
					current_kind = my.tokStream.GetKind(my.buffer[buffer_index])
					my.tokStream.ResetTo(my.tokStream.GetNext(my.buffer[buffer_index]))
				} else {
					if act == my.ERROR_ACTION {
						var configuration = configuration_stack.Pop()
						if configuration == nil {
							act = my.ERROR_ACTION
						} else {
							local_stack_top = configuration.stack_top
							configuration.RetrieveStack(local_stack)
							act = configuration.act
							buffer_index = configuration.curtok
							// no need to execute: action.Reset(configuration.action_length)
							current_kind = my.tokStream.GetKind(my.buffer[buffer_index])
							my.tokStream.ResetTo(my.tokStream.GetNext(my.buffer[buffer_index]))
							continue
						}
						break
					} else {
						if act > my.ACCEPT_ACTION {
							if configuration_stack.FindConfiguration(local_stack, local_stack_top, buffer_index) {
								act = my.ERROR_ACTION
							} else {
								configuration_stack.Push(local_stack, local_stack_top, act+1, buffer_index, 0)
								act = my.BaseAction(act)
							}
							continue
						} else {
							break
						}
					}
				}
			}
		}

		local_stack_top += 1
		if local_stack_top >= len(local_stack) {
			break
		}
		local_stack[local_stack_top] = act
		act = my.TAction(act, current_kind)
	}
	if act == my.ACCEPT_ACTION {
		return MAX_DISTANCE
	} else {
		return buffer_index
	}
}

//
//  This routine is invoked when an error is encountered.  It
// tries to Diagnose the error and recover from it.  If it is
// successful, the state stack, the current token and the buffer
// are readjusted i.e., after a successful recovery,
// state_stack_top points to the location in the state stack
// that contains the state on which to recover current_token
// identifies the symbol on which to recover.
//
// Up to three configurations may be available when my routine
// is invoked. PREV_STACK may contain the sequence of states
// preceding any action on prevtok, STACK always contains the
// sequence of states preceding any action on current_token, and
// NEXT_STACK may contain the sequence of states preceding any
// action on the successor of current_token.
//
func (my *DiagnoseParser) ErrorRecovery(error_token int) *RepairCandidate {
	var prevtok int = my.tokStream.GetPrevious(error_token)

	//
	// Try primary phase recoveries. If not successful, try secondary
	// phase recoveries.  If not successful and we are at end of the
	// file, we issue the end-of-file error and quit. Otherwise, ...
	//
	var candidate *RepairCandidate = my.PrimaryPhase(error_token)
	if candidate.symbol != 0 {
		return candidate
	}
	candidate = my.SecondaryPhase(error_token)
	if candidate.symbol != 0 {
		return candidate
	}
	//
	// At my point, primary and (initial attempt at) secondary
	// recovery did not work.  We will now Get into "panic mode" and
	// keep trying secondary phase recoveries until we either find
	// a successful recovery or have consumed the remaining input
	// tokens.
	//
	if my.tokStream.GetKind(error_token) != my.EOFT_SYMBOL {
		for ;my.tokStream.GetKind(my.buffer[BUFF_UBOUND]) != my.EOFT_SYMBOL; {
			candidate = my.SecondaryPhase(my.buffer[MAX_DISTANCE-MIN_DISTANCE+2])
			if candidate.symbol != 0 {
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
	var scope_repair = NewPrimaryRepairInfo()
	scope_repair.bufferPosition = BUFF_UBOUND
	var top int
	for top = my.stateStackTop; top >= 0; top-- {
		my.scopeTrial(scope_repair, my.stateStack, top)
		if scope_repair.distance > 0 {
			break
		}
	}
	//
	// If any scope repair was successful, emit the message now
	//
	var i int = 0
	for ; i < my.scopeStackTop; i++ {
		my.EmitError(SCOPE_CODE,
			-my.scopeIndex[i],
			my.locationStack[my.scopePosition[i]],
			my.buffer[1],
			my.NonterminalIndex(my.ScopeLhs(my.scopeIndex[i])))
	}

	//
	// If the original error_token was already pointing to the EOF, issue the EOF-reached message.
	//
	if my.tokStream.GetKind(error_token) == my.EOFT_SYMBOL {
		my.EmitError(EOF_CODE,
			my.TerminalIndex(my.EOFT_SYMBOL),
			prevtok,
			prevtok, 0)
	} else {
		//
		// We reached the end of the file while panicking. Delete all
		// remaining tokens in the input.
		//
		var i int = BUFF_UBOUND
		for ; my.tokStream.GetKind(my.buffer[i]) == my.EOFT_SYMBOL; i-- {
		}

		my.EmitError(DELETION_CODE,
			my.TerminalIndex(my.tokStream.GetKind(error_token)),
			error_token,
			my.buffer[i], 0)
	}
	//
	// Create the "failed" candidate and return it.
	//
	candidate.symbol = 0
	candidate.location = my.buffer[BUFF_UBOUND] // point to EOF
	return candidate
}

//
// This function tries primary and scope recovery on each
// available configuration.  If a successful recovery is found
// and no secondary phase recovery can do better, a diagnosis is
// issued, the configuration is updated and the function returns
// "true".  Otherwise, it returns "false".
//
func (my *DiagnoseParser) PrimaryPhase(error_token int) *RepairCandidate {
	//
	// Initialize the buffer.
	//
	var i int
	if my.nextStackTop >= 0 {
		i = 3
	} else {
		i = 2
	}
	my.buffer[i] = error_token
	var j int = i
	for ; j > 0; j-- {
		my.buffer[j-1] = my.tokStream.GetPrevious(my.buffer[j])
	}
	var k int = i + 1
	for ; k < BUFF_SIZE; k++ {
		my.buffer[k] = my.tokStream.GetNext(my.buffer[k-1])
	}

	//
	// If NEXT_STACK_TOP > 0 then the parse was successful on CURRENT_TOKEN
	// and the error was detected on the successor of CURRENT_TOKEN. In
	// that case, first check whether or not primary recovery is
	// possible on next_stack ...
	//
	var repair = NewPrimaryRepairInfo()
	if my.nextStackTop >= 0 {
		repair.bufferPosition = 3
		my.CheckPrimaryDistance(repair, my.nextStack, my.nextStackTop)
	}

	//
	// ... Try primary recovery on the current token and compare
	// the quality of my recovery to the one on the next token...
	//
	var base_repair = NewPrimaryRepairInfoAndClone(repair)

	base_repair.bufferPosition = 2
	my.CheckPrimaryDistance(base_repair, my.stateStack, my.stateStackTop)
	if base_repair.distance > repair.distance || base_repair.misspellIndex > repair.misspellIndex {
		repair = base_repair
	}

	//
	// Finally, if prev_stack_top >= 0 try primary recovery on
	// the prev_stack configuration and compare it to the best
	// recovery computed thus far.
	//
	if my.prevStackTop >= 0 {
		var prev_repair = NewPrimaryRepairInfoAndClone(repair)
		prev_repair.bufferPosition = 1
		my.CheckPrimaryDistance(prev_repair, my.prevStack, my.prevStackTop)
		if prev_repair.distance > repair.distance || prev_repair.misspellIndex > repair.misspellIndex {
			repair = prev_repair
		}
	}

	//
	// Before accepting the best primary phase recovery obtained,
	// ensure that we cannot do better with a similar secondary
	// phase recovery.
	//
	var candidate = NewRepairCandidate()
	if my.nextStackTop >= 0 { // next_stack available

		if my.SecondaryCheck(my.nextStack, my.nextStackTop, 3, repair.distance) {
			return candidate
		}
	} else {
		if my.SecondaryCheck(my.stateStack, my.stateStackTop, 2, repair.distance) {
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
	if repair.code == INVALID_CODE ||
		repair.code == DELETION_CODE ||
		repair.code == SUBSTITUTION_CODE ||
		repair.code == MERGE_CODE {
		repair.distance--
	}

	//
	// ... After adjustment, check if the most successful primary
	// recovery can be applied.  If not, continue with more radical
	// recoveries...
	//
	if repair.distance < MIN_DISTANCE {
		return candidate
	}

	//
	// When processing an insertion error, if the token preceeding
	// the error token is not available, we change the repair code
	// into a BEFORE_CODE to instruct the reporting routine that it
	// indicates that the repair symbol should be inserted before
	// the error token.
	//
	if repair.code == INSERTION_CODE {
		if my.tokStream.GetKind(my.buffer[repair.bufferPosition-1]) == 0 {
			repair.code = BEFORE_CODE
		}
	}

	//
	// Select the proper sequence of states on which to recover,
	// update stack accordingly and call diagnostic routine.
	//
	if repair.bufferPosition == 1 {
		my.stateStackTop = my.prevStackTop
		Arraycopy(my.prevStack, 0, my.stateStack, 0, my.stateStackTop+1)
	} else {
		if my.nextStackTop >= 0 && repair.bufferPosition >= 3 {
			my.stateStackTop = my.nextStackTop
			Arraycopy(my.nextStack, 0, my.stateStack, 0, my.stateStackTop+1)
			my.locationStack[my.stateStackTop] = my.buffer[3]
		}
	}
	return my.PrimaryDiagnosis(repair)
}

//
//  This function checks whether or not a given state has a
// candidate, whose string representaion is a merging of the two
// tokens at positions buffer_position and buffer_position+1 in
// the buffer.  If so, it returns the candidate in question
// otherwise it returns 0.
//
func (my *DiagnoseParser) MergeCandidate(state int, buffer_position int) int {
	var str string = my.tokStream.GetName(my.buffer[buffer_position]) + my.tokStream.GetName(my.buffer[buffer_position+1])
	var k int = my.Asi(state)
	for ; my.asr(k) != 0; k++ {
		var i int = my.TerminalIndex(my.asr(k))
		if len(str) == len(my.Name(i)) {
			if ToLower(str) == ToLower(my.Name(i)) {
				return my.asr(k)
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
// which it Sets based on the best possible recovery that it
// finds in the given configuration.  The effectiveness of a
// a repair is judged based on two criteria:
//
//       1) the number of tokens that can be parsed after the repair
//              is applied: distance.
//       2) how close to perfection is the candidate that is chosen:
//              misspell_index.
//
// When my procedure is entered, distance, misspell_index and
// repair_code are assumed to be initialized.
//

func (my *DiagnoseParser) CheckPrimaryDistance(repair *PrimaryRepairInfo, stck []int, stack_top int) {
	//
	//  First, try scope recovery.
	//
	var scope_repair = NewPrimaryRepairInfoAndClone(repair)
	my.scopeTrial(scope_repair, stck, stack_top)
	if scope_repair.distance > repair.distance {
		repair.copy(scope_repair)
	}

	//
	//  Next, try merging the error token with its successor.
	//
	var symbol int = my.MergeCandidate(stck[stack_top], repair.bufferPosition)
	if symbol != 0 {
		var j int = my.ParseCheck(stck, stack_top, symbol, repair.bufferPosition+2)
		if (j > repair.distance) || (j == repair.distance && repair.misspellIndex < 10) {
			repair.misspellIndex = 10
			repair.symbol = symbol
			repair.distance = j
			repair.code = MERGE_CODE
		}
	}

	//
	// Next, try deletion of the error token.
	//
	var j int = my.ParseCheck(stck,
		stack_top,
		my.tokStream.GetKind(my.buffer[repair.bufferPosition+1]),
		repair.bufferPosition+2)

	var k int
	if  my.tokStream.GetKind(my.buffer[repair.bufferPosition]) == my.EOLT_SYMBOL &&
		my.tokStream.AfterEol(my.buffer[repair.bufferPosition+1]) {
		k = 10
	} else {
		k = 0
	}

	if j > repair.distance || (j == repair.distance && k > repair.misspellIndex) {
		repair.misspellIndex = k
		repair.code = DELETION_CODE
		repair.distance = j
	}

	//
	// Update the error configuration by simulating all reduce and
	// goto actions induced by the error token. Then assign the top
	// most state of the new configuration to next_state.
	//
	var next_state int = stck[stack_top]
	var max_pos int = stack_top
	my.tempStackTop = stack_top - 1

	my.tokStream.ResetTo(my.buffer[repair.bufferPosition+1])
	var tok int = my.tokStream.GetKind(my.buffer[repair.bufferPosition])
	var act int = my.TAction(next_state, tok)
	for act <= my.NUM_RULES {
		for {
			var lhs_symbol int = my.Lhs(act)
			my.tempStackTop -= (my.Rhs(act) - 1)

			if my.tempStackTop > max_pos {
				act = my.tempStack[my.tempStackTop]
			} else {
				act = stck[my.tempStackTop]
			}

			act = my.NtAction(act, lhs_symbol)
			if act <= my.NUM_RULES {
				continue
			} else {
				break
			}
		}
		if !(max_pos < my.tempStackTop) {
			max_pos = my.tempStackTop
		}
		my.tempStack[my.tempStackTop+1] = act
		next_state = act
		act = my.TAction(next_state, tok)
	}

	//
	//  Next, place the list of candidates in proper order.
	//
	var root int = 0
	var i int = my.Asi(next_state)
	for ; my.asr(i) != 0; i++ {
		symbol = my.asr(i)
		if symbol != my.EOFT_SYMBOL && symbol != my.ERROR_SYMBOL {
			if root == 0 {
				my.list[symbol] = symbol
			} else {
				my.list[symbol] = my.list[root]
				my.list[root] = symbol
			}
			root = symbol
		}
	}
	if stck[stack_top] != next_state {
		var i int = my.Asi(stck[stack_top])
		for ; my.asr(i) != 0; i++ {
			symbol = my.asr(i)
			if symbol != my.EOFT_SYMBOL && symbol != my.ERROR_SYMBOL && my.list[symbol] == 0 {
				if root == 0 {
					my.list[symbol] = symbol
				} else {
					my.list[symbol] = my.list[root]
					my.list[root] = symbol
				}
				root = symbol
			}
		}
	}

	var head int = my.list[root]
	my.list[root] = 0
	root = head

	//
	//  Next, try insertion for each possible candidate available in
	// the current state, except EOFT and ERROR_SYMBOL.
	//

	symbol = root
	for symbol != 0 {
		var m int = my.ParseCheck(stck, stack_top, symbol, repair.bufferPosition)
		var n int
		if symbol == my.EOLT_SYMBOL && my.tokStream.AfterEol(my.buffer[repair.bufferPosition]) {
			n = 10
		} else {
			n = 0
		}

		if m > repair.distance ||
			(m == repair.distance && n > repair.misspellIndex) {
			repair.misspellIndex = n
			repair.distance = m
			repair.symbol = symbol
			repair.code = INSERTION_CODE
		}
		symbol = my.list[symbol]
	}

	//
	//  Next, Try substitution for each possible candidate available
	// in the current state, except EOFT and ERROR_SYMBOL.
	//
	symbol = root
	for symbol != 0 {

		var m int = my.ParseCheck(stck, stack_top, symbol, repair.bufferPosition+1)
		var n int
		if symbol == my.EOLT_SYMBOL && my.tokStream.AfterEol(my.buffer[repair.bufferPosition+1]) {
			n = 10
		} else {
			n = my.Misspell(symbol, my.buffer[repair.bufferPosition])
		}

		if m > repair.distance ||
			(m == repair.distance && n > repair.misspellIndex) {
			repair.misspellIndex = n
			repair.distance = m
			repair.symbol = symbol
			repair.code = SUBSTITUTION_CODE
		}
		var s int = symbol
		symbol = my.list[symbol]
		my.list[s] = 0 // Reset element
	}

	//
	// Next, we try to insert a nonterminal candidate in front of the
	// error token, or substituting a nonterminal candidate for the
	// error token. Precedence is given to insertion.
	//
	var nt_index int
	for nt_index = my.Nasi(stck[stack_top]); my.Nasr(nt_index) != 0; nt_index++ {
		symbol = my.Nasr(nt_index) + my.NT_OFFSET
		var n int = my.ParseCheck(stck, stack_top, symbol, repair.bufferPosition+1)
		if n > repair.distance {
			repair.misspellIndex = 0
			repair.distance = n
			repair.symbol = symbol
			repair.code = INVALID_CODE
		}

		n = my.ParseCheck(stck, stack_top, symbol, repair.bufferPosition)
		if n > repair.distance || (n == repair.distance && repair.code == INVALID_CODE) {
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
func (my *DiagnoseParser) PrimaryDiagnosis(repair *PrimaryRepairInfo) *RepairCandidate {
	//
	//  Issue diagnostic.
	//
	var prevtok int = my.buffer[repair.bufferPosition-1]
	var current_token int = my.buffer[repair.bufferPosition]

	switch repair.code {
	case INSERTION_CODE:
	case BEFORE_CODE:
		{
			var name_index int
			if repair.symbol > my.NT_OFFSET {
				name_index = my.GetNtermIndex(my.stateStack[my.stateStackTop],
													 repair.symbol,
													 repair.bufferPosition)
			} else {
				name_index = my.GetTermIndex( my.stateStack,
												my.stateStackTop,
												repair.symbol,
												repair.bufferPosition)
			}

			var tok int
			if repair.code == INSERTION_CODE {
				tok = prevtok
			} else {
				tok = current_token
			}
			my.EmitError(repair.code, name_index, tok, tok, 0)
		}
		break
	case INVALID_CODE:
		{
			var name_index int = my.GetNtermIndex(my.stateStack[my.stateStackTop],
														 repair.symbol,
														 repair.bufferPosition+1)
			my.EmitError(repair.code, name_index, current_token, current_token, 0)
		}
		break
	case SUBSTITUTION_CODE:
		{
			var name_index int
			if repair.misspellIndex >= 6 {
				name_index = my.TerminalIndex(repair.symbol)
			} else {
				name_index = my.GetTermIndex( my.stateStack,
												my.stateStackTop,
												repair.symbol,
												repair.bufferPosition+1)
				if name_index != my.TerminalIndex(repair.symbol) {
					repair.code = INVALID_CODE
				}
			}
			my.EmitError(repair.code, name_index, current_token, current_token, 0)
		}
		break
	case MERGE_CODE:
		my.EmitError(repair.code,
			my.TerminalIndex(repair.symbol),
			current_token,
			my.tokStream.GetNext(current_token), 0)
		break
	case SCOPE_CODE:
		{
			var i int = 0
			for ; i < my.scopeStackTop; i++ {
				my.EmitError(repair.code,
					-my.scopeIndex[i],
					my.locationStack[my.scopePosition[i]],
					prevtok,
					my.NonterminalIndex(my.ScopeLhs(my.scopeIndex[i])))
			}
			repair.symbol = my.ScopeLhs(my.scopeIndex[my.scopeStackTop]) + my.NT_OFFSET
			my.stateStackTop = my.scopePosition[my.scopeStackTop]
			my.EmitError(repair.code,
				-my.scopeIndex[my.scopeStackTop],
				my.locationStack[my.scopePosition[my.scopeStackTop]],
				prevtok,
				my.GetNtermIndex(my.stateStack[my.stateStackTop],
														repair.symbol,
														repair.bufferPosition))
			break
		}
	default: // deletion
		my.EmitError(repair.code, my.TerminalIndex(my.ERROR_SYMBOL), current_token, current_token, 0)
		break
	}

	//
	//  Update buffer.
	//
	var candidate = NewRepairCandidate()
	switch repair.code {
	case INSERTION_CODE:
	case BEFORE_CODE:
	case SCOPE_CODE:
		candidate.symbol = repair.symbol
		candidate.location = my.buffer[repair.bufferPosition]
		my.tokStream.ResetTo(my.buffer[repair.bufferPosition])
		break
	case INVALID_CODE:
	case SUBSTITUTION_CODE:
		candidate.symbol = repair.symbol
		candidate.location = my.buffer[repair.bufferPosition]
		my.tokStream.ResetTo(my.buffer[repair.bufferPosition+1])
		break
	case MERGE_CODE:
		candidate.symbol = repair.symbol
		candidate.location = my.buffer[repair.bufferPosition]
		my.tokStream.ResetTo(my.buffer[repair.bufferPosition+2])
		break
	default: // deletion
		candidate.location = my.buffer[repair.bufferPosition+1]
		candidate.symbol = my.tokStream.GetKind(my.buffer[repair.bufferPosition+1])
		my.tokStream.ResetTo(my.buffer[repair.bufferPosition+2])
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
// the Name index of the highest level nonterminal that can
// directly or indirectly produce the candidate.
//
func (my *DiagnoseParser) GetTermIndex(stck []int, stack_top int, tok int, buffer_position int) int {
	//
	// Initialize stack index of temp_stack and initialize maximum
	// position of state stack that is still useful.
	//
	var act int = stck[stack_top]
	var max_pos int = stack_top
	var highest_symbol int = tok

	my.tempStackTop = stack_top - 1

	//
	// Compute all reduce and associated actions induced by the
	// candidate until a SHIFT or SHIFT-REDUCE is computed. ERROR
	// and ACCEPT actions cannot be computed on the candidate in
	// my context, since we know that it is suitable for recovery.
	//
	my.tokStream.ResetTo(my.buffer[buffer_position])
	act = my.TAction(act, tok)
	for act <= my.NUM_RULES {
		//
		// Process all goto-reduce actions following reduction,
		// until a goto action is computed ...
		//
		for {
			var lhs_symbol int = my.Lhs(act)
			my.tempStackTop -= (my.Rhs(act) - 1)

			if my.tempStackTop > max_pos {
				act = my.tempStack[my.tempStackTop]
			} else {
				act = stck[my.tempStackTop]
			}

			act = my.NtAction(act, lhs_symbol)
			if act <= my.NUM_RULES {
				continue
			} else {
				break
			}
		}
		//
		// Compute new maximum useful position of (STATE_)stack,
		// push goto state into the stack, and compute next
		// action on candidate ...
		//

		if !(max_pos < my.tempStackTop) {
			max_pos = my.tempStackTop
		}
		my.tempStack[my.tempStackTop+1] = act
		act = my.TAction(act, tok)
	}

	//
	// At my stage, we have simulated all actions induced by the
	// candidate and we are ready to shift or shift-reduce it. First,
	// Set tok and next_ptr appropriately and identify the candidate
	// as the initial highest_symbol. If a shift action was computed
	// on the candidate, update the stack and compute the next
	// action. Next, simulate all actions possible on the next input
	// token until we either have to shift it or are about to reduce
	// below the initial starting point in the stack (indicated by
	// max_pos as computed in the previous loop).  At that point,
	// return the highest_symbol computed.
	//
	my.tempStackTop++ // adjust top of stack to reflect last goto
	// next move is shift or shift-reduce.

	var threshold int = my.tempStackTop

	tok = my.tokStream.GetKind(my.buffer[buffer_position])
	my.tokStream.ResetTo(my.buffer[buffer_position+1])

	if act > my.ERROR_ACTION { // shift-reduce on candidate?
		act -= my.ERROR_ACTION
	} else {
		if act < my.ACCEPT_ACTION { // shift on candidate
			my.tempStack[my.tempStackTop+1] = act
			act = my.TAction(act, tok)
		}
	}
	for act <= my.NUM_RULES {
		//
		// Process all goto-reduce actions following reduction,
		// until a goto action is computed ...
		//
		for {
			var lhs_symbol int = my.Lhs(act)
			my.tempStackTop -= (my.Rhs(act) - 1)

			if my.tempStackTop < threshold {

				if highest_symbol > my.NT_OFFSET {
					return my.NonterminalIndex(highest_symbol - my.NT_OFFSET)
				} else {
					my.TerminalIndex(highest_symbol)
				}
			}
			if my.tempStackTop == threshold {
				highest_symbol = lhs_symbol + my.NT_OFFSET
			}
			if my.tempStackTop > max_pos {
				act = my.tempStack[my.tempStackTop]
			} else {
				act = stck[my.tempStackTop]
			}

			act = my.NtAction(act, lhs_symbol)
			if act <= my.NUM_RULES {
				continue
			} else {
				break
			}
		}

		my.tempStack[my.tempStackTop+1] = act
		act = my.TAction(act, tok)

	}
	if highest_symbol > my.NT_OFFSET {
		return my.NonterminalIndex(highest_symbol - my.NT_OFFSET)
	} else {
		return my.TerminalIndex(highest_symbol)
	}

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
func (my *DiagnoseParser) GetNtermIndex(start int, sym int, buffer_position int) int {
	var highest_symbol int = sym - my.NT_OFFSET
	var tok int = my.tokStream.GetKind(my.buffer[buffer_position])
	my.tokStream.ResetTo(my.buffer[buffer_position+1])

	//
	// Initialize stack index of temp_stack and initialize maximum
	// position of state stack that is still useful.
	//
	my.tempStackTop = 0
	my.tempStack[my.tempStackTop] = start

	var act int = my.NtAction(start, highest_symbol)
	if act > my.NUM_RULES { // goto action?
		my.tempStack[my.tempStackTop+1] = act
		act = my.TAction(act, tok)
	}

	for act <= my.NUM_RULES {
		//
		// Process all goto-reduce actions following reduction,
		// until a goto action is computed ...
		//
		for {
			my.tempStackTop -= (my.Rhs(act) - 1)
			if my.tempStackTop < 0 {
				return my.NonterminalIndex(highest_symbol)
			}
			if my.tempStackTop == 0 {
				highest_symbol = my.Lhs(act)
			}
			act = my.NtAction(my.tempStack[my.tempStackTop], my.Lhs(act))
			if act <= my.NUM_RULES {
				continue
			} else {
				break
			}
		}
		my.tempStack[my.tempStackTop+1] = act
		act = my.TAction(act, tok)
	}
	return my.NonterminalIndex(highest_symbol)
}

//
//  Check whether or not there is a high probability that a
// given string is a misspelling of another.
// Certain singleton symbols (such as ":" and "") are also
// considered to be misspellings of each other.
//
func (my *DiagnoseParser) Misspell(sym int, tok int) int {
	//
	// Set up the two strings in question. Note that there is a "0"
	// gate added at the end of each string. This is important as
	// the algorithm assumes that it can "Peek" at the symbol immediately
	// following the one that is being analysed.
	//
	var s1 string = ToLower(my.Name(my.TerminalIndex(sym)))
	var n int = len(s1)
	s1 = AppendRune(s1, '\u0000')

	var s2 string = ToLower(my.tokStream.GetName(tok))
	var m int
	if len(s2) < my.MAX_NAME_LENGTH {
		m = len(s2)
	} else {
		m = my.MAX_NAME_LENGTH
	}
	s2 = SubStr(s2, 0, m)
	s2 = AppendRune(s2, '\u0000')
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
	if n == 1 && m == 1 {
		if (CharAt(s1, 0) == ';' && CharAt(s2, 0) == ',') ||
			(CharAt(s1, 0) == ',' && CharAt(s2, 0) == ';') ||
			(CharAt(s1, 0) == ';' && CharAt(s2, 0) == ':') ||
			(CharAt(s1, 0) == ':' && CharAt(s2, 0) == ';') ||
			(CharAt(s1, 0) == '.' && CharAt(s2, 0) == ',') ||
			(CharAt(s1, 0) == ',' && CharAt(s2, 0) == '.') ||
			(CharAt(s1, 0) == '\'' && CharAt(s2, 0) == '"') ||
			(CharAt(s1, 0) == '"' && CharAt(s2, 0) == '\'') {
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
	var count int = 0
	var prefix_length int = 0
	var num_errors int = 0

	var i int = 0
	var j int = 0

	for (i < n) && (j < m) {
		if CharAt(s1, i) == CharAt(s2, j) {
			count++
			i++
			j++
			if num_errors == 0 {
				prefix_length++
			}
		} else {
			if CharAt(s1, i+1) == CharAt(s2, j) && CharAt(s1, i) == CharAt(s2, j+1) { //transposition

				count += 2
				i += 2
				j += 2
				num_errors++
			} else {
				if CharAt(s1, i+1) == CharAt(s2, j+1) { // mismatch
					i += 2
					j += 2
					num_errors++
				} else {
					if (n - i) > (m - j) {
						i++
					} else {
						if (m - j) > (n - i) {
							j++
						} else {
							i++
							j++
						}
					}
					num_errors++
				}
			}
		}
	}

	if i < n || j < m {
		num_errors++
	}
	var temp int = 1
	if n < m {
		temp = n/6 + 1
	} else {
		temp = m/6 + 1
	}
	if num_errors > temp {
		count = prefix_length
	}
	if n < len(s1) {
		temp = len(s1)
	} else {
		temp = n
	}
	return count * 10 / (temp + num_errors)
}


func (my *DiagnoseParser) ScopeTrialCheck(repair *PrimaryRepairInfo, stack []int, stack_top int, indx int) {

	var i int = my.stateSeen[stack_top]
	for ; i != NIL; i = my.statePool[i].next {
		if my.statePool[i].state == stack[stack_top] {
			return
		}
	}
	var old_state_pool_top int = my.statePoolTop
	my.statePoolTop++
	if my.statePoolTop >= len(my.statePool) {
		my.statePool = StateInfoArraycopy(my.statePool, 0, make([]*StateInfo, my.statePoolTop*2), 0, my.statePoolTop)
	}

	my.statePool[old_state_pool_top] = NewStateInfo(stack[stack_top], my.stateSeen[stack_top])
	my.stateSeen[stack_top] = old_state_pool_top

	var action = NewIntTupleWithEstimate(1 << 3)
	i = 0
	for ; i < my.SCOPE_SIZE; i++ {
		//
		// Compute the action (or Set of actions in case of conflicts) that
		// can be executed on the scope Lookahead symbol. Save the action(s)
		// in the action tuple.
		//
		action.ReSet()
		var act int = my.TAction(stack[stack_top], my.ScopeLa(i))
		if act > my.ACCEPT_ACTION && act < my.ERROR_ACTION {
			// conflicting actions?
			for {

				action.Add(my.BaseAction(act))
				act++
				if my.BaseAction(act) != 0 {
					continue
				} else {
					break
				}
			}
		} else {
			action.Add(act)
		}

		//
		// For each action defined on the scope Lookahead symbol,
		// try scope recovery.
		//
		var action_index int = 0
		for ; action_index < action.Size(); action_index++ {
			my.tokStream.ResetTo(my.buffer[repair.bufferPosition])
			my.tempStackTop = stack_top - 1
			var max_pos int = stack_top

			act = action.Get(action_index)
			for act <= my.NUM_RULES {
				//
				// ... Process all goto-reduce actions following
				// reduction, until a goto action is computed ...
				//
				for {
					var lhs_symbol int = my.Lhs(act)
					my.tempStackTop -= (my.Rhs(act) - 1)

					if my.tempStackTop > max_pos {
						act = my.tempStack[my.tempStackTop]
					} else {
						act = stack[my.tempStackTop]
					}

					act = my.NtAction(act, lhs_symbol)
					if act <= my.NUM_RULES {
						continue
					} else {
						break
					}
				}
				if my.tempStackTop+1 >= len(my.stateStack) {
					return
				}
				if !(max_pos < my.tempStackTop) {
					max_pos = my.tempStackTop
				}
				my.tempStack[my.tempStackTop+1] = act
				act = my.TAction(act, my.ScopeLa(i))
			}
			//
			// If the Lookahead symbol is parsable, then we check
			// whether or not we have a match between the scope
			// prefix and the transition symbols corresponding to
			// the states on top of the stack.
			//
			if act != my.ERROR_ACTION {
				var j int
				var k int = my.ScopePrefix(i)
				for j = my.tempStackTop + 1; j >= (max_pos+1) &&
					my.InSymbol(my.tempStack[j]) == my.ScopeRhs(k); j-- {
					k++
				}
				if j == max_pos {
					for j = max_pos; j >= 1 && my.InSymbol(stack[j]) == my.ScopeRhs(k); j-- {
						k++
					}
				}
				//
				// If the prefix matches, check whether the state
				// newly exposed on top of the stack, (after the
				// corresponding prefix states are popped from the
				// stack), is in the Set of "source states" for the
				// scope in question and that it is at a position
				// below the threshold indicated by MARKED_POS.
				//
				var marked_pos int
				if max_pos < stack_top {
					marked_pos = max_pos + 1
				} else {
					marked_pos = stack_top
				}

				if my.ScopeRhs(k) == 0 && j < marked_pos { // match?
					var stack_position int = j
					for j = my.ScopeStateSet(i); stack[stack_position] != my.ScopeState(j) &&
						my.ScopeState(j) != 0; j++ {
					}
					//
					// If the top state is valid for scope recovery,
					// the left-hand side of the scope is used as
					// starting symbol and we calculate how far the
					// parser can advance within the forward context
					// after parsing the left-hand symbol.
					//
					if my.ScopeState(j) != 0 { // state was found
						var previous_distance int = repair.distance
						var distance int = my.ParseCheck( stack,
															stack_position,
															my.ScopeLhs(i)+my.NT_OFFSET,
															repair.bufferPosition)

						//
						// if the recovery is not successful, we
						// update the stack with all actions induced
						// by the left-hand symbol, and recursively
						// call SCOPE_TRIAL_CHECK to try again.
						// Otherwise, the recovery is successful. If
						// the new distance is greater than the
						// initial SCOPE_DISTANCE, we update
						// SCOPE_DISTANCE and Set scope_stack_top to INDX
						// to indicate the number of scopes that are
						// to be applied for a succesful  recovery.
						// NOTE that my procedure cannot Get into
						// an infinite loop, since each prefix match
						// is guaranteed to take us to a lower point
						// within the stack.
						//
						if (distance - repair.bufferPosition + 1) < MIN_DISTANCE {
							var top int = stack_position
							act = my.NtAction(stack[top], my.ScopeLhs(i))
							for act <= my.NUM_RULES {
								top -= (my.Rhs(act) - 1)
								act = my.NtAction(stack[top], my.Lhs(act))
							}
							top++
							j = act
							act = stack[top] // save
							stack[top] = j   // swap
							my.ScopeTrialCheck(repair, stack, top, indx+1)
							stack[top] = act // restore
						} else {
							if distance > repair.distance {
								my.scopeStackTop = indx
								repair.distance = distance
							}
						}
						//
						// If no other recovery possibility is left (due to
						// backtracking and we are at the end of the input,
						// then we favor a scope recovery over all other kinds
						// of recovery.
						//
						if my.tokStream.GetKind(my.buffer[repair.bufferPosition]) == my.EOFT_SYMBOL &&
							repair.distance == previous_distance {
							my.scopeStackTop = indx
							repair.distance = MAX_DISTANCE
						}
						//
						// If my scope recovery has beaten the
						// previous distance, then we have found a
						// better recovery (or my recovery is one
						// of a list of scope recoveries). Record
						// its information at the proper location
						// (INDX) in SCOPE_INDEX and SCOPE_STACK.
						//
						if repair.distance > previous_distance {
							my.scopeIndex[indx] = i
							my.scopePosition[indx] = stack_position
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
func (my *DiagnoseParser) SecondaryCheck(stack []int, stack_top int, buffer_position int, distance int) bool {
	var top int
	for top = stack_top - 1; top >= 0; top-- {
		var j int = my.ParseCheck(stack,
									top,
									my.tokStream.GetKind(my.buffer[buffer_position]),
									buffer_position+1)
		if ((j - buffer_position + 1) > MIN_DISTANCE) && (j > distance) {
			return true
		}
	}

	var scope_repair *PrimaryRepairInfo = NewPrimaryRepairInfo()
	scope_repair.bufferPosition = buffer_position + 1
	scope_repair.distance = distance
	my.scopeTrial(scope_repair, stack, stack_top)
	return (scope_repair.distance-buffer_position) > MIN_DISTANCE && scope_repair.distance > distance
}


func (my *DiagnoseParser) scopeTrial(repair *PrimaryRepairInfo, stack []int, stack_top int) {
	if len(my.stateSeen) == 0 || len(my.stateSeen) < len(my.stateStack) {
		my.stateSeen = make([]int, len(my.stateStack))
	}
	var i int = 0
	for ; i < len(my.stateStack); i++ {
		my.stateSeen[i] = NIL
	}

	my.statePoolTop = 0
	if len(my.statePool) == 0 || len(my.statePool) < len(my.stateStack) {
		my.statePool = make([]*StateInfo, len(my.stateStack))
	}
	my.ScopeTrialCheck(repair, stack, stack_top, 0)
	repair.code = SCOPE_CODE
	repair.misspellIndex = 10
	return
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
func (my *DiagnoseParser) SecondaryPhase(error_token int) *RepairCandidate {
	var repair = NewSecondaryRepairInfo()
	var misplaced_repair = NewSecondaryRepairInfo()

	//
	// If the next_stack is available, try misplaced and secondary
	// recovery on it first.
	//
	var next_last_index int = 0
	if my.nextStackTop >= 0 {

		var save_location int

		my.buffer[2] = error_token
		my.buffer[1] = my.tokStream.GetPrevious(my.buffer[2])
		my.buffer[0] = my.tokStream.GetPrevious(my.buffer[1])
		var k int = 3
		for ; k < BUFF_UBOUND; k++ {
			my.buffer[k] = my.tokStream.GetNext(my.buffer[k-1])
		}

		my.buffer[BUFF_UBOUND] = my.tokStream.BadToken() // elmt not available
		//
		// If we are at the end of the input stream, compute the
		// index position of the first EOFT symbol (last useful
		// index).
		//
		for next_last_index = MAX_DISTANCE - 1;
			next_last_index >= 1 &&
			my.tokStream.GetKind(my.buffer[next_last_index]) == my.EOFT_SYMBOL;
			next_last_index-- {
		}

		next_last_index = next_last_index + 1

		save_location = my.locationStack[my.nextStackTop]
		my.locationStack[my.nextStackTop] = my.buffer[2]
		misplaced_repair.numDeletions = my.nextStackTop
		my.MisplacementRecovery(misplaced_repair, my.nextStack, my.nextStackTop, next_last_index, true)
		if misplaced_repair.recoveryOnNextStack {
			misplaced_repair.distance++
		}
		repair.numDeletions = my.nextStackTop + BUFF_UBOUND
		my.SecondaryRecovery( repair,
								my.nextStack,
								my.nextStackTop,
								next_last_index, true)

		if repair.recoveryOnNextStack {
			repair.distance++
		}
		my.locationStack[my.nextStackTop] = save_location
	} else { // next_stack not available, initialize ...
		misplaced_repair.numDeletions = my.stateStackTop
		repair.numDeletions = my.stateStackTop + BUFF_UBOUND
	}

	//
	// Try secondary recovery on the "stack" configuration.
	//
	my.buffer[3] = error_token

	my.buffer[2] = my.tokStream.GetPrevious(my.buffer[3])
	my.buffer[1] = my.tokStream.GetPrevious(my.buffer[2])
	my.buffer[0] = my.tokStream.GetPrevious(my.buffer[1])
	var k int = 4
	for ; k < BUFF_SIZE; k++ {
		my.buffer[k] = my.tokStream.GetNext(my.buffer[k-1])
	}

	var last_index int
	for last_index = MAX_DISTANCE - 1;
		last_index >= 1 &&
		my.tokStream.GetKind(my.buffer[last_index]) == my.EOFT_SYMBOL;
		last_index-- {
	}
	last_index++

	my.MisplacementRecovery(misplaced_repair, my.stateStack, my.stateStackTop, last_index, false)

	my.SecondaryRecovery(repair, my.stateStack, my.stateStackTop, last_index, false)

	//
	// If a successful misplaced recovery was found, compare it with
	// the most successful secondary recovery.  If the misplaced
	// recovery either deletes fewer symbols or parse-checks further
	// then it is chosen.
	//
	if misplaced_repair.distance > MIN_DISTANCE {
		if misplaced_repair.numDeletions <= repair.numDeletions ||
			(misplaced_repair.distance-misplaced_repair.numDeletions) >=
				(repair.distance-repair.numDeletions) {
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
	if repair.recoveryOnNextStack {
		my.stateStackTop = my.nextStackTop
		Arraycopy(my.nextStack, 0, my.stateStack, 0, my.stateStackTop+1)

		my.buffer[2] = error_token
		my.buffer[1] = my.tokStream.GetPrevious(my.buffer[2])
		my.buffer[0] = my.tokStream.GetPrevious(my.buffer[1])
		var k int = 3
		for ; k < BUFF_UBOUND; k++ {
			my.buffer[k] = my.tokStream.GetNext(my.buffer[k-1])
		}

		my.buffer[BUFF_UBOUND] = my.tokStream.BadToken() // elmt not available

		my.locationStack[my.nextStackTop] = my.buffer[2]
		last_index = next_last_index
	}

	//
	// Next, try scope recoveries after deletion of one, two, three,
	// four ... buffer_position tokens from the input stream.
	//
	if repair.code == SECONDARY_CODE || repair.code == DELETION_CODE {
		var scope_repair *PrimaryRepairInfo = NewPrimaryRepairInfo()
		for scope_repair.bufferPosition = 2;
			scope_repair.bufferPosition <= repair.bufferPosition &&
			repair.code != SCOPE_CODE;
			scope_repair.bufferPosition++ {
			my.scopeTrial(scope_repair, my.stateStack, my.stateStackTop)
			var j int
			if scope_repair.distance == MAX_DISTANCE {
				j = last_index
			} else {
				j = scope_repair.distance
			}

			var k int = scope_repair.bufferPosition - 1
			if (scope_repair.distance-k) > MIN_DISTANCE && (j-k) > (repair.distance-repair.numDeletions) {
				var i int = my.scopeIndex[my.scopeStackTop] // upper bound
				repair.code = SCOPE_CODE
				repair.symbol = my.ScopeLhs(i) + my.NT_OFFSET
				repair.stackPosition = my.stateStackTop
				repair.bufferPosition = scope_repair.bufferPosition
			}
		}
	}
	//
	// If a successful repair was not found, quit!  Otherwise, issue
	// diagnosis and adjust configuration...
	//
	var candidate = NewRepairCandidate()
	if repair.code == 0 {
		return candidate
	}
	my.SecondaryDiagnosis(repair)

	//
	// Update buffer based on number of elements that are deleted.
	//
	switch repair.code {
	case MISPLACED_CODE:
		candidate.location = my.buffer[2]
		candidate.symbol = my.tokStream.GetKind(my.buffer[2])
		my.tokStream.ResetTo(my.tokStream.GetNext(my.buffer[2]))
		break
	case DELETION_CODE:
		candidate.location = my.buffer[repair.bufferPosition]
		candidate.symbol = my.tokStream.GetKind(my.buffer[repair.bufferPosition])
		my.tokStream.ResetTo(my.tokStream.GetNext(my.buffer[repair.bufferPosition]))
		break
	default: // SCOPE_CODE || SECONDARY_CODE
		candidate.symbol = repair.symbol
		candidate.location = my.buffer[repair.bufferPosition]
		my.tokStream.ResetTo(my.buffer[repair.bufferPosition])
		break
	}
	return candidate
}

//
// This bool function checks whether or not a given
// configuration yields a better misplacement recovery than
// the best misplacement recovery computed previously.
//
func (my *DiagnoseParser) MisplacementRecovery(repair *SecondaryRepairInfo, stack []int, stack_top int,
	last_index int, stack_flag bool) {
	var previous_loc int = my.buffer[2]
	var stack_deletions int = 0
	var top int
	for top = stack_top - 1; top >= 0; top-- {
		if my.locationStack[top] < previous_loc {
			stack_deletions++
		}
		previous_loc = my.locationStack[top]

		var parse_distance int = my.ParseCheck(stack, top, my.tokStream.GetKind(my.buffer[2]), 3)
		var j int
		if parse_distance == MAX_DISTANCE {
			j = last_index
		} else {
			j = parse_distance
		}
		if (parse_distance > MIN_DISTANCE) && (j-stack_deletions) > (repair.distance-repair.numDeletions) {
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
func (my *DiagnoseParser) SecondaryRecovery(repair *SecondaryRepairInfo, stack []int, stack_top int, last_index int, stack_flag bool) {
	var previous_loc int = my.buffer[2]
	var stack_deletions int = 0
	var top int = stack_top
	for ; top >= 0 && repair.numDeletions >= stack_deletions; top-- {
		if my.locationStack[top] < previous_loc {
			stack_deletions++
		}
		previous_loc = my.locationStack[top]
		var i int
		for i = 2; i <= (last_index-MIN_DISTANCE+1) &&
			(repair.numDeletions >= (stack_deletions + i - 1)); i++ {
			var parse_distance int = my.ParseCheck(stack, top, my.tokStream.GetKind(my.buffer[i]), i+1)
			var j int
			if parse_distance == MAX_DISTANCE {
				j = last_index
			} else {
				j = parse_distance
			}
			if (parse_distance - i + 1) > MIN_DISTANCE {
				var k int = stack_deletions + i - 1
				if (k < repair.numDeletions) ||
					(j-k) > (repair.distance-repair.numDeletions) ||
					((repair.code == SECONDARY_CODE) && (j-k) == (repair.distance-repair.numDeletions)) {
					repair.code = DELETION_CODE
					repair.distance = j
					repair.stackPosition = top
					repair.bufferPosition = i
					repair.numDeletions = k
					repair.recoveryOnNextStack = stack_flag
				}
			}
			var l int
			for l = my.Nasi(stack[top]); l >= 0 && my.Nasr(l) != 0; l++ {
				var symbol int = my.Nasr(l) + my.NT_OFFSET
				parse_distance = my.ParseCheck(stack, top, symbol, i)
				if parse_distance == MAX_DISTANCE {
					j = last_index
				} else {
					j = parse_distance
				}

				if (parse_distance - i + 1) > MIN_DISTANCE {
					var k int = stack_deletions + i - 1
					if k < repair.numDeletions || (j-k) > (repair.distance-repair.numDeletions) {
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
func (my *DiagnoseParser) SecondaryDiagnosis(repair *SecondaryRepairInfo) {
	switch repair.code {
	case SCOPE_CODE:
		if repair.stackPosition < my.stateStackTop {
			my.EmitError(DELETION_CODE,
				my.TerminalIndex(my.ERROR_SYMBOL),
				my.locationStack[repair.stackPosition],
				my.buffer[1], 0)
		}
		var i int = 0
		for ; i < my.scopeStackTop; i++ {
			my.EmitError(SCOPE_CODE,
				-my.scopeIndex[i],
				my.locationStack[my.scopePosition[i]],
				my.buffer[1],
				my.NonterminalIndex(my.ScopeLhs(my.scopeIndex[i])))
		}

		repair.symbol = my.ScopeLhs(my.scopeIndex[my.scopeStackTop]) + my.NT_OFFSET
		my.stateStackTop = my.scopePosition[my.scopeStackTop]
		my.EmitError(SCOPE_CODE,
			-my.scopeIndex[my.scopeStackTop],
			my.locationStack[my.scopePosition[my.scopeStackTop]],
			my.buffer[1],
			my.GetNtermIndex(my.stateStack[my.stateStackTop],
				repair.symbol,
				repair.bufferPosition))
		break
	default:
		var name_index int
		if repair.code == SECONDARY_CODE {
			name_index = my.GetNtermIndex(my.stateStack[repair.stackPosition],
												repair.symbol,
												repair.bufferPosition)
		} else {
			name_index = my.TerminalIndex(my.ERROR_SYMBOL)
		}
		my.EmitError(repair.code, name_index,
			my.locationStack[repair.stackPosition],
			my.buffer[repair.bufferPosition-1], 0)
		my.stateStackTop = repair.stackPosition
	}
	return
}

//
// keep looking ahead until we compute a valid action
//
func (my *DiagnoseParser) Lookahead(act int, token int) int {
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
func (my *DiagnoseParser) TAction(act int, sym int) int {
	act = my.prs.TAction(act, sym)
	if act > my.LA_STATE_OFFSET {
		return my.Lookahead(act, my.tokStream.Peek())
	} else {
		return act
	}

}

//
// This method is invoked by an LPG PARSER or a semantic
// routine to process an error message.
//

func (my *DiagnoseParser) EmitError(msg_code int, name_index int, left_token int, right_token int, scope_name_index int) {
	/*
	   var left_token_loc int
	   if left_token > right_token {
	       left_token_loc = right_token
	   } else {
	       left_token_loc = left_token
	   }

	   var  right_token_loc int  = right_token
	*/

	var token_name string
	if name_index >= 0 && !(ToUpper(my.Name(name_index)) == "ERROR") {
		token_name = "\"" + my.Name(name_index) + "\""
	} else {
		token_name = ""
	}

	if msg_code == INVALID_CODE {
		if len(token_name) == 0 {
			msg_code = INVALID_CODE
		} else {
			msg_code = INVALID_TOKEN_CODE
		}
	}
	if msg_code == SCOPE_CODE {
		token_name = "\""
		var i int = my.ScopeSuffix(-name_index)
		for ; my.ScopeRhs(i) != 0; i++ {

			if !my.IsNullable(my.ScopeRhs(i)) {
				var symbol_index int
				if my.ScopeRhs(i) > my.NT_OFFSET {
					symbol_index = my.NonterminalIndex(my.ScopeRhs(i) - my.NT_OFFSET)
				} else {
					symbol_index = my.TerminalIndex(my.ScopeRhs(i))
				}

				if len(my.Name(symbol_index)) > 0 {
					if len(token_name) > 1 { // Not just starting quote?
						token_name += " " // add a space separator
					}
					token_name += my.Name(symbol_index)
				}
			}
		}
		token_name += "\""
	}
	my.tokStream.ReportError(msg_code, left_token, right_token, []string{token_name}, 0)
	return
}