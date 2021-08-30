package lpg2

//
// PrsStream holds an arraylist of tokens "lexed" from the input stream.
//
type PrsStream struct {
  
     iLexStream ILexStream 
     kindMap []int
     tokens TokenArrayList 
     adjuncts TokenArrayList 
     index int
     len int
}

func NewPrsStream(iLexStream ILexStream) *PrsStream{
	this := new(PrsStream)
	if (iLexStream != nil) {
		this.iLexStream = iLexStream
		iLexStream.setPrsStream(this)
		this.resetTokenStream()
	}
	return this
}
func (this * PrsStream) orderedExportedSymbols() []string {
	return nil
}
func (this * PrsStream) remapTerminalSymbols(ordered_parser_symbols []string, eof_symbol int)  {
	// lexStream might be null, maybe only erroneously, but it has happened
	if (!this.iLexStream || this.iLexStream instanceof EscapeStrictPropertyInitializationLexStream ) {
		throw new ReferenceError("PrsStream.remapTerminalSymbols(..)  lexStream is undefined")
	}

	var ordered_lexer_symbols []string = this.iLexStream.orderedExportedSymbols()
	if (ordered_lexer_symbols == undefined) {
		throw new NullTerminalSymbolsException()
	}
	if (ordered_parser_symbols == undefined) {
		throw new NullTerminalSymbolsException()
	}
	var unimplemented_symbols Lpg.Util.ArrayList<int> = new Lpg.Util.ArrayList<int>()
	if (ordered_lexer_symbols != ordered_parser_symbols) {
		this.kindMap = new []int(ordered_lexer_symbols.length)
		var terminal_map =new  Map<string, int>()
		for (var i int = 0 i < ordered_lexer_symbols.length i++) {
			terminal_map.set(ordered_lexer_symbols[i], (i))
		}
		var i int = 0
		for ; i < ordered_parser_symbols.length; i++ {
			var k int = <int>terminal_map.get(ordered_parser_symbols[i])
			if (k != undefined) {
				this.kindMap[k] = i
			} else {
				if (i == eof_symbol) {
					throw new UndefinedEofSymbolException()
				}
				unimplemented_symbols.add(i)
			}
		}
	}
	if (unimplemented_symbols.size() > 0) {
		throw new UnimplementedTerminalsException(unimplemented_symbols)
	}
}
func (this * PrsStream) mapKind(kind int) int {
	return (this.kindMap.length == 0 ? kind  this.kindMap[kind])
}
func (this * PrsStream) resetTokenStream()  {
	this.tokens = new TokenArrayList()
	this.index = 0
	this.adjuncts = new TokenArrayList()
}
func (this * PrsStream) setLexStream(lexStream ILexStream)  {
	this.iLexStream = lexStream
	this.resetTokenStream()
}
func (this * PrsStream) resetLexStream(lexStream ILexStream)  {

	if (lexStream) {
		lexStream.setPrsStream(this)
		this.iLexStream = lexStream
	}
}


