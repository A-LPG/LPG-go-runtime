
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
type JavaKWLexerprs struct{}
func NewJavaKWLexerprs() *JavaKWLexerprs{
    return &JavaKWLexerprs{}
}
const JavaKWLexerprs_ERROR_SYMBOL int = 0
func (my * JavaKWLexerprs) GetErrorSymbol() int {
     return JavaKWLexerprs_ERROR_SYMBOL
}
const JavaKWLexerprs_SCOPE_UBOUND int = 0
func (my * JavaKWLexerprs) GetScopeUbound() int {
     return JavaKWLexerprs_SCOPE_UBOUND
}
const JavaKWLexerprs_SCOPE_SIZE int = 0
func (my * JavaKWLexerprs) GetScopeSize() int {
     return JavaKWLexerprs_SCOPE_SIZE
}
const JavaKWLexerprs_MAX_NAME_LENGTH int = 0
func (my * JavaKWLexerprs) GetMaxNameLength() int {
     return JavaKWLexerprs_MAX_NAME_LENGTH
}
const JavaKWLexerprs_NUM_STATES int = 242
func (my * JavaKWLexerprs) GetNumStates() int {
     return JavaKWLexerprs_NUM_STATES
}
const JavaKWLexerprs_NT_OFFSET int = 56
func (my * JavaKWLexerprs) GetNtOffset() int {
     return JavaKWLexerprs_NT_OFFSET
}
const JavaKWLexerprs_LA_STATE_OFFSET int = 475
func (my * JavaKWLexerprs) GetLaStateOffset() int {
     return JavaKWLexerprs_LA_STATE_OFFSET
}
const JavaKWLexerprs_MAX_LA int = 1
func (my * JavaKWLexerprs) GetMaxLa() int {
     return JavaKWLexerprs_MAX_LA
}
const JavaKWLexerprs_NUM_RULES int = 88
func (my * JavaKWLexerprs) GetNumRules() int {
     return JavaKWLexerprs_NUM_RULES
}
const JavaKWLexerprs_NUM_NONTERMINALS int = 16
func (my * JavaKWLexerprs) GetNumNonterminals() int {
     return JavaKWLexerprs_NUM_NONTERMINALS
}
const JavaKWLexerprs_NUM_SYMBOLS int = 72
func (my * JavaKWLexerprs) GetNumSymbols() int {
     return JavaKWLexerprs_NUM_SYMBOLS
}
const JavaKWLexerprs_START_STATE int = 95
func (my * JavaKWLexerprs) GetStartState() int {
     return JavaKWLexerprs_START_STATE
}
const JavaKWLexerprs_IDENTIFIER_SYMBOL int = 0
func (my * JavaKWLexerprs) getIdentifier_SYMBOL() int {
     return JavaKWLexerprs_IDENTIFIER_SYMBOL
}
const JavaKWLexerprs_EOFT_SYMBOL int = 41
func (my * JavaKWLexerprs) GetEoftSymbol() int {
     return JavaKWLexerprs_EOFT_SYMBOL
}
const JavaKWLexerprs_EOLT_SYMBOL int = 57
func (my * JavaKWLexerprs) GetEoltSymbol() int {
     return JavaKWLexerprs_EOLT_SYMBOL
}
const JavaKWLexerprs_ACCEPT_ACTION int = 386
func (my * JavaKWLexerprs) GetAcceptAction() int {
     return JavaKWLexerprs_ACCEPT_ACTION
}
const JavaKWLexerprs_ERROR_ACTION int = 387
func (my * JavaKWLexerprs) GetErrorAction() int {
     return JavaKWLexerprs_ERROR_ACTION
}
const JavaKWLexerprs_BACKTRACK bool = false
func (my * JavaKWLexerprs) GetBacktrack() bool {
     return JavaKWLexerprs_BACKTRACK
}
func (my * JavaKWLexerprs) GetStartSymbol() int{
    return my.Lhs(0)
}
func (my * JavaKWLexerprs) IsValidForParser() bool{
    return JavaKWLexersym.IsValidForParser
}

var  JavaKWLexerprs_IsNullable []int=[]int{0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,
}
func (my * JavaKWLexerprs) IsNullable(index int)bool{
    return JavaKWLexerprs_IsNullable[index] != 0
}

var  JavaKWLexerprs_ProsthesesIndex []int=[]int{0,
            8,7,6,11,9,10,4,12,13,14,
            16,2,3,5,15,1,
}
func (my * JavaKWLexerprs) ProsthesesIndex(index int)int{
    return JavaKWLexerprs_ProsthesesIndex[index]
}

