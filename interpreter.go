package main

import (
	"fmt"
	"io"
	"kerdo.dev/ava/ast"
	"kerdo.dev/ava/types"
	"log"
	"os"
	"reflect"
)

type Interp struct {
	tree ast.ProgStmt

	environment *Environment[any]
	functions   map[string]types.FunctionDefinition
}

func (i *Interp) VisitExprStmt(stmt ast.ExprStmt) any {
	return i.Visit(stmt.Expr)
}

func NewInterpretator(source io.Reader) *Interp {
	tree := CreateAst(source)

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

	fmt.Println("Source code does not contain main function.")
	os.Exit(1)
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
	panic("Location statements are not yet supported")
}

func (i *Interp) VisitParenExpr(expr ast.ParenExpr) any {
	return i.Visit(expr.Expr)
}

func (i *Interp) VisitFuncCall(call ast.FuncCall) any {
	var l int
	var r int
	var ok bool

	if call.IsArithmetic {
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
		exprs := Map(call.Args, func(expr ast.Expr) any {
			stmt := ast.ExprStmt{
				Expr: expr,
			}
			return i.VisitExprStmt(stmt)
		})
		fmt.Println(exprs)
		return nil
	default:
		if fun, ok := i.functions[call.Name]; ok {
			return i.findAndRunDefinedFunction(call, fun)
		} else {
			return i.findAndRunBuiltInFunction(call)
		}
	}

	return nil
}

func (i *Interp) findAndRunDefinedFunction(call ast.FuncCall, def types.FunctionDefinition) any {
	if len(call.Args) != len(def.Params) {
		fmt.Printf("Function %s expects %d arguments, but got %d\n", def.Name, len(def.Params), len(call.Args))
		os.Exit(1)
	}

	i.environment.EnterBlock()

	for k, param := range def.Params {
		arg := i.Visit(call.Args[k])
		i.environment.DeclareAssign(param.Name, arg)
	}

	returnValue := i.Visit(def.Body)

	i.environment.ExitBlock()

	return returnValue
}

func (i *Interp) findAndRunBuiltInFunction(call ast.FuncCall) any {
	builtins := reflect.ValueOf(AvaBuiltins{})
	m := builtins.MethodByName(call.Name)

	if !m.IsValid() {
		fmt.Printf("Undefined function %s\n", call.Name)
		os.Exit(1)
	}

	args := Map(call.Args, func(arg ast.Expr) any {
		return i.Visit(arg)
	})
	//argTypes := Map(args, func(arg any) reflect.Type {
	//	return reflect.TypeOf(arg)
	//})

	if len(args) != m.Type().NumIn() {
		fmt.Printf("Function %s expects %d arguments, but received %d.\n", call.Name, len(args), m.Type().NumIn())
		os.Exit(1)
	}

	//var x interface{}
	//interfaceType := reflect.TypeOf(x)
	//for i, argType := range argTypes {
	//	if m.Type().In(i) != interfaceType && m.Type().In(i) != argType {
	//		fmt.Printf("Function %s expected argument %d to be of type %v, but got %v\n", call.Name, i+1, m.Type().String(), argType.String())
	//		os.Exit(1)
	//	}
	//}

	argValues := Map(args, func(arg any) reflect.Value {
		return reflect.ValueOf(arg)
	})

	result := m.Call(argValues)
	return result
}

func (i *Interp) VisitFuncDecl(decl ast.FuncDecl) any {
	def := types.FunctionDefinition{
		Name:   decl.Name,
		Params: decl.Params,
		Body:   decl.Body,
	}
	i.functions[decl.Name] = def
	return nil
}

func (i *Interp) VisitConstDecl(decl ast.ConstDecl) any {
	val := i.Visit(decl.Init)
	i.environment.DeclareAssign(decl.Name, val)

	return val
}

func (i *Interp) VisitVarDecl(decl ast.VarDecl) any {
	val := i.Visit(decl.Init)
	i.environment.DeclareAssign(decl.Name, val)
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

func (i *Interp) VisitBlock(block ast.Block) any {
	i.environment.EnterBlock()
	for _, stmt := range block.Stmts {
		i.Visit(stmt)
	}
	i.environment.ExitBlock()
	return nil
}
