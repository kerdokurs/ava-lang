package ast

import (
	"fmt"
	"strings"
)

type FuncCall struct {
	Name string
	Args []Expr
}

func (f FuncCall) String() string {
	strs := make([]string, len(f.Args))
	for i, arg := range f.Args {
		strs[i] = arg.String()
	}
	args := strings.Join(strs, ", ")
	if len(args) > 0 {
		args = ", " + args
	}

	return fmt.Sprintf("FuncCall(%s%s)", f.Name, args)
}

func (f FuncCall) exprNode() {}
