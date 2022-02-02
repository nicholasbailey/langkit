package langkit

// import (
// 	"testing"
// )

// func makeInterpreter() Interpreter {
// 	functions := map[Symbol]InterpreterFunction{
// 		Symbol("="): InterpreterFunction{
// 			arity: 2,
// 			evaluator: func(inputs []interface{}) interface{} {
// 				return inputs[0] == inputs[1]
// 			},
// 		},
// 		Symbol("AND"): InterpreterFunction{
// 			arity: 2,
// 			evaluator: func(inputs []interface{}) interface{} {
// 				left := inputs[0].(bool)
// 				right := inputs[1].(bool)
// 				return left && right
// 			},
// 		},
// 		Symbol("OR"): InterpreterFunction{
// 			arity: 2,
// 			evaluator: func(inputs []interface{}) interface{} {
// 				left := inputs[0].(bool)
// 				right := inputs[1].(bool)
// 				return left || right
// 			},
// 		},
// 		Symbol("NOT"): InterpreterFunction{
// 			arity: 1,
// 			evaluator: func(inputs []interface{}) interface{} {
// 				arg := inputs[0].(bool)
// 				return !arg
// 			},
// 		},
// 	}

// 	return &InterpreterImpl{
// 		Functions: functions,
// 	}
// }

// func TestInterpeterHandlesSingleValue(t *testing.T) {
// 	type testCase struct {
// 		name          string
// 		input         *Token
// 		variables     VariableValues
// 		expectedError error
// 		expectedValue interface{}
// 	}

