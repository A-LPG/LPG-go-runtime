
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
type LPGKWLexerprs struct{}
func NewLPGKWLexerprs() *LPGKWLexerprs{
    return &LPGKWLexerprs{}
}
const LPGKWLexerprs_ERROR_SYMBOL int = 0
func (my * LPGKWLexerprs) GetErrorSymbol() int {
     return LPGKWLexerprs_ERROR_SYMBOL
}
const LPGKWLexerprs_SCOPE_UBOUND int = 0
func (my * LPGKWLexerprs) GetScopeUbound() int {
     return LPGKWLexerprs_SCOPE_UBOUND
}
const LPGKWLexerprs_SCOPE_SIZE int = 0
func (my * LPGKWLexerprs) GetScopeSize() int {
     return LPGKWLexerprs_SCOPE_SIZE
}
const LPGKWLexerprs_MAX_NAME_LENGTH int = 0
func (my * LPGKWLexerprs) GetMaxNameLength() int {
     return LPGKWLexerprs_MAX_NAME_LENGTH
}
const LPGKWLexerprs_NUM_STATES int = 145
func (my * LPGKWLexerprs) GetNumStates() int {
     return LPGKWLexerprs_NUM_STATES
}
const LPGKWLexerprs_NT_OFFSET int = 30
func (my * LPGKWLexerprs) GetNtOffset() int {
     return LPGKWLexerprs_NT_OFFSET
}
const LPGKWLexerprs_LA_STATE_OFFSET int = 208
func (my * LPGKWLexerprs) GetLaStateOffset() int {
     return LPGKWLexerprs_LA_STATE_OFFSET
}
const LPGKWLexerprs_MAX_LA int = 0
func (my * LPGKWLexerprs) GetMaxLa() int {
     return LPGKWLexerprs_MAX_LA
}
const LPGKWLexerprs_NUM_RULES int = 29
func (my * LPGKWLexerprs) GetNumRules() int {
     return LPGKWLexerprs_NUM_RULES
}
const LPGKWLexerprs_NUM_NONTERMINALS int = 3
func (my * LPGKWLexerprs) GetNumNonterminals() int {
     return LPGKWLexerprs_NUM_NONTERMINALS
}
const LPGKWLexerprs_NUM_SYMBOLS int = 33
func (my * LPGKWLexerprs) GetNumSymbols() int {
     return LPGKWLexerprs_NUM_SYMBOLS
}
const LPGKWLexerprs_START_STATE int = 30
func (my * LPGKWLexerprs) GetStartState() int {
     return LPGKWLexerprs_START_STATE
}
const LPGKWLexerprs_IDENTIFIER_SYMBOL int = 0
func (my * LPGKWLexerprs) getIdentifier_SYMBOL() int {
     return LPGKWLexerprs_IDENTIFIER_SYMBOL
}
const LPGKWLexerprs_EOFT_SYMBOL int = 27
func (my * LPGKWLexerprs) GetEoftSymbol() int {
     return LPGKWLexerprs_EOFT_SYMBOL
}
const LPGKWLexerprs_EOLT_SYMBOL int = 31
func (my * LPGKWLexerprs) GetEoltSymbol() int {
     return LPGKWLexerprs_EOLT_SYMBOL
}
const LPGKWLexerprs_ACCEPT_ACTION int = 178
func (my * LPGKWLexerprs) GetAcceptAction() int {
     return LPGKWLexerprs_ACCEPT_ACTION
}
const LPGKWLexerprs_ERROR_ACTION int = 179
func (my * LPGKWLexerprs) GetErrorAction() int {
     return LPGKWLexerprs_ERROR_ACTION
}
const LPGKWLexerprs_BACKTRACK bool = false
func (my * LPGKWLexerprs) GetBacktrack() bool {
     return LPGKWLexerprs_BACKTRACK
}
func (my * LPGKWLexerprs) GetStartSymbol() int{
    return my.Lhs(0)
}
func (my * LPGKWLexerprs) IsValidForParser() bool{
    return LPGKWLexersym.IsValidForParser
}

var  LPGKWLexerprs_IsNullable []int=[]int{0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,
}
func (my * LPGKWLexerprs) IsNullable(index int)bool{
    return LPGKWLexerprs_IsNullable[index] != 0
}

var  LPGKWLexerprs_ProsthesesIndex []int=[]int{0,
            2,3,1,
}
func (my * LPGKWLexerprs) ProsthesesIndex(index int)int{
    return LPGKWLexerprs_ProsthesesIndex[index]
}

