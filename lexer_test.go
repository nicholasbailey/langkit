package langkit

import (
	"fmt"
	"strings"
	"testing"
)

func makeLexer(sourceCode string) *TDOPLexer {
	symbolTable := NewLanguage()
	symbolTable.DefineInfix("AND", 20)
	symbolTable.DefineInfix("OR", 10)
	symbolTable.DefinePrefix("NOT", 40)
	symbolTable.DefineInfix("=", 30)
	symbolTable.DefineInfix("!=", 30)
	symbolTable.DefineParens("(", ")")
	symbolTable.DefineQuotes('"', '"', StringLiteral)
	return NewLexer(strings.NewReader(sourceCode), symbolTable)
}

func TestLexer(t *testing.T) {
	type testCase struct {
		input    string
		expected []*Token
	}

	testCases := []testCase{
		testCase{
			input: "A",
			expected: []*Token{
				&Token{
					Symbol:       "(NAME)",
					Value:        "A",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          1,
				},
				&Token{
					Symbol:       "(EOF)",
					Value:        "",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          15,
				},
			},
		},
		testCase{
			input: "A = \"Hello\"",
			expected: []*Token{
				&Token{
					Symbol:       "(NAME)",
					Value:        "A",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          1,
				},
				&Token{
					Symbol:       "=",
					Value:        "=",
					Arity:        2,
					BindingPower: 30,
					Line:         1,
					Col:          1,
				},
				&Token{
					Symbol:       "(STRING)",
					Value:        "Hello",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          1,
				},
				&Token{
					Symbol:       "(EOF)",
					Value:        "",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          15,
				},
			},
		},
		testCase{
			input: "NOT A",
			expected: []*Token{
				&Token{
					Symbol:       "NOT",
					Value:        "NOT",
					Arity:        1,
					BindingPower: 40,
					Line:         1,
					Col:          1,
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "A",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          1,
				},
				&Token{
					Symbol:       "(EOF)",
					Value:        "",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          15,
				},
			},
		},
		testCase{
			input: "A AND B",
			expected: []*Token{
				&Token{
					Symbol:       "(NAME)",
					Value:        "A",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          1,
				},
				&Token{
					Symbol:       "AND",
					Value:        "AND",
					Arity:        2,
					BindingPower: 20,
					Line:         1,
					Col:          3,
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "B",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          7,
				},
				&Token{
					Symbol:       "(EOF)",
					Value:        "",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          8,
				},
			},
		},
		testCase{
			input: "(A OR B) AND C",
			expected: []*Token{
				&Token{
					Symbol:       "(",
					Value:        "(",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          1,
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "A",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          2,
				},
				&Token{
					Symbol:       "OR",
					Value:        "OR",
					Arity:        2,
					BindingPower: 10,
					Line:         1,
					Col:          4,
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "B",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          7,
				},
				&Token{
					Symbol:       ")",
					Value:        ")",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          8,
				},
				&Token{
					Symbol:       "AND",
					Value:        "AND",
					Arity:        2,
					BindingPower: 20,
					Line:         1,
					Col:          10,
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "C",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          14,
				},
				&Token{
					Symbol:       "(EOF)",
					Value:        "",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          15,
				},
			},
		},
		testCase{
			input: "(A = B) AND NOT (C = D)",
			expected: []*Token{
				&Token{
					Symbol:       "(",
					Value:        "(",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          1,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "A",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          2,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "=",
					Value:        "=",
					Arity:        2,
					BindingPower: 30,
					Line:         1,
					Col:          4,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "B",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          7,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       ")",
					Value:        ")",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          8,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "AND",
					Value:        "AND",
					Arity:        2,
					BindingPower: 20,
					Line:         1,
					Col:          10,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "NOT",
					Value:        "NOT",
					Arity:        1,
					BindingPower: 40,
					Line:         1,
					Col:          10,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(",
					Value:        "(",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          1,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "C",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          2,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "=",
					Value:        "=",
					Arity:        2,
					BindingPower: 30,
					Line:         1,
					Col:          4,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "D",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          7,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       ")",
					Value:        ")",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          8,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(EOF)",
					Value:        "",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          15,
					Children:     []*Token{},
				},
			},
		},
		testCase{
			input: "(ABBA = BRAVO) AND NOT (CAPTAIN = DICK)",
			expected: []*Token{
				&Token{
					Symbol:       "(",
					Value:        "(",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          1,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "ABBA",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          2,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "=",
					Value:        "=",
					Arity:        2,
					BindingPower: 30,
					Line:         1,
					Col:          4,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "BRAVO",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          7,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       ")",
					Value:        ")",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          8,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "AND",
					Value:        "AND",
					Arity:        2,
					BindingPower: 20,
					Line:         1,
					Col:          10,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "NOT",
					Value:        "NOT",
					Arity:        1,
					BindingPower: 40,
					Line:         1,
					Col:          10,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(",
					Value:        "(",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          1,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "CAPTAIN",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          2,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "=",
					Value:        "=",
					Arity:        2,
					BindingPower: 30,
					Line:         1,
					Col:          4,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(NAME)",
					Value:        "DICK",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          7,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       ")",
					Value:        ")",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          8,
					Children:     []*Token{},
				},
				&Token{
					Symbol:       "(EOF)",
					Value:        "",
					Arity:        0,
					BindingPower: 0,
					Line:         1,
					Col:          15,
					Children:     []*Token{},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(c *testing.T) {
			lexer := makeLexer(testCase.input)
			tokens := []*Token{}
			loops := 0
			for {
				token, err := lexer.Next()
				if err != nil {
					t.Fatalf("Unexpected parsing error %v", err)
				}
				tokens = append(tokens, token)
				if token.Symbol == EOF || loops > 100 {
					break
				}
				loops++
			}
			for i, tok := range tokens {
				expectedTok := testCase.expected[i]
				equivalent, message := areTokensEquivalent(expectedTok, tok)
				if !equivalent {
					t.Fatal(fmt.Sprintf("Expected %v got %v, ", expectedTok, tok) + message)
				}
			}
		})
	}
}

func areTokensEquivalent(x *Token, y *Token) (bool, string) {
	if len(x.Children) != len(y.Children) {
		return false, "Children of unequal length"
	}
	for i, xChild := range x.Children {
		areEquivalent, message := areTokensEquivalent(xChild, y.Children[i])
		if !areEquivalent {
			return false, fmt.Sprintf("Children at position %v were unequal with message %v", i, message)
		}
	}
	if x.Arity != y.Arity {
		return false, fmt.Sprintf("Arity %v != %v", x.Arity, y.Arity)
	}
	if x.BindingPower != y.BindingPower {
		return false, fmt.Sprintf("BindingPower %v != %v", x.BindingPower, y.BindingPower)
	}
	if x.Line != y.Line {
		return false, fmt.Sprintf("Line %v != %v", x.Line, y.Line)
	}
	if x.Symbol != y.Symbol {
		return false, fmt.Sprintf("Symbol %v != %v", x.Symbol, y.Symbol)
	}
	if x.Value != y.Value {
		return false, fmt.Sprintf("Value %v != %v", x.Value, y.Value)
	}
	return true, ""
}
