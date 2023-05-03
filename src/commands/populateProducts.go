package main

import (
	"ambassador-app/src/database"
	"ambassador-app/src/models"
	"github.com/bxcodec/faker/v4"
	"math/rand"
)

func main() {
	database.Connect()
	for i := 0; i < 30; i++ {
		product := models.Product{
			Title:       faker.Username(),
			Description: faker.Username(),
			Image:       faker.URL(),
			Price:       float64(rand.Intn(90) + 10),
		}
		database.DBConn.Create(&product)
	}
}
