package lpg2
type ArrayList struct {
	array []interface{}
}
func NewArrayListFrom(array []interface{}) *ArrayList {
	return &ArrayList{
		array: array,
	}
}
func NewArrayListFromCopy(array []interface{}) *ArrayList {
	newArray := make([]interface{}, len(array))
	copy(newArray, array)
	return &ArrayList{
		array: newArray,
	}
}
func NewArrayListSize(size int, cap int) *ArrayList {
	return &ArrayList{
		array: make([]interface{}, size, cap),
	}
}
func NewArrayList() *ArrayList {
	return NewArrayListSize(0,0)
}
func (a *ArrayList) Clone() *ArrayList {
	array := make([]interface{}, len(a.array))
	copy(array, a.array)
	return NewArrayListFrom(array)
}

func (a *ArrayList) Clear() bool{
	if len(a.array) > 0 {
		a.array = make([]interface{}, 0)
	}
	return  true
}
func (a *ArrayList) RemoveAt(index int)  (value interface{}, found bool){
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
func (a *ArrayList) Remove(value interface{}) bool{
	if i := a.Search(value); i != -1 {
		_, found := a.RemoveAt(i)
		return found
	}
	return false
}
func (a *ArrayList) Search(value interface{}) int {

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

func (a *ArrayList) RemoveAll() bool{
	return a.Clear()
}

func (a *ArrayList) ToArray() []interface{} {
	array := make([]interface{}, len(a.array))
	copy(array, a.array)
	return array
}

func (a *ArrayList) Size() int {
	return len(a.array)
}
func (a *ArrayList) Add(elem interface{}) *ArrayList {
	a.array = append(a.array, elem)
	return a
}
func (a *ArrayList) Get(index int) interface{}{
	if index < 0 || index >= len(a.array) {
		return nil
	}
	return a.array[index]
}
func (a *ArrayList) At(index int) (value interface{}) {
	return a.Get(index)
}
func (a *ArrayList) Contains( val interface{}) bool{
	return a.Search(val) != -1
}
func (a *ArrayList) IsEmpty() bool{
	return a.Size() == 0
}
func (a *ArrayList) Set(index int, element interface{}) bool {
	if index < 0 || index >= len(a.array) {
		return  false
	}
	a.array[index] = element
	return  true
}
func (a *ArrayList) IndexOf(elem interface{}) int {
	return a.Search(elem)
}
func (a *ArrayList) LastIndexOf(  elem interface{} ) int{
	var size = a.Size()
	for i:= size; i > 0; i--{
		if a.array[size - i - 1] == elem{
			return size - i - 1
		}

	}
	return -1
}

