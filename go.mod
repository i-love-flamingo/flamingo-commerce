module flamingo.me/flamingo-commerce/v3

require (
	flamingo.me/dingo v0.1.4
	flamingo.me/flamingo/v3 v3.0.0-beta.2.0.20190423070243-a5aa37396b82
	flamingo.me/form v1.0.0-alpha.1
	flamingo.me/pugtemplate v1.0.0-alpha.1
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/go-playground/form v3.1.4+incompatible
	github.com/go-test/deep v1.0.1
	github.com/leekchan/accounting v0.0.0-20180703100437-18a1925d6514
	github.com/lib/pq v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/procfs v0.0.0-20190225181712-6ed1f7e10411 // indirect
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24 // indirect
	github.com/stretchr/testify v1.3.0
	go.opencensus.io v0.19.1
	golang.org/x/crypto v0.0.0-20190225124518-7f87c0fbb88b // indirect
	golang.org/x/net v0.0.0-20190226193003-66a96c8a540e // indirect
	golang.org/x/oauth2 v0.0.0-20190226191147-529b322ea346 // indirect
)

replace (
	github.com/robertkrimen/otto => github.com/thebod/otto v0.0.0-20180101010101-83d297c4b64aeb2de4268d9a54c9a503ae2d8139
	golang.org/x/oauth2 => github.com/Ompluscator/oauth2 v0.0.0-20190101010101-b7626857
)
