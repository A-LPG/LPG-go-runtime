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
func (a *TokenArrayList) clone() *TokenArrayList {
	array := make([]IToken, len(a.array))
	copy(array, a.array)
	return NewTokenArrayListFrom(array)
}

func (a *TokenArrayList) clear() bool{
	if len(a.array) > 0 {
		a.array = make([]IToken, 0)
	}
	return  true
}
func (a *TokenArrayList) removeAt(index int)  (value IToken, found bool){
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
func (a *TokenArrayList) remove(value IToken) bool{
	if i := a.search(value); i != -1 {
		_, found := a.removeAt(i)
		return found
	}
	return false
}
func (a *TokenArrayList) search(value IToken) int {

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

func (a *TokenArrayList) removeAll() bool{
	return a.clear()
}

func (a *TokenArrayList) toArray() []IToken {
	array := make([]IToken, len(a.array))
	copy(array, a.array)
	return array
}

func (a *TokenArrayList) size() int {
	return len(a.array)
}
func (a *TokenArrayList) add(elem IToken) *TokenArrayList{
	a.array = append(a.array, elem)
	return a
}
func (a *TokenArrayList) get(index int) IToken{
	if index < 0 || index >= len(a.array) {
		return nil
	}
	return a.array[index]
}
func (a *TokenArrayList) at(index int) (value IToken) {
	return a.get(index)
}
func (a *TokenArrayList) contains( val IToken) bool{
	return a.search(val) != -1
}
func (a *TokenArrayList) isEmpty() bool{
	return a.size() == 0
}
func (a *TokenArrayList) set(index int, element IToken) bool {
	if index < 0 || index >= len(a.array) {
		return  false
	}
	a.array[index] = element
	return  true
}
func (a *TokenArrayList) indexOf(elem IToken) int {
	return a.search(elem)
}
func (a *TokenArrayList) lastIndexOf(  elem IToken ) int{
	var size = a.size()
	for i:= size; i > 0; i--{
		if a.array[size - i - 1] == elem{
			return size - i - 1
		}

	}
	return -1
}

