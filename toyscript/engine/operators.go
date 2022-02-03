package engine

import (
	"fmt"

	"github.com/nicholasbailey/langkit"
)

func resolveBinaryOperands(interpreter *ToyScriptInterpreter, tree *langkit.Token) (*ToyScriptValue, *ToyScriptValue, error) {
	if len(tree.Children) != 2 {
		return nil, nil, Exception(SyntaxError, fmt.Sprintf("invalid symbol %v", tree.Value), tree)
	}
	left := tree.Children[0]
	right := tree.Children[1]
	leftValue, leftErr := interpreter.Evaluate(left)

	rightValue, rightErr := interpreter.Evaluate(right)
	if leftErr != nil {
		return leftValue, rightValue, leftErr
	} else {
		return leftValue, rightValue, rightErr
	}
}

func (interpreter *ToyScriptInterpreter) doLessThan(tree *langkit.Token) (*ToyScriptValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.Type == TInt && rightValue.Type == TInt {
		return BoolFromGoBoolean(leftValue.Value.(int64) < rightValue.Value.(int64)), nil
	} else if leftValue.Type == TFloat && rightValue.Type == TFloat {
		return BoolFromGoBoolean(leftValue.Value.(float64) < rightValue.Value.(float64)), nil
	} else if leftValue.Type == TString && rightValue.Type == TString {
		return BoolFromGoBoolean(leftValue.Value.(string) < rightValue.Value.(string)), nil
	} else if leftValue.Type == rightValue.Type {
		return nil, Exception(TypeError, fmt.Sprintf("type %v cannot be compared with <", rightValue.Type), tree)
	}
	return nil, Exception(TypeError, "attempted to compare incomparable types with <", tree)
}

func (interpreter *ToyScriptInterpreter) doGreaterThan(tree *langkit.Token) (*ToyScriptValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.Type == TInt && rightValue.Type == TInt {
		return BoolFromGoBoolean(leftValue.Value.(int64) > rightValue.Value.(int64)), nil
	} else if leftValue.Type == TFloat && rightValue.Type == TFloat {
		return BoolFromGoBoolean(leftValue.Value.(float64) > rightValue.Value.(float64)), nil
	} else if leftValue.Type == TString && rightValue.Type == TString {
		return BoolFromGoBoolean(leftValue.Value.(string) > rightValue.Value.(string)), nil
	} else if leftValue.Type == rightValue.Type {
		return nil, Exception(TypeError, fmt.Sprintf("type %v cannot be compared with >", rightValue.Type), tree)
	}
	return nil, Exception(TypeError, "attempted to compare incomparable types with >", tree)
}

func (interpreter *ToyScriptInterpreter) doLessThanOrEqualTo(tree *langkit.Token) (*ToyScriptValue, error) {
	equal, err := interpreter.doEqualityCheck(tree)
	if err != nil {
		return nil, err
	}
	if equal.Value == true {
		return True(), nil
	}
	lessThan, err := interpreter.doLessThan(tree)
	if err != nil {
		return nil, err
	}
	if lessThan.Value == true {
		return True(), nil
	}
	return False(), nil
}

func (interpreter *ToyScriptInterpreter) doGreaterThanOrEqualTo(tree *langkit.Token) (*ToyScriptValue, error) {
	equal, err := interpreter.doEqualityCheck(tree)
	if err != nil {
		return nil, err
	}
	if equal.Value == true {
		return True(), nil
	}
	lessThan, err := interpreter.doGreaterThan(tree)
	if err != nil {
		return nil, err
	}
	if lessThan.Value == true {
		return True(), nil
	}
	return False(), nil
}

func (interpreter *ToyScriptInterpreter) doEqualityCheck(tree *langkit.Token) (*ToyScriptValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.Type != rightValue.Type {
		return False(), nil
	}
	if leftValue.Value == rightValue.Value {
		return True(), nil
	}
	return False(), nil
}

func (interpreter *ToyScriptInterpreter) doInequalityCheck(tree *langkit.Token) (*ToyScriptValue, error) {
	result, err := interpreter.doEqualityCheck(tree)
	if err != nil {
		return nil, err
	}
	if result.Value == false {
		return True(), nil
	} else {
		return False(), nil
	}
}

func (interpreter *ToyScriptInterpreter) doAnd(tree *langkit.Token) (*ToyScriptValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	leftTruthy := Truthiness(leftValue)
	rightTruthy := Truthiness(rightValue)
	if leftTruthy.Value == true && rightTruthy.Value == true {
		return rightValue, nil
	}
	return leftValue, nil
}

func (interpreter *ToyScriptInterpreter) doOr(tree *langkit.Token) (*ToyScriptValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	leftTruthy := Truthiness(leftValue)
	rightTruthy := Truthiness(rightValue)
	if leftTruthy.Value == true || rightTruthy.Value == true {
		return leftValue, nil
	}
	return rightValue, nil
}

