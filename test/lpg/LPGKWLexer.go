
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


    //#line 58 "KeywordTemplateF.gi


    //#line 63 "KeywordTemplateF.gi

type  LPGKWLexer struct{
    *LPGKWLexerprs
    inputChars []rune
    keywordKind []int
}
func (my *LPGKWLexer)  GetKeywordKinds() []int { return my.keywordKind }

func (my *LPGKWLexer)  Lexer(curtok int,  lasttok int)int{
    var current_kind = my.GetKind(my.inputChars[curtok])
    var    act int 

    for act = my.TAction(LPGKWLexerprs_START_STATE, current_kind);
            act > LPGKWLexerprs_NUM_RULES && act < LPGKWLexerprs_ACCEPT_ACTION;
            act = my.TAction(act, current_kind){
        curtok++
        if curtok > lasttok{
            current_kind = LPGKWLexersym.Char_EOF
        }else{
            current_kind =my.GetKind(my.inputChars[curtok])
        }

    }

    if (act > LPGKWLexerprs_ERROR_ACTION){
        curtok++
        act -= LPGKWLexerprs_ERROR_ACTION
    }
    if act == LPGKWLexerprs_ERROR_ACTION  || curtok <= lasttok  {
        return my.keywordKind[0]
    }else{
        return my.keywordKind[act]
    }
}

func (my *LPGKWLexer)  SetInputChars(inputChars []rune) {
    my.inputChars = inputChars
}

    //#line 9 "KWLexerFoldedCaseMapF.gi


// Each upper case letter is mapped into its corresponding
// lower case counterpart. For example, if an 'A' appears
// in the input, it is mapped into LPGKWLexersym.Char_a just
// like 'a'.
//
func  LPGKWLexerinit_tokenKind() []int {

    var tokenKind=make([]int,128)
    tokenKind['$'] = LPGKWLexersym.Char_DollarSign
    tokenKind['%'] = LPGKWLexersym.Char_Percent
    tokenKind['_'] = LPGKWLexersym.Char__

    tokenKind['a'] = LPGKWLexersym.Char_a
    tokenKind['b'] = LPGKWLexersym.Char_b
    tokenKind['c'] = LPGKWLexersym.Char_c
    tokenKind['d'] = LPGKWLexersym.Char_d
    tokenKind['e'] = LPGKWLexersym.Char_e
    tokenKind['f'] = LPGKWLexersym.Char_f
    tokenKind['g'] = LPGKWLexersym.Char_g
    tokenKind['h'] = LPGKWLexersym.Char_h
    tokenKind['i'] = LPGKWLexersym.Char_i
    tokenKind['j'] = LPGKWLexersym.Char_j
    tokenKind['k'] = LPGKWLexersym.Char_k
    tokenKind['l'] = LPGKWLexersym.Char_l
    tokenKind['m'] = LPGKWLexersym.Char_m
    tokenKind['n'] = LPGKWLexersym.Char_n
    tokenKind['o'] = LPGKWLexersym.Char_o
    tokenKind['p'] = LPGKWLexersym.Char_p
    tokenKind['q'] = LPGKWLexersym.Char_q
    tokenKind['r'] = LPGKWLexersym.Char_r
    tokenKind['s'] = LPGKWLexersym.Char_s
    tokenKind['t'] = LPGKWLexersym.Char_t
    tokenKind['u'] = LPGKWLexersym.Char_u
    tokenKind['v'] = LPGKWLexersym.Char_v
    tokenKind['w'] = LPGKWLexersym.Char_w
    tokenKind['x'] = LPGKWLexersym.Char_x
    tokenKind['y'] = LPGKWLexersym.Char_y
    tokenKind['z'] = LPGKWLexersym.Char_z

    tokenKind['A'] = LPGKWLexersym.Char_a
    tokenKind['B'] = LPGKWLexersym.Char_b
    tokenKind['C'] = LPGKWLexersym.Char_c
    tokenKind['D'] = LPGKWLexersym.Char_d
    tokenKind['E'] = LPGKWLexersym.Char_e
    tokenKind['F'] = LPGKWLexersym.Char_f
    tokenKind['G'] = LPGKWLexersym.Char_g
    tokenKind['H'] = LPGKWLexersym.Char_h
    tokenKind['I'] = LPGKWLexersym.Char_i
    tokenKind['J'] = LPGKWLexersym.Char_j
    tokenKind['K'] = LPGKWLexersym.Char_k
    tokenKind['L'] = LPGKWLexersym.Char_l
    tokenKind['M'] = LPGKWLexersym.Char_m
    tokenKind['N'] = LPGKWLexersym.Char_n
    tokenKind['O'] = LPGKWLexersym.Char_o
    tokenKind['P'] = LPGKWLexersym.Char_p
    tokenKind['Q'] = LPGKWLexersym.Char_q
    tokenKind['R'] = LPGKWLexersym.Char_r
    tokenKind['S'] = LPGKWLexersym.Char_s
    tokenKind['T'] = LPGKWLexersym.Char_t
    tokenKind['U'] = LPGKWLexersym.Char_u
    tokenKind['V'] = LPGKWLexersym.Char_v
    tokenKind['W'] = LPGKWLexersym.Char_w
    tokenKind['X'] = LPGKWLexersym.Char_x
    tokenKind['Y'] = LPGKWLexersym.Char_y
    tokenKind['Z'] = LPGKWLexersym.Char_z
    return tokenKind
}
var LPGKWLexertokenKind  =  LPGKWLexerinit_tokenKind() 
func (my *LPGKWLexer) GetKind(c rune) int {
    if c < 128{
        return LPGKWLexertokenKind[c]
    }else{
        return 0
    }
}

    //#line 105 "KeywordTemplateF.gi

