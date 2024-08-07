package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/logger"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/ilyaDyb/go_rest_api/service"
	"github.com/ilyaDyb/go_rest_api/tasks"
	"github.com/ilyaDyb/go_rest_api/utils"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	userService service.UserService
}

func NewAuthController(userService service.UserService) *AuthController {
	return &AuthController{userService: userService}
}

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
func (ctrl *AuthController) RegistrationController(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Errorf("client send invalid data with error: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := utils.ValidateStruct(input)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Infof("user send invalid data for registration with error: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	exists, err := ctrl.userService.UserIsExists(input.Username, input.Email)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "gorm",
		}).Errorf("database server could not do IsExists operation with error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Infof("user want to register by username and email which already exits -> username: %v, email: %v", input.Username, input.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user with such username or email already exists"})
		return
	}
	if !utils.IsValidPassword(input.Password) {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Infof("user tried to set invalid password: %v", input.Password)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be >= 8 chars long and contain number"})
		return
	}
	if !utils.IsValidEmailFormat(input.Email) {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Infof("user tried to set invalid email: %v", input.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email format is invalid"})
		return
	}
	confirmationHash := utils.GetMD5Hash(utils.RandStringRunes(20))
	user := models.User{
		Username:         input.Username,
		Role:             "user",
		Email:            input.Email,
		Sex:              input.Sex,
		Age:              input.Age,
		Country:          input.Country,
		City:             input.City,
		Hobbies:          input.Hobbies,
		Firstname:        input.Firstname,
		Lastname:         input.Lastname,
		ConfirmationHash: confirmationHash,
	}

	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := []byte(fmt.Sprintf("To: recipient@example.net\r\n"+
		"Subject: Tinder-clone!\r\n"+
		"\r\n"+
		"Your link for confirming email %s%s/auth/confirm?hash=%s.\r\n", config.ServerProtocol, config.ServerHost, confirmationHash))
	task, err := tasks.NewEmailDeliveryTask(input.Email, msg)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "asynq",
		}).Errorf("server could not start email delivery task with error: %v", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	info, err := config.Client.Enqueue(task)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "asynq",
		}).Errorf("server could not start email delivery with error: %v", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	if err := ctrl.userService.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Confirm email"})
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
func (ctrl *AuthController) LoginController(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Errorf("client send bad data with error: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := utils.ValidateStruct(input)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Errorf("client send bad data: %v, with error: %v", input, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.userService.GetUserByUsername(input.Username)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "gorm",
		}).Errorf("databse service could not find user by username: %v, with err: %v", input.Username, err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := user.CheckPassword(input.Password); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "user",
		}).Infof("user: %v tried to enter invalid password: %v, with err: %v", input.Username, input.Password, err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, err := utils.GenerateJWT(input.Username)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "JWT",
			"username":  input.Username,
		}).Errorf("JWT service could not generate token by username: %v", input.Username)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(input.Username)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "JWT",
			"username":  input.Username,
		}).Errorf("JWT service could not generate refresh token by username: %v", input.Username)
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
func (ctrl *AuthController) RefreshController(c *gin.Context) {
	var input InputRefresh
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims, err := utils.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "JWT",
		}).Errorf("JWT service could not parse refresh token: %v", input.RefreshToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	newToken, err := utils.GenerateJWT(claims.Username)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "JWT",
			"username":  claims.Username,
		}).Errorf("JWT service could not generate token by username: %v", claims.Username)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

// @Summary Change Password
// @Description Change user password using verification code
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   changePasswordInput body ChangePasswordInput true "Change Password Input"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/change-password [post]
func (ctrl *AuthController) ConfirmEmailController(c *gin.Context) {
	hash := c.Query("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hash is required"})
		return
	}

	user, err := ctrl.userService.GetUserByHash(hash)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "user",
			"service":   "gorm",
		}).Errorf("databse service could not find user by hash: %v, with err: %v", hash, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user with this hash"})
		return
	}

	user.IsActive = true
	user.ConfirmationHash = ""
	if err := ctrl.userService.UpdateUser(user); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "user",
			"service":   "gorm",
		}).Errorf("databse service could not update user: with err: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email was confirmed"})
}

type ConfirmEmailResponse struct {
	Message string `json:"message"`
}

type EmailInput struct {
	Email string `json:"email"`
}

