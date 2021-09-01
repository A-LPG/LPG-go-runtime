package lpg2

// This Tuple type can be used to construct a dynamic
// array of integers. The space for the array is allocated in
// blocks of size 2**LOG_BLKSIZE. In declaring a tuple the user
// may specify an estimate of how many elements he expects.
// Based on that estimate, suitable values will be calculated
// for log_blksize and base_increment. If these estimates are
// found to be off later, more space will be allocated.
//

type ObjectTuple struct {
	top   int
	array []interface{}
}

// NewObjectTuple
// Constructor of a Tuple
//
func NewObjectTuple() *ObjectTuple {
	return NewObjectTupleWithEstimate(10)
}
func NewObjectTupleWithEstimate(estimate int) *ObjectTuple {
	return &ObjectTuple{
		top:   0,
		array: make([]interface{}, estimate),
	}
}

//
// This function is used to Reset the size of a dynamic array without
// allocating or deallocting space. It may be invoked with an integer
// argument n which indicates the  size or with no argument which
// indicates that the size should be Reset to 0.
//
func (my *ObjectTuple) ResetTo(n int) {
	my.top = n
}
func (my *ObjectTuple) Reset() {
	my.top = 0
}
func (my *ObjectTuple) capacity() int {
	return len(my.array)
}

//
// Return size of the dynamic array.
//
func (my *ObjectTuple) Size() int {
	return my.top
}

//
// Return a reference to the ith element of the dynamic array.
//
// Note that no check is made here to ensure that 0 <= i < top.
// Such a check might be useful for debugging and a range exception
// should be thrown if it yields true.
//
func (my *ObjectTuple) Get(i int) interface{} {
	if i < 0 || i >= len(my.array) {
		return nil
	}
	return my.array[i]
}

//
// Insert an element in the dynamic array at the location indicated.
//
func (my *ObjectTuple) Set(index int, value interface{}) {
	if index < 0 || index >= len(my.array) {
		return
	}
	my.array[index] = value
}

//
// Add an element to the dynamic array and return the top index.
//
func (my *ObjectTuple) NextIndex() int {
	var i int = my.top
	my.top += 1
	if i >= len(my.array) {
		my.array = ObjectArraycopy(my.array, 0, make([]interface{}, i*2, i*2), 0, i)
	}

	return i
}

//
// Add an element to the dynamic array and return a reference to
// that  element.
//
func (my *ObjectTuple) Add(element interface{}) {
	var i int = my.NextIndex()
	my.array[i] = element
}
