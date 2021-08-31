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
	self := new(PrsStream)
	self.index=0
	self.len=0
	if iLexStream != nil {
		self.iLexStream = iLexStream
		iLexStream.setPrsStream(self)
		self.resetTokenStream()
	}
	return self
}
func (self * PrsStream) orderedExportedSymbols() []string {
	return nil
}
func (self * PrsStream) remapTerminalSymbols(ordered_parser_symbols []string, eof_symbol int) error {
	// lexStream might be null, maybe only erroneously, but it has happened

	if nil == self.iLexStream {
		return NewNullPointerException("PrsStream.remapTerminalSymbols(..)  lexStream is nil")
	}

	var ordered_lexer_symbols = self.iLexStream.orderedExportedSymbols()
	if ordered_lexer_symbols == nil {
		return NewNullExportedSymbolsException("")
	}
	if ordered_parser_symbols == nil {
		return  NewNullTerminalSymbolsException("")
	}
	var unimplemented_symbols  = NewIntArrayList()
	if StringSliceEqual(ordered_lexer_symbols,ordered_parser_symbols) {
		self.kindMap = make([]int,len(ordered_lexer_symbols))
		var terminal_map =make(map[string]int)
		var i int = 0
		for ; i < len(ordered_lexer_symbols); i++{
			terminal_map[ordered_lexer_symbols[i]]= i
		}
		i  = 0
		for ; i < len(ordered_parser_symbols); i++ {

			var k, ok = terminal_map[ordered_parser_symbols[i]]
			if ok {
				self.kindMap[k] = i
			} else {
				if i == eof_symbol {
					return NewUndefinedEofSymbolException("")
				}
				unimplemented_symbols.add(i)
			}
		}
	}
	if unimplemented_symbols.size() > 0 {
		return NewUnimplementedTerminalsException(unimplemented_symbols)
	}
	return nil
}
func (self * PrsStream) mapKind(kind int) int {
	if len(self.kindMap) == 0  || self.kindMap == nil{
		return  kind
	}else{
		return  self.kindMap[kind]
	}
}
func (self * PrsStream) resetTokenStream()  {
	self.tokens = NewTokenArrayList()
	self.index = 0
	self.adjuncts = NewTokenArrayList()
}
func (self * PrsStream) setLexStream(lexStream ILexStream)  {
	self.iLexStream = lexStream
	self.resetTokenStream()
}
func (self * PrsStream) resetLexStream(lexStream ILexStream)  {

	self.iLexStream = lexStream
	if lexStream != nil {
		lexStream.setPrsStream(self)

	}
}


