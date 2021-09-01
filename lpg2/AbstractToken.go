package lpg2

type AbstractToken struct {
	kind         int
	startOffSet  int
	endOffSet    int
	tokenIndex   int
	adjunctIndex int
	iPrsStream   IPrsStream
}

func (a *AbstractToken) GetPrecedingAdjuncts() []IToken {
	panic("implement me")
}

func (a *AbstractToken) GetFollowingAdjuncts() []IToken {
	panic("implement me")
}

func NewAbstractToken(startOffSet int, endOffSet int, kind int,
	iPrsStream IPrsStream) *AbstractToken {
	return &AbstractToken{
		iPrsStream:   iPrsStream,
		startOffSet:  startOffSet,
		endOffSet:    endOffSet,
		kind:         kind,
		tokenIndex:   0,
		adjunctIndex: 0}
}

func (a *AbstractToken) GetKind() int {
	return a.kind
}
func (a *AbstractToken) SetKind(kind int) {
	a.kind = kind
}
func (a *AbstractToken) GetStartOffSet() int {
	return a.startOffSet
}
func (a *AbstractToken) SetStartOffSet(startOffSet int) {
	a.startOffSet = startOffSet
}

func (a *AbstractToken) GetEndOffSet() int {
	return a.endOffSet
}

func (a *AbstractToken) SetEndOffSet(endOffSet int) {
	a.endOffSet = endOffSet
}
func (a *AbstractToken) GetTokenIndex() int {
	return a.tokenIndex
}
func (a *AbstractToken) SetTokenIndex(tokenIndex int) {
	a.tokenIndex = tokenIndex
}
func (a *AbstractToken) SetAdjunctIndex(adjunctIndex int) {
	a.adjunctIndex = adjunctIndex
}
func (a *AbstractToken) GetAdjunctIndex() int {
	return a.adjunctIndex
}
func (a *AbstractToken) GetIPrsStream() IPrsStream {
	return a.iPrsStream
}
func (a *AbstractToken) GetILexStream() ILexStream {
	if a.iPrsStream == nil {
		return nil
	} else {
		return a.iPrsStream.GetILexStream()
	}
}

func (a *AbstractToken) GetLine() int {
	if a.iPrsStream == nil {
		return 0
	} else {
		return a.iPrsStream.GetILexStream().GetLineNumberOfCharAt(a.startOffSet)
	}

}
func (a *AbstractToken) GetColumn() int {
	if a.iPrsStream == nil {
		return 0
	} else {
		return a.iPrsStream.GetILexStream().GetColumnOfCharAt(a.startOffSet)
	}
}

func (a *AbstractToken) GetEndLine() int {
	if a.iPrsStream == nil {
		return 0
	} else {
		return a.iPrsStream.GetILexStream().GetLineNumberOfCharAt(a.endOffSet)
	}
}

func (a *AbstractToken) GetEndColumn() int {
	if a.iPrsStream == nil {
		return 0
	} else {
		return a.iPrsStream.GetILexStream().GetColumnOfCharAt(a.endOffSet)
	}
}

func (a *AbstractToken) ToString() string {
	if a.iPrsStream == nil {
		return "<ToString>"
	} else {
		return a.iPrsStream.ToString(a, a)
	}
}
