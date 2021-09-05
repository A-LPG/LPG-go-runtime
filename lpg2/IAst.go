package lpg2

type IAst interface {

	GetNextAst()IAst

	SetParent(IAst)

	GetParent()IAst

	GetLeftIToken() IToken

	GetRightIToken() IToken

	GetPrecedingAdjuncts() []IToken

	GetFollowingAdjuncts() []IToken

	GetChildren() *ArrayList

	GetAllChildren() *ArrayList

	Accept(IAstVisitor)
}