func (interpreter *ToyScriptInterpreter) doAssigment(tree *langkit.Token) (*ToyScriptValue, error) {
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
}

func (interpreter *ToyScriptInterpreter) doAddition(tree *langkit.Token) (*ToyScriptValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
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

func (interpreter *ToyScriptInterpreter) doSubtraction(tree *langkit.Token) (*ToyScriptValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.Type == TInt && rightValue.Type == TInt {
		newValue := leftValue.Value.(int64) - rightValue.Value.(int64)
		return &ToyScriptValue{
			Type:  TInt,
			Value: newValue,
		}, nil
	}
	if leftValue.Type == TFloat && rightValue.Type == TFloat {
		newValue := leftValue.Value.(float64) - rightValue.Value.(float64)
		return &ToyScriptValue{
			Type:  TFloat,
			Value: newValue,
		}, nil
	}
	if leftValue.Type == rightValue.Type {
		return nil, fmt.Errorf("typeerror: type %v does not support operator - at line %v, col %v", leftValue.Type, tree.Line, tree.Col)
	}
	return nil, fmt.Errorf("typerror: incompatable types %v and %v with operator - at line %v, col %v", leftValue.Type, rightValue.Type, tree.Line, tree.Col)
}

func (interpreter *ToyScriptInterpreter) doMultiplication(tree *langkit.Token) (*ToyScriptValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.Type == TInt && rightValue.Type == TInt {
		newValue := leftValue.Value.(int64) * rightValue.Value.(int64)
		return &ToyScriptValue{
			Type:  TInt,
			Value: newValue,
		}, nil
	}
	if leftValue.Type == TFloat && rightValue.Type == TFloat {
		newValue := leftValue.Value.(float64) * rightValue.Value.(float64)
		return &ToyScriptValue{
			Type:  TFloat,
			Value: newValue,
		}, nil
	}
	if leftValue.Type == rightValue.Type {
		return nil, fmt.Errorf("typeerror: type %v does not support operator * at line %v, col %v", leftValue.Type, tree.Line, tree.Col)
	}
	return nil, fmt.Errorf("typerror: incompatable types %v and %v with operator * at line %v, col %v", leftValue.Type, rightValue.Type, tree.Line, tree.Col)
}

func (interpreter *ToyScriptInterpreter) doDivision(tree *langkit.Token) (*ToyScriptValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.Type == TInt && rightValue.Type == TInt {
		if rightValue.Value.(int64) == 0 {
			return nil, fmt.Errorf("dividebyzeroerror: integer division by zero at line %v, col %v", tree.Line, tree.Col)
		}
		newValue := leftValue.Value.(int64) / rightValue.Value.(int64)
		return &ToyScriptValue{
			Type:  TInt,
			Value: newValue,
		}, nil
	}
	if leftValue.Type == TFloat && rightValue.Type == TFloat {
		if rightValue.Value.(int64) == 0.0 {
			return nil, fmt.Errorf("dividebyzeroerror: float division by zero at line %v, col %v", tree.Line, tree.Col)
		}
		newValue := leftValue.Value.(float64) / rightValue.Value.(float64)
		return &ToyScriptValue{
			Type:  TFloat,
			Value: newValue,
		}, nil
	}
	if leftValue.Type == rightValue.Type {
		return nil, fmt.Errorf("typeerror: type %v does not support operator / at line %v, col %v", leftValue.Type, tree.Line, tree.Col)
	}
	return nil, fmt.Errorf("typerror: incompatable types %v and %v with operator / at line %v, col %v", leftValue.Type, rightValue.Type, tree.Line, tree.Col)
}

func (interpreter *ToyScriptInterpreter) doModulo(tree *langkit.Token) (*ToyScriptValue, error) {
	leftValue, rightValue, err := resolveBinaryOperands(interpreter, tree)
	if err != nil {
		return nil, err
	}
	if leftValue.Type == TInt && rightValue.Type == TInt {
		if rightValue.Value.(int64) == 0 {
			return nil, fmt.Errorf("dividebyzeroerror: integer modulo by zero at line %v, col %v", tree.Line, tree.Col)
		}
		newValue := leftValue.Value.(int64) % rightValue.Value.(int64)
		return &ToyScriptValue{
			Type:  TInt,
			Value: newValue,
		}, nil
	}
	if leftValue.Type == rightValue.Type {
		return nil, fmt.Errorf("typeerror: type %v does not support operator %% at line %v, col %v", leftValue.Type, tree.Line, tree.Col)
	}
	return nil, fmt.Errorf("typerror: incompatable types %v and %v with operator %% at line %v, col %v", leftValue.Type, rightValue.Type, tree.Line, tree.Col)
}
