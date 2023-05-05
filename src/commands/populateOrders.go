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
		var orderitems []models.OrderItem

		for j := 0; j < rand.Intn(5); j++ {
			price := float64(rand.Intn(90) + 10)
			qty := uint(rand.Intn(5))

			orderitems = append(orderitems, models.OrderItem{
				ProductTitle:      faker.Word(),
				Price:             price,
				Quantity:          qty,
				AdminRevenue:      0.9 * price * float64(qty),
				AmbassadorRevenue: 0.1 * price * float64(qty),
			})

		}
		database.DBConn.Create(&models.Order{
			UserId:          uint(rand.Intn(30) + 1),
			Code:            faker.Username(),
			AmbassadorEmail: faker.Email(),
			FirstName:       faker.FirstName(),
			LastName:        faker.LastName(),
			Email:           faker.Email(),
			Complete:        true,
			OrderItems:      orderitems,
		})
	}
}
