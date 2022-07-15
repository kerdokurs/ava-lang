package ast

import "fmt"

type FloatLit struct {
	Value float64
}

func (f FloatLit) Accept(interp Visitor) any {
	return interp.VisitFloatLit(f)
}

func (f FloatLit) String() string {
	return fmt.Sprintf("FloatLit(%f)", f.Value)
}

func (f FloatLit) exprNode() {}
