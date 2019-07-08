GOLANG_VERSION=1.12.6
MOUNT_POINT=/project

.PHONY: build clean lint
.DEFAULT_GOAL := build

build:
	GOBIN=$(shell pwd)/bin go install -mod=vendor -ldflags '-w -s' -v ./...

clean:
	rm -rf bin

lint:
	docker run --rm \
		--user `id -u`:`id -g` \
		--env GOCACHE=$(MOUNT_POINT)/bin/.cache \
		--volume `pwd`:$(MOUNT_POINT) \
		--volume $(GOPATH):/go \
		--workdir $(MOUNT_POINT) \
		golangci/golangci-lint golangci-lint run
