package main

import (
	"fmt"
	"os"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func main() {
	if os.Getenv("RUN") == "1" {
		info := helper.BootupDemoProject("../config")
		<-info.Running
		fmt.Println("Server existed")
	} else {
		fmt.Println("Generating GraphQl")
		helper.GenerateGraphQL()
	}
}
