
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

package lpg
type __LPGLexersym__ struct{
   Char_CtlCharNotWS int
   Char_LF int
   Char_CR int
   Char_HT int
   Char_FF int
   Char_a int
   Char_b int
   Char_c int
   Char_d int
   Char_e int
   Char_f int
   Char_g int
   Char_h int
   Char_i int
   Char_j int
   Char_k int
   Char_l int
   Char_m int
   Char_n int
   Char_o int
   Char_p int
   Char_q int
   Char_r int
   Char_s int
   Char_t int
   Char_u int
   Char_v int
   Char_w int
   Char_x int
   Char_y int
   Char_z int
   Char__ int
   Char_A int
   Char_B int
   Char_C int
   Char_D int
   Char_E int
   Char_F int
   Char_G int
   Char_H int
   Char_I int
   Char_J int
   Char_K int
   Char_L int
   Char_M int
   Char_N int
   Char_O int
   Char_P int
   Char_Q int
   Char_R int
   Char_S int
   Char_T int
   Char_U int
   Char_V int
   Char_W int
   Char_X int
   Char_Y int
   Char_Z int
   Char_0 int
   Char_1 int
   Char_2 int
   Char_3 int
   Char_4 int
   Char_5 int
   Char_6 int
   Char_7 int
   Char_8 int
   Char_9 int
   Char_AfterASCII int
   Char_Space int
   Char_DoubleQuote int
   Char_SingleQuote int
   Char_Percent int
   Char_VerticalBar int
   Char_Exclamation int
   Char_AtSign int
   Char_BackQuote int
   Char_Tilde int
   Char_Sharp int
   Char_DollarSign int
   Char_Ampersand int
   Char_Caret int
   Char_Colon int
   Char_SemiColon int
   Char_BackSlash int
   Char_LeftBrace int
   Char_RightBrace int
   Char_LeftBracket int
   Char_RightBracket int
   Char_QuestionMark int
   Char_Comma int
   Char_Dot int
   Char_LessThan int
   Char_GreaterThan int
   Char_Plus int
   Char_Minus int
   Char_Slash int
   Char_Star int
   Char_LeftParen int
   Char_RightParen int
   Char_Equal int
   Char_EOF int

   OrderedTerminalSymbols []string

   NumTokenKinds int

   IsValidForParser  bool
}
func New__LPGLexersym__() *__LPGLexersym__{
    my := new(__LPGLexersym__)
   my.Char_CtlCharNotWS = 102
   my.Char_LF = 5
   my.Char_CR = 6
   my.Char_HT = 1
   my.Char_FF = 2
   my.Char_a = 15
   my.Char_b = 40
   my.Char_c = 24
   my.Char_d = 30
   my.Char_e = 7
   my.Char_f = 31
   my.Char_g = 42
   my.Char_h = 50
   my.Char_i = 13
   my.Char_j = 70
   my.Char_k = 44
   my.Char_l = 18
   my.Char_m = 34
   my.Char_n = 28
   my.Char_o = 22
   my.Char_p = 36
   my.Char_q = 56
   my.Char_r = 11
   my.Char_s = 20
   my.Char_t = 9
   my.Char_u = 38
   my.Char_v = 52
   my.Char_w = 53
   my.Char_x = 48
   my.Char_y = 45
   my.Char_z = 68
   my.Char__ = 26
   my.Char_A = 16
   my.Char_B = 41
   my.Char_C = 25
   my.Char_D = 32
   my.Char_E = 8
   my.Char_F = 33
   my.Char_G = 43
   my.Char_H = 51
   my.Char_I = 14
   my.Char_J = 71
   my.Char_K = 46
   my.Char_L = 19
   my.Char_M = 35
   my.Char_N = 29
   my.Char_O = 23
   my.Char_P = 37
   my.Char_Q = 57
   my.Char_R = 12
   my.Char_S = 21
   my.Char_T = 10
   my.Char_U = 39
   my.Char_V = 54
   my.Char_W = 55
   my.Char_X = 49
   my.Char_Y = 47
   my.Char_Z = 69
   my.Char_0 = 58
   my.Char_1 = 59
   my.Char_2 = 60
   my.Char_3 = 61
   my.Char_4 = 62
   my.Char_5 = 63
   my.Char_6 = 64
   my.Char_7 = 65
   my.Char_8 = 66
   my.Char_9 = 67
   my.Char_AfterASCII = 72
   my.Char_Space = 3
   my.Char_DoubleQuote = 97
   my.Char_SingleQuote = 98
   my.Char_Percent = 74
   my.Char_VerticalBar = 76
   my.Char_Exclamation = 77
   my.Char_AtSign = 78
   my.Char_BackQuote = 79
   my.Char_Tilde = 80
   my.Char_Sharp = 92
   my.Char_DollarSign = 100
   my.Char_Ampersand = 81
   my.Char_Caret = 82
   my.Char_Colon = 83
   my.Char_SemiColon = 84
   my.Char_BackSlash = 85
   my.Char_LeftBrace = 86
   my.Char_RightBrace = 87
   my.Char_LeftBracket = 93
   my.Char_RightBracket = 94
   my.Char_QuestionMark = 73
   my.Char_Comma = 4
   my.Char_Dot = 88
   my.Char_LessThan = 99
   my.Char_GreaterThan = 95
   my.Char_Plus = 89
   my.Char_Minus = 27
   my.Char_Slash = 90
   my.Char_Star = 91
   my.Char_LeftParen = 96
   my.Char_RightParen = 75
   my.Char_Equal = 17
   my.Char_EOF = 101

   my.OrderedTerminalSymbols = []string{
                 "",
                 "HT",
                 "FF",
                 "Space",
                 "Comma",
                 "LF",
                 "CR",
                 "e",
                 "E",
                 "t",
                 "T",
                 "r",
                 "R",
                 "i",
                 "I",
                 "a",
                 "A",
                 "Equal",
                 "l",
                 "L",
                 "s",
                 "S",
                 "o",
                 "O",
                 "c",
                 "C",
                 "_",
                 "Minus",
                 "n",
                 "N",
                 "d",
                 "f",
                 "D",
                 "F",
                 "m",
                 "M",
                 "p",
                 "P",
                 "u",
                 "U",
                 "b",
                 "B",
                 "g",
                 "G",
                 "k",
                 "y",
                 "K",
                 "Y",
                 "x",
                 "X",
                 "h",
                 "H",
                 "v",
                 "w",
                 "V",
                 "W",
                 "q",
                 "Q",
                 "0",
                 "1",
                 "2",
                 "3",
                 "4",
                 "5",
                 "6",
                 "7",
                 "8",
                 "9",
                 "z",
                 "Z",
                 "j",
                 "J",
                 "AfterASCII",
                 "QuestionMark",
                 "Percent",
                 "RightParen",
                 "VerticalBar",
                 "Exclamation",
                 "AtSign",
                 "BackQuote",
                 "Tilde",
                 "Ampersand",
                 "Caret",
                 "Colon",
                 "SemiColon",
                 "BackSlash",
                 "LeftBrace",
                 "RightBrace",
                 "Dot",
                 "Plus",
                 "Slash",
                 "Star",
                 "Sharp",
                 "LeftBracket",
                 "RightBracket",
                 "GreaterThan",
                 "LeftParen",
                 "DoubleQuote",
                 "SingleQuote",
                 "LessThan",
                 "DollarSign",
                 "EOF",
                 "CtlCharNotWS",
             }

   my.NumTokenKinds = 102
   my.IsValidForParser = true
   return my
}
var LPGLexersym = New__LPGLexersym__()
