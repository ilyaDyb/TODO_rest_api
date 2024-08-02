package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// Id             uint      `json:"id" gorm:"primary_key"`
	Username         string    `json:"username" gorm:"unique"`
	Email            string    `json:"email"`
	Password         string    `json:"-"`
	Firstname        string    `json:"firstname"`
	Lastname         string    `json:"lastname"`
	Sex              string    `json:"sex"`
	Age              uint8     `json:"age"`
	Country          string    `json:"country"`
	City             string    `json:"city"`
	Lat              float32   `json:"lat"`
	Lon              float32   `json:"lon"`
	Role             string    `json:"role"`
	Bio              string    `json:"bio"`
	Hobbies          string    `json:"hobbies"`
	Photo            []Photo   `json:"photo" gorm:"foreignKey:UserID"`
	RestrictionEnd   time.Time `json:"restriction_end"`
	IsActive         bool      `json:"is_active" gorm:"default:false"`
	ConfirmationHash string    `json:"-"`
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
	UserID    uint   `json:"user_id" gorm:"index"`
	URL       string `json:"url"`
	IsPreview bool   `json:"is_preview"`
}

type UserInteraction struct {
	gorm.Model
	UserID          uint   `json:"user_id"`
	TargetID        uint   `json:"target_id"`
	InteractionType string `json:"interaction_type"`
}

type TemporaryUser struct {
	gorm.Model
	// Id             uint      `json:"id" gorm:"primary_key"`
	Username       string    `json:"name" gorm:"unique"`
	Email          string    `json:"email"`
	Password       string    `json:"-"`
	Firstname      string    `json:"firstname"`
	Lastname       string    `json:"lastname"`
	Sex            string    `json:"sex"`
	Age            uint8     `json:"age"`
	Country        string    `json:"country"`
	City           string    `json:"city"`
	Role           string    `json:"role"`
	Hobbies        string    `json:"hobbies"`
	RestrictionEnd time.Time `json:"restriction_end"`
}

func (tempUser *TemporaryUser) ConvertToUser() error {
	return nil
}
