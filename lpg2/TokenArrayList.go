package lpg2
type TokenArrayList struct {
	array []IToken
}
func NewTokenArrayListFrom(array []IToken) *TokenArrayList {
	return &TokenArrayList{
		array: array,
	}
}
func NewTokenArrayListFromCopy(array []IToken) *TokenArrayList {
	newArray := make([]IToken, len(array))
	copy(newArray, array)
	return &TokenArrayList{
		array: newArray,
	}
}
func NewTokenArrayListSize(size int, cap int) *TokenArrayList {
	return &TokenArrayList{
		array: make([]IToken, size, cap),
	}
}
func  NewTokenArrayList() *TokenArrayList {
	return NewTokenArrayListSize(0,0)
}
func (a *TokenArrayList) Clone() *TokenArrayList {
	array := make([]IToken, len(a.array))
	copy(array, a.array)
	return NewTokenArrayListFrom(array)
}

func (a *TokenArrayList) Clear() bool{
	if len(a.array) > 0 {
		a.array = make([]IToken, 0)
	}
	return  true
}
func (a *TokenArrayList) RemoveAt(index int)  (value IToken, found bool){
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
func (a *TokenArrayList) Remove(value IToken) bool{
	if i := a.Search(value); i != -1 {
		_, found := a.RemoveAt(i)
		return found
	}
	return false
}
func (a *TokenArrayList) Search(value IToken) int {

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

func (a *TokenArrayList) RemoveAll() bool{
	return a.Clear()
}

func (a *TokenArrayList) ToArray() []IToken {
	array := make([]IToken, len(a.array))
	copy(array, a.array)
	return array
}

func (a *TokenArrayList) Size() int {
	return len(a.array)
}
func (a *TokenArrayList) Add(elem IToken) *TokenArrayList{
	a.array = append(a.array, elem)
	return a
}
func (a *TokenArrayList) Get(index int) IToken{
	if index < 0 || index >= len(a.array) {
		return nil
	}
	return a.array[index]
}
func (a *TokenArrayList) At(index int) (value IToken) {
	return a.Get(index)
}
func (a *TokenArrayList) Contains( val IToken) bool{
	return a.Search(val) != -1
}
func (a *TokenArrayList) IsEmpty() bool{
	return a.Size() == 0
}
func (a *TokenArrayList) Set(index int, element IToken) bool {
	if index < 0 || index >= len(a.array) {
		return  false
	}
	a.array[index] = element
	return  true
}
func (a *TokenArrayList) IndexOf(elem IToken) int {
	return a.Search(elem)
}
func (a *TokenArrayList) LastIndexOf(  elem IToken ) int{
	var size = a.Size()
	for i:= size; i > 0; i--{
		if a.array[size - i - 1] == elem{
			return size - i - 1
		}

	}
	return -1
}

