package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id        uint   `json:"id" gorm:"primary_key"`
	Username  string `json:"name" gorm:"unique"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Sex       string `json:"sex"`
	Age       uint8  `json:"age"`
	Country   string `json:"country"`
	Location  string `json:"location"`
	Role      string `json:"role"`
	Bio       string `json:"bio"`
	Hobbies   string `json:"hobbies"`
	Photo     Photo  `json:"photo" gorm:"foreignKey:UserID"`
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type Photo struct {
	gorm.Model
	UserID uint   `gorm:"index"`
	URL    string `json:"url"`
}
