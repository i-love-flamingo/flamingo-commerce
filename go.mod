module flamingo.me/flamingo-commerce/v3

require (
	flamingo.me/dingo v0.1.5
	flamingo.me/flamingo/v3 v3.0.0-beta.2.0.20190515120627-9cabe248cf01
	flamingo.me/form v1.0.0-alpha.1
	flamingo.me/pugtemplate v1.0.0-alpha.1
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
	github.com/stretchr/testify v1.3.0
	go.opencensus.io v0.19.1
	golang.org/x/net v0.0.0-20190311183353-d8887717615a // indirect
	golang.org/x/oauth2 v0.0.0-20190226205417-e64efc72b421 // indirect
	golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6 // indirect
	google.golang.org/grpc v1.19.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1
	gopkg.in/square/go-jose.v2 v2.3.0 // indirect
)

replace (
	github.com/robertkrimen/otto => github.com/thebod/otto v0.0.0-20180101010101-83d297c4b64aeb2de4268d9a54c9a503ae2d8139
	golang.org/x/oauth2 => github.com/Ompluscator/oauth2 v0.0.0-20190101010101-b7626857
)
