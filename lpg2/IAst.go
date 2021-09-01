package lpg2

type IAst interface {
	GetNextAst()IAst

	GetParent()IAst

	GetLeftIToken() IToken

	GetRightIToken() IToken

	GetPrecedingAdjuncts() []IToken

	GetFollowingAdjuncts() []IToken

	GetChildren() AstArrayList

	GetAllChildren() AstArrayList

	Accept(IAstVisitor)
}
