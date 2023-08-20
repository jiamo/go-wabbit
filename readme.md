
I take the course from https://www.dabeaz.com/compiler.html.  
The course was taught by David Beazley using Python. 

I think it is fun to rewrite the code in Golang. (study golang)   
If you are interested in this project, You may want to take part in David's course.    
The course itself is also continuously evolving.

## prepare
    npm install # for wasm

## llvm
    # make sure you have clang
    go run cmd/llvm/llvm_main.go tests/Programs/23_mandel.wb

## wasm
    go run cmd/wasm/wasm_main.go tests/Programs/23_mandel.wb

## wvm
    go run cmd/wvm/wvm_main.go tests/Programs/23_mandel.wb

## interpreter
    go run cmd/interpreter/interpreter_main.go tests/Programs/23_mandel.wb

## test
    go test -v wabbit-go/tests

## TODO 
- [x] refactor code. Such like the implement of context. Do golang have good way to do it?
- [x] error handle
- [x] string
- [x] nest function
- [x] interpreter and wvm is so slowly. Optimize it.
- [x] There should have many todos