func (self * PrsStream) makeToken(startLoc int, endLoc int, kind int)  {
	var token  = NewToken( startLoc, endLoc, self.mapKind(kind),self)
	token.setTokenIndex(self.tokens.size())
	self.tokens.add(token)
	token.setAdjunctIndex(self.adjuncts.size())
}
func (self * PrsStream) removeLastToken()  {
	var last_index int = self.tokens.size() - 1
	var token  =self.tokens.get(last_index)
	var adjuncts_size int = self.adjuncts.size()
	for;adjuncts_size > token.getAdjunctIndex(); {
		adjuncts_size-=1
		self.adjuncts.removeAt(adjuncts_size)
	}
	self.tokens.removeAt(last_index)
}
func (self * PrsStream) makeErrorToken(firsttok int, lasttok int, errortok int, kind int) int {
	var index int = self.tokens.size() // the next index

	//
	// Note that when creating an error token, we do not remap its kind.
	// Since self is not a lexical operation, it is the responsibility of
	// the calling program (a parser driver) to pass to us the proper kind
	// that it wants for an error token.
	//
	var token  = NewErrorToken( self.getIToken(firsttok),
								self.getIToken(lasttok),
								self.getIToken(errortok),
								self.getStartOffset(firsttok),
								self.getEndOffset(lasttok),
								kind)

	token.setTokenIndex(self.tokens.size())
	self.tokens.add(token)
	token.setAdjunctIndex(self.adjuncts.size())

	return index
}
func (self * PrsStream) addToken(token IToken)  {
	token.setTokenIndex(self.tokens.size())
	self.tokens.add(token)
	token.setAdjunctIndex(self.adjuncts.size())
}
func (self * PrsStream) makeAdjunct(startLoc int, endLoc int, kind int)  {
	var token_index int = self.tokens.size() - 1// index of last token processed
	var adjunct  = NewAdjunct(startLoc, endLoc, self.mapKind(kind), self)
	adjunct.setAdjunctIndex(self.adjuncts.size())
	adjunct.setTokenIndex(token_index)
	self.adjuncts.add(adjunct)
}
func (self * PrsStream) addAdjunct(adjunct IToken)  {
	var token_index int = self.tokens.size() - 1// index of last token processed
	adjunct.setTokenIndex(token_index)
	adjunct.setAdjunctIndex(self.adjuncts.size())
	self.adjuncts.add(adjunct)
}
func (self * PrsStream) getTokenText(i int) string {
	var t  = self.tokens.get(i)
	return t.toString()
}
func (self * PrsStream) getStartOffset(i int) int {
	var t =self.tokens.get(i)
	return t.getStartOffset()
}
func (self * PrsStream) getEndOffset(i int) int {
	var t =self.tokens.get(i)
	return t.getEndOffset()
}
func (self * PrsStream) getTokenLength(i int) int {
	var t =self.tokens.get(i)
	return t.getEndOffset() - t.getStartOffset() + 1
}
func (self * PrsStream) getLineNumberOfTokenAt(i int) int {
	if nil == self.iLexStream{ 
		return 0
	}
	var t =self.tokens.get(i)
	return self.iLexStream.getLineNumberOfCharAt(t.getStartOffset())
}
func (self * PrsStream) getEndLineNumberOfTokenAt(i int) int {
	if nil == self.iLexStream{
		return 0
	}
	var t =self.tokens.get(i)
	return self.iLexStream.getLineNumberOfCharAt(t.getEndOffset())
}
func (self * PrsStream) getColumnOfTokenAt(i int) int {
	if nil == self.iLexStream{
		return 0
	}
	var t =self.tokens.get(i)
	return self.iLexStream.getColumnOfCharAt(t.getStartOffset())
}
func (self * PrsStream) getEndColumnOfTokenAt(i int) int {
	if nil == self.iLexStream{
		return 0
	}
	var t =self.tokens.get(i)
	return self.iLexStream.getColumnOfCharAt(t.getEndOffset())
}
func (self * PrsStream) orderedTerminalSymbols() []string {
	return nil
}
func (self * PrsStream) getLineOffset(i int) int {
	if nil == self.iLexStream{
		return 0
	}
	return self.iLexStream.getLineOffset(i)
}
func (self * PrsStream) getLineCount() int {
	if nil == self.iLexStream{
		return 0
	}
	return self.iLexStream.getLineCount()
}
func (self * PrsStream) getLineNumberOfCharAt(i int) int {
	if nil == self.iLexStream{
		return 0
	}
	return self.iLexStream.getLineNumberOfCharAt(i)
}
func (self * PrsStream) getColumnOfCharAt(i int) int {
	return self.getColumnOfCharAt(i)
}
func (self * PrsStream) getFirstErrorToken(i int) int {
	return self.getFirstRealToken(i)
}
func (self * PrsStream) getFirstRealToken(i int) int {
	for;i >= self.len; {
		var temp = self.tokens.get(i)
		var errorToken,_ = temp.(*ErrorToken)
		i = errorToken.getFirstRealToken().getTokenIndex()
	}
	return i
}
func (self * PrsStream) getLastErrorToken(i int) int {
	return self.getLastRealToken(i)
}
func (self * PrsStream) getLastRealToken(i int) int {
	for;i >= self.len; {
		var temp = self.tokens.get(i)
		var errorToken,_ = temp.(*ErrorToken)
		i = errorToken.getLastRealToken().getTokenIndex()
	}
	return i
}
func (self * PrsStream) getInputChars() string {
	if nil == self.iLexStream{
		return ""
	}
	return self.iLexStream.getInputChars()
}

