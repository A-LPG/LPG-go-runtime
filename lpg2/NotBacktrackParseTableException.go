package lpg2
type NotBacktrackParseTableException struct {
    wrap error
}
func  NewNotBacktrackParseTableException(info string) *NotBacktrackParseTableException{
    t := new(NotBacktrackParseTableException)
    if info == ""{
        info="NotBacktrackParseTableException"
    }
    t.wrap = NewErr(info)

    return t
}

func (a *NotBacktrackParseTableException) ToString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *NotBacktrackParseTableException) Error() string {
    return a.wrap.Error()
}

