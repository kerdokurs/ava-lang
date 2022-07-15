package ast

import "fmt"

type BinaryExpr struct {
	X        Expr
	Operator string
	Y        Expr
}

func (b BinaryExpr) String() string {
	return fmt.Sprintf("BinaryExpr(%s, %s, %s)", b.X.String(), b.Operator, b.Y.String())
}

func (b BinaryExpr) exprNode() {
}
