package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strings"
	"unicode"
)

type Lexer struct {
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) ReadAllTokens() []Token {
	tokens := make([]Token, 0)

	for {
		token := l.readNextToken()
		tokens = append(tokens, token)

		if token.Type == EOF {
			break
		}
	}

	return tokens
}

func (l *Lexer) readNextToken() Token {
	for {
		rs, err := l.reader.Peek(1)

		if err != nil {
			if errors.Is(err, io.EOF) {
				return Token{
					Type: EOF,
					Data: "",
				}
			}

			log.Fatalf("Error reading input: %v\n", err)
		}

		r := rune(rs[0])

		if unicode.IsSpace(r) {
			l.reader.ReadRune()
			continue
		} else if r == '/' {
			return l.readDivisionOrComment()
		} else if couldBeOperator(r) {
			return l.readOperator()
		} else if unicode.IsDigit(r) {
			return l.readNumericLiteral()
		} else if unicode.IsLetter(r) {
			return l.readIdentOrKeyword()
		} else if r == '(' {
			return l.readSingleChar(LPAREN)
		} else if r == ')' {
			return l.readSingleChar(RPAREN)
		} else if r == '{' {
			return l.readSingleChar(LCURLY)
		} else if r == '}' {
			return l.readSingleChar(RCURLY)
		} else if r == ';' {
			return l.readSingleChar(SEMI)
		} else if r == ',' {
			return l.readSingleChar(COMMA)
		} else if r == '"' {
			return l.readStrLiteral()
		}

		log.Fatalf("Invalid rune: %s\n", string(r))
	}
}

func (l *Lexer) readDivisionOrComment() Token {
	_, _, _ = l.reader.ReadRune()

	r, _, err := l.reader.ReadRune()
	if err != nil {
		panic(err)
	}

	sb := strings.Builder{}
	if r == '/' {
		for {
			r, _, err = l.reader.ReadRune()
			if err != nil {
				panic(err)
			}

			if r == '\n' {
				break
			}

			sb.WriteRune(r)
		}

		return Token{
			Type: LCOMMENT,
			Data: sb.String(),
		}
	} else if r == '*' {
		for {
			rs, err := l.reader.Peek(2)
			if err != nil {
				panic(err)
			}

			if rs[0] == '*' && rs[1] == '/' {
				break
			}

			r, _, err := l.reader.ReadRune()
			if err != nil {
				panic(err)
			}

			sb.WriteRune(r)
		}

		return Token{
			Type: BCOMMENT,
			Data: sb.String(),
		}
	}

	return Token{
		Type: OPERATOR,
		Data: "/",
	}
}

func (l *Lexer) readStrLiteral() Token {
	_, _, err := l.reader.ReadRune()
	if err != nil {
		panic(err)
	}

	sb := strings.Builder{}

	escaping := false

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			panic(err)
		}

		if r == '\\' {
			escaping = true
			continue
		} else if r == '"' && !escaping {
			break
		}

		sb.WriteRune(r)
	}

	data := sb.String()
	return Token{
		Type: STRING,
		Data: data,
	}
}

func (l *Lexer) readSingleChar(typ TokenType) Token {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		panic(err)
	}

	return Token{
		Type: typ,
		Data: string(r),
	}
}

func (l *Lexer) readIdentOrKeyword() Token {
	sb := strings.Builder{}

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			panic(err)
		}

		if !(unicode.IsDigit(r) || unicode.IsLetter(r)) {
			l.reader.UnreadRune()
			break
		}

		sb.WriteRune(r)
	}

	var typ TokenType = IDENT

	data := sb.String()
	if contains(keywords, data) {
		typ = KEYWORD
	} else if contains(intrinsicTypes, data) {
		typ = ITYPE
	} else if data == "true" || data == "false" {
		typ = BOOL
	}

	return Token{
		Type: typ,
		Data: data,
	}
}

var keywords = []string{
	"if", "while",
	"var", "fun", "const",
	"loc", "use",
	"struct", "impl",
}

var intrinsicTypes = []string{
	"u8", "i8",
	"u16", "i16",
	"u32", "i32",
	"u64", "i64",
	"void",
	"f32", "f64",
	"str",
}

func (l *Lexer) readOperator() Token {
	builder := strings.Builder{}

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			panic(err)
		}

		if !contains(opRunes, r) {
			l.reader.UnreadRune()
			break
		}

		builder.WriteRune(r)
	}

	data := builder.String()

	return Token{
		Type: OPERATOR,
		Data: data,
	}
}

func (l *Lexer) readNumericLiteral() Token {
	var typ TokenType = INT
	sb := strings.Builder{}

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Fatalf("Error reading input: %v\n", err)
			}
			break
		}

		if r == '.' {
			typ = FLOAT
		} else if r == 'x' {
			typ = HEX
		} else if !unicode.IsDigit(r) {
			err = l.reader.UnreadRune()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatalf("Error unreading from input: %v\n", err)
			}
			break
		}

		sb.WriteRune(r)
	}

	data := sb.String()

	// TODO: Add other types
	switch typ {
	case INT:
		if len(data) >= 2 && data[0] == '0' && data[1] == '0' {
			log.Fatalf("Illegal integer: %s\n", data)
		}
	case HEX:
		if len(data) >= 2 && data[0] == '0' && data[1] == '0' {
			log.Fatalf("Illegal integer: %s\n", data)
		}
	case FLOAT:
		if len(data) >= 2 && data[0] == '0' && data[1] == '0' {
			log.Fatalf("Illegal integer: %s\n", data)
		}
	}

	return Token{
		Type: typ,
		Data: data,
	}
}

func contains[T comparable](arr []T, val T) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}

	return false
}

var operators = []string{
	".", "::", "->", ":",
	"+", "-", "*", "/",
	"%", "<", ">", "<=", ">=", "==", "!=",
	"&", "&&", "|", "||",
}
var opRunes = Map(operators, func(op string) rune {
	return rune(op[0])
})

func couldBeOperator(r rune) bool {
	for _, op := range operators {
		if rune(op[0]) == r {
			return true
		}
	}

	return false
}
