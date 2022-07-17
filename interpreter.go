package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

type Interp struct {
	tree ProgStmt

	environment *Environment[AvaVar]
	functions   map[string]FunctionDefinition
	structs     map[string]StructDefinition
}

func (i *Interp) VisitExprStmt(stmt ExprStmt) AvaVal {
	return i.Visit(stmt.Expr)
}

func NewInterpretator(source io.Reader) *Interp {
	tree := CreateAst(source)

	if IsDebug {
		fmt.Println(tree.String())
	}

	return &Interp{
		tree:        tree,
		environment: NewEnvironment[AvaVar](),
		functions:   make(map[string]FunctionDefinition),
		structs:     make(map[string]StructDefinition),
	}
}

func (i *Interp) Run() {
	i.VisitProgStmt(i.tree)

	if _, ok := i.functions["main"]; ok {
		funcCall := FuncCall{
			Name: "main",
			Args: []Expr{},
		}
		i.VisitFuncCall(funcCall)
		return
	}

	fmt.Println("Source code does not contain main function.")
	os.Exit(1)
}

func (i *Interp) Visit(node Node) AvaVal {
	return node.Accept(i)
}

func (i *Interp) VisitProgStmt(stmt ProgStmt) AvaVal {
	for _, glbl := range stmt.Glbls {
		i.Visit(glbl)
	}

	return AvaVal{}
}

func (i *Interp) VisitLocStmt(stmt LocStmt) AvaVal {
	panic("Location statements are not yet supported")
}

func (i *Interp) VisitParenExpr(expr ParenExpr) AvaVal {
	return i.Visit(expr.Expr)
}

func (i *Interp) visitArithmeticCall(call FuncCall) AvaVal {
	if len(call.Args) != 2 {
		fmt.Printf("Arithmetic operation requires exactly 2 arguments, but got %d. (Possible parser bug)\n", len(call.Args))
		os.Exit(1)
	}
	args := Map(call.Args, func(arg Expr) AvaVal {
		return i.Visit(arg)
	})

	a := args[0]
	b := args[1]

	if a.Type != b.Type {
		fmt.Printf("Arithmetic operation arguments must be same! Received types %d and %d\n", a.Type, b.Type)
		os.Exit(1)
	}

	if a.Type != Int {
		fmt.Printf("Custom type arithmetic operation is not implemented yet.")
		os.Exit(1)
	}

	aInt, bInt := 0, 0
	var ok bool
	if aInt, ok = a.Value.(int); !ok {
		fmt.Printf("Could not cast value to int in arithmetic operation %s\n", call.Name)
		os.Exit(1)
	}
	if bInt, ok = b.Value.(int); !ok {
		fmt.Printf("Could not cast value to int in arithmetic operation %s\n", call.Name)
		os.Exit(1)
	}

	val := 0
	switch call.Name {
	case "+":
		val = aInt + bInt
	default:
		fmt.Printf("Unsupported arithmetic operation: %s\n", call.Name)
		os.Exit(1)
	}

	return AvaVal{
		Type:  Int,
		Value: val,
	}
}

func (i *Interp) visitComparisonCall(call FuncCall) AvaVal {
	if len(call.Args) != 2 {
		fmt.Printf("Comparison requires exactly 2 arguments, but got %d. (Possible parser bug)\n", len(call.Args))
		os.Exit(1)
	}
	args := Map(call.Args, func(arg Expr) AvaVal {
		return i.Visit(arg)
	})

	a := args[0]
	b := args[1]

	if a.Type != b.Type {
		fmt.Printf("Comparison arguments must be same! Received types %d and %d\n", a.Type, b.Type)
		os.Exit(1)
	}

	if a.Type != Int {
		fmt.Printf("Custom type comparison is not implemented yet.")
		os.Exit(1)
	}

	aInt, bInt := 0, 0
	var ok bool
	if aInt, ok = a.Value.(int); !ok {
		fmt.Printf("Could not cast value to int in comparison %s\n", call.Name)
		os.Exit(1)
	}
	if bInt, ok = b.Value.(int); !ok {
		fmt.Printf("Could not cast value to int in comparison %s\n", call.Name)
		os.Exit(1)
	}

	val := false
	switch call.Name {
	case "<":
		val = aInt < bInt
	case "==":
		val = aInt == bInt
	default:
		fmt.Printf("Unsupported comparison operator: %s\n", call.Name)
		os.Exit(1)
	}

	return AvaVal{
		Type:  Bool,
		Value: val,
	}
}

