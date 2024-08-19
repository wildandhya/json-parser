package main

import (
	"fmt"
	"unicode"
)

type Lexer struct {
	input []rune
	line  int
	col   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: []rune(input), line: 1, col: 0}
}

func (l *Lexer) GetTokens() ([]Token, error) {
	var tokens []Token

	for len(l.input) > 0 {
		char := l.input[0]

		if unicode.IsSpace(char) {
			l.advanceCursor()
			continue
		}

		token, err := l.getToken(char)
		if err != nil {
			return tokens, err
		} else if token != (Token{}) {
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

func (l *Lexer) getToken(char rune) (Token, error) {
	token, err := l.GetTokenString(char)
	if err == nil && token != (Token{}) {
		return token, nil
	}

	token, err = l.GetTokenNumber(char)
	if err == nil && token != (Token{}) {
		return token, nil
	}

	token, err = l.GetTokenSyntax(char)
	if err == nil && token != (Token{}) {
		return token, nil
	}

	token, err = l.GetTokenBoolean(char)
	if err == nil && token != (Token{}) {
		return token, nil
	}

	token, err = l.GetTokenNull(char)
	if err == nil && token != (Token{}) {
		return token, nil
	}

	return Token{}, l.PrintError(char)
}

func (l *Lexer) GetTokenString(char rune) (Token, error) {
	if char != '"' {
		return Token{}, nil
	}
	l.advanceCursor()

	var value []rune
	escaped := false

	for i, char := range l.input {
		l.col++
		if escaped {
			if validEscape(char) {
				escaped = false
			} else {
				return Token{}, fmt.Errorf("invalid escape character '\\%s' at line %d, col %d", string(char), l.line, l.col)
			}
		} else if char == '\\' {
			escaped = true
		} else if char == '"' {
			l.input = l.input[i+1:]
			return Token{TokenString, string(value), l.line, l.col}, nil
		} else {
			value = append(value, char)
		}
	}

	return Token{}, fmt.Errorf("unterminated string at line %d, col %d", l.line, l.col)
}

func validEscape(char rune) bool {
	switch char {
	case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
		return true
	default:
		return false
	}
}

func (l *Lexer) GetTokenSyntax(char rune) (Token, error) {
	if _, ok := SyntaxChar[char]; ok {
		l.advanceCursor()
		return Token{TokenSyntax, string(char), l.line, l.col}, nil
	}
	return Token{}, nil
}

func (l *Lexer) GetTokenNull(char rune) (Token, error) {
	var value string
	if char == 'n' && len(l.input) >= 4 && string(l.input[:4]) == "null" {
		value = "null"
	} else {
		return Token{}, fmt.Errorf("invalid null value starting with '%s' at line %d, col %d", string(char), l.line, l.col)
	}

	l.input = l.input[len(value):]
	l.col += len(value)

	return Token{TokenNull, value, l.line, l.col}, nil
}

func (l *Lexer) GetTokenBoolean(char rune) (Token, error) {
	var value string
	if char == 't' && len(l.input) >= 4 && string(l.input[:4]) == "true" {
		value = "true"
	} else if char == 'f' && len(l.input) >= 5 && string(l.input[:5]) == "false" {
		value = "false"
	} else {
		return Token{}, fmt.Errorf("invalid boolean value starting with '%s' at line %d, col %d", string(char), l.line, l.col)
	}

	l.input = l.input[len(value):]
	l.col += len(value)

	return Token{TokenBoolean, value, l.line, l.col}, nil
}

func (l *Lexer) GetTokenNumber(char rune) (Token, error) {
	if char == '-' || unicode.IsDigit(char) {
		startCol := l.col
		var value []rune

		if char == '-' {
			value = append(value, char)
			l.advanceCursor()
		}

		hasDecimalPoint := false

		for len(l.input) > 0 {
			char := l.input[0]
			if unicode.IsDigit(char) {
				value = append(value, char)
				l.advanceCursor()
			} else if char == '.' && !hasDecimalPoint {
				// Handle the decimal point for float numbers
				hasDecimalPoint = true
				value = append(value, char)
				l.advanceCursor()
			} else {
				break
			}
		}

		// Ensure the token is not just a sign or a lone decimal point
		if len(value) == 1 && (value[0] == '-' || value[0] == '+' || value[0] == '.') {
			return Token{}, fmt.Errorf("invalid number '%s' at line %d, col %d", string(value), l.line, startCol)
		}

		return Token{TokenNumber, string(value), l.line, startCol}, nil
	}
	return Token{}, nil
}

func (l *Lexer) advanceCursor() {
	if len(l.input) > 0 {
		if l.input[0] == '\n' {
			l.line++
			l.col = 0
		} else {
			l.col++
		}
		l.input = l.input[1:]
	}
}

func (l *Lexer) PrintError(char rune) error {
	return fmt.Errorf("unexpected character '%s' at line %d, column %d", string(char), l.line, l.col)
}
