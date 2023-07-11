package main

import (
	"bufio"
	"fmt"
	"os"
	"wabbit-go/tokenize"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./tokenizer filename")
		os.Exit(1)
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		fmt.Println(statsErr)
		os.Exit(1)
	}

	size := stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	data := string(bytes)

	tokens, _ := tokenize.Tokenize(data)
	fmt.Println("total tokens ", len(tokens))
	for _, tok := range tokens {
		fmt.Println(tok)
	}
}
