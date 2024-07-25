package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ilyaDyb/go_rest_api/config"
	"github.com/ilyaDyb/go_rest_api/models"

	"github.com/gin-gonic/gin"
)

type PersonValidator struct {
	Name       string    `form:"name"`
	Address    string    `form:"address"`
	Birthday   time.Time `form:"birthday" time_format:"2006-01-02"`
	CreateTime time.Time `form:"createTime" time_format:"unixNano"`
	UnixTime   time.Time `form:"unixTime" time_format:"unix"`
}

// receive params
func UserNameController(c *gin.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, "Hello %s", name)
}

func UserNameActionController(c *gin.Context) {
	name := c.Param("name")
	action := c.Param("action")
	message := name + " is " + action
	c.String(http.StatusOK, message)
}

// receive params and default if param is none
func WelcomeController(c *gin.Context) {
	firstname := c.DefaultQuery("firstname", "Guest")
	lastname := c.Query("lastname")
	c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
}

// receive form data with default value
func FormPostController(c *gin.Context) {
	message := c.PostForm("message")
	nick := c.DefaultPostForm("nick", "anonymous")

	//example how to return json response
	c.JSON(http.StatusOK, gin.H{
		"status":  "posted",
		"message": message,
		"nick":    nick,
	})
}

// receive params (def or not) and form's data
func QueryFormPostController(c *gin.Context) {
	id := c.Query("id")
	page := c.DefaultQuery("page", "0")
	name := c.PostForm("name")
	message := c.PostForm("message")

	c.String(http.StatusOK, "id: %s; page: %s; name: %s; message: %s", id, page, name, message)
}

// uploading file
func UploadFile(c *gin.Context) {
	// curl -X POST http://localhost:8080/upload   -F "file=@workdir/testcases.txt"   -H "Content-Type: multipart/form-data"
	file, err := c.FormFile("file")
	if err != nil {
		msg := fmt.Sprintf("Error retrueving the file: %s", err)
		log.Println(msg)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    msg,
		})
	}
	// log.Println(file.Filename, time.Now())
	// dst := fmt.Sprint("users/", config.UploadPath)
	dst := fmt.Sprintf("%s%s", config.DefaultUploadPath, file.Filename)
	log.Println(dst, time.Now())
	// log.Println(dst)
	c.SaveUploadedFile(file, dst)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    fmt.Sprintf("'%s' uploaded successfully", file.Filename),
	})
}

func TestingValidate(c *gin.Context)  {
	// curl -X GET "localhost:8080/testing?name=appleboy&address=xyz&birthday=1992-03-15&createTime=1562400033000000123&unixTime=1562400033"
	var person PersonValidator
	if c.ShouldBind(&person) == nil {
		log.Println(person.Name)
		log.Println(person.Address)
		log.Println(person.Birthday)
		log.Println(person.CreateTime)
		log.Println(person.UnixTime)
	}
	c.String(http.StatusOK, "Success")
}

func TestQueries(c *gin.Context) {
	// all photo
	// var allPhoto []models.Photo
	// config.DB.Find(&allPhoto)
	// c.JSON(http.StatusOK, allPhoto)

	//all photo for special user
	// var allPhotoForSpecialUser []models.Photo
	// var user models.User
	// config.DB.Where("username = ?", "wicki").First(&user)
	// UserID := user.Id
	// config.DB.Where("user_id = ?", UserID).Find(&allPhotoForSpecialUser)
	// c.JSON(http.StatusOK, allPhotoForSpecialUser)

	// All users with their photos
	var allUsers []models.User
	config.DB.Preload("Photo").Find(&allUsers)
	config.DB.Find(&allUsers)
	// for _, usr := range allUsers {
	// 	usr.Sex = "male"
	// 	config.DB.Save(&usr)
	// }
	// config.DB.Save(&allUsers)
	c.JSON(http.StatusOK, allUsers)

	//All interaction
	// var usr models.User
	// config.DB.Where("username = ?", "wicki").First(&usr)
	// usr.Sex = "female"
	// config.DB.Save(&usr)
	// var interactions []models.UserInteraction
	// config.DB.Model(models.UserInteraction{}).Where("user_id = ?", usr.Id).Find(&interactions)
	// c.JSON(http.StatusOK, gin.H{"count": len(interactions), "interactions": interactions})
	// c.JSON(http.StatusOK, usr)
	return
}