var  JavaKWLexerprs_IsKeyword []int=[]int{0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,
}
func (my * JavaKWLexerprs) IsKeyword(index int)bool{
    return JavaKWLexerprs_IsKeyword[index] != 0
}

var  JavaKWLexerprs_BaseCheck []int=[]int{0,
            8,6,7,5,4,4,5,4,5,5,
            8,7,2,6,4,4,7,5,5,7,
            5,3,4,2,10,6,10,3,9,4,
            6,3,4,7,7,9,6,6,5,6,
            8,5,6,12,4,5,6,9,4,3,
            4,8,5,12,10,10,8,9,11,10,
            1,1,1,1,1,1,1,1,1,1,
            1,1,1,1,1,1,1,1,1,1,
            1,1,1,1,1,1,1,1,
}
func (my * JavaKWLexerprs) BaseCheck(index int)int{
    return JavaKWLexerprs_BaseCheck[index]
}
var JavaKWLexerprs_Rhs  = JavaKWLexerprs_BaseCheck
func (my * JavaKWLexerprs) Rhs(index int) int{ return JavaKWLexerprs_Rhs[index] }

var  JavaKWLexerprs_BaseAction []int=[]int{
            12,12,12,12,12,12,12,12,12,12,
            12,12,12,12,12,12,12,12,12,12,
            12,12,12,12,12,12,12,12,12,12,
            12,12,12,12,12,12,12,12,12,12,
            12,12,12,12,12,12,12,12,12,12,
            12,12,12,12,12,12,12,12,12,12,
            12,1,1,13,13,5,5,10,10,7,
            7,14,14,3,3,8,8,11,11,2,
            2,4,4,6,6,15,15,9,9,14,
            50,90,60,115,113,1,110,48,120,117,
            64,134,92,66,164,98,116,97,132,118,
            97,161,106,129,119,133,100,135,136,149,
            152,135,257,145,130,120,141,156,137,127,
            160,123,172,145,165,171,144,177,263,146,
            149,176,150,103,187,181,166,153,202,190,
            194,167,148,178,158,160,165,170,197,57,
            205,214,166,173,199,218,210,223,178,174,
            221,179,58,227,236,186,22,183,233,239,
            188,56,190,242,28,60,247,55,248,256,
            249,189,59,194,260,56,54,70,82,266,
            30,270,35,95,106,44,100,13,6,107,
            77,65,83,276,267,112,278,279,269,117,
            277,206,281,286,150,118,284,289,231,292,
            294,293,295,75,299,301,302,303,310,311,
            124,304,315,316,319,320,324,322,323,326,
            332,333,334,327,339,341,342,345,347,349,
            350,351,353,354,355,358,359,361,365,137,
            367,368,374,369,376,370,384,375,385,138,
            380,387,381,391,398,392,399,401,404,407,
            410,402,408,413,414,415,417,420,418,421,
            423,432,435,436,427,431,441,437,443,448,
            450,451,452,454,455,458,459,462,463,442,
            465,468,470,469,474,473,476,481,486,482,
            484,489,491,492,494,498,495,501,502,504,
            506,507,508,509,514,512,519,522,523,524,
            527,528,533,535,536,539,540,541,203,545,
            546,549,550,551,556,553,561,552,562,563,
            566,569,570,574,577,579,567,580,583,586,
            584,589,590,591,592,387,387,
}
func (my * JavaKWLexerprs) BaseAction(index int)int{
    return JavaKWLexerprs_BaseAction[index]
}
var JavaKWLexerprs_Lhs  = JavaKWLexerprs_BaseAction
func (my * JavaKWLexerprs) Lhs(index int) int{ return JavaKWLexerprs_Lhs[index] }

