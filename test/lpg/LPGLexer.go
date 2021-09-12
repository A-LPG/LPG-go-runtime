
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


    //#line 122 "LexerTemplateF.gi

import ."github.com/A-LPG/LPG-go-runtime/lpg2"

    //#line 7 "LPGLexer.gi



    //#line 128 "LexerTemplateF.gi

type  LPGLexer struct{

    kwLexer  *LPGKWLexer 
    PrintTokens  bool
    lexParser  *LexParser
    lexStream *LPGLexerLpgLexStream 
    prs     ParseTable
}
func NewLPGLexer(filename string,  tab int ,input_chars []rune) *LPGLexer {
    my := new(LPGLexer)
    my.prs =  NewLPGLexerprs()
    my.lexParser   =  NewLexParser()
    my.PrintTokens = false
    var e error
    my.lexStream, e =  NewLPGLexerLpgLexStream(filename,input_chars, tab)
    if e != nil{
        return nil
    }
    my.lexParser.Reset( my.lexStream, my.prs,  my)
    my.ResetKeywordLexer()
    return my
}
    func (my *LPGLexer)  GetParseTable()ParseTable {
     return my.prs 
    }

    func (my *LPGLexer)  GetParser() *LexParser{ 
        return my.lexParser 
    }

    func (my *LPGLexer)  GetToken(i int)int { 
    return my.lexParser.GetToken(i)
     }
    func (my *LPGLexer)  GetRhsFirstTokenIndex(i int)int {
     return my.lexParser.GetFirstTokenAt(i) 
     }
    func (my *LPGLexer)  GetRhsLastTokenIndex(i int) int{
     return my.lexParser.GetLastTokenAt(i) 
     }

    func (my *LPGLexer)  GetLeftSpan()int { 
    return my.lexParser.GetToken(1) 
    }
    func (my *LPGLexer)  GetRightSpan()int {
     return my.lexParser.GetLastToken() 
     }

    func (my *LPGLexer)  ResetKeywordLexer(){
        if my.kwLexer == nil {
            my.kwLexer = NewLPGKWLexer(my.lexStream.GetInputChars(), LPGParsersym.TK_MACRO_NAME)
        }else {
            my.kwLexer.SetInputChars(my.lexStream.GetInputChars())
        }
    }

    func (my *LPGLexer)  Reset(filename string,tab int,input_chars []rune) error{
        var e error
        my.lexStream ,e = NewLPGLexerLpgLexStream(filename,input_chars, tab)
        if e != nil{
            return e
        }
        my.lexParser.Reset(my.lexStream, my.prs, my)
        my.ResetKeywordLexer()
        return nil
    }
    

    func (my *LPGLexer)  GetILexStream()ILexStream { 
        return my.lexStream
    }



    func (my *LPGLexer)  InitializeLexer(prsStream IPrsStream, start_offSet int, end_offSet int) error{
        if my.lexStream.GetInputChars() == nil{
            return NewNullPointerException("LexStream was not initialized")
        }
        my.lexStream.SetPrsStream(prsStream)
        prsStream.SetLexStream(my.lexStream)
        prsStream.MakeToken(start_offSet, end_offSet, 0) // Token list must start with a bad token
        return nil
    }

    func (my *LPGLexer)  AddEOF(prsStream IPrsStream, end_offSet int) {
        prsStream.MakeToken(end_offSet, end_offSet, LPGParsersym.TK_EOF_TOKEN) // and end with the end of file token
        prsStream.SetStreamLength(prsStream.GetSize())
    }

    func (my *LPGLexer) LexerWithPosition(prsStream IPrsStream , start_offSet int , end_offSet int, monitor Monitor) error{
    
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
    
    func (my *LPGLexer) Lexer(prsStream IPrsStream ,  monitor Monitor) error{
    
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
    func (my *LPGLexer)  ReportLexicalError(startLoc int, endLoc int) {
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

    //#line 12 "LPGLexer.gi

 
    //#line 176 "LexerBasicMapF.gi

//
// The Lexer contains an array of characters as the input stream to be parsed.
// There are methods to retrieve and classify characters.
// The lexparser "token" is implemented simply as the index of the next character in the array.
// The Lexer extends the abstract class LpgLexStream with an implementation of the abstract
// method GetKind.  The template defines the Lexer class and the Lexer() method.
// A driver creates the action class, "Lexer", passing an Option object to the constructor.
//

func (my *LPGLexer)  GetKeywordKinds()[]int { 
    return my.kwLexer.GetKeywordKinds()
}

func (my *LPGLexer) MakeToken(left_token int , right_token int , kind int ){
    my.lexStream.MakeToken(left_token, right_token, kind);
}

func (my *LPGLexer) MakeTokenWithKind(kind int){
    var startOffset = my.GetLeftSpan()
    var endOffset = my.GetRightSpan()
    my.lexStream.MakeToken(startOffset, endOffset, kind)
    if my.PrintTokens{
        my.PrintValue(startOffset, endOffset)
    }
}

func (my *LPGLexer) MakeComment(kind int){
    var startOffset = my.GetLeftSpan()
    var endOffset = my.GetRightSpan()
    my.lexStream.GetIPrsStream().MakeAdjunct(startOffset, endOffset, kind);
}

func (my *LPGLexer) SkipToken(){
    if my.PrintTokens{
        my.PrintValue(my.GetLeftSpan(), my.GetRightSpan())
    }
}

func (my *LPGLexer) CheckForKeyWord(){
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

func (my *LPGLexer) CheckForKeyWordWithKind(defaultKind int){
    var startOffset = my.GetLeftSpan()
    var endOffset = my.GetRightSpan()
    var    kwKind = my.kwLexer.Lexer(startOffset, endOffset)
    if kwKind == LPGParsersym.TK_MACRO_NAME{
        kwKind = defaultKind
    }
    my.lexStream.MakeToken(startOffset, endOffset, kwKind)
    if my.PrintTokens{
        my.PrintValue(startOffset, endOffset)
    }
}

func (my *LPGLexer) PrintValue(startOffset int, endOffset int){
    var s = my.lexStream.GetInputChars()[startOffset: endOffset  + 1]
    print(string(s))
}

    //#line 281 "LexerTemplateF.gi

    func (my *LPGLexer)  RuleAction(ruleNumber int){
        switch ruleNumber {

            //
            // Rule 1:  Token ::= white
            //
            case 1: { 
             my.SkipToken()             
                  break
            }

            //
            // Rule 2:  Token ::= singleLineComment
            //
            case 2: { 
             my.MakeComment(LPGParsersym.TK_SINGLE_LINE_COMMENT)             
                  break
            }

            //
            // Rule 4:  Token ::= MacroSymbol
            //
            case 4: { 
             my.CheckForKeyWord()            
                  break
            }

            //
            // Rule 5:  Token ::= Symbol
            //
            case 5: { 
             my.CheckForKeyWordWithKind(LPGParsersym.TK_SYMBOL)            
                  break
            }

            //
            // Rule 6:  Token ::= Block
            //
            case 6: { 
             my.MakeTokenWithKind(LPGParsersym.TK_BLOCK)            
                  break
            }

            //
            // Rule 7:  Token ::= Equivalence
            //
            case 7: { 
             my.MakeTokenWithKind(LPGParsersym.TK_EQUIVALENCE)            
                  break
            }

            //
            // Rule 8:  Token ::= Equivalence ?
            //
            case 8: { 
             my.MakeTokenWithKind(LPGParsersym.TK_PRIORITY_EQUIVALENCE)            
                  break
            }

            //
            // Rule 9:  Token ::= #
            //
            case 9: { 
             my.MakeTokenWithKind(LPGParsersym.TK_SHARP)            
                  break
            }

            //
            // Rule 10:  Token ::= Arrow
            //
            case 10: { 
             my.MakeTokenWithKind(LPGParsersym.TK_ARROW)            
                  break
            }

            //
            // Rule 11:  Token ::= Arrow ?
            //
            case 11: { 
             my.MakeTokenWithKind(LPGParsersym.TK_PRIORITY_ARROW)            
                  break
            }

            //
            // Rule 12:  Token ::= |
            //
            case 12: { 
             my.MakeTokenWithKind(LPGParsersym.TK_OR_MARKER)            
                  break
            }

            //
            // Rule 13:  Token ::= [
            //
            case 13: { 
             my.MakeTokenWithKind(LPGParsersym.TK_LEFT_BRACKET)            
                  break
            }

            //
            // Rule 14:  Token ::= ]
            //
            case 14: { 
             my.MakeTokenWithKind(LPGParsersym.TK_RIGHT_BRACKET)            
                  break
            }

            //
            // Rule 858:  OptionLines ::= OptionLineList
            //
            case 858: { 
            
                  // What ever needs to happen after the options have been 
                  // scanned must happen here.
                    
                  break
            }
      
            //
            // Rule 867:  options ::= % oO pP tT iI oO nN sS
            //
            case 867: { 
            
                  my.MakeToken(my.GetLeftSpan(), my.GetRightSpan(), LPGParsersym.TK_OPTIONS_KEY)
                    
                  break
            }
      
            //
            // Rule 868:  OptionComment ::= singleLineComment
            //
            case 868: { 
             my.MakeComment(LPGParsersym.TK_SINGLE_LINE_COMMENT)             
                  break
            }

            //
            // Rule 892:  separator ::= ,$comma
            //
            case 892: { 
              my.MakeToken(my.GetLeftSpan(), my.GetRightSpan(), LPGParsersym.TK_COMMA)             
                  break
            }

            //
            // Rule 893:  option ::= action_block$ab optionWhite =$eq optionWhite ($lp optionWhite filename$fn optionWhite ,$comma1 optionWhite block_begin$bb optionWhite ,$comma2 optionWhite block_end$be optionWhite )$rp optionWhite
            //
            case 893: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_LEFT_PAREN)
                  my.MakeToken(my.GetRhsFirstTokenIndex(7), my.GetRhsLastTokenIndex(7), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(9), my.GetRhsLastTokenIndex(9), LPGParsersym.TK_COMMA)
                  my.MakeToken(my.GetRhsFirstTokenIndex(11), my.GetRhsLastTokenIndex(11), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(13), my.GetRhsLastTokenIndex(13), LPGParsersym.TK_COMMA)
                  my.MakeToken(my.GetRhsFirstTokenIndex(15), my.GetRhsLastTokenIndex(15), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(17), my.GetRhsLastTokenIndex(17), LPGParsersym.TK_RIGHT_PAREN)
                    
                  break
            }
      
            //
            // Rule 896:  option ::= ast_block$ab optionWhite =$eq optionWhite ($lp optionWhite block_begin$bb optionWhite ,$comma2 optionWhite block_end$be optionWhite )$rp optionWhite
            //
            case 896: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_LEFT_PAREN)

                  my.MakeToken(my.GetRhsFirstTokenIndex(7), my.GetRhsLastTokenIndex(7), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(9), my.GetRhsLastTokenIndex(9), LPGParsersym.TK_COMMA)
                  my.MakeToken(my.GetRhsFirstTokenIndex(11), my.GetRhsLastTokenIndex(11), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(13), my.GetRhsLastTokenIndex(13), LPGParsersym.TK_RIGHT_PAREN)
                    
                  break
            }
      
            //
            // Rule 901:  option ::= ast_directory$ad optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 901: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 904:  option ::= ast_type$at optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 904: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 907:  option ::= attributes$a optionWhite
            //
            case 907: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 908:  option ::= no attributes$a optionWhite
            //
            case 908: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 910:  option ::= automatic_ast$a optionWhite
            //
            case 910: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 911:  option ::= no automatic_ast$a optionWhite
            //
            case 911: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 912:  option ::= automatic_ast$aa optionWhite =$eq optionWhite automatic_ast_value$val optionWhite
            //
            case 912: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 916:  option ::= backtrack$b optionWhite
            //
            case 916: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 917:  option ::= no backtrack$b optionWhite
            //
            case 917: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 919:  option ::= byte$b optionWhite
            //
            case 919: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 920:  option ::= no byte$b optionWhite
            //
            case 920: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 922:  option ::= conflicts$c optionWhite
            //
            case 922: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 923:  option ::= no conflicts$c optionWhite
            //
            case 923: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 925:  option ::= dat_directory$dd optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 925: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 928:  option ::= dat_file$df optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 928: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 930:  option ::= dcl_file$df optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 930: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 932:  option ::= def_file$df optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 932: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 934:  option ::= debug$d optionWhite
            //
            case 934: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 935:  option ::= no debug$d optionWhite
            //
            case 935: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 937:  option ::= edit$e optionWhite
            //
            case 937: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 938:  option ::= no edit$e optionWhite
            //
            case 938: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 940:  option ::= error_maps$e optionWhite
            //
            case 940: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 941:  option ::= no error_maps$e optionWhite
            //
            case 941: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 944:  option ::= escape$e optionWhite =$eq optionWhite anyNonWhiteChar$val optionWhite
            //
            case 944: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                
                    
                  break
            }
      
            //
            // Rule 946:  option ::= export_terminals$et optionWhite =$eq optionWhite filename$fn optionWhite
            //
            case 946: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 947:  option ::= export_terminals$et optionWhite =$eq optionWhite ($lp optionWhite filename$fn optionWhite )$rp optionWhite
            //
            case 947: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_LEFT_PAREN)
                  my.MakeToken(my.GetRhsFirstTokenIndex(7), my.GetRhsLastTokenIndex(7), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(9), my.GetRhsLastTokenIndex(9), LPGParsersym.TK_RIGHT_PAREN)
                    
                  break
            }
      
            //
            // Rule 948:  option ::= export_terminals$et optionWhite =$eq optionWhite ($lp optionWhite filename$fn optionWhite ,$comma optionWhite export_prefix$ep optionWhite )$rp optionWhite
            //
            case 948: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_LEFT_PAREN)
                  my.MakeToken(my.GetRhsFirstTokenIndex(7), my.GetRhsLastTokenIndex(7), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(9), my.GetRhsLastTokenIndex(9), LPGParsersym.TK_COMMA)
                  my.MakeToken(my.GetRhsFirstTokenIndex(11), my.GetRhsLastTokenIndex(11), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(13), my.GetRhsLastTokenIndex(13), LPGParsersym.TK_RIGHT_PAREN)
                    
                  break
            }
      
            //
            // Rule 949:  option ::= export_terminals$et optionWhite =$eq optionWhite ($lp optionWhite filename$fn optionWhite ,$comma1 optionWhite export_prefix$ep optionWhite ,$comma2 optionWhite export_suffix$es optionWhite )$rp optionWhite
            //
            case 949: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_LEFT_PAREN)
                  my.MakeToken(my.GetRhsFirstTokenIndex(7), my.GetRhsLastTokenIndex(7), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(9), my.GetRhsLastTokenIndex(9), LPGParsersym.TK_COMMA)
                  my.MakeToken(my.GetRhsFirstTokenIndex(11), my.GetRhsLastTokenIndex(11), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(13), my.GetRhsLastTokenIndex(13), LPGParsersym.TK_COMMA)
                  my.MakeToken(my.GetRhsFirstTokenIndex(15), my.GetRhsLastTokenIndex(15), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(17), my.GetRhsLastTokenIndex(17), LPGParsersym.TK_RIGHT_PAREN)
                    
                  break
            }
      
            //
            // Rule 954:  option ::= extends_parsetable$e optionWhite
            //
            case 954: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 955:  option ::= no extends_parsetable$e optionWhite
            //
            case 955: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 956:  option ::= extends_parsetable$ep optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 956: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 959:  option ::= factory$f optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 959: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 961:  option ::= file_prefix$fp optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 961: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 964:  option ::= filter$f optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 964: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 966:  option ::= first$f optionWhite
            //
            case 966: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 967:  option ::= no first$f optionWhite
            //
            case 967: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 969:  option ::= follow$f optionWhite
            //
            case 969: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 970:  option ::= no follow$f optionWhite
            //
            case 970: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 972:  option ::= goto_default$g optionWhite
            //
            case 972: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 973:  option ::= no goto_default$g optionWhite
            //
            case 973: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 976:  option ::= headers$h optionWhite =$eq optionWhite ($lp optionWhite filename$fn optionWhite ,$comma1 optionWhite block_begin$bb optionWhite ,$comma2 optionWhite block_end$be optionWhite )$rp optionWhite
            //
            case 976: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_LEFT_PAREN)
                  my.MakeToken(my.GetRhsFirstTokenIndex(7), my.GetRhsLastTokenIndex(7), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(9), my.GetRhsLastTokenIndex(9), LPGParsersym.TK_COMMA)
                  my.MakeToken(my.GetRhsFirstTokenIndex(11), my.GetRhsLastTokenIndex(11), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(13), my.GetRhsLastTokenIndex(13), LPGParsersym.TK_COMMA)
                  my.MakeToken(my.GetRhsFirstTokenIndex(15), my.GetRhsLastTokenIndex(15), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(17), my.GetRhsLastTokenIndex(17), LPGParsersym.TK_RIGHT_PAREN)
                    
                  break
            }
      
            //
            // Rule 978:  option ::= imp_file$if optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 978: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 981:  option ::= import_terminals$it optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 981: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 984:  option ::= include_directory$id optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 984: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 988:  option ::= lalr_level$l optionWhite
            //
            case 988: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 989:  option ::= no lalr_level$l optionWhite
            //
            case 989: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 990:  option ::= lalr_level$l optionWhite =$eq optionWhite number$val optionWhite
            //
            case 990: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 995:  option ::= list$l optionWhite
            //
            case 995: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 996:  option ::= no list$l optionWhite
            //
            case 996: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 998:  option ::= margin$m optionWhite =$eq optionWhite number$val optionWhite
            //
            case 998: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1000:  option ::= max_cases$mc optionWhite =$eq optionWhite number$val optionWhite
            //
            case 1000: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1003:  option ::= names$n optionWhite
            //
            case 1003: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1004:  option ::= no names$n optionWhite
            //
            case 1004: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1005:  option ::= names$n optionWhite =$eq optionWhite names_value$val optionWhite
            //
            case 1005: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1010:  option ::= nt_check$n optionWhite
            //
            case 1010: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1011:  option ::= no nt_check$n optionWhite
            //
            case 1011: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1014:  option ::= or_marker$om optionWhite =$eq optionWhite anyNonWhiteChar$val optionWhite
            //
            case 1014: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1017:  option ::= out_directory$dd optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 1017: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1020:  option ::= parent_saved$ps optionWhite
            //
            case 1020: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1021:  option ::= no parent_saved$ps optionWhite
            //
            case 1021: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1024:  option ::= package$p optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 1024: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1026:  option ::= parsetable_interfaces$pi optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 1026: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1030:  option ::= prefix$p optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 1030: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1032:  option ::= priority$p optionWhite
            //
            case 1032: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1033:  option ::= no priority$p optionWhite
            //
            case 1033: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1035:  option ::= programming_language$pl optionWhite =$eq optionWhite programming_language_value$val optionWhite
            //
            case 1035: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1039:  option ::= prs_file$pf optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 1039: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1042:  option ::= quiet$q optionWhite
            //
            case 1042: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1043:  option ::= no quiet$q optionWhite
            //
            case 1043: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1045:  option ::= read_reduce$r optionWhite
            //
            case 1045: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1046:  option ::= no read_reduce$r optionWhite
            //
            case 1046: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1049:  option ::= remap_terminals$r optionWhite
            //
            case 1049: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1050:  option ::= no remap_terminals$r optionWhite
            //
            case 1050: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            }

            //
            // Rule 1053:  option ::= scopes$s optionWhite
            //
            case 1053: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1054:  option ::= no scopes$s optionWhite
            //
            case 1054: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1056:  option ::= serialize$s optionWhite
            //
            case 1056: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1057:  option ::= no serialize$s optionWhite
            //
            case 1057: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1059:  option ::= shift_default$s optionWhite
            //
            case 1059: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1060:  option ::= no shift_default$s optionWhite
            //
            case 1060: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1063:  option ::= single_productions$s optionWhite
            //
            case 1063: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1064:  option ::= no single_productions$s optionWhite
            //
            case 1064: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1067:  option ::= slr$s optionWhite
            //
            case 1067: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1068:  option ::= no slr$s optionWhite
            //
            case 1068: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1070:  option ::= soft_keywords$s optionWhite
            //
            case 1070: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1071:  option ::= no soft_keywords$s optionWhite
            //
            case 1071: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1075:  option ::= states$s optionWhite
            //
            case 1075: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1076:  option ::= no states$s optionWhite
            //
            case 1076: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1078:  option ::= suffix$s optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 1078: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1080:  option ::= sym_file$sf optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 1080: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1083:  option ::= tab_file$tf optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 1083: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1086:  option ::= template$t optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 1086: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1088:  option ::= trailers$t optionWhite =$eq optionWhite ($lp optionWhite filename$fn optionWhite ,$comma1 optionWhite block_begin$bb optionWhite ,$comma2 optionWhite block_end$be optionWhite )$rp optionWhite
            //
            case 1088: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_LEFT_PAREN)
                  my.MakeToken(my.GetRhsFirstTokenIndex(7), my.GetRhsLastTokenIndex(7), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(9), my.GetRhsLastTokenIndex(9), LPGParsersym.TK_COMMA)
                  my.MakeToken(my.GetRhsFirstTokenIndex(11), my.GetRhsLastTokenIndex(11), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(13), my.GetRhsLastTokenIndex(13), LPGParsersym.TK_COMMA)
                  my.MakeToken(my.GetRhsFirstTokenIndex(15), my.GetRhsLastTokenIndex(15), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(17), my.GetRhsLastTokenIndex(17), LPGParsersym.TK_RIGHT_PAREN)
                    
                  break
            }
      
            //
            // Rule 1090:  option ::= table$t optionWhite
            //
            case 1090: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1091:  option ::= no table$t optionWhite
            //
            case 1091: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1092:  option ::= table$t optionWhite =$eq optionWhite programming_language_value$val optionWhite
            //
            case 1092: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1094:  option ::= trace$t optionWhite
            //
            case 1094: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1095:  option ::= no trace$t optionWhite
            //
            case 1095: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1096:  option ::= trace$t optionWhite =$eq optionWhite trace_value$val optionWhite
            //
            case 1096: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1101:  option ::= variables$v optionWhite
            //
            case 1101: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1102:  option ::= no variables$v optionWhite
            //
            case 1102: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1103:  option ::= variables$v optionWhite =$eq optionWhite variables_value$val optionWhite
            //
            case 1103: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1110:  option ::= verbose$v optionWhite
            //
            case 1110: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1111:  option ::= no verbose$v optionWhite
            //
            case 1111: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1113:  option ::= visitor$v optionWhite
            //
            case 1113: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1114:  option ::= no visitor$v optionWhite
            //
            case 1114: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1115:  option ::= visitor$v optionWhite =$eq optionWhite visitor_value$val optionWhite
            //
            case 1115: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1120:  option ::= visitor_type$vt optionWhite =$eq optionWhite Value$val optionWhite
            //
            case 1120: { 
            
                  my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(3), my.GetRhsLastTokenIndex(3), LPGParsersym.TK_EQUAL)
                  my.MakeToken(my.GetRhsFirstTokenIndex(5), my.GetRhsLastTokenIndex(5), LPGParsersym.TK_SYMBOL)
                    
                  break
            }
      
            //
            // Rule 1123:  option ::= warnings$w optionWhite
            //
            case 1123: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1124:  option ::= no warnings$w optionWhite
            //
            case 1124: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1126:  option ::= xreference$x optionWhite
            //
            case 1126: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(1), my.GetRhsLastTokenIndex(1), LPGParsersym.TK_SYMBOL)             
                  break
            } 

            //
            // Rule 1127:  option ::= no xreference$x optionWhite
            //
            case 1127: { 
              my.MakeToken(my.GetRhsFirstTokenIndex(2), my.GetRhsLastTokenIndex(2), LPGParsersym.TK_SYMBOL)             
                  break
            } 

    //#line 285 "LexerTemplateF.gi

    
            default:
                break
        }
        return
    }

    //#line 3 "LexerBasicMapF.gi
 
