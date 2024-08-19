package main
import (
    "fmt"
    "testing"
)

func TestParser(t *testing.T) {
    tokens := []Token{
        {TokenSyntax, "{", 1, 1},
        {TokenString, "name", 2, 2},
        {TokenSyntax, ":", 2, 7},
        {TokenString, "John", 2, 9},
        {TokenSyntax, ",", 2, 14},
        {TokenString, "age", 3, 2},
        {TokenSyntax, ":", 3, 6},
        {TokenNumber, "30", 3, 8},
        {TokenSyntax, "}", 4, 1},
    }
    parser := NewParser(tokens)
    data, err := parser.Parse()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    expected := map[string]interface{}{
        "name": "John",
        "age":  30,
    }
    if fmt.Sprintf("%v", data) != fmt.Sprintf("%v", expected) {
        t.Errorf("Expected %v, got %v", expected, data)
    }
	fmt.Println(data)
}

func TestParserEmptyObject(t *testing.T) {
    tokens := []Token{
        {TokenSyntax, "{", 1, 1},
        {TokenSyntax, "}", 1, 2},
    }
    parser := NewParser(tokens)
    data, err := parser.Parse()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    expected := map[string]interface{}{}
    if fmt.Sprintf("%v", data) != fmt.Sprintf("%v", expected) {
        t.Errorf("Expected %v, got %v", expected, data)
    }
}

func TestParserArray(t *testing.T) {
    tokens := []Token{
        {TokenSyntax, "[", 1, 1},
        {TokenNumber, "1", 1, 2},
        {TokenSyntax, ",", 1, 3},
        {TokenNumber, "2", 1, 4},
        {TokenSyntax, ",", 1, 5},
        {TokenNumber, "3", 1, 6},
        {TokenSyntax, "]", 1, 7},
    }
    parser := NewParser(tokens)
    data, err := parser.Parse()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    expected := []interface{}{1, 2, 3}
    if fmt.Sprintf("%v", data) != fmt.Sprintf("%v", expected) {
        t.Errorf("Expected %v, got %v", expected, data)
    }
	fmt.Println(data)
}

func TestParserArrayOfObject(t *testing.T) {
    tokens := []Token{
        {TokenSyntax, "[", 1, 1},
        {TokenSyntax, "{", 1, 1},
        {TokenString, "name", 2, 2},
        {TokenSyntax, ":", 2, 7},
        {TokenString, "John", 2, 9},
        {TokenSyntax, ",", 2, 14},
        {TokenString, "age", 3, 2},
        {TokenSyntax, ":", 3, 6},
        {TokenNumber, "30", 3, 8},
        {TokenSyntax, "}", 4, 1},
        {TokenSyntax, "]", 1, 7},
    }
    parser := NewParser(tokens)
    data, err := parser.Parse()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
	expected := []map[string]interface{}{
		{  "name": "John","age":  30},
    }
    if fmt.Sprintf("%v", data) != fmt.Sprintf("%v", expected) {
        t.Errorf("Expected %v, got %v", expected, data)
    }
	fmt.Println(data)
}

func TestParserNestedObject(t *testing.T) {
    tokens := []Token{
        {TokenSyntax, "{", 1, 1},
        {TokenString, "person", 2, 2},
        {TokenSyntax, ":", 2, 9},
        {TokenSyntax, "{", 2, 11},
        {TokenString, "name", 3, 2},
        {TokenSyntax, ":", 3, 7},
        {TokenString, "John", 3, 9},
        {TokenSyntax, "}", 4, 1},
        {TokenSyntax, "}", 5, 1},
    }
    parser := NewParser(tokens)
    data, err := parser.Parse()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    expected := map[string]interface{}{
        "person": map[string]interface{}{
            "name": "John",
        },
    }
    if fmt.Sprintf("%v", data) != fmt.Sprintf("%v", expected) {
        t.Errorf("Expected %v, got %v", expected, data)
    }
	fmt.Println(data)
}

func TestParserInvalidSyntax(t *testing.T) {
    tokens := []Token{
        {TokenSyntax, "{", 1, 1},
        {TokenString, "name", 2, 2},
        {TokenSyntax, ":", 2, 7},
        {TokenString, "John", 2, 9},
        {TokenSyntax, ",", 2, 14},
        {TokenString, "age", 3, 2},
        {TokenSyntax, ":", 3, 6},
        {TokenNumber, "30", 3, 8},
        // Missing closing brace
    }
    parser := NewParser(tokens)
    _, err := parser.Parse()
    if err == nil {
        t.Fatalf("Expected an error due to invalid syntax, got nil")
    }
	
}