
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


    //#line 122 "LexerTemplateF.gi

import ."github.com/A-LPG/LPG-go-runtime/lpg2"

    //#line 128 "LexerTemplateF.gi

type  JavaLexer struct{

    kwLexer  *JavaKWLexer 
    PrintTokens  bool
    lexParser  *LexParser
    lexStream *JavaLexerLpgLexStream 
    prs     ParseTable
}
func NewJavaLexer(filename string,  tab int ,input_chars []rune) *JavaLexer {
    my := new(JavaLexer)
    my.prs =  NewJavaLexerprs()
    my.lexParser   =  NewLexParser()
    my.PrintTokens = false
    var e error
    my.lexStream, e =  NewJavaLexerLpgLexStream(filename,input_chars, tab)
    if e != nil{
        return nil
    }
    my.lexParser.Reset( my.lexStream, my.prs,  my)
    my.ResetKeywordLexer()
    return my
}
    func (my *JavaLexer)  GetParseTable()ParseTable {
     return my.prs 
    }

    func (my *JavaLexer)  GetParser() *LexParser{ 
        return my.lexParser 
    }

    func (my *JavaLexer)  GetToken(i int)int { 
    return my.lexParser.GetToken(i)
     }
    func (my *JavaLexer)  GetRhsFirstTokenIndex(i int)int {
     return my.lexParser.GetFirstTokenAt(i) 
     }
    func (my *JavaLexer)  GetRhsLastTokenIndex(i int) int{
     return my.lexParser.GetLastTokenAt(i) 
     }

    func (my *JavaLexer)  GetLeftSpan()int { 
    return my.lexParser.GetToken(1) 
    }
    func (my *JavaLexer)  GetRightSpan()int {
     return my.lexParser.GetLastToken() 
     }

    func (my *JavaLexer)  ResetKeywordLexer(){
        if my.kwLexer == nil {
            my.kwLexer = NewJavaKWLexer(my.lexStream.GetInputChars(), JavaParsersym.TK_IDENTIFIER)
        }else {
            my.kwLexer.SetInputChars(my.lexStream.GetInputChars())
        }
    }

    func (my *JavaLexer)  Reset(filename string,tab int,input_chars []rune) error{
        var e error
        my.lexStream ,e = NewJavaLexerLpgLexStream(filename,input_chars, tab)
        if e != nil{
            return e
        }
        my.lexParser.Reset(my.lexStream, my.prs, my)
        my.ResetKeywordLexer()
        return nil
    }
    

    func (my *JavaLexer)  GetILexStream()ILexStream { 
        return my.lexStream
    }



    func (my *JavaLexer)  InitializeLexer(prsStream IPrsStream, start_offSet int, end_offSet int) error{
        if my.lexStream.GetInputChars() == nil{
            return NewNullPointerException("LexStream was not initialized")
        }
        my.lexStream.SetPrsStream(prsStream)
        prsStream.SetLexStream(my.lexStream)
        prsStream.MakeToken(start_offSet, end_offSet, 0) // Token list must start with a bad token
        return nil
    }

    func (my *JavaLexer)  AddEOF(prsStream IPrsStream, end_offSet int) {
        prsStream.MakeToken(end_offSet, end_offSet, JavaParsersym.TK_EOF_TOKEN) // and end with the end of file token
        prsStream.SetStreamLength(prsStream.GetSize())
    }

    func (my *JavaLexer) LexerWithPosition(prsStream IPrsStream , start_offSet int , end_offSet int, monitor Monitor) error{
    
        if start_offSet <= 1{
            err := my.InitializeLexer(prsStream, 0, -1)
            if err != nil {
                return err
            }
        }else {
            err := my.InitializeLexer(prsStream, start_offSet - 1, start_offSet - 1)
            if err != nil {
                return err
            }
        }
        my.lexParser.ParseCharacters(start_offSet, end_offSet,monitor)
        var index int
        if end_offSet >= my.lexStream.GetStreamIndex(){
            index =my.lexStream.GetStreamIndex()
        }else{
            index = end_offSet + 1
        }
        my.AddEOF(prsStream,index)
        return nil
    }
    
    func (my *JavaLexer) Lexer(prsStream IPrsStream ,  monitor Monitor) error{
    
        err := my.InitializeLexer(prsStream, 0, -1)
        if err != nil {
            return err
        }
        my.lexParser.ParseCharactersWhitMonitor(monitor)
        my.AddEOF(prsStream, my.lexStream.GetStreamIndex())
        return nil
    }

    /**
     * If a parse stream was not passed to my Lexical analyser then we
     * simply report a lexical error. Otherwise, we produce a bad token.
     */
    func (my *JavaLexer)  ReportLexicalError(startLoc int, endLoc int) {
        var prs_stream = my.lexStream.GetIPrsStream()
        if prs_stream == nil{
            my.lexStream.ReportLexicalErrorPosition(startLoc, endLoc)
        }else {
            //
            // Remove any token that may have been processed that fall in the
            // range of the lexical error... then add one error token that spans
            // the error range.
            //
            var i int = prs_stream.GetSize() - 1
            for ; i > 0 ;i-- {
                if prs_stream.GetStartOffset(i) >= startLoc {
                    prs_stream.RemoveLastToken()
                } else {
                break
                }
            }
            prs_stream.MakeToken(startLoc, endLoc, 0) // add an error token to the prsStream
        }        
    }

    //#line 176 "LexerBasicMapF.gi

