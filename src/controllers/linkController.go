package controllers

import (
	"ambassador-app/src/database"
	"ambassador-app/src/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func Link(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var links []models.Link

	database.DBConn.Find(&links).Where("user_id =? ", id)
	for i, link := range links {
		var orders []models.Order
		database.DBConn.Where("code = ? and complete= true", link.Code).Find(&orders)
		links[i].Orders = orders
	}
	return c.JSON(links)
}