func NewLPGKWLexer(inputChars []rune, identifierKind int)*LPGKWLexer{

    my := new(LPGKWLexer)
    my.LPGKWLexerprs = NewLPGKWLexerprs()
    my.keywordKind = make([]int,29 + 1)
    my.inputChars = inputChars
    my.keywordKind[0] = identifierKind

        //
        // Rule 1:  Keyword ::= KeyPrefix a l i a s
        //
        
        my.keywordKind[1] = (LPGParsersym.TK_ALIAS_KEY)
      
    
        //
        // Rule 2:  Keyword ::= KeyPrefix a s t
        //
        
        my.keywordKind[2] = (LPGParsersym.TK_AST_KEY)
      
    
        //
        // Rule 3:  Keyword ::= KeyPrefix d e f i n e
        //
        
        my.keywordKind[3] = (LPGParsersym.TK_DEFINE_KEY)
      
    
        //
        // Rule 4:  Keyword ::= KeyPrefix d i s j o i n t p r e d e c e s s o r s e t s
        //
        
        my.keywordKind[4] = (LPGParsersym.TK_DISJOINTPREDECESSORSETS_KEY)
      
    
        //
        // Rule 5:  Keyword ::= KeyPrefix d r o p r u l e s
        //
        
        my.keywordKind[5] = (LPGParsersym.TK_DROPRULES_KEY)
      
    
        //
        // Rule 6:  Keyword ::= KeyPrefix d r o p s y m b o l s
        //
        
        my.keywordKind[6] = (LPGParsersym.TK_DROPSYMBOLS_KEY)
      
    
        //
        // Rule 7:  Keyword ::= KeyPrefix e m p t y
        //
        
        my.keywordKind[7] = (LPGParsersym.TK_EMPTY_KEY)
      
    
        //
        // Rule 8:  Keyword ::= KeyPrefix e n d
        //
        
        my.keywordKind[8] = (LPGParsersym.TK_END_KEY)
      
    
        //
        // Rule 9:  Keyword ::= KeyPrefix e r r o r
        //
        
        my.keywordKind[9] = (LPGParsersym.TK_ERROR_KEY)
      
    
        //
        // Rule 10:  Keyword ::= KeyPrefix e o l
        //
        
        my.keywordKind[10] = (LPGParsersym.TK_EOL_KEY)
      
    
        //
        // Rule 11:  Keyword ::= KeyPrefix e o f
        //
        
        my.keywordKind[11] = (LPGParsersym.TK_EOF_KEY)
      
    
        //
        // Rule 12:  Keyword ::= KeyPrefix e x p o r t
        //
        
        my.keywordKind[12] = (LPGParsersym.TK_EXPORT_KEY)
      
    
        //
        // Rule 13:  Keyword ::= KeyPrefix g l o b a l s
        //
        
        my.keywordKind[13] = (LPGParsersym.TK_GLOBALS_KEY)
      
    
        //
        // Rule 14:  Keyword ::= KeyPrefix h e a d e r s
        //
        
        my.keywordKind[14] = (LPGParsersym.TK_HEADERS_KEY)
      
    
        //
        // Rule 15:  Keyword ::= KeyPrefix i d e n t i f i e r
        //
        
        my.keywordKind[15] = (LPGParsersym.TK_IDENTIFIER_KEY)
      
    
        //
        // Rule 16:  Keyword ::= KeyPrefix i m p o r t
        //
        
        my.keywordKind[16] = (LPGParsersym.TK_IMPORT_KEY)
      
    
        //
        // Rule 17:  Keyword ::= KeyPrefix i n c l u d e
        //
        
        my.keywordKind[17] = (LPGParsersym.TK_INCLUDE_KEY)
      
    
        //
        // Rule 18:  Keyword ::= KeyPrefix k e y w o r d s
        //
        
        my.keywordKind[18] = (LPGParsersym.TK_KEYWORDS_KEY)
      
    
        //
        // Rule 19:  Keyword ::= KeyPrefix s o f t k e y w o r d s
        //
        
        my.keywordKind[19] = (LPGParsersym.TK_SOFT_KEYWORDS_KEY)
      
    
        //
        // Rule 20:  Keyword ::= KeyPrefix n a m e s
        //
        
        my.keywordKind[20] = (LPGParsersym.TK_NAMES_KEY)
      
    
        //
        // Rule 21:  Keyword ::= KeyPrefix n o t i c e
        //
        
        my.keywordKind[21] = (LPGParsersym.TK_NOTICE_KEY)
      
    
        //
        // Rule 22:  Keyword ::= KeyPrefix t e r m i n a l s
        //
        
        my.keywordKind[22] = (LPGParsersym.TK_TERMINALS_KEY)
      
    
        //
        // Rule 23:  Keyword ::= KeyPrefix r e c o v e r
        //
        
        my.keywordKind[23] = (LPGParsersym.TK_RECOVER_KEY)
      
    
        //
        // Rule 24:  Keyword ::= KeyPrefix r u l e s
        //
        
        my.keywordKind[24] = (LPGParsersym.TK_RULES_KEY)
      
    
        //
        // Rule 25:  Keyword ::= KeyPrefix s t a r t
        //
        
        my.keywordKind[25] = (LPGParsersym.TK_START_KEY)
      
    
        //
        // Rule 26:  Keyword ::= KeyPrefix t r a i l e r s
        //
        
        my.keywordKind[26] = (LPGParsersym.TK_TRAILERS_KEY)
      
    
        //
        // Rule 27:  Keyword ::= KeyPrefix t y p e s
        //
        
        my.keywordKind[27] = (LPGParsersym.TK_TYPES_KEY)
      
    
    //#line 117 "KeywordTemplateF.gi

    var i int = 0
    for ;i < len(my.keywordKind); i++{
        if my.keywordKind[i] == 0 {
            my.keywordKind[i] = identifierKind
        }
    }
    return my
}

