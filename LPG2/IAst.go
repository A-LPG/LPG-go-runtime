package lpg2

type IAst interface {
    
     getNextAst()
        
     getParent()
        
     getLeftIToken()  IToken
        
     getRightIToken()  IToken
        
     getPrecedingAdjuncts()  []IToken
        
     getFollowingAdjuncts()  []IToken
        
     getChildren()  ArrayList
        
     getAllChildren()  ArrayList
        
     accept(IAstVisitor)   
}