func (i *Interp) VisitFuncCall(call FuncCall) AvaVal {
	if call.IsArithmetic {
		return i.visitArithmeticCall(call)
	} else if call.IsComparison {
		return i.visitComparisonCall(call)
	}

	if fun, ok := i.functions[call.Name]; ok {
		return i.findAndRunDefinedFunction(call, fun)
	}

	return i.findAndRunBuiltInFunction(call)
}

func (i *Interp) findAndRunDefinedFunction(call FuncCall, def FunctionDefinition) AvaVal {
	if len(call.Args) != len(def.Params) {
		fmt.Printf("Function %s expects %d arguments, but got %d\n", def.Name, len(def.Params), len(call.Args))
		os.Exit(1)
	}

	i.environment.EnterBlock()

	for k, param := range def.Params {
		arg := i.Visit(call.Args[k])
		v := AvaVar{
			Type:  arg.Type,
			Value: arg,
		}
		i.environment.DeclareAssign(param.Name, v)
	}

	returnValue := i.Visit(def.Body)

	i.environment.ExitBlock()

	return returnValue
}

func (i *Interp) findAndRunBuiltInFunction(call FuncCall) AvaVal {
	builtins := reflect.ValueOf(AvaBuiltins{})
	m := builtins.MethodByName(call.Name)

	if !m.IsValid() {
		fmt.Printf("Undefined function %s\n", call.Name)
		os.Exit(1)
	}

	args := Map(call.Args, func(arg Expr) AvaVal {
		return i.Visit(arg)
	})
	//argTypes := Map(args, func(arg AvaVar) reflect.Type {
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

	argValues := Map(args, func(arg AvaVal) reflect.Value {
		return reflect.ValueOf(arg.Value)
	})

	returnType := Void
	if m.Type().NumOut() > 0 {
		t := m.Type().Out(0).String()
		switch t {
		case "int":
			returnType = Int
		case "string":
			returnType = String
		}
	}

	// TODO: Returning a reflect.Value could lead to a problem.
	// It is not a pure type anymore. Deal with it later.
	var result reflect.Value
	results := m.Call(argValues)
	if len(results) > 0 {
		result = results[0]
	}
	return AvaVal{
		Type:  returnType,
		Value: result,
	}
}

func (i *Interp) VisitFuncDecl(decl FuncDecl) AvaVal {
	def := FunctionDefinition{
		Name:   decl.Name,
		Params: decl.Params,
		Body:   decl.Body,
	}
	i.functions[decl.Name] = def
	return AvaVal{
		Type: Void,
	}
}

func (i *Interp) inferType(val AvaVal) (typ AvaType) {
	typeName := reflect.TypeOf(val.Value).String()

	switch typeName {
	case "int":
		typ = Int
	}

	return
}

func (i *Interp) VisitConstDecl(decl ConstDecl) AvaVal {
	val := i.Visit(decl.Init)

	// TODO: This logic could cause issues
	typeName := decl.Type
	typ := val.Type
	if typeName == "" && typ == Unknown {
		typ = i.inferType(val)
	}
	// TODO: ref

	if typ != val.Type {
		fmt.Printf("Constant variable %s declared with type %s, but got expression with type %d\n", decl.Name, decl.Type, typ)
		os.Exit(1)
	}

	v := AvaVar{
		Type:    typ,
		Value:   val,
		IsConst: true,
		//IsRef:   isRef,
	}

	i.environment.DeclareAssign(decl.Name, v)

	return val
}

