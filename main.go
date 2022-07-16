package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var IsVerbose bool

func printHelp() {
	fmt.Println(`Ava usage:
	- help - prints this help message
	- com <file> - compiles given file
	  -out - specify output path (default "a.out")
	  -verbose - print out debug information of compilation (default false)
	- run <file> - interprets given file
	- version - prints version`)
	os.Exit(0)
}

func printVersion() {
	fmt.Println("Ava alpha v0.2 arm64")
	os.Exit(0)
}

func getFile(fileName string) (*os.File, error) {
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return os.Open(fileName)
}

func runInterpreter(fileName string) {
	if IsVerbose {
		fmt.Printf("Starting interpretation on %s\n", fileName)
	}

	file, err := getFile(fileName)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", fileName, err.Error())
		os.Exit(1)
	}

	interp := NewInterpretator(file)
	interp.Run()
}

func runCompilation(fileName string, outPath *string) {
}

func main() {
	outPath := flag.String("out", "a.out", "Output file (only allowed with compile mode)")
	verbose := flag.Bool("verbose", false, "Verbose")

	flag.Parse()

	IsVerbose = *verbose

	args := flag.Args()

	if len(args) < 1 {
		printHelp()
		return
	} else if len(args) == 1 {
		switch args[0] {
		case "help":
			printHelp()
		case "version":
			printVersion()
		}
	} else if len(args) == 2 {
		fileName := args[1]
		switch args[0] {
		case "run":
			runInterpreter(fileName)
		case "com":
			runCompilation(fileName, outPath)
		default:
			fmt.Printf("Invalid command: %s\n", args[0])
			os.Exit(1)
		}
		os.Exit(0)
	}

	fmt.Println("Invalid number of arguments.")
	os.Exit(1)
	//reader, _ := os.Open("com.ava")

	//defer func(reader *os.File) {
	//	err := reader.Close()
	//	if err != nil {
	//		log.Printf("Could not close file: %v\n", err)
	//	}
	//}(reader)
	//
	//lexer := NewLexer(reader)
	//tokens := lexer.ReadAllTokens()
	//for i, token := range tokens {
	//	fmt.Printf("%d -> %v\n", i, token)
	//}
	//
	//parser := NewParser(tokens)
	//node := parser.Parse()
	//fmt.Printf("%s\n", node.String())
	//
	////interp := NewInterpretator(node)
	////interp.Compile()
	//
	//compiler := NewCompiler(node)
	//compiler.Compile("com.asm")
}
