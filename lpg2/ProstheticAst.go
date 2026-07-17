package lpg2

// ProstheticAst synthesizes a placeholder AST node for a %Recover nonterminal
// that the backtracking parser replays as an ErrorToken (inserted by scope
// recovery). It is invoked with the error token and returns a freshly built
// node that the parser pushes onto the value stack.
type ProstheticAst func(errorToken IToken) IAst

// ProstheticAstProvider is implemented by generated parser/action types that
// declare %Recover symbols. The backtracking parser type-asserts against it so
// that grammars without %Recover keep their historical throw behavior.
type ProstheticAstProvider interface {
	GetProstheticAst() []ProstheticAst
}

// ProsthesisIndexProvider is implemented by generated parse tables that declare
// %Recover symbols. It maps a replayed nonterminal token kind (NT_OFFSET
// already applied) to a compact slot in ProstheticAstProvider.GetProstheticAst.
type ProsthesisIndexProvider interface {
	GetProsthesisIndex(index int) int
}
