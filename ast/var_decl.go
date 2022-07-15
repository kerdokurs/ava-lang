package ast

import "fmt"

type VarDecl struct {
	Name string
	Type string
	Init Expr
}

func (v VarDecl) Accept(interp Visitor) any {
	return interp.VisitVarDecl(v)
}

func (v VarDecl) String() string {
	typ := v.Type
	if len(typ) == 0 {
		typ = "?"
	}
	exprStr := v.Init.String()

	return fmt.Sprintf("ValDecl(%s, %s, %s)", v.Name, typ, exprStr)
}

func (v VarDecl) stmtNode() {}

func (v VarDecl) glblStmt() {}
