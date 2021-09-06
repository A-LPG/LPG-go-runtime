
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


    //#line 58 "KeywordTemplateF.gi


    //#line 63 "KeywordTemplateF.gi

type  JavaKWLexer struct{
    *JavaKWLexerprs
    inputChars []rune
    keywordKind []int
}
func (my *JavaKWLexer)  GetKeywordKinds() []int { return my.keywordKind }

func (my *JavaKWLexer)  Lexer(curtok int,  lasttok int)int{
    var current_kind = my.GetKind(my.inputChars[curtok])
    var    act int 

    for act = my.TAction(JavaKWLexerprs_START_STATE, current_kind);
            act > JavaKWLexerprs_NUM_RULES && act < JavaKWLexerprs_ACCEPT_ACTION;
            act = my.TAction(act, current_kind){
        curtok++
        if curtok > lasttok{
            current_kind = JavaKWLexersym.Char_EOF
        }else{
            current_kind =my.GetKind(my.inputChars[curtok])
        }

    }

    if (act > JavaKWLexerprs_ERROR_ACTION){
        curtok++
        act -= JavaKWLexerprs_ERROR_ACTION
    }
    if act == JavaKWLexerprs_ERROR_ACTION  || curtok <= lasttok  {
        return my.keywordKind[0]
    }else{
        return my.keywordKind[act]
    }
}

func (my *JavaKWLexer)  SetInputChars(inputChars []rune) {
    my.inputChars = inputChars
}

    //#line 10 "KWLexerMapF.gi


func  JavaKWLexerinit_tokenKind() []int {

    var tokenKind=make([]int,128)
    tokenKind['$'] = JavaKWLexersym.Char_DollarSign
    tokenKind['%'] = JavaKWLexersym.Char_Percent
    tokenKind['_'] = JavaKWLexersym.Char__

    tokenKind['a'] = JavaKWLexersym.Char_a
    tokenKind['b'] = JavaKWLexersym.Char_b
    tokenKind['c'] = JavaKWLexersym.Char_c
    tokenKind['d'] = JavaKWLexersym.Char_d
    tokenKind['e'] = JavaKWLexersym.Char_e
    tokenKind['f'] = JavaKWLexersym.Char_f
    tokenKind['g'] = JavaKWLexersym.Char_g
    tokenKind['h'] = JavaKWLexersym.Char_h
    tokenKind['i'] = JavaKWLexersym.Char_i
    tokenKind['j'] = JavaKWLexersym.Char_j
    tokenKind['k'] = JavaKWLexersym.Char_k
    tokenKind['l'] = JavaKWLexersym.Char_l
    tokenKind['m'] = JavaKWLexersym.Char_m
    tokenKind['n'] = JavaKWLexersym.Char_n
    tokenKind['o'] = JavaKWLexersym.Char_o
    tokenKind['p'] = JavaKWLexersym.Char_p
    tokenKind['q'] = JavaKWLexersym.Char_q
    tokenKind['r'] = JavaKWLexersym.Char_r
    tokenKind['s'] = JavaKWLexersym.Char_s
    tokenKind['t'] = JavaKWLexersym.Char_t
    tokenKind['u'] = JavaKWLexersym.Char_u
    tokenKind['v'] = JavaKWLexersym.Char_v
    tokenKind['w'] = JavaKWLexersym.Char_w
    tokenKind['x'] = JavaKWLexersym.Char_x
    tokenKind['y'] = JavaKWLexersym.Char_y
    tokenKind['z'] = JavaKWLexersym.Char_z

    tokenKind['A'] = JavaKWLexersym.Char_A
    tokenKind['B'] = JavaKWLexersym.Char_B
    tokenKind['C'] = JavaKWLexersym.Char_C
    tokenKind['D'] = JavaKWLexersym.Char_D
    tokenKind['E'] = JavaKWLexersym.Char_E
    tokenKind['F'] = JavaKWLexersym.Char_F
    tokenKind['G'] = JavaKWLexersym.Char_G
    tokenKind['H'] = JavaKWLexersym.Char_H
    tokenKind['I'] = JavaKWLexersym.Char_I
    tokenKind['J'] = JavaKWLexersym.Char_J
    tokenKind['K'] = JavaKWLexersym.Char_K
    tokenKind['L'] = JavaKWLexersym.Char_L
    tokenKind['M'] = JavaKWLexersym.Char_M
    tokenKind['N'] = JavaKWLexersym.Char_N
    tokenKind['O'] = JavaKWLexersym.Char_O
    tokenKind['P'] = JavaKWLexersym.Char_P
    tokenKind['Q'] = JavaKWLexersym.Char_Q
    tokenKind['R'] = JavaKWLexersym.Char_R
    tokenKind['S'] = JavaKWLexersym.Char_S
    tokenKind['T'] = JavaKWLexersym.Char_T
    tokenKind['U'] = JavaKWLexersym.Char_U
    tokenKind['V'] = JavaKWLexersym.Char_V
    tokenKind['W'] = JavaKWLexersym.Char_W
    tokenKind['X'] = JavaKWLexersym.Char_X
    tokenKind['Y'] = JavaKWLexersym.Char_Y
    tokenKind['Z'] = JavaKWLexersym.Char_Z

    return tokenKind
}
var JavaKWLexertokenKind  =  JavaKWLexerinit_tokenKind() 
func (my *JavaKWLexer) GetKind(c rune) int {
    if (uint32(c) & 0xFFFFFF80) == 0{
        return JavaKWLexertokenKind[c]
    }else{
        return 0
    }
}

    //#line 105 "KeywordTemplateF.gi

