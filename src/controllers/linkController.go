package controllers

import (
	"ambassador-app/src/database"
	"ambassador-app/src/middleware"
	"ambassador-app/src/models"
	"github.com/bxcodec/faker/v4"
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

type CreateLinkRequest struct {
	Products []int
}

func CreateLink(c *fiber.Ctx) error {
	var req CreateLinkRequest

	err := c.BodyParser(&req)
	if err != nil {
		return err
	}
	id, _ := middleware.GetUserId(c)
	link := models.Link{
		UserId: id,
		Code:   faker.Username(),
	}
	for _, productId := range req.Products {
		product := models.Product{}
		product.Id = uint(productId)
		link.Products = append(link.Products, product)
	}
	database.DBConn.Create(&link)
	return c.JSON(link)
}

func Stats(c *fiber.Ctx) error {
	id, _ := middleware.GetUserId(c)

	var links []models.Link
	database.DBConn.Find(&links, models.Link{
		UserId: id,
	})
	var result []interface{}

	var orders []models.Order
	for _, link := range links {
		database.DBConn.Preload("OrderItem").Find(&orders, models.Order{
			Code:     link.Code,
			Complete: true,
		})
		revenue := 0.0
		for _, order := range orders {
			revenue += order.GetTotal()

		}
		result = append(result, fiber.Map{
			"code":    link.Code,
			"count":   len(orders),
			"revenue": revenue,
		})
	}
	return c.JSON(result)

}

func GetLink(c *fiber.Ctx) error {
	code := c.Params("code")

	link := models.Link{
		Code: code,
	}
	database.DBConn.Preload("User").Preload("Products").First(&link)
	return c.JSON(link)
}
