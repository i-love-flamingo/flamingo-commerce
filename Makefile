CONTEXT?=dev
.PHONY: up localup update test

test:
	go test -race -v ./...

