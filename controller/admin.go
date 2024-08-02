package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/ilyaDyb/go_rest_api/service"
	"github.com/ilyaDyb/go_rest_api/utils"
)

type AdminController struct {
	userService *service.UserService
}

func NewAdminController(userService *service.UserService) *AdminController {
    return &AdminController{userService: userService}
}

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
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }
    total, err := ctrl.userService.GetUsersCount()
    if err != nil {
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


func (ctrl *AdminController) GetPutPostDeleteUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    method := c.Request.Method
    if err != nil && method != "POST"{
        c.AbortWithStatus(http.StatusBadRequest)
        return
    }

    var user *models.User
    if method != "POST" {
        user, err = ctrl.userService.GetUserByID(uint(id))
        if err != nil {
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
    }
    switch method {
    case "GET":
        c.JSON(http.StatusOK, user)
    case "DELETE":
        if err := ctrl.userService.DeleteUser(user); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.Status(http.StatusNoContent)
    case "PUT":
        type ChangeProfileInput struct {
            Firstname string                `form:"firstname" validate:"max=20"`
            Lastname  string                `form:"lastname" validate:"max=20"`
            Age       uint8                 `form:"age" validate:"min=18,max=99"`
            Country   string                `form:"country" validate:"max=30"`
            City      string                `form:"city" validate:"max=30"`
            Bio       string                `form:"bio" validate:"max=500"`
            Hobbies   string                `form:"hobbies" validate:"max=100"`
        }
        var input ChangeProfileInput
        user.Firstname = input.Firstname
        user.Lastname = input.Lastname
        user.Age = input.Age
        user.Country = input.Country
        user.City = input.City
        user.Bio = input.Bio
        user.Hobbies = input.Hobbies

        if err := ctrl.userService.UpdateUser(user); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update user profile"})
            return
        }
        c.Status(http.StatusNoContent)
    case "POST":
        // TODO
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
        err := c.ShouldBindBodyWithJSON(&input)
        if err != nil {
            log.Println("Error binding JSON:", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        err = utils.ValidateStruct(input)
        if err != nil {
            log.Println("Validation error:", err)
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
            log.Println("Error hashing password:", err)
            c.Status(http.StatusInternalServerError)
            return
        }
        
        err = ctrl.userService.CreateUser(&user)
        if err != nil {
            log.Println("Error creating user:", err)
            c.Status(http.StatusInternalServerError)
            return
        }
        
        c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})        
    }
}
