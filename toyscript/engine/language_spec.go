package engine

import (
	"fmt"

	"github.com/nicholasbailey/langkit"
)

func BuildToyscriptLanguageSpec() langkit.LanguageSpecification {
	spec := langkit.NewLanguage()
	spec.DefineQuotes('"', '"', langkit.StringLiteral)
	spec.DefineQuotes('\'', '\'', langkit.StringLiteral)
	spec.DefineParens("(", ")")
	spec.DefinePrefix("!", 80)
	spec.DefineInfix("&&", 30)
	spec.DefineInfix("||", 20)
	spec.DefineInfix("=", 10)
	spec.DefineInfix("==", 50)
	spec.DefineInfix("!=", 50)
	spec.DefineInfix("<", 50)
	spec.DefineInfix(">", 50)
	spec.DefineInfix("<=", 50)
	spec.DefineInfix(">=", 50)
	spec.DefineInfix("+", 60)
	spec.DefineInfix("-", 60)
	spec.DefineInfix("*", 70)
	spec.DefineInfix("/", 70)
	spec.DefineInfix("%", 70)
	spec.DefineStatementTerminator(";")
	spec.DefineEmpty(",")
	spec.DefineBlock("{", "}")
	spec.DefineEmpty("else")
	ifStd := func(token *langkit.Token, parser *langkit.TDOPParser) (*langkit.Token, error) {
		expression, err := parser.Expression(0)
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, expression)
		peek, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		// TODO - remove this Hack because expression parser is wonky
		if peek.Symbol == ")" {
			_, err = parser.Lexer.Next()
			if err != nil {
				return nil, err
			}
		}
		block, err := parser.Block()
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, block)
		next, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		// TODO: make this symbol references
		if next.Value == "else" {
			parser.Lexer.Next()
			next, err := parser.Lexer.Peek()
			if err != nil {
				return nil, err
			}
			if next.Value == "if" {
				stmt, err := parser.Statement()
				if err != nil {
					return nil, err
				}
				stmt.Symbol = langkit.ElseIf
				token.Children = append(token.Children, stmt)
			} else {
				block, err := parser.Block()
				if err != nil {
					return nil, err
				}
				token.Children = append(token.Children, block)
			}
		}
		return token, nil
	}

	ifNud := func(token *langkit.Token, parser *langkit.TDOPParser) (*langkit.Token, error) {
		expression, err := parser.Expression(0)
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, expression)
		peek, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		// Hack because expression parser is wonky
		if peek.Symbol == ")" {
			_, err = parser.Lexer.Next()
			if err != nil {
				return nil, err
			}
		}
		block, err := parser.Block()
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, block)
		next, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		// TODO: make this symbol references
		if next.Value == "else" {
			parser.Lexer.Next()
			next, err := parser.Lexer.Peek()
			if err != nil {
				return nil, err
			}
			if next.Value == "if" {
				stmt, err := parser.Statement()
				if err != nil {
					return nil, err
				}
				token.Children = append(token.Children, stmt)
			} else {
				block, err := parser.Block()
				if err != nil {
					return nil, err
				}
				token.Children = append(token.Children, block)
			}
		}
		return token, nil
	}

	spec.Define("if", 0, 0, ifNud, nil, ifStd)

	whileStd := func(token *langkit.Token, parser *langkit.TDOPParser) (*langkit.Token, error) {
		expression, err := parser.Expression(0)
		if err != nil {
			return nil, err
		}
		peek, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		// TODO - remove this Hack because expression parser is wonky
		if peek.Symbol == ")" {
			_, err = parser.Lexer.Next()
			if err != nil {
				return nil, err
			}
		}
		token.Children = append(token.Children, expression)
		block, err := parser.Block()
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, block)
		return token, nil
	}

	spec.DefineStatment("while", whileStd)

	openParensLed := func(right *langkit.Token, parser *langkit.TDOPParser, left *langkit.Token) (*langkit.Token, error) {
		if left.Symbol != langkit.Name && left.Symbol != langkit.Symbol("(") {
			return nil, fmt.Errorf("syntaxerror: unexpected ( at line %v, col %v", right.Line, right.Col)
		}
		right.Children = append(right.Children, left)
		t, err := parser.Lexer.Peek()
		if err != nil {
			return nil, err
		}
		if t.Symbol != ")" {
			for {
				expressionResult, err := parser.Expression(0)
				if err != nil {
					return nil, err
				}
				right.Children = append(right.Children, expressionResult)
				right, err := parser.Lexer.Peek()
				if err != nil {
					return nil, err
				}
				if right.Symbol != "," {
					break
				}
				_, err = parser.Lexer.Next()
				if err != nil {
					return nil, err
				}
			}
			close, err := parser.Lexer.Next()
			if err != nil {
				return nil, err
			}
			if close.Symbol != ")" {
				return nil, fmt.Errorf("syntaxerror: unterminated parentheses with symbol %v at line %v, col %v", close.Value, close.Line, close.Col)
			}
		} else {
			_, err = parser.Lexer.Next()
			if err != nil {
				return nil, err
			}
		}
		right.Symbol = langkit.FunctionInvocation
		return right, nil
	}

	spec.Define("(", 90, 0, nil, openParensLed, nil)

	defStd := func(token *langkit.Token, parser *langkit.TDOPParser) (*langkit.Token, error) {
		fmt.Printf("Parsing function definition\n")
		token.Symbol = langkit.FunctionDefinition
		functionName, err := parser.Lexer.Next()
		if err != nil {
			return nil, err
		}
		if functionName.Symbol != langkit.Name {
			return nil, Exception(SyntaxError, fmt.Sprintf("expected identifier, got %v", functionName.Value), functionName)
		}
		token.Children = append(token.Children, functionName)

		openParens, err := parser.Lexer.Next()
		if err != nil {
			return nil, err
		}
		if openParens.Symbol != "(" {
			return nil, Exception(SyntaxError, fmt.Sprintf("expected (, got %v", openParens.Value), openParens)
		}
		parameters := []*langkit.Token{}
		next, err := parser.Lexer.Next()
		if err != nil {
			return nil, err
		}
		if next.Symbol != ")" {
			for {
				if next.Symbol != langkit.Name {
					return nil, err
				}
				parameters = append(parameters, next)
				further, err := parser.Lexer.Peek()
				if err != nil {
					return nil, err
				}
				if further.Symbol != "," {
					break
				}
				_, err = parser.Lexer.Next()
				if err != nil {
					return nil, err
				}
			}
			close, err := parser.Lexer.Next()
			if err != nil {
				return nil, err
			}
			if close.Symbol != ")" {
				return nil, fmt.Errorf("syntaxerror: unterminated parentheses with symbol %v at line %v, col %v", close.Value, close.Line, close.Col)
			}
		} else {
			_, err = parser.Lexer.Next()
			if err != nil {
				return nil, err
			}
		}
		parameterToken := &langkit.Token{
			Symbol:   langkit.FunctionParameters,
			Value:    "(",
			Line:     openParens.Line,
			Arity:    len(parameters),
			Col:      openParens.Col,
			Children: parameters,
		}
		token.Children = append(token.Children, parameterToken)
		block, err := parser.Block()
		if err != nil {
			return nil, err
		}
		token.Children = append(token.Children, block)
		return token, nil
	}

	spec.DefineStatment("def", defStd)

	return spec
}
