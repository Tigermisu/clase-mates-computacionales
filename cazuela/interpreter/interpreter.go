package interpreter

import (
	"clase-mates-computacionales/cazuela/environment"
	"clase-mates-computacionales/cazuela/errorHandler"
	"clase-mates-computacionales/cazuela/lexer"
	"clase-mates-computacionales/cazuela/parser"
	"fmt"
	"math"
	"strconv"
)

var env *environment.Environment

var ShouldPrintAllExpressions = false

type Callable interface {
	arity() int
	Call([]interface{}) interface{}
}

type CazuelaFunction struct {
	declaration parser.FnDecl
}

func (f CazuelaFunction) Call(arguments []interface{}) (response interface{}) {
	localEnv := environment.Environment{make(map[string]interface{}), env}

	for i := 0; i < len(arguments); i++ {
		localEnv.Define(f.declaration.Parameters[i].Lexeme, arguments[i])
	}

	defer func() {
		if returnValue := recover(); returnValue != nil {
			response = returnValue
		}
	}()

	executeBlock(f.declaration.Body, localEnv)

	return
}

func (f CazuelaFunction) arity() int {
	return len(f.declaration.Parameters)
}

func InitEnv() {
	env = &environment.Environment{make(map[string]interface{}), nil}

	env.Define("pi", 3.141592653589793)
	env.Define("e", 2.718281828459045)
}

// Interpret takes an AST and interprets it (magic!)
func Interpret(stmts []parser.Stmt) {
	defer func() {
		if err := recover(); err != nil {
			errorHandler.RaiseError(errorHandler.CodeRuntimeError, "Error interno en tiempo de ejecución", -1, fmt.Sprintf("%v", err), true)
		}
	}()

	for _, s := range stmts {
		execute(s)
	}
}

func execute(s parser.Stmt) {
	if v, ok := s.(parser.Statement); ok {
		evaluateStatement(v)
	} else if v, ok := s.(parser.Print); ok {
		evaluatePrint(v)
	} else if v, ok := s.(parser.Declaration); ok {
		evaluateDeclaration(v)
	} else if v, ok := s.(parser.Block); ok {
		executeBlock(v.Statements, environment.Environment{make(map[string]interface{}), env})
	} else if v, ok := s.(parser.If); ok {
		executeIf(v)
	} else if v, ok := s.(parser.While); ok {
		executeWhile(v)
	} else if v, ok := s.(parser.FnDecl); ok {
		fn := CazuelaFunction{v}
		env.Define(v.Name.Lexeme, fn)
	} else if v, ok := s.(parser.ReturnStmt); ok {
		executeReturn(v)
	}
}

func executeWhile(v parser.While) {
	for isTruthy(evaluate(v.Condition)) {
		execute(v.Body)
	}
}

func executeReturn(v parser.ReturnStmt) {
	var value interface{}
	if v.Value != nil {
		value = evaluate(v.Value)
	}

	panic(value)
}

func executeBlock(statements []parser.Stmt, localEnv environment.Environment) {
	previousEnv := env

	defer func() {
		env = previousEnv
	}()

	env = &localEnv

	for _, statement := range statements {
		execute(statement)
	}

}

func evaluateDeclaration(st parser.Declaration) {
	var value interface{}
	if st.Initializer != nil {
		value = evaluate(st.Initializer)
	}

	env.Define(st.Name.Lexeme, value)
}

func evaluateStatement(st parser.Statement) {
	r := evaluate(st.Expr)

	if ShouldPrintAllExpressions {
		fmt.Printf("<| %v |>\n", r)
	}
}

func evaluatePrint(st parser.Print) {
	value := evaluate(st.Expr)
	fmt.Println(value)
}

func executeIf(ifStmt parser.If) {
	if isTruthy(evaluate(ifStmt.Condition)) {
		execute(ifStmt.ThenBranch)
	} else if ifStmt.ElseBranch != nil {
		execute(ifStmt.ElseBranch)
	}
}

func getLiteralValue(expr parser.LiteralExpression) interface{} {
	return expr.Value
}

func getGroupValue(expr parser.GroupingExpression) interface{} {
	return evaluate(expr.Expression)
}

func getUnaryValue(expr parser.UnaryExpression) interface{} {
	right := evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case lexer.TokenMinus:
		checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	case lexer.TokenNegation:
		return !isTruthy(right)
	}

	return nil
}

