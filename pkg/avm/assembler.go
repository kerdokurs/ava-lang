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

	stringAddrs map[string]int

	startLabel int

	scopeStack []*Scope
}

func NewAssembler() *Assembler {
	return &Assembler{
		funcLabels:    make(map[string]int),
		variableIndex: make(map[string]int),
		stringAddrs:   make(map[string]int),
		scopeStack:    make([]*Scope, 0),
	}
}

func (a *Assembler) Assemble() *AVM {
	a.vm = NewVM()
	a.vm.Static = make([]byte, 0)

	for _, decl := range a.Prog.Decls {
		a.assembleDecl(decl)
	}

	entryLbl := a.vm.Label()
	a.emit(Lbl, entryLbl)
	a.emit(Call, a.startLabel)
	a.emit(Store, int(GPR0))
	a.emit(Hlt)

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

		localCount := d.Body.CountLocals()
		scope := a.Scope()
		scope.LocalCount = localCount
		a.assemblePushFrame(scope)

		fnStart := len(a.bytecode)

		a.assembleBlock(d.Body)

		cleanupLbl := a.vm.Label()
		a.emit(Lbl, cleanupLbl)
		a.assemblePopFrame(scope)
		a.PopScope()

		for i := fnStart; i < len(a.bytecode); i++ {
			if a.bytecode[i].Type == Ret {
				a.bytecode[i] = Instruction{Jmp, cleanupLbl}
			}
		}

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

func (a *Assembler) assemblePushFrame(scope *Scope) {
	// Frame description:
	// GPR0 - return value
	// [0] - old FP
	// [0..n] - arguments, only if function body
	// [n+1..m] - local variables

	// store old FP
	a.emit(Load, int(FP))
	a.emit(Load, int(FP))
	a.emit(LoadImmediate, scope.LocalCount+1)
	a.emit(Add, 0)
	a.emit(Store, int(FP))
	for i := 0; i < scope.LocalCount; i++ {
		a.emit(LoadImmediate, 0)
	}
	// a.emit(Store, int(GPR2))
	// // GPR2 = old fp
	// a.emit(Inc, int(GPR2))
	// a.emit(Load, int(GPR2))
}

func (a *Assembler) assemblePopFrame(scope *Scope) {
	// put push GPR0 to stack
	a.emit(Store, int(GPR0))

	for i := 0; i < scope.LocalCount; i++ {
		a.emit(Pop)
	}

	a.emit(Store, int(FP))
}

func (a *Assembler) assembleStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case ast.VarDecl:
		scope := a.currentScope()
		varIndex := len(scope.LocalVariableIndices)
		scope.LocalVariableIndices[s.Name] = varIndex
		fmt.Printf("%s got index %d\n", s.Name, varIndex)
		if s.HasInit {
			a.assembleExpr(s.Expr)
			a.bytecode = append(a.bytecode, Instruction{
				StoreA, varIndex,
			})
		}
	case ast.ReturnStmt:
		a.assembleExpr(s.Expr)
		a.bytecode = append(a.bytecode, Instruction{
			Ret, 0,
		})
	case ast.ExprStmt:
		a.assembleExpr(s.Expr)
		a.emit(Pop)
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
		scope := a.currentScope()
		varIndex, ok := scope.LocalVariableIndices[s.Variable.Name]
		if !ok {
			fmt.Println("no variable found with name", s.Variable.Name)
			return
		}

		a.bytecode = append(a.bytecode, Instruction{
			StoreA, int(GPR0) + varIndex,
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
		scope := a.currentScope()
		varOffset, ok := scope.LocalVariableIndices[e.Name]
		if !ok {
			fmt.Printf("no variable named %s in current scope\n", e.Name)
			return
		}

		a.emit(LoadA, varOffset)
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
			} else if e.Name == "putstr" {
				a.assembleExpr(e.Args[0])
				a.bytecode = append(a.bytecode, Instruction{
					PutCStr, 0,
				})
				a.bytecode = append(a.bytecode, Instruction{
					LoadImmediate, 0,
				})
			} else if e.Name == "+=" {
				panic("assembling += not implemented")
			} else {
				fmt.Println("no function named", e.Name)
			}
			return
		}

		a.bytecode = append(a.bytecode, Instruction{
			Call, funcLabel,
		})
	case ast.StringLiteral:
		addr := len(a.vm.Static)
		for _, ch := range e.Value {
			a.vm.Static = append(a.vm.Static, byte(ch))
		}
		a.vm.Static = append(a.vm.Static, 0)
		a.bytecode = append(a.bytecode, Instruction{
			Load, addr,
		})
	default:
		fmt.Println("unsupported expr to assemble", e)
	}
}

func (a *Assembler) emit(code InstructionType, arg ...int) {
	opArg := 0
	if len(arg) > 0 {
		opArg = arg[0]
	}
	a.bytecode = append(a.bytecode, Instruction{
		code, opArg,
	})
}

type Scope struct {
	LocalCount           int
	LocalVariableIndices map[string]int
}

func (a *Assembler) Scope() *Scope {
	scope := new(Scope)
	scope.LocalVariableIndices = make(map[string]int)
	a.scopeStack = append(a.scopeStack, scope)
	return scope
}

func (a *Assembler) PopScope() {
	a.scopeStack = a.scopeStack[:len(a.scopeStack)-1]
}

func (a *Assembler) currentScope() *Scope {
	return a.scopeStack[len(a.scopeStack)-1]
}
