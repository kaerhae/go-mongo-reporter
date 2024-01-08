package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func DotenvInit() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error on loading .env file")
	}
}

func GetMongoURI() string {
	DotenvInit()
	return os.Getenv("DATABASE_URI")
}

func GetDBName() string {
	DotenvInit()
	return os.Getenv("DATABASE")
}

func GetSecret() string {
	DotenvInit()
	return os.Getenv("SECRET_KEY")
}
