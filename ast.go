package main

import (
	"io"
	"kerdo.dev/ava/ast"
)

func CreateAst(source io.Reader) ast.ProgStmt {
	lexer := NewLexer(source)
	tokens := lexer.ReadAllTokens()

	parser := NewParser(tokens)
	return parser.Parse()
}
