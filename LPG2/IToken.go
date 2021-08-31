package lpg2

const EOF int = 0xffff

type IToken interface {
	getKind() int

	setKind(kind int)

	getStartOffset() int

	setStartOffset(startOffset int)

	getEndOffset() int

	setEndOffset(endOffset int)

	getTokenIndex() int

	setTokenIndex(i int)

	getAdjunctIndex() int

	setAdjunctIndex(i int)

	getPrecedingAdjuncts() []IToken

	getFollowingAdjuncts() []IToken

	getILexStream() ILexStream

	getIPrsStream() IPrsStream

	getLine() int

	getColumn() int

	getEndLine() int

	getEndColumn() int

	toString() string
}
