package ast

import "fmt"

type BoolLit struct {
	Value bool
}

func (b BoolLit) String() string {
	return fmt.Sprintf("BoolLit(%t)", b.Value)
}

func (b BoolLit) exprNode() {}
