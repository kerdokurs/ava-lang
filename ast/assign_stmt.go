package ast

import "fmt"

type AssignStmt struct {
	Variable string
	Value    Expr
}

func (a AssignStmt) String() string {
	return fmt.Sprintf("AssignStmt(%s, %s)", a.Variable, a.Value.String())
}

func (a AssignStmt) Accept(interp Visitor) any {
	return interp.VisitAssignStmt(a)
}

func (a AssignStmt) stmtNode() {}

func (a AssignStmt) exprNode() {}
