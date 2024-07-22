package utils

import (
	"strings"
	"unicode"

	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/models"
)

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	} 
	hasDigit := false
	for _, r := range password {
		if unicode.IsDigit(r) {
			hasDigit = true
			break
		}
	}
	return hasDigit
}

func IsValidUsernameEmail(username string, email string) bool {
	var existingUser models.User
	if err := config.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return false
	}
	if err := config.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return false
	}
	return true
}

func IsValidPhotoExt(filename string) bool {
	validExtensions := []string{"png", "jpg", "jpeg", "webp"}
	extension := strings.Split(filename, ".")[1]
	valid := false
	for _, ext := range(validExtensions) {
		if ext == extension {
			valid = true
			break
		}
	}
	return valid
}