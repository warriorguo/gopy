package parser

import (
	"fmt"
	"strconv"

	"github.com/warriorguo/gopy/pkg/ast"
	"github.com/warriorguo/gopy/pkg/lexer"
)

type Parser struct {
	tokens   []lexer.Token
	position int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
	}
}

func (p *Parser) currentToken() lexer.Token {
	if p.position >= len(p.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}
	return p.tokens[p.position]
}

func (p *Parser) advance() {
	if p.position < len(p.tokens) {
		p.position++
	}
}

func (p *Parser) expect(tokenType lexer.TokenType) error {
	if p.currentToken().Type != tokenType {
		return fmt.Errorf("expected %s, got %s(%s) at line %d", tokenType, p.currentToken().Type, p.currentToken().Lexeme, p.currentToken().Line)
	}
	p.advance()
	return nil
}

func (p *Parser) skipNewlines() {
	for p.currentToken().Type == lexer.NEWLINE {
		p.advance()
	}
}

func (p *Parser) Parse() (*ast.Module, error) {
	module := &ast.Module{
		Body:     []ast.Stmt{},
		Position: ast.Position{Line: 1, Column: 1},
	}

	p.skipNewlines()

	for p.currentToken().Type != lexer.EOF {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			module.Body = append(module.Body, stmt)
		}
		p.skipNewlines()
	}

	return module, nil
}

func (p *Parser) parseStmt() (ast.Stmt, error) {
	switch p.currentToken().Type {
	case lexer.IF:
		return p.parseIfStmt()
	case lexer.WHILE:
		return p.parseWhileStmt()
	case lexer.FOR:
		return p.parseForStmt()
	case lexer.DEF:
		return p.parseFuncDef()
	case lexer.RETURN:
		return p.parseReturnStmt()
	case lexer.PASS:
		return p.parsePassStmt()
	case lexer.PRINT:
		return p.parsePrintStmt()
	case lexer.NEWLINE:
		p.advance()
		return nil, nil
	default:
		if p.isAssignment() {
			return p.parseAssignStmt()
		}
		return p.parseExprStmt()
	}
}

func (p *Parser) isAssignment() bool {
	saved := p.position
	defer func() { p.position = saved }()

	if p.currentToken().Type == lexer.IDENT {
		p.advance()
		tokenType := p.currentToken().Type
		return tokenType == lexer.ASSIGN || tokenType == lexer.PLUS_ASSIGN || tokenType == lexer.MINUS_ASSIGN
	}
	return false
}

func (p *Parser) parseAssignStmt() (ast.Stmt, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}

	target, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	tokenType := p.currentToken().Type

	switch tokenType {
	case lexer.ASSIGN:
		p.advance()
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.AssignStmt{
			Target:   target,
			Value:    value,
			Position: pos,
		}, nil

	case lexer.PLUS_ASSIGN:
		p.advance()
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.AugAssignStmt{
			Target:   target,
			Op:       "+=",
			Value:    value,
			Position: pos,
		}, nil

	case lexer.MINUS_ASSIGN:
		p.advance()
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.AugAssignStmt{
			Target:   target,
			Op:       "-=",
			Value:    value,
			Position: pos,
		}, nil

	default:
		return nil, fmt.Errorf("expected assignment operator, got %s at line %d", tokenType, p.currentToken().Line)
	}
}

func (p *Parser) parseExprStmt() (ast.Stmt, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.ExprStmt{
		Expr:     expr,
		Position: pos,
	}, nil
}

func (p *Parser) parsePrintStmt() (ast.Stmt, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	p.advance()

	var values []ast.Expr

	if p.currentToken().Type != lexer.NEWLINE && p.currentToken().Type != lexer.EOF {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		values = append(values, expr)

		for p.currentToken().Type == lexer.COMMA {
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			values = append(values, expr)
		}
	}

	return &ast.PrintStmt{
		Values:   values,
		Position: pos,
	}, nil
}

