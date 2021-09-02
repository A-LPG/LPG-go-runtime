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
// The user must subclass LexStreamBase and implement the abstract methods GetKind.
//

type LexStream struct {
     DEFAULT_TAB int
     index int
     streamLength int
     inputChars  []rune

     fileName string
     lineOffSets *IntSegmentedTuple
     tab int 
     prsStream IPrsStream
     errMsg IMessageHandler
}



func NewLexStream(fileName string, inputChars []rune, tab int, lineOffSets *IntSegmentedTuple) (*LexStream,error){
    my := new(LexStream)
    my.DEFAULT_TAB  = 1
    my.index  = -1
    my.streamLength = 0
    my.tab = my.DEFAULT_TAB
    my.lineOffSets = NewIntSegmentedTuple(12,4)
    my.SetLineOffSet(-1)

    err := my.initialize(fileName, inputChars, lineOffSets)
    if err != nil {
        return nil,err
    }
    return  my, nil
}
func (my *LexStream) GetFileString(fileName string) (*string, error) {

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
func (my *LexStream) initialize(fileName string, inputChars []rune, lineOffSets *IntSegmentedTuple)  error {
    if nil == inputChars {
        var str,ex = my.GetFileString(fileName)
        if ex != nil{
            return ex
        }
        inputChars = []rune(*str)
    }

    if nil == inputChars {
        return nil
    }

    my.SetInputChars(inputChars)
    my.SetStreamLength(len(my.inputChars))
    my.SetFileName(fileName)
    if lineOffSets != nil {
        my.lineOffSets = lineOffSets
    } else {
        my.ComputeLineOffSets()
    }
    return nil
}
func (my *LexStream) ComputeLineOffSets()  {
    my.lineOffSets.ReSet()
    my.SetLineOffSet(-1)
    var i int = 0
    var size = len(my.inputChars)
    for ;i < size; i++ {
        if my.inputChars[i] == 0x0A {
            my.SetLineOffSet(i)
        }
    }
}
func (my *LexStream) SetInputChars(inputChars []rune)  {
    my.inputChars = inputChars
    my.index = -1 // Reset the start index to the beginning of the input
}
func (my *LexStream) GetInputChars() []rune {
    return my.inputChars
}
func (my *LexStream) SetFileName(fileName string)  {
    my.fileName = fileName
}
func (my *LexStream) GetFileName() string {
    return my.fileName
}
func (my *LexStream) SetLineOffSets(lineOffSets *IntSegmentedTuple)  {
    my.lineOffSets = lineOffSets
}
func (my *LexStream) GetLineOffSets() *IntSegmentedTuple {
    return my.lineOffSets
}
func (my *LexStream) SetTab(tab int)  {
    my.tab = tab
}
func (my *LexStream) GetTab() int {
    return my.tab
}
func (my *LexStream) SetStreamIndex(index int)  {
    my.index = index
}
func (my *LexStream) GetStreamIndex() int {
    return my.index
}
func (my *LexStream) SetStreamLength(streamLength int)  {
    my.streamLength = streamLength
}
func (my *LexStream) GetStreamLength() int {
    return my.streamLength
}
func (my *LexStream) SetLineOffSet(i int)  {
    my.lineOffSets.Add(i)
}
func (my *LexStream) GetLineOffSet(i int) int {
    return my.lineOffSets.Get(i)
}
func (my *LexStream) SetPrsStream(prsStream IPrsStream)  {
    if nil == prsStream{
        return
    }
    prsStream.SetLexStream(my)
    my.prsStream = prsStream
}
func (my *LexStream) GetIPrsStream() IPrsStream{
    return my.prsStream
}

func (my *LexStream) OrderedExportedSymbols() []string {
    return nil
}
func (my *LexStream) GetCharValue(i int) string {
    return string(my.inputChars[i])
}
func (my *LexStream) GetIntValue(i int) int {
    return int(my.inputChars[i])
}

func (my *LexStream) GetLineCount() int {
    return my.lineOffSets.Size() - 1
}
func (my *LexStream) GetLineNumberOfCharAt(i int) int {
    var index int = my.lineOffSets.BinarySearch(i)
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
func (my *LexStream) GetColumnOfCharAt(i int) int {
    var lineNo int = my.GetLineNumberOfCharAt(i)
    var start int = my.lineOffSets.Get(lineNo - 1)
    if start + 1 >= my.streamLength {
        return 1
    }
    var k int = start + 1
    for  ;k < i ;k++ {
        if my.inputChars[k] == '\t' {
            var offSet int = (k - start) - 1
            start -= ((my.tab - 1) - offSet % my.tab)
        }
    }
    return i - start
}
func (my *LexStream) GetToken() int {
    my.index = my.GetNext(my.index)
    return my.index
}
func (my *LexStream) GetTokenFromEndToken(end_token int ) int {
     if my.index < end_token {
         my.index =my.GetNext(my.index)
     }else{
         my.index =my.streamLength
     }
     return  my.index
}
func (my *LexStream) GetKind(i int) int {
    return 0
}
func (my *LexStream) next(i int) int {
    return my.GetNext(i)
}
func (my *LexStream) GetNext(i int) int {
     i+=1
     if i < my.streamLength {
         return  i
     }else{
        return  my.streamLength
     }
}
func (my *LexStream) previous(i int) int {
    return my.GetPrevious(i)
}
func (my *LexStream) GetPrevious(i int) int {
    if i <= 0 {
        return 0
    }else {
       return i - 1
    }
}
func (my *LexStream) GetName(i int) string {
    if i >= my.GetStreamLength() {
        return ""
    }else{
        return  "" + my.GetCharValue(i)
    }
}
func (my *LexStream) Peek() int {
    return my.GetNext(my.index)
}
func (my *LexStream) ResetTo(i int)  {
    my.index = i - 1
}
func (my *LexStream) Reset()  {
    my.index = -1
}

func (my *LexStream) BadToken() int {
    return 0
}
func (my *LexStream) GetLine(i int) int {
    return my.GetLineNumberOfCharAt(i)
}

func (my *LexStream) GetColumn(i int) int {
    return my.GetColumnOfCharAt(i)
}
func (my *LexStream) GetEndLine(i int) int {
    return my.GetLine(i)
}
func (my *LexStream) GetEndColumn(i int) int {
    return my.GetColumnOfCharAt(i)
}
func (my *LexStream) AfterEol(i int) bool {
    if i < 1 {
        return  true
    } else{
        return my.GetLineNumberOfCharAt(i - 1) < my.GetLineNumberOfCharAt(i)
    }
}
func (my *LexStream) GetFirstErrorToken(i int) int {
    return my.GetFirstRealToken(i)
}
func (my *LexStream) GetFirstRealToken(i int) int {
    return i
}
func (my *LexStream) GetLastErrorToken(i int) int {
    return my.GetLastRealToken(i)
}
func (my *LexStream) GetLastRealToken(i int) int {
    return i
}

//
// Here is where we report errors.  The default method is simply to print the error message to the console.
// However, the user may supply an error message handler to process error messages.  To support that
// a message handler interface is provided that has a single method HandleMessage().  The user has his
// error message handler class implement the IMessageHandler interface and provides an object of this type
// to the runtime using the SetMessageHandler(errorMsg) method. If the message handler object is Set,
// the ReportError methods will invoke its HandleMessage() method.
//
func (my *LexStream) SetMessageHandler(errMsg IMessageHandler)  {
    my.errMsg = errMsg
}
func (my *LexStream) GetMessageHandler() IMessageHandler {
    return my.errMsg
}
func (my *LexStream) MakeToken(startLoc int, endLoc int, kind int)  {
    if my.prsStream == nil {
        my.prsStream.MakeToken(startLoc, endLoc, kind)
    } else {
        my.ReportLexicalErrorPosition(startLoc, endLoc)// make it a lexical error
    }
}

func (my *LexStream) GetLocation(leftLoc int, rightLoc int) []int {
    var endLoc int
    if rightLoc < my.streamLength {
        endLoc = rightLoc
    }else{
        endLoc =my.streamLength - 1
    }
    var length int = endLoc - leftLoc + 1

    return []int{leftLoc,
                 length,
                 my.GetLineNumberOfCharAt(leftLoc),
                 my.GetColumnOfCharAt(leftLoc),
                 my.GetLineNumberOfCharAt(rightLoc),
                 my.GetColumnOfCharAt(rightLoc)}
}
func (my *LexStream) ReportLexicalErrorPosition(leftLoc int, rightLoc int) {

        var errorCode int
        if rightLoc >= my.streamLength {
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
                tokenText= "\"" + my.ToString(leftLoc,  rightLoc+ 1) + "\" "
            }else{
                tokenText = "\"" + my.GetCharValue(leftLoc) + "\" "
            }
        }
        var errorLeftLoc = 0
        var errorRightLoc = 0
        var errorInfo = []string{tokenText}
        my.ReportLexicalError(errorCode, leftLoc, rightLoc, errorLeftLoc, errorRightLoc, errorInfo)

}
func (my *LexStream) ReportLexicalError(leftLoc int, rightLoc int, errorCode int, errorLeftLoc int,
    errorRightLoc int, errorInfo  []string)  {

    if nil == my.errMsg {
       var locationInfo =fmt.Sprintf("%s : %d : %d : %d : %d : %d : %d : %d",
            my.GetFileName(),
            my.GetLineNumberOfCharAt(leftLoc),
            my.GetColumnOfCharAt(leftLoc),
            my.GetLineNumberOfCharAt(rightLoc),
            my.GetColumnOfCharAt(rightLoc),
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
        my.errMsg.HandleMessage(errorCode,
                                  my.GetLocation(leftLoc, rightLoc),
                                  my.GetLocation(errorLeftLoc, errorRightLoc),
                                  my.GetFileName(),
                                  errorInfo)
    }
}

func (my *LexStream) ReportError(errorCode int, leftToken int, rightToken int, errorInfo []string, errorToken int)  {
    my.ReportLexicalError(leftToken, rightToken, errorCode,  errorToken,errorToken, errorInfo)
}

func (my *LexStream) ToString(startOffSet int, endOffSet int) string {
    var length int = endOffSet - startOffSet + 1
    if endOffSet >= len(my.inputChars) {
        return "$EOF"
    } else{
        if length <= 0 {
            return ""
        } else{
            return string(my.inputChars[startOffSet :endOffSet])
        }
    }
}


