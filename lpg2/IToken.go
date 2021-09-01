package lpg2

const EOF int = 0xffff

type IToken interface {
	GetKind() int

	SetKind(kind int)

	GetStartOffSet() int

	SetStartOffSet(startOffSet int)

	GetEndOffSet() int

	SetEndOffSet(endOffSet int)

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
