.PHONY: build clean
.DEFAULT_GOAL := build

build:
	GOBIN=$(shell pwd)/bin go install -mod=vendor -a -ldflags '-w -s' -v ./...

clean:
	rm -rf bin
