package engine

import (
	"fmt"

	"github.com/nicholasbailey/langkit"
)

func (interpreter *ToyScriptInterpreter) callFunction(tree *langkit.Token) (*ToyScriptValue, error) {
	functionName := tree.Children[0].Value
	function, found := interpreter.Functions[functionName]
	if !found {
		return nil, fmt.Errorf("valueerror: unrecognized function name %v at line %v, col %v", functionName, tree.Line, tree.Col)
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
}
