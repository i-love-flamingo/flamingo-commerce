module flamingo.me/flamingo-commerce/v3

go 1.13

require (
	flamingo.me/dingo v0.1.6
	flamingo.me/flamingo/v3 v3.0.3
	flamingo.me/form v1.0.1-0.20191008191024-ff6f3a9330d6
	flamingo.me/graphql v1.0.1
	flamingo.me/pugtemplate v1.0.0
	github.com/99designs/gqlgen v0.9.0
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/go-playground/form v3.1.4+incompatible
	github.com/go-test/deep v1.0.1
	github.com/golang/mock v1.2.0 // indirect
	github.com/golang/protobuf v1.3.0 // indirect
	github.com/gorilla/sessions v1.1.3
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
	github.com/robertkrimen/otto => github.com/thebod/otto v0.0.0-20170712091932-83d297c4b64a
	golang.org/x/oauth2 => github.com/Ompluscator/oauth2 v0.0.0-20190121141151-b76268579942
)
