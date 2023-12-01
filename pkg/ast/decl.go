package ast

import "fmt"

type FuncDecl struct {
	Name       string
	Body       Block
	ReturnType string
}

func (d FuncDecl) declNode() {}
func (d FuncDecl) String() string {
	return fmt.Sprintf("Func(%s, %s\n%s\n)", d.ReturnType, d.Name, d.Body.String())
}
func (d FuncDecl) IsVoid() bool {
	return d.ReturnType == "void" || d.ReturnType == ""
}
