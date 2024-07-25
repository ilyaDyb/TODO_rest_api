package config

import (
	"log"
	"os"
	"strings"

	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

func Connect() {
    DB_HOST := "host=" + os.Getenv("DB_HOST")
    DB_USER := "user=" + os.Getenv("DB_USER")
    DB_PASSWORD := "password=" + os.Getenv("DB_PASSWORD")
    DB_NAME := "dbname=" + os.Getenv("DB_NAME")
    DB_PORT := "port=" + os.Getenv("DB_PORT")
    params := []string{DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT}
    dsn := strings.Join(params, " ")
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database: " + err.Error())
    }
    DB.AutoMigrate(&models.User{})
    DB.AutoMigrate(&models.Photo{})
    DB.AutoMigrate(&models.UserInteraction{})
}