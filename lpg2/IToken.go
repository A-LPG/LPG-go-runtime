package lpg2

const EOF int = 0xffff

type IToken interface {
	GetKind() int

	SetKind(kind int)

	GetStartOffset() int

	SetStartOffset(startOffset int)

	GetEndOffset() int

	SetEndOffset(endOffset int)

	GetTokenIndex() int

	SetTokenIndex(i int)

	GetAdjunctIndex() int

	SetAdjunctIndex(i int)

	GetPrecedingAdjuncts() []IToken

	GetFollowingAdjuncts() []IToken

	GetILexStream() ILexStream

	GetIPrsStream() IPrsStream

	GetLine() int

	GetColumn() int

	GetEndLine() int

	GetEndColumn() int

	ToString() string
}
