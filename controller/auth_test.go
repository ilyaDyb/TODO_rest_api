package controller_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/ilyaDyb/go_rest_api/config"
// 	"github.com/ilyaDyb/go_rest_api/controller"
// 	"github.com/ilyaDyb/go_rest_api/models"
// 	"github.com/ilyaDyb/go_rest_api/service"
// 	"github.com/ilyaDyb/go_rest_api/utils"
// )

// func setupTestDB() {
// 	config.ConnectTestDB()
// 	// Заполняем тестовую базу данных начальными данными
// 	config.DB.Create(&models.User{
// 		Email:    "test@example.com",
// 		IsActive: false,
// 		Password: "oldpasswordhash",
// 	})
// }

// func teardownTestDB() {
// 	db, _ := config.DB.DB()
// 	db.Close()
// }

// func TestDropPasswordController(t *testing.T) {
// 	setupTestDB()
// 	defer teardownTestDB()

// 	// Set up Gin router
// 	router := gin.Default()
// 	userService := &mockUserService{}
// 	authController := controller.NewAuthController(userService)
// 	router.POST("/drop-password", authController.DropPasswordController)

// 	// Create test request
// 	emailInput := controllers.EmailInput{
// 		Email: "test@example.com",
// 	}
// 	body, _ := json.Marshal(emailInput)
// 	req, _ := http.NewRequest("POST", "/drop-password", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	// Record the response
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	// Assert the response
// 	assert.Equal(t, http.StatusAccepted, w.Code)
// 	var response map[string]string
// 	json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.Equal(t, "check your email", response["message"])
// }

// // Mock user service
// type mockUserService struct{}

// func (m *mockUserService) IsExistsEmail(email string) (bool, error) {
// 	if email == "test@example.com" {
// 		return true, nil
// 	}
// 	return false, nil
// }

// func (m *mockUserService) GetUserByEmail(email string) (models.User, error) {
// 	return models.User{ID: 1, Email: email}, nil
// }
