package tests

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
	"wabbit-go/tokenize" // Update this import path
)

var rightProgramPath string
var wrongProgramPath string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rightProgramPath = filepath.Join(wd, "Programs")
	wrongProgramPath = filepath.Join(wd, "ErrorLex")
}

func TestRightProgram(t *testing.T) {
	t.Log(rightProgramPath + "/*.wb")
	rightFiles, _ := filepath.Glob(rightProgramPath + "/*.wb")
	//t.Log("files:", rightFiles)
	for _, rightFile := range rightFiles {
		t.Log("handing ", rightFile)
		tokens, err := tokenize.HandleFile(rightFile) // Make sure Main function exists in your tokenize package

		if err != nil {
			t.Log(tokens)
			t.Errorf(err.Error())
		}
	}
}

func TestErrorProgram(t *testing.T) {
	wrongFiles, _ := filepath.Glob(wrongProgramPath + "/*.wb")
	for _, file := range wrongFiles {
		t.Log("handing ", file)
		_, err := tokenize.HandleFile(file)
		if err == nil {
			t.Errorf("should failed but not failed")
		}
	}
}

func TestSymbols(t *testing.T) {
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

func TestNumbers(t *testing.T) {
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

func TestKeywords(t *testing.T) {
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
