package models

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)
type User struct {
	gorm.Model
	ID		 int	 `json:"id" gorm:"primary_key"`
	Username     string  `json:"name" gorm:"unique"`
	// Email    string  `json:"email"`
	Password string  `json:"password"`
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