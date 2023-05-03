package controllers

import (
	"ambassador-app/src/database"
	"ambassador-app/src/middleware"
	"ambassador-app/src/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"strconv"
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

	user := models.User{
		FirstName:    data["firstname"],
		Lastname:     data["lastname"],
		Email:        data["email"],
		IsAmbassador: false,
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
	var payload = jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))
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
