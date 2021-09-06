
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
type __JavaKWLexersym__ struct{
   Char_DollarSign int
   Char_Percent int
   Char__ int
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
   Char_EOF int

   OrderedTerminalSymbols []string

   NumTokenKinds int

   IsValidForParser  bool
}
func New__JavaKWLexersym__() *__JavaKWLexersym__{
    my := new(__JavaKWLexersym__)
   my.Char_DollarSign = 35
   my.Char_Percent = 42
   my.Char__ = 43
   my.Char_a = 3
   my.Char_b = 19
   my.Char_c = 8
   my.Char_d = 14
   my.Char_e = 2
   my.Char_f = 16
   my.Char_g = 23
   my.Char_h = 15
   my.Char_i = 6
   my.Char_j = 28
   my.Char_k = 29
   my.Char_l = 7
   my.Char_m = 27
   my.Char_n = 4
   my.Char_o = 5
   my.Char_p = 24
   my.Char_q = 44
   my.Char_r = 9
   my.Char_s = 10
   my.Char_t = 1
   my.Char_u = 11
   my.Char_v = 20
   my.Char_w = 25
   my.Char_x = 36
   my.Char_y = 26
   my.Char_z = 37
   my.Char_A = 12
   my.Char_B = 38
   my.Char_C = 21
   my.Char_D = 30
   my.Char_E = 31
   my.Char_F = 45
   my.Char_G = 39
   my.Char_H = 46
   my.Char_I = 17
   my.Char_J = 32
   my.Char_K = 47
   my.Char_L = 33
   my.Char_M = 48
   my.Char_N = 13
   my.Char_O = 18
   my.Char_P = 49
   my.Char_Q = 50
   my.Char_R = 51
   my.Char_S = 52
   my.Char_T = 22
   my.Char_U = 40
   my.Char_V = 34
   my.Char_W = 53
   my.Char_X = 54
   my.Char_Y = 55
   my.Char_Z = 56
   my.Char_EOF = 41

   my.OrderedTerminalSymbols = []string{
                 "",
                 "t",
                 "e",
                 "a",
                 "n",
                 "o",
                 "i",
                 "l",
                 "c",
                 "r",
                 "s",
                 "u",
                 "A",
                 "N",
                 "d",
                 "h",
                 "f",
                 "I",
                 "O",
                 "b",
                 "v",
                 "C",
                 "T",
                 "g",
                 "p",
                 "w",
                 "y",
                 "m",
                 "j",
                 "k",
                 "D",
                 "E",
                 "J",
                 "L",
                 "V",
                 "DollarSign",
                 "x",
                 "z",
                 "B",
                 "G",
                 "U",
                 "EOF",
                 "Percent",
                 "_",
                 "q",
                 "F",
                 "H",
                 "K",
                 "M",
                 "P",
                 "Q",
                 "R",
                 "S",
                 "W",
                 "X",
                 "Y",
                 "Z",
             }

   my.NumTokenKinds = 56
   my.IsValidForParser = true
   return my
}
var JavaKWLexersym = New__JavaKWLexersym__()
