package migrations

import (
	"log"
	"main/configs"
	"main/pkg/db"
	"main/pkg/models"
	"main/pkg/repository"
	"main/pkg/utils"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateAdminUser(option string, adminUser string, adminPass string) {
	switch strings.ToLower(option) {
	case "up":
		createAdminUserUp(adminUser, adminPass)
		log.Default().Println("Successfully created admin user")
		os.Exit(0)
	case "down":
		createAdminUserDown()
		log.Default().Println("Successfully deleted admin user")
		os.Exit(0)
	default:
		log.Fatal("Invalid option for --create-user-admin")
	}
}

func createAdminUserUp(adminUser string, adminPass string) {
	_, client, _, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		log.Fatal("Error on connecting MongoDB", err)
	}
	r := repository.NewUserRepository(client)
	_, err = r.GetSingleUserByUsername(adminUser)
	if err == nil {
		log.Default().Println("migration already done and user exists")
		os.Exit(0)
	}

	adminID, err := primitive.ObjectIDFromHex("000011112222333344445555")
	if err != nil {
		log.Fatal("Error while creating ID ", err)
	}
	hashPwd, err := utils.HashPwd(adminPass)
	if err != nil {
		log.Fatal("Error while hashing password: ", err)
	}

	admin := models.User{
		ID:           adminID,
		Username:     adminUser,
		Email:        "",
		PasswordHash: hashPwd,
		CreatedAt:    time.Now().UTC().String(),
		Permission: models.Permission{
			Admin: true,
			Write: true,
			Read:  true,
		},
		Reports: []primitive.ObjectID{},
	}

	_, err = r.Create(&admin)
	if err != nil {
		log.Fatal("Error while creating admin: ", err)
	}
}

func createAdminUserDown() {
	_, client, _, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		log.Fatal("Error on connecting MongoDB", err)
	}

	r := repository.NewUserRepository(client)
	deleteCount, err := r.DeleteSingleUser("000011112222333344445555")
	if err != nil {
		log.Fatal("Error while deleting admin user: ", err)
	}
	if deleteCount == 0 {
		log.Default().Println("no user found")
		os.Exit(0)
	}

}
