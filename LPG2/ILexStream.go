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

	getIntValue(i int) int

	makeToken(start_loc int, end_loc int, kind int)

	setMessageHandler(handler IMessageHandler)

	getMessageHandler() IMessageHandler

	getLocation(left_loc int, right_loc int) []int

	reportLexicalError(left_loc int, right_loc int,
		error_code int , error_left_loc_arg int , error_right_loc_arg int , error_info []string)

	toString(startOffset int, endOffset int) string
}
