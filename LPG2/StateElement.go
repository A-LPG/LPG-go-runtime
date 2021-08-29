package lpg2

type StateElement struct {
    parent *StateElement 
    children *StateElement 
    siblings *StateElement 
    number int 
}
func NewStateElement() * StateElement{
    return &StateElement{

    }
}