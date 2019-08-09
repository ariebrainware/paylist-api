package endpoint

import (
	"fmt"
	"net/http"
	"strconv"
	
	"github.com/ariebrainware/paylist-api/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/ariebrainware/dgrijalva/jwt-go"  //Used to sign and verify JWT tokens
)

type Token struct {
	ID uint
	Email string `json:"email"`
	Name string `json:"name"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	jwt.StandardClaims
}

// CreateUser function to sign up
func CreateUser(c *gin.Context) {
	users := model.User{
		Email	: c.PostForm("email"),
		Name	: c.PostForm("name"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	fmt.Println(c.PostForm("email"))
	fmt.Println(c.PostForm("name"))
	fmt.Println(c.PostForm("username"))
	fmt.Println(c.PostForm("password"))

	//Password Encryption
	password, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotImplemented, gin.H{
			"message": "password encryption failed",
		})
		c.JSON(http.StatusCreated, gin.H{
			"message": "password encryption success!",
		})
	}
	users.Password = string(password)
	db.Save(&users)
	c.JSON(http.StatusCreated, gin.H{
		"status":      http.StatusCreated,
		"message":     "User created Successfully!",
		"resourcedId": users.ID,
	})
}

// FetchUser function to get list of users
func FetchUser(c *gin.Context) {
	var users []model.User

	db.Find(&users)

	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": users})
}

// UpdateUser function to update user information
func UpdateUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	updateuser := model.User{
		Email	: c.PostForm("email"),
		Name	: c.PostForm("name"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	err := db.Model(&model.User{}).Where("ID = ?", id).Update(&updateuser).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusNotFound,
			"message": "User update unsuccessfully!",
			"error":   err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "User update successfully!",
		"errors":  err,
	})
}

// DeleteUser function to handle user deletion
func DeleteUser(c *gin.Context) {
	var users model.User
	usersID := c.Param("id")

	db.First(&users, usersID)

	if users.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No user found!"})
		return
	}
	db.Delete(&users)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK, "message": "User Delete succcessfully!"})
}

// FetchSingleUser function to get single user
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

// Login function to handle login user
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user := &model.User{}
	if username == ""{
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "please provide username and password"})
		return
	}

	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "wrong username"})
		return
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
	 	c.JSON(http.StatusNotFound, gin.H{"status": false, "message": "wrong password or password doesn't match"})
	 	return
	 }
	 
	 tk := &Token{
		ID: user.ID,
		Email: user.Email,
		Name: user.Name,
		Username:  user.Username,
	}
	//Create JWT token 
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "logged in",
		"token": tokenString,
		"user": user,
	})
}

//FuncAuth Function Authorization to handle authorized
func Auth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected SigningMethod :%v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if token != nil && err == nil {
		fmt.Println("token verified")
	} else {
		result :=gin.H {
			"message": "not authorized",
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
	}
}