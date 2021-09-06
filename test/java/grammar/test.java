package test;



import org.eclipse.imp.lpg.parser.LPGLexer;
import org.eclipse.imp.lpg.parser.LPGParser;

import java.io.File;
import java.io.IOException;

public class test2 {
    public3  static  void main(String[] arg){

        String filename= "C:\\Users\\kuafu\\source\\repos\\lpgRuntimeCpp\\LpgExample\\LpgExample.cpp";
        filename= "E:/LPG2/lpg2/src/jikespg.g";
        String dddd =filename.substring(0,1);
        int num = 3/6 + 6;
        File f= new File(filename);
        if (f.exists()) {
            try {
                LPGLexer lexer= new LPGLexer(); // Create the lexer
                // SMS 12 Feb 2008
                // Previously attempted to get lex stream before it has been set,
                // so added a call to reset(..) here.  This call to reset(..)
                // obviates the need to initialize the lexer later on in the
                // method (so code related to that can be removed)

                lexer.reset(filename, 4);

                LPGParser parser= new LPGParser(lexer.getILexStream()); // Create the parser

                lexer.lexer(null, parser.getIPrsStream());
                LPGParser.ASTNode ast= (LPGParser.ASTNode) parser.parser();
            } catch (java.io.IOException e) {
                // skip this file
            }
        }
    }
}
