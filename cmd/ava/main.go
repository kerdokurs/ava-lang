package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"kerdo.dev/ava-lang/pkg/analyser"
	"kerdo.dev/ava-lang/pkg/avm"
	"kerdo.dev/ava-lang/pkg/lexer"
	"kerdo.dev/ava-lang/pkg/parser"
)

var srcPath = flag.String("src", "main.ava", "source file")
var stdIn = flag.String("stdin", "", "stdin")
var stdFile = flag.String("stdfile", "", "stdfile")

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

	analyser := analyser.NewAnalyser()
	analyser.Program = &program
	err = analyser.Analyse()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Program analysis failed. Exiting.")
		os.Exit(1)
		return
	}
	fmt.Println(program.String())

	avmAssembler := avm.NewAssembler()
	avmAssembler.Prog = program
	vm := avmAssembler.Assemble()
	for i, inst := range vm.Bytecode {
		fmt.Printf("%d -> %s, %d\n", i, avm.InstEnumToString[inst.Type], inst.Value)
	}

	if stdIn != nil && *stdIn != "" {
		vm.StdIn = make([]byte, len(*stdIn))
		for i, c := range *stdIn {
			vm.StdIn[i] = byte(c)
		}
	} else if stdFile != nil && *stdFile != "" {
		f, err := os.Open(*stdFile)
		if err != nil {
			panic(err)
		}
		data, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		f.Close()
		vm.StdIn = data
	} else {
		vm.StdIn = make([]byte, 0)
	}

	var exitCode int
	fmt.Println("Program output:")
	exitCode = vm.Run()
	fmt.Println()

	fmt.Printf("Exit code: %v\n", exitCode)
	fmt.Printf("Labels: %+v\n", vm.Labels)
	fmt.Printf("Registers: %+v\n", vm.Registers)
	fmt.Printf("Stack: %+v\n", vm.Stack)
	fmt.Printf("Heap: %+v\n", vm.Heap)
	fmt.Printf("Static: %+v\n", vm.Static)
}
