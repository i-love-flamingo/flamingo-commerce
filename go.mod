module flamingo.me/flamingo-commerce/v3

replace (
	flamingo.me/flamingo/v3 => ../flamingo
	flamingo.me/form => ../form
	flamingo.me/pugtemplate => ../pugtemplate
	flamingo.me/redirects => ../redirects
	golang.org/x/oauth2 => github.com/Ompluscator/oauth2 v0.0.0-20190101010101-b7626857
)

require (
	flamingo.me/dingo v0.1.3
	flamingo.me/flamingo/v3 v3.0.0-alpha4
	flamingo.me/pugtemplate v0.0.0-20190214131921-b2e86e90c6a5
	github.com/go-playground/form v3.1.4+incompatible
	github.com/go-test/deep v1.0.1
	github.com/gorilla/sessions v1.1.3
	github.com/leebenson/conform v0.0.0-20180615210222-bc2e0311fd85
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
	go.aoe.com/flamingo/form v0.0.0-20190214145643-d0d66a148576
	go.opencensus.io v0.19.0
	golang.org/x/oauth2 v0.0.0-20190212230446-3e8b2be13635
	gopkg.in/go-playground/validator.v9 v9.27.0
)
