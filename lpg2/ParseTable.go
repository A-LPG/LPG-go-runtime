package lpg2

type ParseTable interface {
	BaseCheck(index int) int

	Rhs(index int) int

	BaseAction(index int) int

	Lhs(index int) int

	TermCheck(index int) int

	TermAction(index int) int

	Asb(index int) int

	Asr(index int) int

	Nasb(index int) int

	Nasr(index int) int

	TerminalIndex(index int) int

	NonterminalIndex(index int) int

	ScopePrefix(index int) int

	ScopeSuffix(index int) int

	ScopeLhs(index int) int

	ScopeLa(index int) int

	ScopeStateSet(index int) int

	ScopeRhs(index int) int

	ScopeState(index int) int

	InSymb(index int) int

	Name(index int) string

	OriginalState(state int) int

	Asi(state int) int

	Nasi(state int) int

	InSymbol(state int) int

	NtAction(state int, sym int) int

	TAction(act int, sym int) int

	LookAhead(act int, sym int) int

	GetErrorSymbol() int

	GetScopeUbound() int

	GetScopeSize() int

	GetMaxNameLength() int

	GetNumStates() int

	GetNtOffSet() int

	GetLaStateOffSet() int

	GetMaxLa() int

	GetNumRules() int

	GetNumNonterminals() int

	GetNumSymbols() int

	GetStartState() int

	GetStartSymbol() int

	GetEoftSymbol() int

	GetEoltSymbol() int

	GetAcceptAction() int

	GetErrorAction() int

	IsNullable(symbol int) bool

	IsValidForParser() bool

	GetBacktrack() bool
}
