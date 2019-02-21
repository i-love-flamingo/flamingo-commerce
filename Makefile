CONTEXT?=dev
REPLACE?=-replace flamingo.me/flamingo/v3=../flamingo
DROPREPLACE?=-dropreplace flamingo.me/flamingo/v3

.PHONY: local unlocal

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
