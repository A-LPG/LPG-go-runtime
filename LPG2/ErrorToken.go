package lpg2

type ErrorToken  struct{
    *Token
    firstToken IToken
    lastToken IToken
    errorToken IToken
}


func NewErrorToken( firstToken IToken, lastToken IToken, errorToken IToken, startOffset int, endOffset int,
                 kind int)  *ErrorToken {

    t := new(ErrorToken)
    t.Token = NewToken(startOffset,endOffset,kind,nil)
    t.firstToken = firstToken
    t.lastToken = lastToken
    t.errorToken = errorToken
    return  t
}

func (a *ErrorToken) getFirstToken()  IToken{
    return a.getFirstRealToken()
}
func (a *ErrorToken) getFirstRealToken()  IToken{
    return a.firstToken
}
func (a *ErrorToken) getLastToken()  IToken{
    return a.getLastRealToken()
}
func (a *ErrorToken) getLastRealToken()  IToken{
    return a.lastToken
}
func (a *ErrorToken) getErrorToken()  IToken{
    return a.errorToken
}
func (a *ErrorToken) getPrecedingAdjuncts()  []IToken{
    return a.firstToken.getPrecedingAdjuncts()
}
func (a *ErrorToken) getFollowingAdjuncts()  []IToken {
    return a.lastToken.getFollowingAdjuncts()
}
