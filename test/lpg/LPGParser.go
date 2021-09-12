package lpg


    //#line 131 "dtParserTemplateF.gi

import (
. "github.com/A-LPG/LPG-go-runtime/lpg2"
)

    //#line 8 "LPGParser.g


 
    //#line 139 "dtParserTemplateF.gi

type  LPGParser  struct{
    prsStream   *PrsStream
    dtParser *DeterministicParser
    unimplementedSymbolsWarning  bool
    prsTable  ParseTable
}
func NewLPGParser(lexStream ILexStream) (*LPGParser,error){
    my := new(LPGParser)
    my.prsTable = &LPGParserprs{}
    my.prsStream =  NewPrsStream(lexStream)
    my.unimplementedSymbolsWarning  = false
    var e error
    my.dtParser,e =  NewDeterministicParser(nil, my.prsTable,  my,nil) 
    if e == nil{
        if lexStream != nil{
            err := my.Reset(lexStream)
            if err != nil {
                return nil, err
            }
        }
        return  my,nil
    }
    var _,ok = e.(*NotDeterministicParseTableException)
    if ok{
        e = NewNotDeterministicParseTableException("Regenerate %prs_type.go with -NOBACKTRACK option")
        return  nil,e
    }
    _,ok = e.(*BadParseSymFileException)
    if ok{
        e= NewBadParseSymFileException("Bad Parser Symbol File -- %sym_type.go Regenerate %prs_type.go")
        return  nil,e
    }
    return nil,e
}
    

func (my *LPGParser)  GetParseTable() ParseTable{ 
    return my.prsTable 
}

func (my *LPGParser)  GetParser() *DeterministicParser{
    return my.dtParser 
}

func (my *LPGParser)  SetResult( object interface{}) {
    my.dtParser.SetSym1(object)
}
func (my *LPGParser) GetRhsSym(i int) interface{} { 
    return my.dtParser.GetSym(i) 
}

func (my *LPGParser)  GetRhsTokenIndex(i int)int {
    return my.dtParser.GetToken(i)
    }
func (my *LPGParser)  GetRhsIToken(i int) IToken{
    return my.prsStream.GetIToken(my.GetRhsTokenIndex(i)) 
}

func (my *LPGParser)  GetRhsFirstTokenIndex(i int)int { 
return my.dtParser.GetFirstTokenAt(i)
    }
func (my *LPGParser)  GetRhsFirstIToken(i int) IToken{
    return my.prsStream.GetIToken(my.GetRhsFirstTokenIndex(i)) 
    }

func (my *LPGParser)  GetRhsLastTokenIndex(i int)int {
    return my.dtParser.GetLastTokenAt(i) 
    }
func (my *LPGParser)  GetRhsLastIToken(i int) IToken{ 
return my.prsStream.GetIToken(my.GetRhsLastTokenIndex(i)) 
}

func (my *LPGParser)  GetLeftSpan() int{ 
    return my.dtParser.GetFirstToken() 
}
func (my *LPGParser)  GetLeftIToken() IToken {
    return my.prsStream.GetIToken(my.GetLeftSpan())
}

func (my *LPGParser)  GetRightSpan() int{
    return my.dtParser.GetLastToken() 
}
func (my *LPGParser)  GetRightIToken()IToken {
    return my.prsStream.GetIToken(my.GetRightSpan())
}

func (my *LPGParser)  GetRhsErrorTokenIndex(i int)int{
    var index = my.dtParser.GetToken(i)
    var err = my.prsStream.GetIToken(index)
    var _,ok = err.(*ErrorToken)
    if ok {
        return index
    }else{
        return 0
    }
}
func (my *LPGParser)  GetRhsErrorIToken(i int) *ErrorToken{
    var index = my.dtParser.GetToken(i)
    var err = my.prsStream.GetIToken(index)
    var token,_ = err.(*ErrorToken)
    return token
}

func (my *LPGParser)  Reset(lexStream ILexStream ) error{
    my.prsStream = NewPrsStream(lexStream)
    err := my.dtParser.Reset(my.prsStream,nil,nil,nil)
    if err != nil {
        return err
    }
    var ex = my.prsStream.RemapTerminalSymbols(my.OrderedTerminalSymbols(), my.prsTable.GetEoftSymbol())
    if ex == nil{
        return nil
    }
    var _,ok = ex.(*NullExportedSymbolsException)
    if ok {
        return ex
    }
    _,ok = ex.(*NullTerminalSymbolsException)
    if ok {
        return ex
    }
    var e *UnimplementedTerminalsException
    e,ok = ex.(*UnimplementedTerminalsException)
    if ok {
        if my.unimplementedSymbolsWarning {
            var unimplemented_symbols = e.GetSymbols()
            println("The Lexer will not scan the following token(s):")
            var i int = 0
            for ; i < unimplemented_symbols.Size() ;i++{
                var id = unimplemented_symbols.Get(i)
                println("    " + LPGParsersym.OrderedTerminalSymbols[id])
            }
            println()
        }
        return  ex
    }
    _,ok = ex.(*UndefinedEofSymbolException)
    if ok {
        return NewUndefinedEofSymbolException("The Lexer does not implement the Eof symbol " +
        LPGParsersym.OrderedTerminalSymbols[my.prsTable.GetEoftSymbol()])
    }
    return ex
}


func (my *LPGParser)  NumTokenKinds()int {
        return LPGParsersym.NumTokenKinds 
}
func (my *LPGParser)  OrderedTerminalSymbols()[]string {
    return LPGParsersym.OrderedTerminalSymbols
}
func (my *LPGParser)  GetTokenKindName(kind int) string{
        return LPGParsersym.OrderedTerminalSymbols[kind] 
}
func (my *LPGParser)  GetEOFTokenKind() int{
    return my.prsTable.GetEoftSymbol()
    }
func (my *LPGParser)  GetIPrsStream() IPrsStream{
    return my.prsStream
    }
