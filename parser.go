package langkit

import (
	"fmt"
)

func NewParser(lexer Lexer) *TDOPParser {
	return &TDOPParser{
		Lexer: lexer,
	}
}

type TDOPParser struct {
	Lexer Lexer
}

func (parser *TDOPParser) Statement() (*Token, error) {
	tok, err := parser.Lexer.Peek()
	if err != nil {
		return nil, err
	}
	if tok.Std != nil {
		tok, err = parser.Lexer.Next()
		if err != nil {
			return nil, err
		}
		return tok.Std(tok, parser)
	}
	res, err := parser.Expression(0)
	if err != nil {
		return nil, err
	}
	terminator, err := parser.Lexer.Next()
	if !parser.Lexer.IsStatementTerminator(terminator) {

		return nil, fmt.Errorf("syntaxerror: unterminated statement with %v at line %v, col %v", terminator.Value, terminator.Line, terminator.Col)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (parser *TDOPParser) Statements() ([]*Token, error) {
	statements := []*Token{}
	next, err := parser.Lexer.Peek()
	if err != nil {
		return nil, err
	}
	for next.Symbol != EOF {
		statement, err := parser.Statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
		next, err = parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
	}
	return statements, nil
}

func (parser *TDOPParser) Expression(rightBindingPower int) (*Token, error) {
	var left *Token

	t, err := parser.Lexer.Next()
	if err != nil {
		return nil, err
	}
	if t.Nud == nil {
		return nil, fmt.Errorf("%v is not a valid prefix symbol", t.Symbol)
	}
	left, err = t.Nud(t, parser)
	if err != nil {
		return nil, err
	}
	for {
		peek, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		if rightBindingPower >= peek.BindingPower {
			break
		}
		t, err := parser.Lexer.Next()
		if err != nil {
			return nil, err
		}
		if t.Led == nil {
			return nil, fmt.Errorf("%v is not a valid infix symbol", t.Symbol)
		}
		left, err = t.Led(t, parser, left)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}
