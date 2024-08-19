package main

import (
	"fmt"
	"strconv"
)

type TokenType int

const (
	TokenString TokenType = iota
	TokenNull
	TokenNumber
	TokenBoolean
	TokenSyntax
	EOF
)

var SyntaxChar = map[rune]struct{}{
	'{': {},
	'}': {},
	'[': {},
	']': {},
	':': {},
	',': {},
}

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

func (t *Token) GetTokenType(tokenType TokenType) string {
	switch tokenType {
	case TokenString:
		return "string"
	case TokenNull:
		return "null"
	case TokenNumber:
		return "number"
	case TokenBoolean:
		return "boolean"
	case TokenSyntax:
		return "json syntax"
	default:
		return "invalid"
	}
}

func(t *Token) GetTokenValue(token Token)(value any, err error){
	switch token.Type{
	case TokenString:
		value = token.Value
	case TokenNumber:
		value, err = strconv.ParseFloat(token.Value, 64)
	case TokenBoolean:
		value, err = strconv.ParseBool(token.Value)
	case TokenNull:
		value = nil
	default:
		err = t.GetTokenError(token)
	}

	if err != nil{
		return nil, err
	}

	return value, nil
}

func(t *Token) GetTokenError(token Token)error{
	return fmt.Errorf("unexpected %s token %s at line %d, column %d", t.GetTokenType(token.Type), t.Value, t.Line, t.Col)
}
