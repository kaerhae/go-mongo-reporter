package main

import (
	"flag"
	"log"
	"main/configs"
	"main/migrations"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var adminUser string
	var adminPass string
	flag.StringVar(&adminUser, "admin-user", "", "REPORTER_ROOT_USER")
	flag.StringVar(&adminPass, "admin-pass", "", "REPORTER_ROOT_PASS")
	flag.Parse()
	if len(os.Args) == 1 {
		log.Fatal("Missing argument. Set argument to up or down")
	}
	option := os.Args[1]

	if adminUser == "" {
		adminUser = configs.GetReporterRootUsername()
		if adminUser == "" {
			log.Fatal("No REPORTER_ROOT_USER set")
		}
	}
	if adminPass == "" {
		adminPass = configs.GetReporterRootPass()
		if adminPass == "" {
			log.Fatal("No REPORTER_ROOT_PASSWORD set")
		}
	}

	migrations.CreateAdminUser(option, adminUser, adminPass)

}
