package langkit

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// Represents a distinct symbol
type Symbol string

const (
	NameSymbol    Symbol = "(NAME)"
	EOFSymbol     Symbol = "(EOF)"
	StringLiteral Symbol = "(STRING)"
)

type NudFunction func(right *Token, parser *TDOPParser) *Token
type LedFunction func(right *Token, parser *TDOPParser, left *Token) *Token

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
}

type Lexer interface {
	Next() *Token
	Peek() *Token
}

type SymbolTable interface {
	Defined(symbol Symbol) bool
	Register(symbol Symbol, bindingPower int, arity int, nud NudFunction, led LedFunction)
	GenerateToken(symbol Symbol, value string, line int, col int) *Token
	Eof(line int, col int) *Token
	DefineInfix(symbol Symbol, bindingPower int)
	DefineSymbol(symbol Symbol)
	DefinePrefix(symbol Symbol, bindingPower int)
	DefineParens(openParens Symbol, closeParens Symbol)
}

type SymbolTableImpl map[Symbol]*Token

func (table SymbolTableImpl) Defined(symbol Symbol) bool {
	_, present := table[symbol]
	return present
}

func (table SymbolTableImpl) Register(symbol Symbol, bindingPower int, arity int, nud NudFunction, led LedFunction) {
	token := Token{
		Symbol:       symbol,
		Value:        "",
		BindingPower: bindingPower,
		Arity:        arity,
		Nud:          nud,
		Led:          led,
		Children:     []*Token{},
		Col:          -1,
		Line:         -1,
	}
	table[symbol] = &token
}

func (table SymbolTableImpl) GenerateToken(symbol Symbol, value string, line int, col int) *Token {
	tok, present := table[symbol]
	if !present {
		return nil
	}
	return &Token{
		Symbol:       symbol,
		Value:        value,
		Arity:        tok.Arity,
		BindingPower: tok.BindingPower,
		Nud:          tok.Nud,
		Led:          tok.Led,
		Children:     []*Token{},
		Col:          col,
		Line:         line,
	}
}

func (table SymbolTableImpl) Eof(line int, col int) *Token {
	return &Token{
		Symbol:       EOFSymbol,
		Value:        "",
		Arity:        0,
		BindingPower: 0,
		Nud:          nil,
		Led:          nil,
		Children:     []*Token{},
		Col:          col,
		Line:         line,
	}
}

func (table SymbolTableImpl) DefineInfix(symbol Symbol, bindingPower int) {
	led := func(t *Token, parser *TDOPParser, left *Token) *Token {
		t.Children = append(t.Children, left)
		t.Children = append(t.Children, parser.Expression(t.BindingPower))
		return t
	}
	table.Register(symbol, bindingPower, 2, nil, led)
}

func (table SymbolTableImpl) DefinePrefix(symbol Symbol, bindingPower int) {
	nud := func(t *Token, parser *TDOPParser) *Token {
		t.Children = append(t.Children, parser.Expression(bindingPower))
		return t
	}
	table.Register(symbol, bindingPower, 1, nud, nil)
}

// come up with a better name for this
func (table SymbolTableImpl) DefineSymbol(symbol Symbol) {
	nud := func(t *Token, p *TDOPParser) *Token {
		return t
	}
	table.Register(symbol, 0, 0, nud, nil)
}

func (table SymbolTableImpl) DefineParens(openParens Symbol, closeParens Symbol) {
	nud := func(t *Token, p *TDOPParser) *Token {
		expressionToken := p.Expression(0)
		return expressionToken
	}
	table.Register(openParens, 0, 0, nud, nil)
	table.DefineSymbol(closeParens)
}

func NewSymbolTable() SymbolTable {
	var symbolTable SymbolTableImpl = make(map[Symbol]*Token)
	symbolTable.DefineSymbol(StringLiteral)
	symbolTable.DefineSymbol(NameSymbol)
	return symbolTable
}

func NewLexer(reader io.Reader, symbolTable SymbolTable) *TDOPLexer {
	return &TDOPLexer{
		Reader:      bufio.NewReader(reader),
		SymbolTable: symbolTable,
		CachedToken: nil,
		Line:        1,
		Col:         1,
	}
}

type TDOPLexer struct {
	Reader      io.RuneReader
	SymbolTable SymbolTable
	CachedToken *Token
	Line        int
	Col         int
}

func (lexer *TDOPLexer) Peek() *Token {
	token := lexer.Next()
	lexer.CachedToken = token
	return token
}

