package main

import (
	"kerdo.dev/ava/ast"
	"kerdo.dev/ava/types"
	"log"
)

type Interp struct {
	tree ast.ProgStmt

	environment *Environment[any]
	functions   map[string]types.FunctionDefinition
}

func NewInterpretator(tree ast.ProgStmt) *Interp {
	return &Interp{
		tree:        tree,
		environment: NewEnvironment[any](),
		functions:   make(map[string]types.FunctionDefinition),
	}
}

func (i *Interp) Run() {
	i.VisitProgStmt(i.tree)

	if _, ok := i.functions["main"]; ok {
		funcCall := ast.FuncCall{
			Name: "main",
			Args: []ast.Expr{},
		}
		i.VisitFuncCall(funcCall)
		return
	}

	log.Fatalf("Main function must be defined.")
}

func (i *Interp) Visit(node ast.Node) any {
	return node.Accept(i)
}

func (i *Interp) VisitProgStmt(stmt ast.ProgStmt) any {
	for _, glbl := range stmt.Glbls {
		i.Visit(glbl)
	}

	return 0
}

func (i *Interp) VisitLocStmt(stmt ast.LocStmt) any {
	//TODO implement me
	panic("implement me")
}

func (i *Interp) VisitParenExpr(expr ast.ParenExpr) any {
	return i.Visit(expr.Expr)
}

func (i *Interp) VisitFuncCall(call ast.FuncCall) any {
	var l int
	var r int
	var ok bool

	if len(call.Args) > 0 {
		if l, ok = i.Visit(call.Args[0]).(int); !ok {
			log.Fatalf("Functions can only be called with integer types.\n")
		}
	}

	if len(call.Args) > 1 {
		if r, ok = i.Visit(call.Args[1]).(int); !ok {
			log.Fatalf("Functions can only be called with integer types.\n")
		}
	} else {
		r = l
		l = 0
	}

	switch call.Name {
	case "+":
		return l + r
	case "-":
		return l - r
	case "*":
		return l * r
	case "/":
		return r / l
	case "print":
		return 0
	default:
		if fun, ok := i.functions[call.Name]; ok {
			log.Printf("Running function %s\n", fun.Name)
		} else {
			log.Fatalf("Function with name %s was not found\n", call.Name)
		}
	}

	return nil
}

func (i *Interp) VisitFuncDecl(decl ast.FuncDecl) any {
	def := types.FunctionDefinition{
		Name: decl.Name,
	}
	i.functions[decl.Name] = def
	return nil
}

func (i *Interp) VisitConstDecl(decl ast.ConstDecl) any {
	val := i.Visit(decl.Init)
	i.environment.DeclareAssign(decl.Name, &val)

	return val
}

func (i *Interp) VisitVarDecl(decl ast.VarDecl) any {
	val := i.Visit(decl.Init)
	i.environment.DeclareAssign(decl.Name, &val)
	return val
}

func (i *Interp) VisitVariable(variable ast.Variable) any {
	return i.environment.Get(variable.Name)
}

func (i *Interp) VisitIntLit(lit ast.IntLit) any {
	return lit.Value
}

func (i *Interp) VisitFloatLit(lit ast.FloatLit) any {
	return lit.Value
}

func (i *Interp) VisitBoolLit(lit ast.BoolLit) any {
	return lit.Value
}

func (i *Interp) VisitNilLit(lit ast.NilLit) any {
	return nil
}

func (i *Interp) VisitStrLit(lit ast.StrLit) any {
	return lit.Value
}
