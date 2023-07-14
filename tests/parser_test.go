package tests

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"wabbit-go/interpreter"
	"wabbit-go/parser" // Update this import path
)

//var rightProgramPath string
//var wrongProgramPath string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rightProgramPath = filepath.Join(wd, "Programs")
	wrongProgramPath = filepath.Join(wd, "ErrorLex")
}

func TestParserRightProgram(t *testing.T) {
	t.Log(rightProgramPath + "/*.wb")
	rightFiles, _ := filepath.Glob(rightProgramPath + "/*.wb")
	//t.Log("files:", rightFiles)
	for _, rightFile := range rightFiles {
		t.Log("handing ", rightFile)
		p, err := parser.HandleFile(rightFile) // Make sure Main function exists in your tokenize package

		if err != nil {
			t.Errorf(err.Error())
		}
		interpreter.InterpretProgram(p.Model)
	}
}

// TODO wrong program should report error
//func TestParseErrorProgram(t *testing.T) {
//	wrongFiles, _ := filepath.Glob(wrongProgramPath + "/*.wb")
//	for _, file := range wrongFiles {
//		t.Log("handing ", file)
//		p, err := tokenize.HandleFile(file)
//		if err == nil {
//			t.Errorf("should failed but not failed")
//		}
//
//	}
//}
