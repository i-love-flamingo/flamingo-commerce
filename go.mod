module flamingo.me/flamingo-commerce/v3

replace (
	flamingo.me/flamingo/v3 => ../flamingo
	flamingo.me/form => ../form
	flamingo.me/pugtemplate => ../pugtemplate
	flamingo.me/redirects => ../redirects
)

require (
	flamingo.me/dingo v0.1.3
	flamingo.me/flamingo v0.0.0-20190122075217-ac03fb2ca2e2
	git.apache.org/thrift.git v0.0.0-20180807212849-6e67faa92827
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973
	github.com/boj/redistore v0.0.0-20160128113310-fc113767cd6b
	github.com/coreos/go-oidc v2.0.0+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/etgryphon/stringUp v0.0.0-20121020160746-31534ccd8cac
	github.com/garyburd/redigo v1.6.0
	github.com/ghodss/yaml v0.0.0-20180820084758-c7ce16629ff4
	github.com/go-playground/form v3.1.3+incompatible
	github.com/go-playground/locales v0.12.1
	github.com/go-playground/universal-translator v0.16.0
	github.com/go-test/deep v1.0.1
	github.com/golang/protobuf v1.2.0
	github.com/google/uuid v1.1.0
	github.com/gorilla/context v1.1.1
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.1.3
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/leebenson/conform v0.0.0-20180615210222-bc2e0311fd85
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/pkg/errors v0.8.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910
	github.com/prometheus/common v0.0.0-20180801064454-c7de2306084e
	github.com/prometheus/procfs v0.0.0-20180725123919-05ee40e3a273
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/objx v0.1.1
	github.com/stretchr/testify v1.2.2
	github.com/zemirco/memorystore v0.0.0-20160308183530-ecd57e5134f6
	go.opencensus.io v0.0.0-20180823191657-71e2e3e3082a
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a
	golang.org/x/sync v0.0.0-20181108010431-42b317875d0f
	google.golang.org/api v0.0.0-20180824000442-943e5aafc110
	google.golang.org/appengine v1.3.0
	gopkg.in/go-playground/validator.v9 v9.21.1
	gopkg.in/sourcemap.v1 v1.0.5
	gopkg.in/square/go-jose.v2 v2.1.9
	gopkg.in/yaml.v2 v2.2.2
)
