package lpg2


const  NIL_CODE int  = -1
const  LEX_ERROR_CODE int  = 0
const  ERROR_CODE int = 1
const  BEFORE_CODE int = 2
const  INSERTION_CODE int = 3
const  INVALID_CODE int = 4
const  SUBSTITUTION_CODE int = 5
const  SECONDARY_CODE int = 5
const  DELETION_CODE int = 6
const  MERGE_CODE int = 7
const  MISPLACED_CODE int = 8
const  SCOPE_CODE int = 9
const  EOF_CODE int = 10
const  INVALID_TOKEN_CODE int = 11
const  ERROR_RULE_ERROR_CODE int = 11
const  ERROR_RULE_WARNING_CODE int = 12
const  NO_MESSAGE_CODE int = 13

const  MANUAL_CODE int = 14

var  errorMsgText =[]string{
"unexpected character ignored",     // $NON-NLS-1$
"parsing terminated at this token", // $NON-NLS-1$
" inserted before this token",      // $NON-NLS-1$
" expected after this token",       // $NON-NLS-1$
"unexpected input discarded",      // $NON-NLS-1$
" expected instead of this input", // $NON-NLS-1$
" unexpected token(s): ignored",   // $NON-NLS-1$
" formed from merged tokens",      // $NON-NLS-1$
"misplaced construct(s):",         // $NON-NLS-1$
" inserted to complete scope",  // $NON-NLS-1$
" reached after this token",    // $NON-NLS-1$
" is invalid",                  // $NON-NLS-1$
" is ignored",                  // $NON-NLS-1$
"",                              // $NON-NLS-1$
}
