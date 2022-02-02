package main

import (
	"fmt"
	"os"

	"github.com/nicholasbailey/langkit/toyscript/engine"
)

func main() {
	path := os.Args[1]
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	engine := engine.BuildToyscriptEngine()
	_, err = engine.Execute(file)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