func (my *LPGParser) Parser() (interface{}, error) {
    return my.ParserWithMonitor(0,nil)
}
func (my *LPGParser) ParserWithMonitor(error_repair_count int ,  monitor Monitor) (interface{}, error){

    my.dtParser.SetMonitor(monitor)
    
    var ast,ex= my.dtParser.ParseEntry(0)
    if ex == nil{
        return ast,ex
    }
    var e,ok= ex.(*BadParseException)
    if ok{
        my.prsStream.ResetTo(e.ErrorToken) // point to error token

        var diagnoseParser = NewDiagnoseParser(my.prsStream, my.prsTable,0,0,nil)
        diagnoseParser.Diagnose(e.ErrorToken)
    }
    return ast,ex
}
//
// Additional entry points, if any
//


    //#line 38 "LPGParser.g


 
    //#line 45 "LPGParser.g


 
    //#line 222 "LPGParser.g


    //#line 329 "dtParserTemplateF.gi

   
    func (my *LPGParser)  RuleAction(ruleNumber int){
        switch ruleNumber{

            //
            // Rule 1:  LPG ::= options_segment LPG_INPUT
            //
            case 1: {
               //#line 44 "LPGParser.g"
                my.SetResult(
                    //#line 44 LPGParser.g
                    NewLPG(my, my.GetLeftIToken(), my.GetRightIToken(),
                           //#line 44 LPGParser.g
                           my.GetRhsSym(1).(*option_specList),
                           //#line 44 LPGParser.g
                           my.GetRhsSym(2).(*LPG_itemList)),
                //#line 44 LPGParser.g
                )
                break
            }
            //
            // Rule 2:  LPG_INPUT ::= $Empty
            //
            case 2: {
               //#line 49 "LPGParser.g"
                my.SetResult(
                    //#line 49 LPGParser.g
                    NewLPG_itemList(my.GetLeftIToken(), my.GetRightIToken(), true /* left recursive */),
                //#line 49 LPGParser.g
                )
                break
            }
            //
            // Rule 3:  LPG_INPUT ::= LPG_INPUT LPG_item
            //
            case 3: {
               //#line 50 "LPGParser.g"
                (my.GetRhsSym(1).(*LPG_itemList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 4:  LPG_item ::= ALIAS_KEY$ alias_segment END_KEY_OPT$
            //
            case 4: {
               //#line 53 "LPGParser.g"
                my.SetResult(
                    //#line 53 LPGParser.g
                    NewAliasSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 53 LPGParser.g
                                my.GetRhsSym(2).(*aliasSpecList)),
                //#line 53 LPGParser.g
                )
                break
            }
            //
            // Rule 5:  LPG_item ::= AST_KEY$ ast_segment END_KEY_OPT$
            //
            case 5: {
               //#line 54 "LPGParser.g"
                my.SetResult(
                    //#line 54 LPGParser.g
                    NewAstSeg(my.GetLeftIToken(), my.GetRightIToken(),
                              //#line 54 LPGParser.g
                              my.GetRhsSym(2).(*action_segmentList)),
                //#line 54 LPGParser.g
                )
                break
            }
            //
            // Rule 6:  LPG_item ::= DEFINE_KEY$ define_segment END_KEY_OPT$
            //
            case 6: {
               //#line 55 "LPGParser.g"
                my.SetResult(
                    //#line 55 LPGParser.g
                    NewDefineSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                 //#line 55 LPGParser.g
                                 my.GetRhsSym(2).(*defineSpecList)),
                //#line 55 LPGParser.g
                )
                break
            }
            //
            // Rule 7:  LPG_item ::= EOF_KEY$ eof_segment END_KEY_OPT$
            //
            case 7: {
               //#line 56 "LPGParser.g"
                my.SetResult(
                    //#line 56 LPGParser.g
                    NewEofSeg(my.GetLeftIToken(), my.GetRightIToken(),
                              //#line 56 LPGParser.g
                              my.GetRhsSym(2).(Ieof_segment)),
                //#line 56 LPGParser.g
                )
                break
            }
            //
            // Rule 8:  LPG_item ::= EOL_KEY$ eol_segment END_KEY_OPT$
            //
            case 8: {
               //#line 57 "LPGParser.g"
                my.SetResult(
                    //#line 57 LPGParser.g
                    NewEolSeg(my.GetLeftIToken(), my.GetRightIToken(),
                              //#line 57 LPGParser.g
                              my.GetRhsSym(2).(Ieol_segment)),
                //#line 57 LPGParser.g
                )
                break
            }
            //
            // Rule 9:  LPG_item ::= ERROR_KEY$ error_segment END_KEY_OPT$
            //
            case 9: {
               //#line 58 "LPGParser.g"
                my.SetResult(
                    //#line 58 LPGParser.g
                    NewErrorSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 58 LPGParser.g
                                my.GetRhsSym(2).(Ierror_segment)),
                //#line 58 LPGParser.g
                )
                break
            }
            //
            // Rule 10:  LPG_item ::= EXPORT_KEY$ export_segment END_KEY_OPT$
            //
            case 10: {
               //#line 59 "LPGParser.g"
                my.SetResult(
                    //#line 59 LPGParser.g
                    NewExportSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                 //#line 59 LPGParser.g
                                 my.GetRhsSym(2).(*terminal_symbolList)),
                //#line 59 LPGParser.g
                )
                break
            }
            //
            // Rule 11:  LPG_item ::= GLOBALS_KEY$ globals_segment END_KEY_OPT$
            //
            case 11: {
               //#line 60 "LPGParser.g"
                my.SetResult(
                    //#line 60 LPGParser.g
                    NewGlobalsSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                  //#line 60 LPGParser.g
                                  my.GetRhsSym(2).(*action_segmentList)),
                //#line 60 LPGParser.g
                )
                break
            }
            //
            // Rule 12:  LPG_item ::= HEADERS_KEY$ headers_segment END_KEY_OPT$
            //
            case 12: {
               //#line 61 "LPGParser.g"
                my.SetResult(
                    //#line 61 LPGParser.g
                    NewHeadersSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                  //#line 61 LPGParser.g
                                  my.GetRhsSym(2).(*action_segmentList)),
                //#line 61 LPGParser.g
                )
                break
            }
            //
            // Rule 13:  LPG_item ::= IDENTIFIER_KEY$ identifier_segment END_KEY_OPT$
            //
            case 13: {
               //#line 62 "LPGParser.g"
                my.SetResult(
                    //#line 62 LPGParser.g
                    NewIdentifierSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                     //#line 62 LPGParser.g
                                     my.GetRhsSym(2).(Iidentifier_segment)),
                //#line 62 LPGParser.g
                )
                break
            }
            //
            // Rule 14:  LPG_item ::= IMPORT_KEY$ import_segment END_KEY_OPT$
            //
            case 14: {
               //#line 63 "LPGParser.g"
                my.SetResult(
                    //#line 63 LPGParser.g
                    NewImportSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                 //#line 63 LPGParser.g
                                 my.GetRhsSym(2).(*import_segment)),
                //#line 63 LPGParser.g
                )
                break
            }
            //
            // Rule 15:  LPG_item ::= INCLUDE_KEY$ include_segment END_KEY_OPT$
            //
            case 15: {
               //#line 64 "LPGParser.g"
                my.SetResult(
                    //#line 64 LPGParser.g
                    NewIncludeSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                  //#line 64 LPGParser.g
                                  my.GetRhsSym(2).(*include_segment)),
                //#line 64 LPGParser.g
                )
                break
            }
            //
            // Rule 16:  LPG_item ::= KEYWORDS_KEY$ keywords_segment END_KEY_OPT$
            //
            case 16: {
               //#line 65 "LPGParser.g"
                my.SetResult(
                    //#line 65 LPGParser.g
                    NewKeywordsSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                   //#line 65 LPGParser.g
                                   my.GetRhsSym(2).(*keywordSpecList)),
                //#line 65 LPGParser.g
                )
                break
            }
            //
            // Rule 17:  LPG_item ::= NAMES_KEY$ names_segment END_KEY_OPT$
            //
            case 17: {
               //#line 66 "LPGParser.g"
                my.SetResult(
                    //#line 66 LPGParser.g
                    NewNamesSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 66 LPGParser.g
                                my.GetRhsSym(2).(*nameSpecList)),
                //#line 66 LPGParser.g
                )
                break
            }
            //
            // Rule 18:  LPG_item ::= NOTICE_KEY$ notice_segment END_KEY_OPT$
            //
            case 18: {
               //#line 67 "LPGParser.g"
                my.SetResult(
                    //#line 67 LPGParser.g
                    NewNoticeSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                 //#line 67 LPGParser.g
                                 my.GetRhsSym(2).(*action_segmentList)),
                //#line 67 LPGParser.g
                )
                break
            }
            //
            // Rule 19:  LPG_item ::= RULES_KEY$ rules_segment END_KEY_OPT$
            //
            case 19: {
               //#line 68 "LPGParser.g"
                my.SetResult(
                    //#line 68 LPGParser.g
                    NewRulesSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 68 LPGParser.g
                                my.GetRhsSym(2).(*rules_segment)),
                //#line 68 LPGParser.g
                )
                break
            }
            //
            // Rule 20:  LPG_item ::= SOFT_KEYWORDS_KEY$ keywords_segment END_KEY_OPT$
            //
            case 20: {
               //#line 69 "LPGParser.g"
                my.SetResult(
                    //#line 69 LPGParser.g
                    NewSoftKeywordsSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                       //#line 69 LPGParser.g
                                       my.GetRhsSym(2).(*keywordSpecList)),
                //#line 69 LPGParser.g
                )
                break
            }
            //
            // Rule 21:  LPG_item ::= START_KEY$ start_segment END_KEY_OPT$
            //
            case 21: {
               //#line 70 "LPGParser.g"
                my.SetResult(
                    //#line 70 LPGParser.g
                    NewStartSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 70 LPGParser.g
                                my.GetRhsSym(2).(*start_symbolList)),
                //#line 70 LPGParser.g
                )
                break
            }
            //
            // Rule 22:  LPG_item ::= TERMINALS_KEY$ terminals_segment END_KEY_OPT$
            //
            case 22: {
               //#line 71 "LPGParser.g"
                my.SetResult(
                    //#line 71 LPGParser.g
                    NewTerminalsSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                    //#line 71 LPGParser.g
                                    my.GetRhsSym(2).(*terminalList)),
                //#line 71 LPGParser.g
                )
                break
            }
            //
            // Rule 23:  LPG_item ::= TRAILERS_KEY$ trailers_segment END_KEY_OPT$
            //
            case 23: {
               //#line 72 "LPGParser.g"
                my.SetResult(
                    //#line 72 LPGParser.g
                    NewTrailersSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                   //#line 72 LPGParser.g
                                   my.GetRhsSym(2).(*action_segmentList)),
                //#line 72 LPGParser.g
                )
                break
            }
            //
            // Rule 24:  LPG_item ::= TYPES_KEY$ types_segment END_KEY_OPT$
            //
            case 24: {
               //#line 73 "LPGParser.g"
                my.SetResult(
                    //#line 73 LPGParser.g
                    NewTypesSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 73 LPGParser.g
                                my.GetRhsSym(2).(*type_declarationsList)),
                //#line 73 LPGParser.g
                )
                break
            }
            //
            // Rule 25:  LPG_item ::= RECOVER_KEY$ recover_segment END_KEY_OPT$
            //
            case 25: {
               //#line 74 "LPGParser.g"
                my.SetResult(
                    //#line 74 LPGParser.g
                    NewRecoverSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                  //#line 74 LPGParser.g
                                  my.GetRhsSym(2).(*SYMBOLList)),
                //#line 74 LPGParser.g
                )
                break
            }
            //
            // Rule 26:  LPG_item ::= DISJOINTPREDECESSORSETS_KEY$ predecessor_segment END_KEY_OPT$
            //
            case 26: {
               //#line 75 "LPGParser.g"
                my.SetResult(
                    //#line 75 LPGParser.g
                    NewPredecessorSeg(my.GetLeftIToken(), my.GetRightIToken(),
                                      //#line 75 LPGParser.g
                                      my.GetRhsSym(2).(*symbol_pairList)),
                //#line 75 LPGParser.g
                )
                break
            }
            //
            // Rule 27:  options_segment ::= $Empty
            //
            case 27: {
               //#line 78 "LPGParser.g"
                my.SetResult(
                    //#line 78 LPGParser.g
                    Newoption_specList(my.GetLeftIToken(), my.GetRightIToken(), true /* left recursive */),
                //#line 78 LPGParser.g
                )
                break
            }
            //
            // Rule 28:  options_segment ::= options_segment option_spec
            //
            case 28: {
               //#line 78 "LPGParser.g"
                (my.GetRhsSym(1).(*option_specList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 29:  option_spec ::= OPTIONS_KEY$ option_list
            //
            case 29: {
               //#line 79 "LPGParser.g"
                my.SetResult(
                    //#line 79 LPGParser.g
                    Newoption_spec(my.GetLeftIToken(), my.GetRightIToken(),
                                   //#line 79 LPGParser.g
                                   my.GetRhsSym(2).(*optionList)),
                //#line 79 LPGParser.g
                )
                break
            }
            //
            // Rule 30:  option_list ::= option
            //
            case 30: {
               //#line 80 "LPGParser.g"
                my.SetResult(
                    //#line 80 LPGParser.g
                    NewoptionListFromElement(my.GetRhsSym(1).(*option), true /* left recursive */),
                //#line 80 LPGParser.g
                )
                break
            }
            //
            // Rule 31:  option_list ::= option_list ,$ option
            //
            case 31: {
               //#line 80 "LPGParser.g"
                (my.GetRhsSym(1).(*optionList)).AddElement(my.GetRhsSym(3).(IAst))
                break
            }
            //
            // Rule 32:  option ::= SYMBOL option_value
            //
            case 32: {
               //#line 81 "LPGParser.g"
                my.SetResult(
                    //#line 81 LPGParser.g
                    Newoption(my.GetLeftIToken(), my.GetRightIToken(),
                              //#line 81 LPGParser.g
                              NewASTNodeToken(my.GetRhsIToken(1)),
                              //#line 81 LPGParser.g
                              AnyCastToIoption_value(my.GetRhsSym(2))),
                //#line 81 LPGParser.g
                )
                break
            }
            //
            // Rule 33:  option_value ::= $Empty
            //
            case 33: {
               //#line 82 "LPGParser.g"
                my.SetResult(nil);
                break
            }
            //
            // Rule 34:  option_value ::= =$ SYMBOL
            //
            case 34: {
               //#line 82 "LPGParser.g"
                my.SetResult(
                    //#line 82 LPGParser.g
                    Newoption_value0(my.GetLeftIToken(), my.GetRightIToken(),
                                     //#line 82 LPGParser.g
                                     NewASTNodeToken(my.GetRhsIToken(2))),
                //#line 82 LPGParser.g
                )
                break
            }
            //
            // Rule 35:  option_value ::= =$ ($ symbol_list )$
            //
            case 35: {
               //#line 82 "LPGParser.g"
                my.SetResult(
                    //#line 82 LPGParser.g
                    Newoption_value1(my.GetLeftIToken(), my.GetRightIToken(),
                                     //#line 82 LPGParser.g
                                     my.GetRhsSym(3).(*SYMBOLList)),
                //#line 82 LPGParser.g
                )
                break
            }
            //
            // Rule 36:  symbol_list ::= SYMBOL
            //
            case 36: {
               //#line 84 "LPGParser.g"
                my.SetResult(
                    //#line 84 LPGParser.g
                    NewSYMBOLListFromElement(NewASTNodeToken(my.GetRhsIToken(1)), true /* left recursive */),
                //#line 84 LPGParser.g
                )
                break
            }
            //
            // Rule 37:  symbol_list ::= symbol_list ,$ SYMBOL
            //
            case 37: {
               //#line 85 "LPGParser.g"
                (my.GetRhsSym(1).(*SYMBOLList)).AddElement(NewASTNodeToken(my.GetRhsIToken(3)))
                break
            }
            //
            // Rule 38:  alias_segment ::= aliasSpec
            //
            case 38: {
               //#line 88 "LPGParser.g"
                my.SetResult(
                    //#line 88 LPGParser.g
                    NewaliasSpecListFromElement(my.GetRhsSym(1).(IaliasSpec), true /* left recursive */),
                //#line 88 LPGParser.g
                )
                break
            }
            //
            // Rule 39:  alias_segment ::= alias_segment aliasSpec
            //
            case 39: {
               //#line 88 "LPGParser.g"
                (my.GetRhsSym(1).(*aliasSpecList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 40:  aliasSpec ::= ERROR_KEY produces alias_rhs
            //
            case 40: {
               //#line 90 "LPGParser.g"
                my.SetResult(
                    //#line 90 LPGParser.g
                    NewaliasSpec0(my.GetLeftIToken(), my.GetRightIToken(),
                                  //#line 90 LPGParser.g
                                  NewASTNodeToken(my.GetRhsIToken(1)),
                                  //#line 90 LPGParser.g
                                  my.GetRhsSym(2).(Iproduces),
                                  //#line 90 LPGParser.g
                                  my.GetRhsSym(3).(Ialias_rhs)),
                //#line 90 LPGParser.g
                )
                break
            }
            //
            // Rule 41:  aliasSpec ::= EOL_KEY produces alias_rhs
            //
            case 41: {
               //#line 91 "LPGParser.g"
                my.SetResult(
                    //#line 91 LPGParser.g
                    NewaliasSpec1(my.GetLeftIToken(), my.GetRightIToken(),
                                  //#line 91 LPGParser.g
                                  NewASTNodeToken(my.GetRhsIToken(1)),
                                  //#line 91 LPGParser.g
                                  my.GetRhsSym(2).(Iproduces),
                                  //#line 91 LPGParser.g
                                  my.GetRhsSym(3).(Ialias_rhs)),
                //#line 91 LPGParser.g
                )
                break
            }
            //
            // Rule 42:  aliasSpec ::= EOF_KEY produces alias_rhs
            //
            case 42: {
               //#line 92 "LPGParser.g"
                my.SetResult(
                    //#line 92 LPGParser.g
                    NewaliasSpec2(my.GetLeftIToken(), my.GetRightIToken(),
                                  //#line 92 LPGParser.g
                                  NewASTNodeToken(my.GetRhsIToken(1)),
                                  //#line 92 LPGParser.g
                                  my.GetRhsSym(2).(Iproduces),
                                  //#line 92 LPGParser.g
                                  my.GetRhsSym(3).(Ialias_rhs)),
                //#line 92 LPGParser.g
                )
                break
            }
            //
            // Rule 43:  aliasSpec ::= IDENTIFIER_KEY produces alias_rhs
            //
            case 43: {
               //#line 93 "LPGParser.g"
                my.SetResult(
                    //#line 93 LPGParser.g
                    NewaliasSpec3(my.GetLeftIToken(), my.GetRightIToken(),
                                  //#line 93 LPGParser.g
                                  NewASTNodeToken(my.GetRhsIToken(1)),
                                  //#line 93 LPGParser.g
                                  my.GetRhsSym(2).(Iproduces),
                                  //#line 93 LPGParser.g
                                  my.GetRhsSym(3).(Ialias_rhs)),
                //#line 93 LPGParser.g
                )
                break
            }
            //
            // Rule 44:  aliasSpec ::= SYMBOL produces alias_rhs
            //
            case 44: {
               //#line 94 "LPGParser.g"
                my.SetResult(
                    //#line 94 LPGParser.g
                    NewaliasSpec4(my.GetLeftIToken(), my.GetRightIToken(),
                                  //#line 94 LPGParser.g
                                  NewASTNodeToken(my.GetRhsIToken(1)),
                                  //#line 94 LPGParser.g
                                  my.GetRhsSym(2).(Iproduces),
                                  //#line 94 LPGParser.g
                                  my.GetRhsSym(3).(Ialias_rhs)),
                //#line 94 LPGParser.g
                )
                break
            }
            //
            // Rule 45:  aliasSpec ::= alias_lhs_macro_name produces alias_rhs
            //
            case 45: {
               //#line 95 "LPGParser.g"
                my.SetResult(
                    //#line 95 LPGParser.g
                    NewaliasSpec5(my.GetLeftIToken(), my.GetRightIToken(),
                                  //#line 95 LPGParser.g
                                  my.GetRhsSym(1).(*alias_lhs_macro_name),
                                  //#line 95 LPGParser.g
                                  my.GetRhsSym(2).(Iproduces),
                                  //#line 95 LPGParser.g
                                  my.GetRhsSym(3).(Ialias_rhs)),
                //#line 95 LPGParser.g
                )
                break
            }
            //
            // Rule 46:  alias_lhs_macro_name ::= MACRO_NAME
            //
            case 46: {
               //#line 97 "LPGParser.g"
                my.SetResult(
                    //#line 97 LPGParser.g
                    Newalias_lhs_macro_name(my.GetRhsIToken(1)),
                //#line 97 LPGParser.g
                )
                break
            }
            //
            // Rule 47:  alias_rhs ::= SYMBOL
            //
            case 47: {
               //#line 99 "LPGParser.g"
                my.SetResult(
                    //#line 99 LPGParser.g
                    Newalias_rhs0(my.GetRhsIToken(1)),
                //#line 99 LPGParser.g
                )
                break
            }
            //
            // Rule 48:  alias_rhs ::= MACRO_NAME
            //
            case 48: {
               //#line 100 "LPGParser.g"
                my.SetResult(
                    //#line 100 LPGParser.g
                    Newalias_rhs1(my.GetRhsIToken(1)),
                //#line 100 LPGParser.g
                )
                break
            }
            //
            // Rule 49:  alias_rhs ::= ERROR_KEY
            //
            case 49: {
               //#line 101 "LPGParser.g"
                my.SetResult(
                    //#line 101 LPGParser.g
                    Newalias_rhs2(my.GetRhsIToken(1)),
                //#line 101 LPGParser.g
                )
                break
            }
            //
            // Rule 50:  alias_rhs ::= EOL_KEY
            //
            case 50: {
               //#line 102 "LPGParser.g"
                my.SetResult(
                    //#line 102 LPGParser.g
                    Newalias_rhs3(my.GetRhsIToken(1)),
                //#line 102 LPGParser.g
                )
                break
            }
            //
            // Rule 51:  alias_rhs ::= EOF_KEY
            //
            case 51: {
               //#line 103 "LPGParser.g"
                my.SetResult(
                    //#line 103 LPGParser.g
                    Newalias_rhs4(my.GetRhsIToken(1)),
                //#line 103 LPGParser.g
                )
                break
            }
            //
            // Rule 52:  alias_rhs ::= EMPTY_KEY
            //
            case 52: {
               //#line 104 "LPGParser.g"
                my.SetResult(
                    //#line 104 LPGParser.g
                    Newalias_rhs5(my.GetRhsIToken(1)),
                //#line 104 LPGParser.g
                )
                break
            }
            //
            // Rule 53:  alias_rhs ::= IDENTIFIER_KEY
            //
            case 53: {
               //#line 105 "LPGParser.g"
                my.SetResult(
                    //#line 105 LPGParser.g
                    Newalias_rhs6(my.GetRhsIToken(1)),
                //#line 105 LPGParser.g
                )
                break
            }
            //
            // Rule 54:  ast_segment ::= action_segment_list
            //
            case 54:
                break
            //
            // Rule 55:  define_segment ::= defineSpec
            //
            case 55: {
               //#line 111 "LPGParser.g"
                my.SetResult(
                    //#line 111 LPGParser.g
                    NewdefineSpecListFromElement(my.GetRhsSym(1).(*defineSpec), true /* left recursive */),
                //#line 111 LPGParser.g
                )
                break
            }
            //
            // Rule 56:  define_segment ::= define_segment defineSpec
            //
            case 56: {
               //#line 111 "LPGParser.g"
                (my.GetRhsSym(1).(*defineSpecList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 57:  defineSpec ::= macro_name_symbol macro_segment
            //
            case 57: {
               //#line 112 "LPGParser.g"
                my.SetResult(
                    //#line 112 LPGParser.g
                    NewdefineSpec(my.GetLeftIToken(), my.GetRightIToken(),
                                  //#line 112 LPGParser.g
                                  my.GetRhsSym(1).(Imacro_name_symbol),
                                  //#line 112 LPGParser.g
                                  my.GetRhsSym(2).(*macro_segment)),
                //#line 112 LPGParser.g
                )
                break
            }
            //
            // Rule 58:  macro_name_symbol ::= MACRO_NAME
            //
            case 58: {
               //#line 115 "LPGParser.g"
                my.SetResult(
                    //#line 115 LPGParser.g
                    Newmacro_name_symbol0(my.GetRhsIToken(1)),
                //#line 115 LPGParser.g
                )
                break
            }
            //
            // Rule 59:  macro_name_symbol ::= SYMBOL
            //
            case 59: {
               //#line 116 "LPGParser.g"
                my.SetResult(
                    //#line 116 LPGParser.g
                    Newmacro_name_symbol1(my.GetRhsIToken(1)),
                //#line 116 LPGParser.g
                )
                break
            }
            //
            // Rule 60:  macro_segment ::= BLOCK
            //
            case 60: {
               //#line 117 "LPGParser.g"
                my.SetResult(
                    //#line 117 LPGParser.g
                    Newmacro_segment(my.GetRhsIToken(1)),
                //#line 117 LPGParser.g
                )
                break
            }
            //
            // Rule 61:  eol_segment ::= terminal_symbol
            //
            case 61:
                break
            //
            // Rule 62:  eof_segment ::= terminal_symbol
            //
            case 62:
                break
            //
            // Rule 63:  error_segment ::= terminal_symbol
            //
            case 63:
                break
            //
            // Rule 64:  export_segment ::= terminal_symbol
            //
            case 64: {
               //#line 127 "LPGParser.g"
                my.SetResult(
                    //#line 127 LPGParser.g
                    Newterminal_symbolListFromElement(my.GetRhsSym(1).(Iterminal_symbol), true /* left recursive */),
                //#line 127 LPGParser.g
                )
                break
            }
            //
            // Rule 65:  export_segment ::= export_segment terminal_symbol
            //
            case 65: {
               //#line 127 "LPGParser.g"
                (my.GetRhsSym(1).(*terminal_symbolList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 66:  globals_segment ::= action_segment
            //
            case 66: {
               //#line 130 "LPGParser.g"
                my.SetResult(
                    //#line 130 LPGParser.g
                    Newaction_segmentListFromElement(my.GetRhsSym(1).(*action_segment), true /* left recursive */),
                //#line 130 LPGParser.g
                )
                break
            }
            //
            // Rule 67:  globals_segment ::= globals_segment action_segment
            //
            case 67: {
               //#line 130 "LPGParser.g"
                (my.GetRhsSym(1).(*action_segmentList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 68:  headers_segment ::= action_segment_list
            //
            case 68:
                break
            //
            // Rule 69:  identifier_segment ::= terminal_symbol
            //
            case 69:
                break
            //
            // Rule 70:  import_segment ::= SYMBOL drop_command_list
            //
            case 70: {
               //#line 139 "LPGParser.g"
                my.SetResult(
                    //#line 139 LPGParser.g
                    Newimport_segment(my.GetLeftIToken(), my.GetRightIToken(),
                                      //#line 139 LPGParser.g
                                      NewASTNodeToken(my.GetRhsIToken(1)),
                                      //#line 139 LPGParser.g
                                      my.GetRhsSym(2).(*drop_commandList)),
                //#line 139 LPGParser.g
                )
                break
            }
            //
            // Rule 71:  drop_command_list ::= $Empty
            //
            case 71: {
               //#line 141 "LPGParser.g"
                my.SetResult(
                    //#line 141 LPGParser.g
                    Newdrop_commandList(my.GetLeftIToken(), my.GetRightIToken(), true /* left recursive */),
                //#line 141 LPGParser.g
                )
                break
            }
            //
            // Rule 72:  drop_command_list ::= drop_command_list drop_command
            //
            case 72: {
               //#line 141 "LPGParser.g"
                (my.GetRhsSym(1).(*drop_commandList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 73:  drop_command ::= DROPSYMBOLS_KEY drop_symbols
            //
            case 73: {
               //#line 143 "LPGParser.g"
                my.SetResult(
                    //#line 143 LPGParser.g
                    Newdrop_command0(my.GetLeftIToken(), my.GetRightIToken(),
                                     //#line 143 LPGParser.g
                                     NewASTNodeToken(my.GetRhsIToken(1)),
                                     //#line 143 LPGParser.g
                                     my.GetRhsSym(2).(*SYMBOLList)),
                //#line 143 LPGParser.g
                )
                break
            }
            //
            // Rule 74:  drop_command ::= DROPRULES_KEY drop_rules
            //
            case 74: {
               //#line 144 "LPGParser.g"
                my.SetResult(
                    //#line 144 LPGParser.g
                    Newdrop_command1(my.GetLeftIToken(), my.GetRightIToken(),
                                     //#line 144 LPGParser.g
                                     NewASTNodeToken(my.GetRhsIToken(1)),
                                     //#line 144 LPGParser.g
                                     my.GetRhsSym(2).(*drop_ruleList)),
                //#line 144 LPGParser.g
                )
                break
            }
            //
            // Rule 75:  drop_symbols ::= SYMBOL
            //
            case 75: {
               //#line 146 "LPGParser.g"
                my.SetResult(
                    //#line 146 LPGParser.g
                    NewSYMBOLListFromElement(NewASTNodeToken(my.GetRhsIToken(1)), true /* left recursive */),
                //#line 146 LPGParser.g
                )
                break
            }
            //
            // Rule 76:  drop_symbols ::= drop_symbols SYMBOL
            //
            case 76: {
               //#line 147 "LPGParser.g"
                (my.GetRhsSym(1).(*SYMBOLList)).AddElement(NewASTNodeToken(my.GetRhsIToken(2)))
                break
            }
            //
            // Rule 77:  drop_rules ::= drop_rule
            //
            case 77: {
               //#line 148 "LPGParser.g"
                my.SetResult(
                    //#line 148 LPGParser.g
                    Newdrop_ruleListFromElement(my.GetRhsSym(1).(*drop_rule), true /* left recursive */),
                //#line 148 LPGParser.g
                )
                break
            }
            //
            // Rule 78:  drop_rules ::= drop_rules drop_rule
            //
            case 78: {
               //#line 149 "LPGParser.g"
                (my.GetRhsSym(1).(*drop_ruleList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 79:  drop_rule ::= SYMBOL optMacroName produces ruleList
            //
            case 79: {
               //#line 151 "LPGParser.g"
                my.SetResult(
                    //#line 151 LPGParser.g
                    Newdrop_rule(my.GetLeftIToken(), my.GetRightIToken(),
                                 //#line 151 LPGParser.g
                                 NewASTNodeToken(my.GetRhsIToken(1)),
                                 //#line 151 LPGParser.g
                                 AnyCastTooptMacroName(my.GetRhsSym(2)),
                                 //#line 151 LPGParser.g
                                 my.GetRhsSym(3).(Iproduces),
                                 //#line 151 LPGParser.g
                                 my.GetRhsSym(4).(*ruleList)),
                //#line 151 LPGParser.g
                )
                break
            }
            //
            // Rule 80:  optMacroName ::= $Empty
            //
            case 80: {
               //#line 153 "LPGParser.g"
                my.SetResult(nil);
                break
            }
            //
            // Rule 81:  optMacroName ::= MACRO_NAME
            //
            case 81: {
               //#line 153 "LPGParser.g"
                my.SetResult(
                    //#line 153 LPGParser.g
                    NewoptMacroName(my.GetRhsIToken(1)),
                //#line 153 LPGParser.g
                )
                break
            }
            //
            // Rule 82:  include_segment ::= SYMBOL
            //
            case 82: {
               //#line 156 "LPGParser.g"
                my.SetResult(
                    //#line 156 LPGParser.g
                    Newinclude_segment(my.GetRhsIToken(1)),
                //#line 156 LPGParser.g
                )
                break
            }
            //
            // Rule 83:  keywords_segment ::= keywordSpec
            //
            case 83: {
               //#line 159 "LPGParser.g"
                my.SetResult(
                    //#line 159 LPGParser.g
                    NewkeywordSpecListFromElement(my.GetRhsSym(1).(IkeywordSpec), true /* left recursive */),
                //#line 159 LPGParser.g
                )
                break
            }
            //
            // Rule 84:  keywords_segment ::= keywords_segment keywordSpec
            //
            case 84: {
               //#line 159 "LPGParser.g"
                (my.GetRhsSym(1).(*keywordSpecList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 85:  keywordSpec ::= terminal_symbol
            //
            case 85:
                break
            //
            // Rule 86:  keywordSpec ::= terminal_symbol produces name
            //
            case 86: {
               //#line 161 "LPGParser.g"
                my.SetResult(
                    //#line 161 LPGParser.g
                    NewkeywordSpec(my.GetLeftIToken(), my.GetRightIToken(),
                                   //#line 161 LPGParser.g
                                   my.GetRhsSym(1).(Iterminal_symbol),
                                   //#line 161 LPGParser.g
                                   my.GetRhsSym(2).(Iproduces),
                                   //#line 161 LPGParser.g
                                   my.GetRhsSym(3).(Iname)),
                //#line 161 LPGParser.g
                )
                break
            }
            //
            // Rule 87:  names_segment ::= nameSpec
            //
            case 87: {
               //#line 164 "LPGParser.g"
                my.SetResult(
                    //#line 164 LPGParser.g
                    NewnameSpecListFromElement(my.GetRhsSym(1).(*nameSpec), true /* left recursive */),
                //#line 164 LPGParser.g
                )
                break
            }
            //
            // Rule 88:  names_segment ::= names_segment nameSpec
            //
            case 88: {
               //#line 164 "LPGParser.g"
                (my.GetRhsSym(1).(*nameSpecList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 89:  nameSpec ::= name produces name
            //
            case 89: {
               //#line 165 "LPGParser.g"
                my.SetResult(
                    //#line 165 LPGParser.g
                    NewnameSpec(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 165 LPGParser.g
                                my.GetRhsSym(1).(Iname),
                                //#line 165 LPGParser.g
                                my.GetRhsSym(2).(Iproduces),
                                //#line 165 LPGParser.g
                                my.GetRhsSym(3).(Iname)),
                //#line 165 LPGParser.g
                )
                break
            }
            //
            // Rule 90:  name ::= SYMBOL
            //
            case 90: {
               //#line 167 "LPGParser.g"
                my.SetResult(
                    //#line 167 LPGParser.g
                    Newname0(my.GetRhsIToken(1)),
                //#line 167 LPGParser.g
                )
                break
            }
            //
            // Rule 91:  name ::= MACRO_NAME
            //
            case 91: {
               //#line 168 "LPGParser.g"
                my.SetResult(
                    //#line 168 LPGParser.g
                    Newname1(my.GetRhsIToken(1)),
                //#line 168 LPGParser.g
                )
                break
            }
            //
            // Rule 92:  name ::= EMPTY_KEY
            //
            case 92: {
               //#line 169 "LPGParser.g"
                my.SetResult(
                    //#line 169 LPGParser.g
                    Newname2(my.GetRhsIToken(1)),
                //#line 169 LPGParser.g
                )
                break
            }
            //
            // Rule 93:  name ::= ERROR_KEY
            //
            case 93: {
               //#line 170 "LPGParser.g"
                my.SetResult(
                    //#line 170 LPGParser.g
                    Newname3(my.GetRhsIToken(1)),
                //#line 170 LPGParser.g
                )
                break
            }
            //
            // Rule 94:  name ::= EOL_KEY
            //
            case 94: {
               //#line 171 "LPGParser.g"
                my.SetResult(
                    //#line 171 LPGParser.g
                    Newname4(my.GetRhsIToken(1)),
                //#line 171 LPGParser.g
                )
                break
            }
            //
            // Rule 95:  name ::= IDENTIFIER_KEY
            //
            case 95: {
               //#line 172 "LPGParser.g"
                my.SetResult(
                    //#line 172 LPGParser.g
                    Newname5(my.GetRhsIToken(1)),
                //#line 172 LPGParser.g
                )
                break
            }
            //
            // Rule 96:  notice_segment ::= action_segment
            //
            case 96: {
               //#line 175 "LPGParser.g"
                my.SetResult(
                    //#line 175 LPGParser.g
                    Newaction_segmentListFromElement(my.GetRhsSym(1).(*action_segment), true /* left recursive */),
                //#line 175 LPGParser.g
                )
                break
            }
            //
            // Rule 97:  notice_segment ::= notice_segment action_segment
            //
            case 97: {
               //#line 175 "LPGParser.g"
                (my.GetRhsSym(1).(*action_segmentList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 98:  rules_segment ::= action_segment_list nonTermList
            //
            case 98: {
               //#line 178 "LPGParser.g"
                my.SetResult(
                    //#line 178 LPGParser.g
                    Newrules_segment(my.GetLeftIToken(), my.GetRightIToken(),
                                     //#line 178 LPGParser.g
                                     my.GetRhsSym(1).(*action_segmentList),
                                     //#line 178 LPGParser.g
                                     my.GetRhsSym(2).(*nonTermList)),
                //#line 178 LPGParser.g
                )
                break
            }
            //
            // Rule 99:  nonTermList ::= $Empty
            //
            case 99: {
               //#line 180 "LPGParser.g"
                my.SetResult(
                    //#line 180 LPGParser.g
                    NewnonTermList(my.GetLeftIToken(), my.GetRightIToken(), true /* left recursive */),
                //#line 180 LPGParser.g
                )
                break
            }
            //
            // Rule 100:  nonTermList ::= nonTermList nonTerm
            //
            case 100: {
               //#line 180 "LPGParser.g"
                (my.GetRhsSym(1).(*nonTermList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 101:  nonTerm ::= ruleNameWithAttributes produces ruleList
            //
            case 101: {
               //#line 182 "LPGParser.g"
                my.SetResult(
                    //#line 182 LPGParser.g
                    NewnonTerm(my.GetLeftIToken(), my.GetRightIToken(),
                               //#line 182 LPGParser.g
                               my.GetRhsSym(1).(*RuleName),
                               //#line 182 LPGParser.g
                               my.GetRhsSym(2).(Iproduces),
                               //#line 182 LPGParser.g
                               my.GetRhsSym(3).(*ruleList)),
                //#line 182 LPGParser.g
                )
                break
            }
            //
            // Rule 102:  ruleNameWithAttributes ::= SYMBOL
            //
            case 102: {
               //#line 185 "LPGParser.g"
                my.SetResult(
                    //#line 185 LPGParser.g
                    NewRuleName(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 185 LPGParser.g
                                NewASTNodeToken(my.GetRhsIToken(1)),
                                //#line 185 LPGParser.g
                                nil,
                                //#line 185 LPGParser.g
                                nil),
                //#line 185 LPGParser.g
                )
                break
            }
            //
            // Rule 103:  ruleNameWithAttributes ::= SYMBOL MACRO_NAME$className
            //
            case 103: {
               //#line 186 "LPGParser.g"
                my.SetResult(
                    //#line 186 LPGParser.g
                    NewRuleName(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 186 LPGParser.g
                                NewASTNodeToken(my.GetRhsIToken(1)),
                                //#line 186 LPGParser.g
                                NewASTNodeToken(my.GetRhsIToken(2)),
                                //#line 186 LPGParser.g
                                nil),
                //#line 186 LPGParser.g
                )
                break
            }
            //
            // Rule 104:  ruleNameWithAttributes ::= SYMBOL MACRO_NAME$className MACRO_NAME$arrayElement
            //
            case 104: {
               //#line 187 "LPGParser.g"
                my.SetResult(
                    //#line 187 LPGParser.g
                    NewRuleName(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 187 LPGParser.g
                                NewASTNodeToken(my.GetRhsIToken(1)),
                                //#line 187 LPGParser.g
                                NewASTNodeToken(my.GetRhsIToken(2)),
                                //#line 187 LPGParser.g
                                NewASTNodeToken(my.GetRhsIToken(3))),
                //#line 187 LPGParser.g
                )
                break
            }
            //
            // Rule 105:  ruleList ::= rule
            //
            case 105: {
               //#line 201 "LPGParser.g"
                my.SetResult(
                    //#line 201 LPGParser.g
                    NewruleListFromElement(my.GetRhsSym(1).(*rule), true /* left recursive */),
                //#line 201 LPGParser.g
                )
                break
            }
            //
            // Rule 106:  ruleList ::= ruleList |$ rule
            //
            case 106: {
               //#line 201 "LPGParser.g"
                (my.GetRhsSym(1).(*ruleList)).AddElement(my.GetRhsSym(3).(IAst))
                break
            }
            //
            // Rule 107:  produces ::= ::=
            //
            case 107: {
               //#line 203 "LPGParser.g"
                my.SetResult(
                    //#line 203 LPGParser.g
                    Newproduces0(my.GetRhsIToken(1)),
                //#line 203 LPGParser.g
                )
                break
            }
            //
            // Rule 108:  produces ::= ::=?
            //
            case 108: {
               //#line 204 "LPGParser.g"
                my.SetResult(
                    //#line 204 LPGParser.g
                    Newproduces1(my.GetRhsIToken(1)),
                //#line 204 LPGParser.g
                )
                break
            }
            //
            // Rule 109:  produces ::= ->
            //
            case 109: {
               //#line 205 "LPGParser.g"
                my.SetResult(
                    //#line 205 LPGParser.g
                    Newproduces2(my.GetRhsIToken(1)),
                //#line 205 LPGParser.g
                )
                break
            }
            //
            // Rule 110:  produces ::= ->?
            //
            case 110: {
               //#line 206 "LPGParser.g"
                my.SetResult(
                    //#line 206 LPGParser.g
                    Newproduces3(my.GetRhsIToken(1)),
                //#line 206 LPGParser.g
                )
                break
            }
            //
            // Rule 111:  rule ::= symWithAttrsList opt_action_segment
            //
            case 111: {
               //#line 208 "LPGParser.g"
                my.SetResult(
                    //#line 208 LPGParser.g
                    Newrule(my.GetLeftIToken(), my.GetRightIToken(),
                            //#line 208 LPGParser.g
                            my.GetRhsSym(1).(*symWithAttrsList),
                            //#line 208 LPGParser.g
                            AnyCastToaction_segment(my.GetRhsSym(2))),
                //#line 208 LPGParser.g
                )
                break
            }
            //
            // Rule 112:  symWithAttrsList ::= $Empty
            //
            case 112: {
               //#line 210 "LPGParser.g"
                my.SetResult(
                    //#line 210 LPGParser.g
                    NewsymWithAttrsList(my.GetLeftIToken(), my.GetRightIToken(), true /* left recursive */),
                //#line 210 LPGParser.g
                )
                break
            }
            //
            // Rule 113:  symWithAttrsList ::= symWithAttrsList symWithAttrs
            //
            case 113: {
               //#line 210 "LPGParser.g"
                (my.GetRhsSym(1).(*symWithAttrsList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 114:  symWithAttrs ::= EMPTY_KEY
            //
            case 114: {
               //#line 212 "LPGParser.g"
                my.SetResult(
                    //#line 212 LPGParser.g
                    NewsymWithAttrs0(my.GetRhsIToken(1)),
                //#line 212 LPGParser.g
                )
                break
            }
            //
            // Rule 115:  symWithAttrs ::= SYMBOL optAttrList
            //
            case 115: {
               //#line 213 "LPGParser.g"
                my.SetResult(
                    //#line 213 LPGParser.g
                    NewsymWithAttrs1(my.GetLeftIToken(), my.GetRightIToken(),
                                     //#line 213 LPGParser.g
                                     NewASTNodeToken(my.GetRhsIToken(1)),
                                     //#line 213 LPGParser.g
                                     AnyCastTosymAttrs(my.GetRhsSym(2))),
                //#line 213 LPGParser.g
                )
                break
            }
            //
            // Rule 116:  optAttrList ::= $Empty
            //
            case 116: {
               //#line 216 "LPGParser.g"
                my.SetResult(
                    //#line 216 LPGParser.g
                    NewsymAttrs(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 216 LPGParser.g
                                nil),
                //#line 216 LPGParser.g
                )
                break
            }
            //
            // Rule 117:  optAttrList ::= MACRO_NAME
            //
            case 117: {
               //#line 217 "LPGParser.g"
                my.SetResult(
                    //#line 217 LPGParser.g
                    NewsymAttrs(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 217 LPGParser.g
                                NewASTNodeToken(my.GetRhsIToken(1))),
                //#line 217 LPGParser.g
                )
                break
            }
            //
            // Rule 118:  opt_action_segment ::= $Empty
            //
            case 118: {
               //#line 219 "LPGParser.g"
                my.SetResult(nil);
                break
            }
            //
            // Rule 119:  opt_action_segment ::= action_segment
            //
            case 119:
                break
            //
            // Rule 120:  action_segment ::= BLOCK
            //
            case 120: {
               //#line 221 "LPGParser.g"
                my.SetResult(
                    //#line 221 LPGParser.g
                    Newaction_segment(my, my.GetRhsIToken(1)),
                //#line 221 LPGParser.g
                )
                break
            }
            //
            // Rule 121:  start_segment ::= start_symbol
            //
            case 121: {
               //#line 226 "LPGParser.g"
                my.SetResult(
                    //#line 226 LPGParser.g
                    Newstart_symbolListFromElement(my.GetRhsSym(1).(Istart_symbol), true /* left recursive */),
                //#line 226 LPGParser.g
                )
                break
            }
            //
            // Rule 122:  start_segment ::= start_segment start_symbol
            //
            case 122: {
               //#line 226 "LPGParser.g"
                (my.GetRhsSym(1).(*start_symbolList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 123:  start_symbol ::= SYMBOL
            //
            case 123: {
               //#line 227 "LPGParser.g"
                my.SetResult(
                    //#line 227 LPGParser.g
                    Newstart_symbol0(my.GetRhsIToken(1)),
                //#line 227 LPGParser.g
                )
                break
            }
            //
            // Rule 124:  start_symbol ::= MACRO_NAME
            //
            case 124: {
               //#line 228 "LPGParser.g"
                my.SetResult(
                    //#line 228 LPGParser.g
                    Newstart_symbol1(my.GetRhsIToken(1)),
                //#line 228 LPGParser.g
                )
                break
            }
            //
            // Rule 125:  terminals_segment ::= terminal
            //
            case 125: {
               //#line 231 "LPGParser.g"
                my.SetResult(
                    //#line 231 LPGParser.g
                    NewterminalListFromElement(my.GetRhsSym(1).(*terminal), true /* left recursive */),
                //#line 231 LPGParser.g
                )
                break
            }
            //
            // Rule 126:  terminals_segment ::= terminals_segment terminal
            //
            case 126: {
               //#line 231 "LPGParser.g"
                (my.GetRhsSym(1).(*terminalList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 127:  terminal ::= terminal_symbol optTerminalAlias
            //
            case 127: {
               //#line 234 "LPGParser.g"
                my.SetResult(
                    //#line 234 LPGParser.g
                    Newterminal(my.GetLeftIToken(), my.GetRightIToken(),
                                //#line 234 LPGParser.g
                                my.GetRhsSym(1).(Iterminal_symbol),
                                //#line 234 LPGParser.g
                                AnyCastTooptTerminalAlias(my.GetRhsSym(2))),
                //#line 234 LPGParser.g
                )
                break
            }
            //
            // Rule 128:  optTerminalAlias ::= $Empty
            //
            case 128: {
               //#line 236 "LPGParser.g"
                my.SetResult(nil);
                break
            }
            //
            // Rule 129:  optTerminalAlias ::= produces name
            //
            case 129: {
               //#line 236 "LPGParser.g"
                my.SetResult(
                    //#line 236 LPGParser.g
                    NewoptTerminalAlias(my.GetLeftIToken(), my.GetRightIToken(),
                                        //#line 236 LPGParser.g
                                        my.GetRhsSym(1).(Iproduces),
                                        //#line 236 LPGParser.g
                                        my.GetRhsSym(2).(Iname)),
                //#line 236 LPGParser.g
                )
                break
            }
            //
            // Rule 130:  terminal_symbol ::= SYMBOL
            //
            case 130: {
               //#line 238 "LPGParser.g"
                my.SetResult(
                    //#line 238 LPGParser.g
                    Newterminal_symbol0(my.GetRhsIToken(1)),
                //#line 238 LPGParser.g
                )
                break
            }
            //
            // Rule 131:  terminal_symbol ::= MACRO_NAME
            //
            case 131: {
               //#line 240 "LPGParser.g"
                my.SetResult(
                    //#line 240 LPGParser.g
                    Newterminal_symbol1(my.GetRhsIToken(1)),
                //#line 240 LPGParser.g
                )
                break
            }
            //
            // Rule 132:  trailers_segment ::= action_segment_list
            //
            case 132:
                break
            //
            // Rule 133:  types_segment ::= type_declarations
            //
            case 133: {
               //#line 246 "LPGParser.g"
                my.SetResult(
                    //#line 246 LPGParser.g
                    Newtype_declarationsListFromElement(my.GetRhsSym(1).(*type_declarations), true /* left recursive */),
                //#line 246 LPGParser.g
                )
                break
            }
            //
            // Rule 134:  types_segment ::= types_segment type_declarations
            //
            case 134: {
               //#line 246 "LPGParser.g"
                (my.GetRhsSym(1).(*type_declarationsList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 135:  type_declarations ::= SYMBOL produces barSymbolList opt_action_segment
            //
            case 135: {
               //#line 248 "LPGParser.g"
                my.SetResult(
                    //#line 248 LPGParser.g
                    Newtype_declarations(my.GetLeftIToken(), my.GetRightIToken(),
                                         //#line 248 LPGParser.g
                                         NewASTNodeToken(my.GetRhsIToken(1)),
                                         //#line 248 LPGParser.g
                                         my.GetRhsSym(2).(Iproduces),
                                         //#line 248 LPGParser.g
                                         my.GetRhsSym(3).(*SYMBOLList),
                                         //#line 248 LPGParser.g
                                         AnyCastToaction_segment(my.GetRhsSym(4))),
                //#line 248 LPGParser.g
                )
                break
            }
            //
            // Rule 136:  barSymbolList ::= SYMBOL
            //
            case 136: {
               //#line 249 "LPGParser.g"
                my.SetResult(
                    //#line 249 LPGParser.g
                    NewSYMBOLListFromElement(NewASTNodeToken(my.GetRhsIToken(1)), true /* left recursive */),
                //#line 249 LPGParser.g
                )
                break
            }
            //
            // Rule 137:  barSymbolList ::= barSymbolList |$ SYMBOL
            //
            case 137: {
               //#line 249 "LPGParser.g"
                (my.GetRhsSym(1).(*SYMBOLList)).AddElement(NewASTNodeToken(my.GetRhsIToken(3)))
                break
            }
            //
            // Rule 138:  predecessor_segment ::= $Empty
            //
            case 138: {
               //#line 252 "LPGParser.g"
                my.SetResult(
                    //#line 252 LPGParser.g
                    Newsymbol_pairList(my.GetLeftIToken(), my.GetRightIToken(), true /* left recursive */),
                //#line 252 LPGParser.g
                )
                break
            }
            //
            // Rule 139:  predecessor_segment ::= predecessor_segment symbol_pair
            //
            case 139: {
               //#line 252 "LPGParser.g"
                (my.GetRhsSym(1).(*symbol_pairList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
            //
            // Rule 140:  symbol_pair ::= SYMBOL SYMBOL
            //
            case 140: {
               //#line 254 "LPGParser.g"
                my.SetResult(
                    //#line 254 LPGParser.g
                    Newsymbol_pair(my.GetLeftIToken(), my.GetRightIToken(),
                                   //#line 254 LPGParser.g
                                   NewASTNodeToken(my.GetRhsIToken(1)),
                                   //#line 254 LPGParser.g
                                   NewASTNodeToken(my.GetRhsIToken(2))),
                //#line 254 LPGParser.g
                )
                break
            }
            //
            // Rule 141:  recover_segment ::= $Empty
            //
            case 141: {
               //#line 257 "LPGParser.g"
                my.SetResult(
                    //#line 257 LPGParser.g
                    NewSYMBOLList(my.GetLeftIToken(), my.GetRightIToken(), true /* left recursive */),
                //#line 257 LPGParser.g
                )
                break
            }
            //
            // Rule 142:  recover_segment ::= recover_segment recover_symbol
            //
            case 142: {
               //#line 257 "LPGParser.g"
                my.SetResult(my.GetRhsSym(1).(*SYMBOLList))
                break
            }
            //
            // Rule 143:  recover_symbol ::= SYMBOL
            //
            case 143: {
               //#line 259 "LPGParser.g"
                my.SetResult(
                    //#line 259 LPGParser.g
                    Newrecover_symbol(my.GetRhsIToken(1)),
                //#line 259 LPGParser.g
                )
                break
            }
            //
            // Rule 144:  END_KEY_OPT ::= $Empty
            //
            case 144: {
               //#line 262 "LPGParser.g"
                my.SetResult(nil);
                break
            }
            //
            // Rule 145:  END_KEY_OPT ::= END_KEY
            //
            case 145: {
               //#line 263 "LPGParser.g"
                my.SetResult(
                    //#line 263 LPGParser.g
                    NewEND_KEY_OPT(my.GetRhsIToken(1)),
                //#line 263 LPGParser.g
                )
                break
            }
            //
            // Rule 146:  action_segment_list ::= $Empty
            //
            case 146: {
               //#line 265 "LPGParser.g"
                my.SetResult(
                    //#line 265 LPGParser.g
                    Newaction_segmentList(my.GetLeftIToken(), my.GetRightIToken(), true /* left recursive */),
                //#line 265 LPGParser.g
                )
                break
            }
            //
            // Rule 147:  action_segment_list ::= action_segment_list action_segment
            //
            case 147: {
               //#line 266 "LPGParser.g"
                (my.GetRhsSym(1).(*action_segmentList)).AddElement(my.GetRhsSym(2).(IAst))
                break
            }
    //#line 333 "dtParserTemplateF.gi

    
            default:
                break
        }
        return
    }
type ASTNode struct{
     leftIToken  IToken 
     rightIToken IToken 
     parent IAst
}
func NewASTNode2(leftIToken  IToken , rightIToken  IToken ) *ASTNode{
        my := new(ASTNode)
        my.leftIToken = leftIToken
        my.rightIToken = rightIToken
        return my
}

func NewASTNode(token  IToken) *ASTNode{
        my := new(ASTNode)
        my.leftIToken = token
        my.rightIToken = token
        return my
}

func (my *ASTNode)     GetNextAst() IAst  { return nil }
func (my *ASTNode)      SetParent(parent IAst )  { my.parent = parent }
func (my *ASTNode)      GetParent() IAst { return my.parent }

func (my *ASTNode)     GetLeftIToken()  IToken { return my.leftIToken }
func (my *ASTNode)     GetRightIToken()  IToken { return my.rightIToken }
func (my *ASTNode)     GetPrecedingAdjuncts()  []IToken { return my.leftIToken.GetPrecedingAdjuncts() }
func (my *ASTNode)     GetFollowingAdjuncts()  []IToken { return my.rightIToken.GetFollowingAdjuncts() }

func (my *ASTNode)      ToString()string{
        return my.leftIToken.GetILexStream().ToString(my.leftIToken.GetStartOffset(), my.rightIToken.GetEndOffset())
}

func (my *ASTNode)     Initialize()  {}

    /**
     * A list of all children of my node, excluding the null ones.
     */
func (my *ASTNode)      GetChildren() *ArrayList{
        var list = my.GetAllChildren() 
        var k = -1
        var i = 0
        for ; i < list.Size(); i++{
            var element = list.Get(i)
            if element != nil{
                k += i
                if k != i{
                    list.Set(k, element)
                }
            }
        }
        i = list.Size() - 1
        for ; i > k; i--{ // remove extraneous elements
            list.RemoveAt(i)
        }
        return list
}

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *ASTNode)    GetAllChildren() *ArrayList{return nil}

 func (my *ASTNode) Accept(v IAstVisitor){}


func AnyCastToASTNode(i interface{}) *ASTNode {
	if nil == i{
		return nil
	}else{
		return i.(*ASTNode)
	}
}
type AbstractASTNodeList struct{
     *ASTNode
     leftRecursive bool 
     list *ArrayList 
}
func NewAbstractASTNodeList(leftToken  IToken, rightToken  IToken, leftRecursive bool)*AbstractASTNodeList{
      my := new(AbstractASTNodeList)
      my.ASTNode = NewASTNode2(leftToken, rightToken)
      my.list = NewArrayList()
      my.leftRecursive = leftRecursive
      return my
}

func (my *AbstractASTNodeList)      Size() int { return my.list.Size(); }
func (my *AbstractASTNodeList)      GetList() *ArrayList{ return my.list }
func (my *AbstractASTNodeList)      GetElementAt(i int ) IAst{
    var k int
    if my.leftRecursive {
       k =i
    }else{
       k =my.list.Size() - 1 - i
    }
    return my.list.Get(k).(IAst) 
    }
func (my *AbstractASTNodeList)      GetArrayList() *ArrayList{
        if ! my.leftRecursive{ // reverse the list 
           var i = 0
           var n = my.list.Size() - 1
           for ; i < n;  n--{
                var ith = my.list.Get(i)
                var nth = my.list.Get(n)
                my.list.Set(i, nth)
                my.list.Set(n, ith)
               i++
           }
           my.leftRecursive = true
        }
        return my.list
    }
    /**
     * @deprecated replaced by {@link #AddElement()}
     *
     */
func (my *AbstractASTNodeList)      Add(element IAst) bool {
        my.AddElement(element)
        return true
}

func (my *AbstractASTNodeList)      AddElement(element IAst){
        my.list.Add(element)
        if my.leftRecursive{
             my.rightIToken = element.GetRightIToken()
        }else{
          my.leftIToken = element.GetLeftIToken()
        }
}

    /**
     * Make a copy of the list and return it. Note that we obtain the local list by
     * invoking GetArrayList so as to make sure that the list we return is in proper order.
     */
func (my *AbstractASTNodeList)      GetAllChildren() *ArrayList{
        return my.GetArrayList().Clone()
}



func AnyCastToAbstractASTNodeList(i interface{}) *AbstractASTNodeList {
	if nil == i{
		return nil
	}else{
		return i.(*AbstractASTNodeList)
	}
}
type ASTNodeToken struct{
    *ASTNode
 }
func NewASTNodeToken(token  IToken)*ASTNodeToken{
      my := new(ASTNodeToken)
      my.ASTNode = NewASTNode(token)
      return my
}

func (my *ASTNodeToken)      GetIToken()  IToken{ return my.leftIToken }
func (my *ASTNodeToken)      ToString()  string  { return my.leftIToken.ToString() }

    /**
     * A token class has no children. So, we return the empty list.
     */
func (my *ASTNodeToken)        GetAllChildren()  *ArrayList { return nil }


func (my *ASTNodeToken)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *ASTNodeToken)       Enter(v Visitor){
        v.VisitASTNodeToken(my)
        v.EndVisitASTNodeToken(my)
    }


func AnyCastToASTNodeToken(i interface{}) *ASTNodeToken {
	if nil == i{
		return nil
	}else{
		return i.(*ASTNodeToken)
	}
}
type IRootForLPGParser interface{
     GetLeftIToken()  IToken
     GetRightIToken()  IToken
 Accept(v IAstVisitor)
}

func CastToAnyForLPGParser(i interface{}) interface{}{return i}

func AnyCastToIRootForLPGParser(i interface{}) IRootForLPGParser {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IRootForLPGParser)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by:
 *<b>
 *<ul>
 *<li>alias_lhs_macro_name
 *<li>macro_segment
 *<li>optMacroName
 *<li>include_segment
 *<li>RuleName
 *<li>symAttrs
 *<li>action_segment
 *<li>recover_symbol
 *<li>END_KEY_OPT
 *<li>alias_rhs0
 *<li>alias_rhs1
 *<li>alias_rhs2
 *<li>alias_rhs3
 *<li>alias_rhs4
 *<li>alias_rhs5
 *<li>alias_rhs6
 *<li>macro_name_symbol0
 *<li>macro_name_symbol1
 *<li>name0
 *<li>name1
 *<li>name2
 *<li>name3
 *<li>name4
 *<li>name5
 *<li>produces0
 *<li>produces1
 *<li>produces2
 *<li>produces3
 *<li>symWithAttrs0
 *<li>symWithAttrs1
 *<li>start_symbol0
 *<li>start_symbol1
 *<li>terminal_symbol0
 *<li>terminal_symbol1
 *</ul>
 *</b>
 */
type IASTNodeToken interface{
IRootForLPGParser
}

func AnyCastToIASTNodeToken(i interface{}) IASTNodeToken {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IASTNodeToken)
	  }
}
/**
 * is implemented by <b>LPG</b>
 */
type ILPG interface{
IRootForLPGParser
}

func AnyCastToILPG(i interface{}) ILPG {
	  if nil == i{
		 return nil
	  }else{
		 return i.(ILPG)
	  }
}
/**
 * is implemented by <b>option_specList</b>
 */
type Ioptions_segment interface{
IRootForLPGParser
}

func AnyCastToIoptions_segment(i interface{}) Ioptions_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ioptions_segment)
	  }
}
/**
 * is implemented by <b>LPG_itemList</b>
 */
type ILPG_INPUT interface{
IRootForLPGParser
}

func AnyCastToILPG_INPUT(i interface{}) ILPG_INPUT {
	  if nil == i{
		 return nil
	  }else{
		 return i.(ILPG_INPUT)
	  }
}
/**
 * is implemented by:
 *<b>
 *<ul>
 *<li>AliasSeg
 *<li>AstSeg
 *<li>DefineSeg
 *<li>EofSeg
 *<li>EolSeg
 *<li>ErrorSeg
 *<li>ExportSeg
 *<li>GlobalsSeg
 *<li>HeadersSeg
 *<li>IdentifierSeg
 *<li>ImportSeg
 *<li>IncludeSeg
 *<li>KeywordsSeg
 *<li>NamesSeg
 *<li>NoticeSeg
 *<li>RulesSeg
 *<li>SoftKeywordsSeg
 *<li>StartSeg
 *<li>TerminalsSeg
 *<li>TrailersSeg
 *<li>TypesSeg
 *<li>RecoverSeg
 *<li>PredecessorSeg
 *</ul>
 *</b>
 */
type ILPG_item interface{
IRootForLPGParser
}

func AnyCastToILPG_item(i interface{}) ILPG_item {
	  if nil == i{
		 return nil
	  }else{
		 return i.(ILPG_item)
	  }
}
/**
 * is implemented by <b>aliasSpecList</b>
 */
type Ialias_segment interface{
IRootForLPGParser
}

func AnyCastToIalias_segment(i interface{}) Ialias_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ialias_segment)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by <b>END_KEY_OPT</b>
 */
type IEND_KEY_OPT interface{
IASTNodeToken
}

func AnyCastToIEND_KEY_OPT(i interface{}) IEND_KEY_OPT {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IEND_KEY_OPT)
	  }
}
/**
 * is implemented by <b>action_segmentList</b>
 */
type Iast_segment interface{
IRootForLPGParser
}

func AnyCastToIast_segment(i interface{}) Iast_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iast_segment)
	  }
}
/**
 * is implemented by <b>defineSpecList</b>
 */
type Idefine_segment interface{
IRootForLPGParser
}

func AnyCastToIdefine_segment(i interface{}) Idefine_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Idefine_segment)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by:
 *<b>
 *<ul>
 *<li>terminal_symbol0
 *<li>terminal_symbol1
 *</ul>
 *</b>
 */
type Ieof_segment interface{
IRootForLPGParser
}

func AnyCastToIeof_segment(i interface{}) Ieof_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ieof_segment)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by:
 *<b>
 *<ul>
 *<li>terminal_symbol0
 *<li>terminal_symbol1
 *</ul>
 *</b>
 */
type Ieol_segment interface{
IRootForLPGParser
}

func AnyCastToIeol_segment(i interface{}) Ieol_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ieol_segment)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by:
 *<b>
 *<ul>
 *<li>terminal_symbol0
 *<li>terminal_symbol1
 *</ul>
 *</b>
 */
type Ierror_segment interface{
IRootForLPGParser
}

func AnyCastToIerror_segment(i interface{}) Ierror_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ierror_segment)
	  }
}
/**
 * is implemented by <b>terminal_symbolList</b>
 */
type Iexport_segment interface{
IRootForLPGParser
}

func AnyCastToIexport_segment(i interface{}) Iexport_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iexport_segment)
	  }
}
/**
 * is implemented by <b>action_segmentList</b>
 */
type Iglobals_segment interface{
IRootForLPGParser
}

func AnyCastToIglobals_segment(i interface{}) Iglobals_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iglobals_segment)
	  }
}
/**
 * is implemented by <b>action_segmentList</b>
 */
type Iheaders_segment interface{
IRootForLPGParser
}

func AnyCastToIheaders_segment(i interface{}) Iheaders_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iheaders_segment)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by:
 *<b>
 *<ul>
 *<li>terminal_symbol0
 *<li>terminal_symbol1
 *</ul>
 *</b>
 */
type Iidentifier_segment interface{
IRootForLPGParser
}

func AnyCastToIidentifier_segment(i interface{}) Iidentifier_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iidentifier_segment)
	  }
}
/**
 * is implemented by <b>import_segment</b>
 */
type Iimport_segment interface{
IRootForLPGParser
}

func AnyCastToIimport_segment(i interface{}) Iimport_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iimport_segment)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by <b>include_segment</b>
 */
type Iinclude_segment interface{
IASTNodeToken
}

func AnyCastToIinclude_segment(i interface{}) Iinclude_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iinclude_segment)
	  }
}
/**
 * is implemented by <b>keywordSpecList</b>
 */
type Ikeywords_segment interface{
IRootForLPGParser
}

func AnyCastToIkeywords_segment(i interface{}) Ikeywords_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ikeywords_segment)
	  }
}
/**
 * is implemented by <b>nameSpecList</b>
 */
type Inames_segment interface{
IRootForLPGParser
}

func AnyCastToInames_segment(i interface{}) Inames_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Inames_segment)
	  }
}
/**
 * is implemented by <b>action_segmentList</b>
 */
type Inotice_segment interface{
IRootForLPGParser
}

func AnyCastToInotice_segment(i interface{}) Inotice_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Inotice_segment)
	  }
}
/**
 * is implemented by <b>rules_segment</b>
 */
type Irules_segment interface{
IRootForLPGParser
}

func AnyCastToIrules_segment(i interface{}) Irules_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Irules_segment)
	  }
}
/**
 * is implemented by <b>start_symbolList</b>
 */
type Istart_segment interface{
IRootForLPGParser
}

func AnyCastToIstart_segment(i interface{}) Istart_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Istart_segment)
	  }
}
/**
 * is implemented by <b>terminalList</b>
 */
