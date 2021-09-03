package lpg2

import "fmt"

//
// PrsStream holds an arraylist of tokens "lexed" from the input stream.
//
type PrsStream struct {
     iLexStream ILexStream 
     kindMap []int
     tokens *TokenArrayList
     adjuncts *TokenArrayList
     index int
     len int
}

func NewPrsStream(iLexStream ILexStream) *PrsStream{
	my := new(PrsStream)
	my.index=0
	my.len=0
	if iLexStream != nil {
		my.iLexStream = iLexStream
		iLexStream.SetPrsStream(my)
		my.ReSetTokenStream()
	}
	return my
}
func (my * PrsStream) OrderedExportedSymbols() []string {
	return nil
}
func (my * PrsStream) RemapTerminalSymbols(ordered_parser_symbols []string, eof_symbol int) error {
	// lexStream might be null, maybe only erroneously, but it has happened

	if nil == my.iLexStream {
		return NewNullPointerException("PrsStream.RemapTerminalSymbols(..)  lexStream is nil")
	}

	var ordered_lexer_symbols = my.iLexStream.OrderedExportedSymbols()
	if ordered_lexer_symbols == nil {
		return NewNullExportedSymbolsException("")
	}
	if ordered_parser_symbols == nil {
		return  NewNullTerminalSymbolsException("")
	}
	var unimplemented_symbols  = NewIntArrayList()
	if StringSliceEqual(ordered_lexer_symbols,ordered_parser_symbols) {
		my.kindMap = make([]int,len(ordered_lexer_symbols))
		var terminal_map =make(map[string]int)
		var i int = 0
		for ; i < len(ordered_lexer_symbols); i++{
			terminal_map[ordered_lexer_symbols[i]]= i
		}
		i  = 0
		for ; i < len(ordered_parser_symbols); i++ {

			var k, ok = terminal_map[ordered_parser_symbols[i]]
			if ok {
				my.kindMap[k] = i
			} else {
				if i == eof_symbol {
					return NewUndefinedEofSymbolException("")
				}
				unimplemented_symbols.Add(i)
			}
		}
	}
	if unimplemented_symbols.Size() > 0 {
		return NewUnimplementedTerminalsException(unimplemented_symbols)
	}
	return nil
}
func (my * PrsStream) MapKind(kind int) int {
	if len(my.kindMap) == 0  || my.kindMap == nil{
		return  kind
	}else{
		return  my.kindMap[kind]
	}
}
func (my * PrsStream) ReSetTokenStream()  {
	my.tokens = NewTokenArrayList()
	my.index = 0
	my.adjuncts = NewTokenArrayList()
}
func (my * PrsStream) SetLexStream(lexStream ILexStream)  {
	my.iLexStream = lexStream
	my.ReSetTokenStream()
}
func (my * PrsStream) ResetLexStream(lexStream ILexStream)  {

	my.iLexStream = lexStream
	if lexStream != nil {
		lexStream.SetPrsStream(my)

	}
}


