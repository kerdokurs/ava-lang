package ast

import "fmt"

type UnaryExpr struct {
	Operator string
	Expr     Expr
}

func (u UnaryExpr) String() string {
	return fmt.Sprintf("UnaryExpr(%s, %s)", u.Operator, u.Expr.String())
}

func (u UnaryExpr) exprNode() {}