type Iterminals_segment interface{
IRootForLPGParser
}

func AnyCastToIterminals_segment(i interface{}) Iterminals_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iterminals_segment)
	  }
}
/**
 * is implemented by <b>action_segmentList</b>
 */
type Itrailers_segment interface{
IRootForLPGParser
}

func AnyCastToItrailers_segment(i interface{}) Itrailers_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Itrailers_segment)
	  }
}
/**
 * is implemented by <b>type_declarationsList</b>
 */
type Itypes_segment interface{
IRootForLPGParser
}

func AnyCastToItypes_segment(i interface{}) Itypes_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Itypes_segment)
	  }
}
/**
 * is implemented by <b>SYMBOLList</b>
 */
type Irecover_segment interface{
IRootForLPGParser
}

func AnyCastToIrecover_segment(i interface{}) Irecover_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Irecover_segment)
	  }
}
/**
 * is implemented by <b>symbol_pairList</b>
 */
type Ipredecessor_segment interface{
IRootForLPGParser
}

func AnyCastToIpredecessor_segment(i interface{}) Ipredecessor_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ipredecessor_segment)
	  }
}
/**
 * is implemented by <b>option_spec</b>
 */
type Ioption_spec interface{
IRootForLPGParser
}

func AnyCastToIoption_spec(i interface{}) Ioption_spec {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ioption_spec)
	  }
}
/**
 * is implemented by <b>optionList</b>
 */
type Ioption_list interface{
IRootForLPGParser
}

func AnyCastToIoption_list(i interface{}) Ioption_list {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ioption_list)
	  }
}
/**
 * is implemented by <b>option</b>
 */
type Ioption interface{
IRootForLPGParser
}

func AnyCastToIoption(i interface{}) Ioption {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ioption)
	  }
}
/**
 * is implemented by:
 *<b>
 *<ul>
 *<li>option_value0
 *<li>option_value1
 *</ul>
 *</b>
 */
type Ioption_value interface{
IRootForLPGParser
}

func AnyCastToIoption_value(i interface{}) Ioption_value {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ioption_value)
	  }
}
/**
 * is implemented by <b>SYMBOLList</b>
 */
type Isymbol_list interface{
IRootForLPGParser
}

func AnyCastToIsymbol_list(i interface{}) Isymbol_list {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Isymbol_list)
	  }
}
/**
 * is implemented by:
 *<b>
 *<ul>
 *<li>aliasSpec0
 *<li>aliasSpec1
 *<li>aliasSpec2
 *<li>aliasSpec3
 *<li>aliasSpec4
 *<li>aliasSpec5
 *</ul>
 *</b>
 */
type IaliasSpec interface{
IRootForLPGParser
}

func AnyCastToIaliasSpec(i interface{}) IaliasSpec {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IaliasSpec)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by:
 *<b>
 *<ul>
 *<li>produces0
 *<li>produces1
 *<li>produces2
 *<li>produces3
 *</ul>
 *</b>
 */
type Iproduces interface{
IASTNodeToken
}

func AnyCastToIproduces(i interface{}) Iproduces {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iproduces)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by:
 *<b>
 *<ul>
 *<li>alias_rhs0
 *<li>alias_rhs1
 *<li>alias_rhs2
 *<li>alias_rhs3
 *<li>alias_rhs4
 *<li>alias_rhs5
 *<li>alias_rhs6
 *</ul>
 *</b>
 */
type Ialias_rhs interface{
IASTNodeToken
}

func AnyCastToIalias_rhs(i interface{}) Ialias_rhs {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ialias_rhs)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by <b>alias_lhs_macro_name</b>
 */
type Ialias_lhs_macro_name interface{
IASTNodeToken
}

func AnyCastToIalias_lhs_macro_name(i interface{}) Ialias_lhs_macro_name {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Ialias_lhs_macro_name)
	  }
}
/**
 * is implemented by <b>action_segmentList</b>
 */
type Iaction_segment_list interface{
Iast_segment
Iheaders_segment
Itrailers_segment
}

func AnyCastToIaction_segment_list(i interface{}) Iaction_segment_list {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iaction_segment_list)
	  }
}
/**
 * is implemented by <b>defineSpec</b>
 */
type IdefineSpec interface{
IRootForLPGParser
}

func AnyCastToIdefineSpec(i interface{}) IdefineSpec {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IdefineSpec)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by:
 *<b>
 *<ul>
 *<li>macro_name_symbol0
 *<li>macro_name_symbol1
 *</ul>
 *</b>
 */
type Imacro_name_symbol interface{
IASTNodeToken
}

func AnyCastToImacro_name_symbol(i interface{}) Imacro_name_symbol {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Imacro_name_symbol)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by <b>macro_segment</b>
 */
type Imacro_segment interface{
IASTNodeToken
}

func AnyCastToImacro_segment(i interface{}) Imacro_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Imacro_segment)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by:
 *<b>
 *<ul>
 *<li>terminal_symbol0
 *<li>terminal_symbol1
 *</ul>
 *</b>
 */
type Iterminal_symbol interface{
Ieol_segment
Ieof_segment
Ierror_segment
Iidentifier_segment
IkeywordSpec
IASTNodeToken
}

func AnyCastToIterminal_symbol(i interface{}) Iterminal_symbol {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iterminal_symbol)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by <b>action_segment</b>
 */
type Iaction_segment interface{
Iopt_action_segment
IASTNodeToken
}

func AnyCastToIaction_segment(i interface{}) Iaction_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iaction_segment)
	  }
}
/**
 * is implemented by <b>drop_commandList</b>
 */
type Idrop_command_list interface{
IRootForLPGParser
}

func AnyCastToIdrop_command_list(i interface{}) Idrop_command_list {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Idrop_command_list)
	  }
}
/**
 * is implemented by:
 *<b>
 *<ul>
 *<li>drop_command0
 *<li>drop_command1
 *</ul>
 *</b>
 */
type Idrop_command interface{
IRootForLPGParser
}

func AnyCastToIdrop_command(i interface{}) Idrop_command {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Idrop_command)
	  }
}
/**
 * is implemented by <b>SYMBOLList</b>
 */
type Idrop_symbols interface{
IRootForLPGParser
}

func AnyCastToIdrop_symbols(i interface{}) Idrop_symbols {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Idrop_symbols)
	  }
}
/**
 * is implemented by <b>drop_ruleList</b>
 */
type Idrop_rules interface{
IRootForLPGParser
}

func AnyCastToIdrop_rules(i interface{}) Idrop_rules {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Idrop_rules)
	  }
}
/**
 * is implemented by <b>drop_rule</b>
 */
type Idrop_rule interface{
IRootForLPGParser
}

func AnyCastToIdrop_rule(i interface{}) Idrop_rule {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Idrop_rule)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by <b>optMacroName</b>
 */
type IoptMacroName interface{
IASTNodeToken
}

func AnyCastToIoptMacroName(i interface{}) IoptMacroName {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IoptMacroName)
	  }
}
/**
 * is implemented by <b>ruleList</b>
 */
type IruleList interface{
IRootForLPGParser
}

func AnyCastToIruleList(i interface{}) IruleList {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IruleList)
	  }
}
/**
 * is implemented by:
 *<b>
 *<ul>
 *<li>keywordSpec
 *<li>terminal_symbol0
 *<li>terminal_symbol1
 *</ul>
 *</b>
 */
type IkeywordSpec interface{
IRootForLPGParser
}

func AnyCastToIkeywordSpec(i interface{}) IkeywordSpec {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IkeywordSpec)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by:
 *<b>
 *<ul>
 *<li>name0
 *<li>name1
 *<li>name2
 *<li>name3
 *<li>name4
 *<li>name5
 *</ul>
 *</b>
 */
type Iname interface{
IASTNodeToken
}

func AnyCastToIname(i interface{}) Iname {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iname)
	  }
}
/**
 * is implemented by <b>nameSpec</b>
 */
type InameSpec interface{
IRootForLPGParser
}

func AnyCastToInameSpec(i interface{}) InameSpec {
	  if nil == i{
		 return nil
	  }else{
		 return i.(InameSpec)
	  }
}
/**
 * is implemented by <b>nonTermList</b>
 */
type InonTermList interface{
IRootForLPGParser
}

func AnyCastToInonTermList(i interface{}) InonTermList {
	  if nil == i{
		 return nil
	  }else{
		 return i.(InonTermList)
	  }
}
/**
 * is implemented by <b>nonTerm</b>
 */
type InonTerm interface{
IRootForLPGParser
}

func AnyCastToInonTerm(i interface{}) InonTerm {
	  if nil == i{
		 return nil
	  }else{
		 return i.(InonTerm)
	  }
}
/**
 * is implemented by <b>RuleName</b>
 */
type IruleNameWithAttributes interface{
IASTNodeToken
}

func AnyCastToIruleNameWithAttributes(i interface{}) IruleNameWithAttributes {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IruleNameWithAttributes)
	  }
}
/**
 * is implemented by <b>rule</b>
 */
type Irule interface{
IRootForLPGParser
}

func AnyCastToIrule(i interface{}) Irule {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Irule)
	  }
}
/**
 * is implemented by <b>symWithAttrsList</b>
 */
type IsymWithAttrsList interface{
IRootForLPGParser
}

func AnyCastToIsymWithAttrsList(i interface{}) IsymWithAttrsList {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IsymWithAttrsList)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by <b>action_segment</b>
 */
type Iopt_action_segment interface{
IRootForLPGParser
}

func AnyCastToIopt_action_segment(i interface{}) Iopt_action_segment {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iopt_action_segment)
	  }
}
/**
 * is implemented by:
 *<b>
 *<ul>
 *<li>symWithAttrs0
 *<li>symWithAttrs1
 *</ul>
 *</b>
 */
type IsymWithAttrs interface{
IASTNodeToken
}

func AnyCastToIsymWithAttrs(i interface{}) IsymWithAttrs {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IsymWithAttrs)
	  }
}
/**
 * is implemented by <b>symAttrs</b>
 */
type IoptAttrList interface{
IASTNodeToken
}

func AnyCastToIoptAttrList(i interface{}) IoptAttrList {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IoptAttrList)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by:
 *<b>
 *<ul>
 *<li>start_symbol0
 *<li>start_symbol1
 *</ul>
 *</b>
 */
type Istart_symbol interface{
IASTNodeToken
}

func AnyCastToIstart_symbol(i interface{}) Istart_symbol {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Istart_symbol)
	  }
}
/**
 * is implemented by <b>terminal</b>
 */
type Iterminal interface{
IRootForLPGParser
}

func AnyCastToIterminal(i interface{}) Iterminal {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Iterminal)
	  }
}
/**
 * is implemented by <b>optTerminalAlias</b>
 */
type IoptTerminalAlias interface{
IRootForLPGParser
}

func AnyCastToIoptTerminalAlias(i interface{}) IoptTerminalAlias {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IoptTerminalAlias)
	  }
}
/**
 * is implemented by <b>type_declarations</b>
 */
type Itype_declarations interface{
IRootForLPGParser
}

func AnyCastToItype_declarations(i interface{}) Itype_declarations {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Itype_declarations)
	  }
}
/**
 * is implemented by <b>SYMBOLList</b>
 */
type IbarSymbolList interface{
IRootForLPGParser
}

func AnyCastToIbarSymbolList(i interface{}) IbarSymbolList {
	  if nil == i{
		 return nil
	  }else{
		 return i.(IbarSymbolList)
	  }
}
/**
 * is implemented by <b>symbol_pair</b>
 */
type Isymbol_pair interface{
IRootForLPGParser
}

func AnyCastToIsymbol_pair(i interface{}) Isymbol_pair {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Isymbol_pair)
	  }
}
/**
 * is always implemented by <b>ASTNodeToken</b>. It is also implemented by <b>recover_symbol</b>
 */
type Irecover_symbol interface{
IASTNodeToken
}

func AnyCastToIrecover_symbol(i interface{}) Irecover_symbol {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Irecover_symbol)
	  }
}
/**
 *<b>
*<li>Rule 1:  LPG ::= options_segment LPG_INPUT
 *</b>
 */
type LPG struct{
    *ASTNode
     environment *LPGParser
      _options_segment *option_specList
      _LPG_INPUT *LPG_itemList
}
func (my *LPG)      Getoptions_segment() *option_specList{ return my._options_segment}
func (my *LPG)      Setoptions_segment( _options_segment *option_specList)  { my._options_segment = _options_segment }
func (my *LPG)      GetLPG_INPUT() *LPG_itemList{ return my._LPG_INPUT}
func (my *LPG)      SetLPG_INPUT( _LPG_INPUT *LPG_itemList)  { my._LPG_INPUT = _LPG_INPUT }

func (my *LPG)     GetEnvironment() *LPGParser{ return my.environment }

func NewLPG(environment *LPGParser,leftIToken IToken, rightIToken IToken ,
              _options_segment *option_specList,
              _LPG_INPUT *LPG_itemList)*LPG{
      my := new(LPG)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my.environment = environment
        my._options_segment = _options_segment;
        if nil != _options_segment{
        var trait_ interface{} = _options_segment
         trait_.(IAst).SetParent(my)
}
        my._LPG_INPUT = _LPG_INPUT;
        if nil != _LPG_INPUT{
        var trait_ interface{} = _LPG_INPUT
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *LPG)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._options_segment{  list.Add(my._options_segment) }
        if nil != my._LPG_INPUT{  list.Add(my._LPG_INPUT) }
        return list
    }

func (my *LPG)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *LPG)       Enter(v Visitor){
        var checkChildren = v.VisitLPG(my)
        if checkChildren{
            if nil != my._options_segment{my._options_segment.Accept(v)}
            if nil != my._LPG_INPUT{my._LPG_INPUT.Accept(v)}
        }
        v.EndVisitLPG(my)
    }


func AnyCastToLPG(i interface{}) *LPG {
	if nil == i{
		return nil
	}else{
		return i.(*LPG)
	}
}
/**
 *<b>
*<li>Rule 2:  LPG_INPUT ::= $Empty
*<li>Rule 3:  LPG_INPUT ::= LPG_INPUT LPG_item
 *</b>
 */
type LPG_itemList struct{
    *AbstractASTNodeList
}
func (my *LPG_itemList)      GetLPG_itemAt(i int) ILPG_item{
     var r,_=my.GetElementAt(i).(ILPG_item)
     return r
     }

func NewLPG_itemList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*LPG_itemList{
      my := new(LPG_itemList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  NewLPG_itemListFromElement(element ILPG_item,leftRecursive bool)*LPG_itemList{
        var obj = NewLPG_itemList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *LPG_itemList)      AddElement(_LPG_item IAst){ 
      my.AbstractASTNodeList.AddElement(_LPG_item)
        _LPG_item.SetParent(my)
    }


func (my *LPG_itemList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *LPG_itemList)     Enter(v  Visitor){
        var checkChildren = v.VisitLPG_itemList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.GetLPG_itemAt(i)
                if nil !=element{                    element.Accept(v)
                }
            }
        }
        v.EndVisitLPG_itemList(my)
    }


func AnyCastToLPG_itemList(i interface{}) *LPG_itemList {
	if nil == i{
		return nil
	}else{
		return i.(*LPG_itemList)
	}
}
/**
 *<b>
*<li>Rule 4:  LPG_item ::= ALIAS_KEY$ alias_segment END_KEY_OPT$
 *</b>
 */
type AliasSeg struct{
    *ASTNode
      _alias_segment *aliasSpecList
}
func (my *AliasSeg)      Getalias_segment() *aliasSpecList{ return my._alias_segment}
func (my *AliasSeg)      Setalias_segment( _alias_segment *aliasSpecList)  { my._alias_segment = _alias_segment }

