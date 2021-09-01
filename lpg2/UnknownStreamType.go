package lpg2

type UnknownStreamType struct {
    wrap error
}
func  NewUnknownStreamType(info string) *UnknownStreamType{
    t := new(UnknownStreamType)
    if info == ""{
        info="UnknownStreamType"
    }
    t.wrap = NewErr(info)

    return t
}

func (a *UnknownStreamType) ToString() string  {
    return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *UnknownStreamType) Error() string {
   return a.wrap.Error()
}