func (p *Parser) parseIfStmt() (ast.Stmt, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	p.advance()

	test, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	var orelse []ast.Stmt

	if p.currentToken().Type == lexer.ELIF || p.currentToken().Type == lexer.ELSE {
		if p.currentToken().Type == lexer.ELSE {
			p.advance()
			if err := p.expect(lexer.COLON); err != nil {
				return nil, err
			}
			orelse, err = p.parseBlock()
			if err != nil {
				return nil, err
			}
		} else {
			elifStmt, err := p.parseIfStmt()
			if err != nil {
				return nil, err
			}
			orelse = []ast.Stmt{elifStmt}
		}
	}

	return &ast.IfStmt{
		Test:     test,
		Body:     body,
		Orelse:   orelse,
		Position: pos,
	}, nil
}

func (p *Parser) parseWhileStmt() (ast.Stmt, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	p.advance()

	test, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStmt{
		Test:     test,
		Body:     body,
		Position: pos,
	}, nil
}

func (p *Parser) parseForStmt() (ast.Stmt, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	p.advance()

	target, err := p.parseAtomExpr()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.IN); err != nil {
		return nil, err
	}

	iter, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ast.ForStmt{
		Target:   target,
		Iter:     iter,
		Body:     body,
		Position: pos,
	}, nil
}