func NewAliasSeg(leftIToken IToken, rightIToken IToken ,
              _alias_segment *aliasSpecList)*AliasSeg{
      my := new(AliasSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._alias_segment = _alias_segment;
        if nil != _alias_segment{
        var trait_ interface{} = _alias_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *AliasSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._alias_segment{  list.Add(my._alias_segment) }
        return list
    }

func (my *AliasSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *AliasSeg)       Enter(v Visitor){
        var checkChildren = v.VisitAliasSeg(my)
        if checkChildren{
            if nil != my._alias_segment{my._alias_segment.Accept(v)}
        }
        v.EndVisitAliasSeg(my)
    }


func AnyCastToAliasSeg(i interface{}) *AliasSeg {
	if nil == i{
		return nil
	}else{
		return i.(*AliasSeg)
	}
}
/**
 *<b>
*<li>Rule 5:  LPG_item ::= AST_KEY$ ast_segment END_KEY_OPT$
 *</b>
 */
type AstSeg struct{
    *ASTNode
      _ast_segment *action_segmentList
}
func (my *AstSeg)      Getast_segment() *action_segmentList{ return my._ast_segment}
func (my *AstSeg)      Setast_segment( _ast_segment *action_segmentList)  { my._ast_segment = _ast_segment }

func NewAstSeg(leftIToken IToken, rightIToken IToken ,
              _ast_segment *action_segmentList)*AstSeg{
      my := new(AstSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._ast_segment = _ast_segment;
        if nil != _ast_segment{
        var trait_ interface{} = _ast_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *AstSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._ast_segment{  list.Add(my._ast_segment) }
        return list
    }

func (my *AstSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *AstSeg)       Enter(v Visitor){
        var checkChildren = v.VisitAstSeg(my)
        if checkChildren{
            if nil != my._ast_segment{my._ast_segment.Accept(v)}
        }
        v.EndVisitAstSeg(my)
    }


func AnyCastToAstSeg(i interface{}) *AstSeg {
	if nil == i{
		return nil
	}else{
		return i.(*AstSeg)
	}
}
/**
 *<b>
*<li>Rule 6:  LPG_item ::= DEFINE_KEY$ define_segment END_KEY_OPT$
 *</b>
 */
type DefineSeg struct{
    *ASTNode
      _define_segment *defineSpecList
}
func (my *DefineSeg)      Getdefine_segment() *defineSpecList{ return my._define_segment}
func (my *DefineSeg)      Setdefine_segment( _define_segment *defineSpecList)  { my._define_segment = _define_segment }

func NewDefineSeg(leftIToken IToken, rightIToken IToken ,
              _define_segment *defineSpecList)*DefineSeg{
      my := new(DefineSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._define_segment = _define_segment;
        if nil != _define_segment{
        var trait_ interface{} = _define_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *DefineSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._define_segment{  list.Add(my._define_segment) }
        return list
    }

func (my *DefineSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *DefineSeg)       Enter(v Visitor){
        var checkChildren = v.VisitDefineSeg(my)
        if checkChildren{
            if nil != my._define_segment{my._define_segment.Accept(v)}
        }
        v.EndVisitDefineSeg(my)
    }


func AnyCastToDefineSeg(i interface{}) *DefineSeg {
	if nil == i{
		return nil
	}else{
		return i.(*DefineSeg)
	}
}
/**
 *<b>
*<li>Rule 7:  LPG_item ::= EOF_KEY$ eof_segment END_KEY_OPT$
 *</b>
 */
type EofSeg struct{
    *ASTNode
      _eof_segment Ieof_segment
}
func (my *EofSeg)      Geteof_segment() Ieof_segment{ return my._eof_segment}
func (my *EofSeg)      Seteof_segment( _eof_segment Ieof_segment)  { my._eof_segment = _eof_segment }

func NewEofSeg(leftIToken IToken, rightIToken IToken ,
              _eof_segment Ieof_segment)*EofSeg{
      my := new(EofSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._eof_segment = _eof_segment;
        if nil != _eof_segment{
        var trait_ interface{} = _eof_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *EofSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._eof_segment{  list.Add(my._eof_segment) }
        return list
    }

func (my *EofSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *EofSeg)       Enter(v Visitor){
        var checkChildren = v.VisitEofSeg(my)
        if checkChildren{
            if nil != my._eof_segment{my._eof_segment.Accept(v)}
        }
        v.EndVisitEofSeg(my)
    }


func AnyCastToEofSeg(i interface{}) *EofSeg {
	if nil == i{
		return nil
	}else{
		return i.(*EofSeg)
	}
}
/**
 *<b>
*<li>Rule 8:  LPG_item ::= EOL_KEY$ eol_segment END_KEY_OPT$
 *</b>
 */
type EolSeg struct{
    *ASTNode
      _eol_segment Ieol_segment
}
func (my *EolSeg)      Geteol_segment() Ieol_segment{ return my._eol_segment}
func (my *EolSeg)      Seteol_segment( _eol_segment Ieol_segment)  { my._eol_segment = _eol_segment }

func NewEolSeg(leftIToken IToken, rightIToken IToken ,
              _eol_segment Ieol_segment)*EolSeg{
      my := new(EolSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._eol_segment = _eol_segment;
        if nil != _eol_segment{
        var trait_ interface{} = _eol_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *EolSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._eol_segment{  list.Add(my._eol_segment) }
        return list
    }

func (my *EolSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *EolSeg)       Enter(v Visitor){
        var checkChildren = v.VisitEolSeg(my)
        if checkChildren{
            if nil != my._eol_segment{my._eol_segment.Accept(v)}
        }
        v.EndVisitEolSeg(my)
    }


func AnyCastToEolSeg(i interface{}) *EolSeg {
	if nil == i{
		return nil
	}else{
		return i.(*EolSeg)
	}
}
/**
 *<b>
*<li>Rule 9:  LPG_item ::= ERROR_KEY$ error_segment END_KEY_OPT$
 *</b>
 */
type ErrorSeg struct{
    *ASTNode
      _error_segment Ierror_segment
}
func (my *ErrorSeg)      Geterror_segment() Ierror_segment{ return my._error_segment}
func (my *ErrorSeg)      Seterror_segment( _error_segment Ierror_segment)  { my._error_segment = _error_segment }

func NewErrorSeg(leftIToken IToken, rightIToken IToken ,
              _error_segment Ierror_segment)*ErrorSeg{
      my := new(ErrorSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._error_segment = _error_segment;
        if nil != _error_segment{
        var trait_ interface{} = _error_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *ErrorSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._error_segment{  list.Add(my._error_segment) }
        return list
    }

func (my *ErrorSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *ErrorSeg)       Enter(v Visitor){
        var checkChildren = v.VisitErrorSeg(my)
        if checkChildren{
            if nil != my._error_segment{my._error_segment.Accept(v)}
        }
        v.EndVisitErrorSeg(my)
    }


func AnyCastToErrorSeg(i interface{}) *ErrorSeg {
	if nil == i{
		return nil
	}else{
		return i.(*ErrorSeg)
	}
}
/**
 *<b>
*<li>Rule 10:  LPG_item ::= EXPORT_KEY$ export_segment END_KEY_OPT$
 *</b>
 */
type ExportSeg struct{
    *ASTNode
      _export_segment *terminal_symbolList
}
func (my *ExportSeg)      Getexport_segment() *terminal_symbolList{ return my._export_segment}
func (my *ExportSeg)      Setexport_segment( _export_segment *terminal_symbolList)  { my._export_segment = _export_segment }

func NewExportSeg(leftIToken IToken, rightIToken IToken ,
              _export_segment *terminal_symbolList)*ExportSeg{
      my := new(ExportSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._export_segment = _export_segment;
        if nil != _export_segment{
        var trait_ interface{} = _export_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *ExportSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._export_segment{  list.Add(my._export_segment) }
        return list
    }

func (my *ExportSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *ExportSeg)       Enter(v Visitor){
        var checkChildren = v.VisitExportSeg(my)
        if checkChildren{
            if nil != my._export_segment{my._export_segment.Accept(v)}
        }
        v.EndVisitExportSeg(my)
    }


func AnyCastToExportSeg(i interface{}) *ExportSeg {
	if nil == i{
		return nil
	}else{
		return i.(*ExportSeg)
	}
}
/**
 *<b>
*<li>Rule 11:  LPG_item ::= GLOBALS_KEY$ globals_segment END_KEY_OPT$
 *</b>
 */
type GlobalsSeg struct{
    *ASTNode
      _globals_segment *action_segmentList
}
func (my *GlobalsSeg)      Getglobals_segment() *action_segmentList{ return my._globals_segment}
func (my *GlobalsSeg)      Setglobals_segment( _globals_segment *action_segmentList)  { my._globals_segment = _globals_segment }

func NewGlobalsSeg(leftIToken IToken, rightIToken IToken ,
              _globals_segment *action_segmentList)*GlobalsSeg{
      my := new(GlobalsSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._globals_segment = _globals_segment;
        if nil != _globals_segment{
        var trait_ interface{} = _globals_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *GlobalsSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._globals_segment{  list.Add(my._globals_segment) }
        return list
    }

func (my *GlobalsSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *GlobalsSeg)       Enter(v Visitor){
        var checkChildren = v.VisitGlobalsSeg(my)
        if checkChildren{
            if nil != my._globals_segment{my._globals_segment.Accept(v)}
        }
        v.EndVisitGlobalsSeg(my)
    }


func AnyCastToGlobalsSeg(i interface{}) *GlobalsSeg {
	if nil == i{
		return nil
	}else{
		return i.(*GlobalsSeg)
	}
}
/**
 *<b>
*<li>Rule 12:  LPG_item ::= HEADERS_KEY$ headers_segment END_KEY_OPT$
 *</b>
 */
type HeadersSeg struct{
    *ASTNode
      _headers_segment *action_segmentList
}
func (my *HeadersSeg)      Getheaders_segment() *action_segmentList{ return my._headers_segment}
func (my *HeadersSeg)      Setheaders_segment( _headers_segment *action_segmentList)  { my._headers_segment = _headers_segment }

func NewHeadersSeg(leftIToken IToken, rightIToken IToken ,
              _headers_segment *action_segmentList)*HeadersSeg{
      my := new(HeadersSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._headers_segment = _headers_segment;
        if nil != _headers_segment{
        var trait_ interface{} = _headers_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *HeadersSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._headers_segment{  list.Add(my._headers_segment) }
        return list
    }

func (my *HeadersSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *HeadersSeg)       Enter(v Visitor){
        var checkChildren = v.VisitHeadersSeg(my)
        if checkChildren{
            if nil != my._headers_segment{my._headers_segment.Accept(v)}
        }
        v.EndVisitHeadersSeg(my)
    }


func AnyCastToHeadersSeg(i interface{}) *HeadersSeg {
	if nil == i{
		return nil
	}else{
		return i.(*HeadersSeg)
	}
}
/**
 *<b>
*<li>Rule 13:  LPG_item ::= IDENTIFIER_KEY$ identifier_segment END_KEY_OPT$
 *</b>
 */
type IdentifierSeg struct{
    *ASTNode
      _identifier_segment Iidentifier_segment
}
func (my *IdentifierSeg)      Getidentifier_segment() Iidentifier_segment{ return my._identifier_segment}
func (my *IdentifierSeg)      Setidentifier_segment( _identifier_segment Iidentifier_segment)  { my._identifier_segment = _identifier_segment }

func NewIdentifierSeg(leftIToken IToken, rightIToken IToken ,
              _identifier_segment Iidentifier_segment)*IdentifierSeg{
      my := new(IdentifierSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._identifier_segment = _identifier_segment;
        if nil != _identifier_segment{
        var trait_ interface{} = _identifier_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *IdentifierSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._identifier_segment{  list.Add(my._identifier_segment) }
        return list
    }

func (my *IdentifierSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *IdentifierSeg)       Enter(v Visitor){
        var checkChildren = v.VisitIdentifierSeg(my)
        if checkChildren{
            if nil != my._identifier_segment{my._identifier_segment.Accept(v)}
        }
        v.EndVisitIdentifierSeg(my)
    }


func AnyCastToIdentifierSeg(i interface{}) *IdentifierSeg {
	if nil == i{
		return nil
	}else{
		return i.(*IdentifierSeg)
	}
}
/**
 *<b>
*<li>Rule 14:  LPG_item ::= IMPORT_KEY$ import_segment END_KEY_OPT$
 *</b>
 */
type ImportSeg struct{
    *ASTNode
      _import_segment *import_segment
}
func (my *ImportSeg)      Getimport_segment() *import_segment{ return my._import_segment}
func (my *ImportSeg)      Setimport_segment( _import_segment *import_segment)  { my._import_segment = _import_segment }

func NewImportSeg(leftIToken IToken, rightIToken IToken ,
              _import_segment *import_segment)*ImportSeg{
      my := new(ImportSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._import_segment = _import_segment;
        if nil != _import_segment{
        var trait_ interface{} = _import_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *ImportSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._import_segment{  list.Add(my._import_segment) }
        return list
    }

func (my *ImportSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *ImportSeg)       Enter(v Visitor){
        var checkChildren = v.VisitImportSeg(my)
        if checkChildren{
            if nil != my._import_segment{my._import_segment.Accept(v)}
        }
        v.EndVisitImportSeg(my)
    }


func AnyCastToImportSeg(i interface{}) *ImportSeg {
	if nil == i{
		return nil
	}else{
		return i.(*ImportSeg)
	}
}
/**
 *<b>
*<li>Rule 15:  LPG_item ::= INCLUDE_KEY$ include_segment END_KEY_OPT$
 *</b>
 */
type IncludeSeg struct{
    *ASTNode
      _include_segment *include_segment
}
func (my *IncludeSeg)      Getinclude_segment() *include_segment{ return my._include_segment}
func (my *IncludeSeg)      Setinclude_segment( _include_segment *include_segment)  { my._include_segment = _include_segment }

func NewIncludeSeg(leftIToken IToken, rightIToken IToken ,
              _include_segment *include_segment)*IncludeSeg{
      my := new(IncludeSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._include_segment = _include_segment;
        if nil != _include_segment{
        var trait_ interface{} = _include_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *IncludeSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._include_segment{  list.Add(my._include_segment) }
        return list
    }

func (my *IncludeSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *IncludeSeg)       Enter(v Visitor){
        var checkChildren = v.VisitIncludeSeg(my)
        if checkChildren{
            if nil != my._include_segment{my._include_segment.Accept(v)}
        }
        v.EndVisitIncludeSeg(my)
    }


func AnyCastToIncludeSeg(i interface{}) *IncludeSeg {
	if nil == i{
		return nil
	}else{
		return i.(*IncludeSeg)
	}
}
/**
 *<b>
*<li>Rule 16:  LPG_item ::= KEYWORDS_KEY$ keywords_segment END_KEY_OPT$
 *</b>
 */
type KeywordsSeg struct{
    *ASTNode
      _keywords_segment *keywordSpecList
}
func (my *KeywordsSeg)      Getkeywords_segment() *keywordSpecList{ return my._keywords_segment}
func (my *KeywordsSeg)      Setkeywords_segment( _keywords_segment *keywordSpecList)  { my._keywords_segment = _keywords_segment }

func NewKeywordsSeg(leftIToken IToken, rightIToken IToken ,
              _keywords_segment *keywordSpecList)*KeywordsSeg{
      my := new(KeywordsSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._keywords_segment = _keywords_segment;
        if nil != _keywords_segment{
        var trait_ interface{} = _keywords_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *KeywordsSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._keywords_segment{  list.Add(my._keywords_segment) }
        return list
    }

func (my *KeywordsSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *KeywordsSeg)       Enter(v Visitor){
        var checkChildren = v.VisitKeywordsSeg(my)
        if checkChildren{
            if nil != my._keywords_segment{my._keywords_segment.Accept(v)}
        }
        v.EndVisitKeywordsSeg(my)
    }


func AnyCastToKeywordsSeg(i interface{}) *KeywordsSeg {
	if nil == i{
		return nil
	}else{
		return i.(*KeywordsSeg)
	}
}
/**
 *<b>
*<li>Rule 17:  LPG_item ::= NAMES_KEY$ names_segment END_KEY_OPT$
 *</b>
 */
type NamesSeg struct{
    *ASTNode
      _names_segment *nameSpecList
}
func (my *NamesSeg)      Getnames_segment() *nameSpecList{ return my._names_segment}
func (my *NamesSeg)      Setnames_segment( _names_segment *nameSpecList)  { my._names_segment = _names_segment }

func NewNamesSeg(leftIToken IToken, rightIToken IToken ,
              _names_segment *nameSpecList)*NamesSeg{
      my := new(NamesSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._names_segment = _names_segment;
        if nil != _names_segment{
        var trait_ interface{} = _names_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *NamesSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._names_segment{  list.Add(my._names_segment) }
        return list
    }

func (my *NamesSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *NamesSeg)       Enter(v Visitor){
        var checkChildren = v.VisitNamesSeg(my)
        if checkChildren{
            if nil != my._names_segment{my._names_segment.Accept(v)}
        }
        v.EndVisitNamesSeg(my)
    }


func AnyCastToNamesSeg(i interface{}) *NamesSeg {
	if nil == i{
		return nil
	}else{
		return i.(*NamesSeg)
	}
}
/**
 *<b>
*<li>Rule 18:  LPG_item ::= NOTICE_KEY$ notice_segment END_KEY_OPT$
 *</b>
 */
type NoticeSeg struct{
    *ASTNode
      _notice_segment *action_segmentList
}
func (my *NoticeSeg)      Getnotice_segment() *action_segmentList{ return my._notice_segment}
func (my *NoticeSeg)      Setnotice_segment( _notice_segment *action_segmentList)  { my._notice_segment = _notice_segment }

func NewNoticeSeg(leftIToken IToken, rightIToken IToken ,
              _notice_segment *action_segmentList)*NoticeSeg{
      my := new(NoticeSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._notice_segment = _notice_segment;
        if nil != _notice_segment{
        var trait_ interface{} = _notice_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *NoticeSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._notice_segment{  list.Add(my._notice_segment) }
        return list
    }

func (my *NoticeSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *NoticeSeg)       Enter(v Visitor){
        var checkChildren = v.VisitNoticeSeg(my)
        if checkChildren{
            if nil != my._notice_segment{my._notice_segment.Accept(v)}
        }
        v.EndVisitNoticeSeg(my)
    }


func AnyCastToNoticeSeg(i interface{}) *NoticeSeg {
	if nil == i{
		return nil
	}else{
		return i.(*NoticeSeg)
	}
}
/**
 *<b>
*<li>Rule 19:  LPG_item ::= RULES_KEY$ rules_segment END_KEY_OPT$
 *</b>
 */
type RulesSeg struct{
    *ASTNode
      _rules_segment *rules_segment
}
func (my *RulesSeg)      Getrules_segment() *rules_segment{ return my._rules_segment}
func (my *RulesSeg)      Setrules_segment( _rules_segment *rules_segment)  { my._rules_segment = _rules_segment }

func NewRulesSeg(leftIToken IToken, rightIToken IToken ,
              _rules_segment *rules_segment)*RulesSeg{
      my := new(RulesSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._rules_segment = _rules_segment;
        if nil != _rules_segment{
        var trait_ interface{} = _rules_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *RulesSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._rules_segment{  list.Add(my._rules_segment) }
        return list
    }

func (my *RulesSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *RulesSeg)       Enter(v Visitor){
        var checkChildren = v.VisitRulesSeg(my)
        if checkChildren{
            if nil != my._rules_segment{my._rules_segment.Accept(v)}
        }
        v.EndVisitRulesSeg(my)
    }


func AnyCastToRulesSeg(i interface{}) *RulesSeg {
	if nil == i{
		return nil
	}else{
		return i.(*RulesSeg)
	}
}
/**
 *<b>
*<li>Rule 20:  LPG_item ::= SOFT_KEYWORDS_KEY$ keywords_segment END_KEY_OPT$
 *</b>
 */
type SoftKeywordsSeg struct{
    *ASTNode
      _keywords_segment *keywordSpecList
}
func (my *SoftKeywordsSeg)      Getkeywords_segment() *keywordSpecList{ return my._keywords_segment}
func (my *SoftKeywordsSeg)      Setkeywords_segment( _keywords_segment *keywordSpecList)  { my._keywords_segment = _keywords_segment }

func NewSoftKeywordsSeg(leftIToken IToken, rightIToken IToken ,
              _keywords_segment *keywordSpecList)*SoftKeywordsSeg{
      my := new(SoftKeywordsSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._keywords_segment = _keywords_segment;
        if nil != _keywords_segment{
        var trait_ interface{} = _keywords_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *SoftKeywordsSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._keywords_segment{  list.Add(my._keywords_segment) }
        return list
    }

func (my *SoftKeywordsSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *SoftKeywordsSeg)       Enter(v Visitor){
        var checkChildren = v.VisitSoftKeywordsSeg(my)
        if checkChildren{
            if nil != my._keywords_segment{my._keywords_segment.Accept(v)}
        }
        v.EndVisitSoftKeywordsSeg(my)
    }


func AnyCastToSoftKeywordsSeg(i interface{}) *SoftKeywordsSeg {
	if nil == i{
		return nil
	}else{
		return i.(*SoftKeywordsSeg)
	}
}
/**
 *<b>
*<li>Rule 21:  LPG_item ::= START_KEY$ start_segment END_KEY_OPT$
 *</b>
 */
type StartSeg struct{
    *ASTNode
      _start_segment *start_symbolList
}
func (my *StartSeg)      Getstart_segment() *start_symbolList{ return my._start_segment}
func (my *StartSeg)      Setstart_segment( _start_segment *start_symbolList)  { my._start_segment = _start_segment }

func NewStartSeg(leftIToken IToken, rightIToken IToken ,
              _start_segment *start_symbolList)*StartSeg{
      my := new(StartSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._start_segment = _start_segment;
        if nil != _start_segment{
        var trait_ interface{} = _start_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *StartSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._start_segment{  list.Add(my._start_segment) }
        return list
    }

func (my *StartSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *StartSeg)       Enter(v Visitor){
        var checkChildren = v.VisitStartSeg(my)
        if checkChildren{
            if nil != my._start_segment{my._start_segment.Accept(v)}
        }
        v.EndVisitStartSeg(my)
    }


func AnyCastToStartSeg(i interface{}) *StartSeg {
	if nil == i{
		return nil
	}else{
		return i.(*StartSeg)
	}
}
/**
 *<b>
*<li>Rule 22:  LPG_item ::= TERMINALS_KEY$ terminals_segment END_KEY_OPT$
 *</b>
 */
type TerminalsSeg struct{
    *ASTNode
      _terminals_segment *terminalList
}
func (my *TerminalsSeg)      Getterminals_segment() *terminalList{ return my._terminals_segment}
func (my *TerminalsSeg)      Setterminals_segment( _terminals_segment *terminalList)  { my._terminals_segment = _terminals_segment }

func NewTerminalsSeg(leftIToken IToken, rightIToken IToken ,
              _terminals_segment *terminalList)*TerminalsSeg{
      my := new(TerminalsSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._terminals_segment = _terminals_segment;
        if nil != _terminals_segment{
        var trait_ interface{} = _terminals_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *TerminalsSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._terminals_segment{  list.Add(my._terminals_segment) }
        return list
    }

func (my *TerminalsSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *TerminalsSeg)       Enter(v Visitor){
        var checkChildren = v.VisitTerminalsSeg(my)
        if checkChildren{
            if nil != my._terminals_segment{my._terminals_segment.Accept(v)}
        }
        v.EndVisitTerminalsSeg(my)
    }


func AnyCastToTerminalsSeg(i interface{}) *TerminalsSeg {
	if nil == i{
		return nil
	}else{
		return i.(*TerminalsSeg)
	}
}
/**
 *<b>
*<li>Rule 23:  LPG_item ::= TRAILERS_KEY$ trailers_segment END_KEY_OPT$
 *</b>
 */
type TrailersSeg struct{
    *ASTNode
      _trailers_segment *action_segmentList
}
func (my *TrailersSeg)      Gettrailers_segment() *action_segmentList{ return my._trailers_segment}
func (my *TrailersSeg)      Settrailers_segment( _trailers_segment *action_segmentList)  { my._trailers_segment = _trailers_segment }

func NewTrailersSeg(leftIToken IToken, rightIToken IToken ,
              _trailers_segment *action_segmentList)*TrailersSeg{
      my := new(TrailersSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._trailers_segment = _trailers_segment;
        if nil != _trailers_segment{
        var trait_ interface{} = _trailers_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *TrailersSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._trailers_segment{  list.Add(my._trailers_segment) }
        return list
    }

func (my *TrailersSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *TrailersSeg)       Enter(v Visitor){
        var checkChildren = v.VisitTrailersSeg(my)
        if checkChildren{
            if nil != my._trailers_segment{my._trailers_segment.Accept(v)}
        }
        v.EndVisitTrailersSeg(my)
    }


func AnyCastToTrailersSeg(i interface{}) *TrailersSeg {
	if nil == i{
		return nil
	}else{
		return i.(*TrailersSeg)
	}
}
/**
 *<b>
*<li>Rule 24:  LPG_item ::= TYPES_KEY$ types_segment END_KEY_OPT$
 *</b>
 */
type TypesSeg struct{
    *ASTNode
      _types_segment *type_declarationsList
}
func (my *TypesSeg)      Gettypes_segment() *type_declarationsList{ return my._types_segment}
func (my *TypesSeg)      Settypes_segment( _types_segment *type_declarationsList)  { my._types_segment = _types_segment }

func NewTypesSeg(leftIToken IToken, rightIToken IToken ,
              _types_segment *type_declarationsList)*TypesSeg{
      my := new(TypesSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._types_segment = _types_segment;
        if nil != _types_segment{
        var trait_ interface{} = _types_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *TypesSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._types_segment{  list.Add(my._types_segment) }
        return list
    }

func (my *TypesSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *TypesSeg)       Enter(v Visitor){
        var checkChildren = v.VisitTypesSeg(my)
        if checkChildren{
            if nil != my._types_segment{my._types_segment.Accept(v)}
        }
        v.EndVisitTypesSeg(my)
    }


func AnyCastToTypesSeg(i interface{}) *TypesSeg {
	if nil == i{
		return nil
	}else{
		return i.(*TypesSeg)
	}
}
/**
 *<b>
*<li>Rule 25:  LPG_item ::= RECOVER_KEY$ recover_segment END_KEY_OPT$
 *</b>
 */
type RecoverSeg struct{
    *ASTNode
      _recover_segment *SYMBOLList
}
func (my *RecoverSeg)      Getrecover_segment() *SYMBOLList{ return my._recover_segment}
func (my *RecoverSeg)      Setrecover_segment( _recover_segment *SYMBOLList)  { my._recover_segment = _recover_segment }

func NewRecoverSeg(leftIToken IToken, rightIToken IToken ,
              _recover_segment *SYMBOLList)*RecoverSeg{
      my := new(RecoverSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._recover_segment = _recover_segment;
        if nil != _recover_segment{
        var trait_ interface{} = _recover_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *RecoverSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._recover_segment{  list.Add(my._recover_segment) }
        return list
    }

func (my *RecoverSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *RecoverSeg)       Enter(v Visitor){
        var checkChildren = v.VisitRecoverSeg(my)
        if checkChildren{
            if nil != my._recover_segment{my._recover_segment.Accept(v)}
        }
        v.EndVisitRecoverSeg(my)
    }


func AnyCastToRecoverSeg(i interface{}) *RecoverSeg {
	if nil == i{
		return nil
	}else{
		return i.(*RecoverSeg)
	}
}
/**
 *<b>
*<li>Rule 26:  LPG_item ::= DISJOINTPREDECESSORSETS_KEY$ predecessor_segment END_KEY_OPT$
 *</b>
 */
type PredecessorSeg struct{
    *ASTNode
      _predecessor_segment *symbol_pairList
}
func (my *PredecessorSeg)      Getpredecessor_segment() *symbol_pairList{ return my._predecessor_segment}
func (my *PredecessorSeg)      Setpredecessor_segment( _predecessor_segment *symbol_pairList)  { my._predecessor_segment = _predecessor_segment }

func NewPredecessorSeg(leftIToken IToken, rightIToken IToken ,
              _predecessor_segment *symbol_pairList)*PredecessorSeg{
      my := new(PredecessorSeg)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._predecessor_segment = _predecessor_segment;
        if nil != _predecessor_segment{
        var trait_ interface{} = _predecessor_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *PredecessorSeg)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._predecessor_segment{  list.Add(my._predecessor_segment) }
        return list
    }

func (my *PredecessorSeg)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *PredecessorSeg)       Enter(v Visitor){
        var checkChildren = v.VisitPredecessorSeg(my)
        if checkChildren{
            if nil != my._predecessor_segment{my._predecessor_segment.Accept(v)}
        }
        v.EndVisitPredecessorSeg(my)
    }


func AnyCastToPredecessorSeg(i interface{}) *PredecessorSeg {
	if nil == i{
		return nil
	}else{
		return i.(*PredecessorSeg)
	}
}
/**
 *<b>
*<li>Rule 27:  options_segment ::= $Empty
*<li>Rule 28:  options_segment ::= options_segment option_spec
 *</b>
 */
type option_specList struct{
    *AbstractASTNodeList
}
func (my *option_specList)      Getoption_specAt(i int) *option_spec{
     var r,_=my.GetElementAt(i).(*option_spec)
     return r
     }

func Newoption_specList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*option_specList{
      my := new(option_specList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  Newoption_specListFromElement(element *option_spec,leftRecursive bool)*option_specList{
        var obj = Newoption_specList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *option_specList)      AddElement(_option_spec IAst){ 
      my.AbstractASTNodeList.AddElement(_option_spec)
        _option_spec.SetParent(my)
    }


func (my *option_specList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *option_specList)     Enter(v  Visitor){
        var checkChildren = v.Visitoption_specList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.Getoption_specAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisitoption_specList(my)
    }


func AnyCastTooption_specList(i interface{}) *option_specList {
	if nil == i{
		return nil
	}else{
		return i.(*option_specList)
	}
}
/**
 *<b>
*<li>Rule 29:  option_spec ::= OPTIONS_KEY$ option_list
 *</b>
 */
type option_spec struct{
    *ASTNode
      _option_list *optionList
}
func (my *option_spec)      Getoption_list() *optionList{ return my._option_list}
func (my *option_spec)      Setoption_list( _option_list *optionList)  { my._option_list = _option_list }

func Newoption_spec(leftIToken IToken, rightIToken IToken ,
              _option_list *optionList)*option_spec{
      my := new(option_spec)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._option_list = _option_list;
        if nil != _option_list{
        var trait_ interface{} = _option_list
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *option_spec)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._option_list{  list.Add(my._option_list) }
        return list
    }

func (my *option_spec)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *option_spec)       Enter(v Visitor){
        var checkChildren = v.Visitoption_spec(my)
        if checkChildren{
            if nil != my._option_list{my._option_list.Accept(v)}
        }
        v.EndVisitoption_spec(my)
    }


func AnyCastTooption_spec(i interface{}) *option_spec {
	if nil == i{
		return nil
	}else{
		return i.(*option_spec)
	}
}
/**
 *<b>
*<li>Rule 30:  option_list ::= option
*<li>Rule 31:  option_list ::= option_list ,$ option
 *</b>
 */
type optionList struct{
    *AbstractASTNodeList
}
func (my *optionList)      GetoptionAt(i int) *option{
     var r,_=my.GetElementAt(i).(*option)
     return r
     }

func NewoptionList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*optionList{
      my := new(optionList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  NewoptionListFromElement(element *option,leftRecursive bool)*optionList{
        var obj = NewoptionList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *optionList)      AddElement(_option IAst){ 
      my.AbstractASTNodeList.AddElement(_option)
        _option.SetParent(my)
    }


func (my *optionList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *optionList)     Enter(v  Visitor){
        var checkChildren = v.VisitoptionList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.GetoptionAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisitoptionList(my)
    }


func AnyCastTooptionList(i interface{}) *optionList {
	if nil == i{
		return nil
	}else{
		return i.(*optionList)
	}
}
/**
 *<b>
*<li>Rule 32:  option ::= SYMBOL option_value
 *</b>
 */
type option struct{
    *ASTNode
      _SYMBOL *ASTNodeToken
      _option_value Ioption_value
}
func (my *option)      GetSYMBOL() *ASTNodeToken{ return my._SYMBOL}
func (my *option)      SetSYMBOL( _SYMBOL *ASTNodeToken)  { my._SYMBOL = _SYMBOL }
    /**
     * The value returned by <b>Getoption_value</b> may be <b>null</b>
     */
func (my *option)      Getoption_value() Ioption_value{ return my._option_value}
func (my *option)      Setoption_value( _option_value Ioption_value)  { my._option_value = _option_value }

func Newoption(leftIToken IToken, rightIToken IToken ,
              _SYMBOL *ASTNodeToken,
              _option_value Ioption_value)*option{
      my := new(option)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._SYMBOL = _SYMBOL;
        if nil != _SYMBOL{
        var trait_ interface{} = _SYMBOL
         trait_.(IAst).SetParent(my)
}
        my._option_value = _option_value;
        if nil != _option_value{
        var trait_ interface{} = _option_value
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *option)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._SYMBOL{  list.Add(my._SYMBOL) }
        if nil != my._option_value{  list.Add(my._option_value) }
        return list
    }

func (my *option)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *option)       Enter(v Visitor){
        var checkChildren = v.Visitoption(my)
        if checkChildren{
            if nil != my._SYMBOL{my._SYMBOL.Accept(v)}
            if nil != my._option_value{my._option_value.Accept(v)}
        }
        v.EndVisitoption(my)
    }


func AnyCastTooption(i interface{}) *option {
	if nil == i{
		return nil
	}else{
		return i.(*option)
	}
}
/**
 *<b>
*<li>Rule 36:  symbol_list ::= SYMBOL
*<li>Rule 37:  symbol_list ::= symbol_list ,$ SYMBOL
*<li>Rule 75:  drop_symbols ::= SYMBOL
*<li>Rule 76:  drop_symbols ::= drop_symbols SYMBOL
*<li>Rule 136:  barSymbolList ::= SYMBOL
*<li>Rule 137:  barSymbolList ::= barSymbolList |$ SYMBOL
*<li>Rule 141:  recover_segment ::= $Empty
*<li>Rule 142:  recover_segment ::= recover_segment recover_symbol
 *</b>
 */
type SYMBOLList struct{
    *AbstractASTNodeList
}
func (my *SYMBOLList)      GetSYMBOLAt(i int) *ASTNodeToken{
     var r,_=my.GetElementAt(i).(*ASTNodeToken)
     return r
     }

func NewSYMBOLList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*SYMBOLList{
      my := new(SYMBOLList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  NewSYMBOLListFromElement(element *ASTNodeToken,leftRecursive bool)*SYMBOLList{
        var obj = NewSYMBOLList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *SYMBOLList)      AddElement(_SYMBOL IAst){ 
      my.AbstractASTNodeList.AddElement(_SYMBOL)
        _SYMBOL.SetParent(my)
    }


func (my *SYMBOLList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *SYMBOLList)     Enter(v  Visitor){
        var checkChildren = v.VisitSYMBOLList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.GetSYMBOLAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisitSYMBOLList(my)
    }


func AnyCastToSYMBOLList(i interface{}) *SYMBOLList {
	if nil == i{
		return nil
	}else{
		return i.(*SYMBOLList)
	}
}
/**
 *<b>
*<li>Rule 38:  alias_segment ::= aliasSpec
*<li>Rule 39:  alias_segment ::= alias_segment aliasSpec
 *</b>
 */
type aliasSpecList struct{
    *AbstractASTNodeList
}
func (my *aliasSpecList)      GetaliasSpecAt(i int) IaliasSpec{
     var r,_=my.GetElementAt(i).(IaliasSpec)
     return r
     }

func NewaliasSpecList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*aliasSpecList{
      my := new(aliasSpecList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  NewaliasSpecListFromElement(element IaliasSpec,leftRecursive bool)*aliasSpecList{
        var obj = NewaliasSpecList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *aliasSpecList)      AddElement(_aliasSpec IAst){ 
      my.AbstractASTNodeList.AddElement(_aliasSpec)
        _aliasSpec.SetParent(my)
    }


func (my *aliasSpecList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *aliasSpecList)     Enter(v  Visitor){
        var checkChildren = v.VisitaliasSpecList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.GetaliasSpecAt(i)
                if nil !=element{                    element.Accept(v)
                }
            }
        }
        v.EndVisitaliasSpecList(my)
    }


func AnyCastToaliasSpecList(i interface{}) *aliasSpecList {
	if nil == i{
		return nil
	}else{
		return i.(*aliasSpecList)
	}
}
/**
 *<b>
*<li>Rule 46:  alias_lhs_macro_name ::= MACRO_NAME
 *</b>
 */
type alias_lhs_macro_name struct{
    *ASTNodeToken
}
func (my *alias_lhs_macro_name)      GetMACRO_NAME()IToken{ return my.leftIToken; }

func Newalias_lhs_macro_name(token IToken )*alias_lhs_macro_name{
      my := new(alias_lhs_macro_name)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *alias_lhs_macro_name)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *alias_lhs_macro_name)       Enter(v Visitor){
        v.Visitalias_lhs_macro_name(my)
        v.EndVisitalias_lhs_macro_name(my)
    }


func AnyCastToalias_lhs_macro_name(i interface{}) *alias_lhs_macro_name {
	if nil == i{
		return nil
	}else{
		return i.(*alias_lhs_macro_name)
	}
}
/**
 *<b>
*<li>Rule 55:  define_segment ::= defineSpec
*<li>Rule 56:  define_segment ::= define_segment defineSpec
 *</b>
 */
type defineSpecList struct{
    *AbstractASTNodeList
}
func (my *defineSpecList)      GetdefineSpecAt(i int) *defineSpec{
     var r,_=my.GetElementAt(i).(*defineSpec)
     return r
     }

func NewdefineSpecList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*defineSpecList{
      my := new(defineSpecList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  NewdefineSpecListFromElement(element *defineSpec,leftRecursive bool)*defineSpecList{
        var obj = NewdefineSpecList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *defineSpecList)      AddElement(_defineSpec IAst){ 
      my.AbstractASTNodeList.AddElement(_defineSpec)
        _defineSpec.SetParent(my)
    }


func (my *defineSpecList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *defineSpecList)     Enter(v  Visitor){
        var checkChildren = v.VisitdefineSpecList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.GetdefineSpecAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisitdefineSpecList(my)
    }


func AnyCastTodefineSpecList(i interface{}) *defineSpecList {
	if nil == i{
		return nil
	}else{
		return i.(*defineSpecList)
	}
}
/**
 *<b>
*<li>Rule 57:  defineSpec ::= macro_name_symbol macro_segment
 *</b>
 */
type defineSpec struct{
    *ASTNode
      _macro_name_symbol Imacro_name_symbol
      _macro_segment *macro_segment
}
func (my *defineSpec)      Getmacro_name_symbol() Imacro_name_symbol{ return my._macro_name_symbol}
func (my *defineSpec)      Setmacro_name_symbol( _macro_name_symbol Imacro_name_symbol)  { my._macro_name_symbol = _macro_name_symbol }
func (my *defineSpec)      Getmacro_segment() *macro_segment{ return my._macro_segment}
func (my *defineSpec)      Setmacro_segment( _macro_segment *macro_segment)  { my._macro_segment = _macro_segment }

func NewdefineSpec(leftIToken IToken, rightIToken IToken ,
              _macro_name_symbol Imacro_name_symbol,
              _macro_segment *macro_segment)*defineSpec{
      my := new(defineSpec)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._macro_name_symbol = _macro_name_symbol;
        if nil != _macro_name_symbol{
        var trait_ interface{} = _macro_name_symbol
         trait_.(IAst).SetParent(my)
}
        my._macro_segment = _macro_segment;
        if nil != _macro_segment{
        var trait_ interface{} = _macro_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *defineSpec)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._macro_name_symbol{  list.Add(my._macro_name_symbol) }
        if nil != my._macro_segment{  list.Add(my._macro_segment) }
        return list
    }

func (my *defineSpec)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *defineSpec)       Enter(v Visitor){
        var checkChildren = v.VisitdefineSpec(my)
        if checkChildren{
            if nil != my._macro_name_symbol{my._macro_name_symbol.Accept(v)}
            if nil != my._macro_segment{my._macro_segment.Accept(v)}
        }
        v.EndVisitdefineSpec(my)
    }


func AnyCastTodefineSpec(i interface{}) *defineSpec {
	if nil == i{
		return nil
	}else{
		return i.(*defineSpec)
	}
}
/**
 *<b>
*<li>Rule 60:  macro_segment ::= BLOCK
 *</b>
 */
type macro_segment struct{
    *ASTNodeToken
}
func (my *macro_segment)      GetBLOCK()IToken{ return my.leftIToken; }

func Newmacro_segment(token IToken )*macro_segment{
      my := new(macro_segment)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *macro_segment)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *macro_segment)       Enter(v Visitor){
        v.Visitmacro_segment(my)
        v.EndVisitmacro_segment(my)
    }


func AnyCastTomacro_segment(i interface{}) *macro_segment {
	if nil == i{
		return nil
	}else{
		return i.(*macro_segment)
	}
}
/**
 *<b>
*<li>Rule 64:  export_segment ::= terminal_symbol
*<li>Rule 65:  export_segment ::= export_segment terminal_symbol
 *</b>
 */
type terminal_symbolList struct{
    *AbstractASTNodeList
}
func (my *terminal_symbolList)      Getterminal_symbolAt(i int) Iterminal_symbol{
     var r,_=my.GetElementAt(i).(Iterminal_symbol)
     return r
     }

func Newterminal_symbolList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*terminal_symbolList{
      my := new(terminal_symbolList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  Newterminal_symbolListFromElement(element Iterminal_symbol,leftRecursive bool)*terminal_symbolList{
        var obj = Newterminal_symbolList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *terminal_symbolList)      AddElement(_terminal_symbol IAst){ 
      my.AbstractASTNodeList.AddElement(_terminal_symbol)
        _terminal_symbol.SetParent(my)
    }


func (my *terminal_symbolList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *terminal_symbolList)     Enter(v  Visitor){
        var checkChildren = v.Visitterminal_symbolList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.Getterminal_symbolAt(i)
                if nil !=element{                    element.Accept(v)
                }
            }
        }
        v.EndVisitterminal_symbolList(my)
    }


func AnyCastToterminal_symbolList(i interface{}) *terminal_symbolList {
	if nil == i{
		return nil
	}else{
		return i.(*terminal_symbolList)
	}
}
/**
 *<b>
*<li>Rule 66:  globals_segment ::= action_segment
*<li>Rule 67:  globals_segment ::= globals_segment action_segment
*<li>Rule 96:  notice_segment ::= action_segment
*<li>Rule 97:  notice_segment ::= notice_segment action_segment
*<li>Rule 146:  action_segment_list ::= $Empty
*<li>Rule 147:  action_segment_list ::= action_segment_list action_segment
 *</b>
 */
type action_segmentList struct{
    *AbstractASTNodeList
}
func (my *action_segmentList)      Getaction_segmentAt(i int) *action_segment{
     var r,_=my.GetElementAt(i).(*action_segment)
     return r
     }

func Newaction_segmentList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*action_segmentList{
      my := new(action_segmentList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  Newaction_segmentListFromElement(element *action_segment,leftRecursive bool)*action_segmentList{
        var obj = Newaction_segmentList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *action_segmentList)      AddElement(_action_segment IAst){ 
      my.AbstractASTNodeList.AddElement(_action_segment)
        _action_segment.SetParent(my)
    }


func (my *action_segmentList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *action_segmentList)     Enter(v  Visitor){
        var checkChildren = v.Visitaction_segmentList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.Getaction_segmentAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisitaction_segmentList(my)
    }


func AnyCastToaction_segmentList(i interface{}) *action_segmentList {
	if nil == i{
		return nil
	}else{
		return i.(*action_segmentList)
	}
}
/**
 *<b>
*<li>Rule 70:  import_segment ::= SYMBOL drop_command_list
 *</b>
 */
type import_segment struct{
    *ASTNode
      _SYMBOL *ASTNodeToken
      _drop_command_list *drop_commandList
}
func (my *import_segment)      GetSYMBOL() *ASTNodeToken{ return my._SYMBOL}
func (my *import_segment)      SetSYMBOL( _SYMBOL *ASTNodeToken)  { my._SYMBOL = _SYMBOL }
func (my *import_segment)      Getdrop_command_list() *drop_commandList{ return my._drop_command_list}
func (my *import_segment)      Setdrop_command_list( _drop_command_list *drop_commandList)  { my._drop_command_list = _drop_command_list }

func Newimport_segment(leftIToken IToken, rightIToken IToken ,
              _SYMBOL *ASTNodeToken,
              _drop_command_list *drop_commandList)*import_segment{
      my := new(import_segment)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._SYMBOL = _SYMBOL;
        if nil != _SYMBOL{
        var trait_ interface{} = _SYMBOL
         trait_.(IAst).SetParent(my)
}
        my._drop_command_list = _drop_command_list;
        if nil != _drop_command_list{
        var trait_ interface{} = _drop_command_list
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *import_segment)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._SYMBOL{  list.Add(my._SYMBOL) }
        if nil != my._drop_command_list{  list.Add(my._drop_command_list) }
        return list
    }

func (my *import_segment)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *import_segment)       Enter(v Visitor){
        var checkChildren = v.Visitimport_segment(my)
        if checkChildren{
            if nil != my._SYMBOL{my._SYMBOL.Accept(v)}
            if nil != my._drop_command_list{my._drop_command_list.Accept(v)}
        }
        v.EndVisitimport_segment(my)
    }


func AnyCastToimport_segment(i interface{}) *import_segment {
	if nil == i{
		return nil
	}else{
		return i.(*import_segment)
	}
}
/**
 *<b>
*<li>Rule 71:  drop_command_list ::= $Empty
*<li>Rule 72:  drop_command_list ::= drop_command_list drop_command
 *</b>
 */
type drop_commandList struct{
    *AbstractASTNodeList
}
func (my *drop_commandList)      Getdrop_commandAt(i int) Idrop_command{
     var r,_=my.GetElementAt(i).(Idrop_command)
     return r
     }

func Newdrop_commandList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*drop_commandList{
      my := new(drop_commandList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  Newdrop_commandListFromElement(element Idrop_command,leftRecursive bool)*drop_commandList{
        var obj = Newdrop_commandList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *drop_commandList)      AddElement(_drop_command IAst){ 
      my.AbstractASTNodeList.AddElement(_drop_command)
        _drop_command.SetParent(my)
    }


func (my *drop_commandList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *drop_commandList)     Enter(v  Visitor){
        var checkChildren = v.Visitdrop_commandList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.Getdrop_commandAt(i)
                if nil !=element{                    element.Accept(v)
                }
            }
        }
        v.EndVisitdrop_commandList(my)
    }


func AnyCastTodrop_commandList(i interface{}) *drop_commandList {
	if nil == i{
		return nil
	}else{
		return i.(*drop_commandList)
	}
}
/**
 *<b>
*<li>Rule 77:  drop_rules ::= drop_rule
*<li>Rule 78:  drop_rules ::= drop_rules drop_rule
 *</b>
 */
type drop_ruleList struct{
    *AbstractASTNodeList
}
func (my *drop_ruleList)      Getdrop_ruleAt(i int) *drop_rule{
     var r,_=my.GetElementAt(i).(*drop_rule)
     return r
     }

func Newdrop_ruleList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*drop_ruleList{
      my := new(drop_ruleList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  Newdrop_ruleListFromElement(element *drop_rule,leftRecursive bool)*drop_ruleList{
        var obj = Newdrop_ruleList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *drop_ruleList)      AddElement(_drop_rule IAst){ 
      my.AbstractASTNodeList.AddElement(_drop_rule)
        _drop_rule.SetParent(my)
    }


func (my *drop_ruleList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *drop_ruleList)     Enter(v  Visitor){
        var checkChildren = v.Visitdrop_ruleList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.Getdrop_ruleAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisitdrop_ruleList(my)
    }


func AnyCastTodrop_ruleList(i interface{}) *drop_ruleList {
	if nil == i{
		return nil
	}else{
		return i.(*drop_ruleList)
	}
}
/**
 *<b>
*<li>Rule 79:  drop_rule ::= SYMBOL optMacroName produces ruleList
 *</b>
 */
type drop_rule struct{
    *ASTNode
      _SYMBOL *ASTNodeToken
      _optMacroName *optMacroName
      _produces Iproduces
      _ruleList *ruleList
}
func (my *drop_rule)      GetSYMBOL() *ASTNodeToken{ return my._SYMBOL}
func (my *drop_rule)      SetSYMBOL( _SYMBOL *ASTNodeToken)  { my._SYMBOL = _SYMBOL }
    /**
     * The value returned by <b>GetoptMacroName</b> may be <b>null</b>
     */
func (my *drop_rule)      GetoptMacroName() *optMacroName{ return my._optMacroName}
func (my *drop_rule)      SetoptMacroName( _optMacroName *optMacroName)  { my._optMacroName = _optMacroName }
func (my *drop_rule)      Getproduces() Iproduces{ return my._produces}
func (my *drop_rule)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *drop_rule)      GetruleList() *ruleList{ return my._ruleList}
func (my *drop_rule)      SetruleList( _ruleList *ruleList)  { my._ruleList = _ruleList }

func Newdrop_rule(leftIToken IToken, rightIToken IToken ,
              _SYMBOL *ASTNodeToken,
              _optMacroName *optMacroName,
              _produces Iproduces,
              _ruleList *ruleList)*drop_rule{
      my := new(drop_rule)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._SYMBOL = _SYMBOL;
        if nil != _SYMBOL{
        var trait_ interface{} = _SYMBOL
         trait_.(IAst).SetParent(my)
}
        my._optMacroName = _optMacroName;
        if nil != _optMacroName{
        var trait_ interface{} = _optMacroName
         trait_.(IAst).SetParent(my)
}
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._ruleList = _ruleList;
        if nil != _ruleList{
        var trait_ interface{} = _ruleList
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *drop_rule)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._SYMBOL{  list.Add(my._SYMBOL) }
        if nil != my._optMacroName{  list.Add(my._optMacroName) }
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._ruleList{  list.Add(my._ruleList) }
        return list
    }

func (my *drop_rule)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *drop_rule)       Enter(v Visitor){
        var checkChildren = v.Visitdrop_rule(my)
        if checkChildren{
            if nil != my._SYMBOL{my._SYMBOL.Accept(v)}
            if nil != my._optMacroName{my._optMacroName.Accept(v)}
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._ruleList{my._ruleList.Accept(v)}
        }
        v.EndVisitdrop_rule(my)
    }


func AnyCastTodrop_rule(i interface{}) *drop_rule {
	if nil == i{
		return nil
	}else{
		return i.(*drop_rule)
	}
}
/**
 *<em>
*<li>Rule 80:  optMacroName ::= $Empty
 *</em>
 *<p>
 *<b>
*<li>Rule 81:  optMacroName ::= MACRO_NAME
 *</b>
 */
type optMacroName struct{
    *ASTNodeToken
}
func (my *optMacroName)      GetMACRO_NAME()IToken{ return my.leftIToken; }

func NewoptMacroName(token IToken )*optMacroName{
      my := new(optMacroName)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *optMacroName)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *optMacroName)       Enter(v Visitor){
        v.VisitoptMacroName(my)
        v.EndVisitoptMacroName(my)
    }


func AnyCastTooptMacroName(i interface{}) *optMacroName {
	if nil == i{
		return nil
	}else{
		return i.(*optMacroName)
	}
}
/**
 *<b>
*<li>Rule 82:  include_segment ::= SYMBOL
 *</b>
 */
type include_segment struct{
    *ASTNodeToken
}
func (my *include_segment)      GetSYMBOL()IToken{ return my.leftIToken; }

func Newinclude_segment(token IToken )*include_segment{
      my := new(include_segment)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *include_segment)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *include_segment)       Enter(v Visitor){
        v.Visitinclude_segment(my)
        v.EndVisitinclude_segment(my)
    }


func AnyCastToinclude_segment(i interface{}) *include_segment {
	if nil == i{
		return nil
	}else{
		return i.(*include_segment)
	}
}
/**
 *<b>
*<li>Rule 83:  keywords_segment ::= keywordSpec
*<li>Rule 84:  keywords_segment ::= keywords_segment keywordSpec
 *</b>
 */
type keywordSpecList struct{
    *AbstractASTNodeList
}
func (my *keywordSpecList)      GetkeywordSpecAt(i int) IkeywordSpec{
     var r,_=my.GetElementAt(i).(IkeywordSpec)
     return r
     }

func NewkeywordSpecList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*keywordSpecList{
      my := new(keywordSpecList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  NewkeywordSpecListFromElement(element IkeywordSpec,leftRecursive bool)*keywordSpecList{
        var obj = NewkeywordSpecList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *keywordSpecList)      AddElement(_keywordSpec IAst){ 
      my.AbstractASTNodeList.AddElement(_keywordSpec)
        _keywordSpec.SetParent(my)
    }


func (my *keywordSpecList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *keywordSpecList)     Enter(v  Visitor){
        var checkChildren = v.VisitkeywordSpecList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.GetkeywordSpecAt(i)
                if nil !=element{                    element.Accept(v)
                }
            }
        }
        v.EndVisitkeywordSpecList(my)
    }


func AnyCastTokeywordSpecList(i interface{}) *keywordSpecList {
	if nil == i{
		return nil
	}else{
		return i.(*keywordSpecList)
	}
}
/**
 *<em>
*<li>Rule 85:  keywordSpec ::= terminal_symbol
 *</em>
 *<p>
 *<b>
*<li>Rule 86:  keywordSpec ::= terminal_symbol produces name
 *</b>
 */
type keywordSpec struct{
    *ASTNode
      _terminal_symbol Iterminal_symbol
      _produces Iproduces
      _name Iname
}
func (my *keywordSpec)      Getterminal_symbol() Iterminal_symbol{ return my._terminal_symbol}
func (my *keywordSpec)      Setterminal_symbol( _terminal_symbol Iterminal_symbol)  { my._terminal_symbol = _terminal_symbol }
func (my *keywordSpec)      Getproduces() Iproduces{ return my._produces}
func (my *keywordSpec)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *keywordSpec)      Getname() Iname{ return my._name}
func (my *keywordSpec)      Setname( _name Iname)  { my._name = _name }

func NewkeywordSpec(leftIToken IToken, rightIToken IToken ,
              _terminal_symbol Iterminal_symbol,
              _produces Iproduces,
              _name Iname)*keywordSpec{
      my := new(keywordSpec)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._terminal_symbol = _terminal_symbol;
        if nil != _terminal_symbol{
        var trait_ interface{} = _terminal_symbol
         trait_.(IAst).SetParent(my)
}
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._name = _name;
        if nil != _name{
        var trait_ interface{} = _name
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *keywordSpec)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._terminal_symbol{  list.Add(my._terminal_symbol) }
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._name{  list.Add(my._name) }
        return list
    }

func (my *keywordSpec)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *keywordSpec)       Enter(v Visitor){
        var checkChildren = v.VisitkeywordSpec(my)
        if checkChildren{
            if nil != my._terminal_symbol{my._terminal_symbol.Accept(v)}
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._name{my._name.Accept(v)}
        }
        v.EndVisitkeywordSpec(my)
    }


func AnyCastTokeywordSpec(i interface{}) *keywordSpec {
	if nil == i{
		return nil
	}else{
		return i.(*keywordSpec)
	}
}
/**
 *<b>
*<li>Rule 87:  names_segment ::= nameSpec
*<li>Rule 88:  names_segment ::= names_segment nameSpec
 *</b>
 */
type nameSpecList struct{
    *AbstractASTNodeList
}
func (my *nameSpecList)      GetnameSpecAt(i int) *nameSpec{
     var r,_=my.GetElementAt(i).(*nameSpec)
     return r
     }

func NewnameSpecList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*nameSpecList{
      my := new(nameSpecList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  NewnameSpecListFromElement(element *nameSpec,leftRecursive bool)*nameSpecList{
        var obj = NewnameSpecList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *nameSpecList)      AddElement(_nameSpec IAst){ 
      my.AbstractASTNodeList.AddElement(_nameSpec)
        _nameSpec.SetParent(my)
    }


func (my *nameSpecList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *nameSpecList)     Enter(v  Visitor){
        var checkChildren = v.VisitnameSpecList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.GetnameSpecAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisitnameSpecList(my)
    }


func AnyCastTonameSpecList(i interface{}) *nameSpecList {
	if nil == i{
		return nil
	}else{
		return i.(*nameSpecList)
	}
}
/**
 *<b>
*<li>Rule 89:  nameSpec ::= name produces name
 *</b>
 */
type nameSpec struct{
    *ASTNode
      _name Iname
      _produces Iproduces
      _name3 Iname
}
func (my *nameSpec)      Getname() Iname{ return my._name}
func (my *nameSpec)      Setname( _name Iname)  { my._name = _name }
func (my *nameSpec)      Getproduces() Iproduces{ return my._produces}
func (my *nameSpec)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *nameSpec)      Getname3() Iname{ return my._name3}
func (my *nameSpec)      Setname3( _name3 Iname)  { my._name3 = _name3 }

func NewnameSpec(leftIToken IToken, rightIToken IToken ,
              _name Iname,
              _produces Iproduces,
              _name3 Iname)*nameSpec{
      my := new(nameSpec)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._name = _name;
        if nil != _name{
        var trait_ interface{} = _name
         trait_.(IAst).SetParent(my)
}
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._name3 = _name3;
        if nil != _name3{
        var trait_ interface{} = _name3
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *nameSpec)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._name{  list.Add(my._name) }
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._name3{  list.Add(my._name3) }
        return list
    }

func (my *nameSpec)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *nameSpec)       Enter(v Visitor){
        var checkChildren = v.VisitnameSpec(my)
        if checkChildren{
            if nil != my._name{my._name.Accept(v)}
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._name3{my._name3.Accept(v)}
        }
        v.EndVisitnameSpec(my)
    }


func AnyCastTonameSpec(i interface{}) *nameSpec {
	if nil == i{
		return nil
	}else{
		return i.(*nameSpec)
	}
}
/**
 *<b>
*<li>Rule 98:  rules_segment ::= action_segment_list nonTermList
 *</b>
 */
type rules_segment struct{
    *ASTNode
      _action_segment_list *action_segmentList
      _nonTermList *nonTermList
}
func (my *rules_segment)      Getaction_segment_list() *action_segmentList{ return my._action_segment_list}
func (my *rules_segment)      Setaction_segment_list( _action_segment_list *action_segmentList)  { my._action_segment_list = _action_segment_list }
func (my *rules_segment)      GetnonTermList() *nonTermList{ return my._nonTermList}
func (my *rules_segment)      SetnonTermList( _nonTermList *nonTermList)  { my._nonTermList = _nonTermList }

func Newrules_segment(leftIToken IToken, rightIToken IToken ,
              _action_segment_list *action_segmentList,
              _nonTermList *nonTermList)*rules_segment{
      my := new(rules_segment)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._action_segment_list = _action_segment_list;
        if nil != _action_segment_list{
        var trait_ interface{} = _action_segment_list
         trait_.(IAst).SetParent(my)
}
        my._nonTermList = _nonTermList;
        if nil != _nonTermList{
        var trait_ interface{} = _nonTermList
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *rules_segment)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._action_segment_list{  list.Add(my._action_segment_list) }
        if nil != my._nonTermList{  list.Add(my._nonTermList) }
        return list
    }

func (my *rules_segment)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *rules_segment)       Enter(v Visitor){
        var checkChildren = v.Visitrules_segment(my)
        if checkChildren{
            if nil != my._action_segment_list{my._action_segment_list.Accept(v)}
            if nil != my._nonTermList{my._nonTermList.Accept(v)}
        }
        v.EndVisitrules_segment(my)
    }


func AnyCastTorules_segment(i interface{}) *rules_segment {
	if nil == i{
		return nil
	}else{
		return i.(*rules_segment)
	}
}
/**
 *<b>
*<li>Rule 99:  nonTermList ::= $Empty
*<li>Rule 100:  nonTermList ::= nonTermList nonTerm
 *</b>
 */
type nonTermList struct{
    *AbstractASTNodeList
}
func (my *nonTermList)      GetnonTermAt(i int) *nonTerm{
     var r,_=my.GetElementAt(i).(*nonTerm)
     return r
     }

func NewnonTermList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*nonTermList{
      my := new(nonTermList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  NewnonTermListFromElement(element *nonTerm,leftRecursive bool)*nonTermList{
        var obj = NewnonTermList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *nonTermList)      AddElement(_nonTerm IAst){ 
      my.AbstractASTNodeList.AddElement(_nonTerm)
        _nonTerm.SetParent(my)
    }


func (my *nonTermList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *nonTermList)     Enter(v  Visitor){
        var checkChildren = v.VisitnonTermList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.GetnonTermAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisitnonTermList(my)
    }


func AnyCastTononTermList(i interface{}) *nonTermList {
	if nil == i{
		return nil
	}else{
		return i.(*nonTermList)
	}
}
/**
 *<b>
*<li>Rule 101:  nonTerm ::= ruleNameWithAttributes produces ruleList
 *</b>
 */
type nonTerm struct{
    *ASTNode
      _ruleNameWithAttributes *RuleName
      _produces Iproduces
      _ruleList *ruleList
}
func (my *nonTerm)      GetruleNameWithAttributes() *RuleName{ return my._ruleNameWithAttributes}
func (my *nonTerm)      SetruleNameWithAttributes( _ruleNameWithAttributes *RuleName)  { my._ruleNameWithAttributes = _ruleNameWithAttributes }
func (my *nonTerm)      Getproduces() Iproduces{ return my._produces}
func (my *nonTerm)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *nonTerm)      GetruleList() *ruleList{ return my._ruleList}
func (my *nonTerm)      SetruleList( _ruleList *ruleList)  { my._ruleList = _ruleList }

func NewnonTerm(leftIToken IToken, rightIToken IToken ,
              _ruleNameWithAttributes *RuleName,
              _produces Iproduces,
              _ruleList *ruleList)*nonTerm{
      my := new(nonTerm)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._ruleNameWithAttributes = _ruleNameWithAttributes;
        if nil != _ruleNameWithAttributes{
        var trait_ interface{} = _ruleNameWithAttributes
         trait_.(IAst).SetParent(my)
}
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._ruleList = _ruleList;
        if nil != _ruleList{
        var trait_ interface{} = _ruleList
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *nonTerm)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._ruleNameWithAttributes{  list.Add(my._ruleNameWithAttributes) }
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._ruleList{  list.Add(my._ruleList) }
        return list
    }

func (my *nonTerm)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *nonTerm)       Enter(v Visitor){
        var checkChildren = v.VisitnonTerm(my)
        if checkChildren{
            if nil != my._ruleNameWithAttributes{my._ruleNameWithAttributes.Accept(v)}
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._ruleList{my._ruleList.Accept(v)}
        }
        v.EndVisitnonTerm(my)
    }


func AnyCastTononTerm(i interface{}) *nonTerm {
	if nil == i{
		return nil
	}else{
		return i.(*nonTerm)
	}
}
/**
 *<b>
*<li>Rule 102:  ruleNameWithAttributes ::= SYMBOL
*<li>Rule 103:  ruleNameWithAttributes ::= SYMBOL MACRO_NAME$className
*<li>Rule 104:  ruleNameWithAttributes ::= SYMBOL MACRO_NAME$className MACRO_NAME$arrayElement
 *</b>
 */
type RuleName struct{
    *ASTNode
      _SYMBOL *ASTNodeToken
      _className *ASTNodeToken
      _arrayElement *ASTNodeToken
}
func (my *RuleName)      GetSYMBOL() *ASTNodeToken{ return my._SYMBOL}
    /**
     * The value returned by <b>GetclassName</b> may be <b>null</b>
     */
func (my *RuleName)      GetclassName() *ASTNodeToken{ return my._className}
    /**
     * The value returned by <b>GetarrayElement</b> may be <b>null</b>
     */
func (my *RuleName)      GetarrayElement() *ASTNodeToken{ return my._arrayElement}

func NewRuleName(leftIToken IToken, rightIToken IToken ,
              _SYMBOL *ASTNodeToken,
              _className *ASTNodeToken ,
              _arrayElement *ASTNodeToken )*RuleName{
      my := new(RuleName)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._SYMBOL = _SYMBOL
        if nil != _SYMBOL{CastToAnyForLPGParser(_SYMBOL).(IAst).SetParent(my) }
        my._className = _className
        if nil != _className{CastToAnyForLPGParser(_className).(IAst).SetParent(my) }
        my._arrayElement = _arrayElement
        if nil != _arrayElement{CastToAnyForLPGParser(_arrayElement).(IAst).SetParent(my) }
        my.Initialize()
        return my
      }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *RuleName)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._SYMBOL{  list.Add(my._SYMBOL) }
        if nil != my._className{  list.Add(my._className) }
        if nil != my._arrayElement{  list.Add(my._arrayElement) }
        return list
    }

func (my *RuleName)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *RuleName)       Enter(v Visitor){
        var checkChildren = v.VisitRuleName(my)
        if checkChildren{
            if nil != my._SYMBOL{my._SYMBOL.Accept(v)}
            if nil != my._className{my._className.Accept(v)}
            if nil != my._arrayElement{my._arrayElement.Accept(v)}
        }
        v.EndVisitRuleName(my)
    }


func AnyCastToRuleName(i interface{}) *RuleName {
	if nil == i{
		return nil
	}else{
		return i.(*RuleName)
	}
}
/**
 *<b>
*<li>Rule 105:  ruleList ::= rule
*<li>Rule 106:  ruleList ::= ruleList |$ rule
 *</b>
 */
type ruleList struct{
    *AbstractASTNodeList
}
func (my *ruleList)      GetruleAt(i int) *rule{
     var r,_=my.GetElementAt(i).(*rule)
     return r
     }

func NewruleList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*ruleList{
      my := new(ruleList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  NewruleListFromElement(element *rule,leftRecursive bool)*ruleList{
        var obj = NewruleList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *ruleList)      AddElement(_rule IAst){ 
      my.AbstractASTNodeList.AddElement(_rule)
        _rule.SetParent(my)
    }


func (my *ruleList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *ruleList)     Enter(v  Visitor){
        var checkChildren = v.VisitruleList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.GetruleAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisitruleList(my)
    }


func AnyCastToruleList(i interface{}) *ruleList {
	if nil == i{
		return nil
	}else{
		return i.(*ruleList)
	}
}
/**
 *<b>
*<li>Rule 111:  rule ::= symWithAttrsList opt_action_segment
 *</b>
 */
type rule struct{
    *ASTNode
      _symWithAttrsList *symWithAttrsList
      _opt_action_segment *action_segment
}
func (my *rule)      GetsymWithAttrsList() *symWithAttrsList{ return my._symWithAttrsList}
func (my *rule)      SetsymWithAttrsList( _symWithAttrsList *symWithAttrsList)  { my._symWithAttrsList = _symWithAttrsList }
    /**
     * The value returned by <b>Getopt_action_segment</b> may be <b>null</b>
     */
func (my *rule)      Getopt_action_segment() *action_segment{ return my._opt_action_segment}
func (my *rule)      Setopt_action_segment( _opt_action_segment *action_segment)  { my._opt_action_segment = _opt_action_segment }

func Newrule(leftIToken IToken, rightIToken IToken ,
              _symWithAttrsList *symWithAttrsList,
              _opt_action_segment *action_segment)*rule{
      my := new(rule)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._symWithAttrsList = _symWithAttrsList;
        if nil != _symWithAttrsList{
        var trait_ interface{} = _symWithAttrsList
         trait_.(IAst).SetParent(my)
}
        my._opt_action_segment = _opt_action_segment;
        if nil != _opt_action_segment{
        var trait_ interface{} = _opt_action_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *rule)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._symWithAttrsList{  list.Add(my._symWithAttrsList) }
        if nil != my._opt_action_segment{  list.Add(my._opt_action_segment) }
        return list
    }

func (my *rule)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *rule)       Enter(v Visitor){
        var checkChildren = v.Visitrule(my)
        if checkChildren{
            if nil != my._symWithAttrsList{my._symWithAttrsList.Accept(v)}
            if nil != my._opt_action_segment{my._opt_action_segment.Accept(v)}
        }
        v.EndVisitrule(my)
    }


func AnyCastTorule(i interface{}) *rule {
	if nil == i{
		return nil
	}else{
		return i.(*rule)
	}
}
/**
 *<b>
*<li>Rule 112:  symWithAttrsList ::= $Empty
*<li>Rule 113:  symWithAttrsList ::= symWithAttrsList symWithAttrs
 *</b>
 */
type symWithAttrsList struct{
    *AbstractASTNodeList
}
func (my *symWithAttrsList)      GetsymWithAttrsAt(i int) IsymWithAttrs{
     var r,_=my.GetElementAt(i).(IsymWithAttrs)
     return r
     }

func NewsymWithAttrsList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*symWithAttrsList{
      my := new(symWithAttrsList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  NewsymWithAttrsListFromElement(element IsymWithAttrs,leftRecursive bool)*symWithAttrsList{
        var obj = NewsymWithAttrsList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *symWithAttrsList)      AddElement(_symWithAttrs IAst){ 
      my.AbstractASTNodeList.AddElement(_symWithAttrs)
        _symWithAttrs.SetParent(my)
    }


func (my *symWithAttrsList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *symWithAttrsList)     Enter(v  Visitor){
        var checkChildren = v.VisitsymWithAttrsList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.GetsymWithAttrsAt(i)
                if nil !=element{                    element.Accept(v)
                }
            }
        }
        v.EndVisitsymWithAttrsList(my)
    }


func AnyCastTosymWithAttrsList(i interface{}) *symWithAttrsList {
	if nil == i{
		return nil
	}else{
		return i.(*symWithAttrsList)
	}
}
/**
 *<b>
*<li>Rule 116:  optAttrList ::= $Empty
*<li>Rule 117:  optAttrList ::= MACRO_NAME
 *</b>
 */
type symAttrs struct{
    *ASTNode
      _MACRO_NAME *ASTNodeToken
}
    /**
     * The value returned by <b>GetMACRO_NAME</b> may be <b>null</b>
     */
func (my *symAttrs)      GetMACRO_NAME() *ASTNodeToken{ return my._MACRO_NAME}

func NewsymAttrs(leftIToken IToken, rightIToken IToken ,
              _MACRO_NAME *ASTNodeToken )*symAttrs{
      my := new(symAttrs)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._MACRO_NAME = _MACRO_NAME
        if nil != _MACRO_NAME{CastToAnyForLPGParser(_MACRO_NAME).(IAst).SetParent(my) }
        my.Initialize()
        return my
      }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *symAttrs)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._MACRO_NAME{  list.Add(my._MACRO_NAME) }
        return list
    }

func (my *symAttrs)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *symAttrs)       Enter(v Visitor){
        var checkChildren = v.VisitsymAttrs(my)
        if checkChildren{
            if nil != my._MACRO_NAME{my._MACRO_NAME.Accept(v)}
        }
        v.EndVisitsymAttrs(my)
    }


func AnyCastTosymAttrs(i interface{}) *symAttrs {
	if nil == i{
		return nil
	}else{
		return i.(*symAttrs)
	}
}
/**
 *<b>
*<li>Rule 120:  action_segment ::= BLOCK
 *</b>
 */
type action_segment struct{
    *ASTNodeToken
     environment *LPGParser
}
func (my *action_segment)     GetEnvironment() *LPGParser{ return my.environment }

func (my *action_segment)      GetBLOCK()IToken{ return my.leftIToken; }

func Newaction_segment(environment *LPGParser,token IToken )*action_segment{
      my := new(action_segment)
      my.environment = environment;
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *action_segment)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *action_segment)       Enter(v Visitor){
        v.Visitaction_segment(my)
        v.EndVisitaction_segment(my)
    }


func AnyCastToaction_segment(i interface{}) *action_segment {
	if nil == i{
		return nil
	}else{
		return i.(*action_segment)
	}
}
/**
 *<b>
*<li>Rule 121:  start_segment ::= start_symbol
*<li>Rule 122:  start_segment ::= start_segment start_symbol
 *</b>
 */
type start_symbolList struct{
    *AbstractASTNodeList
}
func (my *start_symbolList)      Getstart_symbolAt(i int) Istart_symbol{
     var r,_=my.GetElementAt(i).(Istart_symbol)
     return r
     }

func Newstart_symbolList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*start_symbolList{
      my := new(start_symbolList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  Newstart_symbolListFromElement(element Istart_symbol,leftRecursive bool)*start_symbolList{
        var obj = Newstart_symbolList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *start_symbolList)      AddElement(_start_symbol IAst){ 
      my.AbstractASTNodeList.AddElement(_start_symbol)
        _start_symbol.SetParent(my)
    }


func (my *start_symbolList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *start_symbolList)     Enter(v  Visitor){
        var checkChildren = v.Visitstart_symbolList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.Getstart_symbolAt(i)
                if nil !=element{                    element.Accept(v)
                }
            }
        }
        v.EndVisitstart_symbolList(my)
    }


func AnyCastTostart_symbolList(i interface{}) *start_symbolList {
	if nil == i{
		return nil
	}else{
		return i.(*start_symbolList)
	}
}
/**
 *<b>
*<li>Rule 125:  terminals_segment ::= terminal
*<li>Rule 126:  terminals_segment ::= terminals_segment terminal
 *</b>
 */
type terminalList struct{
    *AbstractASTNodeList
}
func (my *terminalList)      GetterminalAt(i int) *terminal{
     var r,_=my.GetElementAt(i).(*terminal)
     return r
     }

func NewterminalList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*terminalList{
      my := new(terminalList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  NewterminalListFromElement(element *terminal,leftRecursive bool)*terminalList{
        var obj = NewterminalList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *terminalList)      AddElement(_terminal IAst){ 
      my.AbstractASTNodeList.AddElement(_terminal)
        _terminal.SetParent(my)
    }


func (my *terminalList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *terminalList)     Enter(v  Visitor){
        var checkChildren = v.VisitterminalList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.GetterminalAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisitterminalList(my)
    }


func AnyCastToterminalList(i interface{}) *terminalList {
	if nil == i{
		return nil
	}else{
		return i.(*terminalList)
	}
}
/**
 *<b>
*<li>Rule 127:  terminal ::= terminal_symbol optTerminalAlias
 *</b>
 */
type terminal struct{
    *ASTNode
      _terminal_symbol Iterminal_symbol
      _optTerminalAlias *optTerminalAlias
}
func (my *terminal)      Getterminal_symbol() Iterminal_symbol{ return my._terminal_symbol}
func (my *terminal)      Setterminal_symbol( _terminal_symbol Iterminal_symbol)  { my._terminal_symbol = _terminal_symbol }
    /**
     * The value returned by <b>GetoptTerminalAlias</b> may be <b>null</b>
     */
func (my *terminal)      GetoptTerminalAlias() *optTerminalAlias{ return my._optTerminalAlias}
func (my *terminal)      SetoptTerminalAlias( _optTerminalAlias *optTerminalAlias)  { my._optTerminalAlias = _optTerminalAlias }

func Newterminal(leftIToken IToken, rightIToken IToken ,
              _terminal_symbol Iterminal_symbol,
              _optTerminalAlias *optTerminalAlias)*terminal{
      my := new(terminal)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._terminal_symbol = _terminal_symbol;
        if nil != _terminal_symbol{
        var trait_ interface{} = _terminal_symbol
         trait_.(IAst).SetParent(my)
}
        my._optTerminalAlias = _optTerminalAlias;
        if nil != _optTerminalAlias{
        var trait_ interface{} = _optTerminalAlias
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *terminal)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._terminal_symbol{  list.Add(my._terminal_symbol) }
        if nil != my._optTerminalAlias{  list.Add(my._optTerminalAlias) }
        return list
    }

func (my *terminal)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *terminal)       Enter(v Visitor){
        var checkChildren = v.Visitterminal(my)
        if checkChildren{
            if nil != my._terminal_symbol{my._terminal_symbol.Accept(v)}
            if nil != my._optTerminalAlias{my._optTerminalAlias.Accept(v)}
        }
        v.EndVisitterminal(my)
    }


func AnyCastToterminal(i interface{}) *terminal {
	if nil == i{
		return nil
	}else{
		return i.(*terminal)
	}
}
/**
 *<em>
*<li>Rule 128:  optTerminalAlias ::= $Empty
 *</em>
 *<p>
 *<b>
*<li>Rule 129:  optTerminalAlias ::= produces name
 *</b>
 */
type optTerminalAlias struct{
    *ASTNode
      _produces Iproduces
      _name Iname
}
func (my *optTerminalAlias)      Getproduces() Iproduces{ return my._produces}
func (my *optTerminalAlias)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *optTerminalAlias)      Getname() Iname{ return my._name}
func (my *optTerminalAlias)      Setname( _name Iname)  { my._name = _name }

func NewoptTerminalAlias(leftIToken IToken, rightIToken IToken ,
              _produces Iproduces,
              _name Iname)*optTerminalAlias{
      my := new(optTerminalAlias)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._name = _name;
        if nil != _name{
        var trait_ interface{} = _name
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *optTerminalAlias)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._name{  list.Add(my._name) }
        return list
    }

func (my *optTerminalAlias)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *optTerminalAlias)       Enter(v Visitor){
        var checkChildren = v.VisitoptTerminalAlias(my)
        if checkChildren{
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._name{my._name.Accept(v)}
        }
        v.EndVisitoptTerminalAlias(my)
    }


func AnyCastTooptTerminalAlias(i interface{}) *optTerminalAlias {
	if nil == i{
		return nil
	}else{
		return i.(*optTerminalAlias)
	}
}
/**
 *<b>
*<li>Rule 133:  types_segment ::= type_declarations
*<li>Rule 134:  types_segment ::= types_segment type_declarations
 *</b>
 */
type type_declarationsList struct{
    *AbstractASTNodeList
}
func (my *type_declarationsList)      Gettype_declarationsAt(i int) *type_declarations{
     var r,_=my.GetElementAt(i).(*type_declarations)
     return r
     }

func Newtype_declarationsList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*type_declarationsList{
      my := new(type_declarationsList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  Newtype_declarationsListFromElement(element *type_declarations,leftRecursive bool)*type_declarationsList{
        var obj = Newtype_declarationsList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *type_declarationsList)      AddElement(_type_declarations IAst){ 
      my.AbstractASTNodeList.AddElement(_type_declarations)
        _type_declarations.SetParent(my)
    }


func (my *type_declarationsList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *type_declarationsList)     Enter(v  Visitor){
        var checkChildren = v.Visittype_declarationsList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.Gettype_declarationsAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisittype_declarationsList(my)
    }


func AnyCastTotype_declarationsList(i interface{}) *type_declarationsList {
	if nil == i{
		return nil
	}else{
		return i.(*type_declarationsList)
	}
}
/**
 *<b>
*<li>Rule 135:  type_declarations ::= SYMBOL produces barSymbolList opt_action_segment
 *</b>
 */
type type_declarations struct{
    *ASTNode
      _SYMBOL *ASTNodeToken
      _produces Iproduces
      _barSymbolList *SYMBOLList
      _opt_action_segment *action_segment
}
func (my *type_declarations)      GetSYMBOL() *ASTNodeToken{ return my._SYMBOL}
func (my *type_declarations)      SetSYMBOL( _SYMBOL *ASTNodeToken)  { my._SYMBOL = _SYMBOL }
func (my *type_declarations)      Getproduces() Iproduces{ return my._produces}
func (my *type_declarations)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *type_declarations)      GetbarSymbolList() *SYMBOLList{ return my._barSymbolList}
func (my *type_declarations)      SetbarSymbolList( _barSymbolList *SYMBOLList)  { my._barSymbolList = _barSymbolList }
    /**
     * The value returned by <b>Getopt_action_segment</b> may be <b>null</b>
     */
func (my *type_declarations)      Getopt_action_segment() *action_segment{ return my._opt_action_segment}
func (my *type_declarations)      Setopt_action_segment( _opt_action_segment *action_segment)  { my._opt_action_segment = _opt_action_segment }

func Newtype_declarations(leftIToken IToken, rightIToken IToken ,
              _SYMBOL *ASTNodeToken,
              _produces Iproduces,
              _barSymbolList *SYMBOLList,
              _opt_action_segment *action_segment)*type_declarations{
      my := new(type_declarations)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._SYMBOL = _SYMBOL;
        if nil != _SYMBOL{
        var trait_ interface{} = _SYMBOL
         trait_.(IAst).SetParent(my)
}
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._barSymbolList = _barSymbolList;
        if nil != _barSymbolList{
        var trait_ interface{} = _barSymbolList
         trait_.(IAst).SetParent(my)
}
        my._opt_action_segment = _opt_action_segment;
        if nil != _opt_action_segment{
        var trait_ interface{} = _opt_action_segment
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *type_declarations)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._SYMBOL{  list.Add(my._SYMBOL) }
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._barSymbolList{  list.Add(my._barSymbolList) }
        if nil != my._opt_action_segment{  list.Add(my._opt_action_segment) }
        return list
    }

func (my *type_declarations)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *type_declarations)       Enter(v Visitor){
        var checkChildren = v.Visittype_declarations(my)
        if checkChildren{
            if nil != my._SYMBOL{my._SYMBOL.Accept(v)}
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._barSymbolList{my._barSymbolList.Accept(v)}
            if nil != my._opt_action_segment{my._opt_action_segment.Accept(v)}
        }
        v.EndVisittype_declarations(my)
    }


func AnyCastTotype_declarations(i interface{}) *type_declarations {
	if nil == i{
		return nil
	}else{
		return i.(*type_declarations)
	}
}
/**
 *<b>
*<li>Rule 138:  predecessor_segment ::= $Empty
*<li>Rule 139:  predecessor_segment ::= predecessor_segment symbol_pair
 *</b>
 */
type symbol_pairList struct{
    *AbstractASTNodeList
}
func (my *symbol_pairList)      Getsymbol_pairAt(i int) *symbol_pair{
     var r,_=my.GetElementAt(i).(*symbol_pair)
     return r
     }

func Newsymbol_pairList(leftToken  IToken, rightToken  IToken , leftRecursive bool)*symbol_pairList{
      my := new(symbol_pairList)
      my.AbstractASTNodeList = NewAbstractASTNodeList(leftToken, rightToken, leftRecursive)
      return my
}

    func  Newsymbol_pairListFromElement(element *symbol_pair,leftRecursive bool)*symbol_pairList{
        var obj = Newsymbol_pairList(element.GetLeftIToken(),element.GetRightIToken(), leftRecursive)
        obj.list.Add(element)
        CastToAnyForLPGParser(element).(IAst).SetParent(obj)
        return obj
    }

func (my *symbol_pairList)      AddElement(_symbol_pair IAst){ 
      my.AbstractASTNodeList.AddElement(_symbol_pair)
        _symbol_pair.SetParent(my)
    }


func (my *symbol_pairList)      Accept(v IAstVisitor ){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
    }
func (my *symbol_pairList)     Enter(v  Visitor){
        var checkChildren = v.Visitsymbol_pairList(my)
        if checkChildren{
           var i = 0
           for ; i < my.Size(); i++{
                 var element = my.Getsymbol_pairAt(i)
                if nil !=element{                    if ! v.PreVisit(element){ continue}
                    element.Enter(v)
                    v.PostVisit(element)
                }
            }
        }
        v.EndVisitsymbol_pairList(my)
    }


func AnyCastTosymbol_pairList(i interface{}) *symbol_pairList {
	if nil == i{
		return nil
	}else{
		return i.(*symbol_pairList)
	}
}
/**
 *<b>
*<li>Rule 140:  symbol_pair ::= SYMBOL SYMBOL
 *</b>
 */
type symbol_pair struct{
    *ASTNode
      _SYMBOL *ASTNodeToken
      _SYMBOL2 *ASTNodeToken
}
func (my *symbol_pair)      GetSYMBOL() *ASTNodeToken{ return my._SYMBOL}
func (my *symbol_pair)      SetSYMBOL( _SYMBOL *ASTNodeToken)  { my._SYMBOL = _SYMBOL }
func (my *symbol_pair)      GetSYMBOL2() *ASTNodeToken{ return my._SYMBOL2}
func (my *symbol_pair)      SetSYMBOL2( _SYMBOL2 *ASTNodeToken)  { my._SYMBOL2 = _SYMBOL2 }

func Newsymbol_pair(leftIToken IToken, rightIToken IToken ,
              _SYMBOL *ASTNodeToken,
              _SYMBOL2 *ASTNodeToken)*symbol_pair{
      my := new(symbol_pair)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._SYMBOL = _SYMBOL;
        if nil != _SYMBOL{
        var trait_ interface{} = _SYMBOL
         trait_.(IAst).SetParent(my)
}
        my._SYMBOL2 = _SYMBOL2;
        if nil != _SYMBOL2{
        var trait_ interface{} = _SYMBOL2
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *symbol_pair)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._SYMBOL{  list.Add(my._SYMBOL) }
        if nil != my._SYMBOL2{  list.Add(my._SYMBOL2) }
        return list
    }

func (my *symbol_pair)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *symbol_pair)       Enter(v Visitor){
        var checkChildren = v.Visitsymbol_pair(my)
        if checkChildren{
            if nil != my._SYMBOL{my._SYMBOL.Accept(v)}
            if nil != my._SYMBOL2{my._SYMBOL2.Accept(v)}
        }
        v.EndVisitsymbol_pair(my)
    }


func AnyCastTosymbol_pair(i interface{}) *symbol_pair {
	if nil == i{
		return nil
	}else{
		return i.(*symbol_pair)
	}
}
/**
 *<b>
*<li>Rule 143:  recover_symbol ::= SYMBOL
 *</b>
 */
type recover_symbol struct{
    *ASTNodeToken
}
func (my *recover_symbol)      GetSYMBOL()IToken{ return my.leftIToken; }

func Newrecover_symbol(token IToken )*recover_symbol{
      my := new(recover_symbol)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *recover_symbol)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *recover_symbol)       Enter(v Visitor){
        v.Visitrecover_symbol(my)
        v.EndVisitrecover_symbol(my)
    }


func AnyCastTorecover_symbol(i interface{}) *recover_symbol {
	if nil == i{
		return nil
	}else{
		return i.(*recover_symbol)
	}
}
/**
 *<em>
*<li>Rule 144:  END_KEY_OPT ::= $Empty
 *</em>
 *<p>
 *<b>
*<li>Rule 145:  END_KEY_OPT ::= END_KEY
 *</b>
 */
type END_KEY_OPT struct{
    *ASTNodeToken
}
func (my *END_KEY_OPT)      GetEND_KEY()IToken{ return my.leftIToken; }

func NewEND_KEY_OPT(token IToken )*END_KEY_OPT{
      my := new(END_KEY_OPT)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *END_KEY_OPT)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *END_KEY_OPT)       Enter(v Visitor){
        v.VisitEND_KEY_OPT(my)
        v.EndVisitEND_KEY_OPT(my)
    }


func AnyCastToEND_KEY_OPT(i interface{}) *END_KEY_OPT {
	if nil == i{
		return nil
	}else{
		return i.(*END_KEY_OPT)
	}
}
/**
 *<b>
*<li>Rule 34:  option_value ::= =$ SYMBOL
 *</b>
 */
type option_value0 struct{
    *ASTNode
      _SYMBOL *ASTNodeToken
}
func (my *option_value0)      GetSYMBOL() *ASTNodeToken{ return my._SYMBOL}
func (my *option_value0)      SetSYMBOL( _SYMBOL *ASTNodeToken)  { my._SYMBOL = _SYMBOL }

func Newoption_value0(leftIToken IToken, rightIToken IToken ,
              _SYMBOL *ASTNodeToken)*option_value0{
      my := new(option_value0)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._SYMBOL = _SYMBOL;
        if nil != _SYMBOL{
        var trait_ interface{} = _SYMBOL
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *option_value0)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._SYMBOL{  list.Add(my._SYMBOL) }
        return list
    }

func (my *option_value0)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *option_value0)       Enter(v Visitor){
        var checkChildren = v.Visitoption_value0(my)
        if checkChildren{
            if nil != my._SYMBOL{my._SYMBOL.Accept(v)}
        }
        v.EndVisitoption_value0(my)
    }


func AnyCastTooption_value0(i interface{}) *option_value0 {
	if nil == i{
		return nil
	}else{
		return i.(*option_value0)
	}
}
/**
 *<b>
*<li>Rule 35:  option_value ::= =$ ($ symbol_list )$
 *</b>
 */
type option_value1 struct{
    *ASTNode
      _symbol_list *SYMBOLList
}
func (my *option_value1)      Getsymbol_list() *SYMBOLList{ return my._symbol_list}
func (my *option_value1)      Setsymbol_list( _symbol_list *SYMBOLList)  { my._symbol_list = _symbol_list }

func Newoption_value1(leftIToken IToken, rightIToken IToken ,
              _symbol_list *SYMBOLList)*option_value1{
      my := new(option_value1)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._symbol_list = _symbol_list;
        if nil != _symbol_list{
        var trait_ interface{} = _symbol_list
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *option_value1)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._symbol_list{  list.Add(my._symbol_list) }
        return list
    }

func (my *option_value1)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *option_value1)       Enter(v Visitor){
        var checkChildren = v.Visitoption_value1(my)
        if checkChildren{
            if nil != my._symbol_list{my._symbol_list.Accept(v)}
        }
        v.EndVisitoption_value1(my)
    }


func AnyCastTooption_value1(i interface{}) *option_value1 {
	if nil == i{
		return nil
	}else{
		return i.(*option_value1)
	}
}
/**
 *<b>
*<li>Rule 40:  aliasSpec ::= ERROR_KEY produces alias_rhs
 *</b>
 */
type aliasSpec0 struct{
    *ASTNode
      _ERROR_KEY *ASTNodeToken
      _produces Iproduces
      _alias_rhs Ialias_rhs
}
func (my *aliasSpec0)      GetERROR_KEY() *ASTNodeToken{ return my._ERROR_KEY}
func (my *aliasSpec0)      SetERROR_KEY( _ERROR_KEY *ASTNodeToken)  { my._ERROR_KEY = _ERROR_KEY }
func (my *aliasSpec0)      Getproduces() Iproduces{ return my._produces}
func (my *aliasSpec0)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *aliasSpec0)      Getalias_rhs() Ialias_rhs{ return my._alias_rhs}
func (my *aliasSpec0)      Setalias_rhs( _alias_rhs Ialias_rhs)  { my._alias_rhs = _alias_rhs }

func NewaliasSpec0(leftIToken IToken, rightIToken IToken ,
              _ERROR_KEY *ASTNodeToken,
              _produces Iproduces,
              _alias_rhs Ialias_rhs)*aliasSpec0{
      my := new(aliasSpec0)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._ERROR_KEY = _ERROR_KEY;
        if nil != _ERROR_KEY{
        var trait_ interface{} = _ERROR_KEY
         trait_.(IAst).SetParent(my)
}
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._alias_rhs = _alias_rhs;
        if nil != _alias_rhs{
        var trait_ interface{} = _alias_rhs
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *aliasSpec0)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._ERROR_KEY{  list.Add(my._ERROR_KEY) }
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._alias_rhs{  list.Add(my._alias_rhs) }
        return list
    }

func (my *aliasSpec0)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *aliasSpec0)       Enter(v Visitor){
        var checkChildren = v.VisitaliasSpec0(my)
        if checkChildren{
            if nil != my._ERROR_KEY{my._ERROR_KEY.Accept(v)}
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._alias_rhs{my._alias_rhs.Accept(v)}
        }
        v.EndVisitaliasSpec0(my)
    }


func AnyCastToaliasSpec0(i interface{}) *aliasSpec0 {
	if nil == i{
		return nil
	}else{
		return i.(*aliasSpec0)
	}
}
/**
 *<b>
*<li>Rule 41:  aliasSpec ::= EOL_KEY produces alias_rhs
 *</b>
 */
type aliasSpec1 struct{
    *ASTNode
      _EOL_KEY *ASTNodeToken
      _produces Iproduces
      _alias_rhs Ialias_rhs
}
func (my *aliasSpec1)      GetEOL_KEY() *ASTNodeToken{ return my._EOL_KEY}
func (my *aliasSpec1)      SetEOL_KEY( _EOL_KEY *ASTNodeToken)  { my._EOL_KEY = _EOL_KEY }
func (my *aliasSpec1)      Getproduces() Iproduces{ return my._produces}
func (my *aliasSpec1)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *aliasSpec1)      Getalias_rhs() Ialias_rhs{ return my._alias_rhs}
func (my *aliasSpec1)      Setalias_rhs( _alias_rhs Ialias_rhs)  { my._alias_rhs = _alias_rhs }

func NewaliasSpec1(leftIToken IToken, rightIToken IToken ,
              _EOL_KEY *ASTNodeToken,
              _produces Iproduces,
              _alias_rhs Ialias_rhs)*aliasSpec1{
      my := new(aliasSpec1)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._EOL_KEY = _EOL_KEY;
        if nil != _EOL_KEY{
        var trait_ interface{} = _EOL_KEY
         trait_.(IAst).SetParent(my)
}
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._alias_rhs = _alias_rhs;
        if nil != _alias_rhs{
        var trait_ interface{} = _alias_rhs
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *aliasSpec1)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._EOL_KEY{  list.Add(my._EOL_KEY) }
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._alias_rhs{  list.Add(my._alias_rhs) }
        return list
    }

func (my *aliasSpec1)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *aliasSpec1)       Enter(v Visitor){
        var checkChildren = v.VisitaliasSpec1(my)
        if checkChildren{
            if nil != my._EOL_KEY{my._EOL_KEY.Accept(v)}
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._alias_rhs{my._alias_rhs.Accept(v)}
        }
        v.EndVisitaliasSpec1(my)
    }


func AnyCastToaliasSpec1(i interface{}) *aliasSpec1 {
	if nil == i{
		return nil
	}else{
		return i.(*aliasSpec1)
	}
}
/**
 *<b>
*<li>Rule 42:  aliasSpec ::= EOF_KEY produces alias_rhs
 *</b>
 */
type aliasSpec2 struct{
    *ASTNode
      _EOF_KEY *ASTNodeToken
      _produces Iproduces
      _alias_rhs Ialias_rhs
}
func (my *aliasSpec2)      GetEOF_KEY() *ASTNodeToken{ return my._EOF_KEY}
func (my *aliasSpec2)      SetEOF_KEY( _EOF_KEY *ASTNodeToken)  { my._EOF_KEY = _EOF_KEY }
func (my *aliasSpec2)      Getproduces() Iproduces{ return my._produces}
func (my *aliasSpec2)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *aliasSpec2)      Getalias_rhs() Ialias_rhs{ return my._alias_rhs}
func (my *aliasSpec2)      Setalias_rhs( _alias_rhs Ialias_rhs)  { my._alias_rhs = _alias_rhs }

func NewaliasSpec2(leftIToken IToken, rightIToken IToken ,
              _EOF_KEY *ASTNodeToken,
              _produces Iproduces,
              _alias_rhs Ialias_rhs)*aliasSpec2{
      my := new(aliasSpec2)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._EOF_KEY = _EOF_KEY;
        if nil != _EOF_KEY{
        var trait_ interface{} = _EOF_KEY
         trait_.(IAst).SetParent(my)
}
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._alias_rhs = _alias_rhs;
        if nil != _alias_rhs{
        var trait_ interface{} = _alias_rhs
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *aliasSpec2)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._EOF_KEY{  list.Add(my._EOF_KEY) }
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._alias_rhs{  list.Add(my._alias_rhs) }
        return list
    }

func (my *aliasSpec2)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *aliasSpec2)       Enter(v Visitor){
        var checkChildren = v.VisitaliasSpec2(my)
        if checkChildren{
            if nil != my._EOF_KEY{my._EOF_KEY.Accept(v)}
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._alias_rhs{my._alias_rhs.Accept(v)}
        }
        v.EndVisitaliasSpec2(my)
    }


func AnyCastToaliasSpec2(i interface{}) *aliasSpec2 {
	if nil == i{
		return nil
	}else{
		return i.(*aliasSpec2)
	}
}
/**
 *<b>
*<li>Rule 43:  aliasSpec ::= IDENTIFIER_KEY produces alias_rhs
 *</b>
 */
type aliasSpec3 struct{
    *ASTNode
      _IDENTIFIER_KEY *ASTNodeToken
      _produces Iproduces
      _alias_rhs Ialias_rhs
}
func (my *aliasSpec3)      GetIDENTIFIER_KEY() *ASTNodeToken{ return my._IDENTIFIER_KEY}
func (my *aliasSpec3)      SetIDENTIFIER_KEY( _IDENTIFIER_KEY *ASTNodeToken)  { my._IDENTIFIER_KEY = _IDENTIFIER_KEY }
func (my *aliasSpec3)      Getproduces() Iproduces{ return my._produces}
func (my *aliasSpec3)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *aliasSpec3)      Getalias_rhs() Ialias_rhs{ return my._alias_rhs}
func (my *aliasSpec3)      Setalias_rhs( _alias_rhs Ialias_rhs)  { my._alias_rhs = _alias_rhs }

func NewaliasSpec3(leftIToken IToken, rightIToken IToken ,
              _IDENTIFIER_KEY *ASTNodeToken,
              _produces Iproduces,
              _alias_rhs Ialias_rhs)*aliasSpec3{
      my := new(aliasSpec3)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._IDENTIFIER_KEY = _IDENTIFIER_KEY;
        if nil != _IDENTIFIER_KEY{
        var trait_ interface{} = _IDENTIFIER_KEY
         trait_.(IAst).SetParent(my)
}
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._alias_rhs = _alias_rhs;
        if nil != _alias_rhs{
        var trait_ interface{} = _alias_rhs
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *aliasSpec3)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._IDENTIFIER_KEY{  list.Add(my._IDENTIFIER_KEY) }
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._alias_rhs{  list.Add(my._alias_rhs) }
        return list
    }

func (my *aliasSpec3)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *aliasSpec3)       Enter(v Visitor){
        var checkChildren = v.VisitaliasSpec3(my)
        if checkChildren{
            if nil != my._IDENTIFIER_KEY{my._IDENTIFIER_KEY.Accept(v)}
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._alias_rhs{my._alias_rhs.Accept(v)}
        }
        v.EndVisitaliasSpec3(my)
    }


func AnyCastToaliasSpec3(i interface{}) *aliasSpec3 {
	if nil == i{
		return nil
	}else{
		return i.(*aliasSpec3)
	}
}
/**
 *<b>
*<li>Rule 44:  aliasSpec ::= SYMBOL produces alias_rhs
 *</b>
 */
type aliasSpec4 struct{
    *ASTNode
      _SYMBOL *ASTNodeToken
      _produces Iproduces
      _alias_rhs Ialias_rhs
}
func (my *aliasSpec4)      GetSYMBOL() *ASTNodeToken{ return my._SYMBOL}
func (my *aliasSpec4)      SetSYMBOL( _SYMBOL *ASTNodeToken)  { my._SYMBOL = _SYMBOL }
func (my *aliasSpec4)      Getproduces() Iproduces{ return my._produces}
func (my *aliasSpec4)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *aliasSpec4)      Getalias_rhs() Ialias_rhs{ return my._alias_rhs}
func (my *aliasSpec4)      Setalias_rhs( _alias_rhs Ialias_rhs)  { my._alias_rhs = _alias_rhs }

func NewaliasSpec4(leftIToken IToken, rightIToken IToken ,
              _SYMBOL *ASTNodeToken,
              _produces Iproduces,
              _alias_rhs Ialias_rhs)*aliasSpec4{
      my := new(aliasSpec4)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._SYMBOL = _SYMBOL;
        if nil != _SYMBOL{
        var trait_ interface{} = _SYMBOL
         trait_.(IAst).SetParent(my)
}
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._alias_rhs = _alias_rhs;
        if nil != _alias_rhs{
        var trait_ interface{} = _alias_rhs
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *aliasSpec4)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._SYMBOL{  list.Add(my._SYMBOL) }
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._alias_rhs{  list.Add(my._alias_rhs) }
        return list
    }

func (my *aliasSpec4)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *aliasSpec4)       Enter(v Visitor){
        var checkChildren = v.VisitaliasSpec4(my)
        if checkChildren{
            if nil != my._SYMBOL{my._SYMBOL.Accept(v)}
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._alias_rhs{my._alias_rhs.Accept(v)}
        }
        v.EndVisitaliasSpec4(my)
    }


func AnyCastToaliasSpec4(i interface{}) *aliasSpec4 {
	if nil == i{
		return nil
	}else{
		return i.(*aliasSpec4)
	}
}
/**
 *<b>
*<li>Rule 45:  aliasSpec ::= alias_lhs_macro_name produces alias_rhs
 *</b>
 */
type aliasSpec5 struct{
    *ASTNode
      _alias_lhs_macro_name *alias_lhs_macro_name
      _produces Iproduces
      _alias_rhs Ialias_rhs
}
func (my *aliasSpec5)      Getalias_lhs_macro_name() *alias_lhs_macro_name{ return my._alias_lhs_macro_name}
func (my *aliasSpec5)      Setalias_lhs_macro_name( _alias_lhs_macro_name *alias_lhs_macro_name)  { my._alias_lhs_macro_name = _alias_lhs_macro_name }
func (my *aliasSpec5)      Getproduces() Iproduces{ return my._produces}
func (my *aliasSpec5)      Setproduces( _produces Iproduces)  { my._produces = _produces }
func (my *aliasSpec5)      Getalias_rhs() Ialias_rhs{ return my._alias_rhs}
func (my *aliasSpec5)      Setalias_rhs( _alias_rhs Ialias_rhs)  { my._alias_rhs = _alias_rhs }

func NewaliasSpec5(leftIToken IToken, rightIToken IToken ,
              _alias_lhs_macro_name *alias_lhs_macro_name,
              _produces Iproduces,
              _alias_rhs Ialias_rhs)*aliasSpec5{
      my := new(aliasSpec5)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._alias_lhs_macro_name = _alias_lhs_macro_name;
        if nil != _alias_lhs_macro_name{
        var trait_ interface{} = _alias_lhs_macro_name
         trait_.(IAst).SetParent(my)
}
        my._produces = _produces;
        if nil != _produces{
        var trait_ interface{} = _produces
         trait_.(IAst).SetParent(my)
}
        my._alias_rhs = _alias_rhs;
        if nil != _alias_rhs{
        var trait_ interface{} = _alias_rhs
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *aliasSpec5)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._alias_lhs_macro_name{  list.Add(my._alias_lhs_macro_name) }
        if nil != my._produces{  list.Add(my._produces) }
        if nil != my._alias_rhs{  list.Add(my._alias_rhs) }
        return list
    }

func (my *aliasSpec5)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *aliasSpec5)       Enter(v Visitor){
        var checkChildren = v.VisitaliasSpec5(my)
        if checkChildren{
            if nil != my._alias_lhs_macro_name{my._alias_lhs_macro_name.Accept(v)}
            if nil != my._produces{my._produces.Accept(v)}
            if nil != my._alias_rhs{my._alias_rhs.Accept(v)}
        }
        v.EndVisitaliasSpec5(my)
    }


func AnyCastToaliasSpec5(i interface{}) *aliasSpec5 {
	if nil == i{
		return nil
	}else{
		return i.(*aliasSpec5)
	}
}
/**
 *<b>
*<li>Rule 47:  alias_rhs ::= SYMBOL
 *</b>
 */
type alias_rhs0 struct{
    *ASTNodeToken
}
func (my *alias_rhs0)      GetSYMBOL()IToken{ return my.leftIToken; }

func Newalias_rhs0(token IToken )*alias_rhs0{
      my := new(alias_rhs0)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *alias_rhs0)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *alias_rhs0)       Enter(v Visitor){
        v.Visitalias_rhs0(my)
        v.EndVisitalias_rhs0(my)
    }


func AnyCastToalias_rhs0(i interface{}) *alias_rhs0 {
	if nil == i{
		return nil
	}else{
		return i.(*alias_rhs0)
	}
}
/**
 *<b>
*<li>Rule 48:  alias_rhs ::= MACRO_NAME
 *</b>
 */
type alias_rhs1 struct{
    *ASTNodeToken
}
func (my *alias_rhs1)      GetMACRO_NAME()IToken{ return my.leftIToken; }

func Newalias_rhs1(token IToken )*alias_rhs1{
      my := new(alias_rhs1)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *alias_rhs1)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *alias_rhs1)       Enter(v Visitor){
        v.Visitalias_rhs1(my)
        v.EndVisitalias_rhs1(my)
    }


func AnyCastToalias_rhs1(i interface{}) *alias_rhs1 {
	if nil == i{
		return nil
	}else{
		return i.(*alias_rhs1)
	}
}
/**
 *<b>
*<li>Rule 49:  alias_rhs ::= ERROR_KEY
 *</b>
 */
type alias_rhs2 struct{
    *ASTNodeToken
}
func (my *alias_rhs2)      GetERROR_KEY()IToken{ return my.leftIToken; }

func Newalias_rhs2(token IToken )*alias_rhs2{
      my := new(alias_rhs2)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *alias_rhs2)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *alias_rhs2)       Enter(v Visitor){
        v.Visitalias_rhs2(my)
        v.EndVisitalias_rhs2(my)
    }


func AnyCastToalias_rhs2(i interface{}) *alias_rhs2 {
	if nil == i{
		return nil
	}else{
		return i.(*alias_rhs2)
	}
}
/**
 *<b>
*<li>Rule 50:  alias_rhs ::= EOL_KEY
 *</b>
 */
type alias_rhs3 struct{
    *ASTNodeToken
}
func (my *alias_rhs3)      GetEOL_KEY()IToken{ return my.leftIToken; }

func Newalias_rhs3(token IToken )*alias_rhs3{
      my := new(alias_rhs3)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *alias_rhs3)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *alias_rhs3)       Enter(v Visitor){
        v.Visitalias_rhs3(my)
        v.EndVisitalias_rhs3(my)
    }


func AnyCastToalias_rhs3(i interface{}) *alias_rhs3 {
	if nil == i{
		return nil
	}else{
		return i.(*alias_rhs3)
	}
}
/**
 *<b>
*<li>Rule 51:  alias_rhs ::= EOF_KEY
 *</b>
 */
type alias_rhs4 struct{
    *ASTNodeToken
}
func (my *alias_rhs4)      GetEOF_KEY()IToken{ return my.leftIToken; }

func Newalias_rhs4(token IToken )*alias_rhs4{
      my := new(alias_rhs4)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *alias_rhs4)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *alias_rhs4)       Enter(v Visitor){
        v.Visitalias_rhs4(my)
        v.EndVisitalias_rhs4(my)
    }


func AnyCastToalias_rhs4(i interface{}) *alias_rhs4 {
	if nil == i{
		return nil
	}else{
		return i.(*alias_rhs4)
	}
}
/**
 *<b>
*<li>Rule 52:  alias_rhs ::= EMPTY_KEY
 *</b>
 */
type alias_rhs5 struct{
    *ASTNodeToken
}
func (my *alias_rhs5)      GetEMPTY_KEY()IToken{ return my.leftIToken; }

func Newalias_rhs5(token IToken )*alias_rhs5{
      my := new(alias_rhs5)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *alias_rhs5)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *alias_rhs5)       Enter(v Visitor){
        v.Visitalias_rhs5(my)
        v.EndVisitalias_rhs5(my)
    }


func AnyCastToalias_rhs5(i interface{}) *alias_rhs5 {
	if nil == i{
		return nil
	}else{
		return i.(*alias_rhs5)
	}
}
/**
 *<b>
*<li>Rule 53:  alias_rhs ::= IDENTIFIER_KEY
 *</b>
 */
type alias_rhs6 struct{
    *ASTNodeToken
}
func (my *alias_rhs6)      GetIDENTIFIER_KEY()IToken{ return my.leftIToken; }

func Newalias_rhs6(token IToken )*alias_rhs6{
      my := new(alias_rhs6)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *alias_rhs6)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *alias_rhs6)       Enter(v Visitor){
        v.Visitalias_rhs6(my)
        v.EndVisitalias_rhs6(my)
    }


func AnyCastToalias_rhs6(i interface{}) *alias_rhs6 {
	if nil == i{
		return nil
	}else{
		return i.(*alias_rhs6)
	}
}
/**
 *<b>
*<li>Rule 58:  macro_name_symbol ::= MACRO_NAME
 *</b>
 */
type macro_name_symbol0 struct{
    *ASTNodeToken
}
func (my *macro_name_symbol0)      GetMACRO_NAME()IToken{ return my.leftIToken; }

func Newmacro_name_symbol0(token IToken )*macro_name_symbol0{
      my := new(macro_name_symbol0)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *macro_name_symbol0)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *macro_name_symbol0)       Enter(v Visitor){
        v.Visitmacro_name_symbol0(my)
        v.EndVisitmacro_name_symbol0(my)
    }


func AnyCastTomacro_name_symbol0(i interface{}) *macro_name_symbol0 {
	if nil == i{
		return nil
	}else{
		return i.(*macro_name_symbol0)
	}
}
/**
 *<b>
*<li>Rule 59:  macro_name_symbol ::= SYMBOL
 *</b>
 */
type macro_name_symbol1 struct{
    *ASTNodeToken
}
func (my *macro_name_symbol1)      GetSYMBOL()IToken{ return my.leftIToken; }

func Newmacro_name_symbol1(token IToken )*macro_name_symbol1{
      my := new(macro_name_symbol1)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *macro_name_symbol1)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *macro_name_symbol1)       Enter(v Visitor){
        v.Visitmacro_name_symbol1(my)
        v.EndVisitmacro_name_symbol1(my)
    }


func AnyCastTomacro_name_symbol1(i interface{}) *macro_name_symbol1 {
	if nil == i{
		return nil
	}else{
		return i.(*macro_name_symbol1)
	}
}
/**
 *<b>
*<li>Rule 73:  drop_command ::= DROPSYMBOLS_KEY drop_symbols
 *</b>
 */
type drop_command0 struct{
    *ASTNode
      _DROPSYMBOLS_KEY *ASTNodeToken
      _drop_symbols *SYMBOLList
}
func (my *drop_command0)      GetDROPSYMBOLS_KEY() *ASTNodeToken{ return my._DROPSYMBOLS_KEY}
func (my *drop_command0)      SetDROPSYMBOLS_KEY( _DROPSYMBOLS_KEY *ASTNodeToken)  { my._DROPSYMBOLS_KEY = _DROPSYMBOLS_KEY }
func (my *drop_command0)      Getdrop_symbols() *SYMBOLList{ return my._drop_symbols}
func (my *drop_command0)      Setdrop_symbols( _drop_symbols *SYMBOLList)  { my._drop_symbols = _drop_symbols }

func Newdrop_command0(leftIToken IToken, rightIToken IToken ,
              _DROPSYMBOLS_KEY *ASTNodeToken,
              _drop_symbols *SYMBOLList)*drop_command0{
      my := new(drop_command0)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._DROPSYMBOLS_KEY = _DROPSYMBOLS_KEY;
        if nil != _DROPSYMBOLS_KEY{
        var trait_ interface{} = _DROPSYMBOLS_KEY
         trait_.(IAst).SetParent(my)
}
        my._drop_symbols = _drop_symbols;
        if nil != _drop_symbols{
        var trait_ interface{} = _drop_symbols
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *drop_command0)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._DROPSYMBOLS_KEY{  list.Add(my._DROPSYMBOLS_KEY) }
        if nil != my._drop_symbols{  list.Add(my._drop_symbols) }
        return list
    }

func (my *drop_command0)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *drop_command0)       Enter(v Visitor){
        var checkChildren = v.Visitdrop_command0(my)
        if checkChildren{
            if nil != my._DROPSYMBOLS_KEY{my._DROPSYMBOLS_KEY.Accept(v)}
            if nil != my._drop_symbols{my._drop_symbols.Accept(v)}
        }
        v.EndVisitdrop_command0(my)
    }


func AnyCastTodrop_command0(i interface{}) *drop_command0 {
	if nil == i{
		return nil
	}else{
		return i.(*drop_command0)
	}
}
/**
 *<b>
*<li>Rule 74:  drop_command ::= DROPRULES_KEY drop_rules
 *</b>
 */
type drop_command1 struct{
    *ASTNode
      _DROPRULES_KEY *ASTNodeToken
      _drop_rules *drop_ruleList
}
func (my *drop_command1)      GetDROPRULES_KEY() *ASTNodeToken{ return my._DROPRULES_KEY}
func (my *drop_command1)      SetDROPRULES_KEY( _DROPRULES_KEY *ASTNodeToken)  { my._DROPRULES_KEY = _DROPRULES_KEY }
func (my *drop_command1)      Getdrop_rules() *drop_ruleList{ return my._drop_rules}
func (my *drop_command1)      Setdrop_rules( _drop_rules *drop_ruleList)  { my._drop_rules = _drop_rules }

func Newdrop_command1(leftIToken IToken, rightIToken IToken ,
              _DROPRULES_KEY *ASTNodeToken,
              _drop_rules *drop_ruleList)*drop_command1{
      my := new(drop_command1)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._DROPRULES_KEY = _DROPRULES_KEY;
        if nil != _DROPRULES_KEY{
        var trait_ interface{} = _DROPRULES_KEY
         trait_.(IAst).SetParent(my)
}
        my._drop_rules = _drop_rules;
        if nil != _drop_rules{
        var trait_ interface{} = _drop_rules
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *drop_command1)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._DROPRULES_KEY{  list.Add(my._DROPRULES_KEY) }
        if nil != my._drop_rules{  list.Add(my._drop_rules) }
        return list
    }

func (my *drop_command1)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *drop_command1)       Enter(v Visitor){
        var checkChildren = v.Visitdrop_command1(my)
        if checkChildren{
            if nil != my._DROPRULES_KEY{my._DROPRULES_KEY.Accept(v)}
            if nil != my._drop_rules{my._drop_rules.Accept(v)}
        }
        v.EndVisitdrop_command1(my)
    }


func AnyCastTodrop_command1(i interface{}) *drop_command1 {
	if nil == i{
		return nil
	}else{
		return i.(*drop_command1)
	}
}
/**
 *<b>
*<li>Rule 90:  name ::= SYMBOL
 *</b>
 */
type name0 struct{
    *ASTNodeToken
}
func (my *name0)      GetSYMBOL()IToken{ return my.leftIToken; }

func Newname0(token IToken )*name0{
      my := new(name0)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *name0)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *name0)       Enter(v Visitor){
        v.Visitname0(my)
        v.EndVisitname0(my)
    }


func AnyCastToname0(i interface{}) *name0 {
	if nil == i{
		return nil
	}else{
		return i.(*name0)
	}
}
/**
 *<b>
*<li>Rule 91:  name ::= MACRO_NAME
 *</b>
 */
type name1 struct{
    *ASTNodeToken
}
func (my *name1)      GetMACRO_NAME()IToken{ return my.leftIToken; }

func Newname1(token IToken )*name1{
      my := new(name1)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *name1)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *name1)       Enter(v Visitor){
        v.Visitname1(my)
        v.EndVisitname1(my)
    }


func AnyCastToname1(i interface{}) *name1 {
	if nil == i{
		return nil
	}else{
		return i.(*name1)
	}
}
/**
 *<b>
*<li>Rule 92:  name ::= EMPTY_KEY
 *</b>
 */
type name2 struct{
    *ASTNodeToken
}
func (my *name2)      GetEMPTY_KEY()IToken{ return my.leftIToken; }

func Newname2(token IToken )*name2{
      my := new(name2)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *name2)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *name2)       Enter(v Visitor){
        v.Visitname2(my)
        v.EndVisitname2(my)
    }


func AnyCastToname2(i interface{}) *name2 {
	if nil == i{
		return nil
	}else{
		return i.(*name2)
	}
}
/**
 *<b>
*<li>Rule 93:  name ::= ERROR_KEY
 *</b>
 */
type name3 struct{
    *ASTNodeToken
}
func (my *name3)      GetERROR_KEY()IToken{ return my.leftIToken; }

func Newname3(token IToken )*name3{
      my := new(name3)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *name3)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *name3)       Enter(v Visitor){
        v.Visitname3(my)
        v.EndVisitname3(my)
    }


func AnyCastToname3(i interface{}) *name3 {
	if nil == i{
		return nil
	}else{
		return i.(*name3)
	}
}
/**
 *<b>
*<li>Rule 94:  name ::= EOL_KEY
 *</b>
 */
type name4 struct{
    *ASTNodeToken
}
func (my *name4)      GetEOL_KEY()IToken{ return my.leftIToken; }

func Newname4(token IToken )*name4{
      my := new(name4)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *name4)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *name4)       Enter(v Visitor){
        v.Visitname4(my)
        v.EndVisitname4(my)
    }


func AnyCastToname4(i interface{}) *name4 {
	if nil == i{
		return nil
	}else{
		return i.(*name4)
	}
}
/**
 *<b>
*<li>Rule 95:  name ::= IDENTIFIER_KEY
 *</b>
 */
type name5 struct{
    *ASTNodeToken
}
func (my *name5)      GetIDENTIFIER_KEY()IToken{ return my.leftIToken; }

func Newname5(token IToken )*name5{
      my := new(name5)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *name5)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *name5)       Enter(v Visitor){
        v.Visitname5(my)
        v.EndVisitname5(my)
    }


func AnyCastToname5(i interface{}) *name5 {
	if nil == i{
		return nil
	}else{
		return i.(*name5)
	}
}
/**
 *<b>
*<li>Rule 107:  produces ::= ::=
 *</b>
 */
type produces0 struct{
    *ASTNodeToken
}
func (my *produces0)      GetEQUIVALENCE()IToken{ return my.leftIToken; }

func Newproduces0(token IToken )*produces0{
      my := new(produces0)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *produces0)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *produces0)       Enter(v Visitor){
        v.Visitproduces0(my)
        v.EndVisitproduces0(my)
    }


func AnyCastToproduces0(i interface{}) *produces0 {
	if nil == i{
		return nil
	}else{
		return i.(*produces0)
	}
}
/**
 *<b>
*<li>Rule 108:  produces ::= ::=?
 *</b>
 */
type produces1 struct{
    *ASTNodeToken
}
func (my *produces1)      GetPRIORITY_EQUIVALENCE()IToken{ return my.leftIToken; }

func Newproduces1(token IToken )*produces1{
      my := new(produces1)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *produces1)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *produces1)       Enter(v Visitor){
        v.Visitproduces1(my)
        v.EndVisitproduces1(my)
    }


