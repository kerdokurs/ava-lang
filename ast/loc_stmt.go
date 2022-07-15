package ast

import "fmt"

type LocStmt struct {
	Value string
}

func (l LocStmt) String() string {
	return fmt.Sprintf("Loc(%s)", l.Value)
}

func (l LocStmt) stmtNode() {}
