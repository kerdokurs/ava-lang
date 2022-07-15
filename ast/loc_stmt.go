package ast

import "fmt"

type LocStmt struct {
	Value string
}

func (l LocStmt) Accept(interp Visitor) any {
	return interp.VisitLocStmt(l)
}

func (l LocStmt) String() string {
	return fmt.Sprintf("Loc(%s)", l.Value)
}

func (l LocStmt) stmtNode() {}
