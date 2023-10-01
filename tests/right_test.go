package tests

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"wabbit-go/interpreter"
	"wabbit-go/llvm"
	"wabbit-go/parser" // Update this import path
	"wabbit-go/wasm"
	"wabbit-go/wvm"
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
		p, err := parser.HandleFile(rightFile) // Make sure Main function exists in your tokenize package

		if err != nil {
			t.Errorf(err.Error())
		}
		if !strings.HasSuffix(rightFile, "25_tailrecurisve.wb") {
			// can't handle the tail deep recursive
			interpreter.InterpretProgram(p)
		}

		wvm.Wvm(p)
		wasm.Wasm(p)
		llvm.LLVM(p)
	}
}
