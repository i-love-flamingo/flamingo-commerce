package graphql

import (
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -nometadata -o schema.go -pkg graphql schema.graphql

// Service is the Graphql-Service of this module
type Service struct{}

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Models return the 'Schema name' => 'Go model' mapping of this module
func (*Service) Models() map[string]config.TypeMapEntry {
	return graphql.ModelMap{
		"Commerce_Price": graphql.ModelMapEntry{
			Type: domain.Price{},
			Fields: map[string]string{
				"amount": "FloatAmount",
			},
		},
		"Commerce_Price_Charges":              domain.Charges{},
		"Commerce_Price_Charge":               domain.Charge{},
		"Commerce_Price_ChargeQualifier":      domain.ChargeQualifier{},
		"Commerce_Price_ChargeQualifierInput": domain.ChargeQualifier{},
	}.Models()
}
