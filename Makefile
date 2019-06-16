BINARY=arknights
LDFLAGS=-ldflags "-s -w"

all: build

build:
	rm -f ${BINARY}
	go build ${LDFLAGS} -o ${BINARY}