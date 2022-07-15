package ast

type Node interface {
	String() string
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Decl interface {
	Node
	declNode()
}

type GlblStmt interface {
	Stmt
	glblStmt()
}
