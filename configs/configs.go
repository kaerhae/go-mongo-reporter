package configs

import (
	"os"
)

func GetMongoURI() string {
	return os.Getenv("DATABASE_URI")
}

func GetDBName() string {
	return os.Getenv("DATABASE")
}

func GetSecret() string {
	return os.Getenv("SECRET_KEY")
}

func GetPort() string {
	return os.Getenv("PORT")
}
