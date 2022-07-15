package ast

type NilLit struct {
}

func (n NilLit) String() string {
	return "NilLit()"
}

func (n NilLit) exprNode() {}
