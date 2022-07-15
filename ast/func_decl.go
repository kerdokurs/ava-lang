package ast

import (
	"fmt"
	"strings"
)

type FuncDecl struct {
	Name       string
	ReturnType string
	Params     []FuncParam
}

type FuncParam struct {
	Name string
	Type string
}

func (f FuncDecl) String() string {
	retType := f.ReturnType
	if len(retType) == 0 {
		retType = "?"
	}

	params := make([]string, len(f.Params))
	for i, param := range f.Params {
		s := param.String()
		s = strings.Replace(s, "\n", "\n\t", -1)
		params[i] = "\t" + param.String()
	}
	paramsStr := strings.Join(params, ",\n")

	return fmt.Sprintf("FuncDecl(\n\t%s,\n\t%s,\n%s\n)", f.Name, retType, paramsStr)
}

func (f FuncDecl) stmtNode() {}

func (f FuncDecl) glblStmt() {}

func (f FuncParam) String() string {
	return fmt.Sprintf("FuncParam(%s, %s)", f.Name, f.Type)
}
