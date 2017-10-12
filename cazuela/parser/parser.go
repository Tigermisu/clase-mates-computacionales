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
	GetType() int
}

// A BinaryExpression holds a expression with two other expressions and an operator in the middle
type BinaryExpression struct {
	Left     Expression
	Operator lexer.Token
	Right    Expression
}

// A LiteralExpression holds a simple literal
type LiteralExpression struct {
	Value interface{}
}

// A GroupingExpression holds more expressions inside it :D
type GroupingExpression struct {
	Expression Expression
}

// An UnaryExpression has only one operator and an expression to the right of it
type UnaryExpression struct {
	Operator lexer.Token
	Right    Expression
}

func (be BinaryExpression) GetType() int {
	return TypeBinary
}

func (le LiteralExpression) GetType() int {
	return TypeLiteral
}

func (ge GroupingExpression) GetType() int {
	return TypeGrouping
}

func (ue UnaryExpression) GetType() int {
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
		expr = BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func comparison() Expression {
	expr := addition()

	for match(lexer.TokenGreaterThan, lexer.TokenGreaterEqual, lexer.TokenLessThan, lexer.TokenLessEqual) {
		operator := previous()
		right := addition()
		expr = BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func addition() Expression {
	expr := multiplication()

	for match(lexer.TokenMinus, lexer.TokenPlus) {
		operator := previous()
		right := multiplication()
		expr = BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func multiplication() Expression {
	expr := unary()

	for match(lexer.TokenMult, lexer.TokenDivision, lexer.TokenModulo) {
		operator := previous()
		right := unary()
		expr = BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func unary() Expression {
	if match(lexer.TokenNotEqualTo, lexer.TokenMinus) {
		operator := previous()
		right := unary()
		return UnaryExpression{Operator: operator, Right: right}
	}

	return primary()
}

func primary() Expression {
	if match(lexer.TokenFalse) {
		return LiteralExpression{Value: false}
	}

	if match(lexer.TokenTrue) {
		return LiteralExpression{Value: true}
	}

	if match(lexer.TokenNull) {
		return LiteralExpression{Value: nil}
	}

	if match(lexer.TokenNumber, lexer.TokenString) {
		return LiteralExpression{Value: previous().Literal}
	}

	if match(lexer.TokenLeftParentheses) {
		expr := expression()
		consume(lexer.TokenRightParenteses, "Se buscaba un ')' en la expresión")
		return GroupingExpression{expr}
	}

	errorHandler.RaiseError(errorHandler.CodeSyntaxError, fmt.Sprintf("Elemento desconocido: %v", peek().Lexeme), peek().Line, "[Cocinado]", true)
	return nil
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
