grammar Ava;

program
    :   globalStmts
    ;

globalStmts
    :   globalStmt*
    ;

globalStmt
    :   varDecl ';'
    |   constDecl ';'
    |   funcDef
    ;

blockStmts
    :   blockStmt*
    ;

blockStmt
    :   exprStmt ';'
    |   varDecl ';'
    |   constDecl ';'
    |   assign ';'
    ;

// Function definition
funcDef
    :   'func' name=Ident '(' funcArgsDef? ')' funcRetDef? body=block
    ;

funcArgsDef
    :   funcArgDef (',' funcArgDef)* ','?
    ;

funcArgDef
    :   name=Ident ':' type=varType nullable='?'?
    ;

funcRetDef
    :   '->' type=varType nullable='?'?
    ;

block
    :   '{' blockStmts '}'
    ;

exprStmt
    : expr
    ;

expr
    :   addExpr
    |   compExpr
    ;

compExpr
    :   left=addExpr op=('=='|'!='|'<='|'>='|'<'|'>') right=addExpr #BinaryComp
    |   addExpr                                                     #SimpleComp
    ;

addExpr
    :   left=addExpr op=('+'|'-') right=mulExpr #BinaryAdd
    |   mulExpr                                 #SimpleAdd
    ;

mulExpr
    :   left=mulExpr op=('*'|'/'|'%') right=unaryExpr   #BinaryMul
    |   unaryExpr                                       #SimpleMul
    ;

unaryExpr
    :   op='-' unaryExpr   #UnaryOp
    |   funcExpr        #SimpleUnary
    ;

funcExpr
    :   name=Ident '(' args=argList? ')'                       #FuncCall
    |   '(' expr ')'                                           #Parens
    |   IntLit                                                 #IntLiteral
    |   StrLit                                                 #StrLiteral
    |   BoolLit                                                #BoolLiteral
    |   NilLit                                                 #NilLiteral
    |   Ident                                                  #Variable
    ;

argList
    :   expr (',' expr)* ','?
    ;

varDecl
    :   'var' name=Ident (':' varType)? ('=' init=expr)? #UnmutableVarDecl
    |   'val' name=Ident (':' varType)? ('=' init=expr)? #MutableVarDecl
    ;

constDecl
    :   'const' name=Ident (':' varType)? '=' init=expr
    ;

assign
    :   name=Ident '=' expr
    ;

// Variable definition

varType
    :   ref='&'? IntrinsicType
    |   ref='&'? Ident ('<' genericType=varType '>')?
    ;

// Literals

IntLit
    :   '0'
    |   [1-9][0-9]*
    ;

HexLit
    :   '0x' IntLit
    ;

StrLit
    :   '"' ~["\n]* '"'
    ;

BoolLit
    :   'true'
    |   'false'
    ;

NilLit
    :   'nil'
    ;

// Lexer rules

IntrinsicType
    :   'u8'
    ;

Ident
    :   [A-Za-z_][A-Za-z0-9_]*
    ;

BlockComment
    :   '/*' ( . )*? '*/' -> skip
    ;

LineComment
    :   '//' (~[\n] .)* -> skip
    ;

WS: [\n\t\r ] -> skip;