package main

import (
    "fmt"
    "testing"
)

func TestLexer(t *testing.T) {
    json := `{"key":null}`
    lexer := NewLexer(json)
    tokens, err := lexer.GetTokens()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(tokens) == 0 {
        t.Fatalf("Expected tokens, got none")
    }
    fmt.Println(tokens)
}

func TestLexerWithBoolean(t *testing.T) {
    json := `{"key":true}`
    lexer := NewLexer(json)
    tokens, err := lexer.GetTokens()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(tokens) == 0 {
        t.Fatalf("Expected tokens, got none")
    }
    fmt.Println(tokens)
}

func TestLexerWithNumber(t *testing.T) {
    json := `{"key":123}`
    lexer := NewLexer(json)
    tokens, err := lexer.GetTokens()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(tokens) == 0 {
        t.Fatalf("Expected tokens, got none")
    }
    fmt.Println(tokens)
}

func TestLexerWithString(t *testing.T) {
    json := `{"key":"value"}`
    lexer := NewLexer(json)
    tokens, err := lexer.GetTokens()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(tokens) == 0 {
        t.Fatalf("Expected tokens, got none")
    }
    fmt.Println(tokens)
}

func TestLexerWithArray(t *testing.T) {
    json := `{"key":[1, 2, 3]}`
    lexer := NewLexer(json)
    tokens, err := lexer.GetTokens()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(tokens) == 0 {
        t.Fatalf("Expected tokens, got none")
    }
    fmt.Println(tokens)
}

func TestLexerWithNestedObject(t *testing.T) {
    json := `{"key":{"nestedKey":"nestedValue"}}`
    lexer := NewLexer(json)
    tokens, err := lexer.GetTokens()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(tokens) == 0 {
        t.Fatalf("Expected tokens, got none")
    }
    fmt.Println(tokens)
}