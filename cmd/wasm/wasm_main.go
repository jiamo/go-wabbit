package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"wabbit-go/parser"
	"wabbit-go/wasm"
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
	err = os.WriteFile("out.wat", []byte(wasm.Wasm(prog)), 0644)
	if err != nil {
		log.Fatalf("Failed to write to out.wat: %v", err)
	}
	//log.Debugf("node_modules/wabt/bin/wat2wasm")
	cmd := exec.Command("node_modules/wabt/bin/wat2wasm", "--enable-tail-call", "out.wat")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run wat2wasm: %v", err)
	}
	log.Debugf("node test.js")
	//cmd = exec.Command("/Users/jiamo/.nvm/versions/node/v18.14.2/bin/node", "./test.js")
	//output, err := cmd.Output()
	//if err != nil {
	//	log.Fatalf("Failed to run test.js: %v", err)
	//}
	//fmt.Print(string(output))
	cmd = exec.Command("/Users/jiamo/.nvm/versions/node/v18.14.2/bin/node", "--experimental-wasm-return_call", "./test.js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run test.js: %v", err)
	}
}
