package main

type TokenType int

const (
	EOF = iota

	INT
	HEX
	FLOAT
	STRING
	BOOL
	NIL

	OPERATOR

	IDENT
	ITYPE
	// Intrinsic functions?

	KEYWORD

	LPAREN
	RPAREN
	LCURLY
	RCURLY

	SEMI
	COMMA

	LCOMMENT
	BCOMMENT
)

var enumToString = []string{
	"EOF",
	"INT",
	"HEX",
	"FLOAT",
	"STRING",
	"BOOL",
	"NIL",
	"OPERATOR",
	"IDENTIFIER",
	"INTRINSIC TYPE",
	"KEYWORD",
	"LPAREN",
	"RPAREN",
	"LCURLY",
	"RCURLY",
	"SEMI COLON",
	"COMMA",
	"LINE COMMENT",
	"BLOCK COMMENT",
}

type Token struct {
	Type TokenType
	Data string
}

func (t *Token) Name() string {
	return enumToString[t.Type]
}

func Name(typ TokenType) string {
	return enumToString[typ]
}
