package lpg2
type UnavailableParserInformationException struct {
    wrap error
}
func  NewUnavailableParserInformationException(info string) *UnavailableParserInformationException{
    t := new(UnavailableParserInformationException)
    if info == ""{
        info="Unavailable parser Information Exception"
    }
    t.wrap = NewErr(info)

    return t
}

func (a *UnavailableParserInformationException) ToString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *UnavailableParserInformationException) Error() string {
    return a.wrap.Error()
}

