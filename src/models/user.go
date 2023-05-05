package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Model
	FirstName    string   `json:"firstName"`
	Lastname     string   `json:"lastname"`
	Email        string   `json:"email" gorm:"unique"`
	Password     string   `json:"-"`
	IsAmbassador bool     `json:"-"`
	Revenue      *float64 `json:"revenue,omitempty" gorm:"-"`
}

func (user *User) SetPassword(password string) {
	hashpassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.Password = string(hashpassword)

}

func (user User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

}

type Admin User

func (admin *Admin) CalculateRevenue(db *gorm.DB) {
	var orders []Order
	db.Preload("OrderItems").Find(&orders, &Order{
		UserId:   admin.Id,
		Complete: true,
	})
	var revenue float64 = 0
	for _, order := range orders {
		for _, orderitem := range order.OrderItems {
			revenue += orderitem.AdminRevenue
		}
	}
	admin.Revenue = &revenue
}

type Ambassador User

func (ambassador *Ambassador) CalculateRevenue(db *gorm.DB) {

	var orders []Order
	db.Preload("OrderItems").Find(&orders, &Order{
		UserId:   ambassador.Id,
		Complete: true,
	})
	var revenue float64 = 0
	for _, order := range orders {
		for _, orderitem := range order.OrderItems {
			revenue += orderitem.AmbassadorRevenue
		}
	}
	ambassador.Revenue = &revenue

}
