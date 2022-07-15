package ast

import (
	"fmt"
	"strings"
)

type Block struct {
	Stmts          []Stmt
	ImplicitReturn *Expr
}

func (b Block) String() string {
	parts := make([]string, len(b.Stmts))
	for i, stmt := range b.Stmts {
		parts[i] = stmt.String()
	}

	str := strings.Join(parts, ",\n")

	return fmt.Sprintf("Block(%s)", str)
}

func (b Block) Accept(interp Visitor) any {
	return interp.VisitBlock(b)
}

func (b Block) stmtNode() {}
