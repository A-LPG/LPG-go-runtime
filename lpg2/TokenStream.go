package lpg2

type TokenStream interface {

	GetTokenFromEndToken(end_token int) int

	GetToken() int

	GetKind(i int) int

	GetNext(i int) int

	GetPrevious(i int) int

	GetName(i int) string

	Peek() int
	Reset()
	ResetTo(i int)

	BadToken() int

	GetLine(i int) int

	GetColumn(i int) int

	GetEndLine(i int) int

	GetEndColumn(i int) int

	AfterEol(i int) bool

	GetFileName() string

	GetStreamLength() int

	GetFirstRealToken(i int) int

	GetLastRealToken(i int) int

	ReportError(errorCode int, leftToken int, rightToken int, errorInfo []string, errorToken int)
}
