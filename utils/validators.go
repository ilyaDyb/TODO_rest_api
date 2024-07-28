package utils

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator"
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

func IsValidEmailFormat(email string) bool {
	if len(email) <= 5 {
		return false
	}
	if strings.Count(email, "@") != 1 {
		return false
	}
	if !strings.Contains(strings.Split(email, "@")[1], ".") {
		return false
	}
	return true
}

func ValidateStruct(data interface{}) error {
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		return err
	}
	return nil
}