//
// The Lexer contains an array of characters as the input stream to be parsed.
// There are methods to retrieve and classify characters.
// The lexparser "token" is implemented simply as the index of the next character in the array.
// The Lexer extends the abstract class LpgLexStream with an implementation of the abstract
// method GetKind.  The template defines the Lexer class and the Lexer() method.
// A driver creates the action class, "Lexer", passing an Option object to the constructor.
//

func (my *JavaLexer)  GetKeywordKinds()[]int { 
    return my.kwLexer.GetKeywordKinds()
}

func (my *JavaLexer) MakeToken(left_token int , right_token int , kind int ){
    my.lexStream.MakeToken(left_token, right_token, kind);
}

func (my *JavaLexer) MakeTokenWithKind(kind int){
    var startOffset = my.GetLeftSpan()
    var endOffset = my.GetRightSpan()
    my.lexStream.MakeToken(startOffset, endOffset, kind)
    if my.PrintTokens{
        my.PrintValue(startOffset, endOffset)
    }
}

func (my *JavaLexer) MakeComment(kind int){
    var startOffset = my.GetLeftSpan()
    var endOffset = my.GetRightSpan()
    my.lexStream.GetIPrsStream().MakeAdjunct(startOffset, endOffset, kind);
}

func (my *JavaLexer) SkipToken(){
    if my.PrintTokens{
        my.PrintValue(my.GetLeftSpan(), my.GetRightSpan())
    }
}

func (my *JavaLexer) CheckForKeyWord(){
    var startOffset = my.GetLeftSpan()
    var endOffset = my.GetRightSpan()
    var kwKind = my.kwLexer.Lexer(startOffset, endOffset)
    my.lexStream.MakeToken(startOffset, endOffset, kwKind)
    if my.PrintTokens{
        my.PrintValue(startOffset, endOffset)
    }
}

//
// This flavor of CheckForKeyWord is necessary when the default kind
// (which is returned when the keyword filter doesn't match) is something
// other than _IDENTIFIER.
//

func (my *JavaLexer) CheckForKeyWordWithKind(defaultKind int){
    var startOffset = my.GetLeftSpan()
    var endOffset = my.GetRightSpan()
    var    kwKind = my.kwLexer.Lexer(startOffset, endOffset)
    if kwKind == JavaParsersym.TK_IDENTIFIER{
        kwKind = defaultKind
    }
    my.lexStream.MakeToken(startOffset, endOffset, kwKind)
    if my.PrintTokens{
        my.PrintValue(startOffset, endOffset)
    }
}

