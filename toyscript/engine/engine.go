package engine

import (
	"fmt"
	"strconv"

	"github.com/nicholasbailey/langkit"
)

type ToyscriptType string

const (
	TString ToyscriptType = "string"
	TInt    ToyscriptType = "int"
	TBool   ToyscriptType = "bool"
	TFloat  ToyscriptType = "float"
	TNull   ToyscriptType = "null"
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
		return True(), nil
	case "false":
		return False(), nil
	case langkit.Name:
		value, found := interpreter.GlobalVars[tree.Value]
		if !found {
			return nil, fmt.Errorf("syntaxerror: unbound variable %v at line %v, col %v", tree.Value, tree.Line, tree.Col)
		}
		return value, nil
	// Handle Variable assignment
	case "&&":
		return interpreter.doAnd(tree)
	case "||":
		return interpreter.doOr(tree)
	case "!=":
		return interpreter.doInequalityCheck(tree)
	case "==":
		return interpreter.doEqualityCheck(tree)
	case "=":
		return interpreter.doAssigment(tree)
	case "+":
		return interpreter.doAddition(tree)
	case "-":
		return interpreter.doSubtraction(tree)
	case "*":
		return interpreter.doMultiplication(tree)
	case "/":
		return interpreter.doDivision(tree)
	case "%":
		return interpreter.doModulo(tree)
	case "<":
		return interpreter.doLessThan(tree)
	case ">":
		return interpreter.doGreaterThan(tree)
	case "<=":
		return interpreter.doLessThanOrEqualTo(tree)
	case ">=":
		return interpreter.doGreaterThanOrEqualTo(tree)
	case "while":
		return interpreter.doWhile(tree)
	case "def":
		// return interpreter.defineFunction(tree)
	case langkit.FunctionInvocation:
		return interpreter.callFunction(tree)
	case langkit.Block:
		val, err := interpreter.Execute(tree.Children)
		if err != nil {
			return nil, err
		}
		return val.(*ToyScriptValue), nil
	case "if":
		return interpreter.doIf(tree)
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

func BuildToyscriptEngine() langkit.Engine {
	interpreter := BuildToyscriptInterpreter()
	languageSpec := BuildToyscriptLanguageSpec()
	return langkit.NewEngine(languageSpec, interpreter)
}
