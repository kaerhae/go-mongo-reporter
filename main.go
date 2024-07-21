package main

import (
	"flag"
	"fmt"
	"log"
	"main/cmd/middleware"
	"main/cmd/routes"
	"main/configs"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	useSyslog := flag.Bool("syslog", true, "Determine if app use both stdout and syslog for logging, or only stdout")
	flag.Parse()
	logger := middleware.NewSyslogger(*useSyslog)
	router := routes.SetupRouter(logger)
	log.Fatal(router.Run(fmt.Sprintf("%s:%s", configs.GetAddress(), configs.GetPort())))
}
