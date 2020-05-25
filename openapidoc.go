package flamingo_commerce

import _ "flamingo.me/flamingo-commerce/v3/cart"

//go:generate go run github.com/swaggo/swag/cmd/swag init --generalInfo=openapidoc.go --dir=./ --output=docs/openapi

// ***** Swagger Documentation annotation for generator swag (https://github.com/swaggo/swag#declarative-comments-format) ******
// @title Flamingo Commerce API Spec
// @description Swagger (OpenAPI) Spec of all Flamingo Commerce modules
// @version 1.0
// @contact.name Flamingo
// @contact.url https://gitter.im/i-love-flamingo/community#
// @contact.email flamingo@aoe.com
// @license.name MIT
// *****
// @tag.name v1 Cart ajax API
// @tag.description This Cart APIs are most suitable to be called from a browser, because they rely on the session and cookie headers.