var  JavaKWLexerprs_TermCheck []int=[]int{0,
            0,1,2,3,4,0,6,7,8,9,
            10,6,0,0,14,2,16,4,0,19,
            20,0,10,23,24,25,13,0,1,0,
            9,19,19,4,0,35,15,3,11,5,
            6,7,15,0,31,16,3,0,5,0,
            7,38,25,26,5,0,27,2,15,0,
            11,2,3,0,0,0,3,18,3,0,
            6,12,3,9,0,12,0,12,9,3,
            11,0,0,2,3,11,4,11,41,40,
            31,28,11,28,0,32,0,32,4,0,
            4,7,26,0,5,0,0,2,9,13,
            5,0,6,7,3,0,0,0,1,0,
            9,5,6,0,1,26,7,10,0,14,
            36,3,0,10,0,0,0,0,1,0,
            12,5,7,7,0,30,14,10,0,0,
            6,0,33,4,3,0,8,23,3,0,
            1,17,30,12,0,0,1,12,33,21,
            0,0,8,39,3,0,0,0,8,4,
            0,22,5,12,8,21,0,22,13,0,
            1,21,6,0,1,18,0,21,0,3,
            20,0,0,17,0,0,8,6,12,0,
            6,22,7,0,34,22,3,0,17,21,
            0,17,0,6,4,12,0,5,26,20,
            0,5,0,13,17,0,1,5,0,9,
            18,0,4,34,18,4,0,0,0,3,
            18,13,4,6,13,0,0,22,12,0,
            5,13,0,4,17,0,0,5,0,0,
            5,15,13,18,5,0,0,0,0,1,
            0,6,5,0,8,0,1,19,0,1,
            24,0,0,0,0,1,5,4,0,7,
            0,0,0,0,1,25,4,24,10,0,
            0,11,3,3,0,0,2,16,0,0,
            5,0,0,0,3,0,0,2,10,10,
            7,0,0,0,8,4,14,5,0,1,
            0,0,2,10,0,1,0,6,0,0,
            0,1,0,0,0,9,7,0,0,11,
            0,7,2,6,0,1,0,0,0,0,
            3,5,20,0,0,0,3,2,10,0,
            0,23,29,0,0,2,0,3,8,10,
            0,0,2,19,3,9,27,0,0,2,
            0,0,2,0,1,7,0,0,2,0,
            1,10,0,0,0,0,0,0,1,0,
            0,8,0,9,8,6,0,15,6,9,
            0,0,25,2,0,0,0,3,3,9,
            0,0,0,3,2,9,20,0,1,0,
            0,0,2,0,0,4,7,0,0,1,
            7,0,0,6,0,11,2,0,0,0,
            29,10,0,0,6,0,9,15,9,6,
            0,0,10,0,9,0,1,4,0,8,
            0,0,1,0,0,15,8,0,8,2,
            0,0,1,0,4,0,0,0,0,2,
            16,0,7,0,3,7,23,4,0,1,
            14,0,0,0,3,2,0,0,0,7,
            27,5,0,1,0,0,2,2,0,0,
            0,3,2,16,0,0,1,8,0,0,
            0,0,0,4,10,0,4,2,8,11,
            0,0,0,2,4,0,0,2,0,0,
            8,2,4,0,1,24,0,1,0,0,
            14,2,0,0,6,0,1,5,0,0,
            0,0,2,0,0,0,0,0,0,10,
            0,0,0,0,16,14,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            37,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,
}
func (my * JavaKWLexerprs) TermCheck(index int)int{
    return JavaKWLexerprs_TermCheck[index]
}

