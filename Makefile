all: clean setup install

setup-deps:
	go get -u -v -d ./...

build:
	go build -o ./bin/radium ./cmd/radium/*.go

install:
	go install ./cmd/radium/

setup:
	mkdir -p ./bin

clean:
	rm -rf ./bin


.PHONY:	all	build	setup	clean
