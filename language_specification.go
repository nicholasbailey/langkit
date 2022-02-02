package langkit

import (
	"unicode"
)

// Represents a collection of symbols
// defined for the language along with their
// binding behavior
type LanguageSpecification interface {
	// Checks if a symbol is defined in the language
	IsDefined(symbol Symbol) bool
	// Defines a new symbol
	Define(symbol Symbol, bindingPower int, arity int, nud NudFunction, led LedFunction, std StdFunction)

	DefineInfix(symbol Symbol, bindingPower int)
	DefineValue(symbol Symbol)
	DefinePrefix(symbol Symbol, bindingPower int)
	DefineParens(openParens Symbol, closeParens Symbol)
	DefineQuotes(openQuote rune, closeQuote rune, literalType Symbol)
	// Generates a token for the given symbol and a given value
	GenerateToken(symbol Symbol, value string, line int, col int) *Token
	//
	Eof(line int, col int) *Token
	GetQuoteSpec(openQuote rune) *quoteSpecification
	IsIdentifierCharacter(character rune) bool
	IsIdentifierStartChararacter(character rune) bool
	DefineStatementTerminator(symbol Symbol)
	IsStatementTerminator(symbol Symbol) bool
	DefineEmpty(symbol Symbol)
}

func NewLanguage() LanguageSpecification {
	symbols := make(map[Symbol]*Token)
	quotes := make(map[rune]*quoteSpecification)
	language := &languageSpecificationImpl{
		quoteDefinitions: quotes,
		symbols:          symbols,
	}
	language.DefineValue(Name)
	language.DefineValue(IntLiteral)
	language.DefineValue(FloatLiteral)
	return language
}

type quoteSpecification struct {
	openQuote   rune
	closeQuote  rune
	literalType Symbol
}

type languageSpecificationImpl struct {
	quoteDefinitions     map[rune]*quoteSpecification
	symbols              map[Symbol]*Token
	statementTerminators []Symbol
}

func (spec *languageSpecificationImpl) DefineEmpty(symbol Symbol) {
	spec.Define(symbol, 0, 0, nil, nil, nil)
}

func (spec *languageSpecificationImpl) IsStatementTerminator(symbol Symbol) bool {
	for _, s := range spec.statementTerminators {
		if s == symbol {
			return true
		}
	}
	return false
}

func (spec *languageSpecificationImpl) DefineStatementTerminator(symbol Symbol) {
	spec.statementTerminators = append(spec.statementTerminators, symbol)
	spec.symbols[symbol] = &Token{
		Symbol:       symbol,
		Arity:        0,
		BindingPower: 0,
		Nud:          nil,
		Led:          nil,
		Std:          nil,
		Children:     []*Token{},
	}
}

func (spec *languageSpecificationImpl) IsIdentifierStartChararacter(char rune) bool {
	return spec.IsIdentifierCharacter(char) && !unicode.IsDigit(char)
}

func (spec *languageSpecificationImpl) IsIdentifierCharacter(char rune) bool {
	_, found := spec.symbols[Symbol(string(char))]
	if found {
		return false
	} else {
		return !unicode.IsSpace(char)
	}
}

func (spec *languageSpecificationImpl) GetQuoteSpec(openQuote rune) *quoteSpecification {
	val, _ := spec.quoteDefinitions[openQuote]
	return val
}

func (spec *languageSpecificationImpl) IsDefined(symbol Symbol) bool {
	_, present := spec.symbols[symbol]
	return present
}

func (spec *languageSpecificationImpl) Define(symbol Symbol, bindingPower int, arity int, nud NudFunction, led LedFunction, std StdFunction) {
	existing, found := spec.symbols[symbol]
	if found {
		if nud != nil && existing.Nud == nil {
			existing.Nud = nud
		}
		if led != nil && existing.Led == nil {
			existing.Led = led
		}
		if std != nil && existing.Std == nil {
			existing.Std = std
		}
		if bindingPower > existing.BindingPower {
			existing.BindingPower = bindingPower
		}
	} else {
		token := Token{
			Symbol:       symbol,
			Value:        "",
			BindingPower: bindingPower,
			Arity:        arity,
			Nud:          nud,
			Led:          led,
			Std:          std,
			Children:     []*Token{},
			Col:          -1,
			Line:         -1,
		}
		spec.symbols[symbol] = &token
	}
}

func (spec *languageSpecificationImpl) GenerateToken(symbol Symbol, value string, line int, col int) *Token {
	tok, present := spec.symbols[symbol]
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

func (spec *languageSpecificationImpl) Eof(line int, col int) *Token {
	return &Token{
		Symbol:       EOF,
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

func (spec *languageSpecificationImpl) DefineQuotes(openQuote rune, closeQuote rune, literalType Symbol) {
	spec.quoteDefinitions[openQuote] = &quoteSpecification{
		openQuote:   openQuote,
		closeQuote:  closeQuote,
		literalType: literalType,
	}
	spec.DefineValue(literalType)
}

func (spec *languageSpecificationImpl) DefineInfix(symbol Symbol, bindingPower int) {
	led := func(t *Token, parser *TDOPParser, left *Token) (*Token, error) {
		t.Children = append(t.Children, left)
		exprResult, err := parser.Expression(t.BindingPower)
		if err != nil {
			return nil, err
		}
		t.Children = append(t.Children, exprResult)
		return t, nil
	}
	spec.Define(symbol, bindingPower, 2, nil, led, nil)
}

func (spec *languageSpecificationImpl) DefinePrefix(symbol Symbol, bindingPower int) {
	nud := func(t *Token, parser *TDOPParser) (*Token, error) {
		expResult, err := parser.Expression(bindingPower)
		if err != nil {
			return nil, err
		}
		t.Children = append(t.Children, expResult)
		return t, nil
	}
	spec.Define(symbol, bindingPower, 1, nud, nil, nil)
}

// come up with a better name for this
func (spec *languageSpecificationImpl) DefineValue(symbol Symbol) {
	nud := func(t *Token, p *TDOPParser) (*Token, error) {
		return t, nil
	}
	spec.Define(symbol, 0, 0, nud, nil, nil)
}

func (spec *languageSpecificationImpl) DefineParens(openParens Symbol, closeParens Symbol) {
	nud := func(t *Token, p *TDOPParser) (*Token, error) {
		expressionToken, err := p.Expression(0)
		if err != nil {
			return nil, err
		}
		return expressionToken, nil
	}
	spec.Define(openParens, 0, 0, nud, nil, nil)
	spec.DefineValue(closeParens)
}
