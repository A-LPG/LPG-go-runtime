package lpg2

type IAst interface {
	getNextAst()IAst

	getParent()IAst

	getLeftIToken() IToken

	getRightIToken() IToken

	getPrecedingAdjuncts() []IToken

	getFollowingAdjuncts() []IToken

	getChildren() AstArrayList

	getAllChildren() AstArrayList

	accept(IAstVisitor)
}
