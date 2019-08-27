package endpoint

import (
	//"github.com/jinzhu/gorm"
	"fmt"
	"net/http"
	"time"
	"strconv"
	//u "io/ioutils"

	"github.com/ariebrainware/paylist-api/model"
	jwt "github.com/dgrijalva/jwt-go" //Used to sign and verify JWT tokens
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/ariebrainware/paylist-api/configdb"
	
)

// Token is a struct for token model
type Token struct {
	Username string
	jwt.StandardClaims
}

type user1 struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	Balance int `json:"balance"`
}

// CreateUser function to sign up
func CreateUser(c *gin.Context) {
	balance, _ := strconv.Atoi(c.PostForm("balance"))
	users := model.User{
		Email:    c.PostForm("email"),
		Name:     c.PostForm("name"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Balance: balance,
	}
	// fmt.Println(c.PostForm("email"))
	// fmt.Println(c.PostForm("name"))
	// fmt.Println(c.PostForm("username"))
	// fmt.Println(c.PostForm("password"))

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
	var exists model.User
	if err := configdb.DB.Where("username = ?", users.Username).First(&exists).Error; err == nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"message": "username already exist",
		})
		return
	}
	users.Password = string(password)
	eror := configdb.DB.Save(&users).Error
	if eror != nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"message": "failed to sign up",
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":      http.StatusCreated,
		"message":     "User created Successfully!",
		"resourcedId": users.ID,
		})
}

// FetchAllUser function to get list of users
func FetchAllUser(c *gin.Context) {
	var users []model.User
	var user []user1
	configdb.DB.Find(&users)

	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		return
	}
	for _, item := range users {

		user = append(user, user1{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			DeletedAt: item.DeletedAt,
			Email:     item.Email,
			Name:      item.Name,
			Username:  item.Username,
			Balance: item.Balance,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": user})
}

// UpdateUser function to update user information
func UpdateUser(c *gin.Context) {
	var users model.User
	ID := c.Param("id")
	balance, _ := strconv.Atoi(c.PostForm("balance"))
	updatedUser := model.User{
		Email : c.PostForm("email"),
		Name:   c.PostForm("name"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Balance: balance,
	}
	configdb.DB.First(&users, ID)

	if users.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound, 
			"message": "No ID found!"})
			return
		}
	   
	configdb.DB.Model(&users).Update(&updatedUser)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		 "message": "User updated successfully!"})
}

// DeleteUser function to handle user deletion
func DeleteUser(c *gin.Context) {
	var users model.User
	usersID := c.Param("id")

	configdb.DB.First(&users, usersID)

	if users.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No user found!"})
		return
	}
	configdb.DB.Delete(&users)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK, "message": "User Delete succcessfully!"})
}

// FetchSingleUser function to get single user
func FetchSingleUser(c *gin.Context) {
	var users model.User
	usersID := c.Param("id")
	err := configdb.DB.Model(&model.User{}).Where("ID = ?", usersID).Find(&users).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		return
	}
	user := &user1{
		ID:        users.ID,
		CreatedAt: users.CreatedAt,
		UpdatedAt: users.UpdatedAt,
		DeletedAt: users.DeletedAt,
		Email:     users.Email,
		Name:      users.Name,
		Username:  users.Username,
		Balance : users.Balance,
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": user})
}

// Login function to handle login user
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user := &model.User{}
	if username == "" || password == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "please provide username and password"})
		return
	}

	err := configdb.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "wrong username"})
		return
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "wrong password or password doesn't match",
		})
		return
	}

	tk := &Token{
		Username: user.Username,
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
	configdb.DB.Save(&user)
	users := &user1{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
		Email:     user.Email,
		Name:      user.Name,
		Username:  user.Username,
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "logged in",
		"token":   tokenString,
		"user":    users,
	})
}

// Auth function authorization to handle authorized
func Auth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected SigningMethod :%v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if token != nil && err == nil {
		fmt.Println("token verified")
	} else {
		result := gin.H{
			"message": "not authorized",
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
	}
}

// func Logout(c *gin.Context){
// 	session := sessions.Default()
// 	user := session.

// }
