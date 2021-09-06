package lpg

import "testing"

type LpgVisitor struct {
	*AbstractVisitor
}
func (my *LpgVisitor)     VisitLPG(n  *LPG) bool{
	println(n.ToString())
	println(n.GetAllChildren().Size())
	return true
}
	func NewLpgVisitor() *LpgVisitor {
	t :=new(LpgVisitor)
	t.AbstractVisitor = NewAbstractVisitor(t)
	return t
}
func TestLPG(t *testing.T)  {
	var filename = "jikespg.g"
	var lexer = NewLPGLexer(filename,4,nil) // Create the lexer

	var parser,e = NewLPGParser(lexer.GetILexStream()) // Create the parser
	if e !=nil{
		println(e.Error())
		return
	}
	//lexer.PrintTokens = true
	lexer.Lexer(parser.GetIPrsStream(),nil)
	var ast,err = parser.Parser()
	if err != nil{
		t.Error(err.Error())
	}else{
		var  v = NewLpgVisitor()
		ast.(*LPG).Accept(v)
		//v.Visit(ast.(lpg2.IAst))
		t.Log(ast)
	}

}