func (this * PrsStream) makeToken(startLoc int, endLoc int, kind int)  {
	var token Token = new Token( startLoc, endLoc, this.mapKind(kind),this)
	token.setTokenIndex(this.tokens.size())
	this.tokens.add(token)
	token.setAdjunctIndex(this.adjuncts.size())
}
func (this * PrsStream) removeLastToken()  {
	var last_index int = this.tokens.size() - 1
	var token Token = <Token>this.tokens.get(last_index)
	var adjuncts_size int = this.adjuncts.size()
	while (adjuncts_size > token.getAdjunctIndex()) {
		this.adjuncts.remove(--adjuncts_size)
	}
	this.tokens.remove(last_index)
}
func (this * PrsStream) makeErrorToken(firsttok int, lasttok int, errortok int, kind int) int {
	var index int = this.tokens.size() // the next index

	//
	// Note that when creating an error token, we do not remap its kind.
	// Since this is not a lexical operation, it is the responsibility of
	// the calling program (a parser driver) to pass to us the proper kind
	// that it wants for an error token.
	//
	var token Token = new ErrorToken(  this.getIToken(firsttok),
										this.getIToken(lasttok),
										this.getIToken(errortok),
										this.getStartOffset(firsttok),
										this.getEndOffset(lasttok),
										kind)

	token.setTokenIndex(this.tokens.size())
	this.tokens.add(token)
	token.setAdjunctIndex(this.adjuncts.size())

	return index
}
func (this * PrsStream) addToken(token IToken)  {
	token.setTokenIndex(this.tokens.size())
	this.tokens.add(token)
	token.setAdjunctIndex(this.adjuncts.size())
}
func (this * PrsStream) makeAdjunct(startLoc int, endLoc int, kind int)  {
	var token_index int = this.tokens.size() - 1// index of last token processed
	var adjunct Adjunct = new Adjunct(startLoc, endLoc, this.mapKind(kind), this)
	adjunct.setAdjunctIndex(this.adjuncts.size())
	adjunct.setTokenIndex(token_index)
	this.adjuncts.add(adjunct)
}
func (this * PrsStream) addAdjunct(adjunct IToken)  {
	var token_index int = this.tokens.size() - 1// index of last token processed
	adjunct.setTokenIndex(token_index)
	adjunct.setAdjunctIndex(this.adjuncts.size())
	this.adjuncts.add(adjunct)
}
func (this * PrsStream) getTokenText(i int) string {
	var t IToken = <IToken>this.tokens.get(i)
	return t.toString()
}
func (this * PrsStream) getStartOffset(i int) int {
	var t IToken = <IToken>this.tokens.get(i)
	return t.getStartOffset()
}
func (this * PrsStream) getEndOffset(i int) int {
	var t IToken = <IToken>this.tokens.get(i)
	return t.getEndOffset()
}
func (this * PrsStream) getTokenLength(i int) int {
	var t IToken = <IToken>this.tokens.get(i)
	return t.getEndOffset() - t.getStartOffset() + 1
}
func (this * PrsStream) getLineintOfTokenAt(i int) int {
	if (!this.iLexStream) return 0
	var t IToken = <IToken>this.tokens.get(i)
	return this.iLexStream?.getLineintOfCharAt(t.getStartOffset())
}
func (this * PrsStream) getEndLineintOfTokenAt(i int) int {
	if (!this.iLexStream) return 0
	var t IToken = <IToken>this.tokens.get(i)
	return this.iLexStream?.getLineintOfCharAt(t.getEndOffset())
}
func (this * PrsStream) getColumnOfTokenAt(i int) int {
	if (!this.iLexStream) return 0
	var t IToken = <IToken>this.tokens.get(i)
	return this.iLexStream?.getColumnOfCharAt(t.getStartOffset())
}
func (this * PrsStream) getEndColumnOfTokenAt(i int) int {
	if (!this.iLexStream) return 0
	var t IToken = <IToken>this.tokens.get(i)
	return this.iLexStream?.getColumnOfCharAt(t.getEndOffset())
}
func (this * PrsStream) orderedTerminalSymbols() []string {
	return []
}
func (this * PrsStream) getLineOffset(i int) int {
	if (!this.iLexStream) return 0
	return this.iLexStream?.getLineOffset(i)
}
func (this * PrsStream) getLineCount() int {
	if (!this.iLexStream) return 0
	return this.iLexStream?.getLineCount()
}
func (this * PrsStream) getLineintOfCharAt(i int) int {
	if (!this.iLexStream) return 0
	return this.iLexStream?.getLineintOfCharAt(i)
}
func (this * PrsStream) getColumnOfCharAt(i int) int {
	return this.getColumnOfCharAt(i)
}
func (this * PrsStream) getFirstErrorToken(i int) int {
	return this.getFirstRealToken(i)
}
func (this * PrsStream) getFirstRealToken(i int) int {
	while (i >= this.len) {
		i = (<ErrorToken>this.tokens.get(i)).getFirstRealToken().getTokenIndex()
	}
	return i
}
func (this * PrsStream) getLastErrorToken(i int) int {
	return this.getLastRealToken(i)
}
func (this * PrsStream) getLastRealToken(i int) int {
	while (i >= this.len) {
		i = (<ErrorToken>this.tokens.get(i)).getLastRealToken().getTokenIndex()
	}
	return i
}
func (this * PrsStream) getInputChars() string {
	return (this.iLexStream instanceof LexStream ? (<LexStream>this.iLexStream).getInputChars()  "")
}