func AnyCastToproduces1(i interface{}) *produces1 {
	if nil == i{
		return nil
	}else{
		return i.(*produces1)
	}
}
/**
 *<b>
*<li>Rule 109:  produces ::= ->
 *</b>
 */
type produces2 struct{
    *ASTNodeToken
}
func (my *produces2)      GetARROW()IToken{ return my.leftIToken; }

func Newproduces2(token IToken )*produces2{
      my := new(produces2)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *produces2)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *produces2)       Enter(v Visitor){
        v.Visitproduces2(my)
        v.EndVisitproduces2(my)
    }


func AnyCastToproduces2(i interface{}) *produces2 {
	if nil == i{
		return nil
	}else{
		return i.(*produces2)
	}
}
/**
 *<b>
*<li>Rule 110:  produces ::= ->?
 *</b>
 */
type produces3 struct{
    *ASTNodeToken
}
func (my *produces3)      GetPRIORITY_ARROW()IToken{ return my.leftIToken; }

func Newproduces3(token IToken )*produces3{
      my := new(produces3)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *produces3)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *produces3)       Enter(v Visitor){
        v.Visitproduces3(my)
        v.EndVisitproduces3(my)
    }


func AnyCastToproduces3(i interface{}) *produces3 {
	if nil == i{
		return nil
	}else{
		return i.(*produces3)
	}
}
/**
 *<b>
*<li>Rule 114:  symWithAttrs ::= EMPTY_KEY
 *</b>
 */
