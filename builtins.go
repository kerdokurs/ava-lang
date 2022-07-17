package main

import "fmt"

type AvaBuiltins struct {
}

func (AvaBuiltins) Print(args ...any) {
	fmt.Println(args...)
}

func (AvaBuiltins) Input() string {
	var str string
	fmt.Scanln(&str)
	return str
}
