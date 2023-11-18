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
	LoadImmediate
	Load
	Store
	Add
	Sub
	Lt
	Jmp
	Jz
	Inc
	Lbl
	Call
	Ret

	PutInt
	PutCStr

	NewString
	PutStr
)

var EnumToString = map[InstructionType]string{
	Hlt:           "HALT",
	LoadImmediate: "LOAD IMMEDIATE",
	Load:          "LOAD",
	Store:         "STORE",
	Add:           "ADD",
	Sub:           "SUB",
	Lt:            "LT",
	Jmp:           "JMP",
	Jz:            "JZ",
	Inc:           "INC",
	Lbl:           "LBL",
	Call:          "CALL",
	Ret:           "RET",
}

func (i *Instruction) Execute(vm *AVM) {
	switch i.Type {
	case Hlt:
		fmt.Println("HALT")
	case LoadImmediate:
		vm.push(i.Value)
	case Load:
		reg := Register(i.Value)
		vm.push(vm.Registers[reg])
	case Store:
		value := vm.pop()
		reg := Register(i.Value)
		vm.Registers[reg] = value
	case Add:
		a := vm.pop()
		b := vm.pop()
		vm.push(a + b)
	case Sub:
		a := vm.pop()
		b := vm.pop()
		vm.push(a - b)
	case Lt:
		rhs := vm.pop()
		lhs := vm.pop()
		var value int
		if lhs < rhs {
			value = 1
		}
		vm.push(value)
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
	default:
		fmt.Printf("unsupported operation: %v\n", i)
		os.Exit(1)
	}
}
