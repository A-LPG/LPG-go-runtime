package lpg2

type ILexStream interface {
	TokenStream

	GetIPrsStream() IPrsStream

	SetPrsStream(stream IPrsStream)

	GetLineCount() int

	OrderedExportedSymbols() []string

	GetLineOffSet(i int) int

	GetLineNumberOfCharAt(i int) int

	GetColumnOfCharAt(i int) int

	GetCharValue(i int) string
	GetInputChars() []rune

	GetIntValue(i int) int

	MakeToken(startLoc int, endLoc int, kind int)

	SetMessageHandler(handler IMessageHandler)

	GetMessageHandler() IMessageHandler

	GetLocation(leftLoc int, rightLoc int) []int
	ReportLexicalErrorPosition(leftLoc int, rightLoc int)
	ReportLexicalError(leftLoc int, rightLoc int,
		errorCode int , errorLeftLoc int , errorRightLoc int , errorInfo []string)

	ToString(startOffSet int, endOffSet int) string
}
