package analyser

import (
	"fmt"

	"kerdo.dev/ava-lang/pkg/ast"
)

type Analyser struct {
	Program *ast.ProgStmt
}

func NewAnalyser() *Analyser {
	return &Analyser{}
}

func (a *Analyser) Analyse() error {
	for _, stmt := range a.Program.Decls {
		if fnd, ok := stmt.(ast.FuncDecl); ok {
			_, err := a.analyseFuncDecl(fnd)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Analyser) analyseFuncDecl(d ast.FuncDecl) (string, error) {
	return a.analyseBlock(d.Body)
}

func (a *Analyser) analyseBlock(b ast.Block) (string, error) {
	for _, stmt := range b.Stmts {
		if vdl, ok := stmt.(ast.VarDecl); ok {
			_, err := a.analyseVarDecl(&vdl)
			if err != nil {
				return "", err
			}
		}
	}

	return "", nil
}

func (a *Analyser) analyseVarDecl(d *ast.VarDecl) (string, error) {
	isTyped := d.VarType != ""

	if !isTyped && !d.HasInit {
		return "", fmt.Errorf("variable declaration without type needs to have immediate expression")
	}

	if d.HasInit {
		exprType, err := a.analyseExpr(d.Expr)
		if err != nil {
			return exprType, err
		}

		if isTyped && d.VarType != exprType {
			return "", fmt.Errorf("variable is typed as %s but has immediate expression with type %s", d.VarType, exprType)
		} else if !isTyped {
			d.VarType = exprType
		}
	}

	return "", nil
}

func (a *Analyser) analyseExpr(e ast.Expr) (string, error) {
	switch e := e.(type) {
	case ast.IntLiteral:
		return "int", nil
	default:
		fmt.Println("unimplemented type deduction for expr", e)
		return "void", nil
	}
}
