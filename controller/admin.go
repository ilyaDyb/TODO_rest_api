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
	userService service.UserService
}

func NewAdminController(userService service.UserService) *AdminController {
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
            Email     string     `form:"email"`
            Firstname string     `form:"firstname"`
            Lastname  string     `form:"lastname"`
            Age       string     `form:"age"`
            Country   string     `form:"country"`
            City      string     `form:"city"`
            Bio       string     `form:"bio"`
            Hobbies   string     `form:"hobbies"`
        }
        var input ChangeProfileInput
        err := c.ShouldBindBodyWithJSON(&input)
        if err != nil {
            log.Println(err.Error())
            return
        }
        log.Println(input)
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
        log.Println(user)
        if err := ctrl.userService.UpdateUser(user); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update user profile"})
            return
        }
        c.Status(http.StatusNoContent)
    case "POST":
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
