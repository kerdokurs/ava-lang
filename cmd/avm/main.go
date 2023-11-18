package main

import (
	"fmt"

	"kerdo.dev/ava-lang/pkg/avm"
)

func main() {
	// int x = 0;
	// for (int i = 0; i < count; i++) {
	// 	x += i;
	// }
	static := []byte{
		'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '\n', 0,
	}
	helloWorld := 0
	vm := avm.NewVM()
	vm.Static = static
	loopCmp := vm.Label()
	loopEnd := vm.Label()
	code := []avm.Instruction{
		// count
		{avm.LoadImmediate, 100},
		{avm.Store, int(avm.GPR2)},

		// x
		{avm.LoadImmediate, 0},
		{avm.Store, int(avm.GPR1)},

		// i
		{avm.LoadImmediate, 1},
		{avm.Store, int(avm.GPR0)},

		// i < count
		{avm.Lbl, loopCmp},
		{avm.Load, int(avm.GPR0)},
		{avm.Load, int(avm.GPR2)},
		{avm.Lt, 0},

		// if (i < 10) == 0 goto 10
		{avm.Jz, loopEnd},

		// else
		// load x, add i
		{avm.Load, int(avm.GPR1)},
		{avm.Load, int(avm.GPR0)},
		{avm.Add, 0},

		// store x
		{avm.Store, int(avm.GPR1)},

		// i++
		{avm.Inc, int(avm.GPR0)},

		{avm.PutCStr, helloWorld},
		{avm.NewString, helloWorld},
		{avm.PutStr, 0},

		{avm.Jmp, loopCmp},
		{avm.Lbl, loopEnd},

		// print x
		{avm.Load, int(avm.GPR1)},
		{avm.PutInt, 0},

		// print hello world
		{avm.PutCStr, helloWorld},

		// new string on heap, print
		{avm.NewString, helloWorld},
		{avm.PutStr, 0},

		{avm.Hlt, 0},
	}
	vm.Bytecode = code
	vm.LinkLabels()
	vm.Run()

	fmt.Printf("%+v\n", vm.Labels)
	fmt.Printf("%+v\n", vm.Registers)
	fmt.Printf("%+v\n", vm.Stack)
	fmt.Printf("%+v\n", vm.Heap)
	fmt.Printf("%+v\n", vm.Static)
}
