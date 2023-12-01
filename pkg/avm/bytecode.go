package avm

import (
	"fmt"
	"os"
)

type InstructionType uint

type Instruction struct {
	Type  InstructionType
	Value int
}

type InstructionExecutionFunc func(vm *AVM) int

const (
	Hlt InstructionType = iota
	Trap

	LoadImmediate
	Load
	LoadA
	Store
	StoreA
	Pop
	Add
	Sub
	Mul
	Div
	Mod
	Lt
	Gt
	Eq
	Jmp
	Jz
	Inc
	Lbl
	Call
	Ret

	Not

	PutInt
	PutCStr

	GetInt

	NewString
	PutStr
)

var InstEnumToString = map[InstructionType]string{
	Hlt:           "HALT",
	Trap:          "TRAP",
	LoadImmediate: "LOAD IMMEDIATE",
	Load:          "LOAD",
	LoadA:         "LOADA",
	Store:         "STORE",
	StoreA:        "STOREA",
	Pop:           "POP",
	Add:           "ADD",
	Sub:           "SUB",
	Mul:           "MUL",
	Div:           "DIV",
	Mod:           "MOD",
	Lt:            "LT",
	Gt:            "GT",
	Eq:            "EQ",
	Jmp:           "JMP",
	Jz:            "JZ",
	Inc:           "INC",
	Lbl:           "LBL",
	Call:          "CALL",
	Ret:           "RET",
	Not:           "NOT",
	PutInt:        "PUT INT",
	PutCStr:       "PUT CSTR",
	GetInt:        "GET INT",
}

func (i *Instruction) Execute(vm *AVM) {
	switch i.Type {
	case Hlt:
		fmt.Println("HALT")
	case Trap:
		fmt.Println("TRAP")
		vm.programCounter = len(vm.Bytecode) - 1
	case LoadImmediate:
		vm.push(i.Value)
	case Load:
		reg := Register(i.Value)
		vm.push(vm.Registers[reg])
	case LoadA:
		offset := i.Value
		localLocation := vm.Registers[FP] - offset - 1
		vm.push(vm.Stack[localLocation])
	case Store:
		value := vm.pop()
		reg := Register(i.Value)
		vm.Registers[reg] = value
	case StoreA:
		offset := i.Value
		localLocation := vm.Registers[FP] - offset - 1
		val := vm.pop()
		vm.Stack[localLocation] = val
	case Pop:
		vm.pop()
	case Add:
		a := vm.pop()
		b := vm.pop()
		vm.push(a + b)
	case Sub:
		a := vm.pop()
		b := vm.pop()
		vm.push(a - b)
	case Mul:
		a := vm.pop()
		b := vm.pop()
		vm.push(a * b)
	case Div:
		a := vm.pop()
		b := vm.pop()
		vm.push(a / b)
	case Mod:
		a := vm.pop()
		b := vm.pop()
		vm.push(a % b)
	case Lt:
		lhs := vm.pop()
		rhs := vm.pop()
		var value int
		if lhs < rhs {
			value = 1
		}
		vm.push(value)
	case Gt:
		lhs := vm.pop()
		rhs := vm.pop()
		var value int
		if lhs > rhs {
			value = 1
		}
		vm.push(value)
	case Eq:
		lhs := vm.pop()
		rhs := vm.pop()
		if lhs == rhs {
			vm.push(1)
		} else {
			vm.push(0)
		}
	case Jmp:
		addr := vm.Labels[i.Value]
		vm.programCounter = addr - 1
	case Jz:
		value := vm.pop()
		if value == 0 {
			addr := vm.Labels[i.Value]
			vm.programCounter = addr - 1
		}
	case Call:
		vm.callStack = append(vm.callStack, vm.programCounter)
		addr := vm.Labels[i.Value]
		vm.programCounter = addr - 1
	case Ret:
		pc := vm.callStack[len(vm.callStack)-1]
		vm.callStack = vm.callStack[:len(vm.callStack)-1]
		vm.programCounter = pc

		vm.push(vm.Registers[GPR0])
	case Not:
		value := vm.pop()
		if value == 0 {
			vm.push(1)
		} else {
			vm.push(0)
		}
	case Inc:
		reg := Register(i.Value)
		vm.Registers[reg]++
	case Lbl:
	case PutInt:
		fmt.Print(vm.pop())
	case PutCStr:
		for j := 0; i.Value+j < len(vm.Static) && vm.Static[i.Value+j] != 0; j++ {
			val := vm.Static[i.Value+j]
			fmt.Printf("%c", rune(val))
		}
	case NewString:
		strPtr := vm.NewString(i.Value)
		vm.push(strPtr)
	case PutStr:
		strPtr := vm.pop()
		for j := 0; j+strPtr < len(vm.Heap) && vm.Heap[j+strPtr] != 0; j++ {
			val := vm.Heap[j+strPtr]
			fmt.Printf("%c", rune(val))
		}
	case GetInt:
		if vm.StdInPtr == len(vm.StdIn) {
			vm.push(69)
		} else {
			b := vm.StdIn[vm.StdInPtr]
			vm.StdInPtr++
			s := string(b)
			vm.push(int(s[0]))
			// 	if s == "\n" {
			// 		vm.push(10000)
			// 	} else {
			// 		if v, err := strconv.Atoi(s); err != nil {
			// 			panic(err)
			// 		} else {
			// 			vm.push(v)
			// 		}
			// 	}
		}
	default:
		fmt.Printf("unsupported operation: %v\n", i)
		os.Exit(1)
	}
}
