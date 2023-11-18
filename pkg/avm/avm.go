package avm

import "fmt"

type Register int

const (
	GPR0 Register = iota
	GPR1
	GPR2
	GPR3

	RegisterCount
)

type AVM struct {
	Registers [RegisterCount]int

	Bytecode       []Instruction
	programCounter int

	Labels []int

	Stack []int

	Heap []byte

	Static []byte

	callStack []int
}

func New() *AVM {
	vm := &AVM{
		Stack:     make([]int, 0),
		Heap:      make([]byte, 0),
		callStack: make([]int, 0),
	}
	return vm
}

func (vm *AVM) Label() int {
	lblIndex := len(vm.Labels)
	vm.Labels = append(vm.Labels, -1)
	return lblIndex
}

// Returns a relative pointer to Heap start
func (vm *AVM) Alloc(count int) int {
	start := len(vm.Heap)

	for i := 0; i < count; i++ {
		vm.Heap = append(vm.Heap, 0)
	}

	return start
}

func (vm *AVM) WriteHeap(ptr int, data []byte) {
	for i := 0; i < len(data); i++ {
		vm.Heap[ptr+i] = data[i]
	}
}

func (vm *AVM) staticChunkSize(ptr int) int {
	size := 0

	for ptr+size < len(vm.Static) && vm.Static[ptr+size] != 0 {
		size++
	}

	return size
}

// Returns heap-relative pointer
func (vm *AVM) NewString(cstrPtr int) int {
	size := vm.staticChunkSize(cstrPtr) + 1
	ptr := vm.Alloc(size)

	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = vm.Static[i+cstrPtr]
	}
	vm.WriteHeap(ptr, data)

	return ptr
}

func (vm *AVM) LinkLabels() {
	lblCounter := 0
	for i, inst := range vm.Bytecode {
		if inst.Type == Lbl {
			vm.Labels[lblCounter] = i
			vm.Bytecode[i].Value = i
			lblCounter++
		}
	}
}

func (vm *AVM) Run() {
	fmt.Println("starting VM execution at", vm.programCounter)

	for ; vm.programCounter < len(vm.Bytecode); vm.programCounter += 1 {
		instruction := &vm.Bytecode[vm.programCounter]

		instruction.Execute(vm)
		fmt.Println(vm.Stack)
	}
}

func (vm *AVM) SetStart(start int) {
	vm.programCounter = start
}

func (vm *AVM) push(value int) {
	vm.Stack = append(vm.Stack, value)
}

func (vm *AVM) pop() int {
	value := vm.Stack[len(vm.Stack)-1]
	vm.Stack = vm.Stack[:len(vm.Stack)-1]
	return value
}
