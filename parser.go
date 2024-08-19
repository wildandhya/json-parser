package main

import (
	"fmt"
	"strings"
)

type Parser struct {
	tokens []Token
	json map[string]interface{}
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, json: make(map[string]interface{})}
}

func (p *Parser) Parse() (string, error) {
	if len(p.tokens) == 0 {
		return "", fmt.Errorf("no tokens to parse")
	}

	token := p.tokens[0]
	if token.Value == "{" {
		parsedData, err := p.ParseObject()
		if err != nil {
			return "", err
		}
		return p.toJSONString(parsedData), nil
	} else if token.Value == "[" {
		parsedData, err := p.ParseArray()
		if err != nil {
			return "", err
		}
		return p.toJSONString(parsedData), nil
	}

	return "", fmt.Errorf("unexpected token '%s' at line %d, col %d", token.Value, token.Line, token.Col)
}

func (p *Parser) ParseObject() (interface{}, error) {
	p.tokens = p.tokens[1:] // Skip the opening '{'
	keys := make(map[string]struct{})

	for len(p.tokens) > 0 {
		token := p.tokens[0]

		// Check if we've reached the end of the object
		if token.Type == TokenSyntax && token.Value == "}" {
			p.tokens = p.tokens[1:] // Skip the closing '}'
			return p.json, nil
		}

		// Expect a string token for the key
		if token.Type != TokenString {
			return p.json, fmt.Errorf("expected object key at line %d, col %d, but got '%s'", token.Line, token.Col, token.Value)
		}

		key := token.Value

		// Check for duplicate keys
		if _, exists := keys[key]; exists {
			fmt.Printf("warning: duplicate object key '%s' at line %d, col %d\n", key, token.Line, token.Col)
		}
		keys[key] = struct{}{}

		// Move to the next token after the key
		p.tokens = p.tokens[1:]

		// Expect a colon token after the key
		if len(p.tokens) == 0 || p.tokens[0].Type != TokenSyntax || p.tokens[0].Value != ":" {
			return p.json, fmt.Errorf("expected ':' after object key '%s' at line %d, col %d", key, token.Line, token.Col)
		}

		// Move to the next token after the colon
		p.tokens = p.tokens[1:]

		// Parse the value associated with the key
		value, err := p.ParseValue()
		if err != nil {
			return p.json, err
		}

		// Store the parsed key-value pair in the JSON map
		p.json[key] = value

		// After parsing the value, expect either a comma or the end of the object
		if len(p.tokens) == 0 {
			return p.json, fmt.Errorf("unexpected end of input after key '%s'", key)
		}

		token = p.tokens[0]
		if token.Type == TokenSyntax && token.Value == "," {
			p.tokens = p.tokens[1:] // Skip the comma and continue
		} else if token.Type == TokenSyntax && token.Value == "}" {
			// Reached the end of the object
			p.tokens = p.tokens[1:]
			return p.json, nil
		} else {
			return p.json, fmt.Errorf("expected ',' or '}' after value at line %d, col %d", token.Line, token.Col)
		}
	}

	return p.json, fmt.Errorf("unexpected end of input while parsing object")
}

func (p *Parser) ParseValue() (interface{}, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("unexpected end of input while parsing value")
	}

	token := p.tokens[0]
	switch token.Type {
	case TokenString:
		p.tokens = p.tokens[1:] 
		return token.Value, nil
	case TokenNumber:
		p.tokens = p.tokens[1:] 
		return token.Value, nil
	case TokenBoolean:
		p.tokens = p.tokens[1:] // Consume the token
		// Convert string to a boolean (true/false)
		return token.Value == "true", nil
	case TokenSyntax:
		if token.Value == "{" {
			return p.ParseObject()
		} else if token.Value == "[" {
			return p.ParseArray()
		} else {
			return nil, fmt.Errorf("unexpected syntax token '%s' at line %d, col %d", token.Value, token.Line, token.Col)
		}
	default:
		return nil, fmt.Errorf("unexpected token '%s' at line %d, col %d", token.Value, token.Line, token.Col)
	}
}

func (p *Parser) ParseArray() (interface{}, error) {
	p.tokens = p.tokens[1:] // Skip the opening '['
	var array []interface{}

	for len(p.tokens) > 0 {
		token := p.tokens[0]

		// Check if we've reached the end of the array
		if token.Type == TokenSyntax && token.Value == "]" {
			p.tokens = p.tokens[1:] // Skip the closing ']'
			return array, nil
		}

		// Parse the next value in the array
		value, err := p.ParseValue()
		if err != nil {
			return nil, err
		}
		array = append(array, value)

		// After parsing the value, expect either a comma or the end of the array
		if len(p.tokens) == 0 {
			return nil, fmt.Errorf("unexpected end of input while parsing array")
		}

		token = p.tokens[0]
		if token.Type == TokenSyntax && token.Value == "," {
			p.tokens = p.tokens[1:] // Skip the comma and continue
		} else if token.Type == TokenSyntax && token.Value == "]" {
			// Reached the end of the array
			p.tokens = p.tokens[1:]
			return array, nil
		} else {
			return nil, fmt.Errorf("expected ',' or ']' after array value at line %d, col %d", token.Line, token.Col)
		}
	}

	return nil, fmt.Errorf("unexpected end of input while parsing array")
}

func (p *Parser) toJSONString(data interface{}) string {
	switch v := data.(type) {
	case map[string]interface{}:
		var sb strings.Builder
		sb.WriteString("{")
		for k, val := range v {
			if sb.Len() > 1 {
				sb.WriteString(",")
			}
			sb.WriteString(fmt.Sprintf("\"%s\":%s", k, p.toJSONString(val)))
		}
		sb.WriteString("}")
		return sb.String()
	case []interface{}:
		var sb strings.Builder
		sb.WriteString("[")
		for i, val := range v {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(p.toJSONString(val))
		}
		sb.WriteString("]")
		return sb.String()
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return v.(string)
	}
}








