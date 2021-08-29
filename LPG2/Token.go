package lpg2

type Token  struct{
    *AbstractToken
}
    

func NewToken( startOffset int, endOffset int, kind int, iPrsStream IPrsStream) *Token {
    t := new (Token)
    t.AbstractToken = new(AbstractToken)
    t.startOffset=startOffset
    t.endOffset=endOffset
    t.kind=kind
    t.iPrsStream=iPrsStream
    return  t
}

func (a *Token) getFollowingAdjuncts() []IToken {
    var stream = a.getIPrsStream()
    if nil == stream {
        return nil
    }else{
        return stream.getFollowingAdjuncts(a.getTokenIndex())
    }
}
func (a *Token) getPrecedingAdjuncts() []IToken {
    var stream = a.getIPrsStream()
    if nil == stream {
        return nil
    }else{
        return stream.getPrecedingAdjuncts(a.getTokenIndex())
    }
}


