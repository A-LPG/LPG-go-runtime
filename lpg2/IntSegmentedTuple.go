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


type IntSegmentedTuple struct {
     top int
     _size      int
     logBlksize int
     baseSize   int
     base       [][]int
}


func  NewIntSegmentedTuple(logBlksize int, baseSize int ) *IntSegmentedTuple {
        a := new(IntSegmentedTuple)
        a.top= 0
        a._size= 0

        a.logBlksize = 3
        a.baseSize = 4

        a.logBlksize = logBlksize

        if baseSize <= 0 {
            a.baseSize = 4
        }else {
            a.baseSize = baseSize
        }
        a.base = make([][]int,a.baseSize)
        return a
}
func   NewIntSegmentedTupleDefault() *IntSegmentedTuple {
    return NewIntSegmentedTuple(3,4)
}
    //
    // Allocate another block of storage for the dynamic array.
    //
func (a *IntSegmentedTuple) arraycopy(src [][]int, srcPos int,
    dest [][]int, destPos int, length int) [][]int{
    for i:=0; i < length; i++ {
        dest[destPos+ i] = src[srcPos+ i]
    }
    return dest
}
func (a *IntSegmentedTuple) allocateMoreSpace(){
        //
        // The variable size always indicates the maximum int of
        // elements that has been allocated for the array.
        // Initially, it is Set to 0 to indicate that the array is empty.
        // The pool of available elements is divided into segments of size
        // 2**log_blksize each. Each segment is pointed to by a slot in
        // the array base.
        //
        // By dividing size by the size of the segment we obtain the
        // index for the next segment in base. If base is full, it is
        // reallocated.
        //
        //
        var k int = a._size >> a.logBlksize // which segment?

        //
        // If the base is overflowed, reallocate it and initialize the 
        // elements to NULL.
        // Otherwise, allocate a  segment and place its adjusted address
        // in base[k]. The adjustment allows us to index the segment directly,
        // instead of having to perform a subtraction for each reference.
        // See operator[] below.
        //
        //
        if k == a.baseSize {
            a.baseSize *= 2
            a.base = a.arraycopy(a.base, 0, make([][]int,a.baseSize), 0, k)
        }

        a.base[k] = make([]int,1 << a.logBlksize)

        //
        // Finally, we update SIZE.
        //
        a._size += 1 << a.logBlksize

        return
}
    //
    // This function is invoked with an integer argument n. It ensures
    // that enough space is allocated for n elements in the dynamic array.
    // I.e., that the array will be indexable in the range  (0..n-1){
    //
    // Note that a function can be used as a garbage collector.  When
    // invoked with no argument(or 0){, it frees up all dynamic space that
    // was allocated for the array.
    //
func (a *IntSegmentedTuple) Resize(){
    a.ResizeTo(0)
}
func (a *IntSegmentedTuple) ResizeTo( n int ){
    //
    // If array did not previously contain enough space, allocate
    // the necessary additional space. Otherwise, if the array had
    // more blocks than are needed, release the extra blocks.
    //
    if n > a._size {
        for ;; {
            a.allocateMoreSpace()
            if  n > a._size == false{
                break
            }
        }
    }

    a.top  = n
}
    //
    // This function is used to Reset the size of a dynamic array without
    // allocating or deallocting space. It may be invoked with an integer
    // argument n which indicates the  size or with no argument which
    // indicates that the size should be Reset to 0.
    //
func (a *IntSegmentedTuple) ReSet() {
   a.ReSetTo(0)
}
func (a *IntSegmentedTuple) ReSetTo( n int) {
    a.top = n
}
//
// Return size of the dynamic array.
//

func (a *IntSegmentedTuple) Size() int{
        return a.top
}
//
// Can the tuple be indexed with i?
//
func (a *IntSegmentedTuple) outOfRange( i int) bool {
    return i < 0 || i >= a.top
}
//
// Return a reference to the ith element of the dynamic array.
//
// Note that no check is made here to ensure that 0 <= i < top.
// Such a check might be useful for debugging and a range exception
// should be thrown if it yields true.
//
func (a *IntSegmentedTuple) Get( i int) int {

    return a.base[i>>a.logBlksize][i%(1<<a.logBlksize)]
}
//
// Insert an element in the dynamic array at the location indicated.
//
func (a *IntSegmentedTuple) Set( i int, element int) {

    a.base[i>>a.logBlksize][i%(1<<a.logBlksize)] = element
}

// NextIndex
// Add an element to the dynamic array and return the top index.
//
func (a *IntSegmentedTuple) NextIndex() int {
    var i  = a.top
    a.top++
    if i == a._size {
        a.allocateMoreSpace()
    }
    return i
}
//
// Add an element to the dynamic array and return a reference to
// that  element.
//
func (a *IntSegmentedTuple) Add( element int) {
    var i  = a.NextIndex()
    a.base[i>>a.logBlksize][i%(1<<a.logBlksize)] = element
}
//
// If array is sorted, a function will find the index location
// of a given element if it is contained in the array. Otherwise, it
// will return the negation of the index of the element prior to
// which the  element would be inserted in the array.
//
func (a *IntSegmentedTuple) BinarySearch( element int) int {
    var low = 0
    var high = a.top
    for ;high > low; {
        var mid int= int((high + low) / 2)
        var midElement int = a.Get(mid)
        if element == midElement {
            return mid
        } else {
            if element < midElement {
            high = mid
            } else {
            low = mid + 1
            }
        }
    }

    return -low
}