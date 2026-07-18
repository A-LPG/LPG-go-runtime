package lpg2

// NotGLRParseTableException indicates that GLRParser was given a table that
// was not generated with the -glr option.
type NotGLRParseTableException struct {
	wrap error
}

func NewNotGLRParseTableException(info string) *NotGLRParseTableException {
	if info == "" {
		info = "NotGLRParseTableException"
	}
	return &NotGLRParseTableException{wrap: NewErr(info)}
}

func (e *NotGLRParseTableException) ToString() string {
	return e.wrap.Error()
}

func (e *NotGLRParseTableException) Error() string {
	return e.wrap.Error()
}
