package controllers

import (
	"ambassador-app/src/database"
	"ambassador-app/src/models"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"net/smtp"
)

func Orders(c *fiber.Ctx) error {
	var orders []models.Order
	database.DBConn.Preload("OrderItems").Find(&orders)
	for i, order := range orders {
		orders[i].Name = order.FullName()
		orders[i].Total = order.GetTotal()
	}
	return c.JSON(orders)
}

type CreateOrderRequest struct {
	Code      string
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
	City      string
	Zip       string
	Products  []map[string]int
}

func CreateOrder(c *fiber.Ctx) error {
	var req CreateOrderRequest
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}
	link := models.Link{
		Code: req.Code,
	}
	tx := database.DBConn.Begin()
	database.DBConn.Preload("User").First(&link)
	if link.Id == 0 {
		c.Status(400).JSON(fiber.Map{
			"message": "Invalid Link",
		})
	}
	order := models.Order{
		Code:            link.Code,
		UserId:          link.UserId,
		AmbassadorEmail: link.User.Email,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Address:         req.Address,
		Country:         req.Country,
		City:            req.City,
		Zip:             req.Zip,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.Status(400).JSON(fiber.Map{
			"message": err,
		})
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, reqproduct := range req.Products {
		product := models.Product{}
		product.Id = uint(reqproduct["product_id"])
		database.DBConn.First(&product)

		total := product.Price * float64(reqproduct["quantity"])
		item := models.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(reqproduct["quantity"]),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.Status(400).JSON(fiber.Map{
				"message": err,
			})
		}
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Name:        stripe.String(product.Title),
			Description: stripe.String(product.Description),
			Images:      []*string{stripe.String(product.Image)},
			Amount:      stripe.Int64(100 * int64(product.Price)),
			Currency:    stripe.String("usd"),
			Quantity:    stripe.Int64(int64(reqproduct["product"])),
		})
	}
	stripe.Key = ""
	params := stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String("http://localhost:5000/success?source={CHECKOUT_SESSION_ID}"),
		CancelURL:          stripe.String("http://localhost:5000/error"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
	}
	source, err := session.New(&params)
	if err != nil {
		tx.Rollback()
		c.Status(400).JSON(fiber.Map{
			"message": err,
		})
	}
	order.TransactionId = source.ID
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.Status(400).JSON(fiber.Map{
			"message": err,
		})
	}

	tx.Commit()
	return c.JSON(source)
}

func CompleteOrder(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	order := models.Order{}
	database.DBConn.Preload("OrderItem").First(&order, models.Order{
		TransactionId: data["source"],
	})
	if order.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "order not found",
		})
	}

	order.Complete = true
	database.DBConn.Save(&order)

	go func(order models.Order) {
		ambassadorRevenue := 0.0
		adminRevenue := 0.0

		for _, item := range order.OrderItems {
			adminRevenue += item.AdminRevenue
			ambassadorRevenue += item.AmbassadorRevenue
		}
		user := models.User{}
		user.Id = order.UserId

		database.DBConn.First(&user)
		database.Cache.ZIncrBy(context.Background(), "rankings", ambassadorRevenue, user.Name())

		ambassadorMessage := []byte(fmt.Sprintf("You earned $%f from the link #%s", ambassadorRevenue, order.Code))

		smtp.SendMail("127.0.0.1:1025", nil, "no-reply@email.com", []string{order.AmbassadorEmail}, ambassadorMessage)
		adminMessage := []byte(fmt.Sprintf("Order #%f with a total of %f has been completed", order.Id, adminRevenue))

		smtp.SendMail("127.0.0.1:1025", nil, "no-reply@email.com", []string{order.AmbassadorEmail}, adminMessage)

	}(order)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