type symWithAttrs0 struct{
    *ASTNodeToken
}
func (my *symWithAttrs0)      GetEMPTY_KEY()IToken{ return my.leftIToken; }

func NewsymWithAttrs0(token IToken )*symWithAttrs0{
      my := new(symWithAttrs0)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *symWithAttrs0)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *symWithAttrs0)       Enter(v Visitor){
        v.VisitsymWithAttrs0(my)
        v.EndVisitsymWithAttrs0(my)
    }


func AnyCastTosymWithAttrs0(i interface{}) *symWithAttrs0 {
	if nil == i{
		return nil
	}else{
		return i.(*symWithAttrs0)
	}
}
/**
 *<b>
*<li>Rule 115:  symWithAttrs ::= SYMBOL optAttrList
 *</b>
 */
type symWithAttrs1 struct{
    *ASTNode
      _SYMBOL *ASTNodeToken
      _optAttrList *symAttrs
}
func (my *symWithAttrs1)      GetSYMBOL() *ASTNodeToken{ return my._SYMBOL}
func (my *symWithAttrs1)      SetSYMBOL( _SYMBOL *ASTNodeToken)  { my._SYMBOL = _SYMBOL }
    /**
     * The value returned by <b>GetoptAttrList</b> may be <b>null</b>
     */
