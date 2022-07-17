package main

import (
	"fmt"
	"io"
	"strings"
)

func CreateAst(source io.Reader) ProgStmt {
	lexer := NewLexer(source)
	tokens := lexer.ReadAllTokens()

	if IsDebug {
		for i, token := range tokens {
			fmt.Printf("%d -> %+v\n", i, token)
		}
	}

	parser := NewParser(tokens)
	return parser.Parse()
}

// Base

type Node interface {
	String() string
	Accept(interp Visitor) AvaVal
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

// Assign statement

type AssignStmt struct {
	Variable string
	Value    Expr
}

func (a AssignStmt) String() string {
	return fmt.Sprintf("AssignStmt(%s, %s)", a.Variable, a.Value.String())
}

func (a AssignStmt) Accept(interp Visitor) AvaVal {
	return interp.VisitAssignStmt(a)
}

func (a AssignStmt) stmtNode() {}

func (a AssignStmt) exprNode() {}

// Block statement

type Block struct {
	Stmts          []Stmt
	ImplicitReturn *Expr
}

func (b Block) String() string {
	parts := make([]string, len(b.Stmts))
	for i, stmt := range b.Stmts {
		parts[i] = stmt.String()
	}

	str := strings.Join(parts, ",\n")

	return fmt.Sprintf("Block(%s)", str)
}

func (b Block) Accept(interp Visitor) AvaVal {
	return interp.VisitBlock(b)
}

func (b Block) stmtNode() {}

// Bool literal

type BoolLit struct {
	Value bool
}

func (b BoolLit) Accept(interp Visitor) AvaVal {
	return interp.VisitBoolLit(b)
}

func (b BoolLit) String() string {
	return fmt.Sprintf("BoolLit(%t)", b.Value)
}

func (b BoolLit) exprNode() {}

// Const declaration statement

type ConstDecl struct {
	Name     string
	Type     string
	Init     Expr
	IsGlobal bool
}

func (c ConstDecl) Accept(interp Visitor) AvaVal {
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

func (c ConstDecl) declNode() {}

// Expression statement

type ExprStmt struct {
	Expr Expr
}

func (e ExprStmt) String() string {
	return fmt.Sprintf("ExprStmt(%s)", e.Expr.String())
}

func (e ExprStmt) Accept(interp Visitor) AvaVal {
	return interp.VisitExprStmt(e)
}

func (e ExprStmt) stmtNode() {}

// Float literal

type FloatLit struct {
	Value float64
}

func (f FloatLit) Accept(interp Visitor) AvaVal {
	return interp.VisitFloatLit(f)
}

func (f FloatLit) String() string {
	return fmt.Sprintf("FloatLit(%f)", f.Value)
}

func (f FloatLit) exprNode() {}

// Function call expression

type FuncCall struct {
	Name         string
	IsArithmetic bool
	IsComparison bool
	Args         []Expr
}

func (f FuncCall) Accept(interp Visitor) AvaVal {
	return interp.VisitFuncCall(f)
}

func (f FuncCall) String() string {
	strs := make([]string, len(f.Args))
	for i, arg := range f.Args {
		strs[i] = arg.String()
	}
	args := strings.Join(strs, ", ")
	if len(args) > 0 {
		args = ", " + args
	}

	return fmt.Sprintf("FuncCall(%s%s)", f.Name, args)
}

func (f FuncCall) exprNode() {}

// Function declaration statement

type FuncDecl struct {
	Name       string
	ReturnType string
	Params     []FuncParam
	Body       Block
}

func (f FuncDecl) Accept(interp Visitor) AvaVal {
	return interp.VisitFuncDecl(f)
}

type FuncParam struct {
	Name string
	Type string
}

func (f FuncDecl) String() string {
	retType := f.ReturnType
	if len(retType) == 0 {
		retType = "?"
	}

	params := make([]string, len(f.Params))
	for i, param := range f.Params {
		s := param.String()
		s = strings.Replace(s, "\n", "\n\t", -1)
		params[i] = "\t" + param.String()
	}
	paramsStr := strings.Join(params, ",\n")

	bodyStr := f.Body.String()

	return fmt.Sprintf("FuncDecl(\n\t%s,\n\t%s,\n%s\t%s\n)", f.Name, retType, paramsStr, bodyStr)
}

func (f FuncDecl) stmtNode() {}

func (f FuncDecl) glblStmt() {}

func (f FuncParam) String() string {
	return fmt.Sprintf("FuncParam(%s, %s)", f.Name, f.Type)
}

// If statement

type IfStmt struct {
	Condition Expr
	ThenBody  Block
	HasElse   bool // remove if refactor to pointers
	ElseBody  Block
}

func (i IfStmt) String() string {
	elseBlock := ""
	if i.HasElse {
		elseBlock = ", \t" + i.ElseBody.String()
	}

	return fmt.Sprintf("IfStmt(\n\t%s, \t%s%s)", i.Condition.String(), i.ThenBody.String(), elseBlock)
}

func (i IfStmt) Accept(interp Visitor) AvaVal {
	return interp.VisitIfStmt(i)
}

func (i IfStmt) stmtNode() {}

// Integer literal

type IntLit struct {
	Value int
}

func (i IntLit) Accept(interp Visitor) AvaVal {
	return interp.VisitIntLit(i)
}

func (i IntLit) String() string {
	return fmt.Sprintf("IntLit(%d)", i.Value)
}

func (i IntLit) exprNode() {}

// Location statement

type LocStmt struct {
	Value string
}

func (l LocStmt) Accept(interp Visitor) AvaVal {
	return interp.VisitLocStmt(l)
}

func (l LocStmt) String() string {
	return fmt.Sprintf("Loc(%s)", l.Value)
}

func (l LocStmt) stmtNode() {}

// Parens expression

type ParenExpr struct {
	Expr Expr
}

func (p ParenExpr) Accept(interp Visitor) AvaVal {
	return interp.VisitParenExpr(p)
}

func (p ParenExpr) String() string {
	return fmt.Sprintf("Parens(%s)", p.Expr.String())
}

func (p ParenExpr) exprNode() {}

// Program statement

type ProgStmt struct {
	Loc   LocStmt
	Glbls []GlblStmt
}

func (p ProgStmt) Accept(interp Visitor) AvaVal {
	return interp.VisitProgStmt(p)
}

func (p ProgStmt) String() string {
	sb := strings.Builder{}

	parts := make([]string, 0)

	parts = append(parts, fmt.Sprintf("\t"+p.Loc.String()))

	for _, glbl := range p.Glbls {
		s := glbl.String()
		s = strings.Replace(s, "\n", "\n\t", -1)
		parts = append(parts, "\t"+s)
	}

	body := strings.Join(parts, ",\n")
	sb.WriteString(fmt.Sprintf("Program(\n%s\n)", body))

	return sb.String()
}

// String literal

type StrLit struct {
	Value string
}

func (s StrLit) Accept(interp Visitor) AvaVal {
	return interp.VisitStrLit(s)
}

func (s StrLit) String() string {
	return fmt.Sprintf("StrLit(%s)", s.Value)
}

func (s StrLit) exprNode() {}

// Struct declaration statement

type StructDecl struct {
	Name   string
	Fields []StructField
}

type StructField struct {
	Name string
	Type string
}

func (s StructDecl) String() string {
	fields := make([]string, len(s.Fields))

	for i, f := range s.Fields {
		fields[i] = fmt.Sprintf("%s: %s", f.Name, f.Type)
	}

	return fmt.Sprintf("StructDecl(%s, %s)", s.Name, strings.Join(fields, ", "))
}

func (s StructDecl) Accept(interp Visitor) AvaVal {
	return interp.VisitStructDecl(s)
}

func (s StructDecl) stmtNode() {}

func (s StructDecl) glblStmt() {}

// Variable declaration statement

type VarDecl struct {
	Name string
	Type string
	Init Expr
}

func (v VarDecl) Accept(interp Visitor) AvaVal {
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

func (v VarDecl) declNode() {}

// Variable expression

type Variable struct {
	Name string
}

func (v Variable) Accept(interp Visitor) AvaVal {
	return interp.VisitVariable(v)
}

func (v Variable) String() string {
	return fmt.Sprintf("Variable(%s)", v.Name)
}

func (v Variable) exprNode() {
}

// While statement

type WhileStmt struct {
	Condition Expr
	Body      Block
}

func (w WhileStmt) String() string {
	return fmt.Sprintf("WhileStmt(%s, %s)", w.Condition.String(), w.Body.String())
}

func (w WhileStmt) Accept(interp Visitor) AvaVal {
	return interp.VisitWhileStmt(w)
}

func (w WhileStmt) stmtNode() {}
