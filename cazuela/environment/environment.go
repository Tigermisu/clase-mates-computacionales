package environment

import "clase-mates-computacionales/cazuela/lexer"
import "clase-mates-computacionales/cazuela/errorHandler"
import "fmt"

type env interface {
	Define()
	Get() interface{}
}

type Environment struct {
	Values    map[string]interface{}
	Enclosing *Environment
}

func (e Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e Environment) Get(name lexer.Token) interface{} {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	errorHandler.RaiseError(errorHandler.CodeUndefinedVariable, fmt.Sprintf("Variable %v no definida", name.Lexeme), name.Line, "[Ejecución]", true)
	return nil
}

func (e Environment) Assign(name lexer.Token, value interface{}) interface{} {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return value
	}

	if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	}

	errorHandler.RaiseError(errorHandler.CodeUndefinedVariable, fmt.Sprintf("Variable %v no definida", name.Lexeme), name.Line, "[Ejecución]", true)
	return nil
}
