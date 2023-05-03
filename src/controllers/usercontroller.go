package controllers

import (
	"ambassador-app/src/database"
	"ambassador-app/src/models"
	"github.com/gofiber/fiber/v2"
)

func Ambassadors(c *fiber.Ctx) error {
	var users []models.User
	database.DBConn.Where("is_ambassador = true").Find(&users)
	return c.JSON(users)
}
