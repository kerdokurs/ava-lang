package main

import (
	"flag"
	"fmt"
	"os"

	"kerdo.dev/ava-lang/pkg/avm"
	"kerdo.dev/ava-lang/pkg/lexer"
	"kerdo.dev/ava-lang/pkg/parser"
)

var srcPath = flag.String("src", "main.ava", "source file")

func main() {
	flag.Parse()

	sourceFile, err := os.Open(*srcPath)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	defer sourceFile.Close()

	lexer := lexer.New(sourceFile)
	tokens := lexer.Lex()
	for i, token := range tokens {
		fmt.Printf("%d -> %+v\n", i, token)
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

	fmt.Println("Program output:")
	vm.Run()
	fmt.Println()

	fmt.Printf("Labels: %+v\n", vm.Labels)
	fmt.Printf("Registers: %+v\n", vm.Registers)
	fmt.Printf("Stack: %+v\n", vm.Stack)
	fmt.Printf("Heap: %+v\n", vm.Heap)
	fmt.Printf("Static: %+v\n", vm.Static)
}
