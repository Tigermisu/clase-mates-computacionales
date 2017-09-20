package main

import (
	"clase-mates-computacionales/cazuela/cazuela"
	"clase-mates-computacionales/utilities"
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Uso: cazuela [archivo]")
		cazuela.RaiseErrorWithCode(cazuela.CodeTooManyArguments)
	} else if len(args) == 2 {
		file := utilities.LoadFile(args[1])
		execute(file)
	} else {
		startLineInterpreter()
	}
}

func startLineInterpreter() {
	cazuela.IgnoreFatals = true
	for {
		fmt.Print("<Cazuela># ")
		input := utilities.GetConsoleInput()
		execute(input)
	}
}

func execute(command string) {
	tokens := cazuela.GetTokens(command)
	for _, v := range tokens {
		fmt.Println(v)
	}
}
