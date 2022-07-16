package ast

type Visitor interface {
	Visit(Node) any

	VisitProgStmt(ProgStmt) any
	VisitLocStmt(LocStmt) any

	VisitParenExpr(ParenExpr) any

	VisitBlock(Block) any

	VisitIfStmt(IfStmt) any

	VisitAssignStmt(AssignStmt) any

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