func getBinaryValue(expr parser.BinaryExpression) interface{} {
	left := evaluate(expr.Left)
	right := evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case lexer.TokenMinus:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case lexer.TokenDivision:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case lexer.TokenMult:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	case lexer.TokenModulo:
		checkNumberOperands(expr.Operator, left, right)
		return float64(int(left.(float64)) % int(right.(float64)))
	case lexer.TokenExponentation:
		checkNumberOperands(expr.Operator, left, right)
		return math.Pow(left.(float64), right.(float64))
	case lexer.TokenPlus:
		l, isLFloat := left.(float64)
		r, isRFloat := right.(float64)

		if isLFloat && isRFloat {
			return l + r
		}

		leftString, isLString := left.(string)
		rightString, isRString := right.(string)

		if isLString || isRString {
			if isLString && !isRString {
				return leftString + strconv.FormatFloat(r, 'f', -1, 64)
			} else if !isLString && isRString {
				return strconv.FormatFloat(l, 'f', -1, 64) + rightString
			} else {
				return leftString + rightString
			}
		}

		errorHandler.RaiseError(errorHandler.CodeRuntimeError, fmt.Sprintf("Se esperaba números o cadenas para %v", expr.Operator.Lexeme), expr.Operator.Line, "[Suma]", true)

		return nil
	case lexer.TokenGreaterThan:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case lexer.TokenGreaterEqual:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case lexer.TokenLessThan:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case lexer.TokenLessEqual:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case lexer.TokenNotEqualTo:
		return !isEqual(left, right)
	case lexer.TokenEqualEqual:
		return isEqual(left, right)

	}

	return nil
}

func checkNumberOperand(operator lexer.Token, operand interface{}) {
	if v, ok := operand.(float64); !ok {
		errorHandler.RaiseError(errorHandler.CodeRuntimeError, fmt.Sprintf("Se esperaba un número para %v, se obtuvo %v", operator.Lexeme, v), operator.Line, "[Unaria]", true)
	}
}

func checkNumberOperands(operator lexer.Token, left interface{}, right interface{}) {
	l, isNum := left.(float64)
	r, isNumR := right.(float64)

	if !(isNum && isNumR) {
		errorHandler.RaiseError(errorHandler.CodeRuntimeError, fmt.Sprintf("Se esperaban números para %v, se obtuvo %v y %v", operator.Lexeme, l, r), operator.Line, "[Binaria]", true)
	}
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return a == b
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}

	if b, ok := v.(bool); ok {
		return b
	}

	return true
}

func evaluateLogicalExpression(expr parser.LogicalExpression) interface{} {
	left := evaluate(expr.Left)

	if expr.Operator.TokenType == lexer.TokenOr {
		if isTruthy(left) {
			return left
		}
	} else if !isTruthy(left) {
		return left
	}

	return evaluate(expr.Right)
}

func evaluateCallExpression(expr parser.CallExpression) interface{} {
	callee := evaluate(expr.Callee)

	arguments := make([]interface{}, 0)

	for _, arg := range expr.Arguments {
		arguments = append(arguments, evaluate(arg))
	}

	if fn, ok := callee.(Callable); ok {
		if len(arguments) != fn.arity() {
			errorHandler.RaiseError(errorHandler.CodeRuntimeError, fmt.Sprintf("Se esperaban %d argumentos pero se recibieron %d", fn.arity(), len(arguments)), expr.ClosingParenteses.Line, "Función", true)
		}
		received := fn.Call(arguments)
		return received
	}

	errorHandler.RaiseError(errorHandler.CodeRuntimeError, "Se intentó llamar algo que no es una función", expr.ClosingParenteses.Line, "Función", true)
	return nil

}

func evaluate(expr parser.Expression) interface{} {
	if v, ok := expr.(parser.LiteralExpression); ok {
		return getLiteralValue(v)
	} else if v, ok := expr.(parser.GroupingExpression); ok {
		return getGroupValue(v)
	} else if v, ok := expr.(parser.UnaryExpression); ok {
		return getUnaryValue(v)
	} else if v, ok := expr.(parser.BinaryExpression); ok {
		return getBinaryValue(v)
	} else if v, ok := expr.(parser.VariableExpression); ok {
		return env.Get(v.Name)
	} else if v, ok := expr.(parser.AssignmentExpression); ok {
		value := evaluate(v.Value)
		return env.Assign(v.Name, value)
	} else if v, ok := expr.(parser.LogicalExpression); ok {
		return evaluateLogicalExpression(v)
	} else if v, ok := expr.(parser.CallExpression); ok {
		return evaluateCallExpression(v)
	}

	return nil
}
