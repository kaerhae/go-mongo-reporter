package main

import (
	"flag"
	"fmt"
	"log"
	"main/api"
	"main/configs"
	"main/pkg/middleware"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	useSyslog := flag.Bool("syslog", true, "Determine if app use both stdout and syslog for logging, or only stdout")
	flag.Parse()

	mongoUser := os.Getenv("MONGO_USER")
	mongoPass := os.Getenv("MONGO_PASS")
	mongoIP := os.Getenv("MONGO_IP")
	mongoPort := os.Getenv("MONGO_PORT")
	db := os.Getenv("DATABASE")
	if mongoUser == "" || mongoPass == "" || mongoIP == "" || mongoPort == "" || db == "" {
		log.Fatal("MongoDB environment variables are missing")
	}

	logger := middleware.NewSyslogger(*useSyslog)
	router := api.SetupRouter(logger)
	log.Fatal(router.Run(fmt.Sprintf("%s:%s", configs.GetAddress(), configs.GetPort())))
}
