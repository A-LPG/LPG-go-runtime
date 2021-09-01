package lpg2

type ErrorToken  struct{
    *Token
    firstToken IToken
    lastToken IToken
    errorToken IToken
}


func NewErrorToken( firstToken IToken, lastToken IToken, errorToken IToken, startOffSet int, endOffSet int,
                 kind int)  *ErrorToken {

    t := new(ErrorToken)
    t.Token = NewToken(startOffSet,endOffSet,kind,nil)
    t.firstToken = firstToken
    t.lastToken = lastToken
    t.errorToken = errorToken
    return  t
}

func (a *ErrorToken) GetFirstToken()  IToken{
    return a.GetFirstRealToken()
}
func (a *ErrorToken) GetFirstRealToken()  IToken{
    return a.firstToken
}
func (a *ErrorToken) GetLastToken()  IToken{
    return a.GetLastRealToken()
}
func (a *ErrorToken) GetLastRealToken()  IToken{
    return a.lastToken
}
func (a *ErrorToken) GetErrorToken()  IToken{
    return a.errorToken
}
func (a *ErrorToken) GetPrecedingAdjuncts()  []IToken{
    return a.firstToken.GetPrecedingAdjuncts()
}
func (a *ErrorToken) GetFollowingAdjuncts()  []IToken {
    return a.lastToken.GetFollowingAdjuncts()
}
