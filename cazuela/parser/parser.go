package parser

import (
	"clase-mates-computacionales/cazuela/errorHandler"
	"clase-mates-computacionales/cazuela/lexer"
	"fmt"
)

/*
Based on the formal grammar of Lox (http://www.craftinginterpreters.com/)

program     → declaration* eof ;

declaration → varDecl
			  | fnDecl
              | stmt ;

varDecl		→ "var" IDENTIFIER ( "=" expression )? ";" ;

fnDecl  	→ "fn" function ;
function 	→ IDENTIFIER "(" parameters? ")" block ;
parameters → IDENTIFIER ( "," IDENTIFIER )* ;

stmt   		   → statement
				| print
				| block
				| while
				| if
				| forLoop
				| ReturnStmt

forLoop		  → "por" "(" ( varDecl | exprStmt | ";" )
				expression? ";"
				expression? ")" statement ;

ReturnStmt 	  → "sazonar" expression? ";" ;

if			   → "si" "(" expression ")" statement ( "nope" statement )? ;
block     	   → "{" declaration* "}" ;
while 	       → "mientras" "(" expression ")" statement ;

expression     → assignment ;

assigment      → identifier "=" assignment\
				| logic_or ;

logic_or  	   → logic_and ( "o" logic_and )* ;
logic_and  	   → equality ( "y" equality )* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → addition ( ( ">" | ">=" | "<" | "<=" ) addition )* ;
addition       → multiplication ( ( "-" | "+" ) multiplication )* ;
multiplication → exponentiation ( ( "/" | "*" | "%" ) exponentiation )* ;
exponentiation → unary ( ( "^" ) unary )* ;
unary          → ( "!" | "-" ) unary ;
				 | call ;
call		   → primary ( "(" arguments? ")" )* ;
arguments 	   → expression ( "," expression )* ;
primary        → "verdadero" | "falso" | "nulo"
				 | NUMBER | STRING
				 | "(" expression ")"
				 | IDENTIFIER ;

*/

const (
	TypeBinary     = 0x10
	TypeLiteral    = 0x11
	TypeGrouping   = 0x12
	TypeUnary      = 0x13
	TypeVariable   = 0x14
	TypeAssignment = 0x15
	TypeLogical    = 0x16
	TypeWhile      = 0x17
	TypeCall       = 0x18
	TypeFn         = 0x19
	TypeReturn     = 0x1A

	TypeStatement   = 0x20
	TypePrint       = 0x21
	TypeDeclaration = 0x22
	TypeBlock       = 0x23
	TypeIf          = 0x24
)

type Stmt interface {
	GetStmtType() int
}

type Statement struct {
	Expr Expression
}

type Print struct {
	Expr Expression
}

type Declaration struct {
	Name        lexer.Token
	Initializer Expression
}

type Block struct {
	Statements []Stmt
}

type If struct {
	Condition  Expression
	ThenBranch Stmt
	ElseBranch Stmt
}

type While struct {
	Condition Expression
	Body      Stmt
}

type FnDecl struct {
	Name       lexer.Token
	Parameters []lexer.Token
	Body       []Stmt
}

type ReturnStmt struct {
	Keyword lexer.Token
	Value   Expression
}

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

type VariableExpression struct {
	Name lexer.Token
}

type AssignmentExpression struct {
	Name  lexer.Token
	Value Expression
}

type LogicalExpression struct {
	Left     Expression
	Operator lexer.Token
	Right    Expression
}

type CallExpression struct {
	Callee            Expression
	ClosingParenteses lexer.Token
	Arguments         []Expression
}

func (st Statement) GetStmtType() int {
	return TypeStatement
}

func (st Print) GetStmtType() int {
	return TypePrint
}

func (st Declaration) GetStmtType() int {
	return TypeDeclaration
}

func (st Block) GetStmtType() int {
	return TypeBlock
}

func (st If) GetStmtType() int {
	return TypeIf
}

