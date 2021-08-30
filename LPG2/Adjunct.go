package lpg2



type Adjunct struct {
    *AbstractToken
}
    

func NewAdjunct( start_offset int, end_offset int, kind int, prs_stream IPrsStream) *Adjunct{
     t := new(Adjunct)
     t.AbstractToken=NewAbstractToken(start_offset, end_offset, kind, prs_stream)
     return t
}
    
func (a *Adjunct)  getFollowingAdjuncts() []IToken{
    return nil
}


func (a *Adjunct)  getPrecedingAdjuncts() []IToken {
    return nil
}

