package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"wabbit-go/llvm"
	"wabbit-go/parser"
)

func init() {
	log.SetLevel(log.DebugLevel)
	wd, _ := os.Getwd()
	log.Debugf("os.Getwd() %s", wd)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./tokenizer filename")
		os.Exit(1)
	}
	log.SetLevel(log.DebugLevel)
	filename := os.Args[1]
	prog, err := parser.HandleFile(filename)
	if err != nil {
		log.Errorf("wrong program %v", err)
	}
	err = os.WriteFile("out.ll", []byte(llvm.LLVM(prog)), 0644)
	if err != nil {
		log.Fatalf("Failed to write to out.ll: %v", err)
	}

	cmd := exec.Command("clang", "-O3", "runtime.c", "out.ll")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run clang: %v", err)
	}
	cmd = exec.Command("./a.out")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	//if err := cmd.Run(); err != nil {
	//	log.Fatalf("Failed to run ./a.out: %v", err)
	//}
}
