package lpg2


type IntArrayList struct {
	array []int
}



func NewIntArrayListFrom(array []int) *IntArrayList {
	return &IntArrayList{
		array: array,
	}
}
func NewIntArrayListFromCopy(array []int) *IntArrayList {
	newArray := make([]int, len(array))
	copy(newArray, array)
	return &IntArrayList{
		array: newArray,
	}
}
func NewIntArrayListSize(size int, cap int) *IntArrayList {
	return &IntArrayList{
		array: make([]int, size, cap),
	}
}
func  NewIntArrayList() *IntArrayList {
	return NewIntArrayListSize(0,0)
}
func (a *IntArrayList) Clone() *IntArrayList {
	array := make([]int, len(a.array))
	copy(array, a.array)
	return NewIntArrayListFrom(array)
}

func (a *IntArrayList) Clear() bool{
	if len(a.array) > 0 {
		a.array = make([]int, 0)
	}
	return  true
}
func (a *IntArrayList) RemoveAt(index int)  (value int, found bool){
	if index < 0 || index >= len(a.array) {
		return -1, false
	}
	// Determine array boundaries when deleting to improve deletion efficiency.
	if index == 0 {
		value := a.array[0]
		a.array = a.array[1:]
		return value, true
	} else if index == len(a.array)-1 {
		value := a.array[index]
		a.array = a.array[:index]
		return value, true
	}
	// If it is a non-boundary delete,
	// it will involve the creation of an array,
	// then the deletion is less efficient.
	value = a.array[index]
	a.array = append(a.array[:index], a.array[index+1:]...)
	return value, true
}
func (a *IntArrayList) Remove(value int) bool{
	if i := a.Search(value); i != -1 {
		_, found := a.RemoveAt(i)
		return found
	}
	return false
}
func (a *IntArrayList) Search(value int) int {

	if len(a.array) == 0 {
		return -1
	}
	result := -1
	for index, v := range a.array {
		if v == value {
			result = index
			break
		}
	}
	return result
}

func (a *IntArrayList) RemoveAll() bool{
	return a.Clear()
}

func (a *IntArrayList) ToArray() []int {
	array := make([]int, len(a.array))
	copy(array, a.array)
	return array
}

func (a *IntArrayList) Size() int {
	return len(a.array)
}
func (a *IntArrayList) Add(elem int) *IntArrayList{
	a.array = append(a.array, elem)
	return a
}
func (a *IntArrayList) Get(index int) int{
	if index < 0 || index >= len(a.array) {
		return -1
	}
	return a.array[index]
}
func (a *IntArrayList) At(index int) (value int) {
	return a.Get(index)
}
func (a *IntArrayList) Contains( val int) bool{
	return a.Search(val) != -1
}
func (a *IntArrayList) IsEmpty() bool{
	return a.Size() == 0
}
func (a *IntArrayList) Set(index int, element int) bool {
	if index < 0 || index >= len(a.array) {
		return  false
	}
	a.array[index] = element
	return  true
}
func (a *IntArrayList) IndexOf(elem int) int {
	return a.Search(elem)
}
func (a *IntArrayList) LastIndexOf(  elem int ) int{
	var size = a.Size()
	for i:= size; i > 0; i--{
		if a.array[size - i - 1] == elem{
			return size - i - 1
		}

	}
	return -1
}