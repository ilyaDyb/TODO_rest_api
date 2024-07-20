package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, gin.H{"username": username})
}
