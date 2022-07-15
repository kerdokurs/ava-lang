package ast

import "fmt"

type Variable struct {
	Name string
}

func (v Variable) String() string {
	return fmt.Sprintf("Variable(%s)", v.Name)
}

func (v Variable) exprNode() {
}