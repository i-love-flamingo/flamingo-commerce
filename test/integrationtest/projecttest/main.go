package main

import (
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
	"fmt"
	"os"
)

func main() {
	if os.Getenv("RUN") == "1" {
		info := helper.BootupDemoProject()
		<-info.Running
		fmt.Println("Server existed")
	} else {
		fmt.Println("Generating GraphQl")
		helper.GenerateGraphQL()
	}
}
