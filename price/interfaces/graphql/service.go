package graphql

import (
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

type Service struct{}

func (*Service) Schema() []byte {
	//language=graphql
	return []byte(`
type Commerce_Price{
	amount: Float
	currency: String!
}
`)
}

func (*Service) Models() map[string]config.TypeMapEntry {
	return graphql.ModelMap{
		"Commerce_Price": graphql.ModelMapEntry{
			Type: domain.Price{},
			Fields: map[string]string{
				"amount": "FloatAmount",
			},
		},
	}.Models()
}
