package ast

import "fmt"

type FuncDecl struct {
	Name string
	Body Block
}

func (d FuncDecl) declNode() {}
func (d FuncDecl) String() string {
	return fmt.Sprintf("Func(%s\n%s\n)", d.Name, d.Body.String())
}
