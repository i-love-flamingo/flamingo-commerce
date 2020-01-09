// +build integration

package main

import (
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
	"fmt"
)

func main() {
	fmt.Println("Generating GraphQl")
	helper.GenerateGraphQL()
}
