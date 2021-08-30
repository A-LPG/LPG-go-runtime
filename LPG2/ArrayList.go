package lpg2
type ArrayList struct {
	array []IAst
}
func NewArrayListFrom(array []IAst) *ArrayList {
	return &ArrayList{
		array: array,
	}
}
func NewArrayListFromCopy(array []IAst, safe ...bool) *ArrayList {
	newArray := make([]IAst, len(array))
	copy(newArray, array)
	return &ArrayList{
		array: newArray,
	}
}
func NewArrayListSize(size int, cap int) *ArrayList {
	return &ArrayList{
		array: make([]IAst, size, cap),
	}
}
func  NewArrayList() *ArrayList {
	return NewArrayListSize(0,0)
}
func (a *ArrayList) clone() *ArrayList {
	array := make([]IAst, len(a.array))
	copy(array, a.array)
	return NewArrayListFrom(array)
}

func (a *ArrayList) clear() bool{
	if len(a.array) > 0 {
		a.array = make([]IAst, 0)
	}
	return  true
}
func (a *ArrayList) removeAt(index int)  (value IAst, found bool){
	if index < 0 || index >= len(a.array) {
		return nil, false
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
func (a *ArrayList) remove(value IAst) bool{
	if i := a.search(value); i != -1 {
		_, found := a.removeAt(i)
		return found
	}
	return false
}
func (a *ArrayList) search(value IAst) int {

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

func (a *ArrayList) removeAll() bool{
	return a.clear()
}

func (a *ArrayList) toArray() []IAst {
	array := make([]IAst, len(a.array))
	copy(array, a.array)
	return array
}

func (a *ArrayList) size() int {
	return len(a.array)
}
func (a *ArrayList) add(elem IAst) *ArrayList{
	a.array = append(a.array, elem)
	return a
}
func (a *ArrayList) get(index int) IAst{
	if index < 0 || index >= len(a.array) {
		return nil
	}
	return a.array[index]
}
func (a *ArrayList) at(index int) (value IAst) {
	return a.get(index)
}
func (a *ArrayList) contains( val IAst) bool{
	return a.search(val) != -1
}
func (a *ArrayList) isEmpty() bool{
	return a.size() == 0
}
func (a *ArrayList) set(index int, element IAst) bool {
	if index < 0 || index >= len(a.array) {
		return  false
	}
	a.array[index] = element
	return  true
}
func (a *ArrayList) indexOf(elem IAst) int {
	return a.search(elem)
}
func (a *ArrayList) lastIndexOf(  elem IAst ) int{
	var size = a.size()
	for i:= size; i > 0; i--{
		if a.array[size - i - 1] == elem{
			return size - i - 1
		}

	}
	return -1
}

