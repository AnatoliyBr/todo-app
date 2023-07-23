.PHONY: build
build:
			go build -v ./cmd/app

.PHONY: test
test:
			go test -v -race ./...

.DEFAULT_GOAL := build