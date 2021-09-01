package lpg2
type NullPointerException struct {
	wrap error
}
func  NewNullPointerException(info string) *NullPointerException{
	t := new(NullPointerException)
	if info == ""{
		info="NullPointerException"
	}
	t.wrap = NewErr(info)

	return t
}

func (a *NullPointerException) ToString() string  {
	return a.wrap.Error()
}

// Error implements the interface of Error, it returns all the error as string.
func (a *NullPointerException) Error() string {
	return a.wrap.Error()
}
