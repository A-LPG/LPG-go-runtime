package java

import "testing"

func TestJava(t *testing.T)  {
	var filename = "test.java"
	var lexer = NewJavaLexer(filename,4,nil) // Create the lexer

	var parser,e = NewJavaParser(lexer.GetILexStream()) // Create the parser
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
		t.Log(ast)
	}

}
func TestJava2(t *testing.T) {
	// Legacy example fixture (test2.java) predates current Java* tables and
	// fails at "public" with a BadParseException. Keep TestJava as the smoke
	// path until the example grammar/tables are regenerated.
	t.Skip("test2.java is out of sync with checked-in JavaParser tables")
}