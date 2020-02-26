CONTEXT?=dev
REPLACE?=-replace flamingo.me/flamingo/v3=../flamingo -replace flamingo.me/form=../form
DROPREPLACE?=-dropreplace flamingo.me/flamingo/v3 -dropreplace flamingo.me/form

.PHONY: local unlocal test

local:
	git config filter.gomod-flamingo-commerce.smudge 'go mod edit -fmt -print $(REPLACE) /dev/stdin'
	git config filter.gomod-flamingo-commerce.clean 'go mod edit -fmt -print $(DROPREPLACE) /dev/stdin'
	git config filter.gomod-flamingo-commerce.required true
	go mod edit -fmt $(REPLACE)

unlocal:
	git config filter.gomod-flamingo-commerce.smudge ''
	git config filter.gomod-flamingo-commerce.clean ''
	git config filter.gomod-flamingo-commerce.required false
	go mod edit -fmt $(DROPREPLACE)

test:
	go test -race -v ./...
	gofmt -l -e -d .
	golint ./...
	misspell -error .
	ineffassign .

integrationtest:
	go test -v ./test/integrationtest/... -tags=integration

generateprojecttestgraphql:
	rm -f test/integrationtest/projecttest/graphql/generated.go
	cd test/integrationtest/projecttest && go run -tags graphql main.go

fix:
	gofmt -l -w .

rundemoprojecttest:
	export RUN="1" && cd test/integrationtest/projecttest/tests && go run ../main.go