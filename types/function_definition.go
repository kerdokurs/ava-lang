package types

import "kerdo.dev/ava/ast"

type FunctionDefinition struct {
	Name string
	Body ast.Block
}
