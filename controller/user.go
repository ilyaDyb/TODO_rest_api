package controller

import (
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/models"
	"github.com/ilyaDyb/go_rest_api/utils"
	_ "github.com/ilyaDyb/go_rest_api/utils"
)

// @Summary  User profile
// @Accept   json
// @Produce  json
// @Param Authorization header string true "With the Bearer started"
// @Param username path string false "Username"
// @Success  200 {object} utils.ModelResponse
// @Failure  403 {object} utils.ErrorResponse
// @Failure  404 {object} utils.ErrorResponse
// @Router   /u/profile [get]
// @Router   /u/profile/{username} [get]
func ProfileController(c *gin.Context) {
	username := c.Param("username")

	username = username[1:]
	if username == "" || username == "/" {
		usernameFromToken, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden!"})
			return
		}
		username = usernameFromToken.(string)
	}
	var user models.User
	if err := config.DB.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

type ChangeProfileInput struct {
	Firstname string                `json:"firstname"`
	Lastname  string                `json:"lastname"`
	Age       uint8                 `json:"age" binding:"min=18,max=99"`
	Country   string                `json:"country"`
	Bio       string                `json:"bio"`
	Hobbies   string                `json:"hobbies"`
	Photo     *multipart.FileHeader `json:"photo"`
}

// EditProfileController edits user profile
// @Summary Edit user profile
// @Description Edit user profile details including uploading a profile photo
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "With the Bearer started"
// @Param firstname formData string false "First Name"
// @Param lastname formData string false "Last Name"
// @Param age formData uint8 false "Age"
// @Param country formData string false "Country"
// @Param bio formData string false "Bio"
// @Param hobbies formData string false "Hobbies"
// @Param photo formData file false "Profile Photo"
// @Success 202 {object} utils.MessageResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /u/profile [put]
func EditProfileController(c *gin.Context) {
	currentUsername := c.MustGet("username")
	file, err := c.FormFile("photo")
	if err != nil {
		if err.Error() != "http: no such file" {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
	}
	if !utils.IsValidPhotoExt(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid extension"})
		return
	}

	filename := filepath.Base(file.Filename)
	filepath := filepath.Join(config.UserPhotoPath, filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input ChangeProfileInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Table("users").Where("username = ?", currentUsername).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Firstname = input.Firstname
	user.Lastname = input.Lastname
	user.Age = input.Age
	user.Country = input.Country
	user.Bio = input.Bio
	user.Hobbies = input.Hobbies

	photo := models.Photo{
		UserID: uint(user.ID),
		URL:    filepath,
	}
	if err := config.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Profile updated successfully"})
}
