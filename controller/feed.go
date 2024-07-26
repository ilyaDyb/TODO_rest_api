package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/rosberry/go-pagination"
)

func GetUsersList(userID uint, role string, paginator *pagination.Paginator) []models.User {
	var curUser models.User
	config.DB.First(&curUser, userID)

	var users []models.User

	var interactedIDs []uint
	config.DB.Model(&models.UserInteraction{}).Where("user_id = ?", userID).Pluck("target_id", &interactedIDs)

	ageLower := curUser.Age - 3
	ageUpper := curUser.Age + 100
	gender := "male"
	if curUser.Sex == "male" {
		gender = "female"
	} else {
		gender = "male"
	}

	q := config.DB.Preload("Photo").Model(&models.User{}).
		Where("role = ?", role).
		Where("id != ?", userID).
		Where("age BETWEEN ? and ?", ageLower, ageUpper).
		Where("sex = ?", gender).
		Where("id NOT IN (?)", interactedIDs)

	err := paginator.Find(q, &users)
	if err != nil {
		log.Println(err)
		return nil
	}
	return users
}

type usersListResponse struct {
	Result    bool                 `json:"result"`
	Users      []models.User        `json:"users"`
	Pagination *pagination.PageInfo `json:"pagination"`
}

// @Summary Get profile
// @Accept json
// @Produce json
// @Param Authorization header string true "With the Bearer started"
// @Router /feed/get-profiles [get]
func GetProfiles(c *gin.Context) {
	username := c.MustGet("username").(string)
	var user models.User
	if err := config.DB.Model(models.User{}).Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	userID := user.Id
	paginator, err := pagination.New(pagination.Options{
		GinContext:    c,
		DB:            config.DB,
		Model:         &models.User{},
		Limit:         2,
		CustomRequest: &pagination.RequestOptions{
			Cursor: func(c *gin.Context) (query string) {
				return c.Query("cursor")
			},
			After: func(c *gin.Context) (query string) {
				return c.Query("after")
			},
			Before: func(c *gin.Context) (query string) {
				return c.Query("before")
			},
		},
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	users := GetUsersList(userID, "user", paginator)

	c.JSON(http.StatusOK, usersListResponse{
		Result: true,
		Users: users,
		Pagination: paginator.PageInfo,
	})
}

type InputGrade struct {
	TargetID  uint
	InterType string
}

// @Summary to Grade profiles
// @Accept json
// @Produce json
// @Param Authorization header string true "With the Bearer started"
// @Param InputGrade body InputGrade true "Input for Grade other profile"
// @Router /feed/grade [post]
func GradeProfile(c *gin.Context) {
	username := c.MustGet("username").(string)
	var curUser models.User
	if err := config.DB.Table("users").Where("username = ?", username).First(&curUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	if curUser.RestrictionEnd.After(time.Now()) && !(curUser.RestrictionEnd.IsZero()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("you have restriction for interaction expire at %02d-%02d", curUser.RestrictionEnd.Month(), curUser.RestrictionEnd.Day())})
		return
	}
	var input InputGrade
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	InterType := input.InterType
	if InterType != "like" && InterType != "dislike" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Interaction should be 'like' or 'dislike'"})
		return
	}
	targetId := input.TargetID
	var interaction models.UserInteraction
	interaction.TargetID = targetId
	interaction.UserID = curUser.Id
	interaction.InteractionType = InterType
	if err := config.DB.Create(&interaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	var countOfInteraction int64
	config.DB.Model(&models.UserInteraction{}).Where("user_id = ?", curUser.Id).Count(&countOfInteraction)
	// Check if the user is subscribed if not
	// if subscriber {limit = 100} else {limit = 10} countOfInteraction%limit == 0 && != 0
	if countOfInteraction%10 == 0 && countOfInteraction != 0 {
		RestrictionEnd := time.Now().Add(24 * time.Hour)
		curUser.RestrictionEnd = RestrictionEnd
		if err := config.DB.Save(&curUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.Status(http.StatusOK)
}
