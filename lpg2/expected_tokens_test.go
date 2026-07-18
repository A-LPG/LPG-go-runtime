package lpg2

import (
	"reflect"
	"testing"
)

type mockTable struct{}

func (mockTable) GetErrorAction() int { return 0 }
func (mockTable) GetNtOffset() int    { return 4 }
func (mockTable) TAction(state, sym int) int {
	if state == 0 && (sym == 1 || sym == 2) {
		return 1
	}
	return 0
}
func (mockTable) TerminalIndex(sym int) int { return sym }
func (mockTable) Name(index int) string {
	if index == 1 {
		return "a"
	}
	if index == 2 {
		return "b"
	}
	return ""
}

func (mockTable) BaseCheck(int) int             { return 0 }
func (mockTable) Rhs(int) int                   { return 0 }
func (mockTable) BaseAction(int) int            { return 0 }
func (mockTable) Lhs(int) int                   { return 0 }
func (mockTable) TermCheck(int) int             { return 0 }
func (mockTable) TermAction(int) int            { return 0 }
func (mockTable) Asb(int) int                   { return 0 }
func (mockTable) Asr(int) int                   { return 0 }
func (mockTable) Nasb(int) int                  { return 0 }
func (mockTable) Nasr(int) int                  { return 0 }
func (mockTable) NonterminalIndex(int) int      { return 0 }
func (mockTable) ScopePrefix(int) int           { return 0 }
func (mockTable) ScopeSuffix(int) int           { return 0 }
func (mockTable) ScopeLhs(int) int              { return 0 }
func (mockTable) ScopeLa(int) int               { return 0 }
func (mockTable) ScopeStateSet(int) int         { return 0 }
func (mockTable) ScopeRhs(int) int              { return 0 }
func (mockTable) ScopeState(int) int            { return 0 }
func (mockTable) InSymb(int) int                { return 0 }
func (mockTable) OriginalState(int) int         { return 0 }
func (mockTable) Asi(int) int                   { return 0 }
func (mockTable) Nasi(int) int                  { return 0 }
func (mockTable) InSymbol(int) int              { return 0 }
func (mockTable) NtAction(int, int) int         { return 0 }
func (mockTable) LookAhead(int, int) int        { return 0 }
func (mockTable) GetErrorSymbol() int           { return 0 }
func (mockTable) GetScopeUbound() int           { return 0 }
func (mockTable) GetScopeSize() int             { return 0 }
func (mockTable) GetMaxNameLength() int         { return 0 }
func (mockTable) GetNumStates() int             { return 0 }
func (mockTable) GetLaStateOffset() int         { return 0 }
func (mockTable) GetMaxLa() int                   { return 0 }
func (mockTable) GetNumRules() int              { return 0 }
func (mockTable) GetNumNonterminals() int       { return 0 }
func (mockTable) GetNumSymbols() int            { return 0 }
func (mockTable) GetStartState() int            { return 0 }
func (mockTable) GetStartSymbol() int           { return 0 }
func (mockTable) GetEoftSymbol() int            { return 0 }
func (mockTable) GetEoltSymbol() int            { return 0 }
func (mockTable) GetAcceptAction() int          { return 0 }
func (mockTable) IsNullable(int) bool           { return false }
func (mockTable) IsValidForParser() bool        { return true }
func (mockTable) GetBacktrack() bool            { return false }

func TestExpectedTerminalNames(t *testing.T) {
	var prs mockTable
	got := ExpectedTerminalNames(prs, 0)
	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExpectedTerminalNames() = %v, want %v", got, want)
	}
}

func TestMismatchParseIssue(t *testing.T) {
	var prs mockTable
	issue := MismatchParseIssue(prs, 0, ERROR_CODE, SourceSpan{1, 1}, "x")
	if issue.Code != ERROR_CODE || issue.Got != "x" {
		t.Fatalf("unexpected issue: %+v", issue)
	}
	if !reflect.DeepEqual(issue.Expected, []string{"a", "b"}) {
		t.Fatalf("expected %v, got %v", []string{"a", "b"}, issue.Expected)
	}
}
