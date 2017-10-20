package lexer

import (
	"clase-mates-computacionales/cazuela/errorHandler"
	"clase-mates-computacionales/utilities"
	"fmt"
	"strconv"
)

// token types
const (
	// Math tokens
	TokenPlus          = 0x00
	TokenMinus         = 0x01
	TokenMult          = 0x02
	TokenDivision      = 0x03
	TokenModulo        = 0x04
	TokenExponentation = 0x05

	// Control Tokens
	TokenEqual           = 0x10
	TokenComma           = 0x11
	TokenLeftParentheses = 0x12
	TokenRightParenteses = 0x13
	TokenLeftBrace       = 0x14
	TokenRightBrace      = 0x15
	TokenEOF             = 0x16
	TokenSemiColon       = 0x17
	TokenNegation        = 0x18

	// Comparison Tokens
	TokenLessThan     = 0x20
	TokenGreaterThan  = 0x21
	TokenLessEqual    = 0x22
	TokenGreaterEqual = 0x23
	TokenEqualEqual   = 0x24
	TokenNotEqualTo   = 0x25

	// Literals
	TokenIdentifier = 0x40
	TokenString     = 0x41
	TokenNumber     = 0x42

	// Keywords
	TokenNull     = 0x80
	TokenLet      = 0x81
	TokenTrue     = 0x82
	TokenFalse    = 0x83
	TokenIf       = 0x84
	TokenElse     = 0x85
	TokenFunction = 0x86
	TokenFor      = 0x87
	TokenWhile    = 0x88
	TokenReturn   = 0x89
	TokenPrint    = 0x8A
	TokenAnd      = 0x8B
	TokenOr       = 0x8C
)

var keywords = map[string]int{
	"nulo":      TokenNull,
	"var":       TokenLet,
	"verdadero": TokenTrue,
	"falso":     TokenFalse,
	"si":        TokenIf,
	"nope":      TokenElse,
	"fn":        TokenFunction,
	"por":       TokenFor,
	"mientras":  TokenWhile,
	"sazonar":   TokenReturn,
	"servir":    TokenPrint,
	"y":         TokenAnd,
	"o":         TokenOr,
}

// A Token represents a token as interpreted by the lexer
type Token struct {
	TokenType int
	Lexeme    string
	Literal   interface{}
	Line      int
}

// Scan control and progress variables
var line, start, currentPosition, end int
var runes []rune
var rawCommand string
var tokens []Token

func (t Token) String() string {
	return fmt.Sprintf("(0x%X) %v - %v", t.TokenType, t.Lexeme, t.Literal)
}

// GetTokens takes a command string, and returns an array of all the tokens identified.
func GetTokens(command string) []Token {
	tokens = []Token{}
	line = 1
	start = 0
	currentPosition = 0
	end = len(command)
	runes = []rune(command)
	rawCommand = command

	for !atEndOfCommand() {
		start = currentPosition

		scanNextToken()
	}

	tokens = append(tokens, Token{TokenEOF, "~EOF~", nil, line})

	return tokens
}

func scanNextToken() {
	character := runes[currentPosition]
	currentPosition++

	switch character {
	case ';':
		addToken(TokenSemiColon)
		break
	case '+':
		addToken(TokenPlus)
		break
	case '-':
		addToken(TokenMinus)
		break
	case '*':
		addToken(TokenMult)
		break
	case '%':
		addToken(TokenModulo)
		break
	case '^':
		addToken(TokenExponentation)
		break
	case '(':
		addToken(TokenLeftParentheses)
		break
	case ')':
		addToken(TokenRightParenteses)
		break
	case '{':
		addToken(TokenLeftBrace)
		break
	case '}':
		addToken(TokenRightBrace)
		break
	case ',':
		addToken(TokenComma)
		break
	case '!':
		addTokenIfMatch('=', TokenNotEqualTo, TokenNegation)
		break
	case '=':
		addTokenIfMatch('=', TokenEqualEqual, TokenEqual)
		break
	case '<':
		addTokenIfMatch('=', TokenLessEqual, TokenLessThan)
		break
	case '>':
		addTokenIfMatch('=', TokenGreaterEqual, TokenGreaterThan)
		break
	case '/':
		if runes[currentPosition] == '/' { // This is a comment
			for peek() != '\n' && !atEndOfCommand() {
				currentPosition++
			}
		} else {
			addToken(TokenDivision)
		}
		break
	case '"':
		parseStringLexeme()
		break
	case ' ':
	case '\r':
	case '\t':
		break // eat whitespace
	case '\n':
		line++ // Mark line increase
		break
	default:
		if isDigit(character) {
			parseNumberLexeme()
		} else if isAlpha(character) {
			parseIdentifier()
		} else {
			errorHandler.RaiseError(errorHandler.CodeSyntaxError, fmt.Sprintf("Caracter desconocido: %c", character), line, "[Preparado]", true)
		}

	}
}

func addToken(TokenType int) {
	addTokenWithLiteral(TokenType, nil)
}

func addTokenIfMatch(m rune, tokenIfMatch int, tokenElse int) {
	if match(m) {
		addToken(tokenIfMatch)
	} else {
		addToken(tokenElse)
	}
}

func addTokenWithLiteral(TokenType int, Literal interface{}) {
	Lexeme := rawCommand[start:currentPosition]
	tokens = append(tokens, Token{TokenType, Lexeme, Literal, line})
}

func match(m rune) bool {
	if atEndOfCommand() || runes[currentPosition] != m {
		return false
	}
	currentPosition++
	return true
}

func peek() rune {
	if atEndOfCommand() {
		return 3 // 3 == EOF
	}
	return runes[currentPosition]
}

func peekNext() rune {
	if currentPosition+1 >= end {
		return 3
	}
	return runes[currentPosition+1]
}

func atEndOfCommand() bool {
	return currentPosition >= end
}

func isDigit(character rune) bool {
	return character >= '0' && character <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(character rune) bool {
	return isDigit(character) || isAlpha(character)
}

func parseNumberLexeme() {
	for isDigit(peek()) {
		currentPosition++
	}

	if peek() == '.' && isDigit(peekNext()) {
		currentPosition++

		for isDigit(peek()) {
			currentPosition++
		}
	}

	Literal, err := strconv.ParseFloat(rawCommand[start:currentPosition], 64)
	utilities.AssertError(err)

	addTokenWithLiteral(TokenNumber, Literal)
}

func parseIdentifier() {
	for isAlphaNumeric(peek()) {
		currentPosition++
	}

	possibleKeyword := keywords[rawCommand[start:currentPosition]]

	if possibleKeyword != 0 { // is 0 when the value is not in the map
		addToken(possibleKeyword)
	} else {
		addToken(TokenIdentifier)
	}
}

func parseStringLexeme() {
	for peek() != '"' && !atEndOfCommand() {
		if peek() == '\n' {
			line++
		}
		currentPosition++
	}

	if atEndOfCommand() {
		errorHandler.RaiseError(errorHandler.CodeUnexpectedEOF,
			"Se esperaba terminar una cadena, pero el archivo se acabó.", line, "", true)
	} else {
		currentPosition++

		Literal := rawCommand[start+1 : currentPosition-1]
		addTokenWithLiteral(TokenString, Literal)
	}
}
