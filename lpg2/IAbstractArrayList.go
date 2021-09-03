package lpg2

type IAbstractArrayList interface {
      Size() int 
      GetElementAt(i int) IAst
      GetList() *ArrayList
      Add(elt IAst)bool
}