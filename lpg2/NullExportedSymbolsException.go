package lpg2
type NullExportedSymbolsException struct {
    wrap error
}
func  NewNullExportedSymbolsException(info string) *NullExportedSymbolsException{
    t := new(NullExportedSymbolsException)
    if info == ""{
        info="NullExportedSymbolsException"
    }
    t.wrap = NewErr(info)

    return t
}

func (a *NullExportedSymbolsException) ToString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *NullExportedSymbolsException) Error() string {
    return a.wrap.Error()
}