func (lexer *TDOPLexer) Next() *Token {
	// fmt.Printf("Called Next\n")
	if lexer.CachedToken != nil {
		//fmt.Print("Getting cached token from previous peek\n")
		token := lexer.CachedToken
		lexer.CachedToken = nil
		//fmt.Printf("Returning token %v\n", token)
		return token
	}
	var builder strings.Builder

	char, size, err := lexer.Reader.ReadRune()
	for size > 0 && err == nil {

		// fmt.Printf("Reading rune %v\n", string(char))
		// Handle whitespace
		if unicode.IsSpace(char) {
			// fmt.Println("Reading space character")
			currentSymbol := Symbol(builder.String())
			// fmt.Printf("Current symbol '%v'\n", currentSymbol)
			var token *Token
			if currentSymbol != "" {
				if lexer.SymbolTable.Defined(currentSymbol) {
					// fmt.Printf("Symbol %v found in symbol table\n", currentSymbol)
					token = lexer.SymbolTable.GenerateToken(currentSymbol, string(currentSymbol), lexer.Line, lexer.Col)
				} else {
					// fmt.Printf("Symbol '%v' not found in symbol table\n", currentSymbol)
					token = lexer.SymbolTable.GenerateToken(NameSymbol, string(currentSymbol), lexer.Line, lexer.Col)
				}
				lexer.Col += len(currentSymbol)
			}
			if char == '\n' {
				lexer.Line++
				lexer.Col = 1
			} else {
				lexer.Col++
			}

			if token != nil {
				//fmt.Printf("Returning token %v", token)
				return token
			}

		} else if char == '"' {
			currentSymbol := Symbol(builder.String())
			var stringLiteral strings.Builder
			char, size, err = lexer.Reader.ReadRune()
			for size > 0 && err == nil {
				if char == '"' {

					stringValue := stringLiteral.String()
					stringToken := lexer.SymbolTable.GenerateToken(StringLiteral, stringValue, lexer.Line, lexer.Col)
					if lexer.SymbolTable.Defined(currentSymbol) {
						lexer.CachedToken = stringToken
						token := lexer.SymbolTable.GenerateToken(currentSymbol, string(currentSymbol), lexer.Line, lexer.Col)
						lexer.Col += len(currentSymbol)
						// fmt.Printf("Returning token %v\n", token)
						return token
					} else {
						return stringToken
					}
				} else {
					stringLiteral.WriteRune(char)
				}
				char, size, err = lexer.Reader.ReadRune()
			}

			if lexer.SymbolTable.Defined(currentSymbol) {
				token := lexer.SymbolTable.GenerateToken(currentSymbol, string(currentSymbol), lexer.Line, lexer.Col)
				lexer.Col += len(currentSymbol)
				// fmt.Printf("Returning token %v\n", token)
				return token
			}
		} else {
			// fmt.Println("Reading text character")
			asSymbol := Symbol(string(char))
			if lexer.SymbolTable.Defined(asSymbol) {
				currentSymbol := Symbol(builder.String())
				if currentSymbol != "" {
					cachedToken := lexer.SymbolTable.GenerateToken(asSymbol, string(asSymbol), lexer.Line, lexer.Col)
					// fmt.Printf("Setting cached Token %v", cachedToken)
					lexer.CachedToken = cachedToken
					lexer.Col += 1
					if lexer.SymbolTable.Defined(currentSymbol) {
						token := lexer.SymbolTable.GenerateToken(currentSymbol, string(currentSymbol), lexer.Line, lexer.Col)
						// fmt.Printf("Symbol %v found in symbol table\n", currentSymbol)
						lexer.Col += len(currentSymbol)
						// fmt.Printf("Returning token %v\n", token)
						return token
					} else {
						// fmt.Printf("Symbol '%v' not found in symbol table\n", currentSymbol)
						token := lexer.SymbolTable.GenerateToken(NameSymbol, string(currentSymbol), lexer.Line, lexer.Col)
						lexer.Col += len(currentSymbol)
						// fmt.Printf("Returning token %v\n", token)
						return token
					}
				} else {
					token := lexer.SymbolTable.GenerateToken(asSymbol, string(asSymbol), lexer.Line, lexer.Col)
					lexer.Col += 1
					return token
				}
			} else {
				builder.WriteRune(char)
				currentSymbol := Symbol(builder.String())
				if lexer.SymbolTable.Defined(currentSymbol) {
					token := lexer.SymbolTable.GenerateToken(currentSymbol, string(currentSymbol), lexer.Line, lexer.Col)
					lexer.Col += len(currentSymbol)
					// fmt.Printf("Returning token %v\n", token)
					return token
				}
			}
		}
		char, size, err = lexer.Reader.ReadRune()
	}
	if errors.Is(err, io.EOF) {
		currentSymbol := Symbol(builder.String())
		if currentSymbol != "" {
			if lexer.SymbolTable.Defined(currentSymbol) {
				token := lexer.SymbolTable.GenerateToken(currentSymbol, string(currentSymbol), lexer.Line, lexer.Col)
				//fmt.Printf("Symbol %v found in symbol table\n", currentSymbol)
				lexer.Col += len(currentSymbol)
				fmt.Printf("Returning token %v\n", token)
				return token
			} else {
				//fmt.Printf("Symbol '%v' not found in symbol table\n", currentSymbol)
				token := lexer.SymbolTable.GenerateToken(NameSymbol, string(currentSymbol), lexer.Line, lexer.Col)
				lexer.Col += len(currentSymbol)
				//fmt.Printf("Returning token %v\n", token)
				return token
			}
		}
		//fmt.Printf("Reached (EOF)")
		return lexer.SymbolTable.Eof(lexer.Line, lexer.Col)
	}

	panic(fmt.Sprintf("Invalid character at line: %v col: %v", lexer.Line, lexer.Col))
}

type TDOPParser struct {
	Lexer Lexer
}

func (parser *TDOPParser) Expression(rightBindingPower int) *Token {
	var left *Token

	t := parser.Lexer.Next()
	//fmt.Printf("Calling Nud on Token %v\n", t)
	left = t.Nud(t, parser)
	for rightBindingPower < parser.Lexer.Peek().BindingPower {
		t := parser.Lexer.Next()
		left = t.Led(t, parser, left)
	}
	return left
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