var  JavaKWLexerprs_TermAction []int=[]int{0,
            387,176,203,207,198,387,200,199,205,195,
            184,246,387,387,204,110,202,90,387,206,
            138,387,244,201,197,122,90,387,215,387,
            210,245,92,225,387,89,211,231,214,228,
            230,229,216,387,110,411,240,387,237,387,
            238,92,213,212,113,387,226,217,239,387,
            98,117,115,387,387,387,134,113,164,387,
            252,115,220,251,387,134,387,164,219,250,
            218,387,387,222,223,275,253,249,386,98,
            117,132,221,161,387,132,387,161,233,387,
            106,234,437,387,243,387,387,236,242,106,
            235,387,248,247,257,387,387,387,267,387,
            256,261,262,387,282,241,118,268,387,100,
            232,119,387,283,387,387,387,387,315,387,
            119,307,127,308,387,100,120,316,387,387,
            141,387,118,266,135,387,130,123,137,387,
            144,141,120,135,387,387,148,137,127,130,
            387,387,146,123,145,387,387,387,149,103,
            387,144,170,145,150,146,387,148,103,387,
            160,149,153,387,165,170,387,150,387,444,
            158,387,387,153,387,387,174,167,444,387,
            173,160,264,387,158,165,166,387,167,174,
            387,173,387,178,445,166,387,179,407,186,
            387,183,387,445,178,387,188,190,387,409,
            179,387,443,186,183,447,387,387,387,442,
            190,443,446,189,447,387,387,188,442,387,
            194,446,387,441,189,387,387,209,387,387,
            224,208,441,194,227,387,387,387,387,259,
            387,254,258,387,263,387,265,260,387,270,
            255,387,387,387,387,274,271,272,387,273,
            13,387,387,387,284,419,279,269,276,387,
            387,277,280,281,387,387,285,278,387,387,
            286,387,387,387,290,387,387,436,287,288,
            289,387,387,387,293,291,438,292,387,294,
            387,387,295,432,387,297,387,296,387,387,
            387,301,387,387,387,298,300,387,387,299,
            28,420,305,304,387,306,387,387,387,387,
            309,410,302,387,387,387,310,312,311,387,
            387,417,303,387,387,402,387,314,318,317,
            387,387,393,313,319,395,403,387,387,392,
            387,387,321,387,322,320,387,387,440,387,
            323,324,387,387,387,387,387,387,426,387,
            387,327,387,429,328,329,387,326,331,330,
            387,387,325,332,387,387,387,333,334,336,
            387,387,387,337,339,338,335,387,408,387,
            387,387,405,387,387,341,340,387,387,397,
            342,387,387,344,387,343,345,387,387,387,
            391,396,46,387,348,387,346,394,347,349,
            387,387,434,387,350,387,351,425,387,427,
            387,387,353,387,387,430,424,387,352,418,
            387,387,413,387,356,19,387,387,387,401,
            355,387,358,387,362,360,354,361,387,389,
            359,387,387,387,363,365,387,387,387,364,
            357,366,387,368,387,387,422,421,387,387,
            387,369,371,367,387,387,399,370,387,387,
            387,387,387,390,404,387,374,439,373,372,
            387,387,387,376,375,387,387,378,387,387,
            377,398,379,387,388,428,387,435,387,387,
            423,416,387,387,380,387,382,381,387,387,
            387,387,384,387,387,387,387,387,387,412,
            387,387,387,387,414,431,387,387,387,387,
            387,387,387,387,387,387,387,387,387,387,
            383,
}
func (my * JavaKWLexerprs) TermAction(index int)int{
    return JavaKWLexerprs_TermAction[index]
}
func (my * JavaKWLexerprs) Asb(index int) int{ return 0 }
func (my * JavaKWLexerprs) Asr(index int) int{ return 0 }
func (my * JavaKWLexerprs) Nasb(index int) int{ return 0 }
func (my * JavaKWLexerprs) Nasr(index int) int{ return 0 }
func (my * JavaKWLexerprs) TerminalIndex(index int) int{ return 0 }
func (my * JavaKWLexerprs) NonterminalIndex(index int) int{ return 0 }
func (my * JavaKWLexerprs) ScopePrefix(index int) int{ return 0 }
func (my * JavaKWLexerprs) ScopeSuffix(index int) int{ return 0 }
func (my * JavaKWLexerprs) ScopeLhs(index int) int{ return 0 }
func (my * JavaKWLexerprs) ScopeLa(index int) int{ return 0 }
func (my * JavaKWLexerprs) ScopeStateSet(index int) int{ return 0 }
func (my * JavaKWLexerprs) ScopeRhs(index int) int{ return 0 }
func (my * JavaKWLexerprs) ScopeState(index int) int{ return 0 }
func (my * JavaKWLexerprs) InSymb(index int) int{ return 0 }
func (my * JavaKWLexerprs) Name(index int)   string{ return "" }
func (my * JavaKWLexerprs) OriginalState(state int) int{
    return 0
}
func (my * JavaKWLexerprs) Asi(state int) int{
    return 0
}
func (my * JavaKWLexerprs) Nasi(state int ) int{
    return 0
}
func (my * JavaKWLexerprs) InSymbol(state int) int{
    return 0
}

    /**
     * assert(! goto_default);
     */
    func (my * JavaKWLexerprs) NtAction(state int,  sym int) int{
        return JavaKWLexerprs_BaseAction[state + sym]
    }

    /**
     * assert(! shift_default);
     */
    func (my * JavaKWLexerprs) TAction(state int,  sym int)int{
        var i = JavaKWLexerprs_BaseAction[state]
        var k = i + sym
        var index int
        if JavaKWLexerprs_TermCheck[k] == sym {
           index = k
        }else{
           index = i
        }
        return JavaKWLexerprs_TermAction[index]
    }
    func (my * JavaKWLexerprs) LookAhead(la_state int , sym int)int{
        var k = la_state + sym
        var index int
        if JavaKWLexerprs_TermCheck[k] == sym {
           index = k
        }else{
           index = la_state
        }
        return JavaKWLexerprs_TermAction[ index]
    }

