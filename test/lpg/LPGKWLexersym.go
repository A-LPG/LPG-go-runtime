
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
type __LPGKWLexersym__ struct{
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
   Char_EOF int

   OrderedTerminalSymbols []string

   NumTokenKinds int

   IsValidForParser  bool
}
func New__LPGKWLexersym__() *__LPGKWLexersym__{
    my := new(__LPGKWLexersym__)
   my.Char_DollarSign = 20
   my.Char_Percent = 21
   my.Char__ = 28
   my.Char_a = 8
   my.Char_b = 17
   my.Char_c = 14
   my.Char_d = 9
   my.Char_e = 1
   my.Char_f = 15
   my.Char_g = 22
   my.Char_h = 23
   my.Char_i = 6
   my.Char_j = 24
   my.Char_k = 18
   my.Char_l = 7
   my.Char_m = 12
   my.Char_n = 10
   my.Char_o = 4
   my.Char_p = 11
   my.Char_q = 29
   my.Char_r = 3
   my.Char_s = 2
   my.Char_t = 5
   my.Char_u = 16
   my.Char_v = 25
   my.Char_w = 19
   my.Char_x = 26
   my.Char_y = 13
   my.Char_z = 30
   my.Char_EOF = 27

   my.OrderedTerminalSymbols = []string{
                 "",
                 "e",
                 "s",
                 "r",
                 "o",
                 "t",
                 "i",
                 "l",
                 "a",
                 "d",
                 "n",
                 "p",
                 "m",
                 "y",
                 "c",
                 "f",
                 "u",
                 "b",
                 "k",
                 "w",
                 "DollarSign",
                 "Percent",
                 "g",
                 "h",
                 "j",
                 "v",
                 "x",
                 "EOF",
                 "_",
                 "q",
                 "z",
             }

   my.NumTokenKinds = 30
   my.IsValidForParser = true
   return my
}
var LPGKWLexersym = New__LPGKWLexersym__()
