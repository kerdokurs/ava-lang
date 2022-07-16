package types

import "kerdo.dev/ava/ast"

type FunctionDefinition struct {
	Name   string
	Params []ast.FuncParam
	Body   ast.Block
}
