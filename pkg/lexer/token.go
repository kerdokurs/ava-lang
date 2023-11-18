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

	Ident
	Keyword

	LParen
	RParen
	LCurly
	RCurly

	Semi
	Comma

	Operator
)

var enumToString = []string{
	"EOF",
	"INT",
	"FLOAT",
	"IDENTIFIER",
	"KEYWORD",
	"LPAREN",
	"RPAREN",
	"LCURLY",
	"RCURLY",
	"SEMI COLON",
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
