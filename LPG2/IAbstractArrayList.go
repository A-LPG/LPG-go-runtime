package lpg2

type IAbstractArrayList interface {
      size() int 
      getElementAt(i int)  interface{}
      getList() AstArrayList
      add(elt interface{})
}