package main

import (
	"ambassador-app/src/database"
	"ambassador-app/src/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

const PORT = ":8100"

func main() {
	database.Connect()
	database.AutoMigrate()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	routes.Setup(app)

	log.Fatal(app.Listen(PORT))
}
