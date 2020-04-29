package graphql

import (
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"

	"flamingo.me/flamingo-commerce/v3/customer/domain"
	"flamingo.me/flamingo-commerce/v3/customer/interfaces/graphql/dtocustomer"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -nometadata -o fs.go -pkg graphql schema.graphql

// Service is the Graphql-Service of this module
type Service struct{}

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Models return the 'Schema name' => 'Go model' mapping of this module
func (*Service) Models() map[string]config.TypeMapEntry {
	return graphql.ModelMap{
		"Commerce_Customer_Status_Result": dtocustomer.CustomerStatusResult{},
		"Commerce_Customer_Result":        dtocustomer.CustomerResult{},
		"Commerce_Customer_PersonData":    domain.PersonData{},
		"Commerce_Customer_Address":       domain.Address{},
	}.Models()
}
