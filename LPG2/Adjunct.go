package lpg2

type Adjunct struct {
	*AbstractToken
}

func NewAdjunct(startOffset int, endOffset int, kind int, prsStream IPrsStream) *Adjunct {
	t := new(Adjunct)
	t.AbstractToken = NewAbstractToken(startOffset, endOffset, kind, prsStream)
	return t
}

func (a *Adjunct) getFollowingAdjuncts() []IToken {
	return nil
}

func (a *Adjunct) getPrecedingAdjuncts() []IToken {
	return nil
}