func (my *JavaLexer) PrintValue(startOffset int, endOffset int){
    var s = my.lexStream.GetInputChars()[startOffset: endOffset  + 1]
    print(string(s))
}

    //#line 281 "LexerTemplateF.gi

    func (my *JavaLexer)  RuleAction(ruleNumber int){
        switch ruleNumber {

            //
            // Rule 1:  Token ::= Identifier
            //
            case 1: { 
                my.CheckForKeyWord()
                  
                  break
            }
    
            //
            // Rule 2:  Token ::= " SLBody "
            //
            case 2: { 
                my.MakeTokenWithKind(JavaParsersym.TK_StringLiteral)
                  
                  break
            }
    
            //
            // Rule 3:  Token ::= ' NotSQ '
            //
            case 3: { 
                my.MakeTokenWithKind(JavaParsersym.TK_CharacterLiteral)
                  
                  break
            }
    
            //
            // Rule 4:  Token ::= IntegerLiteral
            //
            case 4: { 
                my.MakeTokenWithKind(JavaParsersym.TK_IntegerLiteral)
                  
                  break
            }
    
            //
            // Rule 5:  Token ::= FloatingPointLiteral
            //
            case 5: { 
                my.MakeTokenWithKind(JavaParsersym.TK_FloatingPointLiteral)
                  
                  break
            }
    
            //
            // Rule 6:  Token ::= DoubleLiteral
            //
            case 6: { 
                my.MakeTokenWithKind(JavaParsersym.TK_DoubleLiteral)
                  
                  break
            }
    
            //
            // Rule 7:  Token ::= / * Inside Stars /
            //
            case 7: { 
                my.SkipToken()
                  
                  break
            }
    
            //
            // Rule 8:  Token ::= SLC
            //
            case 8: { 
                my.SkipToken()
                  
                  break
            }
    
            //
            // Rule 9:  Token ::= WS
            //
            case 9: { 
                my.SkipToken()
                  
                  break
            }
    
            //
            // Rule 10:  Token ::= +
            //
            case 10: { 
                my.MakeTokenWithKind(JavaParsersym.TK_PLUS)
                  
                  break
            }
    
            //
            // Rule 11:  Token ::= -
            //
            case 11: { 
                my.MakeTokenWithKind(JavaParsersym.TK_MINUS)
                  
                  break
            }
    
            //
            // Rule 12:  Token ::= *
            //
            case 12: { 
                my.MakeTokenWithKind(JavaParsersym.TK_MULTIPLY)
                  
                  break
            }
    
            //
            // Rule 13:  Token ::= /
            //
            case 13: { 
                my.MakeTokenWithKind(JavaParsersym.TK_DIVIDE)
                  
                  break
            }
    
            //
            // Rule 14:  Token ::= (
            //
            case 14: { 
                my.MakeTokenWithKind(JavaParsersym.TK_LPAREN)
                  
                  break
            }
    
            //
            // Rule 15:  Token ::= )
            //
            case 15: { 
                my.MakeTokenWithKind(JavaParsersym.TK_RPAREN)
                  
                  break
            }
    
            //
            // Rule 16:  Token ::= =
            //
            case 16: { 
                my.MakeTokenWithKind(JavaParsersym.TK_EQUAL)
                  
                  break
            }
    
            //
            // Rule 17:  Token ::= ,
            //
            case 17: { 
                my.MakeTokenWithKind(JavaParsersym.TK_COMMA)
                  
                  break
            }
    
            //
            // Rule 18:  Token ::= :
            //
            case 18: { 
                my.MakeTokenWithKind(JavaParsersym.TK_COLON)
                  
                  break
            }
    
            //
            // Rule 19:  Token ::= ;
            //
            case 19: { 
                my.MakeTokenWithKind(JavaParsersym.TK_SEMICOLON)
                  
                  break
            }
    
            //
            // Rule 20:  Token ::= ^
            //
            case 20: { 
                my.MakeTokenWithKind(JavaParsersym.TK_XOR)
                  
                  break
            }
    
            //
            // Rule 21:  Token ::= %
            //
            case 21: { 
                my.MakeTokenWithKind(JavaParsersym.TK_REMAINDER)
                  
                  break
            }
    
            //
            // Rule 22:  Token ::= ~
            //
            case 22: { 
                my.MakeTokenWithKind(JavaParsersym.TK_TWIDDLE)
                  
                  break
            }
    
            //
            // Rule 23:  Token ::= |
            //
            case 23: { 
                my.MakeTokenWithKind(JavaParsersym.TK_OR)
                  
                  break
            }
    
            //
            // Rule 24:  Token ::= &
            //
            case 24: { 
                my.MakeTokenWithKind(JavaParsersym.TK_AND)
                  
                  break
            }
    
            //
            // Rule 25:  Token ::= <
            //
            case 25: { 
                my.MakeTokenWithKind(JavaParsersym.TK_LESS)
                  
                  break
            }
    
            //
            // Rule 26:  Token ::= >
            //
            case 26: { 
                my.MakeTokenWithKind(JavaParsersym.TK_GREATER)
                  
                  break
            }
    
            //
            // Rule 27:  Token ::= .
            //
            case 27: { 
                my.MakeTokenWithKind(JavaParsersym.TK_DOT)
                  
                  break
            }
    
            //
            // Rule 28:  Token ::= !
            //
            case 28: { 
                my.MakeTokenWithKind(JavaParsersym.TK_NOT)
                  
                  break
            }
    
            //
            // Rule 29:  Token ::= [
            //
            case 29: { 
                my.MakeTokenWithKind(JavaParsersym.TK_LBRACKET)
                  
                  break
            }
    
            //
            // Rule 30:  Token ::= ]
            //
            case 30: { 
                my.MakeTokenWithKind(JavaParsersym.TK_RBRACKET)
                  
                  break
            }
    
            //
            // Rule 31:  Token ::= {
            //
            case 31: { 
                my.MakeTokenWithKind(JavaParsersym.TK_LBRACE)
                  
                  break
            }
    
            //
            // Rule 32:  Token ::= }
            //
            case 32: { 
                my.MakeTokenWithKind(JavaParsersym.TK_RBRACE)
                  
                  break
            }
    
            //
            // Rule 33:  Token ::= ?
            //
            case 33: { 
                my.MakeTokenWithKind(JavaParsersym.TK_QUESTION)
                  
                  break
            }
    
            //
            // Rule 34:  Token ::= @
            //
            case 34: { 
                my.MakeTokenWithKind(JavaParsersym.TK_AT)
                  
                  break
            }
    
            //
            // Rule 35:  Token ::= + +
            //
            case 35: { 
                my.MakeTokenWithKind(JavaParsersym.TK_PLUS_PLUS)
                  
                  break
            }
    
            //
            // Rule 36:  Token ::= - -
            //
            case 36: { 
                my.MakeTokenWithKind(JavaParsersym.TK_MINUS_MINUS)
                  
                  break
            }
    
            //
            // Rule 37:  Token ::= = =
            //
            case 37: { 
                my.MakeTokenWithKind(JavaParsersym.TK_EQUAL_EQUAL)
                  
                  break
            }
    
            //
            // Rule 38:  Token ::= < =
            //
            case 38: { 
                my.MakeTokenWithKind(JavaParsersym.TK_LESS_EQUAL)
                  
                  break
            }
    
            //
            // Rule 39:  Token ::= ! =
            //
            case 39: { 
                my.MakeTokenWithKind(JavaParsersym.TK_NOT_EQUAL)
                  
                  break
            }
    
            //
            // Rule 40:  Token ::= < <
            //
            case 40: { 
                my.MakeTokenWithKind(JavaParsersym.TK_LEFT_SHIFT)
                  
                  break
            }
    
            //
            // Rule 41:  Token ::= + =
            //
            case 41: { 
                my.MakeTokenWithKind(JavaParsersym.TK_PLUS_EQUAL)
                  
                  break
            }
    
            //
            // Rule 42:  Token ::= - =
            //
            case 42: { 
                my.MakeTokenWithKind(JavaParsersym.TK_MINUS_EQUAL)
                  
                  break
            }
    
            //
            // Rule 43:  Token ::= * =
            //
            case 43: { 
                my.MakeTokenWithKind(JavaParsersym.TK_MULTIPLY_EQUAL)
                  
                  break
            }
    
            //
            // Rule 44:  Token ::= / =
            //
            case 44: { 
                my.MakeTokenWithKind(JavaParsersym.TK_DIVIDE_EQUAL)
                  
                  break
            }
    
            //
            // Rule 45:  Token ::= & =
            //
            case 45: { 
                my.MakeTokenWithKind(JavaParsersym.TK_AND_EQUAL)
                  
                  break
            }
    
            //
            // Rule 46:  Token ::= | =
            //
            case 46: { 
                my.MakeTokenWithKind(JavaParsersym.TK_OR_EQUAL)
                  
                  break
            }
    
            //
            // Rule 47:  Token ::= ^ =
            //
            case 47: { 
                my.MakeTokenWithKind(JavaParsersym.TK_XOR_EQUAL)
                  
                  break
            }
    
            //
            // Rule 48:  Token ::= % =
            //
            case 48: { 
                my.MakeTokenWithKind(JavaParsersym.TK_REMAINDER_EQUAL)
                  
                  break
            }
    
            //
            // Rule 49:  Token ::= < < =
            //
            case 49: { 
                my.MakeTokenWithKind(JavaParsersym.TK_LEFT_SHIFT_EQUAL)
                  
                  break
            }
    
            //
            // Rule 50:  Token ::= | |
            //
            case 50: { 
                my.MakeTokenWithKind(JavaParsersym.TK_OR_OR)
                  
                  break
            }
    
            //
            // Rule 51:  Token ::= & &
            //
            case 51: { 
                my.MakeTokenWithKind(JavaParsersym.TK_AND_AND)
                  
                  break
            }
    
            //
            // Rule 52:  Token ::= . . .
            //
            case 52: { 
                my.MakeTokenWithKind(JavaParsersym.TK_ELLIPSIS)
                  
                  break
            }
    
    //#line 285 "LexerTemplateF.gi

    
            default:
                break
        }
        return
    }

    //#line 3 "LexerBasicMapF.gi
 
