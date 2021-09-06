//
// This is the grammar specification from the Final Draft of the generic spec.
//
////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2007 IBM Corporation.
// All rights reserved. This program and the accompanying materials
// are made available under the terms of the Eclipse Public License v1.0
// which accompanies this distribution, and is available at
// http://www.eclipse.org/legal/epl-v10.html
//
//Contributors:
//    Philippe Charles (pcharles@us.ibm.com) - initial API and implementation

////////////////////////////////////////////////////////////////////////////////

package java
type __JavaParsersym__ struct{
   TK_ClassBodyDeclarationsoptMarker int
   TK_LPGUserActionMarker int
   TK_IntegerLiteral int
   TK_LongLiteral int
   TK_FloatingPointLiteral int
   TK_DoubleLiteral int
   TK_CharacterLiteral int
   TK_StringLiteral int
   TK_MINUS_MINUS int
   TK_OR int
   TK_MINUS int
   TK_MINUS_EQUAL int
   TK_NOT int
   TK_NOT_EQUAL int
   TK_REMAINDER int
   TK_REMAINDER_EQUAL int
   TK_AND int
   TK_AND_AND int
   TK_AND_EQUAL int
   TK_LPAREN int
   TK_RPAREN int
   TK_MULTIPLY int
   TK_MULTIPLY_EQUAL int
   TK_COMMA int
   TK_DOT int
   TK_DIVIDE int
   TK_DIVIDE_EQUAL int
   TK_COLON int
   TK_SEMICOLON int
   TK_QUESTION int
   TK_AT int
   TK_LBRACKET int
   TK_RBRACKET int
   TK_XOR int
   TK_XOR_EQUAL int
   TK_LBRACE int
   TK_OR_OR int
   TK_OR_EQUAL int
   TK_RBRACE int
   TK_TWIDDLE int
   TK_PLUS int
   TK_PLUS_PLUS int
   TK_PLUS_EQUAL int
   TK_LESS int
   TK_LEFT_SHIFT int
   TK_LEFT_SHIFT_EQUAL int
   TK_LESS_EQUAL int
   TK_EQUAL int
   TK_EQUAL_EQUAL int
   TK_GREATER int
   TK_GREATER_EQUAL int
   TK_RIGHT_SHIFT int
   TK_RIGHT_SHIFT_EQUAL int
   TK_UNSIGNED_RIGHT_SHIFT int
   TK_UNSIGNED_RIGHT_SHIFT_EQUAL int
   TK_ELLIPSIS int
   TK_BeginAction int
   TK_EndAction int
   TK_BeginJava int
   TK_EndJava int
   TK_NoAction int
   TK_NullAction int
   TK_BadAction int
   TK_abstract int
   TK_assert int
   TK_boolean int
   TK_break int
   TK_byte int
   TK_case int
   TK_catch int
   TK_char int
   TK_class int
   TK_const int
   TK_continue int
   TK_default int
   TK_do int
   TK_double int
   TK_enum int
   TK_else int
   TK_extends int
   TK_false int
   TK_final int
   TK_finally int
   TK_float int
   TK_for int
   TK_goto int
   TK_if int
   TK_implements int
   TK_import int
   TK_instanceof int
   TK_int int
   TK_interface int
   TK_long int
   TK_native int
   TK_new int
   TK_null int
   TK_package int
   TK_private int
   TK_protected int
   TK_public int
   TK_return int
   TK_short int
   TK_static int
   TK_strictfp int
   TK_super int
   TK_switch int
   TK_synchronized int
   TK_this int
   TK_throw int
   TK_throws int
   TK_transient int
   TK_true int
   TK_try int
   TK_void int
   TK_volatile int
   TK_while int
   TK_EOF_TOKEN int
   TK_IDENTIFIER int
   TK_ERROR_TOKEN int

   OrderedTerminalSymbols []string

   NumTokenKinds int

   IsValidForParser  bool
}
func New__JavaParsersym__() *__JavaParsersym__{
    my := new(__JavaParsersym__)
   my.TK_ClassBodyDeclarationsoptMarker = 102
   my.TK_LPGUserActionMarker = 103
   my.TK_IntegerLiteral = 32
   my.TK_LongLiteral = 33
   my.TK_FloatingPointLiteral = 34
   my.TK_DoubleLiteral = 35
   my.TK_CharacterLiteral = 36
   my.TK_StringLiteral = 37
   my.TK_MINUS_MINUS = 26
   my.TK_OR = 86
   my.TK_MINUS = 46
   my.TK_MINUS_EQUAL = 72
   my.TK_NOT = 48
   my.TK_NOT_EQUAL = 87
   my.TK_REMAINDER = 88
   my.TK_REMAINDER_EQUAL = 73
   my.TK_AND = 68
   my.TK_AND_AND = 89
   my.TK_AND_EQUAL = 74
   my.TK_LPAREN = 3
   my.TK_RPAREN = 20
   my.TK_MULTIPLY = 69
   my.TK_MULTIPLY_EQUAL = 75
   my.TK_COMMA = 43
   my.TK_DOT = 42
   my.TK_DIVIDE = 90
   my.TK_DIVIDE_EQUAL = 76
   my.TK_COLON = 50
   my.TK_SEMICOLON = 4
   my.TK_QUESTION = 91
   my.TK_AT = 1
   my.TK_LBRACKET = 23
   my.TK_RBRACKET = 53
   my.TK_XOR = 92
   my.TK_XOR_EQUAL = 77
   my.TK_LBRACE = 27
   my.TK_OR_OR = 95
   my.TK_OR_EQUAL = 78
   my.TK_RBRACE = 45
   my.TK_TWIDDLE = 49
   my.TK_PLUS = 47
   my.TK_PLUS_PLUS = 28
   my.TK_PLUS_EQUAL = 79
   my.TK_LESS = 24
   my.TK_LEFT_SHIFT = 70
   my.TK_LEFT_SHIFT_EQUAL = 80
   my.TK_LESS_EQUAL = 81
   my.TK_EQUAL = 51
   my.TK_EQUAL_EQUAL = 93
   my.TK_GREATER = 44
   my.TK_GREATER_EQUAL = 112
   my.TK_RIGHT_SHIFT = 113
   my.TK_RIGHT_SHIFT_EQUAL = 114
   my.TK_UNSIGNED_RIGHT_SHIFT = 115
   my.TK_UNSIGNED_RIGHT_SHIFT_EQUAL = 116
   my.TK_ELLIPSIS = 96
   my.TK_BeginAction = 104
   my.TK_EndAction = 105
   my.TK_BeginJava = 106
   my.TK_EndJava = 107
   my.TK_NoAction = 108
   my.TK_NullAction = 109
   my.TK_BadAction = 110
   my.TK_abstract = 17
   my.TK_assert = 57
   my.TK_boolean = 5
   my.TK_break = 58
   my.TK_byte = 6
   my.TK_case = 71
   my.TK_catch = 97
   my.TK_char = 7
   my.TK_class = 31
   my.TK_const = 117
   my.TK_continue = 59
   my.TK_default = 67
   my.TK_do = 60
   my.TK_double = 8
   my.TK_enum = 41
   my.TK_else = 94
   my.TK_extends = 82
   my.TK_false = 38
   my.TK_final = 19
   my.TK_finally = 98
   my.TK_float = 9
   my.TK_for = 61
   my.TK_goto = 118
   my.TK_if = 62
   my.TK_implements = 111
   my.TK_import = 99
   my.TK_instanceof = 83
   my.TK_int = 10
   my.TK_interface = 21
   my.TK_long = 11
   my.TK_native = 84
   my.TK_new = 29
   my.TK_null = 39
   my.TK_package = 100
   my.TK_private = 14
   my.TK_protected = 15
   my.TK_public = 12
   my.TK_return = 63
   my.TK_short = 13
   my.TK_static = 16
   my.TK_strictfp = 18
   my.TK_super = 25
   my.TK_switch = 64
   my.TK_synchronized = 52
   my.TK_this = 30
   my.TK_throw = 65
   my.TK_throws = 101
   my.TK_transient = 54
   my.TK_true = 40
   my.TK_try = 66
   my.TK_void = 22
   my.TK_volatile = 55
   my.TK_while = 56
   my.TK_EOF_TOKEN = 85
   my.TK_IDENTIFIER = 2
   my.TK_ERROR_TOKEN = 119

   my.OrderedTerminalSymbols = []string{
                 "",
                 "AT",
                 "IDENTIFIER",
                 "LPAREN",
                 "SEMICOLON",
                 "boolean",
                 "byte",
                 "char",
                 "double",
                 "float",
                 "int",
                 "long",
                 "public",
                 "short",
                 "private",
                 "protected",
                 "static",
                 "abstract",
                 "strictfp",
                 "final",
                 "RPAREN",
                 "interface",
                 "void",
                 "LBRACKET",
                 "LESS",
                 "super",
                 "MINUS_MINUS",
                 "LBRACE",
                 "PLUS_PLUS",
                 "new",
                 "this",
                 "class",
                 "IntegerLiteral",
                 "LongLiteral",
                 "FloatingPointLiteral",
                 "DoubleLiteral",
                 "CharacterLiteral",
                 "StringLiteral",
                 "false",
                 "null",
                 "true",
                 "enum",
                 "DOT",
                 "COMMA",
                 "GREATER",
                 "RBRACE",
                 "MINUS",
                 "PLUS",
                 "NOT",
                 "TWIDDLE",
                 "COLON",
                 "EQUAL",
                 "synchronized",
                 "RBRACKET",
                 "transient",
                 "volatile",
                 "while",
                 "assert",
                 "break",
                 "continue",
                 "do",
                 "for",
                 "if",
                 "return",
                 "switch",
                 "throw",
                 "try",
                 "default",
                 "AND",
                 "MULTIPLY",
                 "LEFT_SHIFT",
                 "case",
                 "MINUS_EQUAL",
                 "REMAINDER_EQUAL",
                 "AND_EQUAL",
                 "MULTIPLY_EQUAL",
                 "DIVIDE_EQUAL",
                 "XOR_EQUAL",
                 "OR_EQUAL",
                 "PLUS_EQUAL",
                 "LEFT_SHIFT_EQUAL",
                 "LESS_EQUAL",
                 "extends",
                 "instanceof",
                 "native",
                 "EOF_TOKEN",
                 "OR",
                 "NOT_EQUAL",
                 "REMAINDER",
                 "AND_AND",
                 "DIVIDE",
                 "QUESTION",
                 "XOR",
                 "EQUAL_EQUAL",
                 "else",
                 "OR_OR",
                 "ELLIPSIS",
                 "catch",
                 "finally",
                 "import",
                 "package",
                 "throws",
                 "ClassBodyDeclarationsoptMarker",
                 "LPGUserActionMarker",
                 "BeginAction",
                 "EndAction",
                 "BeginJava",
                 "EndJava",
                 "NoAction",
                 "NullAction",
                 "BadAction",
                 "implements",
                 "GREATER_EQUAL",
                 "RIGHT_SHIFT",
                 "RIGHT_SHIFT_EQUAL",
                 "UNSIGNED_RIGHT_SHIFT",
                 "UNSIGNED_RIGHT_SHIFT_EQUAL",
                 "const",
                 "goto",
                 "ERROR_TOKEN",
             }

   my.NumTokenKinds = 119
   my.IsValidForParser = true
   return my
}
var JavaParsersym = New__JavaParsersym__()
