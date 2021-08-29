package lpg2

type TokenStream interface {

	getToken(end_token int) int

	getKind(i int) int

	getNext(i int) int

	getPrevious(i int) int

	getName(i int) string

	peek() int
	resetDefault()
	reset(i int)

	badToken() int

	getLine(i int) int

	getColumn(i int) int

	getEndLine(i int) int

	getEndColumn(i int) int

	afterEol(i int) int

	getFileName() string

	getStreamLength() int

	getFirstRealToken(i int) int

	getLastRealToken(i int) int

	reportError(errorCode int, leftToken int, rightToken int, errorInfo string, errorToken int)
}
