package main

import "kerdo.dev/ava/ast"

type Visitor interface {
	visitProgramStmt(stmt ast.ProgStmt) any
}
