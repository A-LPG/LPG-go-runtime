package lpg2

type Token  struct{
    *AbstractToken
}
    

func NewToken( startOffset int, endOffset int, kind int, iPrsStream IPrsStream) *Token {
    t := new (Token)
    t.AbstractToken = NewAbstractToken(startOffset,endOffset,kind,iPrsStream)
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


