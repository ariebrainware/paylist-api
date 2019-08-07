package endpoints

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/ariebrainware/paylist-api/model"
	"golang.org/x/crypto/bcrypt"
)

//createUser
func CreateUser(c *gin.Context) {
	users := model.User{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	fmt.Println(c.PostForm("username"))
	fmt.Println(c.PostForm("password"))

	db.Save(&users)
	c.JSON(http.StatusCreated, gin.H{
		"status":	http.StatusCreated,
		"message":	"User created Successfully!",
		"resourcedId": users.ID,
	})
}

//GetUser
func FetchUser(c *gin.Context) {
	var users []model.User

	db.Find(&users)

	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message":"No user found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": users})
}

//UpdateUser
func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	updateuser := model.User{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	err := db.Model(&model.User{}).Where("ID = ?", id).Update(&updateuser).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusNotFound,
			"message": "User update unsuccessfully!",
			"error": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "User update successfully!",
		"errors": err,
	})
}

//deleteuser
func DeleteUser(c *gin.Context) {
	var users model.User
	usersID := c.Param("id")

	db.First(&users, usersID)

	if users.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H {
			"message": "No user found!"})
		return
	}
	db.Delete(&users)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK, "message": "User Delete succcessfully!"})
}

func FetchSingleUser(c *gin.Context) {
	var users model.User
	usersID := c.Param("id")
	err := db.Model(&model.User{}).Where("ID = ?", usersID).Find(&users).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": users})
}

func Login(c *gin.Context) {
	//username := c.PostForm("username")
	password := c.PostForm("password")
	
	user := &model.User{}
	db.Find(&user)
	if user.Username == "" {
		c.JSON(http.StatusNotFound, gin.H{"status":false, "message": "invalid request"})
		return
	}

	err := db.Where("username = ? and password = ?", user.Username, user.Password).First(user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status":false, "message": "wrong username"})
		return 
	}
	
	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		c.JSON(http.StatusNotFound, gin.H{"status":false, "message": "invalid login"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": false, "message": "logged in"})
}