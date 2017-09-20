package parser

import "clase-mates-computacionales/cazuela/lexer"

/*
Based on the formal grammar of Lox (http://www.craftinginterpreters.com/)

expression → literal
           | unary
           | binary
           | grouping ;

literal    → NUMBER | STRING | "true" | "false" | "nil" ;
grouping   → "(" expression ")" ;
unary      → ( "-" | "!" ) expression ;
binary     → expression operator expression ;
operator   → "==" | "!=" | "<" | "<=" | ">" | ">="
		   | "+"  | "-"  | "*" | "/" ;

*/

const (
	TypeBinary   = 0x10
	TypeLiteral  = 0x11
	TypeGrouping = 0x12
	TypeUnary    = 0x13
)

// Expression is the base interface for all expressions
type Expression interface {
	exprType int
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

func test() {
	binary := LiteralExpression{value: 32, exprType: TypeLiteral}
	binary.exprType = 2
}
