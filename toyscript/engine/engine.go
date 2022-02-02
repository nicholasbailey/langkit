package engine

import (
	"fmt"
	"strconv"

	"github.com/nicholasbailey/langkit"
)

type ToyscriptType string

const (
	TString ToyscriptType = "string"
	TInt                  = "int"
	TBool                 = "bool"
	TFloat                = "float"
	TNull                 = "null"
)

type ToyScriptValue struct {
	Type  ToyscriptType
	Value interface{}
}

func Null() *ToyScriptValue {
	return &ToyScriptValue{
		Type:  TNull,
		Value: nil,
	}
}

type ToyscriptFunction func(values []*ToyScriptValue) *ToyScriptValue

type ToyScriptInterpreter struct {
	GlobalVars map[string]*ToyScriptValue
	Functions  map[string]ToyscriptFunction
}

func (interpreter *ToyScriptInterpreter) EvaluateWithValues(tree *langkit.Token, variableValues langkit.VariableValues) (langkit.Value, langkit.Exception) {
	val, err := interpreter.Evaluate(tree)
	return val, err
}

func (interpreter *ToyScriptInterpreter) Execute(statements []*langkit.Token) (langkit.Value, langkit.Exception) {
	var value interface{} = nil
	var err error = nil
	for _, statement := range statements {
		value, err = interpreter.Evaluate(statement)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

func (interpreter *ToyScriptInterpreter) Evaluate(tree *langkit.Token) (*ToyScriptValue, langkit.Exception) {
	switch tree.Symbol {
	case langkit.StringLiteral:
		return &ToyScriptValue{
			Type:  TString,
			Value: tree.Value,
		}, nil
	case langkit.IntLiteral:
		parsedInt, err := strconv.ParseInt(tree.Value, 0, 64)
		if err != nil {
			// TODO: Make this a proper exception
			return nil, err
		}
		return &ToyScriptValue{
			Type:  TInt,
			Value: parsedInt,
		}, nil
	case langkit.FloatLiteral:
		parsedFloat, err := strconv.ParseFloat(tree.Value, 64)
		if err != nil {
			return nil, err
		}
		return &ToyScriptValue{
			Type:  TFloat,
			Value: parsedFloat,
		}, nil
	case "true":
		return &ToyScriptValue{
			Type:  TBool,
			Value: true,
		}, nil
	case "false":
		return &ToyScriptValue{
			Type:  TBool,
			Value: true,
		}, nil
	case langkit.Name:
		value, found := interpreter.GlobalVars[tree.Value]
		if !found {
			return nil, fmt.Errorf("syntaxerror: unbound variable %v at line %v, col %v", tree.Value, tree.Line, tree.Col)
		}
		return value, nil
	// Handle Variable assignment
	case "=":
		if len(tree.Children) != 2 {
			// TODO make more detailed
			return nil, fmt.Errorf("syntaxerror: invald assignment expression at line %v, col %v", tree.Line, tree.Col)
		}
		left := tree.Children[0]
		right := tree.Children[1]
		if left.Symbol != langkit.Name {
			return nil, fmt.Errorf("syntaxerror: invalid assigment expression at line %v, col %v", tree.Line, tree.Col)
		}
		rightValue, err := interpreter.Evaluate(right)
		if err != nil {
			return nil, err
		}
		interpreter.GlobalVars[left.Value] = rightValue
		return rightValue, nil
	case langkit.FunctionInvocation:
		functionName := tree.Children[0].Value
		function, found := interpreter.Functions[functionName]
		if !found {
			return nil, fmt.Errorf("syntaxerror: unrecognized function name %v at line %v, col %v", functionName, tree.Line, tree.Col)
		}
		// TODO - optimize memory allocation here
		childValues := []*ToyScriptValue{}
		for _, childToken := range tree.Children[1:] {
			childValue, err := interpreter.Evaluate(childToken)
			if err != nil {
				return nil, err
			}
			childValues = append(childValues, childValue)
		}
		value := function(childValues)
		return value, nil
	case "+":
		if len(tree.Children) != 2 {
			return nil, fmt.Errorf("syntaxerror: invald addition expression at line %v, col %v", tree.Line, tree.Col)
		}
		left := tree.Children[0]
		right := tree.Children[1]
		leftValue, err := interpreter.Evaluate(left)
		if err != nil {
			return nil, err
		}
		rightValue, err := interpreter.Evaluate(right)
		if err != nil {
			return nil, err
		}
		if leftValue.Type == TInt && rightValue.Type == TInt {
			newValue := leftValue.Value.(int64) + rightValue.Value.(int64)
			return &ToyScriptValue{
				Type:  TInt,
				Value: newValue,
			}, nil
		}
		if leftValue.Type == TFloat && rightValue.Type == TFloat {
			newValue := leftValue.Value.(float64) + rightValue.Value.(float64)
			return &ToyScriptValue{
				Type:  TFloat,
				Value: newValue,
			}, nil
		}
		if leftValue.Type == TString && rightValue.Type == TString {
			newValue := leftValue.Value.(string) + rightValue.Value.(string)
			return &ToyScriptValue{
				Type:  TFloat,
				Value: newValue,
			}, nil
		}
		if leftValue.Type == rightValue.Type {
			return nil, fmt.Errorf("typeerror: type %v does not support operator + at line %v, col %v", leftValue.Type, tree.Line, tree.Col)
		}
		return nil, fmt.Errorf("typerror: incompatable types %v and %v with operator + at line %v, col %v", leftValue.Type, rightValue.Type, tree.Line, tree.Col)

	}
	return nil, fmt.Errorf("syntaxerror: unrecognized symbol %v at line %v, col %v", tree.Value, tree.Line, tree.Col)
}

func BuildToyscriptInterpreter() *ToyScriptInterpreter {
	print := func(values []*ToyScriptValue) *ToyScriptValue {
		for _, value := range values {
			switch value.Type {
			case TString:
				fmt.Print(value.Value.(string))
			case TInt:
				// TODO - move away from builtin
				fmt.Print(value.Value.(int64))
			case TBool:
				fmt.Print(value.Value.(bool))
			case TFloat:
				fmt.Print(value.Value.(float64))
			case TNull:
				fmt.Print("<null>")
			}
			fmt.Print(" ")
		}
		fmt.Print("\n")
		return Null()
	}
	interpreter := &ToyScriptInterpreter{
		Functions:  map[string]ToyscriptFunction{},
		GlobalVars: map[string]*ToyScriptValue{},
	}
	interpreter.Functions["print"] = print
	return interpreter
}

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
	spec.DefineInfix("+", 60)
	spec.DefineInfix("-", 60)
	spec.DefineInfix("*", 70)
	spec.DefineInfix("/", 70)
	spec.DefineInfix("%", 70)
	spec.DefineStatementTerminator(";")
	spec.DefineEmpty(",")

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
	return spec
}

func BuildToyscriptEngine() langkit.Engine {
	interpreter := BuildToyscriptInterpreter()
	languageSpec := BuildToyscriptLanguageSpec()
	return langkit.NewEngine(languageSpec, interpreter)
}
