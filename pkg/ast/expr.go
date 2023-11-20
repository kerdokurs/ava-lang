package ast

import (
	"fmt"
	"strings"

	"kerdo.dev/ava-lang/pkg/utils"
)

type Block struct {
	Stmts []Stmt

	// count of local variables in the current scope. does not include local variables in child blocks
	LocalCount int
}

func (b Block) exprNode() {}
func (b Block) String() string {
	stmtStrs := strings.Join(utils.Map[Stmt, string](b.Stmts, func(s Stmt) string {
		return s.String()
	}), ",")
	return fmt.Sprintf("Block(%s)", stmtStrs)
}

func (b Block) CountLocals() int {
	var count int
	for _, stmt := range b.Stmts {
		if _, ok := stmt.(VarDecl); ok {
			count += 1
		}
	}
	return count
}

type Call struct {
	Name         string
	IsArithmetic bool
	Args         []Expr
}

func (l Call) exprNode() {}
func (l Call) String() string {
	exprStrs := strings.Join(utils.Map[Expr, string](l.Args, Expr.String), ",")
	return fmt.Sprintf("Call(%s, %s)", l.Name, exprStrs)
}

type Variable struct {
	Name string
}

func (v Variable) exprNode() {}
func (v Variable) String() string {
	return fmt.Sprintf("Variable(%s)", v.Name)
}

type IntLiteral struct {
	Value int
}

func (l IntLiteral) exprNode() {}
func (l IntLiteral) String() string {
	return fmt.Sprintf("Int(%d)", l.Value)
}

type StringLiteral struct {
	Value string
}

func (l StringLiteral) exprNode() {}
func (l StringLiteral) String() string {
	return fmt.Sprintf("String(\"%s\")", l.Value)
}
