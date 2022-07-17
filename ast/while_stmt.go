package ast

import "fmt"

type WhileStmt struct {
	Condition Expr
	Body      Block
}

func (w WhileStmt) String() string {
	return fmt.Sprintf("WhileStmt(%s, %s)", w.Condition.String(), w.Body.String())
}

func (w WhileStmt) Accept(interp Visitor) any {
	return interp.VisitWhileStmt(w)
}

func (w WhileStmt) stmtNode() {}
