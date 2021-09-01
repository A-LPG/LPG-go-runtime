package lpg2

type IAbstractArrayList interface {
      Size() int 
      GetElementAt(i int)  interface{}
      GetList() AstArrayList
      Add(elt interface{})
}