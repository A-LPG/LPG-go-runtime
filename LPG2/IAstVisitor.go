package lpg2

type IAstVisitor interface {
      preVisit(element IAst)  bool 
      postVisit(element IAst)    
}