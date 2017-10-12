package main

import (
	"clase-mates-computacionales/cazuela/errorHandler"
	"clase-mates-computacionales/cazuela/interpreter"
	"clase-mates-computacionales/cazuela/lexer"
	"clase-mates-computacionales/cazuela/parser"
	"clase-mates-computacionales/utilities"
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Uso: cazuela [archivo]")
		errorHandler.RaiseErrorWithCode(errorHandler.CodeTooManyArguments)
	} else if len(args) == 2 {
		file := utilities.LoadFile(args[1])
		execute(file)
	} else {
		startLineInterpreter()
	}
}

func startLineInterpreter() {
	errorHandler.IgnoreFatals = true
	for {
		fmt.Print("<Cazuela># ")
		input := utilities.GetConsoleInput()
		execute(input)
	}
}

func execute(command string) {
	tokens := lexer.GetTokens(command)

	if errorHandler.HasFatalled {
		errorHandler.HasFatalled = false
		return
	}

	expr := parser.Parse(tokens)

	if errorHandler.HasFatalled {
		errorHandler.HasFatalled = false
		return
	}

	interpreter.Interpret(expr)
}