func (i *Interp) VisitVarDecl(decl VarDecl) AvaVal {
	val := i.Visit(decl.Init)

	// TODO: This logic could cause issues
	typeName := decl.Type
	typ := val.Type
	if typeName == "" && typ == Unknown {
		typ = i.inferType(val)
	}
	// TODO: ref

	if typ != val.Type {
		fmt.Printf("Variable %s declared with type %s, but got expression with type %d\n", decl.Name, decl.Type, typ)
		os.Exit(1)
	}

	v := AvaVar{
		Type:    typ,
		Value:   val,
		IsConst: false,
		//IsRef:   isRef,
	}

	i.environment.DeclareAssign(decl.Name, v)
	return AvaVal{
		Type: Void,
	}
}

func (i *Interp) VisitVariable(variable Variable) AvaVal {
	v := i.environment.Get(variable.Name)
	return v.Value
}

func (i *Interp) VisitIntLit(lit IntLit) AvaVal {
	return AvaVal{
		Type:  Int,
		Value: lit.Value,
	}
}

func (i *Interp) VisitFloatLit(lit FloatLit) AvaVal {
	return AvaVal{
		Type:  Float,
		Value: lit.Value,
	}
}

func (i *Interp) VisitBoolLit(lit BoolLit) AvaVal {
	return AvaVal{
		Type:  Bool,
		Value: lit.Value,
	}
}

func (i *Interp) VisitStrLit(lit StrLit) AvaVal {
	return AvaVal{
		Type:  String,
		Value: lit.Value,
	}
}

func (i *Interp) VisitBlock(block Block) AvaVal {
	i.environment.EnterBlock()
	for _, stmt := range block.Stmts {
		i.Visit(stmt)
	}
	i.environment.ExitBlock()
	return AvaVal{}
}

func (i *Interp) VisitIfStmt(stmt IfStmt) AvaVal {
	cond := i.Visit(stmt.Condition)
	ok := false

	switch cond.Type {
	case Bool:
		ok = cond.Value.(bool)
	}

	if ok {
		return i.Visit(stmt.ThenBody)
	}

	if stmt.HasElse {
		return i.Visit(stmt.ElseBody)
	}

	return AvaVal{}
}

func (i *Interp) VisitWhileStmt(stmt WhileStmt) AvaVal {
	for {
		cond := i.Visit(stmt.Condition)

		if cond.Type != Bool {
			fmt.Printf("Condition must be bool")
			os.Exit(1)
		}
		val := cond.Value.(bool)

		if !val {
			break
		}

		i.Visit(stmt.Body)
	}

	return AvaVal{}
}

func (i *Interp) VisitAssignStmt(stmt AssignStmt) AvaVal {
	// TODO: If not found??

	variable := i.environment.Get(stmt.Variable)
	if variable.Type == Zero {
		fmt.Printf("Variable %s is not declared.\n", stmt.Variable)
		os.Exit(1)
	}

	if variable.IsConst {
		fmt.Printf("Assignment to constant variable %s\n", stmt.Variable)
		os.Exit(1)
	}

	val := i.Visit(stmt.Value)

	if variable.Type != val.Type {
		fmt.Printf("Trying to assign invalid typed value to variable %s\n", stmt.Variable)
	}

	variable.Type = val.Type
	variable.Value = val
	i.environment.Assign(stmt.Variable, variable)

	return AvaVal{
		Type: Void,
	}
}

func (i *Interp) VisitStructDecl(decl StructDecl) AvaVal {
	if _, ok := i.structs[decl.Name]; ok {
		fmt.Printf("Redefining struct %s is not allowed.\n", decl.Name)
		os.Exit(1)
	}

	v := StructDefinition{
		Name: decl.Name,
	}
	i.structs[decl.Name] = v

	return AvaVal{
		Type: Void,
	}
}
