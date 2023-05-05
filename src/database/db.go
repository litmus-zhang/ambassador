package database

import (
	"ambassador-app/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func Connect() {
	var err error
	DBConn, err = gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:1000)/ambassador"), &gorm.Config{})

	if err != nil {
		panic("Could not connect with database")
	}
}

func AutoMigrate() {

	result := DBConn.AutoMigrate(models.User{},
		models.Product{},
		models.Link{},
		models.Order{},
		models.OrderItem{})
	if result != nil {
		panic(result.Error())
	}

}
