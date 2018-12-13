CONTEXT?=dev
.PHONY: up localup update test

up:
	rm -rf vendor/
	dep ensure -v -vendor-only

update:
	rm -rf vendor/
	dep ensure -v -update flamingo.me/flamingo

localup: up local
	
local:
	rm -rf vendor/flamingo.me/flamingo
	ln -sf ../../../flamingo vendor/flamingo.me/flamingo
	rm -rf vendor/flamingo.me/flamingo/vendor
	
test:
	go test -race -v ./...

updateTools:
	go get -v -u github.com/golang/dep/cmd/dep
