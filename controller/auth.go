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
	Username  string `json:"username" binding:"required" validate:"max=50"`
	Email     string `json:"email" binding:"required,email" validate:"max=100"`
	Password  string `json:"password" binding:"required" validate:"min=8,max=100"`
	Firstname string `json:"firstname" binding:"required" validate:"max=50"`
	Lastname  string `json:"lastname" binding:"required" validate:"max=50"`
	Sex       string `json:"sex" binding:"required" validate:"oneof=male female"`
	Age       uint8  `json:"age" binding:"required" validate:"min=18,max=99"`
	Country   string `json:"country" binding:"required" validate:"max=50"`
	City      string `json:"city" binding:"required" validate:"max=50"`
	Hobbies   string `json:"hobbies" validate:"max=100"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required" validate:"max=50"`
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
	err := utils.ValidateStruct(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE username = ? OR email = ?)" 
	var exists bool
	if err := config.DB.Raw(query, input.Username, input.Email).Scan(&exists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user with such username or email already exists"})
		return
	}
	if !utils.IsValidPassword(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be >= 8 chars long and contain number"})
		return
	}
	if !utils.IsValidEmailFormat(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email format is invalid"})
		return
	}
	user := models.User{
		Username: input.Username, Role: "user", Email: input.Email, Sex: input.Sex,
		Age: input.Age, Country: input.Country, City: input.City, Hobbies: input.Hobbies,
		Firstname: input.Firstname, Lastname: input.Lastname,
	}	
	
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
	err := utils.ValidateStruct(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
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

	refreshToken, err := utils.GenerateRefreshToken(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": refreshToken})
}

type InputRefresh struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
// @Summary      Refreshing access Token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        InputRefresh  body      InputRefresh  true  "InputRefresh"
// @Success      200         {object}  MessageResponse
// @Router       /auth/refresh [post]
func Refresh(c *gin.Context) {
	var input InputRefresh
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := utils.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	newToken, err := utils.GenerateJWT(claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}