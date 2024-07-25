package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/models"
)

func GetProfiles(c *gin.Context) {
	// username := c.MustGet("username").(string)
	// pageStr := c.DefaultQuery("page", "1")
	// limit := 10
	// page, err := strconv.Atoi(pageStr)
	// if err != nil || page < 1 {
	// 	page = 1
	// }
	// var user models.User
	// if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
	// }

	// var interactionCount int64
	// timeAllow := time.Now().Add(24 * time.Hour)
	// TODO
	// В другом контроллере добавлять лайк или дизлайк а тут проверять если кол-во
	// в процентном делении на 10 == 0, то ставим дату когда можно было еще свайпать
	// и сразу тут же проверять чтобы дата была валидная только тогда возвращаем пользователей
	// Есть вариант добавить подписку, чтобы можно было листать до 50 пользователей

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
	if countOfInteraction%10 == 0 && countOfInteraction != 0 {
		RestrictionEnd := time.Now().Add(24 * time.Hour)
		curUser.RestrictionEnd = RestrictionEnd
		if err := config.DB.Save(&curUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.Status(http.StatusOK)
	// Need to test
}
