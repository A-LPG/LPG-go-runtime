package lpg2

type NullTerminalSymbolsException struct {
    wrap error
}
func  NewNullTerminalSymbolsException(info string) *NullTerminalSymbolsException{
    t := new(NullTerminalSymbolsException)
    if info == ""{
        info="NullTerminalSymbolsException"
    }
    t.wrap = NewErr(info)

    return t
}

func (a *NullTerminalSymbolsException) ToString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *NullTerminalSymbolsException) Error() string {
    return a.wrap.Error()
}


