package lpg2
type BadParseSymFileException struct {
    wrap error
}
func  NewBadParseSymFileException(info string) *BadParseSymFileException{
    t := new(BadParseSymFileException)
    if info == ""{
        info="BadParseSymFileException"
    }
    t.wrap = NewErr(info)

    return t
}

func (a *BadParseSymFileException) toString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *BadParseSymFileException) Error() string {
    return a.wrap.Error()
}

