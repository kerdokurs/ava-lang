package main

import (
	"fmt"
	"kerdo.dev/ava/ast"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Compiler struct {
	tree ast.ProgStmt

	staticEnv map[string]any
}

func NewCompiler(tree ast.ProgStmt) *Compiler {
	return &Compiler{
		tree:      tree,
		staticEnv: make(map[string]any),
	}
}

func (c *Compiler) Compile(fileName string) {
	base := filepath.Base(fileName)

	source := ".align 2\n\n"

	source += ".equ SYS_write, 64\n"
	source += ".equ SYS_stdout, 1\n"

	c.VisitProgStmt(c.tree)

	source += ".global _start\n"
	source += "_start:\n"
	source += "\tmov x8, SYS_write\n"
	source += "\tmov x0, SYS_stdout\n"
	source += "\tmov x1, const_a\n"
	source += "\tmov x2, 13\n"
	source += "\tsyscall\n"
	source += "\tret\n\n"

	for i, val := range c.staticEnv {
		op := ""
		switch val.(type) {
		case string:
			op = ".string"
		case int8:
			op = ".byte"
		case uint8:
			op = ".byte"
		case int32:
			op = ".int"
		case uint32:
			op = ".int"
		default:
			log.Fatalf("Type %T is currently not supported\n", val)
		}
		source += fmt.Sprintf("const_%s: %s \"%v\"\n", i, op, val)
	}

	f, err := os.Create(fileName)
	if err != nil {
		f, _ = os.Open(fileName)
	}
	f.WriteString(source)
	f.Close()

	// bin: bin.o
	//	ld -o bin bin.o -lSystem -syslibroot `xcrun -sdk macosx --show-sdk-path` -e _start -arch arm64
	//
	//bin.o: bin.s
	//	as -arch arm64 -o bin.o bin.s

	exec.Command("as -arch arm64 -o " + base + ".o " + base + ".asm")
	exec.Command("ld", "-o", base, base+".bin", "-lSystem", " -syslibroot `xcrun -sdk macosx --show-sdk-path`", "-e _start", "-arch arm64")
}

func (c *Compiler) Visit(node ast.Node) any {
	return node.Accept(c)
}

func (c *Compiler) VisitProgStmt(stmt ast.ProgStmt) any {
	for _, glbl := range stmt.Glbls {
		fmt.Printf("Visiting %+v\n", glbl)
		c.Visit(glbl)
	}

	return nil
}

func (c *Compiler) VisitLocStmt(stmt ast.LocStmt) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitParenExpr(expr ast.ParenExpr) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitBlock(block ast.Block) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitIfStmt(stmt ast.IfStmt) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitAssignStmt(stmt ast.AssignStmt) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitExprStmt(stmt ast.ExprStmt) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitFuncCall(call ast.FuncCall) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitFuncDecl(decl ast.FuncDecl) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitConstDecl(decl ast.ConstDecl) any {
	if !decl.IsGlobal {
		log.Fatalf("Non-global const declarations are not supported yet.")
	}

	if val, ok := c.Visit(decl.Init).(int); ok {
		var typedVal any

		switch decl.Type {
		case "u8":
			typedVal = uint8(val)
		default:
			log.Fatalf("Const variable of type %s is not supported yet.\n", decl.Type)
		}

		c.staticEnv[decl.Name] = typedVal
	} else if val, ok := c.Visit(decl.Init).(string); ok {
		var typedVal any

		switch decl.Type {
		case "str":
			typedVal = val
		default:
			log.Fatalf("Const variable of type %s is not supported yet.", decl.Type)
		}

		c.staticEnv[decl.Name] = typedVal
	}

	return nil
}

func (c *Compiler) VisitVarDecl(decl ast.VarDecl) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitVariable(variable ast.Variable) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitIntLit(lit ast.IntLit) any {
	return lit.Value
}

func (c *Compiler) VisitFloatLit(lit ast.FloatLit) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitBoolLit(lit ast.BoolLit) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitNilLit(lit ast.NilLit) any {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitStrLit(lit ast.StrLit) any {
	return lit.Value
}
