
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
type __JavaLexersym__ struct{
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
func New__JavaLexersym__() *__JavaLexersym__{
    my := new(__JavaLexersym__)
   my.Char_CtlCharNotWS = 102
   my.Char_LF = 100
   my.Char_CR = 101
   my.Char_HT = 37
   my.Char_FF = 38
   my.Char_a = 19
   my.Char_b = 15
   my.Char_c = 20
   my.Char_d = 12
   my.Char_e = 16
   my.Char_f = 11
   my.Char_g = 39
   my.Char_h = 40
   my.Char_i = 41
   my.Char_j = 42
   my.Char_k = 43
   my.Char_l = 25
   my.Char_m = 44
   my.Char_n = 26
   my.Char_o = 45
   my.Char_p = 46
   my.Char_q = 47
   my.Char_r = 27
   my.Char_s = 48
   my.Char_t = 28
   my.Char_u = 29
   my.Char_v = 49
   my.Char_w = 50
   my.Char_x = 32
   my.Char_y = 51
   my.Char_z = 52
   my.Char__ = 53
   my.Char_A = 21
   my.Char_B = 22
   my.Char_C = 23
   my.Char_D = 13
   my.Char_E = 17
   my.Char_F = 14
   my.Char_G = 54
   my.Char_H = 55
   my.Char_I = 56
   my.Char_J = 57
   my.Char_K = 58
   my.Char_L = 30
   my.Char_M = 59
   my.Char_N = 60
   my.Char_O = 61
   my.Char_P = 62
   my.Char_Q = 63
   my.Char_R = 64
   my.Char_S = 65
   my.Char_T = 66
   my.Char_U = 67
   my.Char_V = 68
   my.Char_W = 69
   my.Char_X = 33
   my.Char_Y = 70
   my.Char_Z = 71
   my.Char_0 = 1
   my.Char_1 = 2
   my.Char_2 = 3
   my.Char_3 = 4
   my.Char_4 = 5
   my.Char_5 = 6
   my.Char_6 = 7
   my.Char_7 = 8
   my.Char_8 = 9
   my.Char_9 = 10
   my.Char_AfterASCII = 72
   my.Char_Space = 73
   my.Char_DoubleQuote = 34
   my.Char_SingleQuote = 24
   my.Char_Percent = 81
   my.Char_VerticalBar = 74
   my.Char_Exclamation = 82
   my.Char_AtSign = 83
   my.Char_BackQuote = 97
   my.Char_Tilde = 84
   my.Char_Sharp = 98
   my.Char_DollarSign = 75
   my.Char_Ampersand = 76
   my.Char_Caret = 85
   my.Char_Colon = 86
   my.Char_SemiColon = 87
   my.Char_BackSlash = 77
   my.Char_LeftBrace = 88
   my.Char_RightBrace = 89
   my.Char_LeftBracket = 90
   my.Char_RightBracket = 91
   my.Char_QuestionMark = 92
   my.Char_Comma = 93
   my.Char_Dot = 31
   my.Char_LessThan = 78
   my.Char_GreaterThan = 94
   my.Char_Plus = 35
   my.Char_Minus = 36
   my.Char_Slash = 79
   my.Char_Star = 80
   my.Char_LeftParen = 95
   my.Char_RightParen = 96
   my.Char_Equal = 18
   my.Char_EOF = 99

   my.OrderedTerminalSymbols = []string{
                 "",
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
                 "f",
                 "d",
                 "D",
                 "F",
                 "b",
                 "e",
                 "E",
                 "Equal",
                 "a",
                 "c",
                 "A",
                 "B",
                 "C",
                 "SingleQuote",
                 "l",
                 "n",
                 "r",
                 "t",
                 "u",
                 "L",
                 "Dot",
                 "x",
                 "X",
                 "DoubleQuote",
                 "Plus",
                 "Minus",
                 "HT",
                 "FF",
                 "g",
                 "h",
                 "i",
                 "j",
                 "k",
                 "m",
                 "o",
                 "p",
                 "q",
                 "s",
                 "v",
                 "w",
                 "y",
                 "z",
                 "_",
                 "G",
                 "H",
                 "I",
                 "J",
                 "K",
                 "M",
                 "N",
                 "O",
                 "P",
                 "Q",
                 "R",
                 "S",
                 "T",
                 "U",
                 "V",
                 "W",
                 "Y",
                 "Z",
                 "AfterASCII",
                 "Space",
                 "VerticalBar",
                 "DollarSign",
                 "Ampersand",
                 "BackSlash",
                 "LessThan",
                 "Slash",
                 "Star",
                 "Percent",
                 "Exclamation",
                 "AtSign",
                 "Tilde",
                 "Caret",
                 "Colon",
                 "SemiColon",
                 "LeftBrace",
                 "RightBrace",
                 "LeftBracket",
                 "RightBracket",
                 "QuestionMark",
                 "Comma",
                 "GreaterThan",
                 "LeftParen",
                 "RightParen",
                 "BackQuote",
                 "Sharp",
                 "EOF",
                 "LF",
                 "CR",
                 "CtlCharNotWS",
             }

   my.NumTokenKinds = 102
   my.IsValidForParser = true
   return my
}
var JavaLexersym = New__JavaLexersym__()
