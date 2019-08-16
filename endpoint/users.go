package endpoint

import (
	"fmt"
	"net/http"
	//"strconv"
	"time"

	"github.com/ariebrainware/paylist-api/model"
	jwt "github.com/dgrijalva/jwt-go" //Used to sign and verify JWT tokens
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Token is a struct for token model
type Token struct {
	ID uint
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
}

// CreateUser function to sign up
func CreateUser(c *gin.Context) {
	users := model.User{
		Email:    c.PostForm("email"),
		Name:     c.PostForm("name"),
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

// FetchAllUser function to get list of users
func FetchAllUser(c *gin.Context) {
	var users []model.User
	var user []user1
	db.Find(&users)

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
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": user})
}

// UpdateUser function to update user information
func UpdateUser(c *gin.Context) {
	var users model.User
	ID := c.Param("id")
	updatedUser := model.User{
		Email : c.PostForm("email"),
		Name:   c.PostForm("name"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	db.First(&users, ID)

	if users.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound, 
			"message": "No ID found!"})
			return
		}
	   
	db.Model(&users).Update(&updatedUser)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		 "message": "User updated successfully!"})
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
	user := &user1{
		ID:        users.ID,
		CreatedAt: users.CreatedAt,
		UpdatedAt: users.UpdatedAt,
		DeletedAt: users.DeletedAt,
		Email:     users.Email,
		Name:      users.Name,
		Username:  users.Username,
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

	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "wrong username or password"})
		return
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	fmt.Println(user.Password)
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "wrong password or password doesn't match",
		})
		return
	}

	tk := &Token{
		ID: user.ID,
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
