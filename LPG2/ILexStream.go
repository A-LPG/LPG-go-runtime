package lpg2

type ILexStream interface {
	TokenStream

	getIPrsStream() IPrsStream

	setPrsStream(stream IPrsStream)

	getLineCount() int

	orderedExportedSymbols() []string

	getLineOffset(i int) int

	getLineNumberOfCharAt(i int) int

	getColumnOfCharAt(i int) int

	getCharValue(i int) string
	getInputChars() string

	getIntValue(i int) int

	makeToken(startLoc int, endLoc int, kind int)

	setMessageHandler(handler IMessageHandler)

	getMessageHandler() IMessageHandler

	getLocation(leftLoc int, rightLoc int) []int
	reportLexicalErrorPosition(leftLoc int, rightLoc int)
	reportLexicalError(leftLoc int, rightLoc int,
		errorCode int , errorLeftLoc int , errorRightLoc int , errorInfo []string)

	toString(startOffset int, endOffset int) string
}
