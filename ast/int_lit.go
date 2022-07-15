package ast

import "fmt"

type IntLit struct {
	Value int
}

func (i IntLit) String() string {
	return fmt.Sprintf("IntLit(%d)", i.Value)
}

func (i IntLit) exprNode() {}
