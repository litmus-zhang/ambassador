package main

import (
	"ambassador-app/src/database"
	"ambassador-app/src/models"
	"github.com/bxcodec/faker/v4"
)

func main() {
	database.Connect()
	for i := 0; i < 30; i++ {
		ambassador := models.User{
			FirstName:    faker.LastName(),
			Lastname:     faker.LastName(),
			Email:        faker.Email(),
			IsAmbassador: true,
		}
		ambassador.SetPassword("123456")
		database.DBConn.Create(&ambassador)
	}
}
