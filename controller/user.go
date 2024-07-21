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
// @Param username path string false "Username"
// @Success  200 		{object}  utils.ModelResponse
// @Router   /u/profile/{username} [get]
func ProfileController(c *gin.Context) {
	username := c.Param("username")
	var user models.User
	if err := config.DB.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)

}

// type ChangeProfileInput struct {
// 	gorm.Model
// 	ID       int      `json:"id" gorm:"primary_key"`
// 	Username string   `json:"name" gorm:"unique"`
// 	Email    string   `json:"email"`
// 	Password string   `json:"-"`
// 	Sex      string   `json:"sex"`
// 	Age      uint8    `json:"age"`
// 	Country  string   `json:"country"`
// 	Location string   `json:"location"`
// 	Role     string   `json:"role"`
// 	Bio      string   `json:"bio"`
// 	Hobbies  string `json:"hobbies"`
// 	Photos   []Photo  `json:"photos" gorm:"foreignKey:UserID"`
// }

func EditProfileController(c *gin.Context) {
	currentUsername := c.MustGet("username")
	pathUsername := c.Param("username")
	if currentUsername != pathUsername {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden!"})
		return
	}

}