func (p *Parser) parseFuncDef() (ast.Stmt, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	p.advance()

	if p.currentToken().Type != lexer.IDENT {
		return nil, fmt.Errorf("expected function name at line %d", p.currentToken().Line)
	}
	name := p.currentToken().Lexeme
	p.advance()

	if err := p.expect(lexer.LPAREN); err != nil {
		return nil, err
	}

	var args []string
	if p.currentToken().Type == lexer.IDENT {
		args = append(args, p.currentToken().Lexeme)
		p.advance()

		for p.currentToken().Type == lexer.COMMA {
			p.advance()
			if p.currentToken().Type != lexer.IDENT {
				return nil, fmt.Errorf("expected parameter name at line %d", p.currentToken().Line)
			}
			args = append(args, p.currentToken().Lexeme)
			p.advance()
		}
	}

	if err := p.expect(lexer.RPAREN); err != nil {
		return nil, err
	}

	if err := p.expect(lexer.COLON); err != nil {
		return nil, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ast.FuncDef{
		Name:     name,
		Args:     args,
		Body:     body,
		Position: pos,
	}, nil
}

func (p *Parser) parseReturnStmt() (ast.Stmt, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	p.advance()

	var value ast.Expr
	if p.currentToken().Type != lexer.NEWLINE && p.currentToken().Type != lexer.EOF {
		var err error
		value, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	return &ast.ReturnStmt{
		Value:    value,
		Position: pos,
	}, nil
}

func (p *Parser) parsePassStmt() (ast.Stmt, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	p.advance() // consume 'pass' token

	return &ast.PassStmt{
		Position: pos,
	}, nil
}

func (p *Parser) parseBlock() ([]ast.Stmt, error) {
	p.skipNewlines()

	if err := p.expect(lexer.INDENT); err != nil {
		return nil, err
	}

	var stmts []ast.Stmt
	for p.currentToken().Type != lexer.DEDENT && p.currentToken().Type != lexer.EOF {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
		p.skipNewlines()
	}

	if p.currentToken().Type == lexer.DEDENT {
		p.advance()
	} else if p.currentToken().Type != lexer.EOF {
		return nil, fmt.Errorf("expected DEDENT or EOF, got %s at line %d",
			p.currentToken().Type, p.currentToken().Line)
	}

	return stmts, nil
}

func (p *Parser) parseExpr() (ast.Expr, error) {
	return p.parseOrExpr()
}

func (p *Parser) parseOrExpr() (ast.Expr, error) {
	left, err := p.parseAndExpr()
	if err != nil {
		return nil, err
	}

	if p.currentToken().Type == lexer.OR {
		pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
		values := []ast.Expr{left}

		for p.currentToken().Type == lexer.OR {
			p.advance()
			right, err := p.parseAndExpr()
			if err != nil {
				return nil, err
			}
			values = append(values, right)
		}

		return &ast.BoolOp{
			Op:       "or",
			Values:   values,
			Position: pos,
		}, nil
	}

	return left, nil
}

func (p *Parser) parseAndExpr() (ast.Expr, error) {
	left, err := p.parseNotExpr()
	if err != nil {
		return nil, err
	}

	if p.currentToken().Type == lexer.AND {
		pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
		values := []ast.Expr{left}

		for p.currentToken().Type == lexer.AND {
			p.advance()
			right, err := p.parseNotExpr()
			if err != nil {
				return nil, err
			}
			values = append(values, right)
		}

		return &ast.BoolOp{
			Op:       "and",
			Values:   values,
			Position: pos,
		}, nil
	}

	return left, nil
}

func (p *Parser) parseNotExpr() (ast.Expr, error) {
	if p.currentToken().Type == lexer.NOT {
		pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
		p.advance()
		expr, err := p.parseNotExpr()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryOp{
			Op:       "not",
			Expr:     expr,
			Position: pos,
		}, nil
	}
	return p.parseCompareExpr()
}

func (p *Parser) parseCompareExpr() (ast.Expr, error) {
	left, err := p.parseArithExpr()
	if err != nil {
		return nil, err
	}

	if p.isCompOp() {
		pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
		var ops []string
		var rights []ast.Expr

		for p.isCompOp() {
			op := p.getCompOp()
			ops = append(ops, op)
			p.advance()
			right, err := p.parseArithExpr()
			if err != nil {
				return nil, err
			}
			rights = append(rights, right)
		}

		return &ast.Compare{
			Left:     left,
			Ops:      ops,
			Right:    rights,
			Position: pos,
		}, nil
	}

	return left, nil
}

func (p *Parser) isCompOp() bool {
	switch p.currentToken().Type {
	case lexer.EQ, lexer.NOT_EQ, lexer.LT, lexer.GT, lexer.LTE, lexer.GTE, lexer.IN:
		return true
	}
	return false
}

func (p *Parser) getCompOp() string {
	switch p.currentToken().Type {
	case lexer.EQ:
		return "=="
	case lexer.NOT_EQ:
		return "!="
	case lexer.LT:
		return "<"
	case lexer.GT:
		return ">"
	case lexer.LTE:
		return "<="
	case lexer.GTE:
		return ">="
	case lexer.IN:
		return "in"
	}
	return ""
}

func (p *Parser) parseArithExpr() (ast.Expr, error) {
	left, err := p.parseTermExpr()
	if err != nil {
		return nil, err
	}

	for p.currentToken().Type == lexer.PLUS || p.currentToken().Type == lexer.MINUS {
		pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
		op := p.currentToken().Lexeme
		p.advance()
		right, err := p.parseTermExpr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryOp{
			Left:     left,
			Op:       op,
			Right:    right,
			Position: pos,
		}
	}

	return left, nil
}

func (p *Parser) parseTermExpr() (ast.Expr, error) {
	left, err := p.parseFactorExpr()
	if err != nil {
		return nil, err
	}

	for p.currentToken().Type == lexer.MULTIPLY || p.currentToken().Type == lexer.DIVIDE || p.currentToken().Type == lexer.MODULO {
		pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
		op := p.currentToken().Lexeme
		p.advance()
		right, err := p.parseFactorExpr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryOp{
			Left:     left,
			Op:       op,
			Right:    right,
			Position: pos,
		}
	}

	return left, nil
}

func (p *Parser) parseFactorExpr() (ast.Expr, error) {
	if p.currentToken().Type == lexer.MINUS || p.currentToken().Type == lexer.PLUS {
		pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
		op := p.currentToken().Lexeme
		p.advance()
		expr, err := p.parseFactorExpr()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryOp{
			Op:       op,
			Expr:     expr,
			Position: pos,
		}, nil
	}

	return p.parsePowerExpr()
}

func (p *Parser) parsePowerExpr() (ast.Expr, error) {
	return p.parseCallExpr()
}

func (p *Parser) parseCallExpr() (ast.Expr, error) {
	expr, err := p.parseAtomExpr()
	if err != nil {
		return nil, err
	}

	for {
		switch p.currentToken().Type {
		case lexer.LPAREN:
			pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
			p.advance()
			var args []ast.Expr

			if p.currentToken().Type != lexer.RPAREN {
				arg, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)

				for p.currentToken().Type == lexer.COMMA {
					p.advance()
					arg, err := p.parseExpr()
					if err != nil {
						return nil, err
					}
					args = append(args, arg)
				}
			}

			if err := p.expect(lexer.RPAREN); err != nil {
				return nil, err
			}

			expr = &ast.Call{
				Func:     expr,
				Args:     args,
				Position: pos,
			}

		case lexer.LBRACKET:
			pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
			p.advance()
			slice, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.RBRACKET); err != nil {
				return nil, err
			}

			expr = &ast.Subscript{
				Value:    expr,
				Slice:    slice,
				Position: pos,
			}

		default:
			return expr, nil
		}
	}
}