func (this * PrsStream) getInputBytes() Int8Array {
	//  return (this.iLexStream instanceof Utf8LexStream ? (<Utf8LexStream>this.iLexStream).getInputBytes()  undefined)
	return new Int8Array(0)
}
func (this * PrsStream) toStringFromIndex(first_token int, last_token int) string {
	return this.toString(<IToken>this.tokens.get(first_token), <IToken>this.tokens.get(last_token))
}
func (this * PrsStream) toString(t1 IToken, t2 IToken) string {
	if (!this.iLexStream) return ""
	return this.iLexStream?.toString(t1.getStartOffset(), t2.getEndOffset())
}
func (this * PrsStream) getSize() int {
	return this.tokens.size()
}
func (this * PrsStream) setSize()  {
	this.len = this.tokens.size()
}
func (this * PrsStream) getTokenIndexAtCharacter(offset int) int {
	var low int = 0, high int = this.tokens.size()
	while (high > low) {
		var mid int = Math.floor((high + low) / 2)
		var mid_element IToken = <IToken>this.tokens.get(mid)
		if (offset >= mid_element.getStartOffset() && offset <= mid_element.getEndOffset()) {
			return mid
		} else {
			if (offset < mid_element.getStartOffset()) {
				high = mid
			} else {
				low = mid + 1
			}
		}
	}
	return -(low - 1)
}
func (this * PrsStream) getTokenAtCharacter(offset int) IToken | undefined {
	var tokenIndex int = this.getTokenIndexAtCharacter(offset)
	return (tokenIndex < 0) ? undefined  this.getTokenAt(tokenIndex)
}
func (this * PrsStream) getTokenAt(i int) IToken {
	return <IToken>this.tokens.get(i)
}
func (this * PrsStream) getIToken(i int) IToken {
	return <IToken>this.tokens.get(i)
}
func (this * PrsStream) getTokens() TokenArrayList {
	return this.tokens
}
func (this * PrsStream) getStreamIndex() int {
	return this.index
}
func (this * PrsStream) getStreamLength() int {
	return this.len
}
func (this * PrsStream) setStreamIndex(index int)  {
	this.index = index
}
func (this * PrsStream) setStreamLength2()  {
	this.len = this.tokens.size()
}
func (this * PrsStream) setStreamLength(len? int)  {
	if (typeof len == 'undefined') {
		this.setStreamLength2()
		return
	}
	this.len = len
}
func (this * PrsStream) getILexStream() ILexStream {
	return this.iLexStream
}
func (this * PrsStream) getLexStream() ILexStream {
	return this.iLexStream
}
func (this * PrsStream) dumpTokens()  {
	if (this.getSize() <= 2) {
		return
	}
	Lpg.Lang.System.Out.println(" Kind \tOffset \tLen \tLine \tCol \tText\n")
	for (var i int = 1 i < this.getSize() - 1 i++) {
		this.dumpToken(i)
	}
}
func (this * PrsStream) dumpToken(i int)  {
	console.log(" (" + this.getKind(i) + ")")
	console.log(" \t" + this.getStartOffset(i))
	console.log(" \t" + this.getTokenLength(i))
	console.log(" \t" + this.getLineintOfTokenAt(i))
	console.log(" \t" + this.getColumnOfTokenAt(i))
	console.log(" \t" + this.getTokenText(i))
	console.log("\n")
}
	getAdjunctsFromIndex(i int) IToken[] {
	var start_index int = (<IToken>this.tokens.get(i)).getAdjunctIndex(),
		end_index int = (i + 1 == this.tokens.size()
									? this.adjuncts.size()
										(<IToken>this.tokens.get(this.getNext(i))).getAdjunctIndex()),
		size int = end_index - start_index
	var slice IToken[] = new Array<IToken>(size)
	for (var j int = start_index, k int = 0 j < end_index j++, k++) {
		slice[k] = <IToken>this.adjuncts.get(j)
	}
	return slice
}
func (this * PrsStream) getFollowingAdjuncts(i int) IToken[] {
	return this.getAdjunctsFromIndex(i)
}
func (this * PrsStream) getPrecedingAdjuncts(i int) IToken[] {
	return this.getAdjunctsFromIndex(this.getPrevious(i))
}
func (this * PrsStream) getAdjuncts() TokenArrayList {
	return this.adjuncts
}
func (this * PrsStream) getToken2() int {
	this.index = this.getNext(this.index)
	return this.index
}
func (this * PrsStream) getToken(end_token? int ) int {
	if (!end_token) {
		return this.getToken2()
	}
	return this.index = (this.index < end_token ? this.getNext(this.index)  this.len - 1)
}
func (this * PrsStream) getKind(i int) int {
	var t IToken = <IToken>this.tokens.get(i)
	return t.getKind()
}
func (this * PrsStream) getNext(i int) int {
	return (++i < this.len ? i  this.len - 1)
}
func (this * PrsStream) getPrevious(i int) int {
	return (i <= 0 ? 0  i - 1)
}
func (this * PrsStream) getName(i int) string {
	return this.getTokenText(i)
}
func (this * PrsStream) peek() int {
	return this.getNext(this.index)
}
func (this * PrsStream)   reset1()  
{
	this.index = 0
}
func (this * PrsStream)   reset2(i  int)  
{
	this.index = this.getPrevious(i)
}
func (this * PrsStream) reset(i? int)  {
	if (!i) 
	{
		this.reset1()
	}
	else{
		this.reset2(i)
	}
}

