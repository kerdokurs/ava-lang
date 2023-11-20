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
	Name    string
	Expr    Expr
	HasInit bool
}

func (d VarDecl) stmtNode() {}
func (d VarDecl) String() string {
	exprStr := "nil"
	if d.HasInit {
		exprStr = d.Expr.String()
	}
	return fmt.Sprintf("VarDecl(%s, %s)", d.Name, exprStr)
}

type ExprStmt struct {
	Expr Expr
}

func (s ExprStmt) stmtNode() {}
func (s ExprStmt) String() string {
	return fmt.Sprintf("ExprStmt(%s)", s.Expr.String())
}

type WhileStmt struct {
	Condition Expr
	Body      Block
}

func (s WhileStmt) stmtNode() {}
func (s WhileStmt) String() string {
	return fmt.Sprintf("While(%s, %s)", s.Condition.String(), s.Body.String())
}

type ReturnStmt struct {
	Expr Expr
}

func (s ReturnStmt) stmtNode() {}
func (s ReturnStmt) String() string {
	return fmt.Sprintf("Return(%s)", s.Expr.String())
}

type AssignStmt struct {
	Variable Variable
	Expr     Expr
}

func (s AssignStmt) stmtNode() {}
func (s AssignStmt) String() string {
	return fmt.Sprintf("Assign(%s, %s)", s.Variable.String(), s.Expr.String())
}
