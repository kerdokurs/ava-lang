package ast

import (
	"fmt"
	"strings"

	"kerdo.dev/ava-lang/pkg/utils"
)

type Block struct {
	Stmts []Stmt
}

func (b Block) exprNode() {}
func (b Block) String() string {
	stmtStrs := strings.Join(utils.Map[Stmt, string](b.Stmts, func(s Stmt) string {
		return s.String()
	}), ",")
	return fmt.Sprintf("Block(%s)", stmtStrs)
}

type Call struct {
	Name         string
	IsArithmetic bool
	Args         []Expr
}

func (l Call) exprNode() {}
func (l Call) String() string {
	return fmt.Sprintf("Call(%s)", l.Name)
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
