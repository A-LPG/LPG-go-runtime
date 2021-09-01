package lpg2
type BadParseException struct {
    wrap error
    error_token int
}
func  NewBadParseException(error_token int) *BadParseException{
    t := new(BadParseException)
    t.wrap = NewErr("BadParseException")
    t.error_token = error_token
    return t
}

func (a *BadParseException) ToString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *BadParseException) Error() string {
    return a.wrap.Error()
}