// 	testCases := []testCase{
// 		{
// 			name: "Evaluate a variable",
// 			input: &Token{
// 				Symbol:       "(NAME)",
// 				Value:        "A",
// 				Arity:        0,
// 				BindingPower: 0,
// 				Line:         1,
// 				Col:          1,
// 				Children:     []*Token{},
// 			},
// 			variables: map[string]interface{}{
// 				"A": true,
// 			},
// 			expectedError: nil,
// 			expectedValue: true,
// 		},
// 		{
// 			name: "Evaluate Equality values equal",
// 			input: &Token{
// 				Symbol:       "=",
// 				Value:        "=",
// 				Arity:        2,
// 				BindingPower: 0,
// 				Line:         1,
// 				Col:          1,
// 				Children: []*Token{
// 					{
// 						Symbol:       "(NAME)",
// 						Value:        "A",
// 						Arity:        0,
// 						BindingPower: 0,
// 						Line:         1,
// 						Col:          1,
// 						Children:     []*Token{},
// 					},
// 					{
// 						Symbol:       "(NAME)",
// 						Value:        "B",
// 						Arity:        0,
// 						BindingPower: 0,
// 						Line:         1,
// 						Col:          1,
// 						Children:     []*Token{},
// 					},
// 				},
// 			},
// 			variables: map[string]interface{}{
// 				"A": "Kitty",
// 				"B": "Kitty",
// 			},
// 			expectedValue: true,
// 			expectedError: nil,
// 		},
// 		{
// 			name: "Evaluate Equality values not equal",
// 			input: &Token{
// 				Symbol:       "=",
// 				Value:        "=",
// 				Arity:        2,
// 				BindingPower: 0,
// 				Line:         1,
// 				Col:          1,
// 				Children: []*Token{
// 					{
// 						Symbol:       "(NAME)",
// 						Value:        "A",
// 						Arity:        0,
// 						BindingPower: 0,
// 						Line:         1,
// 						Col:          1,
// 						Children:     []*Token{},
// 					},
// 					{
// 						Symbol:       "(NAME)",
// 						Value:        "B",
// 						Arity:        0,
// 						BindingPower: 0,
// 						Line:         1,
// 						Col:          1,
// 						Children:     []*Token{},
// 					},
// 				},
// 			},
// 			variables: map[string]interface{}{
// 				"A": "Kitty",
// 				"B": "Doggy",
// 			},
// 			expectedValue: false,
// 			expectedError: nil,
// 		},
// 		{
// 			name: "Evaluate And (true, true)",
// 			input: &Token{
// 				Symbol:       "AND",
// 				Value:        "AND",
// 				Arity:        2,
// 				BindingPower: 10,
// 				Line:         1,
// 				Col:          1,
// 				Children: []*Token{
// 					{
// 						Symbol:       "=",
// 						Value:        "=",
// 						Arity:        2,
// 						BindingPower: 0,
// 						Line:         1,
// 						Col:          1,
// 						Children: []*Token{
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "A",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "B",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 						},
// 					},
// 					{
// 						Symbol:       "=",
// 						Value:        "=",
// 						Arity:        2,
// 						BindingPower: 0,
// 						Line:         1,
// 						Col:          1,
// 						Children: []*Token{
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "C",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "D",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			variables: map[string]interface{}{
// 				"A": "Kitty",
// 				"B": "Kitty",
// 				"C": "Doggy",
// 				"D": "Doggy",
// 			},
// 			expectedValue: true,
// 			expectedError: nil,
// 		},
// 		{
// 			name: "Evaluate And (true, false)",
// 			input: &Token{
// 				Symbol:       "AND",
// 				Value:        "AND",
// 				Arity:        2,
// 				BindingPower: 10,
// 				Line:         1,
// 				Col:          1,
// 				Children: []*Token{
// 					{
// 						Symbol:       "=",
// 						Value:        "=",
// 						Arity:        2,
// 						BindingPower: 0,
// 						Line:         1,
// 						Col:          1,
// 						Children: []*Token{
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "A",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "B",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 						},
// 					},
// 					{
// 						Symbol:       "=",
// 						Value:        "=",
// 						Arity:        2,
// 						BindingPower: 0,
// 						Line:         1,
// 						Col:          1,
// 						Children: []*Token{
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "C",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "D",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			variables: map[string]interface{}{
// 				"A": "Kitty",
// 				"B": "Kitty",
// 				"C": "Doggy",
// 				"D": "Birdy",
// 			},
// 			expectedValue: false,
// 			expectedError: nil,
// 		},
// 		{
// 			name: "Evaluate And (false, false)",
// 			input: &Token{
// 				Symbol:       "AND",
// 				Value:        "AND",
// 				Arity:        2,
// 				BindingPower: 10,
// 				Line:         1,
// 				Col:          1,
// 				Children: []*Token{
// 					{
// 						Symbol:       "=",
// 						Value:        "=",
// 						Arity:        2,
// 						BindingPower: 0,
// 						Line:         1,
// 						Col:          1,
// 						Children: []*Token{
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "A",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "B",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 						},
// 					},
// 					{
// 						Symbol:       "=",
// 						Value:        "=",
// 						Arity:        2,
// 						BindingPower: 0,
// 						Line:         1,
// 						Col:          1,
// 						Children: []*Token{
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "C",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 							{
// 								Symbol:       "(NAME)",
// 								Value:        "D",
// 								Arity:        0,
// 								BindingPower: 0,
// 								Line:         1,
// 								Col:          1,
// 								Children:     []*Token{},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			variables: map[string]interface{}{
// 				"A": "Fishy",
// 				"B": "Kitty",
// 				"C": "Doggy",
// 				"D": "Birdy",
// 			},
// 			expectedValue: false,
// 			expectedError: nil,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			interpreter := makeInterpreter()
// 			value, err := interpreter.EvaluateWithValues(tc.input, tc.variables)
// 			if err != tc.expectedError {
// 				t.Fatalf("Expected error %v got %v", tc.expectedError, err)
// 			}
// 			if value != tc.expectedValue {
// 				t.Fatalf("Expected %v, got %v", tc.expectedValue, value)
// 			}
// 		})
// 	}
// }