var  LPGKWLexerprs_IsKeyword []int=[]int{0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
}
func (my * LPGKWLexerprs) IsKeyword(index int)bool{
    return LPGKWLexerprs_IsKeyword[index] != 0
}

var  LPGKWLexerprs_BaseCheck []int=[]int{0,
            6,4,7,24,10,12,6,4,6,4,
            4,7,8,8,11,7,8,9,13,6,
            7,10,8,6,6,9,6,1,1,
}
func (my * LPGKWLexerprs) BaseCheck(index int)int{
    return LPGKWLexerprs_BaseCheck[index]
}
var LPGKWLexerprs_Rhs  = LPGKWLexerprs_BaseCheck
func (my * LPGKWLexerprs) Rhs(index int) int{ return LPGKWLexerprs_Rhs[index] }

var  LPGKWLexerprs_BaseAction []int=[]int{
            1,1,1,1,1,1,1,1,1,1,
            1,1,1,1,1,1,1,1,1,1,
            1,1,1,1,1,1,1,1,2,2,
            26,33,34,37,1,39,12,44,45,62,
            22,65,67,17,35,51,68,14,5,70,
            30,71,72,73,77,80,42,79,82,86,
            85,87,54,88,95,96,97,100,99,103,
            105,109,112,117,113,115,120,121,125,124,
            123,130,131,8,132,133,134,136,139,143,
            145,146,147,149,148,156,153,157,161,162,
            165,159,167,166,176,178,180,184,185,186,
            174,57,168,190,191,194,196,198,201,203,
            206,205,207,210,215,213,217,211,219,221,
            225,228,229,230,223,233,239,234,241,244,
            245,247,248,250,253,237,259,262,255,263,
            265,267,271,273,276,277,278,281,282,283,
            288,286,292,296,298,293,303,287,305,307,
            308,311,313,312,317,319,320,179,179,
}
func (my * LPGKWLexerprs) BaseAction(index int)int{
    return LPGKWLexerprs_BaseAction[index]
}
var LPGKWLexerprs_Lhs  = LPGKWLexerprs_BaseAction
func (my * LPGKWLexerprs) Lhs(index int) int{ return LPGKWLexerprs_Lhs[index] }

var  LPGKWLexerprs_TermCheck []int=[]int{0,
            0,1,2,3,0,5,6,0,8,9,
            10,0,1,0,3,11,0,10,18,3,
            4,0,22,23,13,0,10,14,12,0,
            9,10,3,12,0,1,0,3,0,1,
            6,0,26,0,0,20,21,4,4,5,
            0,8,2,0,16,14,0,7,2,3,
            7,0,1,27,0,1,0,0,15,0,
            0,0,0,7,7,5,0,8,0,0,
            8,0,1,12,0,0,0,0,4,11,
            3,15,13,8,0,0,0,11,0,0,
            4,2,0,9,0,0,11,5,0,1,
            6,0,0,15,0,4,0,1,6,0,
            0,1,0,0,0,6,12,3,5,0,
            0,0,0,0,4,0,7,4,0,4,
            9,19,0,5,0,0,0,0,0,17,
            2,6,0,11,8,0,0,2,0,7,
            0,0,6,2,0,0,0,0,24,5,
            4,4,25,0,14,0,18,0,3,0,
            1,16,5,0,0,0,13,3,3,0,
            0,8,2,0,1,0,1,0,0,10,
            0,1,0,1,0,0,0,10,3,0,
            0,5,0,9,0,6,0,3,0,7,
            0,5,0,13,0,1,6,0,0,0,
            3,3,0,0,16,13,0,8,0,1,
            0,9,2,0,0,2,0,0,15,0,
            0,2,0,7,0,19,12,10,0,7,
            2,0,0,1,0,0,0,6,2,5,
            0,17,0,1,4,0,0,0,2,4,
            0,0,0,3,3,0,0,0,11,7,
            3,0,0,2,9,0,1,0,0,2,
            14,9,0,1,0,1,0,0,2,2,
            0,0,0,2,4,3,0,1,0,0,
            0,2,0,5,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,
}
func (my * LPGKWLexerprs) TermCheck(index int)int{
    return LPGKWLexerprs_TermCheck[index]
}

