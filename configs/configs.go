package configs

import (
	"fmt"
	"os"
)

func GetMongoURI() string {
	user := os.Getenv("MONGO_USER")
	pass := os.Getenv("MONGO_PASS")
	ip_addr := os.Getenv("MONGO_IP")
	port := os.Getenv("MONGO_PORT")
	mongo_uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, ip_addr, port)
	return mongo_uri
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
