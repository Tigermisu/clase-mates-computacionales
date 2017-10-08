package parser

import "clase-mates-computacionales/cazuela/lexer"
import "clase-mates-computacionales/cazuela/errorHandler"
import "fmt"

/*
Based on the formal grammar of Lox (http://www.craftinginterpreters.com/)

expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → addition ( ( ">" | ">=" | "<" | "<=" ) addition )* ;
addition       → multiplication ( ( "-" | "+" ) multiplication )* ;
multiplication → unary ( ( "/" | "*" | "%" ) unary )* ;
unary          → ( "!" | "-" ) unary ;
               | primary ;
primary        → NUMBER | STRING | "false" | "true" | "nil"
               | "(" expression ")" ;

*/

const (
	TypeBinary   = 0x10
	TypeLiteral  = 0x11
	TypeGrouping = 0x12
	TypeUnary    = 0x13
)

// Expression is the base interface for all expressions
type Expression interface {
	getType() int
}

// A BinaryExpression holds a expression with two other expressions and an operator in the middle
type BinaryExpression struct {
	left     Expression
	operator lexer.Token
	right    Expression
}

// A LiteralExpression holds a simple literal
type LiteralExpression struct {
	value interface{}
}

// A GroupingExpression holds more expressions inside it :D
type GroupingExpression struct {
	expression Expression
}

// An UnaryExpression has only one operator and an expression to the right of it
type UnaryExpression struct {
	operator lexer.Token
	right    Expression
}

func (be BinaryExpression) getType() int {
	return TypeBinary
}

func (le LiteralExpression) getType() int {
	return TypeLiteral
}

func (ge GroupingExpression) getType() int {
	return TypeGrouping
}

func (ue UnaryExpression) getType() int {
	return TypeUnary
}

var current int
var tokens []lexer.Token

// Parse takes a series of tokens and returns an AST
func Parse(t []lexer.Token) Expression {
	current = 0
	tokens = t

	return expression()
}

// Set each rule of the grammar as a function
func expression() Expression {
	return equality()
}

func equality() Expression {
	expr := comparison()

	for match(lexer.TokenNotEqualTo, lexer.TokenEqualEqual) {
		operator := previous()
		right := comparison()
		expr = BinaryExpression{left: expr, operator: operator, right: right}
	}

	return expr
}

func comparison() Expression {
	expr := addition()

	for match(lexer.TokenGreaterThan, lexer.TokenGreaterEqual, lexer.TokenLessThan, lexer.TokenLessEqual) {
		operator := previous()
		right := addition()
		expr = BinaryExpression{left: expr, operator: operator, right: right}
	}

	return expr
}

func addition() Expression {
	expr := multiplication()

	for match(lexer.TokenMinus, lexer.TokenPlus) {
		operator := previous()
		right := multiplication()
		expr = BinaryExpression{left: expr, operator: operator, right: right}
	}

	return expr
}

func multiplication() Expression {
	expr := unary()

	for match(lexer.TokenMult, lexer.TokenDivision, lexer.TokenModulo) {
		operator := previous()
		right := unary()
		expr = BinaryExpression{left: expr, operator: operator, right: right}
	}

	return expr
}

func unary() Expression {
	if match(lexer.TokenNotEqualTo, lexer.TokenMinus) {
		operator := previous()
		right := unary()
		return UnaryExpression{operator: operator, right: right}
	}

	return primary()
}

func primary() Expression {
	if match(lexer.TokenFalse) {
		return LiteralExpression{value: false}
	}

	if match(lexer.TokenTrue) {
		return LiteralExpression{value: true}
	}

	if match(lexer.TokenNull) {
		return LiteralExpression{value: nil}
	}

	if match(lexer.TokenNumber, lexer.TokenString) {
		return LiteralExpression{value: previous().Literal}
	}

	if match(lexer.TokenLeftParentheses) {
		expr := expression()
		consume(lexer.TokenRightParenteses, "Se buscaba un ')' en la expresión")
		return GroupingExpression{expr}
	}

	panic("Unknown token")
}

func consume(tokenType int, message string) lexer.Token {
	if check(tokenType) {
		return advance()
	}

	token := peek()

	errorHandler.RaiseError(errorHandler.CodeSyntaxError, fmt.Sprintf("%d en %v: %v", token.Line, token.Lexeme, message), token.Line, token.Lexeme, true)
	panic("consume error")
}

func match(types ...int) bool {
	for _, v := range types {
		if check(v) {
			advance()
			return true
		}
	}
	return false
}

func check(tokenType int) bool {
	if isAtEnd() {
		return false
	}
	return peek().TokenType == tokenType
}

func advance() lexer.Token {
	if !isAtEnd() {
		current++
	}
	return previous()
}

func isAtEnd() bool {
	return peek().TokenType == lexer.TokenEOF
}

func peek() lexer.Token {
	return tokens[current]
}

func previous() lexer.Token {
	return tokens[current-1]
}
