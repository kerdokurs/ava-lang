package ast

type Interpreter interface {
	Visit(Node) any

	VisitProgStmt(ProgStmt) any
	VisitLocStmt(LocStmt) any

	VisitParenExpr(ParenExpr) any

	VisitBlock(Block) any

	VisitExprStmt(ExprStmt) any

	VisitFuncCall(FuncCall) any

	VisitFuncDecl(FuncDecl) any
	VisitConstDecl(ConstDecl) any
	VisitVarDecl(VarDecl) any

	VisitVariable(Variable) any

	VisitIntLit(IntLit) any
	VisitFloatLit(FloatLit) any
	VisitBoolLit(BoolLit) any
	VisitNilLit(NilLit) any
	VisitStrLit(StrLit) any
}
