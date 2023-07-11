GOPATH:=$(shell go env GOPATH)

.PHONY: run build build-mac build-linux build-windows test clean

run:
	go build . # && go run main.go

build: build-mac build-linux build-windows

build-mac: clean
	GOOS=darwin go build -trimpath -o dist/mac/wabiit .

build-linux: clean
	GOOS=linux go build -trimpath -o dist/linux/wabbit .

#build-wasm: clean
#	@cd src && \
#	GOOS=js GOARCH=wasm go build -trimpath -o ../dist/wasm/ghost.wasm wasm/wasm.go

build-windows: clean
	GOOS=windows go build -trimpath -o dist/windows/wabbit.exe .

test:
	MONKEY_LIBPATH=`pwd`/libs go test -v -race -timeout 5s `go list ./...` | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

clean:
	@rm -rf dist/mac
	@rm -rf dist/linux
	@rm -rf dist/windows