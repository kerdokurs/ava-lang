package ast

type NilLit struct {
}

func (n NilLit) Accept(interp Interpreter) any {
	return interp.VisitNilLit(n)
}

func (n NilLit) String() string {
	return "NilLit()"
}

func (n NilLit) exprNode() {}
