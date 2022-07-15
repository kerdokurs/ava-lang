package ast

import "fmt"

type ConstDecl struct {
	Name string
	Type string
	Init Expr
}

func (c ConstDecl) Accept(interp Interpreter) any {
	return interp.VisitConstDecl(c)
}

func (c ConstDecl) String() string {
	typ := c.Type
	if len(typ) == 0 {
		typ = "?"
	}
	exprStr := c.Init.String()

	return fmt.Sprintf("ConstDecl(%s, %s, %s)", c.Name, typ, exprStr)
}

func (c ConstDecl) stmtNode() {}

func (c ConstDecl) glblStmt() {}
