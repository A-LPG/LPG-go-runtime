package lpg2

type IPrsStream interface {
	TokenStream

	GetMessageHandler() IMessageHandler

	SetMessageHandler(handler IMessageHandler)

	GetILexStream() ILexStream

	SetLexStream(lexStream ILexStream)

	// /**
	// * @deprecated replaced by @link //GetFirstRealToken()
	// *
	// */
	// GetFirstErrorToken(i int)

	// /**
	// * @deprecated replaced by @link //GetLastRealToken()
	// *
	// */
	// GetLastErrorToken(i int)

	MakeToken(startLoc int, endLoc int, kind int)

	MakeAdjunct(startLoc int, endLoc int, kind int)

	RemoveLastToken()

	GetLineCount() int

	GetSize() int

	RemapTerminalSymbols(ordered_parser_symbols []string, eof_symbol int) error

	OrderedTerminalSymbols() []string

	MapKind(kind int) int

	ResetTokenStream()

	GetStreamIndex() int

	ResetStreamLength()
	SetStreamIndex(index int)

	SetStreamLength(length int )

	AddToken(token IToken)

	AddAdjunct(adjunct IToken)

	OrderedExportedSymbols() []string

	GetTokens() *TokenArrayList

	GetAdjuncts() *TokenArrayList

	GetFollowingAdjuncts(i int) []IToken

	GetPrecedingAdjuncts(i int) []IToken

	GetIToken(i int) IToken

	GetTokenText(i int) string

	GetStartOffset(i int) int
	
	GetEndOffSet(i int) int

	GetLineOffSet(i int) int

	GetLineNumberOfCharAt(i int) int

	GetColumnOfCharAt(i int) int

	GetTokenLength(i int) int

	GetLineNumberOfTokenAt(i int) int

	GetEndLineNumberOfTokenAt(i int) int

	GetColumnOfTokenAt(i int) int

	GetEndColumnOfTokenAt(i int) int

	GetInputChars() []rune

	ToStringFromIndex(first_token int, last_token int) string

	ToString(t1 IToken, t2 IToken) string

	GetTokenIndexAtCharacter(offSet int) int

	GetTokenAtCharacter(offSet int) IToken

	GetTokenAt(i int) IToken

	DumpTokens()

	DumpToken(i int)

	MakeErrorToken(first int, last int, error int, kind int) int
}
