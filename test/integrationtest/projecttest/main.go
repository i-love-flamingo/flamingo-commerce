// +build integration

package main

import (
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Generating GraphQl")
	helper.GenerateGraphQL()
	if os.Getenv("RUN") == "1" {
		info := helper.BootupDemoProject()
		<-info.Running
		fmt.Println("Server existed")
	}
}
