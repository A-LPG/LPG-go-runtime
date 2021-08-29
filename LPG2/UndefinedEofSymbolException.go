package lpg2

type UndefinedEofSymbolException struct {
    wrap error
}
func  NewUndefinedEofSymbolException(info string) *UndefinedEofSymbolException{
    t := new(UndefinedEofSymbolException)
    if info == ""{
        info="UndefinedEofSymbolException"
    }
    t.wrap = NewErr(info)

    return t
}

func (a *UndefinedEofSymbolException) toString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *UndefinedEofSymbolException) Error() string {
    return a.wrap.Error()
}
