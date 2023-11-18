package avm

import (
	"fmt"

	"kerdo.dev/ava-lang/pkg/ast"
)

type Assembler struct {
	Prog ast.ProgStmt

	vm       *AVM
	bytecode []Instruction

	funcLabels    map[string]int
	variableIndex map[string]int

	startLabel int
}

func NewAssembler() *Assembler {
	return &Assembler{
		funcLabels:    make(map[string]int),
		variableIndex: make(map[string]int),
	}
}

func (a *Assembler) Assemble() *AVM {
	a.vm = New()

	for _, decl := range a.Prog.Decls {
		a.assembleDecl(decl)
	}

	a.vm.Bytecode = a.bytecode
	a.vm.LinkLabels()
	a.vm.programCounter = a.vm.Labels[a.startLabel]
	return a.vm
}

func (a *Assembler) assembleDecl(decl ast.Decl) {
	switch d := decl.(type) {
	case ast.FuncDecl:
		funcLabel := a.vm.Label()
		if d.Name == "main" {
			a.startLabel = funcLabel
		}

		a.funcLabels[d.Name] = funcLabel
		a.bytecode = append(a.bytecode, Instruction{
			Lbl, funcLabel,
		})
		a.assembleBlock(d.Body)
	default:
		fmt.Println("unsupported decl to assemble", d)
	}
}

func (a *Assembler) assembleBlock(block ast.Block) {
	for _, stmt := range block.Stmts {
		a.assembleStmt(stmt)
	}
}

func (a *Assembler) assembleStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case ast.VarDecl:
		a.assembleExpr(s.Expr)
		varIndex := len(a.variableIndex)
		a.variableIndex[s.Name] = varIndex
		a.bytecode = append(a.bytecode, Instruction{
			Store, int(GPR0) + varIndex,
		})
	case ast.ReturnStmt:
		a.assembleExpr(s.Expr)
		a.bytecode = append(a.bytecode, Instruction{
			Ret, 0,
		})
	case ast.ExprStmt:
		a.assembleExpr(s.Expr)
	default:
		fmt.Println("unsupported stmt to assemble", s)
	}
}

func (a *Assembler) assembleExpr(expr ast.Expr) {
	switch e := expr.(type) {
	case ast.IntLiteral:
		a.bytecode = append(a.bytecode, Instruction{
			LoadImmediate, e.Value,
		})
	case ast.Variable:
		varIndex, ok := a.variableIndex[e.Name]
		if !ok {
			fmt.Println("no variable named", e.Name)
			return
		}

		a.bytecode = append(a.bytecode, Instruction{
			Load, int(GPR0) + varIndex,
		})
	case ast.Call:
		funcLabel, ok := a.funcLabels[e.Name]
		if !ok {
			if e.Name == "+" || e.Name == "-" {
				a.assembleExpr(e.Args[1])
				a.assembleExpr(e.Args[0])

				op := Add
				if e.Name == "-" {
					op = Sub
				}

				a.bytecode = append(a.bytecode, Instruction{
					op, 0,
				})
			} else {
				fmt.Println("no function named", e.Name)
			}
			return
		}

		a.bytecode = append(a.bytecode, Instruction{
			Call, funcLabel,
		})
	default:
		fmt.Println("unsupported expr to assemble", e)
	}
}