func (my * PrsStream) MakeToken(startLoc int, endLoc int, kind int)  {
	var token  = NewToken( startLoc, endLoc, my.MapKind(kind),my)
	token.SetTokenIndex(my.tokens.Size())
	my.tokens.Add(token)
	token.SetAdjunctIndex(my.adjuncts.Size())
}
func (my * PrsStream) RemoveLastToken()  {
	var last_index int = my.tokens.Size() - 1
	var token  =my.tokens.Get(last_index)
	var adjuncts_size int = my.adjuncts.Size()
	for;adjuncts_size > token.GetAdjunctIndex(); {
		adjuncts_size-=1
		my.adjuncts.RemoveAt(adjuncts_size)
	}
	my.tokens.RemoveAt(last_index)
}
func (my * PrsStream) MakeErrorToken(firsttok int, lasttok int, errortok int, kind int) int {
	var index int = my.tokens.Size() // the next index

	//
	// Note that when creating an error token, we do not remap its kind.
	// Since my is not a lexical operation, it is the responsibility of
	// the calling program (a parser driver) to pass to us the proper kind
	// that it wants for an error token.
	//
	var token  = NewErrorToken( my.GetIToken(firsttok),
								my.GetIToken(lasttok),
								my.GetIToken(errortok),
								my.GetStartOffset(firsttok),
								my.GetEndOffSet(lasttok),
								kind)

	token.SetTokenIndex(my.tokens.Size())
	my.tokens.Add(token)
	token.SetAdjunctIndex(my.adjuncts.Size())

	return index
}
func (my * PrsStream) AddToken(token IToken)  {
	token.SetTokenIndex(my.tokens.Size())
	my.tokens.Add(token)
	token.SetAdjunctIndex(my.adjuncts.Size())
}
func (my * PrsStream) MakeAdjunct(startLoc int, endLoc int, kind int)  {
	var token_index int = my.tokens.Size() - 1 // index of last token processed
	var adjunct  = NewAdjunct(startLoc, endLoc, my.MapKind(kind), my)
	adjunct.SetAdjunctIndex(my.adjuncts.Size())
	adjunct.SetTokenIndex(token_index)
	my.adjuncts.Add(adjunct)
}
func (my * PrsStream) AddAdjunct(adjunct IToken)  {
	var token_index int = my.tokens.Size() - 1 // index of last token processed
	adjunct.SetTokenIndex(token_index)
	adjunct.SetAdjunctIndex(my.adjuncts.Size())
	my.adjuncts.Add(adjunct)
}
func (my * PrsStream) GetTokenText(i int) string {
	var t  = my.tokens.Get(i)
	return t.ToString()
}
func (my * PrsStream) GetStartOffset(i int) int {
	var t =my.tokens.Get(i)
	return t.GetStartOffset()
}
func (my * PrsStream) GetEndOffSet(i int) int {
	var t =my.tokens.Get(i)
	return t.GetEndOffset()
}
func (my * PrsStream) GetTokenLength(i int) int {
	var t =my.tokens.Get(i)
	return t.GetEndOffset() - t.GetStartOffset() + 1
}
func (my * PrsStream) GetLineNumberOfTokenAt(i int) int {
	if nil == my.iLexStream{ 
		return 0
	}
	var t =my.tokens.Get(i)
	return my.iLexStream.GetLineNumberOfCharAt(t.GetStartOffset())
}
func (my * PrsStream) GetEndLineNumberOfTokenAt(i int) int {
	if nil == my.iLexStream{
		return 0
	}
	var t =my.tokens.Get(i)
	return my.iLexStream.GetLineNumberOfCharAt(t.GetEndOffset())
}
func (my * PrsStream) GetColumnOfTokenAt(i int) int {
	if nil == my.iLexStream{
		return 0
	}
	var t =my.tokens.Get(i)
	return my.iLexStream.GetColumnOfCharAt(t.GetStartOffset())
}
func (my * PrsStream) GetEndColumnOfTokenAt(i int) int {
	if nil == my.iLexStream{
		return 0
	}
	var t =my.tokens.Get(i)
	return my.iLexStream.GetColumnOfCharAt(t.GetEndOffset())
}
func (my * PrsStream) OrderedTerminalSymbols() []string {
	return nil
}
func (my * PrsStream) GetLineOffSet(i int) int {
	if nil == my.iLexStream{
		return 0
	}
	return my.iLexStream.GetLineOffSet(i)
}
func (my * PrsStream) GetLineCount() int {
	if nil == my.iLexStream{
		return 0
	}
	return my.iLexStream.GetLineCount()
}
func (my * PrsStream) GetLineNumberOfCharAt(i int) int {
	if nil == my.iLexStream{
		return 0
	}
	return my.iLexStream.GetLineNumberOfCharAt(i)
}
func (my * PrsStream) GetColumnOfCharAt(i int) int {
	return my.GetColumnOfCharAt(i)
}
func (my * PrsStream) GetFirstErrorToken(i int) int {
	return my.GetFirstRealToken(i)
}
func (my * PrsStream) GetFirstRealToken(i int) int {
	for;i >= my.len; {
		var temp = my.tokens.Get(i)
		var errorToken,_ = temp.(*ErrorToken)
		i = errorToken.GetFirstRealToken().GetTokenIndex()
	}
	return i
}
func (my * PrsStream) GetLastErrorToken(i int) int {
	return my.GetLastRealToken(i)
}
func (my * PrsStream) GetLastRealToken(i int) int {
	for;i >= my.len; {
		var temp = my.tokens.Get(i)
		var errorToken,_ = temp.(*ErrorToken)
		i = errorToken.GetLastRealToken().GetTokenIndex()
	}
	return i
}
func (my * PrsStream) GetInputChars() []rune {
	if nil == my.iLexStream{
		return nil
	}
	return my.iLexStream.GetInputChars()
}

