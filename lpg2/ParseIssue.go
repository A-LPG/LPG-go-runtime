package lpg2

// SourceSpan is a byte/char offset range for structured parse diagnostics.
type SourceSpan struct {
	StartOffset int
	EndOffset   int
}

// ParseIssue is the unified parse-error shape: code / span / expected / got.
type ParseIssue struct {
	Code     int
	Span     SourceSpan
	Expected []string
	Got      string
}

// MismatchParseIssue builds a mismatch issue with expected filled from state.
func MismatchParseIssue(prs ParseTable, state, code int, span SourceSpan, got string) ParseIssue {
	return ParseIssue{
		Code:     code,
		Span:     span,
		Expected: ExpectedTerminalNames(prs, state),
		Got:      got,
	}
}
