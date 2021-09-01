package lpg2
type MismatchedInputCharsException struct {
    wrap error
}
func  NewMismatchedInputCharsException(info string) *MismatchedInputCharsException{
    t := new(MismatchedInputCharsException)
    if info == ""{
        info="MismatchedInputCharsException"
    }
    t.wrap = NewErr(info)

    return t
}

func (a *MismatchedInputCharsException) ToString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *MismatchedInputCharsException) Error() string {
    return a.wrap.Error()
}