func (my * PrsStream) ToStringFromIndex(first_token int, last_token int) string {
	return my.ToString(my.tokens.Get(first_token), my.tokens.Get(last_token))
}
func (my * PrsStream) ToString(t1 IToken, t2 IToken) string {
	if nil == my.iLexStream{
		return ""
	}
	return my.iLexStream.ToString(t1.GetStartOffset(), t2.GetEndOffset())
}
func (my * PrsStream) GetSize() int {
	return my.tokens.Size()
}
func (my * PrsStream) SetSize()  {
	my.len = my.tokens.Size()
}
func (my * PrsStream) GetTokenIndexAtCharacter(offSet int) int {
	var low int = 0
	var high int = my.tokens.Size()
	for;high > low; {
		var mid int = int((high + low) / 2)
		var mid_element =my.tokens.Get(mid)
		if offSet >= mid_element.GetStartOffset() && offSet <= mid_element.GetEndOffset() {
			return mid
		} else {
			if offSet < mid_element.GetStartOffset() {
				high = mid
			} else {
				low = mid + 1
			}
		}
	}
	return -(low - 1)
}
func (my * PrsStream) GetTokenAtCharacter(offSet int) IToken {
	var tokenIndex int = my.GetTokenIndexAtCharacter(offSet)
	if tokenIndex < 0 {
		return nil
	}else{
		return my.GetTokenAt(tokenIndex)
	}
}
func (my * PrsStream) GetTokenAt(i int) IToken {
	return my.tokens.Get(i)
}
func (my * PrsStream) GetIToken(i int) IToken {
	return my.tokens.Get(i)
}
func (my * PrsStream) GetTokens() *TokenArrayList {
	return my.tokens
}
func (my * PrsStream) GetStreamIndex() int {
	return my.index
}
func (my * PrsStream) GetStreamLength() int {
	return my.len
}
func (my * PrsStream) SetStreamIndex(index int)  {
	my.index = index
}
func (my * PrsStream) ReSetStreamLength()  {
	my.len = my.tokens.Size()
}
func (my * PrsStream) SetStreamLength(len int)  {

	my.len = len
}
func (my * PrsStream) GetILexStream() ILexStream {
	return my.iLexStream
}
func (my * PrsStream) GetLexStream() ILexStream {
	return my.iLexStream
}
func (my * PrsStream) DumpTokens()  {
	if my.GetSize() <= 2 {
		return
	}
	println(" Kind \tOffSet \tLen \tLine \tCol \tText\n")

	var i int = 1
	for ;i < my.GetSize() - 1 ;i++{
		my.DumpToken(i)
	}
}
func (my * PrsStream) DumpToken(i int)  {
	fmt.Printf(" ( %d )",my.GetKind(i))
	fmt.Printf(" \t%d" , my.GetStartOffset(i))
	fmt.Printf(" \t%d" , my.GetTokenLength(i))
	fmt.Printf(" \t%d" , my.GetLineNumberOfTokenAt(i))
	fmt.Printf(" \t%d" , my.GetColumnOfTokenAt(i))
	fmt.Printf(" \t%s" , my.GetTokenText(i))
	fmt.Printf("\n")
}
func (my * PrsStream) GetAdjunctsFromIndex(i int) []IToken {
	var start_index int = (my.tokens.Get(i)).GetAdjunctIndex()
	var end_index int
		if i + 1 == my.tokens.Size(){
			end_index = my.adjuncts.Size()
		}else{
			end_index =my.tokens.Get(my.GetNext(i)).GetAdjunctIndex()
		}

	var size int = end_index - start_index
	var slice []IToken = make([]IToken,size)
	var j int = start_index
	var k int = 0
	for ;j < end_index; {
		slice[k] = my.adjuncts.Get(j)
		j++
		k++
	}
	return slice
}
func (my * PrsStream) GetFollowingAdjuncts(i int) []IToken {
	return my.GetAdjunctsFromIndex(i)
}
func (my * PrsStream) GetPrecedingAdjuncts(i int) []IToken {
	return my.GetAdjunctsFromIndex(my.GetPrevious(i))
}
func (my * PrsStream) GetAdjuncts() *TokenArrayList {
	return my.adjuncts
}
func (my * PrsStream) GetToken() int {
	my.index = my.GetNext(my.index)
	return my.index
}
func (my * PrsStream) GetTokenFromEndToken(end_token int ) int {
	if my.index < end_token {
		my.index = my.GetNext(my.index)
	}else{
		my.index =my.len - 1
	}
	return  my.index
}
func (my * PrsStream) GetKind(i int) int {
	var t =my.tokens.Get(i)
	return t.GetKind()
}
func (my * PrsStream) GetNext(i int) int {
	i+=1
	if i < my.len {
		return  i
	}else{
		return my.len - 1
	}

}
func (my * PrsStream) GetPrevious(i int) int {
	if i <= 0 {
		return 0
	}else{
		return i - 1
	}
}
func (my * PrsStream) GetName(i int) string {
	return my.GetTokenText(i)
}
func (my * PrsStream) Peek() int {
	return my.GetNext(my.index)
}
func (my * PrsStream)   Reset() {
	my.index = 0
}
func (my * PrsStream)  ResetTo(i  int) {
	my.index = my.GetPrevious(i)
}

