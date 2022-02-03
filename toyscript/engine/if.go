package engine

import (
	"github.com/nicholasbailey/langkit"
)

func (interpreter *ToyScriptInterpreter) doIf(tree *langkit.Token) (*ToyScriptValue, error) {
	if len(tree.Children) < 2 {
		Exception(SyntaxError, "invalid if expression", tree)
	}
	condition := tree.Children[0]
	conditionValue, err := interpreter.Evaluate(condition)
	if err != nil {
		return nil, err
	}
	executeCondition := Truthiness(conditionValue)
	if executeCondition.Value == true {
		block := tree.Children[1]
		return interpreter.Evaluate(block)
	}
	for _, child := range tree.Children[2:] {
		if child.Symbol == langkit.ElseIf {
			return interpreter.doIf(child)
		} else {
			return interpreter.Evaluate(child)
		}

	}
	return Null(), nil
}
