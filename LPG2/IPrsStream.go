package lpg2

type IPrsStream interface {
	TokenStream

	getMessageHandler() IMessageHandler

	setMessageHandler(handler IMessageHandler)

	getILexStream() ILexStream

	setLexStream(lexStream ILexStream)

	// /**
	// * @deprecated replaced by @link //getFirstRealToken()
	// *
	// */
	// getFirstErrorToken(i int)

	// /**
	// * @deprecated replaced by @link //getLastRealToken()
	// *
	// */
	// getLastErrorToken(i int)

	makeToken(startLoc int, endLoc int, kind int)

	makeAdjunct(startLoc int, endLoc int, kind int)

	removeLastToken()

	getLineCount() int

	getSize() int

	remapTerminalSymbols(ordered_parser_symbols []string, eof_symbol int) error

	orderedTerminalSymbols() []string

	mapKind(kind int) int

	resetTokenStream()

	getStreamIndex() int

	resetStreamLength()
	setStreamIndex(index int)

	setStreamLength(length int )

	addToken(token IToken)

	addAdjunct(adjunct IToken)

	orderedExportedSymbols() []string

	getTokens() *TokenArrayList

	getAdjuncts() *TokenArrayList

	getFollowingAdjuncts(i int) []IToken

	getPrecedingAdjuncts(i int) []IToken

	getIToken(i int) IToken

	getTokenText(i int) string

	getStartOffset(i int) int

	getEndOffset(i int) int

	getLineOffset(i int) int

	getLineNumberOfCharAt(i int) int

	getColumnOfCharAt(i int) int

	getTokenLength(i int) int

	getLineNumberOfTokenAt(i int) int

	getEndLineNumberOfTokenAt(i int) int

	getColumnOfTokenAt(i int) int

	getEndColumnOfTokenAt(i int) int

	getInputChars() string

	toStringFromIndex(first_token int, last_token int) string

	toString(t1 IToken, t2 IToken) string

	getTokenIndexAtCharacter(offset int) int

	getTokenAtCharacter(offset int) IToken

	getTokenAt(i int) IToken

	dumpTokens()

	dumpToken(i int)

	makeErrorToken(first int, last int, error int, kind int) int
}
