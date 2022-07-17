package main

type Visitor interface {
	Visit(Node) AvaVal

	VisitProgStmt(ProgStmt) AvaVal
	VisitLocStmt(LocStmt) AvaVal

	VisitParenExpr(ParenExpr) AvaVal

	VisitBlock(Block) AvaVal

	VisitIfStmt(IfStmt) AvaVal
	VisitWhileStmt(WhileStmt) AvaVal

	VisitStructDecl(StructDecl) AvaVal

	VisitAssignStmt(AssignStmt) AvaVal

	VisitExprStmt(ExprStmt) AvaVal

	VisitFuncCall(FuncCall) AvaVal

	VisitFuncDecl(FuncDecl) AvaVal
	VisitConstDecl(ConstDecl) AvaVal
	VisitVarDecl(VarDecl) AvaVal

	VisitVariable(Variable) AvaVal

	VisitIntLit(IntLit) AvaVal
	VisitFloatLit(FloatLit) AvaVal
	VisitBoolLit(BoolLit) AvaVal
	VisitStrLit(StrLit) AvaVal
}
