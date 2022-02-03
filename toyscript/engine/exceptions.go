package engine

import (
	"fmt"

	"github.com/nicholasbailey/langkit"
)

type ExceptionType string

const (
	SyntaxError       ExceptionType = "SyntaxError"
	DivideByZeroError ExceptionType = "DivideByZero"
	TypeError         ExceptionType = "TypeError"
)

func Exception(
	exceptionType ExceptionType,
	message string,
	token *langkit.Token) langkit.Exception {
	return fmt.Errorf("%v: %v at %v:%v", exceptionType, message, token.Line, token.Col)
}
