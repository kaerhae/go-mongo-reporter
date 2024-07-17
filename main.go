package main

import (
	"fmt"
	"log"
	"main/cmd/routes"
	"main/configs"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	router := routes.SetupRouter()
	log.Fatal(router.Run(fmt.Sprintf("localhost:%s", configs.GetPort())))
}
