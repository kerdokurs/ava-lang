package main

import (
	. "kerdo.dev/ava/ast"
	"log"
	"math/big"
	"strconv"
	"strings"
)

type Parser struct {
	tokens []Token
	i      int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens, i: 0,
	}
}

func (p *Parser) Parse() ProgStmt {
	loc := p.locStmt()

	glblStmts := make([]GlblStmt, 0)
	for {
		if !p.isOfAnyType([]TokenType{KEYWORD}) {
			break
		}

		glblStmt := p.glblStmt()
		glblStmts = append(glblStmts, glblStmt)
	}

	p.done()

	prog := ProgStmt{
		Loc:   loc,
		Glbls: glblStmts,
	}

	return prog
}

func (p *Parser) glblStmt() GlblStmt {
	p.expectAny([]string{"fun", "const", "var"})
	t := p.consume()

	if t.Data == "fun" {
		return p.funcDecl()
	} else if t.Data == "const" {
		return p.constDecl()
	} else if t.Data == "var" {
		return p.varDecl()
	}

	return FuncDecl{}
}

func (p *Parser) stmt() Stmt {
	t := p.cur()

	if t.Data == "var" {
		p.consume()
		return p.varDecl()
	}

	expr := p.expr()
	p.expectAndConsume(SEMI, "")
	return ExprStmt{
		Expr: expr,
	}
}

func (p *Parser) varDecl() VarDecl {
	t := p.expectAndConsume(IDENT, "")
	name := t.Data

	typ := p.varType()

	p.expectAndConsume(OPERATOR, "=")

	var init Expr
	if p.cur().Type != SEMI {
		init = p.expr()
	}

	p.expectAndConsume(SEMI, "")

	return VarDecl{
		Name: name,
		Type: typ,
		Init: init,
	}
}

func (p *Parser) constDecl() ConstDecl {
	t := p.expectAndConsume(IDENT, "")
	name := t.Data

	typ := p.varType()

	p.expectAndConsume(OPERATOR, "=")

	// Expression

	var init Expr
	if p.cur().Type != SEMI {
		init = p.expr()
	}

	p.expectAndConsume(SEMI, "")

	return ConstDecl{
		Name: name,
		Type: typ,
		Init: init,
	}
}

func (p *Parser) expr() Expr {
	return p.addExpr()
}

/// Additive
///  : Multiplicative
///  | Additive (+|-) Multiplicative
func (p *Parser) addExpr() Expr {
	l := p.mulExpr()

	for {
		n := p.cur()
		if !(n.Type == OPERATOR && (n.Data == "+" || n.Data == "-")) {
			break
		}

		op := p.consume()

		r := p.mulExpr()
		l = FuncCall{
			Name:         op.Data,
			IsArithmetic: true,
			Args:         []Expr{l, r},
		}
	}

	return l
}

/// Multiplicative
///  : Literal
///  | Multiplicative (*|/) Literal
func (p *Parser) mulExpr() Expr {
	l := p.primaryExpr()

	for {
		n := p.cur()
		if !(n.Type == OPERATOR && (n.Data == "*" || n.Data == "/")) {
			break
		}

		op := p.consume()

		r := p.primaryExpr()
		l = FuncCall{
			Name:         op.Data,
			IsArithmetic: true,
			Args:         []Expr{l, r},
		}
	}

	return l
}

func (p *Parser) primaryExpr() Expr {
	n := p.cur()
	if n.Type != OPERATOR {
		return p.funcExpr()
	}

	p.expectAny([]string{"-"})
	n = p.consume()
	return FuncCall{
		Name:         n.Data,
		IsArithmetic: true,
		Args:         []Expr{p.expr()},
	}
}

func (p *Parser) funcExpr() Expr {
	p.expectAnyType([]TokenType{INT, HEX, FLOAT, STRING, BOOL, NIL, IDENT, LPAREN})
	t := p.consume()

	switch t.Type {
	case LPAREN:
		return p.parenExpr(t)
	case INT:
		fallthrough
	case HEX:
		return p.intLit(t)
	case FLOAT:
		return p.floatLit(t)
	case STRING:
		return p.stringLit(t)
	case BOOL:
		return p.boolLit(t)
	case NIL:
		return p.nilLit(t)
	case IDENT:
		return p.variableOrFuncCall(t)
	}

	panic("WHAT THE SHIT")
}

func (p *Parser) parenExpr(_ Token) ParenExpr {
	e := p.expr()

	p.expectAndConsume(RPAREN, "")

	return ParenExpr{
		Expr: e,
	}
}

func (p *Parser) variableOrFuncCall(t Token) Expr {
	name := t.Data

	next := p.cur()
	if next.Type != LPAREN {
		return Variable{
			Name: name,
		}
	}
	p.consume()

	args := make([]Expr, 0)

	for {
		n := p.cur()
		if n.Type == RPAREN {
			break
		}

		arg := p.expr()
		args = append(args, arg)

		n = p.cur()
		if n.Type == RPAREN {
			break
		}

		p.expectAndConsume(COMMA, "")
	}

	p.expect(RPAREN, "")
	p.consume()

	return FuncCall{
		Name: name,
		Args: args,
	}
}

func (p *Parser) nilLit(_ Token) NilLit {
	return NilLit{}
}

func (p *Parser) boolLit(t Token) BoolLit {
	value := t.Data == "true"

	return BoolLit{
		Value: value,
	}
}

