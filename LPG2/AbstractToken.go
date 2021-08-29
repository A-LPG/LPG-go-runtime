package lpg2
type AbstractToken  struct{
	kind   int
	startOffset   int
	endOffset   int
	tokenIndex   int
	adjunctIndex   int
	iPrsStream  IPrsStream
}

func (a *AbstractToken) getPrecedingAdjuncts() []IToken {
	panic("implement me")
}

func (a *AbstractToken) getFollowingAdjuncts() []IToken {
	panic("implement me")
}

func  NewAbstractToken(startOffset int , endOffset int, kind int ,
	iPrsStream IPrsStream) *AbstractToken{
     return  &AbstractToken{
	iPrsStream :iPrsStream,
	startOffset : startOffset,
	endOffset : endOffset,
	kind : kind,
	tokenIndex  : 0,
	adjunctIndex  : 0}
}

func (a *AbstractToken) getKind() int {
	return a.kind
}
func (a *AbstractToken) setKind( kind int){
	a.kind = kind
}
func (a *AbstractToken) getStartOffset() int {
	return a.startOffset
}
func (a *AbstractToken) setStartOffset( startOffset int){
	a.startOffset = startOffset
}


func (a *AbstractToken) getEndOffset() int {
	return a.endOffset
}

func (a *AbstractToken) setEndOffset( endOffset int){
	a.endOffset = endOffset
}
func (a *AbstractToken) getTokenIndex() int {
	return a.tokenIndex
}
func (a *AbstractToken) setTokenIndex( tokenIndex int){
	a.tokenIndex = tokenIndex
}
func (a *AbstractToken) setAdjunctIndex( adjunctIndex int) {
	a.adjunctIndex = adjunctIndex
}
func (a *AbstractToken) getAdjunctIndex() int {
	return a.adjunctIndex
}
func (a *AbstractToken) getIPrsStream() IPrsStream{
	return a.iPrsStream
}
func (a *AbstractToken) getILexStream() ILexStream {
	if a.iPrsStream == nil {
		return nil
	}else{
		return a.iPrsStream.getILexStream()
	}
}


func (a *AbstractToken) getLine() int {
	if a.iPrsStream == nil {
		return 0
	}else{
		return a.iPrsStream.getILexStream().getLineNumberOfCharAt(a.startOffset)
	}

}
func (a *AbstractToken) getColumn() int {
	if a.iPrsStream == nil {
		return 0
	}else{
		return a.iPrsStream.getILexStream().getColumnOfCharAt(a.startOffset)
	}
}

func (a *AbstractToken) getEndLine() int {
	if a.iPrsStream == nil {
		return 0
	}else{
		return a.iPrsStream.getILexStream().getLineNumberOfCharAt(a.endOffset)
	}
}

func (a *AbstractToken) getEndColumn() int {
	if a.iPrsStream == nil {
		return 0
	}else{
		return a.iPrsStream.getILexStream().getColumnOfCharAt(a.endOffset)
	}
}

func (a *AbstractToken) toString()  string {
	if a.iPrsStream == nil {
		return "<toString>"
	}else{
		return a.iPrsStream.toString(a,a)
	}
}
