package lpg2


//
// This Tuple type can be used to construct a dynamic
// array of integers. The space for the array is allocated in
// blocks of size 2**LOG_BLKSIZE. In declaring a tuple the user
// may specify an estimate of how many elements he expects.
// Based on that estimate, suitable values will be calculated
// for log_blksize and base_increment. If these estimates are
// found to be off later, more space will be allocated.
//


type IntTuple struct {
     array  []int
     top int
}

//
// Constructor of a Tuple
//
func NewIntTuple() *IntTuple{
    return NewIntTupleWithEstimate(10)
}
func NewIntTupleWithEstimate(estimate int ) *IntTuple{
    return &IntTuple{
        top: 0,
        array: make([]int, estimate, estimate),
    }
}
//
// This function is used to reset the size of a dynamic array without
// allocating or deallocting space. It may be invoked with an integer
// argument n which indicates the  size or with no argument which
// indicates that the size should be reset to 0.
//
func(a *IntTuple) resetFrom( n int ){
    a.top = n
}
func(a *IntTuple) reset(){
    a.top = 0
}
func(a *IntTuple) capacity() int {
    return len(a.array)
}
//
// Return size of the dynamic array.
//
func(a *IntTuple) size() int {
    return a.top
}
//
// Return a reference to the ith element of the dynamic array.
//
// Note that no check is made here to ensure that 0 <= i < top.
// Such a check might be useful for debugging and a range exception
// should be thrown if it yields true.
//
func(a *IntTuple) get(i int) int{
    return a.array[i]
}
//
// Insert an element in the dynamic array at the location indicated.
//
func(a *IntTuple) set( index int, value int) {
    if index < 0 || index >= len(a.array) {
        return
    }
    a.array[index] = value
}


//
// Add an element to the dynamic array and return the top index.
//
func(a *IntTuple) nextIndex() int {
    var i int = a.top
    a.top += 1
    if i >=  len(a.array){
        a.array = arraycopy(a.array,0, make([]int, i * 2, i * 2),0,i)
    }


    return i
}
//
// Add an element to the dynamic array and return a reference to
// that  element.
//
func(a *IntTuple) add( element int){
    var i int = a.nextIndex()
    a.array[i] = element
}