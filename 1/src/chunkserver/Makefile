common=$(shell readlink -f ../../)
all:
	mkdir -p build
	GOPATH=${common}:${GOPATH} go build -o build/chunkserver *.go  