func (my *symWithAttrs1)      GetoptAttrList() *symAttrs{ return my._optAttrList}
func (my *symWithAttrs1)      SetoptAttrList( _optAttrList *symAttrs)  { my._optAttrList = _optAttrList }

func NewsymWithAttrs1(leftIToken IToken, rightIToken IToken ,
              _SYMBOL *ASTNodeToken,
              _optAttrList *symAttrs)*symWithAttrs1{
      my := new(symWithAttrs1)
      my.ASTNode = NewASTNode2(leftIToken, rightIToken)
        my._SYMBOL = _SYMBOL;
        if nil != _SYMBOL{
        var trait_ interface{} = _SYMBOL
         trait_.(IAst).SetParent(my)
}
        my._optAttrList = _optAttrList;
        if nil != _optAttrList{
        var trait_ interface{} = _optAttrList
         trait_.(IAst).SetParent(my)
}
        my.Initialize()
        return my
    }

    /**
     * A list of all children of my node, don't including the null ones.
     */
func (my *symWithAttrs1)        GetAllChildren() * ArrayList{
        var list = NewArrayList()
        if nil != my._SYMBOL{  list.Add(my._SYMBOL) }
        if nil != my._optAttrList{  list.Add(my._optAttrList) }
        return list
    }

func (my *symWithAttrs1)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *symWithAttrs1)       Enter(v Visitor){
        var checkChildren = v.VisitsymWithAttrs1(my)
        if checkChildren{
            if nil != my._SYMBOL{my._SYMBOL.Accept(v)}
            if nil != my._optAttrList{my._optAttrList.Accept(v)}
        }
        v.EndVisitsymWithAttrs1(my)
    }


func AnyCastTosymWithAttrs1(i interface{}) *symWithAttrs1 {
	if nil == i{
		return nil
	}else{
		return i.(*symWithAttrs1)
	}
}
/**
 *<b>
*<li>Rule 123:  start_symbol ::= SYMBOL
 *</b>
 */
type start_symbol0 struct{
    *ASTNodeToken
}
func (my *start_symbol0)      GetSYMBOL()IToken{ return my.leftIToken; }

func Newstart_symbol0(token IToken )*start_symbol0{
      my := new(start_symbol0)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *start_symbol0)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *start_symbol0)       Enter(v Visitor){
        v.Visitstart_symbol0(my)
        v.EndVisitstart_symbol0(my)
    }


func AnyCastTostart_symbol0(i interface{}) *start_symbol0 {
	if nil == i{
		return nil
	}else{
		return i.(*start_symbol0)
	}
}
/**
 *<b>
*<li>Rule 124:  start_symbol ::= MACRO_NAME
 *</b>
 */
type start_symbol1 struct{
    *ASTNodeToken
}
func (my *start_symbol1)      GetMACRO_NAME()IToken{ return my.leftIToken; }

func Newstart_symbol1(token IToken )*start_symbol1{
      my := new(start_symbol1)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *start_symbol1)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *start_symbol1)       Enter(v Visitor){
        v.Visitstart_symbol1(my)
        v.EndVisitstart_symbol1(my)
    }


func AnyCastTostart_symbol1(i interface{}) *start_symbol1 {
	if nil == i{
		return nil
	}else{
		return i.(*start_symbol1)
	}
}
/**
 *<b>
*<li>Rule 130:  terminal_symbol ::= SYMBOL
 *</b>
 */
type terminal_symbol0 struct{
    *ASTNodeToken
}
func (my *terminal_symbol0)      GetSYMBOL()IToken{ return my.leftIToken; }

func Newterminal_symbol0(token IToken )*terminal_symbol0{
      my := new(terminal_symbol0)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *terminal_symbol0)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *terminal_symbol0)       Enter(v Visitor){
        v.Visitterminal_symbol0(my)
        v.EndVisitterminal_symbol0(my)
    }


func AnyCastToterminal_symbol0(i interface{}) *terminal_symbol0 {
	if nil == i{
		return nil
	}else{
		return i.(*terminal_symbol0)
	}
}
/**
 *<b>
*<li>Rule 131:  terminal_symbol ::= MACRO_NAME
 *</b>
 */
type terminal_symbol1 struct{
    *ASTNodeToken
}
func (my *terminal_symbol1)      GetMACRO_NAME()IToken{ return my.leftIToken; }

func Newterminal_symbol1(token IToken )*terminal_symbol1{
      my := new(terminal_symbol1)
      my.ASTNodeToken = NewASTNodeToken(token)
      my.Initialize()
      return my
    }

func (my *terminal_symbol1)       Accept(v IAstVisitor){
        if ! v.PreVisit(my){ return }
        var _ctor ,_ = v.(Visitor)
        my.Enter(_ctor)
        v.PostVisit(my)
}

func (my *terminal_symbol1)       Enter(v Visitor){
        v.Visitterminal_symbol1(my)
        v.EndVisitterminal_symbol1(my)
    }


func AnyCastToterminal_symbol1(i interface{}) *terminal_symbol1 {
	if nil == i{
		return nil
	}else{
		return i.(*terminal_symbol1)
	}
}
type Visitor interface{
  IAstVisitor
    Visit(n  IAst) bool
    EndVisit(n IAst)

    VisitASTNodeToken(n *ASTNodeToken) bool
    EndVisitASTNodeToken(n *ASTNodeToken)

    VisitLPG(n *LPG) bool
    EndVisitLPG(n *LPG)

    VisitLPG_itemList(n *LPG_itemList) bool
    EndVisitLPG_itemList(n *LPG_itemList)

    VisitAliasSeg(n *AliasSeg) bool
    EndVisitAliasSeg(n *AliasSeg)

    VisitAstSeg(n *AstSeg) bool
    EndVisitAstSeg(n *AstSeg)

    VisitDefineSeg(n *DefineSeg) bool
    EndVisitDefineSeg(n *DefineSeg)

    VisitEofSeg(n *EofSeg) bool
    EndVisitEofSeg(n *EofSeg)

    VisitEolSeg(n *EolSeg) bool
    EndVisitEolSeg(n *EolSeg)

    VisitErrorSeg(n *ErrorSeg) bool
    EndVisitErrorSeg(n *ErrorSeg)

    VisitExportSeg(n *ExportSeg) bool
    EndVisitExportSeg(n *ExportSeg)

    VisitGlobalsSeg(n *GlobalsSeg) bool
    EndVisitGlobalsSeg(n *GlobalsSeg)

    VisitHeadersSeg(n *HeadersSeg) bool
    EndVisitHeadersSeg(n *HeadersSeg)

    VisitIdentifierSeg(n *IdentifierSeg) bool
    EndVisitIdentifierSeg(n *IdentifierSeg)

    VisitImportSeg(n *ImportSeg) bool
    EndVisitImportSeg(n *ImportSeg)

    VisitIncludeSeg(n *IncludeSeg) bool
    EndVisitIncludeSeg(n *IncludeSeg)

    VisitKeywordsSeg(n *KeywordsSeg) bool
    EndVisitKeywordsSeg(n *KeywordsSeg)

    VisitNamesSeg(n *NamesSeg) bool
    EndVisitNamesSeg(n *NamesSeg)

    VisitNoticeSeg(n *NoticeSeg) bool
    EndVisitNoticeSeg(n *NoticeSeg)

    VisitRulesSeg(n *RulesSeg) bool
    EndVisitRulesSeg(n *RulesSeg)

    VisitSoftKeywordsSeg(n *SoftKeywordsSeg) bool
    EndVisitSoftKeywordsSeg(n *SoftKeywordsSeg)

    VisitStartSeg(n *StartSeg) bool
    EndVisitStartSeg(n *StartSeg)

    VisitTerminalsSeg(n *TerminalsSeg) bool
    EndVisitTerminalsSeg(n *TerminalsSeg)

    VisitTrailersSeg(n *TrailersSeg) bool
    EndVisitTrailersSeg(n *TrailersSeg)

    VisitTypesSeg(n *TypesSeg) bool
    EndVisitTypesSeg(n *TypesSeg)

    VisitRecoverSeg(n *RecoverSeg) bool
    EndVisitRecoverSeg(n *RecoverSeg)

    VisitPredecessorSeg(n *PredecessorSeg) bool
    EndVisitPredecessorSeg(n *PredecessorSeg)

    Visitoption_specList(n *option_specList) bool
    EndVisitoption_specList(n *option_specList)

    Visitoption_spec(n *option_spec) bool
    EndVisitoption_spec(n *option_spec)

    VisitoptionList(n *optionList) bool
    EndVisitoptionList(n *optionList)

    Visitoption(n *option) bool
    EndVisitoption(n *option)

    VisitSYMBOLList(n *SYMBOLList) bool
    EndVisitSYMBOLList(n *SYMBOLList)

    VisitaliasSpecList(n *aliasSpecList) bool
    EndVisitaliasSpecList(n *aliasSpecList)

    Visitalias_lhs_macro_name(n *alias_lhs_macro_name) bool
    EndVisitalias_lhs_macro_name(n *alias_lhs_macro_name)

    VisitdefineSpecList(n *defineSpecList) bool
    EndVisitdefineSpecList(n *defineSpecList)

    VisitdefineSpec(n *defineSpec) bool
    EndVisitdefineSpec(n *defineSpec)

    Visitmacro_segment(n *macro_segment) bool
    EndVisitmacro_segment(n *macro_segment)

    Visitterminal_symbolList(n *terminal_symbolList) bool
    EndVisitterminal_symbolList(n *terminal_symbolList)

    Visitaction_segmentList(n *action_segmentList) bool
    EndVisitaction_segmentList(n *action_segmentList)

    Visitimport_segment(n *import_segment) bool
    EndVisitimport_segment(n *import_segment)

    Visitdrop_commandList(n *drop_commandList) bool
    EndVisitdrop_commandList(n *drop_commandList)

    Visitdrop_ruleList(n *drop_ruleList) bool
    EndVisitdrop_ruleList(n *drop_ruleList)

    Visitdrop_rule(n *drop_rule) bool
    EndVisitdrop_rule(n *drop_rule)

    VisitoptMacroName(n *optMacroName) bool
    EndVisitoptMacroName(n *optMacroName)

    Visitinclude_segment(n *include_segment) bool
    EndVisitinclude_segment(n *include_segment)

    VisitkeywordSpecList(n *keywordSpecList) bool
    EndVisitkeywordSpecList(n *keywordSpecList)

    VisitkeywordSpec(n *keywordSpec) bool
    EndVisitkeywordSpec(n *keywordSpec)

    VisitnameSpecList(n *nameSpecList) bool
    EndVisitnameSpecList(n *nameSpecList)

    VisitnameSpec(n *nameSpec) bool
    EndVisitnameSpec(n *nameSpec)

    Visitrules_segment(n *rules_segment) bool
    EndVisitrules_segment(n *rules_segment)

    VisitnonTermList(n *nonTermList) bool
    EndVisitnonTermList(n *nonTermList)

    VisitnonTerm(n *nonTerm) bool
    EndVisitnonTerm(n *nonTerm)

    VisitRuleName(n *RuleName) bool
    EndVisitRuleName(n *RuleName)

    VisitruleList(n *ruleList) bool
    EndVisitruleList(n *ruleList)

    Visitrule(n *rule) bool
    EndVisitrule(n *rule)

    VisitsymWithAttrsList(n *symWithAttrsList) bool
    EndVisitsymWithAttrsList(n *symWithAttrsList)

    VisitsymAttrs(n *symAttrs) bool
    EndVisitsymAttrs(n *symAttrs)

    Visitaction_segment(n *action_segment) bool
    EndVisitaction_segment(n *action_segment)

    Visitstart_symbolList(n *start_symbolList) bool
    EndVisitstart_symbolList(n *start_symbolList)

    VisitterminalList(n *terminalList) bool
    EndVisitterminalList(n *terminalList)

    Visitterminal(n *terminal) bool
    EndVisitterminal(n *terminal)

    VisitoptTerminalAlias(n *optTerminalAlias) bool
    EndVisitoptTerminalAlias(n *optTerminalAlias)

    Visittype_declarationsList(n *type_declarationsList) bool
    EndVisittype_declarationsList(n *type_declarationsList)

    Visittype_declarations(n *type_declarations) bool
    EndVisittype_declarations(n *type_declarations)

    Visitsymbol_pairList(n *symbol_pairList) bool
    EndVisitsymbol_pairList(n *symbol_pairList)

    Visitsymbol_pair(n *symbol_pair) bool
    EndVisitsymbol_pair(n *symbol_pair)

    Visitrecover_symbol(n *recover_symbol) bool
    EndVisitrecover_symbol(n *recover_symbol)

    VisitEND_KEY_OPT(n *END_KEY_OPT) bool
    EndVisitEND_KEY_OPT(n *END_KEY_OPT)

    Visitoption_value0(n *option_value0) bool
    EndVisitoption_value0(n *option_value0)

    Visitoption_value1(n *option_value1) bool
    EndVisitoption_value1(n *option_value1)

    VisitaliasSpec0(n *aliasSpec0) bool
    EndVisitaliasSpec0(n *aliasSpec0)

    VisitaliasSpec1(n *aliasSpec1) bool
    EndVisitaliasSpec1(n *aliasSpec1)

    VisitaliasSpec2(n *aliasSpec2) bool
    EndVisitaliasSpec2(n *aliasSpec2)

    VisitaliasSpec3(n *aliasSpec3) bool
    EndVisitaliasSpec3(n *aliasSpec3)

    VisitaliasSpec4(n *aliasSpec4) bool
    EndVisitaliasSpec4(n *aliasSpec4)

    VisitaliasSpec5(n *aliasSpec5) bool
    EndVisitaliasSpec5(n *aliasSpec5)

    Visitalias_rhs0(n *alias_rhs0) bool
    EndVisitalias_rhs0(n *alias_rhs0)

    Visitalias_rhs1(n *alias_rhs1) bool
    EndVisitalias_rhs1(n *alias_rhs1)

    Visitalias_rhs2(n *alias_rhs2) bool
    EndVisitalias_rhs2(n *alias_rhs2)

    Visitalias_rhs3(n *alias_rhs3) bool
    EndVisitalias_rhs3(n *alias_rhs3)

    Visitalias_rhs4(n *alias_rhs4) bool
    EndVisitalias_rhs4(n *alias_rhs4)

    Visitalias_rhs5(n *alias_rhs5) bool
    EndVisitalias_rhs5(n *alias_rhs5)

    Visitalias_rhs6(n *alias_rhs6) bool
    EndVisitalias_rhs6(n *alias_rhs6)

    Visitmacro_name_symbol0(n *macro_name_symbol0) bool
    EndVisitmacro_name_symbol0(n *macro_name_symbol0)

    Visitmacro_name_symbol1(n *macro_name_symbol1) bool
    EndVisitmacro_name_symbol1(n *macro_name_symbol1)

    Visitdrop_command0(n *drop_command0) bool
    EndVisitdrop_command0(n *drop_command0)

    Visitdrop_command1(n *drop_command1) bool
    EndVisitdrop_command1(n *drop_command1)

    Visitname0(n *name0) bool
    EndVisitname0(n *name0)

    Visitname1(n *name1) bool
    EndVisitname1(n *name1)

    Visitname2(n *name2) bool
    EndVisitname2(n *name2)

    Visitname3(n *name3) bool
    EndVisitname3(n *name3)

    Visitname4(n *name4) bool
    EndVisitname4(n *name4)

    Visitname5(n *name5) bool
    EndVisitname5(n *name5)

    Visitproduces0(n *produces0) bool
    EndVisitproduces0(n *produces0)

    Visitproduces1(n *produces1) bool
    EndVisitproduces1(n *produces1)

    Visitproduces2(n *produces2) bool
    EndVisitproduces2(n *produces2)

    Visitproduces3(n *produces3) bool
    EndVisitproduces3(n *produces3)

    VisitsymWithAttrs0(n *symWithAttrs0) bool
    EndVisitsymWithAttrs0(n *symWithAttrs0)

    VisitsymWithAttrs1(n *symWithAttrs1) bool
    EndVisitsymWithAttrs1(n *symWithAttrs1)

    Visitstart_symbol0(n *start_symbol0) bool
    EndVisitstart_symbol0(n *start_symbol0)

    Visitstart_symbol1(n *start_symbol1) bool
    EndVisitstart_symbol1(n *start_symbol1)

    Visitterminal_symbol0(n *terminal_symbol0) bool
    EndVisitterminal_symbol0(n *terminal_symbol0)

    Visitterminal_symbol1(n *terminal_symbol1) bool
    EndVisitterminal_symbol1(n *terminal_symbol1)

}

func AnyCastToVisitor(i interface{}) Visitor {
	  if nil == i{
		 return nil
	  }else{
		 return i.(Visitor)
	  }
}
type AbstractVisitor struct{
   dispatch Visitor
   }
func NewAbstractVisitor(dispatch Visitor) *AbstractVisitor{
         my := new(AbstractVisitor)
         if dispatch != nil{
           my.dispatch = dispatch
         }else{
           my.dispatch = my
         }
        return my
}

func (my *AbstractVisitor)     UnimplementedVisitor(s  string)bool { return true }

func (my *AbstractVisitor)     PreVisit(element IAst) bool{ return true }

func (my *AbstractVisitor)     PostVisit(element  IAst) {}

func (my *AbstractVisitor)     VisitASTNodeToken(n  *ASTNodeToken) bool{ return my.UnimplementedVisitor("Visit(*ASTNodeToken)") }
func (my *AbstractVisitor)     EndVisitASTNodeToken(n  *ASTNodeToken) { my.UnimplementedVisitor("EndVisit(*ASTNodeToken)") }

func (my *AbstractVisitor)     VisitLPG(n  *LPG) bool{ return my.UnimplementedVisitor("Visit(*LPG)") }
func (my *AbstractVisitor)     EndVisitLPG(n  *LPG) { my.UnimplementedVisitor("EndVisit(*LPG)") }

func (my *AbstractVisitor)     VisitLPG_itemList(n  *LPG_itemList) bool{ return my.UnimplementedVisitor("Visit(*LPG_itemList)") }
func (my *AbstractVisitor)     EndVisitLPG_itemList(n  *LPG_itemList) { my.UnimplementedVisitor("EndVisit(*LPG_itemList)") }

func (my *AbstractVisitor)     VisitAliasSeg(n  *AliasSeg) bool{ return my.UnimplementedVisitor("Visit(*AliasSeg)") }
func (my *AbstractVisitor)     EndVisitAliasSeg(n  *AliasSeg) { my.UnimplementedVisitor("EndVisit(*AliasSeg)") }

func (my *AbstractVisitor)     VisitAstSeg(n  *AstSeg) bool{ return my.UnimplementedVisitor("Visit(*AstSeg)") }
func (my *AbstractVisitor)     EndVisitAstSeg(n  *AstSeg) { my.UnimplementedVisitor("EndVisit(*AstSeg)") }

func (my *AbstractVisitor)     VisitDefineSeg(n  *DefineSeg) bool{ return my.UnimplementedVisitor("Visit(*DefineSeg)") }
func (my *AbstractVisitor)     EndVisitDefineSeg(n  *DefineSeg) { my.UnimplementedVisitor("EndVisit(*DefineSeg)") }

func (my *AbstractVisitor)     VisitEofSeg(n  *EofSeg) bool{ return my.UnimplementedVisitor("Visit(*EofSeg)") }
func (my *AbstractVisitor)     EndVisitEofSeg(n  *EofSeg) { my.UnimplementedVisitor("EndVisit(*EofSeg)") }

func (my *AbstractVisitor)     VisitEolSeg(n  *EolSeg) bool{ return my.UnimplementedVisitor("Visit(*EolSeg)") }
func (my *AbstractVisitor)     EndVisitEolSeg(n  *EolSeg) { my.UnimplementedVisitor("EndVisit(*EolSeg)") }

func (my *AbstractVisitor)     VisitErrorSeg(n  *ErrorSeg) bool{ return my.UnimplementedVisitor("Visit(*ErrorSeg)") }
func (my *AbstractVisitor)     EndVisitErrorSeg(n  *ErrorSeg) { my.UnimplementedVisitor("EndVisit(*ErrorSeg)") }

func (my *AbstractVisitor)     VisitExportSeg(n  *ExportSeg) bool{ return my.UnimplementedVisitor("Visit(*ExportSeg)") }
func (my *AbstractVisitor)     EndVisitExportSeg(n  *ExportSeg) { my.UnimplementedVisitor("EndVisit(*ExportSeg)") }

func (my *AbstractVisitor)     VisitGlobalsSeg(n  *GlobalsSeg) bool{ return my.UnimplementedVisitor("Visit(*GlobalsSeg)") }
func (my *AbstractVisitor)     EndVisitGlobalsSeg(n  *GlobalsSeg) { my.UnimplementedVisitor("EndVisit(*GlobalsSeg)") }

func (my *AbstractVisitor)     VisitHeadersSeg(n  *HeadersSeg) bool{ return my.UnimplementedVisitor("Visit(*HeadersSeg)") }
func (my *AbstractVisitor)     EndVisitHeadersSeg(n  *HeadersSeg) { my.UnimplementedVisitor("EndVisit(*HeadersSeg)") }

func (my *AbstractVisitor)     VisitIdentifierSeg(n  *IdentifierSeg) bool{ return my.UnimplementedVisitor("Visit(*IdentifierSeg)") }
func (my *AbstractVisitor)     EndVisitIdentifierSeg(n  *IdentifierSeg) { my.UnimplementedVisitor("EndVisit(*IdentifierSeg)") }

func (my *AbstractVisitor)     VisitImportSeg(n  *ImportSeg) bool{ return my.UnimplementedVisitor("Visit(*ImportSeg)") }
func (my *AbstractVisitor)     EndVisitImportSeg(n  *ImportSeg) { my.UnimplementedVisitor("EndVisit(*ImportSeg)") }

func (my *AbstractVisitor)     VisitIncludeSeg(n  *IncludeSeg) bool{ return my.UnimplementedVisitor("Visit(*IncludeSeg)") }
func (my *AbstractVisitor)     EndVisitIncludeSeg(n  *IncludeSeg) { my.UnimplementedVisitor("EndVisit(*IncludeSeg)") }

func (my *AbstractVisitor)     VisitKeywordsSeg(n  *KeywordsSeg) bool{ return my.UnimplementedVisitor("Visit(*KeywordsSeg)") }
func (my *AbstractVisitor)     EndVisitKeywordsSeg(n  *KeywordsSeg) { my.UnimplementedVisitor("EndVisit(*KeywordsSeg)") }

func (my *AbstractVisitor)     VisitNamesSeg(n  *NamesSeg) bool{ return my.UnimplementedVisitor("Visit(*NamesSeg)") }
func (my *AbstractVisitor)     EndVisitNamesSeg(n  *NamesSeg) { my.UnimplementedVisitor("EndVisit(*NamesSeg)") }

func (my *AbstractVisitor)     VisitNoticeSeg(n  *NoticeSeg) bool{ return my.UnimplementedVisitor("Visit(*NoticeSeg)") }
func (my *AbstractVisitor)     EndVisitNoticeSeg(n  *NoticeSeg) { my.UnimplementedVisitor("EndVisit(*NoticeSeg)") }

func (my *AbstractVisitor)     VisitRulesSeg(n  *RulesSeg) bool{ return my.UnimplementedVisitor("Visit(*RulesSeg)") }
func (my *AbstractVisitor)     EndVisitRulesSeg(n  *RulesSeg) { my.UnimplementedVisitor("EndVisit(*RulesSeg)") }

func (my *AbstractVisitor)     VisitSoftKeywordsSeg(n  *SoftKeywordsSeg) bool{ return my.UnimplementedVisitor("Visit(*SoftKeywordsSeg)") }
func (my *AbstractVisitor)     EndVisitSoftKeywordsSeg(n  *SoftKeywordsSeg) { my.UnimplementedVisitor("EndVisit(*SoftKeywordsSeg)") }

func (my *AbstractVisitor)     VisitStartSeg(n  *StartSeg) bool{ return my.UnimplementedVisitor("Visit(*StartSeg)") }
func (my *AbstractVisitor)     EndVisitStartSeg(n  *StartSeg) { my.UnimplementedVisitor("EndVisit(*StartSeg)") }

func (my *AbstractVisitor)     VisitTerminalsSeg(n  *TerminalsSeg) bool{ return my.UnimplementedVisitor("Visit(*TerminalsSeg)") }
func (my *AbstractVisitor)     EndVisitTerminalsSeg(n  *TerminalsSeg) { my.UnimplementedVisitor("EndVisit(*TerminalsSeg)") }

func (my *AbstractVisitor)     VisitTrailersSeg(n  *TrailersSeg) bool{ return my.UnimplementedVisitor("Visit(*TrailersSeg)") }
func (my *AbstractVisitor)     EndVisitTrailersSeg(n  *TrailersSeg) { my.UnimplementedVisitor("EndVisit(*TrailersSeg)") }

func (my *AbstractVisitor)     VisitTypesSeg(n  *TypesSeg) bool{ return my.UnimplementedVisitor("Visit(*TypesSeg)") }
func (my *AbstractVisitor)     EndVisitTypesSeg(n  *TypesSeg) { my.UnimplementedVisitor("EndVisit(*TypesSeg)") }

func (my *AbstractVisitor)     VisitRecoverSeg(n  *RecoverSeg) bool{ return my.UnimplementedVisitor("Visit(*RecoverSeg)") }
func (my *AbstractVisitor)     EndVisitRecoverSeg(n  *RecoverSeg) { my.UnimplementedVisitor("EndVisit(*RecoverSeg)") }

func (my *AbstractVisitor)     VisitPredecessorSeg(n  *PredecessorSeg) bool{ return my.UnimplementedVisitor("Visit(*PredecessorSeg)") }
func (my *AbstractVisitor)     EndVisitPredecessorSeg(n  *PredecessorSeg) { my.UnimplementedVisitor("EndVisit(*PredecessorSeg)") }

func (my *AbstractVisitor)     Visitoption_specList(n  *option_specList) bool{ return my.UnimplementedVisitor("Visit(*option_specList)") }
func (my *AbstractVisitor)     EndVisitoption_specList(n  *option_specList) { my.UnimplementedVisitor("EndVisit(*option_specList)") }

func (my *AbstractVisitor)     Visitoption_spec(n  *option_spec) bool{ return my.UnimplementedVisitor("Visit(*option_spec)") }
func (my *AbstractVisitor)     EndVisitoption_spec(n  *option_spec) { my.UnimplementedVisitor("EndVisit(*option_spec)") }

func (my *AbstractVisitor)     VisitoptionList(n  *optionList) bool{ return my.UnimplementedVisitor("Visit(*optionList)") }
func (my *AbstractVisitor)     EndVisitoptionList(n  *optionList) { my.UnimplementedVisitor("EndVisit(*optionList)") }

func (my *AbstractVisitor)     Visitoption(n  *option) bool{ return my.UnimplementedVisitor("Visit(*option)") }
func (my *AbstractVisitor)     EndVisitoption(n  *option) { my.UnimplementedVisitor("EndVisit(*option)") }

func (my *AbstractVisitor)     VisitSYMBOLList(n  *SYMBOLList) bool{ return my.UnimplementedVisitor("Visit(*SYMBOLList)") }
func (my *AbstractVisitor)     EndVisitSYMBOLList(n  *SYMBOLList) { my.UnimplementedVisitor("EndVisit(*SYMBOLList)") }

func (my *AbstractVisitor)     VisitaliasSpecList(n  *aliasSpecList) bool{ return my.UnimplementedVisitor("Visit(*aliasSpecList)") }
func (my *AbstractVisitor)     EndVisitaliasSpecList(n  *aliasSpecList) { my.UnimplementedVisitor("EndVisit(*aliasSpecList)") }

func (my *AbstractVisitor)     Visitalias_lhs_macro_name(n  *alias_lhs_macro_name) bool{ return my.UnimplementedVisitor("Visit(*alias_lhs_macro_name)") }
func (my *AbstractVisitor)     EndVisitalias_lhs_macro_name(n  *alias_lhs_macro_name) { my.UnimplementedVisitor("EndVisit(*alias_lhs_macro_name)") }

func (my *AbstractVisitor)     VisitdefineSpecList(n  *defineSpecList) bool{ return my.UnimplementedVisitor("Visit(*defineSpecList)") }
func (my *AbstractVisitor)     EndVisitdefineSpecList(n  *defineSpecList) { my.UnimplementedVisitor("EndVisit(*defineSpecList)") }

func (my *AbstractVisitor)     VisitdefineSpec(n  *defineSpec) bool{ return my.UnimplementedVisitor("Visit(*defineSpec)") }
func (my *AbstractVisitor)     EndVisitdefineSpec(n  *defineSpec) { my.UnimplementedVisitor("EndVisit(*defineSpec)") }

func (my *AbstractVisitor)     Visitmacro_segment(n  *macro_segment) bool{ return my.UnimplementedVisitor("Visit(*macro_segment)") }
func (my *AbstractVisitor)     EndVisitmacro_segment(n  *macro_segment) { my.UnimplementedVisitor("EndVisit(*macro_segment)") }

func (my *AbstractVisitor)     Visitterminal_symbolList(n  *terminal_symbolList) bool{ return my.UnimplementedVisitor("Visit(*terminal_symbolList)") }
func (my *AbstractVisitor)     EndVisitterminal_symbolList(n  *terminal_symbolList) { my.UnimplementedVisitor("EndVisit(*terminal_symbolList)") }

func (my *AbstractVisitor)     Visitaction_segmentList(n  *action_segmentList) bool{ return my.UnimplementedVisitor("Visit(*action_segmentList)") }
func (my *AbstractVisitor)     EndVisitaction_segmentList(n  *action_segmentList) { my.UnimplementedVisitor("EndVisit(*action_segmentList)") }

func (my *AbstractVisitor)     Visitimport_segment(n  *import_segment) bool{ return my.UnimplementedVisitor("Visit(*import_segment)") }
func (my *AbstractVisitor)     EndVisitimport_segment(n  *import_segment) { my.UnimplementedVisitor("EndVisit(*import_segment)") }

func (my *AbstractVisitor)     Visitdrop_commandList(n  *drop_commandList) bool{ return my.UnimplementedVisitor("Visit(*drop_commandList)") }
func (my *AbstractVisitor)     EndVisitdrop_commandList(n  *drop_commandList) { my.UnimplementedVisitor("EndVisit(*drop_commandList)") }

func (my *AbstractVisitor)     Visitdrop_ruleList(n  *drop_ruleList) bool{ return my.UnimplementedVisitor("Visit(*drop_ruleList)") }
func (my *AbstractVisitor)     EndVisitdrop_ruleList(n  *drop_ruleList) { my.UnimplementedVisitor("EndVisit(*drop_ruleList)") }

func (my *AbstractVisitor)     Visitdrop_rule(n  *drop_rule) bool{ return my.UnimplementedVisitor("Visit(*drop_rule)") }
func (my *AbstractVisitor)     EndVisitdrop_rule(n  *drop_rule) { my.UnimplementedVisitor("EndVisit(*drop_rule)") }

func (my *AbstractVisitor)     VisitoptMacroName(n  *optMacroName) bool{ return my.UnimplementedVisitor("Visit(*optMacroName)") }
func (my *AbstractVisitor)     EndVisitoptMacroName(n  *optMacroName) { my.UnimplementedVisitor("EndVisit(*optMacroName)") }

func (my *AbstractVisitor)     Visitinclude_segment(n  *include_segment) bool{ return my.UnimplementedVisitor("Visit(*include_segment)") }
func (my *AbstractVisitor)     EndVisitinclude_segment(n  *include_segment) { my.UnimplementedVisitor("EndVisit(*include_segment)") }

func (my *AbstractVisitor)     VisitkeywordSpecList(n  *keywordSpecList) bool{ return my.UnimplementedVisitor("Visit(*keywordSpecList)") }
func (my *AbstractVisitor)     EndVisitkeywordSpecList(n  *keywordSpecList) { my.UnimplementedVisitor("EndVisit(*keywordSpecList)") }

func (my *AbstractVisitor)     VisitkeywordSpec(n  *keywordSpec) bool{ return my.UnimplementedVisitor("Visit(*keywordSpec)") }
func (my *AbstractVisitor)     EndVisitkeywordSpec(n  *keywordSpec) { my.UnimplementedVisitor("EndVisit(*keywordSpec)") }

func (my *AbstractVisitor)     VisitnameSpecList(n  *nameSpecList) bool{ return my.UnimplementedVisitor("Visit(*nameSpecList)") }
func (my *AbstractVisitor)     EndVisitnameSpecList(n  *nameSpecList) { my.UnimplementedVisitor("EndVisit(*nameSpecList)") }

func (my *AbstractVisitor)     VisitnameSpec(n  *nameSpec) bool{ return my.UnimplementedVisitor("Visit(*nameSpec)") }
func (my *AbstractVisitor)     EndVisitnameSpec(n  *nameSpec) { my.UnimplementedVisitor("EndVisit(*nameSpec)") }

func (my *AbstractVisitor)     Visitrules_segment(n  *rules_segment) bool{ return my.UnimplementedVisitor("Visit(*rules_segment)") }
func (my *AbstractVisitor)     EndVisitrules_segment(n  *rules_segment) { my.UnimplementedVisitor("EndVisit(*rules_segment)") }

func (my *AbstractVisitor)     VisitnonTermList(n  *nonTermList) bool{ return my.UnimplementedVisitor("Visit(*nonTermList)") }
func (my *AbstractVisitor)     EndVisitnonTermList(n  *nonTermList) { my.UnimplementedVisitor("EndVisit(*nonTermList)") }

func (my *AbstractVisitor)     VisitnonTerm(n  *nonTerm) bool{ return my.UnimplementedVisitor("Visit(*nonTerm)") }
func (my *AbstractVisitor)     EndVisitnonTerm(n  *nonTerm) { my.UnimplementedVisitor("EndVisit(*nonTerm)") }

func (my *AbstractVisitor)     VisitRuleName(n  *RuleName) bool{ return my.UnimplementedVisitor("Visit(*RuleName)") }
func (my *AbstractVisitor)     EndVisitRuleName(n  *RuleName) { my.UnimplementedVisitor("EndVisit(*RuleName)") }

func (my *AbstractVisitor)     VisitruleList(n  *ruleList) bool{ return my.UnimplementedVisitor("Visit(*ruleList)") }
func (my *AbstractVisitor)     EndVisitruleList(n  *ruleList) { my.UnimplementedVisitor("EndVisit(*ruleList)") }

func (my *AbstractVisitor)     Visitrule(n  *rule) bool{ return my.UnimplementedVisitor("Visit(*rule)") }
func (my *AbstractVisitor)     EndVisitrule(n  *rule) { my.UnimplementedVisitor("EndVisit(*rule)") }

