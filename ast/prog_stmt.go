package ast

import (
	"fmt"
	"strings"
)

type ProgStmt struct {
	Loc   LocStmt
	Glbls []GlblStmt
}

func (p ProgStmt) Accept(interp Visitor) any {
	return interp.VisitProgStmt(p)
}

func (p ProgStmt) String() string {
	sb := strings.Builder{}

	parts := make([]string, 0)

	parts = append(parts, fmt.Sprintf("\t"+p.Loc.String()))

	for _, glbl := range p.Glbls {
		s := glbl.String()
		s = strings.Replace(s, "\n", "\n\t", -1)
		parts = append(parts, "\t"+s)
	}

	body := strings.Join(parts, ",\n")
	sb.WriteString(fmt.Sprintf("Program(\n%s\n)", body))

	return sb.String()
}
