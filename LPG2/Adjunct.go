package lpg2



type Adjunct struct {
    *AbstractToken
}
    

    func (a *Adjunct)  __init__( start_offset int, end_offset int, kind int, prs_stream: IPrsStream = nil){

    }
        super().__init__(start_offset, end_offset, kind, prs_stream)

    func (a *Adjunct)  getFollowingAdjuncts() []IToken{

    }


    func (a *Adjunct)  getPrecedingAdjuncts() []IToken {
        return nil
    }

