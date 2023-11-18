package ast

import (
	"fmt"
	"strings"
)

type ProgStmt struct {
	Decls []Decl
}

func (s *ProgStmt) stmtNode() {}
func (s *ProgStmt) String() string {
	declStrs := strings.Builder{}

	for _, decl := range s.Decls {
		declStrs.WriteString(decl.String())
		declStrs.WriteRune('\n')
	}

	return fmt.Sprintf("Prog(\n%s)", declStrs.String())
}

type VarDecl struct {
	Name string
	Expr Expr
}

func (d VarDecl) stmtNode() {}
func (d VarDecl) String() string {
	return fmt.Sprintf("VarDecl(%s, %s)", d.Name, d.Expr.String())
}

type ExprStmt struct {
	Expr Expr
}

func (s ExprStmt) stmtNode() {}
func (s ExprStmt) String() string {
	return fmt.Sprintf("Stmt(%s)", s.Expr.String())
}

type ReturnStmt struct {
	Expr Expr
}

func (s ReturnStmt) stmtNode() {}
func (s ReturnStmt) String() string {
	return fmt.Sprintf("Return(%s)", s.Expr.String())
}