func (self * PrsStream) toStringFromIndex(first_token int, last_token int) string {
	return self.toString(self.tokens.get(first_token), self.tokens.get(last_token))
}
func (self * PrsStream) toString(t1 IToken, t2 IToken) string {
	if nil == self.iLexStream{
		return ""
	}
	return self.iLexStream.toString(t1.getStartOffset(), t2.getEndOffset())
}
func (self * PrsStream) getSize() int {
	return self.tokens.size()
}
func (self * PrsStream) setSize()  {
	self.len = self.tokens.size()
}
func (self * PrsStream) getTokenIndexAtCharacter(offset int) int {
	var low int = 0
	var high int = self.tokens.size()
	for;high > low; {
		var mid int = int((high + low) / 2)
		var mid_element =self.tokens.get(mid)
		if offset >= mid_element.getStartOffset() && offset <= mid_element.getEndOffset() {
			return mid
		} else {
			if offset < mid_element.getStartOffset() {
				high = mid
			} else {
				low = mid + 1
			}
		}
	}
	return -(low - 1)
}
func (self * PrsStream) getTokenAtCharacter(offset int) IToken {
	var tokenIndex int = self.getTokenIndexAtCharacter(offset)
	if tokenIndex < 0 {
		return nil
	}else{
		return self.getTokenAt(tokenIndex)
	}
}
func (self * PrsStream) getTokenAt(i int) IToken {
	return self.tokens.get(i)
}
func (self * PrsStream) getIToken(i int) IToken {
	return self.tokens.get(i)
}
func (self * PrsStream) getTokens() *TokenArrayList {
	return self.tokens
}
func (self * PrsStream) getStreamIndex() int {
	return self.index
}
func (self * PrsStream) getStreamLength() int {
	return self.len
}
func (self * PrsStream) setStreamIndex(index int)  {
	self.index = index
}
func (self * PrsStream) resetStreamLength()  {
	self.len = self.tokens.size()
}
func (self * PrsStream) setStreamLength(len int)  {

	self.len = len
}
func (self * PrsStream) getILexStream() ILexStream {
	return self.iLexStream
}
func (self * PrsStream) getLexStream() ILexStream {
	return self.iLexStream
}
func (self * PrsStream) dumpTokens()  {
	if self.getSize() <= 2 {
		return
	}
	println(" Kind \tOffset \tLen \tLine \tCol \tText\n")

	var i int = 1
	for ;i < self.getSize() - 1 ;i++{
		self.dumpToken(i)
	}
}
func (self * PrsStream) dumpToken(i int)  {
	fmt.Printf(" ( %d )",self.getKind(i))
	fmt.Printf(" \t%d" , self.getStartOffset(i))
	fmt.Printf(" \t%d" , self.getTokenLength(i))
	fmt.Printf(" \t%d" , self.getLineNumberOfTokenAt(i))
	fmt.Printf(" \t%d" , self.getColumnOfTokenAt(i))
	fmt.Printf(" \t%d" , self.getTokenText(i))
	fmt.Printf("\n")
}
func (self * PrsStream) getAdjunctsFromIndex(i int) []IToken {
	var start_index int = (self.tokens.get(i)).getAdjunctIndex()
	var end_index int
		if i + 1 == self.tokens.size(){
			end_index = self.adjuncts.size()
		}else{
			end_index =self.tokens.get(self.getNext(i)).getAdjunctIndex()
		}

	var size int = end_index - start_index
	var slice []IToken = make([]IToken,size)
	var j int = start_index
	var k int = 0
	for ;j < end_index; {
		slice[k] = self.adjuncts.get(j)
		j++
		k++
	}
	return slice
}
func (self * PrsStream) getFollowingAdjuncts(i int) []IToken {
	return self.getAdjunctsFromIndex(i)
}
func (self * PrsStream) getPrecedingAdjuncts(i int) []IToken {
	return self.getAdjunctsFromIndex(self.getPrevious(i))
}
func (self * PrsStream) getAdjuncts() *TokenArrayList {
	return self.adjuncts
}
func (self * PrsStream) getToken() int {
	self.index = self.getNext(self.index)
	return self.index
}
func (self * PrsStream) getTokenFromEndToken(end_token int ) int {
	if self.index < end_token {
		self.index = self.getNext(self.index)
	}else{
		self.index =self.len - 1
	}
	return  self.index
}
func (self * PrsStream) getKind(i int) int {
	var t =self.tokens.get(i)
	return t.getKind()
}
func (self * PrsStream) getNext(i int) int {
	i+=1
	if i < self.len {
		return  i
	}else{
		return self.len - 1
	}

}
func (self * PrsStream) getPrevious(i int) int {
	if i <= 0 {
		return 0
	}else{
		return i - 1
	}
}
func (self * PrsStream) getName(i int) string {
	return self.getTokenText(i)
}
func (self * PrsStream) peek() int {
	return self.getNext(self.index)
}
func (self * PrsStream)   reset() {
	self.index = 0
}
func (self * PrsStream)  resetTo(i  int) {
	self.index = self.getPrevious(i)
}

func (self * PrsStream) badToken() int {
	return 0
}
func (self * PrsStream) getLine(i int) int {
	return self.getLineNumberOfTokenAt(i)
}
func (self * PrsStream) getColumn(i int) int {
	return self.getColumnOfTokenAt(i)
}
func (self * PrsStream) getEndLine(i int) int {
	return self.getEndLineNumberOfTokenAt(i)
}
func (self * PrsStream) getEndColumn(i int) int {
	return self.getEndColumnOfTokenAt(i)
}
func (self * PrsStream) afterEol(i int) bool {
	if  i < 1 {
		return  true
	} else{
		return self.getEndLineNumberOfTokenAt(i - 1) < self.getLineNumberOfTokenAt(i)
	}
}
func (self * PrsStream) getFileName() string {
	if nil == self.iLexStream{
		return ""
	}
	return self.iLexStream.getFileName()
}

//
// Here is where we report errors.  The default method is simply to fmt.Printf the error message to the console.
// However, the user may supply an error message handler to process error messages.  To support that
// a message handler interface is provided that has a single method handleMessage().  The user has his
// error message handler class implement the IMessageHandler interface and provides an object of self type
// to the runtime using the setMessageHandler(errorMsg) method. If the message handler object is set,
// the reportError methods will invoke its handleMessage() method.
//
// IMessageHandler errMsg = null // the error message handler object is declared in LexStream
//
func (self * PrsStream) setMessageHandler(errMsg IMessageHandler)  {
	self.iLexStream.setMessageHandler(errMsg)
}
func (self * PrsStream) getMessageHandler() IMessageHandler  {
	return self.iLexStream.getMessageHandler()
}

func (self * PrsStream) reportError(errorCode int, leftToken int, rightToken int, errorInfo []string, errorToken int)  {
	self.iLexStream.reportLexicalError( self.getStartOffset(leftToken),
										self.getEndOffset(rightToken),
										errorCode,
										self.getStartOffset(errorToken),
										self.getEndOffset(errorToken),
										errorInfo)
}


