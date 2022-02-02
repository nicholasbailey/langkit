package engine

import (
	"fmt"
	"os"
	"testing"

	"github.com/nicholasbailey/langkit"
)

func TestToyscriptLexer(t *testing.T) {
	f, _ := os.Open("./test_scripts/test_simple.toy")
	spec := BuildToyscriptLanguageSpec()
	lexer := langkit.NewLexer(f, spec)
	token, _ := lexer.Next()
	for token.Symbol != langkit.EOF {
		fmt.Printf("%v", token.TreeString(0))
		token, _ = lexer.Next()
	}
}

func TestToyscriptParser(t *testing.T) {
	f, _ := os.Open("./test_scripts/test_simple.toy")
	spec := BuildToyscriptLanguageSpec()
	lexer := langkit.NewLexer(f, spec)
	parser := langkit.NewParser(lexer)
	tokens, err := parser.Statements()
	if err != nil {
		t.Fatalf("%v", err)
	} else {
		for _, token := range tokens {
			fmt.Printf("%v\n", token.TreeString(0))
		}
	}
}

func TestToyscriptInterpreter(t *testing.T) {
	f, _ := os.Open("./test_scripts/test_simple.toy")
	engine := BuildToyscriptEngine()
	val, err := engine.Execute(f)
	fmt.Printf("%v err: %v", val, err)
}
