package lpg2
type AstArrayList struct {
	array []IAst
}
func NewAstArrayListFrom(array []IAst) *AstArrayList {
	return &AstArrayList{
		array: array,
	}
}
func NewAstArrayListFromCopy(array []IAst) *AstArrayList {
	newArray := make([]IAst, len(array))
	copy(newArray, array)
	return &AstArrayList{
		array: newArray,
	}
}
func NewAstArrayListSize(size int, cap int) *AstArrayList {
	return &AstArrayList{
		array: make([]IAst, size, cap),
	}
}
func NewAstArrayList() *AstArrayList {
	return NewAstArrayListSize(0,0)
}
func (a *AstArrayList) clone() *AstArrayList {
	array := make([]IAst, len(a.array))
	copy(array, a.array)
	return NewAstArrayListFrom(array)
}

func (a *AstArrayList) clear() bool{
	if len(a.array) > 0 {
		a.array = make([]IAst, 0)
	}
	return  true
}
func (a *AstArrayList) removeAt(index int)  (value IAst, found bool){
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
func (a *AstArrayList) remove(value IAst) bool{
	if i := a.search(value); i != -1 {
		_, found := a.removeAt(i)
		return found
	}
	return false
}
func (a *AstArrayList) search(value IAst) int {

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

func (a *AstArrayList) removeAll() bool{
	return a.clear()
}

func (a *AstArrayList) toArray() []IAst {
	array := make([]IAst, len(a.array))
	copy(array, a.array)
	return array
}

func (a *AstArrayList) size() int {
	return len(a.array)
}
func (a *AstArrayList) add(elem IAst) *AstArrayList {
	a.array = append(a.array, elem)
	return a
}
func (a *AstArrayList) get(index int) IAst{
	if index < 0 || index >= len(a.array) {
		return nil
	}
	return a.array[index]
}
func (a *AstArrayList) at(index int) (value IAst) {
	return a.get(index)
}
func (a *AstArrayList) contains( val IAst) bool{
	return a.search(val) != -1
}
func (a *AstArrayList) isEmpty() bool{
	return a.size() == 0
}
func (a *AstArrayList) set(index int, element IAst) bool {
	if index < 0 || index >= len(a.array) {
		return  false
	}
	a.array[index] = element
	return  true
}
func (a *AstArrayList) indexOf(elem IAst) int {
	return a.search(elem)
}
func (a *AstArrayList) lastIndexOf(  elem IAst ) int{
	var size = a.size()
	for i:= size; i > 0; i--{
		if a.array[size - i - 1] == elem{
			return size - i - 1
		}

	}
	return -1
}

