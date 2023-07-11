package tests

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"wabbit-go/tokenize" // Update this import path
)

var rightProgramPath = filepath.Join(filepath.Dir(os.Args[0]), "Programs")
var wrongProgramPath = filepath.Join(filepath.Dir(os.Args[0]), "ErrorLex")

func testRightProgram(t *testing.T) {
	rightFiles, _ := filepath.Glob(rightProgramPath + "/*.wb")

	for _, rightFile := range rightFiles {
		_, err := tokenize.HandleFile(rightFile) // Make sure Main function exists in your tokenize package
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func testErrorProgram(t *testing.T) {
	wrongFiles, _ := filepath.Glob(wrongProgramPath + "/*.wb")
	for _, file := range wrongFiles {
		_, err := tokenize.HandleFile(file)
		if err == nil {
			t.Errorf("should failed but not failed")
		}
	}
}

func testSymbols(t *testing.T) {
	tokens, err := tokenize.Tokenize("+ - * / < > <= >= == != = && || , ; ( ) { } !") // Make sure Tokenize function exists in your tokenize package
	if err != nil {
		t.Errorf(err.Error())
	}

	tokTypes := make([]string, len(tokens))
	for i, tok := range tokens {
		tokTypes[i] = tok.Type
	}

	expected := []string{"PLUS", "MINUS", "TIMES", "DIVIDE", "LT", "GT", "LE", "GE", "EQ", "NE", "ASSIGN", "LAND", "LOR", "COMMA", "SEMI", "LPAREN", "RPAREN", "LBRACE", "RBRACE", "LNOT", "EOF"}
	assert.Equal(t, expected, tokTypes)
}

func testNumbers(t *testing.T) {
	tokens, err := tokenize.Tokenize("123 123.45") // Make sure Tokenize function exists in your tokenize package
	if err != nil {
		t.Errorf(err.Error())
	}

	tokTypes := make([]string, len(tokens))
	tokValues := make([]string, len(tokens))
	for i, tok := range tokens {
		tokTypes[i] = tok.Type
		tokValues[i] = tok.Value
	}

	expectedTypes := []string{"INTEGER", "FLOAT", "EOF"}
	expectedValues := []string{"123", "123.45", "EOF"}
	assert.Equal(t, expectedTypes, tokTypes)
	assert.Equal(t, expectedValues, tokValues)
}

func testKeywords(t *testing.T) {
	tokens, err := tokenize.Tokenize("if else while var const break continue print func return true false") // Make sure Tokenize function exists in your tokenize package
	if err != nil {
		t.Errorf(err.Error())
	}

	tokTypes := make([]string, len(tokens))
	for i, tok := range tokens {
		tokTypes[i] = tok.Type
	}

	expected := []string{"IF", "ELSE", "WHILE", "VAR", "CONST", "BREAK", "CONTINUE", "PRINT", "FUNC", "RETURN", "TRUE", "FALSE", "EOF"}
	assert.Equal(t, expected, tokTypes)
}
