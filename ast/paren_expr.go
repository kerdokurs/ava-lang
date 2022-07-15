package ast

import "fmt"

type ParenExpr struct {
	Expr Expr
}

func (p ParenExpr) String() string {
	return fmt.Sprintf("Parens(%s)", p.Expr.String())
}

func (p ParenExpr) exprNode() {}