func (p *Parser) stringLit(t Token) StrLit {
	value := t.Data

	return StrLit{
		Value: value,
	}
}

func (p *Parser) floatLit(t Token) FloatLit {
	value, err := strconv.ParseFloat(t.Data, 64)
	if err != nil {
		log.Fatalf("Invalid float literal %s at %d\n", t.Data, p.i-1)
	}

	return FloatLit{
		Value: value,
	}
}

func (p *Parser) intLit(t Token) IntLit {
	var value int
	var err error = nil
	if len(t.Data) > 2 {
		// HEX
		i := new(big.Int)
		i.SetString(t.Data[2:], 16)
		value = int(i.Int64())
	} else {
		value, err = strconv.Atoi(t.Data)
	}

	if err != nil {
		log.Fatalf("Invalid int literal %s at %d\n", t.Data, p.i-1)
	}

	return IntLit{
		Value: value,
	}
}

func (p *Parser) varType() string {
	// Possible type
	pt := p.cur()
	typ := ""
	if pt.Type == OPERATOR && pt.Data == ":" {
		p.consume()

		isRef := false
		// Possible reference type
		prt := p.cur()
		if prt.Type == OPERATOR && prt.Data == "&" {
			isRef = true
			if isRef {
				typ = "&"
			}
			p.consume()
		}

		p.expectAnyType([]TokenType{ITYPE, IDENT})
		t := p.consume()
		typ += t.Data
	}

	return typ
}

func (p *Parser) funcDecl() FuncDecl {
	t := p.expectAndConsume(IDENT, "")

	name := t.Data
	returnType := ""

	p.expectAndConsume(LPAREN, "")

	params := make([]FuncParam, 0)
	for {
		n := p.cur()
		if n.Type == RPAREN {
			break
		}

		param := p.funcParam()
		params = append(params, param)

		n = p.cur()
		if n.Type == RPAREN {
			break
		}

		p.expectAndConsume(COMMA, "")
	}

	p.expectAndConsume(RPAREN, "")

	// Possible return type
	prt := p.cur()
	if prt.Type == OPERATOR {
		p.expect(OPERATOR, "->")

		p.consume()

		p.expectAnyType([]TokenType{ITYPE, IDENT})

		rt := p.consume()
		returnType = rt.Data
	}

	body := p.block()

	return FuncDecl{
		Name:       name,
		ReturnType: returnType,
		Params:     params,
		Body:       body,
	}
}

func (p *Parser) block() Block {
	stmts := make([]Stmt, 0)
	p.expectAndConsume(LCURLY, "")

	for {
		n := p.cur()
		if n.Type == RCURLY {
			break
		}

		stmt := p.stmt()
		stmts = append(stmts, stmt)
	}

	p.expectAndConsume(RCURLY, "")

	return Block{
		Stmts: stmts,
	}
}

func (p *Parser) funcParam() FuncParam {
	n := p.expectAndConsume(IDENT, "")
	name := n.Data
	p.expectAndConsume(OPERATOR, ":")
	p.expectAnyType([]TokenType{ITYPE, IDENT})
	n = p.consume()
	typ := n.Data

	return FuncParam{
		Name: name,
		Type: typ,
	}
}

func (p *Parser) locStmt() LocStmt {
	p.expect(KEYWORD, "loc")
	p.consume()

	sb := strings.Builder{}
	for {
		p.expect(IDENT, "")
		current := p.consume()
		sb.WriteString(current.Data)

		next := p.cur()
		if next.Type != OPERATOR {
			break
		}

		next = p.consume()
		sb.WriteString(next.Data)
	}

	p.expect(SEMI, ";")
	p.consume()

	return LocStmt{
		Value: sb.String(),
	}
}

func (p *Parser) isOfAnyType(typ []TokenType) bool {
	t := p.cur()
	for _, ty := range typ {
		if t.Type == ty {
			return true
		}
	}

	return false
}

func (p *Parser) expectAndConsume(typ TokenType, value string) Token {
	p.expect(typ, value)
	return p.consume()
}

func (p *Parser) expect(typ TokenType, value string) {
	token := p.cur()
	if token.Type != typ {
		log.Fatalf("Expected type %s, but got type %s at %d\n", Name(typ), token.Name(), p.i)
	}

	if len(value) > 0 {
		if value != token.Data {
			log.Fatalf("Expected %s, but got %s at %d\n", value, token.Data, p.i)
		}
	}
}

func (p *Parser) expectAnyType(typs []TokenType) {
	token := p.cur()
	for _, typ := range typs {
		if token.Type == typ {
			return
		}
	}

	typsStr := strings.Join(Map(typs, func(t TokenType) string {
		return Name(t)
	}), ", ")
	log.Fatalf("Expected one of types %s, but got type %s at %d\n", typsStr, token.Name(), p.i)
}

func (p *Parser) expectAny(values []string) {
	token := p.cur()
	for _, value := range values {
		if token.Data == value {
			return
		}
	}

	vStr := strings.Join(values, ", ")
	log.Fatalf("Expected one of %s, but got %s at %d\n", vStr, token.Name(), p.i)
}

func (p *Parser) done() {
	p.expectAndConsume(EOF, "")
}

func (p *Parser) cur() Token {
	return p.tokens[p.i]
}

func (p *Parser) consume() Token {
	tok := p.tokens[p.i]
	p.i++
	return tok
}

func (p *Parser) next() Token {
	return p.tokens[p.i+1]
}
