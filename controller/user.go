package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/models"
	_ "github.com/ilyaDyb/go_rest_api/utils"
)

// @Summary  User profile
// @Accept   json
// @Produce  json
// @Param Authorization header string true "With the Bearer started"
// @Success  200 		{object}  utils.MessageResponse
// @Router   /u/profile [get]
func ProfileController(c *gin.Context) {
	username := c.MustGet("username").(string)
	var user models.User
	if err := config.DB.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)

}

// func EditProfileController(c *gin.Context) {

// }
