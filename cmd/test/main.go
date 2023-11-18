package main

import (
	"fmt"
	"strings"

	"kerdo.dev/ava-lang/pkg/avm"
	"kerdo.dev/ava-lang/pkg/lexer"
	"kerdo.dev/ava-lang/pkg/parser"
)

func main() {
	code := `fun test() {
		let x = 1;
		return x + 2;
	}

	fun main() {
		let y = test() - 1;
	}`

	lexer := lexer.New(strings.NewReader(code))
	tokens := lexer.Lex()
	for _, token := range tokens {
		fmt.Printf("%+v\n", token)
	}

	parser := parser.FromTokenStream(tokens)
	program := parser.Parse()
	fmt.Println(program.String())

	avmAssembler := avm.NewAssembler()
	avmAssembler.Prog = program
	vm := avmAssembler.Assemble()
	for i, inst := range vm.Bytecode {
		fmt.Printf("%d -> %s, %d\n", i, avm.EnumToString[inst.Type], inst.Value)
	}

	vm.Run()
	fmt.Printf("%+v\n", vm.Labels)
	fmt.Printf("%+v\n", vm.Registers)
	fmt.Printf("%+v\n", vm.Stack)
	fmt.Printf("%+v\n", vm.Heap)
	fmt.Printf("%+v\n", vm.Static)
}
