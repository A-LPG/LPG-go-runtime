package lpg2

type IAstVisitor interface {
      PreVisit(element IAst)  bool 
      PostVisit(element IAst)    
}