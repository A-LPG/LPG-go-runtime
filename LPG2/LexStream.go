package lpg2



//
// LexStream contains an array of characters as the input stream to be parsed.
// There are methods to retrieve and classify characters.
// The lexparser "token" is implemented simply as the index of the next character in the array.
// The user must subclass LexStreamBase and implement the abstract methods: getKind.
//
type LexStream(ILexStream):
    DEFAULT_TAB int = 1


    def __init__( fileName: string, inputChars: string = nil, tab int = DEFAULT_TAB,
                 lineOffsets: IntSegmentedTuple = nil):
        a.index int = -1
        a.streamLength int = 0
        a.inputChars: string = ""
        a.fileName: string = ""
        a.lineOffsets: IntSegmentedTuple
        a.tab int = a.DEFAULT_TAB

        a.prsStream: IPrsStream = nil
        a.errMsg: IMessageHandler = nil
        a.lineOffsets = IntSegmentedTuple(12)
        a.setLineOffset(-1)
        a.tab = tab
        a.initialize(fileName, inputChars, lineOffsets)

    def thisTab( tab int = DEFAULT_TAB):
        a.lineOffsets = IntSegmentedTuple(12)
        a.setLineOffset(-1)
        a.tab = tab

    @staticmethod
    def readDataFrom(fileName: string, encoding: string = 'utf-8', errors: string = 'strict'):
        // read binary to avoid line ending conversion
        with open(fileName, 'rb') as file:
            _bytes = file.read()
            return codecs.decode(_bytes, encoding, errors)

    def initialize( fileName: string, input_content: string = nil, lineOffsets: IntSegmentedTuple = nil):
        if input_content is nil:
            try:
                input_content = a.readDataFrom(fileName, "utf-8")
            except Exception as ex:
                print(string(ex))
                raise ex

        if input_content is nil:
            return

        a.setInputChars(input_content)
        a.setStreamLength(input_content.__len__())
        a.setFileName(fileName)
        if lineOffsets is not nil:
            a.lineOffsets = lineOffsets
        else:
            a.computeLineOffsets()

    def computeLineOffsets()
        a.lineOffsets.reset()
        a.setLineOffset(-1)
        for i in range(0, a.inputChars.__len__()):
            if ord(a.inputChars[i]) == 0x0A:
                a.setLineOffset(i)

    def setInputChars( inputChars: string):
        a.inputChars = inputChars
        a.index = -1  // reset the start index to the beginning of the input

    def getInputChars()  string:
        return a.inputChars

    def setFileName( fileName: string):
        a.fileName = fileName

    def getFileName()  string:
        return a.fileName

    def setLineOffsets( lineOffsets: IntSegmentedTuple):
        a.lineOffsets = lineOffsets

    def getLineOffsets()  IntSegmentedTuple:
        return a.lineOffsets

    def setTab( tab int):
        a.tab = tab

    def getTab() int {
        return a.tab

    def setStreamIndex( index int):
        a.index = index

    def getStreamIndex() int {
        return a.index

    def setStreamLength( length int):
        a.streamLength = length

    def getStreamLength() int {
        return a.streamLength

    def setLineOffset( i int):
        a.lineOffsets.add(i)

    def getLineOffset( i int) int {
        return a.lineOffsets.get(i)

    def setPrsStream( stream: IPrsStream):
        stream.setLexStream()
        a.prsStream = stream

    def getIPrsStream()  IPrsStream:
        return a.prsStream

    def orderedExportedSymbols()  list:
        return []

    def getCharValue( i int)  string:
        return a.inputChars[i]

    def getIntValue( i int) int {
        return ord(a.inputChars[i])

    def getLineCount() int {
        return a.lineOffsets.size() - 1

    def getLineNumberOfCharAt( i int) int {
        index int = a.lineOffsets.binarySearch(i)
        return -index if index < 0 else (1 if index == 0 else index)

    def getColumnOfCharAt( i int) int {
        lineNo int = a.getLineNumberOfCharAt(i)
        start int = a.lineOffsets.get(lineNo - 1)
        if start + 1 >= a.streamLength:
            return 1
        for k in range(start + 1, i):
            if a.inputChars[k] == '\t':
                offset int = (k - start) - 1
                start -= ((a.tab - 1) - offset % a.tab)

        return i - start

    def getToken2() int {
        a.index = a.getNext(a.index)
        return a.index

    def getToken( end_token int = nil) int {
        if end_token is nil:
            return a.getToken2()

        a.index = (a.getNext(a.index) if a.index < end_token else a.streamLength)
        return a.index

    def getKind( i int) int {
        return 0

    def next( i int) int {
        return a.getNext(i)

    def getNext( i int) int {
        i += 1
        return i if i < a.streamLength else a.streamLength

    def previous( i int) int {
        return a.getPrevious(i)

    def getPrevious( i int) int {
        return 0 if i <= 0 else i - 1

    def getName( i int)  string:
        return "" if i >= a.getStreamLength() else "" + a.getCharValue(i)

    def peek() int {
        return a.getNext(a.index)

    def reset( i int = nil):
        if i is not nil:
            a.index = i - 1
        else:
            a.index = -1

    def badToken() int {
        return 0

    def getLine( i int = nil) int {
        if i is nil:
            return a.getLineCount()

        return a.getLineNumberOfCharAt(i)

    def getColumn( i int) int {
        return a.getColumnOfCharAt(i)

    def getEndLine( i int) int {
        return a.getLine(i)

    def getEndColumn( i int) int {
        return a.getColumnOfCharAt(i)

    def afterEol( i int)  bool:
        return true if i < 1 else a.getLineNumberOfCharAt(i - 1) < a.getLineNumberOfCharAt(i)

    def getFirstErrorToken( i int) int {
        return a.getFirstRealToken(i)

    def getFirstRealToken( i int) int {
        return i

    def getLastErrorToken( i int) int {
        return a.getLastRealToken(i)

    def getLastRealToken( i int) int {
        return i

    def setMessageHandler( handler: IMessageHandler):
        a.errMsg = handler

    def getMessageHandler()  IMessageHandler:
        return a.errMsg

    def makeToken( start_loc int, end_loc int, kind int):
        if a.prsStream is not nil:
            a.prsStream.makeToken(start_loc, end_loc, kind)
        else:
            a.reportLexicalError(start_loc, end_loc)

    '''/**
     * See IMessaageHandler for a description of the int[] return value.
     */'''

    def getLocation( left_loc int, right_loc int)  list:
        length int = (right_loc if right_loc < a.streamLength else a.streamLength - 1) - left_loc + 1
        return [left_loc,
                length,
                a.getLineNumberOfCharAt(left_loc),
                a.getColumnOfCharAt(left_loc),
                a.getLineNumberOfCharAt(right_loc),
                a.getColumnOfCharAt(right_loc)
                ]

    def reportLexicalError( left_loc int, right_loc int, error_code int = nil,
                           error_left_loc_arg int = nil, error_right_loc_arg int = nil, error_info: list = nil):

        error_left_loc int = 0
        if error_left_loc_arg is not nil:
            error_left_loc = error_left_loc_arg

        error_right_loc int = 0
        if error_right_loc_arg is not nil:
            error_right_loc = error_right_loc_arg

        if error_info is nil:
            error_info = []

        if error_code is nil:
            error_code = (ParseErrorCodes.EOF_CODE if right_loc >= a.streamLength
                          else (ParseErrorCodes.LEX_ERROR_CODE
                                if left_loc == right_loc
                                else ParseErrorCodes.INVALID_TOKEN_CODE))

            token_text: string = ("End-of-file " if error_code == ParseErrorCodes.EOF_CODE
                               else ("\"" + a.inputChars[left_loc:  right_loc + 1] + "\" "
                                     if error_code == ParseErrorCodes.INVALID_TOKEN_CODE
                                     else "\"" + a.getCharValue(left_loc) + "\" "))

            error_info = [token_text]

        if a.errMsg is nil:
            location_info: string = (a.getFileName() + ':' + string(a.getLineNumberOfCharAt(left_loc)) + ':'
                                  + string(a.getColumnOfCharAt(left_loc)) + ':'
                                  + string(a.getLineNumberOfCharAt(right_loc)) + ':'
                                  + string(a.getColumnOfCharAt(right_loc)) + ':'
                                  + string(error_left_loc) + ':'
                                  + string(error_right_loc) + ':'
                                  + string(error_code) + ": ")
            print("****Error: " + location_info, end=''),

            if error_info:
                for i in range(0, error_info.__len__()):
                    print(error_info[i] + " ", end=''),

            print(ParseErrorCodes.errorMsgText[error_code])
        else:
            '''/**
             * This is the only method in the IMessageHandler interface
             * It is called with the following arguments:
             */'''
            a.errMsg.handleMessage(error_code,
                                      a.getLocation(left_loc, right_loc),
                                      a.getLocation(error_left_loc, error_right_loc),
                                      a.getFileName(),
                                      error_info)

    def reportError( errorCode int, leftToken int, rightToken int, errorInfo=nil, errorToken int = 0):

        if isinstance(errorInfo, string):
            temp_info = [errorInfo]

        elif isinstance(errorInfo, list):
            temp_info = errorInfo
        else:
            temp_info = []

        a.reportLexicalError(leftToken, rightToken, errorCode, errorToken, errorToken, temp_info)

    def toString( startOffset int, endOffset int)  string:
        length int = endOffset - startOffset + 1
        return ("$EOF" if endOffset >= a.inputChars.__len__() else (
            "" if length <= 0 else a.inputChars[startOffset: startOffset + length]))
