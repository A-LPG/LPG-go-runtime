package lpg2
type NotDeterministicParseTableException struct {
    wrap error
}
func  NewNotDeterministicParseTableException(info string) *NotDeterministicParseTableException{
    t := new(NotDeterministicParseTableException)
    if info == ""{
        info="NotDeterministicParseTableException"
    }
    t.wrap = NewErr(info)

    return t
}

func (a *NotDeterministicParseTableException) toString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *NotDeterministicParseTableException) Error() string {
    return a.wrap.Error()
}

