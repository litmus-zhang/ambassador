package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

const secretKey = "secretkey"

type ClaimsWithScope struct {
	jwt.StandardClaims
	Scope string
}

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}
	payload := token.Claims.(*ClaimsWithScope)
	isAmbassador := strings.Contains(c.Path(), "/api/ambassador")
	if (payload.Scope == "admin" && isAmbassador) || (payload.Scope == "ambassador" && !isAmbassador) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	return c.Next()
}

func GetUserId(c *fiber.Ctx) (uint, error) {

	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return 0, err
	}
	payload := token.Claims.(*ClaimsWithScope)
	id, _ := strconv.Atoi(payload.Subject)
	return uint(id), nil
}

func GenerateJWT(id uint, scope string) (string, error) {
	var payload = ClaimsWithScope{}
	payload.Subject = strconv.Itoa(int(id))
	payload.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	payload.Scope = scope

	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(secretKey))
}
