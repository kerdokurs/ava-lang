package lexer

type Token struct {
	Type  TokenType
	Value string
}

type TokenType int

const (
	EOF TokenType = iota

	Int
	Float
	String

	Ident
	Keyword

	LParen
	RParen
	LCurly
	RCurly

	Semi
	Colon
	Comma

	Operator
)

var enumToString = []string{
	"EOF",
	"INT",
	"FLOAT",
	"STRING",
	"IDENTIFIER",
	"KEYWORD",
	"LPAREN",
	"RPAREN",
	"LCURLY",
	"RCURLY",
	"SEMI COLON",
	"COLON",
	"COMMA",
	"OPERATOR",
}

func (t *Token) Name() string {
	return enumToString[t.Type]
}

func Name(typ TokenType) string {
	return enumToString[typ]
}

var keywords = []string{
	"let", "fun", "return", "while",
}