type  JavaLexerLpgLexStream struct{
    *LexStream
}
func  NewJavaLexerLpgLexStream( fileName string, inputChars []rune, tab int) (*JavaLexerLpgLexStream,error){
        t := new(JavaLexerLpgLexStream)
        var e error
        t.LexStream ,e = NewLexStreamExt(t,fileName,inputChars,tab,nil)
        return t,e
}

var JavaLexerLpgLexStreamtokenKind =[]int{
        JavaLexersym.Char_CtlCharNotWS,    // 000    0x00
        JavaLexersym.Char_CtlCharNotWS,    // 001    0x01
        JavaLexersym.Char_CtlCharNotWS,    // 002    0x02
        JavaLexersym.Char_CtlCharNotWS,    // 003    0x03
        JavaLexersym.Char_CtlCharNotWS,    // 004    0x04
        JavaLexersym.Char_CtlCharNotWS,    // 005    0x05
        JavaLexersym.Char_CtlCharNotWS,    // 006    0x06
        JavaLexersym.Char_CtlCharNotWS,    // 007    0x07
        JavaLexersym.Char_CtlCharNotWS,    // 008    0x08
        JavaLexersym.Char_HT,              // 009    0x09
        JavaLexersym.Char_LF,              // 010    0x0A
        JavaLexersym.Char_CtlCharNotWS,    // 011    0x0B
        JavaLexersym.Char_FF,              // 012    0x0C
        JavaLexersym.Char_CR,              // 013    0x0D
        JavaLexersym.Char_CtlCharNotWS,    // 014    0x0E
        JavaLexersym.Char_CtlCharNotWS,    // 015    0x0F
        JavaLexersym.Char_CtlCharNotWS,    // 016    0x10
        JavaLexersym.Char_CtlCharNotWS,    // 017    0x11
        JavaLexersym.Char_CtlCharNotWS,    // 018    0x12
        JavaLexersym.Char_CtlCharNotWS,    // 019    0x13
        JavaLexersym.Char_CtlCharNotWS,    // 020    0x14
        JavaLexersym.Char_CtlCharNotWS,    // 021    0x15
        JavaLexersym.Char_CtlCharNotWS,    // 022    0x16
        JavaLexersym.Char_CtlCharNotWS,    // 023    0x17
        JavaLexersym.Char_CtlCharNotWS,    // 024    0x18
        JavaLexersym.Char_CtlCharNotWS,    // 025    0x19
        JavaLexersym.Char_CtlCharNotWS,    // 026    0x1A
        JavaLexersym.Char_CtlCharNotWS,    // 027    0x1B
        JavaLexersym.Char_CtlCharNotWS,    // 028    0x1C
        JavaLexersym.Char_CtlCharNotWS,    // 029    0x1D
        JavaLexersym.Char_CtlCharNotWS,    // 030    0x1E
        JavaLexersym.Char_CtlCharNotWS,    // 031    0x1F
        JavaLexersym.Char_Space,           // 032    0x20
        JavaLexersym.Char_Exclamation,     // 033    0x21
        JavaLexersym.Char_DoubleQuote,     // 034    0x22
        JavaLexersym.Char_Sharp,           // 035    0x23
        JavaLexersym.Char_DollarSign,      // 036    0x24
        JavaLexersym.Char_Percent,         // 037    0x25
        JavaLexersym.Char_Ampersand,       // 038    0x26
        JavaLexersym.Char_SingleQuote,     // 039    0x27
        JavaLexersym.Char_LeftParen,       // 040    0x28
        JavaLexersym.Char_RightParen,      // 041    0x29
        JavaLexersym.Char_Star,            // 042    0x2A
        JavaLexersym.Char_Plus,            // 043    0x2B
        JavaLexersym.Char_Comma,           // 044    0x2C
        JavaLexersym.Char_Minus,           // 045    0x2D
        JavaLexersym.Char_Dot,             // 046    0x2E
        JavaLexersym.Char_Slash,           // 047    0x2F
        JavaLexersym.Char_0,               // 048    0x30
        JavaLexersym.Char_1,               // 049    0x31
        JavaLexersym.Char_2,               // 050    0x32
        JavaLexersym.Char_3,               // 051    0x33
        JavaLexersym.Char_4,               // 052    0x34
        JavaLexersym.Char_5,               // 053    0x35
        JavaLexersym.Char_6,               // 054    0x36
        JavaLexersym.Char_7,               // 055    0x37
        JavaLexersym.Char_8,               // 056    0x38
        JavaLexersym.Char_9,               // 057    0x39
        JavaLexersym.Char_Colon,           // 058    0x3A
        JavaLexersym.Char_SemiColon,       // 059    0x3B
        JavaLexersym.Char_LessThan,        // 060    0x3C
        JavaLexersym.Char_Equal,           // 061    0x3D
        JavaLexersym.Char_GreaterThan,     // 062    0x3E
        JavaLexersym.Char_QuestionMark,    // 063    0x3F
        JavaLexersym.Char_AtSign,          // 064    0x40
        JavaLexersym.Char_A,               // 065    0x41
        JavaLexersym.Char_B,               // 066    0x42
        JavaLexersym.Char_C,               // 067    0x43
        JavaLexersym.Char_D,               // 068    0x44
        JavaLexersym.Char_E,               // 069    0x45
        JavaLexersym.Char_F,               // 070    0x46
        JavaLexersym.Char_G,               // 071    0x47
        JavaLexersym.Char_H,               // 072    0x48
        JavaLexersym.Char_I,               // 073    0x49
        JavaLexersym.Char_J,               // 074    0x4A
        JavaLexersym.Char_K,               // 075    0x4B
        JavaLexersym.Char_L,               // 076    0x4C
        JavaLexersym.Char_M,               // 077    0x4D
        JavaLexersym.Char_N,               // 078    0x4E
        JavaLexersym.Char_O,               // 079    0x4F
        JavaLexersym.Char_P,               // 080    0x50
        JavaLexersym.Char_Q,               // 081    0x51
        JavaLexersym.Char_R,               // 082    0x52
        JavaLexersym.Char_S,               // 083    0x53
        JavaLexersym.Char_T,               // 084    0x54
        JavaLexersym.Char_U,               // 085    0x55
        JavaLexersym.Char_V,               // 086    0x56
        JavaLexersym.Char_W,               // 087    0x57
        JavaLexersym.Char_X,               // 088    0x58
        JavaLexersym.Char_Y,               // 089    0x59
        JavaLexersym.Char_Z,               // 090    0x5A
        JavaLexersym.Char_LeftBracket,     // 091    0x5B
        JavaLexersym.Char_BackSlash,       // 092    0x5C
        JavaLexersym.Char_RightBracket,    // 093    0x5D
        JavaLexersym.Char_Caret,           // 094    0x5E
        JavaLexersym.Char__,               // 095    0x5F
        JavaLexersym.Char_BackQuote,       // 096    0x60
        JavaLexersym.Char_a,               // 097    0x61
        JavaLexersym.Char_b,               // 098    0x62
        JavaLexersym.Char_c,               // 099    0x63
        JavaLexersym.Char_d,               // 100    0x64
        JavaLexersym.Char_e,               // 101    0x65
        JavaLexersym.Char_f,               // 102    0x66
        JavaLexersym.Char_g,               // 103    0x67
        JavaLexersym.Char_h,               // 104    0x68
        JavaLexersym.Char_i,               // 105    0x69
        JavaLexersym.Char_j,               // 106    0x6A
        JavaLexersym.Char_k,               // 107    0x6B
        JavaLexersym.Char_l,               // 108    0x6C
        JavaLexersym.Char_m,               // 109    0x6D
        JavaLexersym.Char_n,               // 110    0x6E
        JavaLexersym.Char_o,               // 111    0x6F
        JavaLexersym.Char_p,               // 112    0x70
        JavaLexersym.Char_q,               // 113    0x71
        JavaLexersym.Char_r,               // 114    0x72
        JavaLexersym.Char_s,               // 115    0x73
        JavaLexersym.Char_t,               // 116    0x74
        JavaLexersym.Char_u,               // 117    0x75
        JavaLexersym.Char_v,               // 118    0x76
        JavaLexersym.Char_w,               // 119    0x77
        JavaLexersym.Char_x,               // 120    0x78
        JavaLexersym.Char_y,               // 121    0x79
        JavaLexersym.Char_z,               // 122    0x7A
        JavaLexersym.Char_LeftBrace,       // 123    0x7B
        JavaLexersym.Char_VerticalBar,     // 124    0x7C
        JavaLexersym.Char_RightBrace,      // 125    0x7D
        JavaLexersym.Char_Tilde,           // 126    0x7E

        JavaLexersym.Char_AfterASCII,      // for all chars in range 128..65534
        JavaLexersym.Char_EOF ,             // for '\uffff' or 65535 
}

    
func (my *JavaLexerLpgLexStream) GetKind(i int) int {  // Classify character at ith location
        var c int 
        if i >= my.GetStreamLength() {
            c = 0xffff
        }else{
            c = my.GetIntValue(i)
        }
        if c < 128 {//ASCII Character
            return JavaLexerLpgLexStreamtokenKind[c]
        }else{
            if c == 0xffff{
                return JavaLexersym.Char_EOF
            }else{
                return JavaLexersym.Char_AfterASCII
            }
        }

}
func (my *JavaLexerLpgLexStream)  OrderedExportedSymbols() []string{
    return JavaParsersym.OrderedTerminalSymbols 
}   

