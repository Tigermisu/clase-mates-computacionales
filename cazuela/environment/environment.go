package environment

import "clase-mates-computacionales/cazuela/lexer"
import "clase-mates-computacionales/cazuela/errorHandler"

type env interface {
	Define()
	Get() interface{}
}

type Environment struct {
	Values map[string]interface{}
}

func (e Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e Environment) Get(name lexer.Token) interface{} {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val
	}

	errorHandler.RaiseError(errorHandler.CodeUndefinedVariable, "Variable no definida", name.Line, "[Ejecución]", true)
	return nil
}
