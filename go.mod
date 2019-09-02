module flamingo.me/flamingo-commerce/v3

require (
	flamingo.me/dingo v0.1.6
	flamingo.me/flamingo/v3 v3.0.1
	flamingo.me/form v1.0.0
	flamingo.me/pugtemplate v1.0.0
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/go-playground/form v3.1.4+incompatible
	github.com/go-test/deep v1.0.1
	github.com/golang/protobuf v1.3.0 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/leekchan/accounting v0.0.0-20180703100437-18a1925d6514
	github.com/lib/pq v1.1.1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/procfs v0.0.0-20190306233201-d0f344d83b0c // indirect
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24 // indirect
	github.com/stretchr/testify v1.4.0
	go.opencensus.io v0.20.2
	gopkg.in/go-playground/assert.v1 v1.2.1
	gopkg.in/square/go-jose.v2 v2.3.0 // indirect
)

replace (
	github.com/robertkrimen/otto => github.com/thebod/otto v0.0.0-20180101010101-83d297c4b64aeb2de4268d9a54c9a503ae2d8139
	golang.org/x/oauth2 => github.com/Ompluscator/oauth2 v0.0.0-20190101010101-b7626857
)
