package main

import (
	"ambassador-app/src/database"
	"ambassador-app/src/models"
	"context"
	"github.com/go-redis/redis/v8"
)

func main() {
	database.Connect()
	database.SetupRedis()

	ctx := context.Background()

	var users []models.User
	database.DBConn.Find(&users, models.User{
		IsAmbassador: true,
	})
	for _, user := range users {
		ambassador := models.Ambassador(user)
		ambassador.CalculateRevenue(database.DBConn)

		database.Cache.ZAdd(ctx, "rankings", &redis.Z{
			Score:  *ambassador.Revenue,
			Member: user.Name(),
		})

	}
}
