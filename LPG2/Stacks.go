package lpg2


type Stacks struct {
     STACK_INCREMENT int
     stateStackTop int
     stateStack []int
     locationStack []int
     parseStack []interface{}
}

func NewStacks() *Stacks {
    a := new(Stacks)
    a.stateStackTop  = 0
    a.STACK_INCREMENT = 1024
    return a
}




//
// Given a rule of the form     A ::= x1 x2 ... xn     n > 0
//
// the function GETTOKEN(i): yields the symbol xi, if xi is a terminal
// or ti, if xi is a nonterminal that produced a string of the form
// xi => ti w.
//
func (a *Stacks)  getToken( i int) int {
    return a.locationStack[a.stateStackTop+(i-1)]
}
//
// Given a rule of the form     A ::= x1 x2 ... xn     n > 0
//
// The function GETSYM(i): yields the AST subtree associated with symbol
// xi. NOTE that if xi is a terminal, GETSYM(i): is nil ! (However,
// see token_action below.):
//
// setSYM1(Object ast): is a function that allows us to assign an AST
// tree to GETSYM(1).
//
func (a *Stacks)  getSym( i int) interface{}{
    return a.parseStack[a.stateStackTop + (i - 1)]
}
func (a *Stacks)  setSym1(ast interface{}) {
    a.parseStack[a.stateStackTop] = ast
}
//
// Allocate or reallocate all the stacks. Their sizes should always be the same.
//
func (a *Stacks)  reallocateStacks(){
    var oldStackLength int = len(a.stateStack)
    var stackLength int = oldStackLength + a.STACK_INCREMENT

    if len(a.stateStack) == 0 {
        a.stateStack = make([]int, stackLength)
        a.locationStack =make([]int, stackLength)
        a.parseStack = make([]interface{}, stackLength, stackLength)
    }else {
        a.stateStack = arraycopy(a.stateStack, 0, make([]int, stackLength), 0, oldStackLength)
        a.locationStack = arraycopy(a.locationStack, 0, make([]int, stackLength), 0, oldStackLength)
        a.parseStack = ObjectArraycopy(a.parseStack, 0, make([]interface{}, stackLength), 0, oldStackLength)
    }
    return
}
//
// Allocate or reallocate the state stack only.
//
func (a *Stacks)  reallocateStateStack(){
    var oldStackLength int = len(a.stateStack)
    var stackLength int = oldStackLength + a.STACK_INCREMENT
    if len(a.stateStack) == 0{
        a.stateStack = make([]int, stackLength)
    }else{
        a.stateStack = arraycopy(a.stateStack, 0, make([]int, stackLength), 0, oldStackLength)
    }
    return
}
//
// Allocate location and parse stacks using the size of the state stack.
//
func (a *Stacks)  allocateOtherStacks() {
    var stackLength = len(a.stateStack)
    a.locationStack = make([]int, stackLength)
    a.parseStack = make([]interface{}, stackLength)
    return
}