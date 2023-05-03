package controllers

import (
	"ambassador-app/src/database"
	"ambassador-app/src/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func Products(c *fiber.Ctx) error {
	var products []models.Product
	database.DBConn.Find(&products)
	return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return err
	}
	database.DBConn.Create(&product)

	return c.JSON(product)

}

func GetProduct(c *fiber.Ctx) error {
	var product models.Product

	id, _ := strconv.Atoi(c.Params("id"))
	product.Id = uint(id)
	database.DBConn.Find(&product)
	return c.JSON(product)
}
func UpdateProduct(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Params("id"))
	product := models.Product{}
	product.Id = uint(id)
	if err := c.BodyParser(&product); err != nil {
		return err
	}
	database.DBConn.Model(&product).Updates(&product)
	return c.JSON(product)
}
func DeleteProduct(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Params("id"))
	product := models.Product{}
	product.Id = uint(id)

	if err := c.BodyParser(&product); err != nil {
		return err
	}
	database.DBConn.Delete(&product)
	return c.JSON(fiber.Map{
		"message": "Object deleted success",
	})
}
