package main

import (
	"fmt"
	"io"
	"kerdo.dev/ava/ast"
)

const Debug = true

func CreateAst(source io.Reader) ast.ProgStmt {
	lexer := NewLexer(source)
	tokens := lexer.ReadAllTokens()

	if Debug {
		for i, token := range tokens {
			fmt.Printf("%d -> %+v\n", i, token)
		}
	}

	parser := NewParser(tokens)
	return parser.Parse()
}
