BINARY=arknights
LDFLAGS=-ldflags "-s -w"

all: build

build:
	rm -f arknights
	go build ${LDFLAGS} -o ${BINARY}