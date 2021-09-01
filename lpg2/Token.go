package lpg2

type Token  struct{
    *AbstractToken
}
    

func NewToken( startOffSet int, endOffSet int, kind int, iPrsStream IPrsStream) *Token {
    t := new (Token)
    t.AbstractToken = NewAbstractToken(startOffSet,endOffSet,kind,iPrsStream)
    return  t
}

func (a *Token) GetFollowingAdjuncts() []IToken {
    var stream = a.GetIPrsStream()
    if nil == stream {
        return nil
    }else{
        return stream.GetFollowingAdjuncts(a.GetTokenIndex())
    }
}
func (a *Token) GetPrecedingAdjuncts() []IToken {
    var stream = a.GetIPrsStream()
    if nil == stream {
        return nil
    }else{
        return stream.GetPrecedingAdjuncts(a.GetTokenIndex())
    }
}


