package ast

import "fmt"

type ExprStmt struct {
	Expr Expr
}

func (e ExprStmt) String() string {
	return fmt.Sprintf("ExprStmt(%s)", e.Expr.String())
}

func (e ExprStmt) Accept(interp Interpreter) any {
	return interp.VisitExprStmt(e)
}

func (e ExprStmt) stmtNode() {}
