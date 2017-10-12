package interpreter

import (
	"clase-mates-computacionales/cazuela/errorHandler"
	"clase-mates-computacionales/cazuela/lexer"
	"clase-mates-computacionales/cazuela/parser"
	"fmt"
	"math"
)

// Interpret takes an AST and interprets it (magic!)
func Interpret(expr parser.Expression) {
	fmt.Println(evaluate(expr))
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
	case lexer.TokenNotEqualTo:
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

		if isLString && isRString {
			return leftString + rightString
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

func evaluate(expr parser.Expression) interface{} {
	if v, ok := expr.(parser.LiteralExpression); ok {
		return getLiteralValue(v)
	} else if v, ok := expr.(parser.GroupingExpression); ok {
		return getGroupValue(v)
	} else if v, ok := expr.(parser.UnaryExpression); ok {
		return getUnaryValue(v)
	} else if v, ok := expr.(parser.BinaryExpression); ok {
		return getBinaryValue(v)
	}

	return nil
}