var  LPGKWLexerprs_TermAction []int=[]int{0,
            179,43,38,35,179,36,40,179,45,44,
            37,179,50,179,49,73,179,105,39,63,
            62,179,42,41,48,179,64,72,65,179,
            58,56,75,57,179,68,179,66,179,47,
            67,179,61,179,179,34,34,51,54,53,
            179,52,69,179,46,81,179,70,127,128,
            189,179,55,178,179,59,179,179,190,179,
            179,179,179,60,71,76,179,74,179,179,
            78,179,83,77,179,179,179,179,85,82,
            87,79,80,84,179,179,179,86,179,179,
            89,90,179,187,179,179,88,181,179,93,
            92,179,179,91,179,94,179,95,96,179,
            179,99,179,179,179,98,97,100,101,179,
            179,179,179,179,104,179,103,108,179,109,
            106,102,179,110,179,179,179,179,179,107,
            203,113,179,111,114,179,179,206,179,116,
            179,179,117,199,179,179,179,179,112,204,
            120,129,115,179,118,179,119,179,122,179,
            124,121,123,179,179,179,186,126,188,179,
            179,125,180,179,131,179,132,179,179,130,
            179,200,179,134,179,179,179,133,135,179,
            179,195,179,136,179,137,179,138,179,139,
            179,191,179,140,179,182,142,179,179,179,
            202,143,179,179,141,145,179,144,179,196,
            179,146,193,179,179,192,179,179,147,179,
            179,205,179,149,179,152,148,150,179,151,
            197,179,179,155,179,179,179,153,201,156,
            179,154,179,158,157,179,179,179,184,159,
            179,179,179,161,194,179,179,179,160,162,
            163,179,179,185,164,179,165,179,179,198,
            168,166,179,167,179,169,179,179,170,171,
            179,179,179,174,172,173,179,175,179,179,
            179,183,179,176,
}
func (my * LPGKWLexerprs) TermAction(index int)int{
    return LPGKWLexerprs_TermAction[index]
}
func (my * LPGKWLexerprs) Asb(index int) int{ return 0 }
func (my * LPGKWLexerprs) Asr(index int) int{ return 0 }
func (my * LPGKWLexerprs) Nasb(index int) int{ return 0 }
func (my * LPGKWLexerprs) Nasr(index int) int{ return 0 }
func (my * LPGKWLexerprs) TerminalIndex(index int) int{ return 0 }
func (my * LPGKWLexerprs) NonterminalIndex(index int) int{ return 0 }
func (my * LPGKWLexerprs) ScopePrefix(index int) int{ return 0 }
func (my * LPGKWLexerprs) ScopeSuffix(index int) int{ return 0 }
func (my * LPGKWLexerprs) ScopeLhs(index int) int{ return 0 }
func (my * LPGKWLexerprs) ScopeLa(index int) int{ return 0 }
func (my * LPGKWLexerprs) ScopeStateSet(index int) int{ return 0 }
func (my * LPGKWLexerprs) ScopeRhs(index int) int{ return 0 }
func (my * LPGKWLexerprs) ScopeState(index int) int{ return 0 }
func (my * LPGKWLexerprs) InSymb(index int) int{ return 0 }
func (my * LPGKWLexerprs) Name(index int)   string{ return "" }
func (my * LPGKWLexerprs) OriginalState(state int) int{
    return 0
}
func (my * LPGKWLexerprs) Asi(state int) int{
    return 0
}
func (my * LPGKWLexerprs) Nasi(state int ) int{
    return 0
}
func (my * LPGKWLexerprs) InSymbol(state int) int{
    return 0
}

    /**
     * assert(! goto_default);
     */
    func (my * LPGKWLexerprs) NtAction(state int,  sym int) int{
        return LPGKWLexerprs_BaseAction[state + sym]
    }

    /**
     * assert(! shift_default);
     */
    func (my * LPGKWLexerprs) TAction(state int,  sym int)int{
        var i = LPGKWLexerprs_BaseAction[state]
        var k = i + sym
        var index int
        if LPGKWLexerprs_TermCheck[k] == sym {
           index = k
        }else{
           index = i
        }
        return LPGKWLexerprs_TermAction[index]
    }
    func (my * LPGKWLexerprs) LookAhead(la_state int , sym int)int{
        var k = la_state + sym
        var index int
        if LPGKWLexerprs_TermCheck[k] == sym {
           index = k
        }else{
           index = la_state
        }
        return LPGKWLexerprs_TermAction[ index]
    }

