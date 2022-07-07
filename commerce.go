package commerce

/*
	Flamingo Commerce Modules.
	The subpackage represent Flamingo Commerce modules - please refer to the documentations for the individual modules.
*/

//go:generate rm -rf docs/openapi
//go:generate go run github.com/swaggo/swag/cmd/swag@v1.6.6-0.20200603163350-20638f327979 init -p pascalcase --generalInfo=commerce.go --dir=./ --output=docs/openapi

// Swagger Documentation used for generator swag (https://github.com/swaggo/swag#declarative-comments-format)
// @title Flamingo Commerce API Spec
// @description Swagger (OpenAPI) Spec of all Flamingo Commerce modules
// @version 1.0
// @contact.name Flamingo
// @contact.url https://gitter.im/i-love-flamingo/community#
// @contact.email flamingo@aoe.com
// @license.name MIT
// @tag.name Cart
// @tag.description All Cart related APIs endpoints, most suitable to be called from a browser, because they rely on the session and cookie headers.
// @tag.name Payment
// @tag.description All Payment related APIs endpoints, most suitable to be called from a browser, because they rely on the session and cookie headers.
// @tag.name Product
// @tag.description All Product related APIs endpoints.
// @tag.name Checkout
// @tag.description  All Checkout related APIs endpoints, most suitable to be called from a browser, because they rely on the session and cookie headers.
