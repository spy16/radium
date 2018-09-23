VERSION:=$(shell git describe --abbrev=0 --tags)
COMMIT:=$(shell git rev-parse HEAD)

all: clean setup install

setup-deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v

build:
	go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT)" -o ./bin/radium ./cmd/radium/*.go

install:
	go install -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT)"  ./cmd/radium/

setup:
	mkdir -p ./bin

clean:
	rm -rf ./bin


.PHONY:	all	build	setup	clean
