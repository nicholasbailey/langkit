package langkit

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	type testCase struct {
		input  string
		output *Token
	}

	cases := []testCase{
		testCase{
			input: "A",
			output: &Token{
				Symbol:       "(NAME)",
				Value:        "A",
				Line:         1,
				Col:          1,
				Arity:        0,
				BindingPower: 0,
				Children:     []*Token{},
			},
		},
		testCase{
			input: "Name = \"Hello\"",
			output: &Token{
				Symbol:       "=",
				Value:        "=",
				Line:         1,
				Col:          1,
				Arity:        2,
				BindingPower: 30,
				Children: []*Token{
					&Token{
						Symbol:       "(NAME)",
						Value:        "Name",
						Line:         1,
						Arity:        0,
						BindingPower: 0,
						Children:     []*Token{},
					},
					&Token{
						Symbol:       "(STRING)",
						Value:        "Hello",
						Line:         1,
						Arity:        0,
						BindingPower: 0,
						Children:     []*Token{},
					},
				},
			},
		},
		testCase{
			input: "A = C AND B = D",
			output: &Token{
				Symbol:       "AND",
				Value:        "AND",
				Line:         1,
				Arity:        2,
				BindingPower: 20,
				Children: []*Token{
					&Token{
						Symbol:       "=",
						Value:        "=",
						Line:         1,
						Arity:        2,
						BindingPower: 30,
						Children: []*Token{
							&Token{
								Symbol:       "(NAME)",
								Value:        "A",
								Line:         1,
								Arity:        0,
								BindingPower: 0,
								Children:     []*Token{},
							},
							&Token{
								Symbol:       "(NAME)",
								Value:        "C",
								Line:         1,
								Arity:        0,
								BindingPower: 0,
								Children:     []*Token{},
							},
						},
					},
					&Token{
						Symbol:       "=",
						Value:        "=",
						Line:         1,
						Arity:        2,
						BindingPower: 30,
						Children: []*Token{
							&Token{
								Symbol:       "(NAME)",
								Value:        "B",
								Line:         1,
								Arity:        0,
								BindingPower: 0,
								Children:     []*Token{},
							},
							&Token{
								Symbol:       "(NAME)",
								Value:        "D",
								Line:         1,
								Arity:        0,
								BindingPower: 0,
								Children:     []*Token{},
							},
						},
					},
				},
			},
		},
		testCase{
			input: "A AND B",
			output: &Token{
				Symbol:       "AND",
				Value:        "AND",
				Line:         1,
				Arity:        2,
				BindingPower: 20,
				Children: []*Token{
					&Token{
						Symbol:       "(NAME)",
						Value:        "A",
						Line:         1,
						Arity:        0,
						BindingPower: 0,
						Children:     []*Token{},
					},
					&Token{
						Symbol:       "(NAME)",
						Value:        "B",
						Line:         1,
						Arity:        0,
						BindingPower: 0,
						Children:     []*Token{},
					},
				},
			},
		},
		testCase{
			input: "(A AND B)",
			output: &Token{
				Symbol:       "AND",
				Value:        "AND",
				Line:         1,
				Arity:        2,
				BindingPower: 20,
				Children: []*Token{
					&Token{
						Symbol:       "(NAME)",
						Value:        "A",
						Line:         1,
						Arity:        0,
						BindingPower: 0,
						Children:     []*Token{},
					},
					&Token{
						Symbol:       "(NAME)",
						Value:        "B",
						Line:         1,
						Arity:        0,
						BindingPower: 0,
						Children:     []*Token{},
					},
				},
			},
		},
		testCase{
			input: "(A AND (B OR C))",
			output: &Token{
				Symbol:       "AND",
				Value:        "AND",
				Line:         1,
				Arity:        2,
				BindingPower: 20,
				Children: []*Token{
					&Token{
						Symbol:       "(NAME)",
						Value:        "A",
						Line:         1,
						Arity:        0,
						BindingPower: 0,
						Children:     []*Token{},
					},
					&Token{
						Symbol:       "OR",
						Value:        "OR",
						Line:         1,
						Arity:        2,
						BindingPower: 10,
						Children: []*Token{
							&Token{
								Symbol:       "(NAME)",
								Value:        "B",
								Line:         1,
								Arity:        0,
								BindingPower: 0,
								Children:     []*Token{},
							},
							&Token{
								Symbol:       "(NAME)",
								Value:        "C",
								Line:         1,
								Arity:        0,
								BindingPower: 0,
								Children:     []*Token{},
							},
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			lexer := makeLexer(c.input)
			parser := TDOPParser{
				Lexer: lexer,
			}
			tree, err := parser.Expression(0)
			if err != nil {
				t.Fatalf("Unexpected parsing error %v", err)
			}
			equivalent, message := areTokensEquivalent(c.output, tree)
			if !equivalent {
				fmt.Println("Expected")
				fmt.Println(c.output.TreeString(0))
				fmt.Println(tree.TreeString(0))
				fmt.Println("Actual")
				t.Fatal(message)
			}
		})

	}
}
