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
	a.vm = NewVM()

	for _, decl := range a.Prog.Decls {
		a.assembleDecl(decl)
	}

	entryLbl := a.vm.Label()
	a.bytecode = append(a.bytecode, Instruction{
		Lbl, entryLbl,
	})
	a.bytecode = append(a.bytecode, Instruction{
		Call, a.startLabel,
	})
	a.bytecode = append(a.bytecode, Instruction{
		Hlt, 0,
	})

	a.vm.Bytecode = a.bytecode
	a.vm.LinkLabels()
	a.vm.programCounter = a.vm.Labels[entryLbl]
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
		a.bytecode = append(a.bytecode, Instruction{
			LoadImmediate, 0,
		})
		a.bytecode = append(a.bytecode, Instruction{
			Ret, 0,
		})
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
		a.bytecode = append(a.bytecode, Instruction{
			Pop, 0,
		})
	case ast.WhileStmt:
		loopLbl := a.vm.Label()
		endLbl := a.vm.Label()
		a.bytecode = append(a.bytecode, Instruction{
			Lbl, loopLbl,
		})
		a.assembleExpr(s.Condition)
		a.bytecode = append(a.bytecode, Instruction{
			Jz, endLbl,
		})
		a.assembleBlock(s.Body)
		a.bytecode = append(a.bytecode, Instruction{
			Jmp, loopLbl,
		})
		a.bytecode = append(a.bytecode, Instruction{
			Lbl, endLbl,
		})
	case ast.AssignStmt:
		// compute expr
		// store
		a.assembleExpr(s.Expr)
		varIndex, ok := a.variableIndex[s.Variable.Name]
		if !ok {
			fmt.Println("no variable found with name", s.Variable.Name)
			return
		}

		a.bytecode = append(a.bytecode, Instruction{
			Store, int(GPR0) + varIndex,
		})
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
			if e.Name == "+" || e.Name == "-" || e.Name == "<" {
				a.assembleExpr(e.Args[1])
				a.assembleExpr(e.Args[0])

				op := Add
				if e.Name == "-" {
					op = Sub
				} else if e.Name == "<" {
					op = Lt
				}

				a.bytecode = append(a.bytecode, Instruction{
					op, 0,
				})
			} else if e.Name == "putint" {
				for i := len(e.Args) - 1; i >= 0; i-- {
					a.assembleExpr(e.Args[i])
				}

				a.bytecode = append(a.bytecode, Instruction{
					PutInt, 0,
				})
				a.bytecode = append(a.bytecode, Instruction{
					LoadImmediate, 0,
				})
			} else if e.Name == "+=" {

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
