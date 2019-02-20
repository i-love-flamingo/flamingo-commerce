module flamingo.me/flamingo-commerce/v3

require (
	flamingo.me/dingo v0.1.3
	flamingo.me/flamingo/v3 v3.0.0-alpha5
	flamingo.me/pugtemplate v0.0.0-20190214131921-b2e86e90c6a5
	github.com/go-playground/form v3.1.4+incompatible
	github.com/go-test/deep v1.0.1
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gorilla/sessions v1.1.3
	github.com/kisielk/errcheck v1.2.0 // indirect
	github.com/leebenson/conform v0.0.0-20180615210222-bc2e0311fd85
	github.com/pkg/errors v0.8.1
	github.com/prometheus/procfs v0.0.0-20190219184716-e4d4a2206da0 // indirect
	github.com/stretchr/testify v1.3.0
	go.aoe.com/flamingo/form v0.0.0-20190214145643-d0d66a148576
	go.opencensus.io v0.19.0
	go4.org v0.0.0-20190218023631-ce4c26f7be8e // indirect
	golang.org/x/build v0.0.0-20190219204446-88cd9dd98818 // indirect
	golang.org/x/crypto v0.0.0-20190219172222-a4c6cb3142f2 // indirect
	golang.org/x/oauth2 v0.0.0-20190219183015-4b83411ed2b3
	golang.org/x/sys v0.0.0-20190219203350-90b0e4468f99 // indirect
	golang.org/x/tools v0.0.0-20190219185102-9394956cfdc5 // indirect
	google.golang.org/genproto v0.0.0-20190219182410-082222b4a5c5 // indirect
	gopkg.in/go-playground/validator.v9 v9.27.0
	honnef.co/go/tools v0.0.0-20190215041234-466a0476246c // indirect
)

replace (
	flamingo.me/flamingo-commerce/v3 => ../flamingo-commerce
	flamingo.me/flamingo/v3 => ../flamingo
	github.com/robertkrimen/otto => github.com/thebod/otto v0.0.0-20180101010101-83d297c4b64aeb2de4268d9a54c9a503ae2d8139
	golang.org/x/oauth2 => github.com/Ompluscator/oauth2 v0.0.0-20190101010101-b7626857
)
