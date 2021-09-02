package lpg2
type BadParseException struct {
    wrap       error
    ErrorToken int

}
func  NewBadParseException(errorToken int) *BadParseException{
    t := new(BadParseException)
    t.wrap = NewErr("BadParseException")
    t.ErrorToken = errorToken
    return t
}

func (a *BadParseException) ToString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *BadParseException) Error() string {
    return a.wrap.Error()
}






