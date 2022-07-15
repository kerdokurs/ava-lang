package ast

import "fmt"

type Variable struct {
	Name string
}

func (v Variable) Accept(interp Interpreter) any {
	return interp.VisitVariable(v)
}

func (v Variable) String() string {
	return fmt.Sprintf("Variable(%s)", v.Name)
}

func (v Variable) exprNode() {
}
