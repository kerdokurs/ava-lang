package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	reader, _ := os.Open("dev.ava")
	defer func(reader *os.File) {
		err := reader.Close()
		if err != nil {
			log.Printf("Could not close file: %v\n", err)
		}
	}(reader)

	lexer := NewLexer(reader)
	tokens := lexer.ReadAllTokens()
	for i, token := range tokens {
		fmt.Printf("%d -> %v\n", i, token)
	}

	parser := NewParser(tokens)
	node := parser.Parse()
	fmt.Printf("%s\n", node.String())

	interp := NewInterpretator(node)
	interp.Run()
}
