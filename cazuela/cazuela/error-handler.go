package cazuela

import (
	"fmt"
	"os"
)

// error codes
const (
	CodeAllGood          = 0x00
	CodeSyntaxError      = 0x01
	CodeTooManyArguments = 0x02
	CodeRuntimeError     = 0x03
	CodeUnexpectedEOF    = 0x04
)

// IgnoreFatals when true prevents the program from exiting during fatal errors
var IgnoreFatals = false

// A MenudoError represents a general error during the interpretation of the program
type MenudoError struct {
	code    int
	line    int
	message string
	context string
}

func (e MenudoError) String() string {
	return fmt.Sprintf("[%d] Error %v: %v", e.line, e.context, e.message)
}

// HaltExecutionWithError stops the program reporting the error message
func HaltExecutionWithError(mError MenudoError) {
	fmt.Printf("\nLa cazuela se vació con el código: %X\n", mError.code)
	fmt.Printf("\t%v", mError)
	os.Exit(mError.code)
}

// RaiseError creates a new error and outputs, halting if needed
func RaiseError(code int, message string, line int, context string, fatal bool) {
	mError := MenudoError{code, line, message, context}

	if fatal && !IgnoreFatals {
		HaltExecutionWithError(mError)
	} else {
		fmt.Println(mError)
	}
}

// RaiseErrorWithCode creates a new error giving just a code, inferring the message
func RaiseErrorWithCode(code int) {
	message := getErrorCodeDescription(code)
	RaiseError(code, message, 0, "", false)

}

func getErrorCodeDescription(errorCode int) string {
	switch errorCode {
	case CodeAllGood:
		return "Ejecución normal"
	case CodeSyntaxError:
		return "Error de sintaxis"
	case CodeTooManyArguments:
		return "Demasiados argumentos durante inicialización"
	case CodeRuntimeError:
		return "Error en tiempo de ejecución"
	}
	return "Error desconocido"
}