func (p *Parser) parseAtomExpr() (ast.Expr, error) {
	switch p.currentToken().Type {
	case lexer.INT:
		return p.parseNumber()
	case lexer.FLOAT:
		return p.parseNumber()
	case lexer.STRING:
		return p.parseString()
	case lexer.TRUE, lexer.FALSE, lexer.NONE:
		return p.parseNameConstant()
	case lexer.IDENT, lexer.RANGE:
		return p.parseName()
	case lexer.LPAREN:
		return p.parseParenExpr()
	case lexer.LBRACKET:
		return p.parseList()
	case lexer.LBRACE:
		return p.parseDict()
	default:
		return nil, fmt.Errorf("unexpected token %s at line %d", p.currentToken().Type, p.currentToken().Line)
	}
}

func (p *Parser) parseNumber() (ast.Expr, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	lexeme := p.currentToken().Lexeme
	tokenType := p.currentToken().Type
	p.advance()

	var value interface{}
	var err error

	if tokenType == lexer.INT {
		value, err = strconv.Atoi(lexeme)
	} else {
		value, err = strconv.ParseFloat(lexeme, 64)
	}

	if err != nil {
		return nil, fmt.Errorf("invalid number %s at line %d", lexeme, pos.Line)
	}

	return &ast.Num{
		N:        value,
		Position: pos,
	}, nil
}

func (p *Parser) parseString() (ast.Expr, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	value := p.currentToken().Lexeme
	p.advance()

	return &ast.Str{
		S:        value,
		Position: pos,
	}, nil
}

func (p *Parser) parseNameConstant() (ast.Expr, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}

	var value interface{}
	switch p.currentToken().Type {
	case lexer.TRUE:
		value = true
	case lexer.FALSE:
		value = false
	case lexer.NONE:
		value = nil
	}
	p.advance()

	return &ast.NameConstant{
		Value:    value,
		Position: pos,
	}, nil
}

func (p *Parser) parseName() (ast.Expr, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	name := p.currentToken().Lexeme
	p.advance()

	return &ast.Name{
		Id:       name,
		Position: pos,
	}, nil
}

func (p *Parser) parseParenExpr() (ast.Expr, error) {
	p.advance()
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if err := p.expect(lexer.RPAREN); err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) parseList() (ast.Expr, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	p.advance()

	var elts []ast.Expr
	if p.currentToken().Type != lexer.RBRACKET {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		elts = append(elts, expr)

		for p.currentToken().Type == lexer.COMMA {
			p.advance()
			if p.currentToken().Type == lexer.RBRACKET {
				break
			}
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			elts = append(elts, expr)
		}
	}

	if err := p.expect(lexer.RBRACKET); err != nil {
		return nil, err
	}

	return &ast.List{
		Elts:     elts,
		Position: pos,
	}, nil
}

func (p *Parser) parseDict() (ast.Expr, error) {
	pos := ast.Position{Line: p.currentToken().Line, Column: p.currentToken().Column}
	p.advance()

	var keys []ast.Expr
	var values []ast.Expr

	if p.currentToken().Type != lexer.RBRACE {
		key, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expect(lexer.COLON); err != nil {
			return nil, err
		}
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
		values = append(values, value)

		for p.currentToken().Type == lexer.COMMA {
			p.advance()
			if p.currentToken().Type == lexer.RBRACE {
				break
			}
			key, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if err := p.expect(lexer.COLON); err != nil {
				return nil, err
			}
			value, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
			values = append(values, value)
		}
	}

	if err := p.expect(lexer.RBRACE); err != nil {
		return nil, err
	}

	return &ast.Dict{
		Keys:     keys,
		Values:   values,
		Position: pos,
	}, nil
}

func Parse(tokens []lexer.Token) (*ast.Module, error) {
	parser := NewParser(tokens)
	return parser.Parse()
}
