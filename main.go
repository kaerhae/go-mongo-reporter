package main

import (
	"fmt"
	"log"
	"main/cmd/routes"
	"main/configs"
)

func main() {

	router := routes.SetupRouter()
	log.Fatal(router.Run(fmt.Sprintf("localhost:%s", configs.GetPort())))
}
