package lpg2
type UnimplementedTerminalsException struct {
    wrap error
    symbols *IntArrayList
}
func  NewUnimplementedTerminalsException(symbols *IntArrayList) *UnimplementedTerminalsException{
    t := new(UnimplementedTerminalsException)
    t.wrap = NewErr("UnimplementedTerminalsException")
    t.symbols = symbols
    return t
}

func (a *UnimplementedTerminalsException) ToString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *UnimplementedTerminalsException) Error() string {
    return a.wrap.Error()
}



func (a *UnimplementedTerminalsException) GetSymbols()  *IntArrayList{
    return a.symbols
}

