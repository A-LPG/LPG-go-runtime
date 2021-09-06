%options automatic_ast=nested,variables=nt,visitor=default
%options programming_language=go
%options template= btParserTemplateF.gi
%options package=lpg
$Terminals
 IntegerLiteral
 PLUS ::= +
 MULTIPLY ::= *
 LPAREN ::= (
 RPAREN ::= )
$end
$Rules
 E ::= E + T
 | T
 T ::= T * F
 | F
 F ::= IntegerLiteral
 F$ParenExpr ::= ( E )
$End