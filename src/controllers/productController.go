package controllers

import (
	"ambassador-app/src/database"
	"ambassador-app/src/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"sort"
	"strconv"
	"strings"
	"time"
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
	go database.ClearCache("products_frontend", "products_backend")

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
func ProductsFrontend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()
	results, err := database.Cache.Get(ctx, "products_frontend").Result()
	if err != nil {
		fmt.Println(err.Error())
		database.DBConn.Find(&products)

		bytes, err := json.Marshal(products)
		if err != nil {
			panic(err.Error())
		}

		errKey := database.Cache.Set(ctx, "products_frontend", bytes, 30*time.Minute).Err()
		if errKey != nil {
			panic(errKey.Error())

		}
	} else {
		json.Unmarshal([]byte(results), &products)
	}
	return c.JSON(results)

}

func ProductsBackend(c *fiber.Ctx) error {
	var products []models.Product
	var ctx = context.Background()
	results, err := database.Cache.Get(ctx, "products_frontend").Result()
	if err != nil {
		fmt.Println(err.Error())
		database.DBConn.Find(&products)

		bytes, err := json.Marshal(products)
		if err != nil {
			panic(err.Error())
		}

		database.Cache.Set(ctx, "products_backend", bytes, 30*time.Minute).Err()

	} else {
		json.Unmarshal([]byte(results), &products)
	}
	var searchProducts []models.Product
	if s := c.Query("s"); s != "" {
		lower := strings.ToLower(s)
		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), lower) || strings.Contains(strings.ToLower(product.Description), lower) {
				searchProducts = append(searchProducts, product)
			}
		}

	} else {
		searchProducts = products
	}
	if s := c.Query("sort"); s != "" {
		sortLower := strings.ToLower(s)
		if sortLower == "asc" {
			sort.Slice(searchProducts, func(i, j int) bool {
				return searchProducts[i].Price < searchProducts[j].Price

			})
		} else if sortLower == "desc" {
			sort.Slice(searchProducts, func(i, j int) bool {
				return searchProducts[i].Price > searchProducts[j].Price

			})

		}

	}
	var total = len(searchProducts)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	var data []models.Product
	perPage := 9

	if total <= page*perPage && total >= (page-1)*perPage {
		data = searchProducts[(page-1)*perPage : total]
	} else if total >= page*perPage {
		data = searchProducts[(page-1)*perPage : page*perPage]
	} else {
		data = []models.Product{}
	}

	return c.JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": total/perPage + 1,
	})

}
