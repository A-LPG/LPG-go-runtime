package lpg2

type TokenStreamNotIPrsStreamException struct {
    wrap error
}
func  NewTokenStreamNotIPrsStreamException(info string) *TokenStreamNotIPrsStreamException{
    t := new(TokenStreamNotIPrsStreamException)
    if info == ""{
        info="TokenStreamNotIPrsStreamException"
    }
    t.wrap = NewErr(info)

    return t
}

func (a *TokenStreamNotIPrsStreamException) toString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *TokenStreamNotIPrsStreamException) Error() string {
    return a.wrap.Error()
}