func (my *AbstractVisitor)     VisitsymWithAttrsList(n  *symWithAttrsList) bool{ return my.UnimplementedVisitor("Visit(*symWithAttrsList)") }
func (my *AbstractVisitor)     EndVisitsymWithAttrsList(n  *symWithAttrsList) { my.UnimplementedVisitor("EndVisit(*symWithAttrsList)") }

func (my *AbstractVisitor)     VisitsymAttrs(n  *symAttrs) bool{ return my.UnimplementedVisitor("Visit(*symAttrs)") }
func (my *AbstractVisitor)     EndVisitsymAttrs(n  *symAttrs) { my.UnimplementedVisitor("EndVisit(*symAttrs)") }

func (my *AbstractVisitor)     Visitaction_segment(n  *action_segment) bool{ return my.UnimplementedVisitor("Visit(*action_segment)") }
func (my *AbstractVisitor)     EndVisitaction_segment(n  *action_segment) { my.UnimplementedVisitor("EndVisit(*action_segment)") }

func (my *AbstractVisitor)     Visitstart_symbolList(n  *start_symbolList) bool{ return my.UnimplementedVisitor("Visit(*start_symbolList)") }
func (my *AbstractVisitor)     EndVisitstart_symbolList(n  *start_symbolList) { my.UnimplementedVisitor("EndVisit(*start_symbolList)") }

func (my *AbstractVisitor)     VisitterminalList(n  *terminalList) bool{ return my.UnimplementedVisitor("Visit(*terminalList)") }
func (my *AbstractVisitor)     EndVisitterminalList(n  *terminalList) { my.UnimplementedVisitor("EndVisit(*terminalList)") }

func (my *AbstractVisitor)     Visitterminal(n  *terminal) bool{ return my.UnimplementedVisitor("Visit(*terminal)") }
func (my *AbstractVisitor)     EndVisitterminal(n  *terminal) { my.UnimplementedVisitor("EndVisit(*terminal)") }

func (my *AbstractVisitor)     VisitoptTerminalAlias(n  *optTerminalAlias) bool{ return my.UnimplementedVisitor("Visit(*optTerminalAlias)") }
func (my *AbstractVisitor)     EndVisitoptTerminalAlias(n  *optTerminalAlias) { my.UnimplementedVisitor("EndVisit(*optTerminalAlias)") }

func (my *AbstractVisitor)     Visittype_declarationsList(n  *type_declarationsList) bool{ return my.UnimplementedVisitor("Visit(*type_declarationsList)") }
func (my *AbstractVisitor)     EndVisittype_declarationsList(n  *type_declarationsList) { my.UnimplementedVisitor("EndVisit(*type_declarationsList)") }

func (my *AbstractVisitor)     Visittype_declarations(n  *type_declarations) bool{ return my.UnimplementedVisitor("Visit(*type_declarations)") }
func (my *AbstractVisitor)     EndVisittype_declarations(n  *type_declarations) { my.UnimplementedVisitor("EndVisit(*type_declarations)") }

func (my *AbstractVisitor)     Visitsymbol_pairList(n  *symbol_pairList) bool{ return my.UnimplementedVisitor("Visit(*symbol_pairList)") }
func (my *AbstractVisitor)     EndVisitsymbol_pairList(n  *symbol_pairList) { my.UnimplementedVisitor("EndVisit(*symbol_pairList)") }

func (my *AbstractVisitor)     Visitsymbol_pair(n  *symbol_pair) bool{ return my.UnimplementedVisitor("Visit(*symbol_pair)") }
func (my *AbstractVisitor)     EndVisitsymbol_pair(n  *symbol_pair) { my.UnimplementedVisitor("EndVisit(*symbol_pair)") }

func (my *AbstractVisitor)     Visitrecover_symbol(n  *recover_symbol) bool{ return my.UnimplementedVisitor("Visit(*recover_symbol)") }
func (my *AbstractVisitor)     EndVisitrecover_symbol(n  *recover_symbol) { my.UnimplementedVisitor("EndVisit(*recover_symbol)") }

func (my *AbstractVisitor)     VisitEND_KEY_OPT(n  *END_KEY_OPT) bool{ return my.UnimplementedVisitor("Visit(*END_KEY_OPT)") }
func (my *AbstractVisitor)     EndVisitEND_KEY_OPT(n  *END_KEY_OPT) { my.UnimplementedVisitor("EndVisit(*END_KEY_OPT)") }

func (my *AbstractVisitor)     Visitoption_value0(n  *option_value0) bool{ return my.UnimplementedVisitor("Visit(*option_value0)") }
func (my *AbstractVisitor)     EndVisitoption_value0(n  *option_value0) { my.UnimplementedVisitor("EndVisit(*option_value0)") }

func (my *AbstractVisitor)     Visitoption_value1(n  *option_value1) bool{ return my.UnimplementedVisitor("Visit(*option_value1)") }
func (my *AbstractVisitor)     EndVisitoption_value1(n  *option_value1) { my.UnimplementedVisitor("EndVisit(*option_value1)") }

func (my *AbstractVisitor)     VisitaliasSpec0(n  *aliasSpec0) bool{ return my.UnimplementedVisitor("Visit(*aliasSpec0)") }
func (my *AbstractVisitor)     EndVisitaliasSpec0(n  *aliasSpec0) { my.UnimplementedVisitor("EndVisit(*aliasSpec0)") }

func (my *AbstractVisitor)     VisitaliasSpec1(n  *aliasSpec1) bool{ return my.UnimplementedVisitor("Visit(*aliasSpec1)") }
func (my *AbstractVisitor)     EndVisitaliasSpec1(n  *aliasSpec1) { my.UnimplementedVisitor("EndVisit(*aliasSpec1)") }

func (my *AbstractVisitor)     VisitaliasSpec2(n  *aliasSpec2) bool{ return my.UnimplementedVisitor("Visit(*aliasSpec2)") }
func (my *AbstractVisitor)     EndVisitaliasSpec2(n  *aliasSpec2) { my.UnimplementedVisitor("EndVisit(*aliasSpec2)") }

func (my *AbstractVisitor)     VisitaliasSpec3(n  *aliasSpec3) bool{ return my.UnimplementedVisitor("Visit(*aliasSpec3)") }
func (my *AbstractVisitor)     EndVisitaliasSpec3(n  *aliasSpec3) { my.UnimplementedVisitor("EndVisit(*aliasSpec3)") }

func (my *AbstractVisitor)     VisitaliasSpec4(n  *aliasSpec4) bool{ return my.UnimplementedVisitor("Visit(*aliasSpec4)") }
func (my *AbstractVisitor)     EndVisitaliasSpec4(n  *aliasSpec4) { my.UnimplementedVisitor("EndVisit(*aliasSpec4)") }

func (my *AbstractVisitor)     VisitaliasSpec5(n  *aliasSpec5) bool{ return my.UnimplementedVisitor("Visit(*aliasSpec5)") }
func (my *AbstractVisitor)     EndVisitaliasSpec5(n  *aliasSpec5) { my.UnimplementedVisitor("EndVisit(*aliasSpec5)") }

func (my *AbstractVisitor)     Visitalias_rhs0(n  *alias_rhs0) bool{ return my.UnimplementedVisitor("Visit(*alias_rhs0)") }
func (my *AbstractVisitor)     EndVisitalias_rhs0(n  *alias_rhs0) { my.UnimplementedVisitor("EndVisit(*alias_rhs0)") }

func (my *AbstractVisitor)     Visitalias_rhs1(n  *alias_rhs1) bool{ return my.UnimplementedVisitor("Visit(*alias_rhs1)") }
func (my *AbstractVisitor)     EndVisitalias_rhs1(n  *alias_rhs1) { my.UnimplementedVisitor("EndVisit(*alias_rhs1)") }

func (my *AbstractVisitor)     Visitalias_rhs2(n  *alias_rhs2) bool{ return my.UnimplementedVisitor("Visit(*alias_rhs2)") }
func (my *AbstractVisitor)     EndVisitalias_rhs2(n  *alias_rhs2) { my.UnimplementedVisitor("EndVisit(*alias_rhs2)") }

func (my *AbstractVisitor)     Visitalias_rhs3(n  *alias_rhs3) bool{ return my.UnimplementedVisitor("Visit(*alias_rhs3)") }
func (my *AbstractVisitor)     EndVisitalias_rhs3(n  *alias_rhs3) { my.UnimplementedVisitor("EndVisit(*alias_rhs3)") }

func (my *AbstractVisitor)     Visitalias_rhs4(n  *alias_rhs4) bool{ return my.UnimplementedVisitor("Visit(*alias_rhs4)") }
func (my *AbstractVisitor)     EndVisitalias_rhs4(n  *alias_rhs4) { my.UnimplementedVisitor("EndVisit(*alias_rhs4)") }

func (my *AbstractVisitor)     Visitalias_rhs5(n  *alias_rhs5) bool{ return my.UnimplementedVisitor("Visit(*alias_rhs5)") }
func (my *AbstractVisitor)     EndVisitalias_rhs5(n  *alias_rhs5) { my.UnimplementedVisitor("EndVisit(*alias_rhs5)") }

func (my *AbstractVisitor)     Visitalias_rhs6(n  *alias_rhs6) bool{ return my.UnimplementedVisitor("Visit(*alias_rhs6)") }
func (my *AbstractVisitor)     EndVisitalias_rhs6(n  *alias_rhs6) { my.UnimplementedVisitor("EndVisit(*alias_rhs6)") }

func (my *AbstractVisitor)     Visitmacro_name_symbol0(n  *macro_name_symbol0) bool{ return my.UnimplementedVisitor("Visit(*macro_name_symbol0)") }
func (my *AbstractVisitor)     EndVisitmacro_name_symbol0(n  *macro_name_symbol0) { my.UnimplementedVisitor("EndVisit(*macro_name_symbol0)") }

func (my *AbstractVisitor)     Visitmacro_name_symbol1(n  *macro_name_symbol1) bool{ return my.UnimplementedVisitor("Visit(*macro_name_symbol1)") }
func (my *AbstractVisitor)     EndVisitmacro_name_symbol1(n  *macro_name_symbol1) { my.UnimplementedVisitor("EndVisit(*macro_name_symbol1)") }

func (my *AbstractVisitor)     Visitdrop_command0(n  *drop_command0) bool{ return my.UnimplementedVisitor("Visit(*drop_command0)") }
func (my *AbstractVisitor)     EndVisitdrop_command0(n  *drop_command0) { my.UnimplementedVisitor("EndVisit(*drop_command0)") }

func (my *AbstractVisitor)     Visitdrop_command1(n  *drop_command1) bool{ return my.UnimplementedVisitor("Visit(*drop_command1)") }
func (my *AbstractVisitor)     EndVisitdrop_command1(n  *drop_command1) { my.UnimplementedVisitor("EndVisit(*drop_command1)") }

func (my *AbstractVisitor)     Visitname0(n  *name0) bool{ return my.UnimplementedVisitor("Visit(*name0)") }
func (my *AbstractVisitor)     EndVisitname0(n  *name0) { my.UnimplementedVisitor("EndVisit(*name0)") }

func (my *AbstractVisitor)     Visitname1(n  *name1) bool{ return my.UnimplementedVisitor("Visit(*name1)") }
func (my *AbstractVisitor)     EndVisitname1(n  *name1) { my.UnimplementedVisitor("EndVisit(*name1)") }

func (my *AbstractVisitor)     Visitname2(n  *name2) bool{ return my.UnimplementedVisitor("Visit(*name2)") }
func (my *AbstractVisitor)     EndVisitname2(n  *name2) { my.UnimplementedVisitor("EndVisit(*name2)") }

func (my *AbstractVisitor)     Visitname3(n  *name3) bool{ return my.UnimplementedVisitor("Visit(*name3)") }
func (my *AbstractVisitor)     EndVisitname3(n  *name3) { my.UnimplementedVisitor("EndVisit(*name3)") }

func (my *AbstractVisitor)     Visitname4(n  *name4) bool{ return my.UnimplementedVisitor("Visit(*name4)") }
func (my *AbstractVisitor)     EndVisitname4(n  *name4) { my.UnimplementedVisitor("EndVisit(*name4)") }

func (my *AbstractVisitor)     Visitname5(n  *name5) bool{ return my.UnimplementedVisitor("Visit(*name5)") }
func (my *AbstractVisitor)     EndVisitname5(n  *name5) { my.UnimplementedVisitor("EndVisit(*name5)") }

func (my *AbstractVisitor)     Visitproduces0(n  *produces0) bool{ return my.UnimplementedVisitor("Visit(*produces0)") }
func (my *AbstractVisitor)     EndVisitproduces0(n  *produces0) { my.UnimplementedVisitor("EndVisit(*produces0)") }

func (my *AbstractVisitor)     Visitproduces1(n  *produces1) bool{ return my.UnimplementedVisitor("Visit(*produces1)") }
func (my *AbstractVisitor)     EndVisitproduces1(n  *produces1) { my.UnimplementedVisitor("EndVisit(*produces1)") }

func (my *AbstractVisitor)     Visitproduces2(n  *produces2) bool{ return my.UnimplementedVisitor("Visit(*produces2)") }
func (my *AbstractVisitor)     EndVisitproduces2(n  *produces2) { my.UnimplementedVisitor("EndVisit(*produces2)") }

func (my *AbstractVisitor)     Visitproduces3(n  *produces3) bool{ return my.UnimplementedVisitor("Visit(*produces3)") }
func (my *AbstractVisitor)     EndVisitproduces3(n  *produces3) { my.UnimplementedVisitor("EndVisit(*produces3)") }

func (my *AbstractVisitor)     VisitsymWithAttrs0(n  *symWithAttrs0) bool{ return my.UnimplementedVisitor("Visit(*symWithAttrs0)") }
func (my *AbstractVisitor)     EndVisitsymWithAttrs0(n  *symWithAttrs0) { my.UnimplementedVisitor("EndVisit(*symWithAttrs0)") }

func (my *AbstractVisitor)     VisitsymWithAttrs1(n  *symWithAttrs1) bool{ return my.UnimplementedVisitor("Visit(*symWithAttrs1)") }
func (my *AbstractVisitor)     EndVisitsymWithAttrs1(n  *symWithAttrs1) { my.UnimplementedVisitor("EndVisit(*symWithAttrs1)") }

func (my *AbstractVisitor)     Visitstart_symbol0(n  *start_symbol0) bool{ return my.UnimplementedVisitor("Visit(*start_symbol0)") }
func (my *AbstractVisitor)     EndVisitstart_symbol0(n  *start_symbol0) { my.UnimplementedVisitor("EndVisit(*start_symbol0)") }

func (my *AbstractVisitor)     Visitstart_symbol1(n  *start_symbol1) bool{ return my.UnimplementedVisitor("Visit(*start_symbol1)") }
func (my *AbstractVisitor)     EndVisitstart_symbol1(n  *start_symbol1) { my.UnimplementedVisitor("EndVisit(*start_symbol1)") }

func (my *AbstractVisitor)     Visitterminal_symbol0(n  *terminal_symbol0) bool{ return my.UnimplementedVisitor("Visit(*terminal_symbol0)") }
func (my *AbstractVisitor)     EndVisitterminal_symbol0(n  *terminal_symbol0) { my.UnimplementedVisitor("EndVisit(*terminal_symbol0)") }

func (my *AbstractVisitor)     Visitterminal_symbol1(n  *terminal_symbol1) bool{ return my.UnimplementedVisitor("Visit(*terminal_symbol1)") }
func (my *AbstractVisitor)     EndVisitterminal_symbol1(n  *terminal_symbol1) { my.UnimplementedVisitor("EndVisit(*terminal_symbol1)") }


func (my *AbstractVisitor)     Visit(n IAst) bool{
     switch n2 := n.(type) {
        case *ASTNodeToken:{
            return my.dispatch.VisitASTNodeToken(n2)
        }
        case *LPG:{
            return my.dispatch.VisitLPG(n2)
        }
        case *LPG_itemList:{
            return my.dispatch.VisitLPG_itemList(n2)
        }
        case *AliasSeg:{
            return my.dispatch.VisitAliasSeg(n2)
        }
        case *AstSeg:{
            return my.dispatch.VisitAstSeg(n2)
        }
        case *DefineSeg:{
            return my.dispatch.VisitDefineSeg(n2)
        }
        case *EofSeg:{
            return my.dispatch.VisitEofSeg(n2)
        }
        case *EolSeg:{
            return my.dispatch.VisitEolSeg(n2)
        }
        case *ErrorSeg:{
            return my.dispatch.VisitErrorSeg(n2)
        }
        case *ExportSeg:{
            return my.dispatch.VisitExportSeg(n2)
        }
        case *GlobalsSeg:{
            return my.dispatch.VisitGlobalsSeg(n2)
        }
        case *HeadersSeg:{
            return my.dispatch.VisitHeadersSeg(n2)
        }
        case *IdentifierSeg:{
            return my.dispatch.VisitIdentifierSeg(n2)
        }
        case *ImportSeg:{
            return my.dispatch.VisitImportSeg(n2)
        }
        case *IncludeSeg:{
            return my.dispatch.VisitIncludeSeg(n2)
        }
        case *KeywordsSeg:{
            return my.dispatch.VisitKeywordsSeg(n2)
        }
        case *NamesSeg:{
            return my.dispatch.VisitNamesSeg(n2)
        }
        case *NoticeSeg:{
            return my.dispatch.VisitNoticeSeg(n2)
        }
        case *RulesSeg:{
            return my.dispatch.VisitRulesSeg(n2)
        }
        case *SoftKeywordsSeg:{
            return my.dispatch.VisitSoftKeywordsSeg(n2)
        }
        case *StartSeg:{
            return my.dispatch.VisitStartSeg(n2)
        }
        case *TerminalsSeg:{
            return my.dispatch.VisitTerminalsSeg(n2)
        }
        case *TrailersSeg:{
            return my.dispatch.VisitTrailersSeg(n2)
        }
        case *TypesSeg:{
            return my.dispatch.VisitTypesSeg(n2)
        }
        case *RecoverSeg:{
            return my.dispatch.VisitRecoverSeg(n2)
        }
        case *PredecessorSeg:{
            return my.dispatch.VisitPredecessorSeg(n2)
        }
        case *option_specList:{
            return my.dispatch.Visitoption_specList(n2)
        }
        case *option_spec:{
            return my.dispatch.Visitoption_spec(n2)
        }
        case *optionList:{
            return my.dispatch.VisitoptionList(n2)
        }
        case *option:{
            return my.dispatch.Visitoption(n2)
        }
        case *SYMBOLList:{
            return my.dispatch.VisitSYMBOLList(n2)
        }
        case *aliasSpecList:{
            return my.dispatch.VisitaliasSpecList(n2)
        }
        case *alias_lhs_macro_name:{
            return my.dispatch.Visitalias_lhs_macro_name(n2)
        }
        case *defineSpecList:{
            return my.dispatch.VisitdefineSpecList(n2)
        }
        case *defineSpec:{
            return my.dispatch.VisitdefineSpec(n2)
        }
        case *macro_segment:{
            return my.dispatch.Visitmacro_segment(n2)
        }
        case *terminal_symbolList:{
            return my.dispatch.Visitterminal_symbolList(n2)
        }
        case *action_segmentList:{
            return my.dispatch.Visitaction_segmentList(n2)
        }
        case *import_segment:{
            return my.dispatch.Visitimport_segment(n2)
        }
        case *drop_commandList:{
            return my.dispatch.Visitdrop_commandList(n2)
        }
        case *drop_ruleList:{
            return my.dispatch.Visitdrop_ruleList(n2)
        }
        case *drop_rule:{
            return my.dispatch.Visitdrop_rule(n2)
        }
        case *optMacroName:{
            return my.dispatch.VisitoptMacroName(n2)
        }
        case *include_segment:{
            return my.dispatch.Visitinclude_segment(n2)
        }
        case *keywordSpecList:{
            return my.dispatch.VisitkeywordSpecList(n2)
        }
        case *keywordSpec:{
            return my.dispatch.VisitkeywordSpec(n2)
        }
        case *nameSpecList:{
            return my.dispatch.VisitnameSpecList(n2)
        }
        case *nameSpec:{
            return my.dispatch.VisitnameSpec(n2)
        }
        case *rules_segment:{
            return my.dispatch.Visitrules_segment(n2)
        }
        case *nonTermList:{
            return my.dispatch.VisitnonTermList(n2)
        }
        case *nonTerm:{
            return my.dispatch.VisitnonTerm(n2)
        }
        case *RuleName:{
            return my.dispatch.VisitRuleName(n2)
        }
        case *ruleList:{
            return my.dispatch.VisitruleList(n2)
        }
        case *rule:{
            return my.dispatch.Visitrule(n2)
        }
        case *symWithAttrsList:{
            return my.dispatch.VisitsymWithAttrsList(n2)
        }
        case *symAttrs:{
            return my.dispatch.VisitsymAttrs(n2)
        }
        case *action_segment:{
            return my.dispatch.Visitaction_segment(n2)
        }
        case *start_symbolList:{
            return my.dispatch.Visitstart_symbolList(n2)
        }
        case *terminalList:{
            return my.dispatch.VisitterminalList(n2)
        }
        case *terminal:{
            return my.dispatch.Visitterminal(n2)
        }
        case *optTerminalAlias:{
            return my.dispatch.VisitoptTerminalAlias(n2)
        }
        case *type_declarationsList:{
            return my.dispatch.Visittype_declarationsList(n2)
        }
        case *type_declarations:{
            return my.dispatch.Visittype_declarations(n2)
        }
        case *symbol_pairList:{
            return my.dispatch.Visitsymbol_pairList(n2)
        }
        case *symbol_pair:{
            return my.dispatch.Visitsymbol_pair(n2)
        }
        case *recover_symbol:{
            return my.dispatch.Visitrecover_symbol(n2)
        }
        case *END_KEY_OPT:{
            return my.dispatch.VisitEND_KEY_OPT(n2)
        }
        case *option_value0:{
            return my.dispatch.Visitoption_value0(n2)
        }
        case *option_value1:{
            return my.dispatch.Visitoption_value1(n2)
        }
        case *aliasSpec0:{
            return my.dispatch.VisitaliasSpec0(n2)
        }
        case *aliasSpec1:{
            return my.dispatch.VisitaliasSpec1(n2)
        }
        case *aliasSpec2:{
            return my.dispatch.VisitaliasSpec2(n2)
        }
        case *aliasSpec3:{
            return my.dispatch.VisitaliasSpec3(n2)
        }
        case *aliasSpec4:{
            return my.dispatch.VisitaliasSpec4(n2)
        }
        case *aliasSpec5:{
            return my.dispatch.VisitaliasSpec5(n2)
        }
        case *alias_rhs0:{
            return my.dispatch.Visitalias_rhs0(n2)
        }
        case *alias_rhs1:{
            return my.dispatch.Visitalias_rhs1(n2)
        }
        case *alias_rhs2:{
            return my.dispatch.Visitalias_rhs2(n2)
        }
        case *alias_rhs3:{
            return my.dispatch.Visitalias_rhs3(n2)
        }
        case *alias_rhs4:{
            return my.dispatch.Visitalias_rhs4(n2)
        }
        case *alias_rhs5:{
            return my.dispatch.Visitalias_rhs5(n2)
        }
        case *alias_rhs6:{
            return my.dispatch.Visitalias_rhs6(n2)
        }
        case *macro_name_symbol0:{
            return my.dispatch.Visitmacro_name_symbol0(n2)
        }
        case *macro_name_symbol1:{
            return my.dispatch.Visitmacro_name_symbol1(n2)
        }
        case *drop_command0:{
            return my.dispatch.Visitdrop_command0(n2)
        }
        case *drop_command1:{
            return my.dispatch.Visitdrop_command1(n2)
        }
        case *name0:{
            return my.dispatch.Visitname0(n2)
        }
        case *name1:{
            return my.dispatch.Visitname1(n2)
        }
        case *name2:{
            return my.dispatch.Visitname2(n2)
        }
        case *name3:{
            return my.dispatch.Visitname3(n2)
        }
        case *name4:{
            return my.dispatch.Visitname4(n2)
        }
        case *name5:{
            return my.dispatch.Visitname5(n2)
        }
        case *produces0:{
            return my.dispatch.Visitproduces0(n2)
        }
        case *produces1:{
            return my.dispatch.Visitproduces1(n2)
        }
        case *produces2:{
            return my.dispatch.Visitproduces2(n2)
        }
        case *produces3:{
            return my.dispatch.Visitproduces3(n2)
        }
        case *symWithAttrs0:{
            return my.dispatch.VisitsymWithAttrs0(n2)
        }
        case *symWithAttrs1:{
            return my.dispatch.VisitsymWithAttrs1(n2)
        }
        case *start_symbol0:{
            return my.dispatch.Visitstart_symbol0(n2)
        }
        case *start_symbol1:{
            return my.dispatch.Visitstart_symbol1(n2)
        }
        case *terminal_symbol0:{
            return my.dispatch.Visitterminal_symbol0(n2)
        }
        case *terminal_symbol1:{
            return my.dispatch.Visitterminal_symbol1(n2)
        }
        default:{ return false}
     }
}
func (my *AbstractVisitor)     EndVisit(n  IAst){
     switch n2 := n.(type) {
        case *ASTNodeToken:{
            my.dispatch.EndVisitASTNodeToken(n2)
        }
        case *LPG:{
            my.dispatch.EndVisitLPG(n2)
        }
        case *LPG_itemList:{
            my.dispatch.EndVisitLPG_itemList(n2)
        }
        case *AliasSeg:{
            my.dispatch.EndVisitAliasSeg(n2)
        }
        case *AstSeg:{
            my.dispatch.EndVisitAstSeg(n2)
        }
        case *DefineSeg:{
            my.dispatch.EndVisitDefineSeg(n2)
        }
        case *EofSeg:{
            my.dispatch.EndVisitEofSeg(n2)
        }
        case *EolSeg:{
            my.dispatch.EndVisitEolSeg(n2)
        }
        case *ErrorSeg:{
            my.dispatch.EndVisitErrorSeg(n2)
        }
        case *ExportSeg:{
            my.dispatch.EndVisitExportSeg(n2)
        }
        case *GlobalsSeg:{
            my.dispatch.EndVisitGlobalsSeg(n2)
        }
        case *HeadersSeg:{
            my.dispatch.EndVisitHeadersSeg(n2)
        }
        case *IdentifierSeg:{
            my.dispatch.EndVisitIdentifierSeg(n2)
        }
        case *ImportSeg:{
            my.dispatch.EndVisitImportSeg(n2)
        }
        case *IncludeSeg:{
            my.dispatch.EndVisitIncludeSeg(n2)
        }
        case *KeywordsSeg:{
            my.dispatch.EndVisitKeywordsSeg(n2)
        }
        case *NamesSeg:{
            my.dispatch.EndVisitNamesSeg(n2)
        }
        case *NoticeSeg:{
            my.dispatch.EndVisitNoticeSeg(n2)
        }
        case *RulesSeg:{
            my.dispatch.EndVisitRulesSeg(n2)
        }
        case *SoftKeywordsSeg:{
            my.dispatch.EndVisitSoftKeywordsSeg(n2)
        }
        case *StartSeg:{
            my.dispatch.EndVisitStartSeg(n2)
        }
        case *TerminalsSeg:{
            my.dispatch.EndVisitTerminalsSeg(n2)
        }
        case *TrailersSeg:{
            my.dispatch.EndVisitTrailersSeg(n2)
        }
        case *TypesSeg:{
            my.dispatch.EndVisitTypesSeg(n2)
        }
        case *RecoverSeg:{
            my.dispatch.EndVisitRecoverSeg(n2)
        }
        case *PredecessorSeg:{
            my.dispatch.EndVisitPredecessorSeg(n2)
        }
        case *option_specList:{
            my.dispatch.EndVisitoption_specList(n2)
        }
        case *option_spec:{
            my.dispatch.EndVisitoption_spec(n2)
        }
        case *optionList:{
            my.dispatch.EndVisitoptionList(n2)
        }
        case *option:{
            my.dispatch.EndVisitoption(n2)
        }
        case *SYMBOLList:{
            my.dispatch.EndVisitSYMBOLList(n2)
        }
        case *aliasSpecList:{
            my.dispatch.EndVisitaliasSpecList(n2)
        }
        case *alias_lhs_macro_name:{
            my.dispatch.EndVisitalias_lhs_macro_name(n2)
        }
        case *defineSpecList:{
            my.dispatch.EndVisitdefineSpecList(n2)
        }
        case *defineSpec:{
            my.dispatch.EndVisitdefineSpec(n2)
        }
        case *macro_segment:{
            my.dispatch.EndVisitmacro_segment(n2)
        }
        case *terminal_symbolList:{
            my.dispatch.EndVisitterminal_symbolList(n2)
        }
        case *action_segmentList:{
            my.dispatch.EndVisitaction_segmentList(n2)
        }
        case *import_segment:{
            my.dispatch.EndVisitimport_segment(n2)
        }
        case *drop_commandList:{
            my.dispatch.EndVisitdrop_commandList(n2)
        }
        case *drop_ruleList:{
            my.dispatch.EndVisitdrop_ruleList(n2)
        }
        case *drop_rule:{
            my.dispatch.EndVisitdrop_rule(n2)
        }
        case *optMacroName:{
            my.dispatch.EndVisitoptMacroName(n2)
        }
        case *include_segment:{
            my.dispatch.EndVisitinclude_segment(n2)
        }
        case *keywordSpecList:{
            my.dispatch.EndVisitkeywordSpecList(n2)
        }
        case *keywordSpec:{
            my.dispatch.EndVisitkeywordSpec(n2)
        }
        case *nameSpecList:{
            my.dispatch.EndVisitnameSpecList(n2)
        }
        case *nameSpec:{
            my.dispatch.EndVisitnameSpec(n2)
        }
        case *rules_segment:{
            my.dispatch.EndVisitrules_segment(n2)
        }
        case *nonTermList:{
            my.dispatch.EndVisitnonTermList(n2)
        }
        case *nonTerm:{
            my.dispatch.EndVisitnonTerm(n2)
        }
        case *RuleName:{
            my.dispatch.EndVisitRuleName(n2)
        }
        case *ruleList:{
            my.dispatch.EndVisitruleList(n2)
        }
        case *rule:{
            my.dispatch.EndVisitrule(n2)
        }
        case *symWithAttrsList:{
            my.dispatch.EndVisitsymWithAttrsList(n2)
        }
        case *symAttrs:{
            my.dispatch.EndVisitsymAttrs(n2)
        }
        case *action_segment:{
            my.dispatch.EndVisitaction_segment(n2)
        }
        case *start_symbolList:{
            my.dispatch.EndVisitstart_symbolList(n2)
        }
        case *terminalList:{
            my.dispatch.EndVisitterminalList(n2)
        }
        case *terminal:{
            my.dispatch.EndVisitterminal(n2)
        }
        case *optTerminalAlias:{
            my.dispatch.EndVisitoptTerminalAlias(n2)
        }
        case *type_declarationsList:{
            my.dispatch.EndVisittype_declarationsList(n2)
        }
        case *type_declarations:{
            my.dispatch.EndVisittype_declarations(n2)
        }
        case *symbol_pairList:{
            my.dispatch.EndVisitsymbol_pairList(n2)
        }
        case *symbol_pair:{
            my.dispatch.EndVisitsymbol_pair(n2)
        }
        case *recover_symbol:{
            my.dispatch.EndVisitrecover_symbol(n2)
        }
        case *END_KEY_OPT:{
            my.dispatch.EndVisitEND_KEY_OPT(n2)
        }
        case *option_value0:{
            my.dispatch.EndVisitoption_value0(n2)
        }
        case *option_value1:{
            my.dispatch.EndVisitoption_value1(n2)
        }
        case *aliasSpec0:{
            my.dispatch.EndVisitaliasSpec0(n2)
        }
        case *aliasSpec1:{
            my.dispatch.EndVisitaliasSpec1(n2)
        }
        case *aliasSpec2:{
            my.dispatch.EndVisitaliasSpec2(n2)
        }
        case *aliasSpec3:{
            my.dispatch.EndVisitaliasSpec3(n2)
        }
        case *aliasSpec4:{
            my.dispatch.EndVisitaliasSpec4(n2)
        }
        case *aliasSpec5:{
            my.dispatch.EndVisitaliasSpec5(n2)
        }
        case *alias_rhs0:{
            my.dispatch.EndVisitalias_rhs0(n2)
        }
        case *alias_rhs1:{
            my.dispatch.EndVisitalias_rhs1(n2)
        }
        case *alias_rhs2:{
            my.dispatch.EndVisitalias_rhs2(n2)
        }
        case *alias_rhs3:{
            my.dispatch.EndVisitalias_rhs3(n2)
        }
        case *alias_rhs4:{
            my.dispatch.EndVisitalias_rhs4(n2)
        }
        case *alias_rhs5:{
            my.dispatch.EndVisitalias_rhs5(n2)
        }
        case *alias_rhs6:{
            my.dispatch.EndVisitalias_rhs6(n2)
        }
        case *macro_name_symbol0:{
            my.dispatch.EndVisitmacro_name_symbol0(n2)
        }
        case *macro_name_symbol1:{
            my.dispatch.EndVisitmacro_name_symbol1(n2)
        }
        case *drop_command0:{
            my.dispatch.EndVisitdrop_command0(n2)
        }
        case *drop_command1:{
            my.dispatch.EndVisitdrop_command1(n2)
        }
        case *name0:{
            my.dispatch.EndVisitname0(n2)
        }
        case *name1:{
            my.dispatch.EndVisitname1(n2)
        }
        case *name2:{
            my.dispatch.EndVisitname2(n2)
        }
        case *name3:{
            my.dispatch.EndVisitname3(n2)
        }
        case *name4:{
            my.dispatch.EndVisitname4(n2)
        }
        case *name5:{
            my.dispatch.EndVisitname5(n2)
        }
        case *produces0:{
            my.dispatch.EndVisitproduces0(n2)
        }
        case *produces1:{
            my.dispatch.EndVisitproduces1(n2)
        }
        case *produces2:{
            my.dispatch.EndVisitproduces2(n2)
        }
        case *produces3:{
            my.dispatch.EndVisitproduces3(n2)
        }
        case *symWithAttrs0:{
            my.dispatch.EndVisitsymWithAttrs0(n2)
        }
        case *symWithAttrs1:{
            my.dispatch.EndVisitsymWithAttrs1(n2)
        }
        case *start_symbol0:{
            my.dispatch.EndVisitstart_symbol0(n2)
        }
        case *start_symbol1:{
            my.dispatch.EndVisitstart_symbol1(n2)
        }
        case *terminal_symbol0:{
            my.dispatch.EndVisitterminal_symbol0(n2)
        }
        case *terminal_symbol1:{
            my.dispatch.EndVisitterminal_symbol1(n2)
        }
        default:{ }
     }
}
func AnyCastToAbstractVisitor(i interface{}) *AbstractVisitor {
	if nil == i{
		return nil
	}else{
		return i.(*AbstractVisitor)
	}
}

