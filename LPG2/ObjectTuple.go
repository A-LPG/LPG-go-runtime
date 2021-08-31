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
// This function is used to reset the size of a dynamic array without
// allocating or deallocting space. It may be invoked with an integer
// argument n which indicates the  size or with no argument which
// indicates that the size should be reset to 0.
//
func (self *ObjectTuple) resetTo(n int) {
	self.top = n
}
func (self *ObjectTuple) reset() {
	self.top = 0
}
func (self *ObjectTuple) capacity() int {
	return len(self.array)
}

//
// Return size of the dynamic array.
//
func (self *ObjectTuple) size() int {
	return self.top
}

//
// Return a reference to the ith element of the dynamic array.
//
// Note that no check is made here to ensure that 0 <= i < top.
// Such a check might be useful for debugging and a range exception
// should be thrown if it yields true.
//
func (self *ObjectTuple) get(i int) interface{} {
	if i < 0 || i >= len(self.array) {
		return nil
	}
	return self.array[i]
}

//
// Insert an element in the dynamic array at the location indicated.
//
func (self *ObjectTuple) set(index int, value interface{}) {
	if index < 0 || index >= len(self.array) {
		return
	}
	self.array[index] = value
}

//
// Add an element to the dynamic array and return the top index.
//
func (self *ObjectTuple) nextIndex() int {
	var i int = self.top
	self.top += 1
	if i >= len(self.array) {
		self.array = ObjectArraycopy(self.array, 0, make([]interface{}, i*2, i*2), 0, i)
	}

	return i
}

//
// Add an element to the dynamic array and return a reference to
// that  element.
//
func (self *ObjectTuple) add(element interface{}) {
	var i int = self.nextIndex()
	self.array[i] = element
}
