package ast

type Node interface {
	String() string
	Accept(interp Interpreter) any
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
