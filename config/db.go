package config

import (
	"log"
	"os"
	"strings"

	"github.com/ilyaDyb/go_rest_api/logger"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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
        logger.Log.WithFields(logrus.Fields{
			"service": "postgres",
		}).Fatalf("could not run postgreSQL sever: %v", err)
        panic("failed to connect database: " + err.Error())
    }
    logger.Log.WithFields(logrus.Fields{
        "service": "postgres",
    }).Info("Postgres was started successfully")
    DB.AutoMigrate(
        &models.User{},
        &models.Photo{},
        &models.UserInteraction{},
        &models.Chat{},
        &models.Message{},
    )
}

func ConnectTestDB()  {
    var err error
    dsn := "file:test.db?mode=memory&cache=shared"
    DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("could not connect to the test database: %v", err.Error())
    }
    log.Println("Test database connected successfully")
    DB.AutoMigrate(&models.User{})
}