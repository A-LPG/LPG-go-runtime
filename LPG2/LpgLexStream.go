package lpg2


type LpgLexStream(LexStream):
    

    def __init__( fileName: string, inputChars: string = nil, tab int = LexStream.DEFAULT_TAB,
                 lineOffsets: IntSegmentedTuple = nil):
        super().__init__(fileName, inputChars, tab, lineOffsets)