// @Summary Request Password Reset
// @Description Request password reset via email
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   email body EmailInput true "Email"
// @Success 202 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/drop-password [post]
func (ctrl *AuthController) DropPasswordController(c *gin.Context) {
	var input EmailInput
	if err := c.ShouldBind(&input); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Errorf("client send bad data with error: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := input.Email
	if !utils.IsValidEmailFormat(email) {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Debug("client send invalid email with error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}
	existsEmail, err := ctrl.userService.IsExistsEmail(email)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "user",
			"service":   "gorm",
		}).Errorf("databse service could not find user by email: %v, with error: %v", email, err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	if !existsEmail {
		c.Status(http.StatusBadRequest)
		return
	}
	code, err := utils.GenerateRandomCode(6)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "user",
		}).Errorf("server could not generate code with error: %v", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	msg := []byte(fmt.Sprintf("To: recipient@example.net\r\n"+
		"Subject: Tinder-clone!\r\n"+
		"\r\n"+
		"Your code for confirming email %v.\r\n", code))
	task, err := tasks.NewEmailDeliveryTask(input.Email, msg)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "asynq",
		}).Errorf("server could not start email delivery with error: %v", err.Error())
		log.Println(task)
		c.Status(http.StatusInternalServerError)
		return
	}
	info, err := config.Client.Enqueue(task)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "asynq",
		}).Errorf("server could not start email delivery with error: %v", err.Error())
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	taskInfo := fmt.Sprintf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	logger.Log.WithFields(logrus.Fields{
		"component": "auth",
		"service":   "asynq",
	}).Infof("email delivery send mail successfully INF: %v", taskInfo)

	user, err := ctrl.userService.GetUserByEmail(email)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "gorm",
		}).Errorf("server could not find user by email with error: %v", err.Error())
		log.Printf("Error getting user by email: %v\n", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	ID := strconv.Itoa(int(user.ID))

	err = utils.SetCache(config.RedisClient, code, ID, time.Minute*3)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "redis",
		}).Errorf("error setting cache with error: %v", err.Error())
		log.Printf("Error setting cache: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	logger.Log.WithFields(logrus.Fields{
		"component": "auth",
	}).Infof("email sent successfully with code: %v to: %v, with id: %v", code, email, ID)
	c.JSON(http.StatusAccepted, gin.H{"message": "check your email"})
}

type ChangePasswordInput struct {
	Code        string `json:"code"`
	NewPassword string `json:"password"`
}

// @Summary Change Password
// @Description Change user password using verification code
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   changePasswordInput body ChangePasswordInput true "Change Password Input"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/change-password [post]
func (ctrl *AuthController) ChangePassword(c *gin.Context) {
	var input ChangePasswordInput
	err := c.ShouldBind(&input)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Errorf("client sent bad data: %v, with error: %v", input, err.Error())
		log.Printf("Error binding input: %v\n", err)
		c.Status(http.StatusBadRequest)
		return
	}
	log.Printf("Received code: %s, new password: %s\n", input.Code, input.NewPassword)

	userID, err := utils.GetCache(config.RedisClient, input.Code)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"service":   "redis",
		}).Errorf("getting cache was failed by key: %v or 3 minutes have passed since the code was sent to the email, with error: %v", input.Code, err.Error())
		c.JSON(http.StatusConflict, gin.H{"error": "internal server error or 3 minutes have passed since the code was sent to the email"})
		return
	}

	var user models.User
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).First(&user).Error; err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Errorf("client sent bad data: %v, with error: %v", input, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	if err := user.CheckPassword(input.NewPassword); err == nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
		}).Info("New password is the same as the old password")
		log.Println("New password is the same as the old password")
		c.JSON(http.StatusBadRequest, gin.H{"error": "the new password must be different from the old password"})
		return
	}

	if err := user.HashPassword(input.NewPassword); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"user_id":   userID,
		}).Errorf("Failed to hash password: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := config.DB.Save(&user).Error; err != nil {
		logger.Log.WithFields(logrus.Fields{
			"component": "auth",
			"user_id":   userID,
		}).Errorf("Failed to save user: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"component": "auth",
		"user_id":   userID,
	}).Info("User successfully changed password")

	c.JSON(http.StatusOK, gin.H{"message": "the password was updated successfully"})
}
