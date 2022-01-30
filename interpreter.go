package langkit

import (
	"errors"
	"fmt"
	"io"
)

type Value interface{}
type Exception error

type Engine interface {
	Execute(source io.Reader) (Value, Exception)
	ExecuteWithGlobals(source io.Reader, globals VariableValues) (Value, Exception)
}

type Interpreter interface {
	Evaluate(tree *Token, variableValues VariableValues) (Value, Exception)
}

func NewEngine(symbolTable SymbolTable, intepreter Interpreter) Engine {
	lexerFactory := func(source io.Reader) Lexer {
		return NewLexer(source, symbolTable)
	}
	return &EngineImpl{
		lexerFactory: lexerFactory,
		interpreter:  intepreter,
	}
}

type EngineImpl struct {
	lexerFactory func(io.Reader) Lexer
	interpreter  Interpreter
}

func (engine *EngineImpl) ExecuteWithGlobals(source io.Reader, globals VariableValues) (Value, Exception) {
	lexer := engine.lexerFactory(source)
	parser := TDOPParser{
		Lexer: lexer,
	}
	tree := parser.Expression(0)
	return engine.interpreter.Evaluate(tree, globals)
}

func (engine *EngineImpl) Execute(source io.Reader) (Value, Exception) {
	return engine.ExecuteWithGlobals(source, map[string]interface{}{})
}

type InterpreterFunction struct {
	evaluator func([]interface{}) interface{}
	arity     int
}
type VariableValues map[string]interface{}

type BuiltinsRegistry map[Symbol]InterpreterFunction

func (registry BuiltinsRegistry) Define(name string, arity int, f func([]interface{}) interface{}) {
	registry[Symbol(name)] = InterpreterFunction{
		arity:     arity,
		evaluator: f,
	}
}

func NewBuiltinsRegistry() BuiltinsRegistry {
	return map[Symbol]InterpreterFunction{}
}

type InterpreterImpl struct {
	Functions BuiltinsRegistry
}

func (interpreter *InterpreterImpl) Evaluate(tree *Token, variableValues VariableValues) (Value, Exception) {
	if len(tree.Children) == 0 {
		if tree.Symbol == NameSymbol {
			val, found := variableValues[tree.Value]
			if !found {
				return nil, errors.New(fmt.Sprintf("Unbound variable %v at line %v, col %v", tree.Value, tree.Line, tree.Col))
			}
			return val, nil
		} else if tree.Symbol == StringLiteral {
			return tree.Value, nil
		} else {
			return nil, errors.New(fmt.Sprintf("Unrecognized literal %v at line %v, col %v", tree.Value, tree.Line, tree.Col))
		}

	}
	f, found := interpreter.Functions[tree.Symbol]
	if !found {
		return nil, errors.New(fmt.Sprintf("Undefined function or operator %v at line %v, col %v", tree.Symbol, tree.Line, tree.Col))
	}
	if f.arity != len(tree.Children) {
		return nil, errors.New(fmt.Sprintf("Incorrect number of arguments to function or operator %v, expected %v arguments, got %v at line %v, col %v", tree.Symbol, f.arity, tree.Arity, tree.Line, tree.Col))
	}
	childValues := []interface{}{}
	for _, subTree := range tree.Children {
		value, err := interpreter.Evaluate(subTree, variableValues)
		if err != nil {
			return nil, err
		}
		childValues = append(childValues, value)
	}
	return f.evaluator(childValues), nil
}
