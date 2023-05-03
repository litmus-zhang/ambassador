package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Model
	FirstName    string `json:"firstName"`
	Lastname     string `json:"lastname"`
	Email        string `json:"email" gorm:"unique"`
	Password     string `json:"-"`
	IsAmbassador bool   `json:"-"`
}

func (user *User) SetPassword(password string) {
	hashpassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.Password = string(hashpassword)

}

func (user User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

}
