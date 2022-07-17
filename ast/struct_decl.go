package ast

import (
	"fmt"
	"strings"
)

type StructDecl struct {
	Name   string
	Fields []StructField
}

type StructField struct {
	Name string
	Type string
}

func (s StructDecl) String() string {
	fields := make([]string, len(s.Fields))

	for i, f := range s.Fields {
		fields[i] = fmt.Sprintf("%s: %s", f.Name, f.Type)
	}

	return fmt.Sprintf("StructDecl(%s, %s)", s.Name, strings.Join(fields, ", "))
}

func (s StructDecl) Accept(interp Visitor) any {
	return interp.VisitStructDecl(s)
}

func (s StructDecl) stmtNode() {}

func (s StructDecl) glblStmt() {}
