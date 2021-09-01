package lpg2

type Adjunct struct {
	*AbstractToken
}

func NewAdjunct(startOffSet int, endOffSet int, kind int, prsStream IPrsStream) *Adjunct {
	t := new(Adjunct)
	t.AbstractToken = NewAbstractToken(startOffSet, endOffSet, kind, prsStream)
	return t
}

func (a *Adjunct) GetFollowingAdjuncts() []IToken {
	return nil
}

func (a *Adjunct) GetPrecedingAdjuncts() []IToken {
	return nil
}