func (this * PrsStream) badToken() int {
	return 0
}
func (this * PrsStream) getLine(i int) int {
	return this.getLineintOfTokenAt(i)
}
func (this * PrsStream) getColumn(i int) int {
	return this.getColumnOfTokenAt(i)
}
func (this * PrsStream) getEndLine(i int) int {
	return this.getEndLineintOfTokenAt(i)
}
func (this * PrsStream) getEndColumn(i int) int {
	return this.getEndColumnOfTokenAt(i)
}
func (this * PrsStream) afterEol(i int) boolean {
	return (i < 1 ? true  this.getEndLineintOfTokenAt(i - 1) < this.getLineintOfTokenAt(i))
}
func (this * PrsStream) getFileName() string {
	if (!this.iLexStream) return""
	return this.iLexStream?.getFileName()
}

//
// Here is where we report errors.  The default method is simply to print the error message to the console.
// However, the user may supply an error message handler to process error messages.  To support that
// a message handler interface is provided that has a single method handleMessage().  The user has his
// error message handler class implement the IMessageHandler interface and provides an object of this type
// to the runtime using the setMessageHandler(errorMsg) method. If the message handler object is set,
// the reportError methods will invoke its handleMessage() method.
//
// IMessageHandler errMsg = null // the error message handler object is declared in LexStream
//
func (this * PrsStream) setMessageHandler(errMsg IMessageHandler)  {
	this.iLexStream?.setMessageHandler(errMsg)
}
func (this * PrsStream) getMessageHandler() IMessageHandler | undefined {
	return this.iLexStream?.getMessageHandler()
}

func (this * PrsStream) reportError(errorCode int, leftToken int, rightToken int, errorInfo string | []string, errorToken int = 0)  {
	var tempInfo []string
	if (typeof errorInfo == "string") {
		tempInfo = [errorInfo]
	}
	else if (Array.isArray(errorInfo)) {
		tempInfo = errorInfo
	}
	else {
		tempInfo = []
	}
	this.iLexStream?.reportLexicalError(this.getStartOffset(leftToken), this.getEndOffset(rightToken),errorCode, this.getStartOffset(errorToken), this.getEndOffset(errorToken), tempInfo)
}