type  LPGLexerLpgLexStream struct{
    *LexStream
}
func  NewLPGLexerLpgLexStream( fileName string, inputChars []rune, tab int) (*LPGLexerLpgLexStream,error){
        t := new(LPGLexerLpgLexStream)
        var e error
        t.LexStream ,e = NewLexStream(fileName,inputChars,tab,nil)
        return t,e
}

var LPGLexerLpgLexStreamtokenKind =[]int{
        LPGLexersym.Char_CtlCharNotWS,    // 000    0x00
        LPGLexersym.Char_CtlCharNotWS,    // 001    0x01
        LPGLexersym.Char_CtlCharNotWS,    // 002    0x02
        LPGLexersym.Char_CtlCharNotWS,    // 003    0x03
        LPGLexersym.Char_CtlCharNotWS,    // 004    0x04
        LPGLexersym.Char_CtlCharNotWS,    // 005    0x05
        LPGLexersym.Char_CtlCharNotWS,    // 006    0x06
        LPGLexersym.Char_CtlCharNotWS,    // 007    0x07
        LPGLexersym.Char_CtlCharNotWS,    // 008    0x08
        LPGLexersym.Char_HT,              // 009    0x09
        LPGLexersym.Char_LF,              // 010    0x0A
        LPGLexersym.Char_CtlCharNotWS,    // 011    0x0B
        LPGLexersym.Char_FF,              // 012    0x0C
        LPGLexersym.Char_CR,              // 013    0x0D
        LPGLexersym.Char_CtlCharNotWS,    // 014    0x0E
        LPGLexersym.Char_CtlCharNotWS,    // 015    0x0F
        LPGLexersym.Char_CtlCharNotWS,    // 016    0x10
        LPGLexersym.Char_CtlCharNotWS,    // 017    0x11
        LPGLexersym.Char_CtlCharNotWS,    // 018    0x12
        LPGLexersym.Char_CtlCharNotWS,    // 019    0x13
        LPGLexersym.Char_CtlCharNotWS,    // 020    0x14
        LPGLexersym.Char_CtlCharNotWS,    // 021    0x15
        LPGLexersym.Char_CtlCharNotWS,    // 022    0x16
        LPGLexersym.Char_CtlCharNotWS,    // 023    0x17
        LPGLexersym.Char_CtlCharNotWS,    // 024    0x18
        LPGLexersym.Char_CtlCharNotWS,    // 025    0x19
        LPGLexersym.Char_CtlCharNotWS,    // 026    0x1A
        LPGLexersym.Char_CtlCharNotWS,    // 027    0x1B
        LPGLexersym.Char_CtlCharNotWS,    // 028    0x1C
        LPGLexersym.Char_CtlCharNotWS,    // 029    0x1D
        LPGLexersym.Char_CtlCharNotWS,    // 030    0x1E
        LPGLexersym.Char_CtlCharNotWS,    // 031    0x1F
        LPGLexersym.Char_Space,           // 032    0x20
        LPGLexersym.Char_Exclamation,     // 033    0x21
        LPGLexersym.Char_DoubleQuote,     // 034    0x22
        LPGLexersym.Char_Sharp,           // 035    0x23
        LPGLexersym.Char_DollarSign,      // 036    0x24
        LPGLexersym.Char_Percent,         // 037    0x25
        LPGLexersym.Char_Ampersand,       // 038    0x26
        LPGLexersym.Char_SingleQuote,     // 039    0x27
        LPGLexersym.Char_LeftParen,       // 040    0x28
        LPGLexersym.Char_RightParen,      // 041    0x29
        LPGLexersym.Char_Star,            // 042    0x2A
        LPGLexersym.Char_Plus,            // 043    0x2B
        LPGLexersym.Char_Comma,           // 044    0x2C
        LPGLexersym.Char_Minus,           // 045    0x2D
        LPGLexersym.Char_Dot,             // 046    0x2E
        LPGLexersym.Char_Slash,           // 047    0x2F
        LPGLexersym.Char_0,               // 048    0x30
        LPGLexersym.Char_1,               // 049    0x31
        LPGLexersym.Char_2,               // 050    0x32
        LPGLexersym.Char_3,               // 051    0x33
        LPGLexersym.Char_4,               // 052    0x34
        LPGLexersym.Char_5,               // 053    0x35
        LPGLexersym.Char_6,               // 054    0x36
        LPGLexersym.Char_7,               // 055    0x37
        LPGLexersym.Char_8,               // 056    0x38
        LPGLexersym.Char_9,               // 057    0x39
        LPGLexersym.Char_Colon,           // 058    0x3A
        LPGLexersym.Char_SemiColon,       // 059    0x3B
        LPGLexersym.Char_LessThan,        // 060    0x3C
        LPGLexersym.Char_Equal,           // 061    0x3D
        LPGLexersym.Char_GreaterThan,     // 062    0x3E
        LPGLexersym.Char_QuestionMark,    // 063    0x3F
        LPGLexersym.Char_AtSign,          // 064    0x40
        LPGLexersym.Char_A,               // 065    0x41
        LPGLexersym.Char_B,               // 066    0x42
        LPGLexersym.Char_C,               // 067    0x43
        LPGLexersym.Char_D,               // 068    0x44
        LPGLexersym.Char_E,               // 069    0x45
        LPGLexersym.Char_F,               // 070    0x46
        LPGLexersym.Char_G,               // 071    0x47
        LPGLexersym.Char_H,               // 072    0x48
        LPGLexersym.Char_I,               // 073    0x49
        LPGLexersym.Char_J,               // 074    0x4A
        LPGLexersym.Char_K,               // 075    0x4B
        LPGLexersym.Char_L,               // 076    0x4C
        LPGLexersym.Char_M,               // 077    0x4D
        LPGLexersym.Char_N,               // 078    0x4E
        LPGLexersym.Char_O,               // 079    0x4F
        LPGLexersym.Char_P,               // 080    0x50
        LPGLexersym.Char_Q,               // 081    0x51
        LPGLexersym.Char_R,               // 082    0x52
        LPGLexersym.Char_S,               // 083    0x53
        LPGLexersym.Char_T,               // 084    0x54
        LPGLexersym.Char_U,               // 085    0x55
        LPGLexersym.Char_V,               // 086    0x56
        LPGLexersym.Char_W,               // 087    0x57
        LPGLexersym.Char_X,               // 088    0x58
        LPGLexersym.Char_Y,               // 089    0x59
        LPGLexersym.Char_Z,               // 090    0x5A
        LPGLexersym.Char_LeftBracket,     // 091    0x5B
        LPGLexersym.Char_BackSlash,       // 092    0x5C
        LPGLexersym.Char_RightBracket,    // 093    0x5D
        LPGLexersym.Char_Caret,           // 094    0x5E
        LPGLexersym.Char__,               // 095    0x5F
        LPGLexersym.Char_BackQuote,       // 096    0x60
        LPGLexersym.Char_a,               // 097    0x61
        LPGLexersym.Char_b,               // 098    0x62
        LPGLexersym.Char_c,               // 099    0x63
        LPGLexersym.Char_d,               // 100    0x64
        LPGLexersym.Char_e,               // 101    0x65
        LPGLexersym.Char_f,               // 102    0x66
        LPGLexersym.Char_g,               // 103    0x67
        LPGLexersym.Char_h,               // 104    0x68
        LPGLexersym.Char_i,               // 105    0x69
        LPGLexersym.Char_j,               // 106    0x6A
        LPGLexersym.Char_k,               // 107    0x6B
        LPGLexersym.Char_l,               // 108    0x6C
        LPGLexersym.Char_m,               // 109    0x6D
        LPGLexersym.Char_n,               // 110    0x6E
        LPGLexersym.Char_o,               // 111    0x6F
        LPGLexersym.Char_p,               // 112    0x70
        LPGLexersym.Char_q,               // 113    0x71
        LPGLexersym.Char_r,               // 114    0x72
        LPGLexersym.Char_s,               // 115    0x73
        LPGLexersym.Char_t,               // 116    0x74
        LPGLexersym.Char_u,               // 117    0x75
        LPGLexersym.Char_v,               // 118    0x76
        LPGLexersym.Char_w,               // 119    0x77
        LPGLexersym.Char_x,               // 120    0x78
        LPGLexersym.Char_y,               // 121    0x79
        LPGLexersym.Char_z,               // 122    0x7A
        LPGLexersym.Char_LeftBrace,       // 123    0x7B
        LPGLexersym.Char_VerticalBar,     // 124    0x7C
        LPGLexersym.Char_RightBrace,      // 125    0x7D
        LPGLexersym.Char_Tilde,           // 126    0x7E

        LPGLexersym.Char_AfterASCII,      // for all chars in range 128..65534
        LPGLexersym.Char_EOF ,             // for '\uffff' or 65535 
}

    
func (my *LPGLexerLpgLexStream) GetKind(i int) int {  // Classify character at ith location
        var c int 
        if i >= my.GetStreamLength() {
            c = 0xffff
        }else{
            c = my.GetIntValue(i)
        }
        if c < 128 {//ASCII Character
            return LPGLexerLpgLexStreamtokenKind[c]
        }else{
            if c == 0xffff{
                return LPGLexersym.Char_EOF
            }else{
                return LPGLexersym.Char_AfterASCII
            }
        }

}
func (my *LPGLexerLpgLexStream)  OrderedExportedSymbols() []string{
    return LPGParsersym.OrderedTerminalSymbols 
}   

