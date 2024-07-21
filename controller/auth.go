package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/ilyaDyb/go_rest_api/utils"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type MessageResponse = utils.MessageResponse
type ErrorResponse = utils.ErrorResponse

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user by providing a username and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        RegisterInput  body      RegisterInput  true  "Register Input"
// @Success      200            {object}  MessageResponse
// @Failure      400            {object}  ErrorResponse
// @Failure      500            {object}  ErrorResponse
// @Router       /auth/registration [post]
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Error when retrieving data")
		return
	}
	var existingUser models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
		return
	}
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
		return
	}
	user := models.User{Username: input.Username, Role: "user", Email: input.Email}
	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("Error when hashing password")
		return
	}
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		log.Println("Error when creating user")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registration success"})
}

// Login godoc
// @Summary      Login a user
// @Description  Login a user by providing a username and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginInput  body      LoginInput  true  "Login Input"
// @Success      200         {object}  MessageResponse
// @Failure      400         {object}  ErrorResponse
// @Failure      401         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	result := config.DB.Where("username = ?", input.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := utils.GenerateJWT(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
