package engine

import "github.com/nicholasbailey/langkit"

func (interpreter *ToyScriptInterpreter) doFor(tree *langkit.Token) (*ToyScriptValue, error) {
	return nil, nil
}

func (interpreter *ToyScriptInterpreter) doWhile(tree *langkit.Token) (*ToyScriptValue, error) {
	if len(tree.Children) != 2 {
		return nil, Exception(SyntaxError, "invalid while block", tree)
	}
	expression := tree.Children[0]
	block := tree.Children[1]
	retVal := Null()
	for {
		expressionRes, err := interpreter.Evaluate(expression)
		if err != nil {
			return nil, err
		}
		expressionTruthiness := Truthiness(expressionRes)
		if expressionTruthiness.Value == false {
			break
		}
		retVal, err = interpreter.Evaluate(block)
		if err != nil {
			return nil, err
		}
	}
	return retVal, nil
}
