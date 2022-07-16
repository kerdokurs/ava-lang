package main

import "fmt"

type AvaBuiltins struct {
}

func (AvaBuiltins) Print(args ...any) {
	fmt.Println(args...)
}
