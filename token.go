package langkit

import (
	"fmt"
	"strings"
)

// Represents a distinct symbol
type Symbol string

const (
	Name               Symbol = "(NAME)"
	EOF                Symbol = "(EOF)"
	StringLiteral      Symbol = "(STRING)"
	IntLiteral         Symbol = "(INT)"
	FloatLiteral       Symbol = "(FLOAT)"
	FunctionInvocation Symbol = "(FUNCTIONINVOCATION)"
	Block              Symbol = "(BLOCK)"
	ElseIf             Symbol = "(ELSEIF)"
	FunctionDefinition Symbol = "(FUNCTIONDEFINITION)"
	FunctionParameters Symbol = "(FUNCTIONPARAMETERS)"
)

type NudFunction func(right *Token, parser *TDOPParser) (*Token, error)
type LedFunction func(right *Token, parser *TDOPParser, left *Token) (*Token, error)
type StdFunction func(*Token, *TDOPParser) (*Token, error)

type Token struct {
	Symbol       Symbol
	Value        string
	Arity        int
	BindingPower int
	Line         int
	Col          int
	Children     []*Token
	Nud          NudFunction
	Led          LedFunction
	Std          StdFunction
}

func (token *Token) TreeString(indentLevel int) string {
	var builder strings.Builder
	for i := 0; i < indentLevel; i++ {
		builder.WriteString("  ")
	}

	builder.WriteString(fmt.Sprintf("{symbol:%v,value:%v,bindingPower:%v,arity:%v}:\n", token.Symbol, token.Value, token.BindingPower, token.Arity))
	for _, child := range token.Children {
		builder.WriteString(child.TreeString(indentLevel + 1))
	}
	return builder.String()
}
