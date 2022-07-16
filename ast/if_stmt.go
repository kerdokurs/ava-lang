package ast

import (
	"fmt"
)

type IfStmt struct {
	Condition Expr
	ThenBody  Block
	HasElse   bool // remove if refactor to pointers
	ElseBody  Block
}

func (i IfStmt) String() string {
	elseBlock := ""
	if i.HasElse {
		elseBlock = ", \t" + i.ElseBody.String()
	}

	return fmt.Sprintf("IfStmt(\n\t%s, \t%s%s)", i.Condition.String(), i.ThenBody.String(), elseBlock)
}

func (i IfStmt) Accept(interp Visitor) any {
	return interp.VisitIfStmt(i)
}

func (i IfStmt) stmtNode() {}
