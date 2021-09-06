package lpg
type __LPGParsersym__ struct{
   TK_EQUIVALENCE int
   TK_PRIORITY_EQUIVALENCE int
   TK_ARROW int
   TK_PRIORITY_ARROW int
   TK_OR_MARKER int
   TK_EQUAL int
   TK_COMMA int
   TK_LEFT_PAREN int
   TK_RIGHT_PAREN int
   TK_LEFT_BRACKET int
   TK_RIGHT_BRACKET int
   TK_SHARP int
   TK_ALIAS_KEY int
   TK_AST_KEY int
   TK_DEFINE_KEY int
   TK_DISJOINTPREDECESSORSETS_KEY int
   TK_DROPRULES_KEY int
   TK_DROPSYMBOLS_KEY int
   TK_EMPTY_KEY int
   TK_END_KEY int
   TK_ERROR_KEY int
   TK_EOL_KEY int
   TK_EOF_KEY int
   TK_EXPORT_KEY int
   TK_GLOBALS_KEY int
   TK_HEADERS_KEY int
   TK_IDENTIFIER_KEY int
   TK_IMPORT_KEY int
   TK_INCLUDE_KEY int
   TK_KEYWORDS_KEY int
   TK_NAMES_KEY int
   TK_NOTICE_KEY int
   TK_OPTIONS_KEY int
   TK_RECOVER_KEY int
   TK_RULES_KEY int
   TK_SOFT_KEYWORDS_KEY int
   TK_START_KEY int
   TK_TERMINALS_KEY int
   TK_TRAILERS_KEY int
   TK_TYPES_KEY int
   TK_EOF_TOKEN int
   TK_SINGLE_LINE_COMMENT int
   TK_MACRO_NAME int
   TK_SYMBOL int
   TK_BLOCK int
   TK_VBAR int
   TK_ERROR_TOKEN int

   OrderedTerminalSymbols []string

   NumTokenKinds int

   IsValidForParser  bool
}
func New__LPGParsersym__() *__LPGParsersym__{
    my := new(__LPGParsersym__)
   my.TK_EQUIVALENCE = 5
   my.TK_PRIORITY_EQUIVALENCE = 6
   my.TK_ARROW = 7
   my.TK_PRIORITY_ARROW = 8
   my.TK_OR_MARKER = 14
   my.TK_EQUAL = 38
   my.TK_COMMA = 37
   my.TK_LEFT_PAREN = 39
   my.TK_RIGHT_PAREN = 40
   my.TK_LEFT_BRACKET = 42
   my.TK_RIGHT_BRACKET = 43
   my.TK_SHARP = 44
   my.TK_ALIAS_KEY = 15
   my.TK_AST_KEY = 16
   my.TK_DEFINE_KEY = 17
   my.TK_DISJOINTPREDECESSORSETS_KEY = 18
   my.TK_DROPRULES_KEY = 19
   my.TK_DROPSYMBOLS_KEY = 20
   my.TK_EMPTY_KEY = 12
   my.TK_END_KEY = 3
   my.TK_ERROR_KEY = 9
   my.TK_EOL_KEY = 10
   my.TK_EOF_KEY = 13
   my.TK_EXPORT_KEY = 21
   my.TK_GLOBALS_KEY = 22
   my.TK_HEADERS_KEY = 23
   my.TK_IDENTIFIER_KEY = 11
   my.TK_IMPORT_KEY = 24
   my.TK_INCLUDE_KEY = 25
   my.TK_KEYWORDS_KEY = 26
   my.TK_NAMES_KEY = 27
   my.TK_NOTICE_KEY = 28
   my.TK_OPTIONS_KEY = 41
   my.TK_RECOVER_KEY = 29
   my.TK_RULES_KEY = 30
   my.TK_SOFT_KEYWORDS_KEY = 31
   my.TK_START_KEY = 32
   my.TK_TERMINALS_KEY = 33
   my.TK_TRAILERS_KEY = 34
   my.TK_TYPES_KEY = 35
   my.TK_EOF_TOKEN = 36
   my.TK_SINGLE_LINE_COMMENT = 45
   my.TK_MACRO_NAME = 2
   my.TK_SYMBOL = 1
   my.TK_BLOCK = 4
   my.TK_VBAR = 46
   my.TK_ERROR_TOKEN = 47

   my.OrderedTerminalSymbols = []string{
                 "",
                 "SYMBOL",
                 "MACRO_NAME",
                 "END_KEY",
                 "BLOCK",
                 "EQUIVALENCE",
                 "PRIORITY_EQUIVALENCE",
                 "ARROW",
                 "PRIORITY_ARROW",
                 "ERROR_KEY",
                 "EOL_KEY",
                 "IDENTIFIER_KEY",
                 "EMPTY_KEY",
                 "EOF_KEY",
                 "OR_MARKER",
                 "ALIAS_KEY",
                 "AST_KEY",
                 "DEFINE_KEY",
                 "DISJOINTPREDECESSORSETS_KEY",
                 "DROPRULES_KEY",
                 "DROPSYMBOLS_KEY",
                 "EXPORT_KEY",
                 "GLOBALS_KEY",
                 "HEADERS_KEY",
                 "IMPORT_KEY",
                 "INCLUDE_KEY",
                 "KEYWORDS_KEY",
                 "NAMES_KEY",
                 "NOTICE_KEY",
                 "RECOVER_KEY",
                 "RULES_KEY",
                 "SOFT_KEYWORDS_KEY",
                 "START_KEY",
                 "TERMINALS_KEY",
                 "TRAILERS_KEY",
                 "TYPES_KEY",
                 "EOF_TOKEN",
                 "COMMA",
                 "EQUAL",
                 "LEFT_PAREN",
                 "RIGHT_PAREN",
                 "OPTIONS_KEY",
                 "LEFT_BRACKET",
                 "RIGHT_BRACKET",
                 "SHARP",
                 "SINGLE_LINE_COMMENT",
                 "VBAR",
                 "ERROR_TOKEN",
             }

   my.NumTokenKinds = 47
   my.IsValidForParser = true
   return my
}
var LPGParsersym = New__LPGParsersym__()
