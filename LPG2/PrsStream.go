package lpg2


//
// PrsStream holds an arraylist of tokens "lexed" from the input stream.
//
type PrsStream(IPrsStream):


    def __init__( lex: ILexStream = nil):
        a.iLexStream: ILexStream = nil
        a.kindMap: list = []
        a.tokens: ArrayList = ArrayList()
        a.adjuncts: ArrayList = ArrayList()
        a.index int = 0
        a.len int = 0
        if lex is not nil:
            a.iLexStream = lex
            lex.setPrsStream()
            a.resetTokenStream()

    def orderedExportedSymbols()  list:
        return []

    def remapTerminalSymbols( ordered_parser_symbols: list, eof_symbol int):
        // lexStream might be nil, maybe only erroneously, but it has happened
        if a.iLexStream is nil:
            raise ReferenceError("PrsStream.remapTerminalSymbols(..):  lexStream is nil")

        ordered_lexer_symbols: list = a.iLexStream.orderedExportedSymbols()
        if ordered_lexer_symbols is nil:
            raise NullTerminalSymbolsException()

        if ordered_parser_symbols is nil:
            raise NullTerminalSymbolsException()

        unimplemented_symbols: ArrayList = ArrayList()
        if ordered_lexer_symbols != ordered_parser_symbols:
            a.kindMap = [0] * (ordered_lexer_symbols.__len__())
            terminal_map = dict()
            for i in range(0, ordered_lexer_symbols.__len__()):
                terminal_map[ordered_lexer_symbols[i]] = i

            for i in range(0, ordered_parser_symbols.__len__()):
                k int = terminal_map.get(ordered_parser_symbols[i], nil)
                if k is not nil:
                    a.kindMap[k] = i
                else:
                    if i == eof_symbol:
                        raise UndefinedEofSymbolException()

                    unimplemented_symbols.add(i)

        if unimplemented_symbols.size() > 0:
            raise UnimplementedTerminalsException(unimplemented_symbols)

    def mapKind( kind int) int {
        return kind if a.kindMap is nil or a.kindMap.__len__() == 0 else a.kindMap[kind]

    def resetTokenStream()
        a.tokens = ArrayList()
        a.index = 0
        a.adjuncts = ArrayList()

    def setLexStream( lexStream: ILexStream):
        a.iLexStream = lexStream
        a.resetTokenStream()

    def resetLexStream( lexStream: ILexStream):

        if lexStream is not nil:
            lexStream.setPrsStream()
            a.iLexStream = lexStream

    def makeToken( startLoc int, endLoc int, kind int):
        token: Token = Token(startLoc, endLoc, a.mapKind(kind), a)
        token.setTokenIndex(a.tokens.size())
        a.tokens.add(token)
        token.setAdjunctIndex(a.adjuncts.size())

    def removeLastToken()
        last_index int = a.tokens.size() - 1
        token: Token = a.tokens.get(last_index)
        adjuncts_size int = a.adjuncts.size()
        while adjuncts_size > token.getAdjunctIndex()
            adjuncts_size -= 1
            a.adjuncts.remove(adjuncts_size)

        a.tokens.remove(last_index)

    def makeErrorToken( firsttok int, lasttok int, errortok int, kind int) int {
        index int = a.tokens.size()  // the next index

        //
        // Note that when creating an error token, we do not remap its kind.
        // Since a is not a lexical operation, it is the responsibility of
        // the calling program (a parser driver): to pass to us the proper kind
        // that it wants for an error token.
        //
        token: Token = ErrorToken(a.getIToken(firsttok),
                                  a.getIToken(lasttok),
                                  a.getIToken(errortok),
                                  a.getStartOffset(firsttok),
                                  a.getEndOffset(lasttok),
                                  kind)

        token.setTokenIndex(a.tokens.size())
        a.tokens.add(token)
        token.setAdjunctIndex(a.adjuncts.size())

        return index

    def addToken( token: IToken):
        token.setTokenIndex(a.tokens.size())
        a.tokens.add(token)
        token.setAdjunctIndex(a.adjuncts.size())

    def makeAdjunct( startLoc int, endLoc int, kind int):
        token_index int = a.tokens.size() - 1  // index of last token processed
        adjunct: Adjunct = Adjunct(startLoc, endLoc, a.mapKind(kind), a)
        adjunct.setAdjunctIndex(a.adjuncts.size())
        adjunct.setTokenIndex(token_index)
        a.adjuncts.add(adjunct)

    def addAdjunct( adjunct: IToken):
        token_index int = a.tokens.size() - 1  // index of last token processed
        adjunct.setTokenIndex(token_index)
        adjunct.setAdjunctIndex(a.adjuncts.size())
        a.adjuncts.add(adjunct)

    def getTokenText( i int)  string:
        t: IToken = a.tokens.get(i)
        return t.toString()

    def getStartOffset( i int) int {
        t: IToken = a.tokens.get(i)
        return t.getStartOffset()

    def getEndOffset( i int) int {
        t: IToken = a.tokens.get(i)
        return t.getEndOffset()

    def getTokenLength( i int) int {
        t: IToken = a.tokens.get(i)
        return t.getEndOffset() - t.getStartOffset() + 1

    def getLineNumberOfTokenAt( i int) int {
        if not a.iLexStream:
            return 0
        t: IToken = a.tokens.get(i)
        return a.iLexStream.getLineNumberOfCharAt(t.getStartOffset())

    def getEndLineNumberOfTokenAt( i int) int {
        if not a.iLexStream: return 0
        t: IToken = a.tokens.get(i)
        return a.iLexStream.getLineNumberOfCharAt(t.getEndOffset())

    def getColumnOfTokenAt( i int) int {
        if not a.iLexStream: return 0
        t: IToken = a.tokens.get(i)
        return a.iLexStream.getColumnOfCharAt(t.getStartOffset())

    def getEndColumnOfTokenAt( i int) int {
        if not a.iLexStream: return 0
        t: IToken = a.tokens.get(i)
        return a.iLexStream.getColumnOfCharAt(t.getEndOffset())

    def orderedTerminalSymbols()  list:
        return []

    def getLineOffset( i int) int {
        if not a.iLexStream: return 0
        return a.iLexStream.getLineOffset(i)

    def getLineCount() int {
        if not a.iLexStream: return 0
        return a.iLexStream.getLineCount()

    def getLineNumberOfCharAt( i int) int {
        if not a.iLexStream: return 0
        return a.iLexStream.getLineNumberOfCharAt(i)

    def getColumnOfCharAt( i int) int {
        return a.getColumnOfCharAt(i)

    def getFirstErrorToken( i int) int {
        return a.getFirstRealToken(i)

    def getFirstRealToken( i int) int {
        while i >= a.len:
            i = (a.tokens.get(i)).getFirstRealToken().getTokenIndex()

        return i

    def getLastErrorToken( i int) int {
        return a.getLastRealToken(i)

    def getLastRealToken( i int) int {
        while i >= a.len:
            i = (a.tokens.get(i)).getLastRealToken().getTokenIndex()

        return i

    def getInputChars()  string:
        return a.iLexStream.getInputChars() if isinstance(a.iLexStream, LexStream) else ""

    def getInputBytes()
        //  return (a.iLexStream instanceof Utf8LexStream ? (<Utf8LexStream>a.iLexStream).getInputBytes() : nil):
        return ""

    def toStringFromIndex( first_token int, last_token int)  string:
        return a.toString(a.tokens.get(first_token), a.tokens.get(last_token))

    def toString( t1: IToken, t2: IToken)  string:
        if not a.iLexStream: return ""
        return a.iLexStream.toString(t1.getStartOffset(), t2.getEndOffset())

    def getSize() int {
        return a.tokens.size()

    def setSize()
        a.len = a.tokens.size()

    def getTokenIndexAtCharacter( offset int) int {
        low int = 0
        high int = a.tokens.size()
        while high > low:
            mid int = (high + low) // 2
            mid_element: IToken = a.tokens.get(mid)
            if offset >= mid_element.getStartOffset() and offset <= mid_element.getEndOffset()
                return mid
            else:
                if offset < mid_element.getStartOffset()
                    high = mid
                else:
                    low = mid + 1

        return -(low - 1)

    def getTokenAtCharacter( offset int)  IToken:
        token_index int = a.getTokenIndexAtCharacter(offset)
        return nil if (token_index < 0) else a.getTokenAt(token_index)

    def getTokenAt( i int)  IToken:
        return a.tokens.get(i)

    def getIToken( i int)  IToken:
        return a.tokens.get(i)

    def getTokens()  ArrayList:
        return a.tokens

    def getStreamIndex() int {
        return a.index

    def getStreamLength() int {
        return a.len

    def setStreamIndex( index int):
        a.index = index

    def setStreamLength2()
        a.len = a.tokens.size()

    def setStreamLength( length int = nil):
        if length is nil:
            a.setStreamLength2()
            return

        a.len = length

    def getILexStream()  ILexStream:
        return a.iLexStream

    def getLexStream()  ILexStream:
        return a.iLexStream

    def dumpTokens()
        if a.getSize() <= 2:
            return

        print(" Kind \tOffset \tLen \tLine \tCol \tText\n")
        for i in range(1, a.getSize() - 1):
            a.dumpToken(i)

    def dumpToken( i int):
        print(" (" + string(a.getKind(i)) + "):", end='')
        print(" \t" + string(a.getStartOffset(i)), end='')
        print(" \t" + string(a.getTokenLength(i)), end='')
        print(" \t" + string(a.getLineNumberOfTokenAt(i)), end='')
        print(" \t" + string(a.getColumnOfTokenAt(i)), end='')
        print(" \t" + string(a.getTokenText(i)))


    def getAdjunctsFromIndex( i int)  list:
        start_index int = (a.tokens.get(i)).getAdjunctIndex()
        end_index int = (a.adjuncts.size()
                          if (i + 1 == a.tokens.size())
                          else (a.tokens.get(a.getNext(i))).getAdjunctIndex())
        size int = end_index - start_index
        token_slice: list = [nil] * size
        j int = start_index
        k int = 0
        while j < end_index:
            token_slice[k] = a.adjuncts.get(j)
            k += 1
            j += 1

        return token_slice

    //
    // Return an iterator for the adjuncts that follow token i.
    //
    def getFollowingAdjuncts( i int)  list:
        return a.getAdjunctsFromIndex(i)

    def getPrecedingAdjuncts( i int)  list:
        return a.getAdjunctsFromIndex(a.getPrevious(i))

    def getAdjuncts()  ArrayList:
        return a.adjuncts

    def getToken2() int {
        a.index = a.getNext(a.index)
        return a.index

    def getToken( end_token int = nil) int {
        if end_token is nil:
            return a.getToken2()
        a.index = (a.getNext(a.index) if a.index < end_token else a.len - 1)
        return a.index

    def getKind( i int) int {
        t: IToken = a.tokens.get(i)
        return t.getKind()

    def getNext( i int) int {
        i += 1
        return i if i < a.len else a.len - 1

    def getPrevious( i int) int {
        return 0 if i <= 0 else i - 1

    def getName( i int)  string:
        return a.getTokenText(i)

    def peek() int {
        return a.getNext(a.index)

    def reset1()
        a.index = 0

    def reset2( i int):
        a.index = a.getPrevious(i)

    def reset( i int = nil):
        if i is nil:
            a.reset1()
        else:
            a.reset2(i)

    def badToken() int {
        return 0

    def getLine( i int) int {
        return a.getLineNumberOfTokenAt(i)

    def getColumn( i int) int {
        return a.getColumnOfTokenAt(i)

    def getEndLine( i int) int {
        return a.getEndLineNumberOfTokenAt(i)

    def getEndColumn( i int) int {
        return a.getEndColumnOfTokenAt(i)

    def afterEol( i int)  bool:
        return true if i < 1 else a.getEndLineNumberOfTokenAt(i - 1) < a.getLineNumberOfTokenAt(i)

    def getFileName()  string:
        if not a.iLexStream:
            return ""
        return a.iLexStream.getFileName()

    //
    // Here is where we report errors.  The default method is simply to print the error message to the console.
    // However, the user may supply an error message handler to process error messages.  To support that
    // a message handler type is provided that has a single method handleMessage().  The user has his
    // error message handler type implement the IMessageHandler type and provides an object of a type
    // to the runtime using the setMessageHandler(errorMsg): method. If the message handler object is set,
    // the reportError methods will invoke its handleMessage() method.
    //
    // IMessageHandler errMsg = nil // the error message handler object is declared in LexStream
    //
    def setMessageHandler( handler: IMessageHandler = nil):
        a.iLexStream.setMessageHandler(handler)

    def getMessageHandler()  IMessageHandler:
        return a.iLexStream.getMessageHandler()

    def reportError( errorCode int, leftToken int, rightToken int, errorInfo=nil, errorToken int = 0):
        temp_info: list
        if isinstance(errorInfo, string):
            temp_info = [errorInfo]
        elif isinstance(errorInfo, list):
            temp_info = errorInfo
        else:
            temp_info = []

        a.iLexStream.reportLexicalError(a.getStartOffset(leftToken), a.getEndOffset(rightToken),
                                           errorCode, a.getStartOffset(errorToken), a.getEndOffset(errorToken),
                                           temp_info)
