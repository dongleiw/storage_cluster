lib=$(shell readlink -f ../../)
all:
	mkdir -p build
	GOPATH=${lib}:${GOPATH} go build -o build/tool *.go