func (my * PrsStream) BadToken() int {
	return 0
}
func (my * PrsStream) GetLine(i int) int {
	return my.GetLineNumberOfTokenAt(i)
}
func (my * PrsStream) GetColumn(i int) int {
	return my.GetColumnOfTokenAt(i)
}
func (my * PrsStream) GetEndLine(i int) int {
	return my.GetEndLineNumberOfTokenAt(i)
}
func (my * PrsStream) GetEndColumn(i int) int {
	return my.GetEndColumnOfTokenAt(i)
}
func (my * PrsStream) AfterEol(i int) bool {
	if  i < 1 {
		return  true
	} else{
		return my.GetEndLineNumberOfTokenAt(i - 1) < my.GetLineNumberOfTokenAt(i)
	}
}
func (my * PrsStream) GetFileName() string {
	if nil == my.iLexStream{
		return ""
	}
	return my.iLexStream.GetFileName()
}

//
// Here is where we report errors.  The default method is simply to fmt.Printf the error message to the console.
// However, the user may supply an error message handler to process error messages.  To support that
// a message handler interface is provided that has a single method HandleMessage().  The user has his
// error message handler class implement the IMessageHandler interface and provides an object of my type
// to the runtime using the SetMessageHandler(errorMsg) method. If the message handler object is Set,
// the ReportError methods will invoke its HandleMessage() method.
//
// IMessageHandler errMsg = null // the error message handler object is declared in LexStream
//
func (my * PrsStream) SetMessageHandler(errMsg IMessageHandler)  {
	my.iLexStream.SetMessageHandler(errMsg)
}
func (my * PrsStream) GetMessageHandler() IMessageHandler  {
	return my.iLexStream.GetMessageHandler()
}

func (my * PrsStream) ReportError(errorCode int, leftToken int, rightToken int, errorInfo []string, errorToken int)  {
	my.iLexStream.ReportLexicalError( my.GetStartOffset(leftToken),
										my.GetEndOffSet(rightToken),
										errorCode,
										my.GetStartOffset(errorToken),
										my.GetEndOffSet(errorToken),
										errorInfo)
}


