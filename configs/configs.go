package configs

import (
	"fmt"
	"os"
)

func GetMongoURI() string {
	user := os.Getenv("MONGO_USER")
	pass := os.Getenv("MONGO_PASS")
	ipAddr := os.Getenv("MONGO_IP")
	port := os.Getenv("MONGO_PORT")
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, ipAddr, port)
	return mongoURI
}

func GetDBName() string {
	return os.Getenv("DATABASE")
}

func GetSecret() string {
	return os.Getenv("SECRET_KEY")
}

func GetAddress() string {
	return os.Getenv("IP_ADDR")
}

func GetPort() string {
	return os.Getenv("PORT")
}

func GetReporterRootUsername() string {
	return os.Getenv("REPORTER_ROOT_USER")
}

func GetReporterRootPass() string {
	return os.Getenv("REPORTER_ROOT_PASSWORD")
}
