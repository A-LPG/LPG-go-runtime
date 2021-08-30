package lpg2

import (
    "bytes"
    "io"
    "os"
    "fmt"
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

func NewLexStream(fileName string, inputChars *string, tab int, lineOffsets *IntSegmentedTuple) *LexStream{
    this := new(LexStream)
    this.lineOffsets = NewIntSegmentedTuple(12,4)
    this.setLineOffset(-1)
    this.tab = tab
    this.initialize(fileName, inputChars, lineOffsets)
    return  this
}
func (this *LexStream) GetFileString(fileName string) (*string, error) {

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
func (this *LexStream) initialize(fileName string, inputChars *string, lineOffsets *IntSegmentedTuple)  error {
    if nil == inputChars {
        var str,ex = this.GetFileString(fileName)
        if ex != nil{
            return ex
        }
        inputChars = str
    }

    if nil == inputChars {
        return nil
    }

    this.setInputChars(*inputChars)
    this.setStreamLength(len(*inputChars))
    this.setFileName(fileName)
    if lineOffsets != nil {
        this.lineOffsets = lineOffsets
    } else {
        this.computeLineOffsets()
    }
    return nil
}
func (this *LexStream) computeLineOffsets()  {
    this.lineOffsets.reset()
    this.setLineOffset(-1)
    var i int = 0
    for ;i < len(this.inputChars); i++ {
        if this.inputChars[i] == 0x0A {
            this.setLineOffset(i)
        }
    }
}
func (this *LexStream) setInputChars(inputChars string)  {
    this.inputChars = []rune(inputChars)
    this.index = -1 // reset the start index to the beginning of the input
}
func (this *LexStream) getInputChars() string {
    return string(this.inputChars)
}
func (this *LexStream) setFileName(fileName string)  {
    this.fileName = fileName
}
func (this *LexStream) getFileName() string {
    return this.fileName
}
func (this *LexStream) setLineOffsets(lineOffsets *IntSegmentedTuple)  {
    this.lineOffsets = lineOffsets
}
func (this *LexStream) getLineOffsets() *IntSegmentedTuple {
    return this.lineOffsets
}
func (this *LexStream) setTab(tab int)  {
    this.tab = tab
}
func (this *LexStream) getTab() int {
    return this.tab
}
func (this *LexStream) setStreamIndex(index int)  {
    this.index = index
}
func (this *LexStream) getStreamIndex() int {
    return this.index
}
func (this *LexStream) setStreamLength(streamLength int)  {
    this.streamLength = streamLength
}
func (this *LexStream) getStreamLength() int {
    return this.streamLength
}
func (this *LexStream) setLineOffset(i int)  {
    this.lineOffsets.add(i)
}
func (this *LexStream) getLineOffset(i int) int {
    return this.lineOffsets.get(i)
}
func (this *LexStream) setPrsStream(prsStream IPrsStream)  {
    prsStream.setLexStream(this)
    this.prsStream = prsStream
}
func (this *LexStream) getIPrsStream() IPrsStream{
    return this.prsStream
}

func (this *LexStream) orderedExportedSymbols() []string {
    return nil
}
func (this *LexStream) getCharValue(i int) string {
    return string(this.inputChars[i])
}
func (this *LexStream) getIntValue(i int) int {
    return int(this.inputChars[i])
}

func (this *LexStream) getLineCount() int {
    return this.lineOffsets.size() - 1
}
func (this *LexStream) getLineNumberOfCharAt(i int) int {
    var index int = this.lineOffsets.binarySearch(i)
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
func (this *LexStream) getColumnOfCharAt(i int) int {
    var lineNo int = this.getLineNumberOfCharAt(i)
     var   start int = this.lineOffsets.get(lineNo - 1)
    if start + 1 >= this.streamLength {
        return 1
    }
    var k int = start + 1
    for  ;k < i ;k++ {
        if this.inputChars[k] == '\t' {
            var offset int = (k - start) - 1
            start -= ((this.tab - 1) - offset % this.tab)
        }
    }
    return i - start
}
func (this *LexStream) getToken() int {
    this.index = this.getNext(this.index)
    return this.index
}
func (this *LexStream) getTokenFromEndToken(end_token int ) int {
     if this.index < end_token {
         this.index =this.getNext(this.index)
     }else{
         this.index =this.streamLength
     }
     return  this.index
}
func (this *LexStream) getKind(i int) int {
    return 0
}
func (this *LexStream) next(i int) int {
    return this.getNext(i)
}
func (this *LexStream) getNext(i int) int {
     i+=1
     if i < this.streamLength {
         return  i
     }else{
        return  this.streamLength
     }
}
func (this *LexStream) previous(i int) int {
    return this.getPrevious(i)
}
func (this *LexStream) getPrevious(i int) int {
    if i <= 0 {
        return 0
    }else {
       return i - 1
    }
}
func (this *LexStream) getName(i int) string {
    if i >= this.getStreamLength() {
        return ""
    }else{
        return  "" + this.getCharValue(i)
    }
}
func (this *LexStream) peek() int {
    return this.getNext(this.index)
}
func (this *LexStream) resetTo(i int)  {
    this.index = i - 1
}
func (this *LexStream) reset()  {
    this.index = -1
}

func (this *LexStream) badToken() int {
    return 0
}
func (this *LexStream) getLine(i int) int {
    return this.getLineNumberOfCharAt(i)
}

func (this *LexStream) getColumn(i int) int {
    return this.getColumnOfCharAt(i)
}
func (this *LexStream) getEndLine(i int) int {
    return this.getLine(i)
}
func (this *LexStream) getEndColumn(i int) int {
    return this.getColumnOfCharAt(i)
}
func (this *LexStream) afterEol(i int) bool {
    if i < 1 {
        return  true
    } else{
        return this.getLineNumberOfCharAt(i - 1) < this.getLineNumberOfCharAt(i)
    }
}
func (this *LexStream) getFirstErrorToken(i int) int {
    return this.getFirstRealToken(i)
}
func (this *LexStream) getFirstRealToken(i int) int {
    return i
}
func (this *LexStream) getLastErrorToken(i int) int {
    return this.getLastRealToken(i)
}
func (this *LexStream) getLastRealToken(i int) int {
    return i
}
func (this *LexStream) setMessageHandler(errMsg IMessageHandler)  {
    this.errMsg = errMsg
}
func (this *LexStream) getMessageHandler() IMessageHandler {
    return this.errMsg
}
func (this *LexStream) makeToken(startLoc int, endLoc int, kind int)  {
    if this.prsStream == nil {
        this.prsStream.makeToken(startLoc, endLoc, kind)
    } else {
        this.reportLexicalError(startLoc, endLoc,0,0,0,nil)
    }
}

func (this *LexStream) getLocation(left_loc int, right_loc int) []int {
    var end_loc int
    if right_loc < this.streamLength {
        end_loc = right_loc
    }else{
        end_loc=this.streamLength - 1
    }
    var length int = end_loc - left_loc + 1

    return []int{left_loc, length, this.getLineNumberOfCharAt(left_loc), this.getColumnOfCharAt(left_loc),
        this.getLineNumberOfCharAt(right_loc), this.getColumnOfCharAt(right_loc)}
}
func (this *LexStream) reportLexicalError(left_loc int, right_loc int, errorCode int, error_left_loc int,
    error_right_loc int, errorInfo  []string)  {

    if NIL_CODE == errorCode && errorInfo == nil {

            if right_loc >= this.streamLength {
                errorCode =EOF_CODE
            } else{
                if  left_loc == right_loc{
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
                    tokenText= "\"" + this.toString(left_loc,  right_loc  + 1) + "\" "
                }else{
                    tokenText = "\"" + this.getCharValue(left_loc) + "\" ";
                }
            }
        error_left_loc = 0
        error_right_loc = 0
        errorInfo = []string{tokenText}
    }


    if nil == this.errMsg {
       var locationInfo =fmt.Sprintf("%s : %d : %d : %d : %d : %d : %d : %d",this.getFileName(),this.getLineNumberOfCharAt(left_loc),
            this.getColumnOfCharAt(left_loc),this.getLineNumberOfCharAt(right_loc),this.getColumnOfCharAt(right_loc),
            error_left_loc,error_right_loc,errorCode)

        print("****Error " + locationInfo)
        var i int = 0
            for  ;i < len(errorInfo); i++ {
                print(errorInfo[i] + " ")
            }

        println(errorMsgText[errorCode])
    } else {
        this.errMsg.handleMessage(errorCode, this.getLocation(left_loc, right_loc), this.getLocation(error_left_loc, error_right_loc), this.getFileName(), errorInfo)
    }
}

func (this *LexStream) reportError(errorCode int, leftToken int, rightToken int, errorInfo []string, errorToken int)  {
    this.reportLexicalError(leftToken, rightToken, errorCode,  errorToken,errorToken, errorInfo)
}

func (this *LexStream) toString(startOffset int, endOffset int) string {
    var length int = endOffset - startOffset + 1
    if endOffset >= len(this.inputChars) {
        return "$EOF"
    } else{
        if length <= 0 {
            return ""
        } else{
            return string(this.inputChars[startOffset : startOffset+ length])
        }
    }
}


