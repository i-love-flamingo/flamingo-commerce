package docs

import (
	// embed openapi/swagger.json
	_ "embed"
)

//go:embed openapi/swagger.json
var swaggerJSON []byte

// Asset returns Flamingo Commerce swagger json definition
func Asset(_ string) ([]byte, error) {
	return swaggerJSON, nil
}
