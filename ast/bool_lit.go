package ast

import "fmt"

type BoolLit struct {
	Value bool
}

func (b BoolLit) Accept(interp Interpreter) any {
	return interp.VisitBoolLit(b)
}

func (b BoolLit) String() string {
	return fmt.Sprintf("BoolLit(%t)", b.Value)
}

func (b BoolLit) exprNode() {}
