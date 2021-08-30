package lpg2

type TokenStream interface {

	getTokenFromEndToken(end_token int) int

	getToken() int

	getKind(i int) int

	getNext(i int) int

	getPrevious(i int) int

	getName(i int) string

	peek() int
	reset()
	resetTo(i int)

	badToken() int

	getLine(i int) int

	getColumn(i int) int

	getEndLine(i int) int

	getEndColumn(i int) int

	afterEol(i int) bool

	getFileName() string

	getStreamLength() int

	getFirstRealToken(i int) int

	getLastRealToken(i int) int

	reportError(errorCode int, leftToken int, rightToken int, errorInfo []string, errorToken int)
}