func (st While) GetStmtType() int {
	return TypeWhile
}

func (st FnDecl) GetStmtType() int {
	return TypeFn
}

func (st ReturnStmt) GetStmtType() int {
	return TypeReturn
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

func (ve VariableExpression) GetType() int {
	return TypeVariable
}

func (ae AssignmentExpression) GetType() int {
	return TypeAssignment
}

func (lo LogicalExpression) GetType() int {
	return TypeLogical
}

func (ca CallExpression) GetType() int {
	return TypeCall
}

var current int
var tokens []lexer.Token

// Parse takes a series of tokens and returns an AST
func Parse(t []lexer.Token) []Stmt {
	current = 0
	tokens = t

	statements := make([]Stmt, 1)

	for !isAtEnd() {
		statements = append(statements, declaration())
	}

	return statements
}

func declaration() Stmt {
	defer func() {
		if err := recover(); err != nil {
			errorHandler.RaiseError(errorHandler.CodeRuntimeError, "Error en cocinado", 0, "[Cocinado]", true)
		}
	}()

	if match(lexer.TokenLet) {
		return varDeclaration()
	}

	if match(lexer.TokenFunction) {
		return function()
	}

	return statement()
}

func function() FnDecl {
	name := consume(lexer.TokenIdentifier, "Se esperaba un nombre de función.")
	consume(lexer.TokenLeftParentheses, "Se esperaba un ( después del nombre de la función.")
	parameters := make([]lexer.Token, 1)

	if !check(lexer.TokenRightParenteses) {
		parameters = append(parameters, consume(lexer.TokenIdentifier, "Se esperaba nombre del parámetro"))
		for match(lexer.TokenComma) {
			parameters = append(parameters, consume(lexer.TokenIdentifier, "Se esperaba nombre del parámetro"))
		}
	}
	consume(lexer.TokenRightParenteses, "Se esperaba un ) después de los parámetros de la función")

	consume(lexer.TokenLeftBrace, "Se esperaba un { al empezar la función.")

	body := block()

	return FnDecl{name, parameters, body}
}

func varDeclaration() Stmt {
	name := consume(lexer.TokenIdentifier, "Se esperaba un nombre de variable.")

	var initializer Expression

	if match(lexer.TokenEqual) {
		initializer = expression()
	}

	consume(lexer.TokenSemiColon, "Se esperaba un ; después de la declaración de la variable.")

	return Declaration{Name: name, Initializer: initializer}

}

func statement() Stmt {
	if match(lexer.TokenIf) {
		return ifStatement()
	}

	if match(lexer.TokenLeftBrace) {
		return Block{block()}
	}

	if match(lexer.TokenPrint) {
		return printStatement()
	}

	if match(lexer.TokenWhile) {
		return whileStatement()
	}

	if match(lexer.TokenFor) {
		return forStatement()
	}

	if match(lexer.TokenReturn) {
		return returnStatement()
	}

	return expressionStatement()
}

func block() []Stmt {
	statements := make([]Stmt, 1)

	for !check(lexer.TokenRightBrace) && !isAtEnd() {
		statements = append(statements, declaration())
	}

	consume(lexer.TokenRightBrace, "Se esperaba un } al final del bloque.")
	return statements
}

func returnStatement() Stmt {
	keyword := previous()

	var value Expression

	if !check(lexer.TokenSemiColon) {
		value = expression()
	}

	consume(lexer.TokenSemiColon, "Se esperaba un ; al final del sazonado")

	return ReturnStmt{keyword, value}
}

func printStatement() Stmt {
	value := expression()

	consume(lexer.TokenSemiColon, "Se buscaba un ; al final.")

	return Print{value}
}

// Caramelizer for whiles
func forStatement() Stmt {
	consume(lexer.TokenLeftParentheses, "Se esperaba un ( después de 'por'.")
	var initializer Stmt
	if match(lexer.TokenLet) {
		initializer = varDeclaration()
	} else if !match(lexer.TokenSemiColon) {
		initializer = expressionStatement()
	}

	var condition Expression
	if !check(lexer.TokenSemiColon) {
		condition = expression()
	}
	consume(lexer.TokenSemiColon, "Se esperaba un ; después de la condición.")

	var increment Expression
	if !check(lexer.TokenRightParenteses) {
		increment = expression()
	}

	consume(lexer.TokenRightParenteses, "Se esperaba un ) después del 'por'.")

	body := statement()

	if increment != nil {
		body = Block{[]Stmt{body, Statement{increment}}}
	}

	if condition == nil {
		condition = LiteralExpression{true}
	}

	body = While{condition, body}

	if initializer != nil {
		body = Block{[]Stmt{initializer, body}}
	}

	return body
}

func expressionStatement() Stmt {
	expr := expression()

	consume(lexer.TokenSemiColon, "Se buscaba un ; al final.")

	return Statement{expr}
}

func ifStatement() Stmt {
	consume(lexer.TokenLeftParentheses, "Se esperaba un ( en la condición.")
	condition := expression()
	consume(lexer.TokenRightParenteses, "Se esperaba un ) al final de la condición.")

	thenBranch := statement()
	var elseBranch Stmt

	if match(lexer.TokenElse) {
		elseBranch = statement()
	}

	return If{condition, thenBranch, elseBranch}
}

func whileStatement() Stmt {
	consume(lexer.TokenLeftParentheses, "Se esperaba un ( en la condición.")
	condition := expression()
	consume(lexer.TokenRightParenteses, "Se esperaba un ) al final de la condición.")

	body := statement()

	return While{condition, body}
}

func expression() Expression {
	return assignment()
}

func assignment() Expression {
	expr := or()

	if match(lexer.TokenEqual) {
		equals := previous()
		value := assignment()

		if v, ok := expr.(VariableExpression); ok {
			name := v.Name
			return AssignmentExpression{Name: name, Value: value}
		}

		errorHandler.RaiseError(errorHandler.CodeSyntaxError, "Lado izquierdo de asignación inválido.", equals.Line, "Cocinado", true)
	}

	return expr
}

func or() Expression {
	expr := and()

	for match(lexer.TokenOr) {
		operator := previous()
		right := and()
		expr = LogicalExpression{expr, operator, right}
	}

	return expr
}

func and() Expression {
	expr := equality()

	for match(lexer.TokenAnd) {
		operator := previous()
		right := equality()
		expr = LogicalExpression{expr, operator, right}
	}

	return expr
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
	expr := exponentiation()

	for match(lexer.TokenMult, lexer.TokenDivision, lexer.TokenModulo) {
		operator := previous()
		right := exponentiation()
		expr = BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func exponentiation() Expression {
	expr := unary()

	for match(lexer.TokenExponentation) {
		operator := previous()
		right := unary()
		expr = BinaryExpression{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func unary() Expression {
	if match(lexer.TokenNegation, lexer.TokenMinus) {
		operator := previous()
		right := unary()
		return UnaryExpression{Operator: operator, Right: right}
	}

	return call()
}

func call() Expression {
	expr := primary()

	for {
		if match(lexer.TokenLeftParentheses) {
			expr = finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func finishCall(callee Expression) Expression {
	arguments := make([]Expression, 1)

	if !check(lexer.TokenRightParenteses) {
		arguments = append(arguments, expression())
		for match(lexer.TokenComma) {
			arguments = append(arguments, expression())
		}
	}

	paren := consume(lexer.TokenRightParenteses, "Se esperaba un ) al final de la llamada.")

	return CallExpression{callee, paren, arguments}
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

	if match(lexer.TokenIdentifier) {
		return VariableExpression{previous()}
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

	errorHandler.RaiseError(errorHandler.CodeSyntaxError, message, token.Line, token.Lexeme, true)
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
