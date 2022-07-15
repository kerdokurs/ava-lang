package ast

import "fmt"

type StrLit struct {
	Value string
}

func (s StrLit) String() string {
	return fmt.Sprintf("StrLit(%s)", s.Value)
}

func (s StrLit) exprNode() {}
