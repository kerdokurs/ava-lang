package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"kerdo.dev/ava-lang/pkg/utils"
)

type Lexer struct {
	reader *bufio.Reader
}

func New(reader io.Reader) Lexer {
	return Lexer{
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) Lex() []Token {
	tokens := make([]Token, 0)

	for {
		token := l.nextToken()
		tokens = append(tokens, token)

		if token.Type == EOF {
			break
		}
	}

	return tokens
}

func (l *Lexer) nextToken() Token {
	for {
		rs, err := l.reader.Peek(1)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return Token{}
			}

			panic(err)
		}

		r := rune(rs[0])

		if unicode.IsSpace(r) {
			l.reader.ReadRune()
			continue
		} else if couldBeOperator(r) {
			return l.readOperator()
		} else if unicode.IsNumber(r) {
			return l.readNumericLiteral()
		} else if unicode.IsLetter(r) {
			return l.readIdentOrKeyword()
		} else if r == '(' {
			return l.readSingleChar(LParen)
		} else if r == ')' {
			return l.readSingleChar(RParen)
		} else if r == '{' {
			return l.readSingleChar(LCurly)
		} else if r == '}' {
			return l.readSingleChar(RCurly)
		} else if r == ';' {
			return l.readSingleChar(Semi)
		}

		fmt.Printf("unimplemented: %c\n", r)
		return Token{}
	}
}

func (l *Lexer) readNumericLiteral() Token {
	var tokenType TokenType = Int
	sb := strings.Builder{}

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			panic(err)
		}

		if !unicode.IsNumber(r) {
			if err := l.reader.UnreadRune(); err != nil {
				panic(err)
			}
			break
		}

		sb.WriteRune(r)
	}

	return Token{
		Type:  tokenType,
		Value: sb.String(),
	}
}

func (l *Lexer) readIdentOrKeyword() Token {
	sb := strings.Builder{}

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			panic(err)
		}

		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			if err := l.reader.UnreadRune(); err != nil {
				panic(err)
			}

			break
		}

		sb.WriteRune(r)
	}

	var tokenType TokenType = Ident

	return Token{
		Type:  tokenType,
		Value: sb.String(),
	}
}

func (l *Lexer) readSingleChar(tokenType TokenType) Token {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		panic(err)
	}

	return Token{
		Type:  tokenType,
		Value: string(r),
	}
}

func (l *Lexer) readOperator() Token {
	sb := strings.Builder{}

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			panic(err)
		}

		if !utils.Contains(opRunes, r) {
			if err := l.reader.UnreadRune(); err != nil {
				panic(err)
			}
			break
		}

		sb.WriteRune(r)
	}

	return Token{
		Type:  Operator,
		Value: sb.String(),
	}
}

var operators = []string{
	".", "->",
	"+", "-", "*", "/",
	"%", "<", ">", "<=", ">=", "==", "!=",
	"&", "&&", "|", "||",
}
var opRunes = utils.Map(operators, func(op string) rune {
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
