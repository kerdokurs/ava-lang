package parser

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"kerdo.dev/ava-lang/pkg/ast"
	"kerdo.dev/ava-lang/pkg/lexer"
)

type Parser struct {
	Tokens []lexer.Token
	index  int
}

func FromTokenStream(tokens []lexer.Token) Parser {
	return Parser{
		Tokens: tokens,
	}
}

func (p *Parser) Parse() ast.ProgStmt {
	decls := make([]ast.Decl, 0)

	for {
		if p.cur().Type == lexer.EOF {
			break
		}

		decl := p.decl()
		decls = append(decls, decl)
	}

	p.done()

	return ast.ProgStmt{
		Decls: decls,
	}
}

func (p *Parser) decl() ast.Decl {
	p.expectAnyKeyword("fun", "let")

	tok := p.consume()

	if tok.Value == "fun" {
		return p.funcDecl()
	}

	fail("unimplemented %s, %s", tok.Name(), tok.Value)
	return nil
}

func (p *Parser) varDecl() ast.VarDecl {
	p.expectAndConsume(lexer.Ident, "let")

	tok := p.expectAndConsume(lexer.Ident, "")
	name := tok.Value

	var expr ast.Expr
	if p.cur().Type != lexer.Semi {
		// has assign expr
		p.expectAndConsume(lexer.Operator, "=")
		expr = p.expr()
	}

	p.expectAndConsume(lexer.Semi, "")

	return ast.VarDecl{
		Name: name,
		Expr: expr,
	}
}

func (p *Parser) funcDecl() ast.FuncDecl {
	// assume `fun` keyword is consumed

	tok := p.expectAndConsume(lexer.Ident, "")
	name := tok.Value

	p.expectAndConsume(lexer.LParen, "")
	p.expectAndConsume(lexer.RParen, "")

	body := p.block()

	return ast.FuncDecl{
		Name: name,
		Body: body,
	}
}

func (p *Parser) block() ast.Block {
	p.expectAndConsume(lexer.LCurly, "")

	stmts := make([]ast.Stmt, 0)
	for p.cur().Type != lexer.RCurly {
		stmts = append(stmts, p.stmt())
	}

	p.expectAndConsume(lexer.RCurly, "")

	return ast.Block{
		Stmts: stmts,
	}
}

func (p *Parser) stmt() ast.Stmt {
	switch p.cur().Value {
	case "let":
		return p.varDecl()
	case "return":
		return p.returnStmt()
	}

	expr := p.expr()
	p.expectAndConsume(lexer.Semi, "")
	return ast.ExprStmt{
		Expr: expr,
	}
}

func (p *Parser) returnStmt() ast.ReturnStmt {
	p.expectAndConsume(lexer.Ident, "return")

	expr := p.expr()
	p.expectAndConsume(lexer.Semi, "")

	return ast.ReturnStmt{
		Expr: expr,
	}
}

func (p *Parser) expr() ast.Expr {
	return p.compExpr()
}

var compOps = []string{
	"<", ">", "<=", ">=", "==", "!=",
	"&", "&&", "|", "||",
}

func (p *Parser) compExpr() ast.Expr {
	l := p.addExpr()

	if p.cur().Type != lexer.Operator {
		return l
	}

	op := p.consume()
	r := p.addExpr()
	return ast.Call{
		Name:         op.Value,
		IsArithmetic: true,
		Args:         []ast.Expr{l, r},
	}
}

// / Additive
// /  : Multiplicative
// /  | Additive (+|-) Multiplicative
func (p *Parser) addExpr() ast.Expr {
	l := p.mulExpr()

	for {
		cur := p.cur()
		if !(cur.Type == lexer.Operator && (cur.Value == "+" || cur.Value == "-")) {
			break
		}

		op := p.consume()

		r := p.mulExpr()
		l = ast.Call{
			Name:         op.Value,
			IsArithmetic: true,
			Args:         []ast.Expr{l, r},
		}
	}

	return l
}

// / Multiplicative
// /  : Literal
// /  | Multiplicative (*|/) Literal
func (p *Parser) mulExpr() ast.Expr {
	l := p.primaryExpr()

	for {
		cur := p.cur()
		if !(cur.Type == lexer.Operator && (cur.Value == "*" || cur.Value == "/")) {
			break
		}

		op := p.consume()

		r := p.primaryExpr()
		l = ast.Call{
			Name:         op.Value,
			IsArithmetic: true,
			Args:         []ast.Expr{l, r},
		}
	}

	return l
}

func (p *Parser) primaryExpr() ast.Expr {
	t := p.cur()

	switch t.Type {
	case lexer.Int:
		return p.intLiteral()
	case lexer.Ident:
		return p.variableOrFunctionCall()
	}

	fmt.Println("unimplemented primaryExpr", t.Type, t.Value)
	os.Exit(1)
	return nil
}

func (p *Parser) variableOrFunctionCall() ast.Expr {
	name := p.consume()

	if p.cur().Type != lexer.LParen {
		return ast.Variable{
			Name: name.Value,
		}
	}

	p.expectAndConsume(lexer.LParen, "")
	p.expectAndConsume(lexer.RParen, "")

	return ast.Call{
		Name: name.Value,
	}
}

func (p *Parser) intLiteral() ast.IntLiteral {
	var value int
	var err error

	tok := p.consume()
	value, err = strconv.Atoi(tok.Value)
	if err != nil {
		panic(err)
	}

	return ast.IntLiteral{
		Value: value,
	}
}

func (p *Parser) expectAnyKeyword(keywords ...string) {
	tok := p.cur()

	for _, keyword := range keywords {
		if tok.Value == keyword {
			return
		}
	}

	fail("expected something: %s, got: %s\n", strings.Join(keywords, ","), tok.Value)
}

func (p *Parser) expectAndConsume(typ lexer.TokenType, value string) lexer.Token {
	p.expect(typ, value)
	return *p.consume()
}

func (p *Parser) expect(typ lexer.TokenType, value string) {
	token := p.cur()
	if token.Type != typ {
		log.Fatalf("Expected type %s, but got type %s at %d\n", lexer.Name(typ), token.Name(), p.index)
	}

	if len(value) > 0 {
		if value != token.Value {
			log.Fatalf("Expected %s, but got %s at %d\n", value, token.Value, p.index)
		}
	}
}

func (p *Parser) done() {
	p.expectAndConsume(lexer.EOF, "")
}

func (p *Parser) cur() *lexer.Token {
	return &p.Tokens[p.index]
}

func fail(fmtString string, args ...any) {
	fmt.Fprintf(os.Stderr, fmtString, args...)
	os.Exit(1)
}

func (p *Parser) consume() *lexer.Token {
	tok := &p.Tokens[p.index]
	p.index++
	return tok
}

func (p *Parser) next() *lexer.Token {
	return &p.Tokens[p.index+1]
}
