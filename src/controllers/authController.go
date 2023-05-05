package controllers

import (
	"ambassador-app/src/database"
	"ambassador-app/src/middleware"
	"ambassador-app/src/models"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

type Message map[string]string

func Register(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"Message": "Password do not match",
		})
	}
	query := database.DBConn.Where("email =?", data["email"]).Find(&models.User{})

	if query.Error != nil {
		return c.JSON(fiber.Map{
			"message": "Email already in use",
			"err":     query.Error.Error(),
		})
	}

	user := models.User{
		FirstName:    data["firstname"],
		Lastname:     data["lastname"],
		Email:        data["email"],
		IsAmbassador: strings.Contains(c.Path(), "/api/ambassador"),
	}
	user.SetPassword(data["password"])

	database.DBConn.Create(&user)
	return c.JSON(fiber.Map{
		"message": "User registration successful",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}
	var user models.User

	database.DBConn.Where("email = ?", data["email"]).First(&user)
	if user.Id == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid credentials",
		})
	}
	err = user.ComparePassword(data["password"])
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid credentials",
		})
	}
	isAmbassador := strings.Contains(c.Path(), "/api/ambassador")
	var scope string
	if isAmbassador {
		scope = "ambassador"
	} else {
		scope = "admin"
	}

	if !isAmbassador && user.IsAmbassador {
		return c.Status(400).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}
	token, err := middleware.GenerateJWT(user.Id, scope)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "invalid credential",
		})
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.Status(200).JSON(fiber.Map{
		"message": "login successful",
	})
}

func User(c *fiber.Ctx) error {
	id, _ := middleware.GetUserId(c)
	var user models.User
	database.DBConn.Where("id = ?", id).First(&user)

	if strings.Contains(c.Path(), "/api/ambassador") {
		ambassador := models.Ambassador(user)
		ambassador.CalculateRevenue(database.DBConn)
		return c.JSON(ambassador)
	}
	return c.JSON(user)

}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})

}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}
	id, _ := middleware.GetUserId(c)

	user := models.User{
		FirstName: data["firstname"],
		Lastname:  data["lastname"],
		Email:     data["email"],
	}
	user.Id = id
	database.DBConn.Model(&user).Updates(&user)
	return c.JSON(user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"Message": "Password do not match",
		})
	}

	id, _ := middleware.GetUserId(c)

	user := models.User{}
	user.Id = id
	user.SetPassword(data["password"])
	database.DBConn.Model(&user).Updates(&user)
	return c.JSON(user)
}
