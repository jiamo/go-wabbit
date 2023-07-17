package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"wabbit-go/parser"
	"wabbit-go/wvm"
)

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
	wvm.Wvm(prog)
}
