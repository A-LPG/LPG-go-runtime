package lpg2

import (
    "bytes"
    "fmt"
    "io"
    "os"
)

//
// LexStream contains an array of characters as the input stream to be parsed.
// There are methods to retrieve and classify characters.
// The lexparser "token" is implemented simply as the index of the next character in the array.
// The user must subclass LexStreamBase and implement the abstract methods getKind.
//

type LexStream struct {
     DEFAULT_TAB int
     index int
     streamLength int
     inputChars  []rune

     fileName string
     lineOffsets *IntSegmentedTuple
     tab int 
     prsStream IPrsStream
     errMsg IMessageHandler
}



func NewLexStream(fileName string, inputChars *string, tab int, lineOffsets *IntSegmentedTuple) (*LexStream,error){
    self := new(LexStream)
    self.DEFAULT_TAB  = 1
    self.index  = -1
    self.streamLength = 0
    self.tab = self.DEFAULT_TAB
    self.lineOffsets = NewIntSegmentedTuple(12,4)
    self.setLineOffset(-1)

    err := self.initialize(fileName, inputChars, lineOffsets)
    if err != nil {
        return nil,err
    }
    return  self, nil
}
func (self *LexStream) GetFileString(fileName string) (*string, error) {

    buf := bytes.NewBuffer(nil)

    f, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    _, err = io.Copy(buf, f)
    if err != nil {
        return nil, err
    }

    s := string(buf.Bytes())

    return &s, nil

}
func (self *LexStream) initialize(fileName string, inputChars *string, lineOffsets *IntSegmentedTuple)  error {
    if nil == inputChars {
        var str,ex = self.GetFileString(fileName)
        if ex != nil{
            return ex
        }
        inputChars = str
    }

    if nil == inputChars {
        return nil
    }

    self.setInputChars(*inputChars)
    self.setStreamLength(len(self.inputChars))
    self.setFileName(fileName)
    if lineOffsets != nil {
        self.lineOffsets = lineOffsets
    } else {
        self.computeLineOffsets()
    }
    return nil
}
func (self *LexStream) computeLineOffsets()  {
    self.lineOffsets.reset()
    self.setLineOffset(-1)
    var i int = 0
    var size = len(self.inputChars)
    for ;i < size; i++ {
        if self.inputChars[i] == 0x0A {
            self.setLineOffset(i)
        }
    }
}
func (self *LexStream) setInputChars(inputChars string)  {
    self.inputChars = []rune(inputChars)
    self.index = -1 // reset the start index to the beginning of the input
}
func (self *LexStream) getInputChars() string {
    return string(self.inputChars)
}
func (self *LexStream) setFileName(fileName string)  {
    self.fileName = fileName
}
func (self *LexStream) getFileName() string {
    return self.fileName
}
func (self *LexStream) setLineOffsets(lineOffsets *IntSegmentedTuple)  {
    self.lineOffsets = lineOffsets
}
func (self *LexStream) getLineOffsets() *IntSegmentedTuple {
    return self.lineOffsets
}
func (self *LexStream) setTab(tab int)  {
    self.tab = tab
}
func (self *LexStream) getTab() int {
    return self.tab
}
func (self *LexStream) setStreamIndex(index int)  {
    self.index = index
}
func (self *LexStream) getStreamIndex() int {
    return self.index
}
func (self *LexStream) setStreamLength(streamLength int)  {
    self.streamLength = streamLength
}
func (self *LexStream) getStreamLength() int {
    return self.streamLength
}
func (self *LexStream) setLineOffset(i int)  {
    self.lineOffsets.add(i)
}
func (self *LexStream) getLineOffset(i int) int {
    return self.lineOffsets.get(i)
}
func (self *LexStream) setPrsStream(prsStream IPrsStream)  {
    if nil == prsStream{
        return
    }
    prsStream.setLexStream(self)
    self.prsStream = prsStream
}
func (self *LexStream) getIPrsStream() IPrsStream{
    return self.prsStream
}

func (self *LexStream) orderedExportedSymbols() []string {
    return nil
}
func (self *LexStream) getCharValue(i int) string {
    return string(self.inputChars[i])
}
func (self *LexStream) getIntValue(i int) int {
    return int(self.inputChars[i])
}

func (self *LexStream) getLineCount() int {
    return self.lineOffsets.size() - 1
}
func (self *LexStream) getLineNumberOfCharAt(i int) int {
    var index int = self.lineOffsets.binarySearch(i)
    if  index < 0 {
        return  -index
    } else{
        if index == 0 {
            return 1
        } else{
            return index
        }
    }
}
func (self *LexStream) getColumnOfCharAt(i int) int {
    var lineNo int = self.getLineNumberOfCharAt(i)
    var start int = self.lineOffsets.get(lineNo - 1)
    if start + 1 >= self.streamLength {
        return 1
    }
    var k int = start + 1
    for  ;k < i ;k++ {
        if self.inputChars[k] == '\t' {
            var offset int = (k - start) - 1
            start -= ((self.tab - 1) - offset % self.tab)
        }
    }
    return i - start
}
func (self *LexStream) getToken() int {
    self.index = self.getNext(self.index)
    return self.index
}
func (self *LexStream) getTokenFromEndToken(end_token int ) int {
     if self.index < end_token {
         self.index =self.getNext(self.index)
     }else{
         self.index =self.streamLength
     }
     return  self.index
}
func (self *LexStream) getKind(i int) int {
    return 0
}
func (self *LexStream) next(i int) int {
    return self.getNext(i)
}
func (self *LexStream) getNext(i int) int {
     i+=1
     if i < self.streamLength {
         return  i
     }else{
        return  self.streamLength
     }
}
func (self *LexStream) previous(i int) int {
    return self.getPrevious(i)
}
func (self *LexStream) getPrevious(i int) int {
    if i <= 0 {
        return 0
    }else {
       return i - 1
    }
}
func (self *LexStream) getName(i int) string {
    if i >= self.getStreamLength() {
        return ""
    }else{
        return  "" + self.getCharValue(i)
    }
}
func (self *LexStream) peek() int {
    return self.getNext(self.index)
}
func (self *LexStream) resetTo(i int)  {
    self.index = i - 1
}
func (self *LexStream) reset()  {
    self.index = -1
}

