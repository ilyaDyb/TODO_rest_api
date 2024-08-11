package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/logger"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/ilyaDyb/go_rest_api/service"
	"github.com/ilyaDyb/go_rest_api/utils"
	"github.com/sirupsen/logrus"
)

type AdminController struct {
	userService service.UserService
    chatService service.ChatService
}

func NewAdminController(userService service.UserService, chatService service.ChatService) *AdminController {
    return &AdminController{userService: userService, chatService: chatService}
}

// UsersList godoc
// @Summary Get a list of users
// @Description Get a list of users with pagination
// @Tags admin
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /admin/users [get]
func (ctrl *AdminController) UsersList(c *gin.Context) {
    limitStr := c.DefaultQuery("limit", "10")
    pageStr := c.DefaultQuery("page", "1")
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        limit = 10
    }
    page, err := strconv.Atoi(pageStr)
    if err != nil {
        page = 1
    }

    users, err := ctrl.userService.GetAllUsers(limit, (page-1)*limit)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("failed to get all users with error: %v", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }
    total, err := ctrl.userService.GetUsersCount()
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("failed to get users count with error: %v", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users count"})
        return
    }
    totalPages := (total + limit - 1) / limit

    c.JSON(http.StatusOK, gin.H{
        "users":      users,
        "page":       page,
        "totalPages": totalPages,
        "total":      total,
    })
}

// GetPutPostDeleteUser godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags admin
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/user/{id} [get]
func (ctrl *AdminController) GetUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("with error: %v", err.Error())
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }

    user, err := ctrl.userService.GetUserByID(uint(id))
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("with error: %v", err.Error())
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Delete a user by ID
// @Tags admin
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/user/{id} [delete]
func (ctrl *AdminController) DeleteUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("with error: %v", err.Error())
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }

    user, err := ctrl.userService.GetUserByID(uint(id))
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("with error: %v", err.Error())
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    if err := ctrl.userService.DeleteUser(user); err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("with error: %v", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}

// UpdateUser godoc
// @Summary Update a user by ID
// @Description Update a user by ID
// @Tags admin
// @Produce json
// @Param id path int true "User ID"
// @Param input body models.User true "User info"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/user/{id} [put]
func (ctrl *AdminController) UpdateUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("with error: %v", err.Error())
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }

    user, err := ctrl.userService.GetUserByID(uint(id))
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("with error: %v", err.Error())
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    type ChangeProfileInput struct {
        Email     string `form:"email"`
        Firstname string `form:"firstname"`
        Lastname  string `form:"lastname"`
        Age       string `form:"age"`
        Country   string `form:"country"`
        City      string `form:"city"`
        Bio       string `form:"bio"`
        Hobbies   string `form:"hobbies"`
    }
    var input ChangeProfileInput
    err = c.ShouldBindJSON(&input)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("with error: %v", err.Error())
        return
    }

    if input.Firstname != "" {
        user.Firstname = input.Firstname
    }
    if input.Email != "" {
        user.Email = input.Email
    }
    if input.Lastname != "" {
        user.Lastname = input.Lastname
    }
    if input.Age != "0" {
        age, err := strconv.Atoi(input.Age)
        if err != nil {
            c.Status(http.StatusBadRequest)
            return
        }
        user.Age = uint8(age)
    }
    if input.Country != "" {
        user.Country = input.Country
    }
    if input.City != "" {
        user.City = input.City
    }
    if input.Bio != "" {
        user.Bio = input.Bio
    }
    if input.Hobbies != "" {
        user.Hobbies = input.Hobbies
    }

    if err := ctrl.userService.UpdateUser(user); err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("with error: %v", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update user profile"})
        return
    }
    c.Status(http.StatusNoContent)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags admin
// @Produce json
// @Param input body models.User true "User info"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/user [post]
func (ctrl *AdminController) CreateUser(c *gin.Context) {
    type createUserInput struct {
        Username  string `json:"username" validate:"max=50"`
        Email     string `json:"email" validate:"max=100"`
        Password  string `json:"password" validate:"min=8,max=100"`
        Firstname string `json:"firstname"`
        Lastname  string `json:"lastname"`
        Sex       string `json:"sex" validate:"oneof=male female"`
        Role      string `json:"role" validate:"oneof=admin user"`
        Age       uint8  `json:"age"`
        Country   string `json:"country"`
        City      string `json:"city"`
        Hobbies   string `json:"hobbies"`
    }

    var input createUserInput
    err := c.ShouldBindJSON(&input)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = utils.ValidateStruct(input)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
        return
    }

    user := models.User{
        Username:  input.Username,
        Email:     input.Email,
        Password:  input.Password,
        Sex:       input.Sex,
        Role:      input.Role,
        Age:       input.Age,
        Country:   input.Country,
        City:      input.City,
        Hobbies:   input.Hobbies,
        Firstname: input.Firstname,
        Lastname:  input.Lastname,
        IsActive:  true,
    }

    err = user.HashPassword(input.Password)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("failed to hash password error: %v", err.Error())
        c.Status(http.StatusInternalServerError)
        return
    }

    err = ctrl.userService.CreateUser(&user)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "component": "admin",
        }).Errorf("failed to create user with error: %v", err.Error())
        c.Status(http.StatusInternalServerError)
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// @Summary Get absolutely all chats
// @Description Route which return all chats
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/chats [get]
func (ctrl *AdminController) GetAllChats(c *gin.Context) {
    chats, _ := ctrl.chatService.GetAllChats()
    c.JSON(200, chats)
}