func NewJavaKWLexer(inputChars []rune, identifierKind int)*JavaKWLexer{

    my := new(JavaKWLexer)
    my.JavaKWLexerprs = NewJavaKWLexerprs()
    my.keywordKind = make([]int,88 + 1)
    my.inputChars = inputChars
    my.keywordKind[0] = identifierKind

        //
        // Rule 1:  KeyWord ::= a b s t r a c t
        //
        
        my.keywordKind[1] = (JavaParsersym.TK_abstract)
      
    
        //
        // Rule 2:  KeyWord ::= a s s e r t
        //
        
        my.keywordKind[2] = (JavaParsersym.TK_assert)
      
    
        //
        // Rule 3:  KeyWord ::= b o o l e a n
        //
        
        my.keywordKind[3] = (JavaParsersym.TK_boolean)
      
    
        //
        // Rule 4:  KeyWord ::= b r e a k
        //
        
        my.keywordKind[4] = (JavaParsersym.TK_break)
      
    
        //
        // Rule 5:  KeyWord ::= b y t e
        //
        
        my.keywordKind[5] = (JavaParsersym.TK_byte)
      
    
        //
        // Rule 6:  KeyWord ::= c a s e
        //
        
        my.keywordKind[6] = (JavaParsersym.TK_case)
      
    
        //
        // Rule 7:  KeyWord ::= c a t c h
        //
        
        my.keywordKind[7] = (JavaParsersym.TK_catch)
      
    
        //
        // Rule 8:  KeyWord ::= c h a r
        //
        
        my.keywordKind[8] = (JavaParsersym.TK_char)
      
    
        //
        // Rule 9:  KeyWord ::= c l a s s
        //
        
        my.keywordKind[9] = (JavaParsersym.TK_class)
      
    
        //
        // Rule 10:  KeyWord ::= c o n s t
        //
        
        my.keywordKind[10] = (JavaParsersym.TK_const)
      
    
        //
        // Rule 11:  KeyWord ::= c o n t i n u e
        //
        
        my.keywordKind[11] = (JavaParsersym.TK_continue)
      
    
        //
        // Rule 12:  KeyWord ::= d e f a u l t
        //
        
        my.keywordKind[12] = (JavaParsersym.TK_default)
      
    
        //
        // Rule 13:  KeyWord ::= d o
        //
        
        my.keywordKind[13] = (JavaParsersym.TK_do)
      
    
        //
        // Rule 14:  KeyWord ::= d o u b l e
        //
        
        my.keywordKind[14] = (JavaParsersym.TK_double)
      
    
        //
        // Rule 15:  KeyWord ::= e l s e
        //
        
        my.keywordKind[15] = (JavaParsersym.TK_else)
      
    
        //
        // Rule 16:  KeyWord ::= e n u m
        //
        
        my.keywordKind[16] = (JavaParsersym.TK_enum)
      
    
        //
        // Rule 17:  KeyWord ::= e x t e n d s
        //
        
        my.keywordKind[17] = (JavaParsersym.TK_extends)
      
    
        //
        // Rule 18:  KeyWord ::= f a l s e
        //
        
        my.keywordKind[18] = (JavaParsersym.TK_false)
      
    
        //
        // Rule 19:  KeyWord ::= f i n a l
        //
        
        my.keywordKind[19] = (JavaParsersym.TK_final)
      
    
        //
        // Rule 20:  KeyWord ::= f i n a l l y
        //
        
        my.keywordKind[20] = (JavaParsersym.TK_finally)
      
    
        //
        // Rule 21:  KeyWord ::= f l o a t
        //
        
        my.keywordKind[21] = (JavaParsersym.TK_float)
      
    
        //
        // Rule 22:  KeyWord ::= f o r
        //
        
        my.keywordKind[22] = (JavaParsersym.TK_for)
      
    
        //
        // Rule 23:  KeyWord ::= g o t o
        //
        
        my.keywordKind[23] = (JavaParsersym.TK_goto)
      
    
        //
        // Rule 24:  KeyWord ::= i f
        //
        
        my.keywordKind[24] = (JavaParsersym.TK_if)
      
    
        //
        // Rule 25:  KeyWord ::= i m p l e m e n t s
        //
        
        my.keywordKind[25] = (JavaParsersym.TK_implements)
      
    
        //
        // Rule 26:  KeyWord ::= i m p o r t
        //
        
        my.keywordKind[26] = (JavaParsersym.TK_import)
      
    
        //
        // Rule 27:  KeyWord ::= i n s t a n c e o f
        //
        
        my.keywordKind[27] = (JavaParsersym.TK_instanceof)
      
    
        //
        // Rule 28:  KeyWord ::= i n t
        //
        
        my.keywordKind[28] = (JavaParsersym.TK_int)
      
    
        //
        // Rule 29:  KeyWord ::= i n t e r f a c e
        //
        
        my.keywordKind[29] = (JavaParsersym.TK_interface)
      
    
        //
        // Rule 30:  KeyWord ::= l o n g
        //
        
        my.keywordKind[30] = (JavaParsersym.TK_long)
      
    
        //
        // Rule 31:  KeyWord ::= n a t i v e
        //
        
        my.keywordKind[31] = (JavaParsersym.TK_native)
      
    
        //
        // Rule 32:  KeyWord ::= n e w
        //
        
        my.keywordKind[32] = (JavaParsersym.TK_new)
      
    
        //
        // Rule 33:  KeyWord ::= n u l l
        //
        
        my.keywordKind[33] = (JavaParsersym.TK_null)
      
    
        //
        // Rule 34:  KeyWord ::= p a c k a g e
        //
        
        my.keywordKind[34] = (JavaParsersym.TK_package)
      
    
        //
        // Rule 35:  KeyWord ::= p r i v a t e
        //
        
        my.keywordKind[35] = (JavaParsersym.TK_private)
      
    
        //
        // Rule 36:  KeyWord ::= p r o t e c t e d
        //
        
        my.keywordKind[36] = (JavaParsersym.TK_protected)
      
    
        //
        // Rule 37:  KeyWord ::= p u b l i c
        //
        
        my.keywordKind[37] = (JavaParsersym.TK_public)
      
    
        //
        // Rule 38:  KeyWord ::= r e t u r n
        //
        
        my.keywordKind[38] = (JavaParsersym.TK_return)
      
    
        //
        // Rule 39:  KeyWord ::= s h o r t
        //
        
        my.keywordKind[39] = (JavaParsersym.TK_short)
      
    
        //
        // Rule 40:  KeyWord ::= s t a t i c
        //
        
        my.keywordKind[40] = (JavaParsersym.TK_static)
      
    
        //
        // Rule 41:  KeyWord ::= s t r i c t f p
        //
        
        my.keywordKind[41] = (JavaParsersym.TK_strictfp)
      
    
        //
        // Rule 42:  KeyWord ::= s u p e r
        //
        
        my.keywordKind[42] = (JavaParsersym.TK_super)
      
    
        //
        // Rule 43:  KeyWord ::= s w i t c h
        //
        
        my.keywordKind[43] = (JavaParsersym.TK_switch)
      
    
        //
        // Rule 44:  KeyWord ::= s y n c h r o n i z e d
        //
        
        my.keywordKind[44] = (JavaParsersym.TK_synchronized)
      
    
        //
        // Rule 45:  KeyWord ::= t h i s
        //
        
        my.keywordKind[45] = (JavaParsersym.TK_this)
      
    
        //
        // Rule 46:  KeyWord ::= t h r o w
        //
        
        my.keywordKind[46] = (JavaParsersym.TK_throw)
      
    
        //
        // Rule 47:  KeyWord ::= t h r o w s
        //
        
        my.keywordKind[47] = (JavaParsersym.TK_throws)
      
    
        //
        // Rule 48:  KeyWord ::= t r a n s i e n t
        //
        
        my.keywordKind[48] = (JavaParsersym.TK_transient)
      
    
        //
        // Rule 49:  KeyWord ::= t r u e
        //
        
        my.keywordKind[49] = (JavaParsersym.TK_true)
      
    
        //
        // Rule 50:  KeyWord ::= t r y
        //
        
        my.keywordKind[50] = (JavaParsersym.TK_try)
      
    
        //
        // Rule 51:  KeyWord ::= v o i d
        //
        
        my.keywordKind[51] = (JavaParsersym.TK_void)
      
    
        //
        // Rule 52:  KeyWord ::= v o l a t i l e
        //
        
        my.keywordKind[52] = (JavaParsersym.TK_volatile)
      
    
        //
        // Rule 53:  KeyWord ::= w h i l e
        //
        
        my.keywordKind[53] = (JavaParsersym.TK_while)
      
    
        //
        // Rule 54:  KeyWord ::= $ bB eE gG iI nN aA cC tT iI oO nN
        //
        
        my.keywordKind[54] = (JavaParsersym.TK_BeginAction)
      
    
        //
        // Rule 55:  KeyWord ::= $ bB eE gG iI nN jJ aA vV aA
        //
        
        my.keywordKind[55] = (JavaParsersym.TK_BeginJava)
      
    
        //
        // Rule 56:  KeyWord ::= $ eE nN dD aA cC tT iI oO nN
        //
        
        my.keywordKind[56] = (JavaParsersym.TK_EndAction)
      
    
        //
        // Rule 57:  KeyWord ::= $ eE nN dD jJ aA vV aA
        //
        
        my.keywordKind[57] = (JavaParsersym.TK_EndJava)
      
    
        //
        // Rule 58:  KeyWord ::= $ nN oO aA cC tT iI oO nN
        //
        
        my.keywordKind[58] = (JavaParsersym.TK_NoAction)
      
    
        //
        // Rule 59:  KeyWord ::= $ nN uU lL lL aA cC tT iI oO nN
        //
        
        my.keywordKind[59] = (JavaParsersym.TK_NullAction)
      
    
        //
        // Rule 60:  KeyWord ::= $ bB aA dD aA cC tT iI oO nN
        //
        
        my.keywordKind[60] = (JavaParsersym.TK_BadAction)
      
    
    //#line 117 "KeywordTemplateF.gi

    var i int = 0
    for ;i < len(my.keywordKind); i++{
        if my.keywordKind[i] == 0 {
            my.keywordKind[i] = identifierKind
        }
    }
    return my
}