func (self *LexStream) badToken() int {
    return 0
}
func (self *LexStream) getLine(i int) int {
    return self.getLineNumberOfCharAt(i)
}

func (self *LexStream) getColumn(i int) int {
    return self.getColumnOfCharAt(i)
}
func (self *LexStream) getEndLine(i int) int {
    return self.getLine(i)
}
func (self *LexStream) getEndColumn(i int) int {
    return self.getColumnOfCharAt(i)
}
func (self *LexStream) afterEol(i int) bool {
    if i < 1 {
        return  true
    } else{
        return self.getLineNumberOfCharAt(i - 1) < self.getLineNumberOfCharAt(i)
    }
}
func (self *LexStream) getFirstErrorToken(i int) int {
    return self.getFirstRealToken(i)
}
func (self *LexStream) getFirstRealToken(i int) int {
    return i
}
func (self *LexStream) getLastErrorToken(i int) int {
    return self.getLastRealToken(i)
}
func (self *LexStream) getLastRealToken(i int) int {
    return i
}

//
// Here is where we report errors.  The default method is simply to print the error message to the console.
// However, the user may supply an error message handler to process error messages.  To support that
// a message handler interface is provided that has a single method handleMessage().  The user has his
// error message handler class implement the IMessageHandler interface and provides an object of this type
// to the runtime using the setMessageHandler(errorMsg) method. If the message handler object is set,
// the reportError methods will invoke its handleMessage() method.
//
func (self *LexStream) setMessageHandler(errMsg IMessageHandler)  {
    self.errMsg = errMsg
}
func (self *LexStream) getMessageHandler() IMessageHandler {
    return self.errMsg
}
func (self *LexStream) makeToken(startLoc int, endLoc int, kind int)  {
    if self.prsStream == nil {
        self.prsStream.makeToken(startLoc, endLoc, kind)
    } else {
        self.reportLexicalErrorPosition(startLoc, endLoc)// make it a lexical error
    }
}

func (self *LexStream) getLocation(leftLoc int, rightLoc int) []int {
    var endLoc int
    if rightLoc < self.streamLength {
        endLoc = rightLoc
    }else{
        endLoc =self.streamLength - 1
    }
    var length int = endLoc - leftLoc + 1

    return []int{leftLoc,
                 length,
                 self.getLineNumberOfCharAt(leftLoc),
                 self.getColumnOfCharAt(leftLoc),
                 self.getLineNumberOfCharAt(rightLoc),
                 self.getColumnOfCharAt(rightLoc)}
}
func (self *LexStream) reportLexicalErrorPosition(leftLoc int, rightLoc int) {

        var errorCode int
        if rightLoc >= self.streamLength {
            errorCode =EOF_CODE
        } else{
            if  leftLoc == rightLoc {
                errorCode=LEX_ERROR_CODE
            }else{
                errorCode=INVALID_TOKEN_CODE
            }
        }
        var tokenText string
        if errorCode == EOF_CODE {
            tokenText = "End-of-file "
        }else{
            if   errorCode == INVALID_TOKEN_CODE{
                tokenText= "\"" + self.toString(leftLoc,  rightLoc+ 1) + "\" "
            }else{
                tokenText = "\"" + self.getCharValue(leftLoc) + "\" ";
            }
        }
        var errorLeftLoc = 0
        var errorRightLoc = 0
        var errorInfo = []string{tokenText}
        self.reportLexicalError(errorCode, leftLoc, rightLoc, errorLeftLoc, errorRightLoc, errorInfo);

}
func (self *LexStream) reportLexicalError(leftLoc int, rightLoc int, errorCode int, errorLeftLoc int,
    errorRightLoc int, errorInfo  []string)  {

    if nil == self.errMsg {
       var locationInfo =fmt.Sprintf("%s : %d : %d : %d : %d : %d : %d : %d",
            self.getFileName(),
            self.getLineNumberOfCharAt(leftLoc),
            self.getColumnOfCharAt(leftLoc),
            self.getLineNumberOfCharAt(rightLoc),
            self.getColumnOfCharAt(rightLoc),
            errorLeftLoc,
            errorRightLoc,
            errorCode)

        print("****Error " + locationInfo)
        var i int = 0
        for ;i < len(errorInfo); i++ {
            print(errorInfo[i] + " ")
        }

        println(errorMsgText[errorCode])
    } else {
        /**
         * This is the only method in the IMessageHandler interface
         * It is called with the following arguments:
         */
        self.errMsg.handleMessage(errorCode,
                                  self.getLocation(leftLoc, rightLoc),
                                  self.getLocation(errorLeftLoc, errorRightLoc),
                                  self.getFileName(),
                                  errorInfo)
    }
}

func (self *LexStream) reportError(errorCode int, leftToken int, rightToken int, errorInfo []string, errorToken int)  {
    self.reportLexicalError(leftToken, rightToken, errorCode,  errorToken,errorToken, errorInfo)
}

func (self *LexStream) toString(startOffset int, endOffset int) string {
    var length int = endOffset - startOffset + 1
    if endOffset >= len(self.inputChars) {
        return "$EOF"
    } else{
        if length <= 0 {
            return ""
        } else{
            return string(self.inputChars[startOffset :endOffset])
        }
    }
}


