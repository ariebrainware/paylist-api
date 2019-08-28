package endpoint

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ariebrainware/paylist-api/config"
	"github.com/ariebrainware/paylist-api/model"
	jwt "github.com/dgrijalva/jwt-go" //Used to sign and verify JWT tokens
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/ariebrainware/paylist-api/util"
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
		util.CallServerError(c, "password encryption failed", err)
		return
	}
	util.CallSuccessOK(c, "password encryption success!", password)

	users.Password = string(password)
	err = config.DB.Save(&users).Error
	if err != nil {
		util.CallServerError(c, "Failed Create User!", err)
	}
	util.CallSuccessOK(c, "User created Successfully!", users.ID)
}

// FetchAllUser function to get list of users
func FetchAllUser(c *gin.Context) {
	var users []model.User
	var user []user1
	config.DB.Find(&users)

	if len(users) <= 0 {
		util.CallErrorNotFound(c, "No User Found!", nil)
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
	util.CallSuccessOK(c, "Fetch All User Data ", user)
}

// UpdateUser function to update user information
func UpdateUser(c *gin.Context) {
	var users model.User
	ID := c.Param("id")
	updatedUser := model.User{
		Email:    c.PostForm("email"),
		Name:     c.PostForm("name"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	config.DB.First(&users, ID)

	if users.ID == 0 {
		util.CallErrorNotFound(c, "Paylist not found, make sure to specify the ID", nil)
		return
	}

	err := config.DB.Model(&users).Update(&updatedUser).Error
	if err != nil {
		util.CallServerError(c, "Failed to update user", err)
	}
	util.CallSuccessOK(c, "User successfully updated!", users)
}

// DeleteUser function to handle user deletion
func DeleteUser(c *gin.Context) {
	var users model.User
	usersID := c.Param("id")

	config.DB.First(&users, usersID)

	if users.ID == 0 {
		util.CallErrorNotFound(c, "user not found", nil)
		return
	}
	err := config.DB.Delete(&users).Error
	if err != nil {
		util.CallServerError(c, "failed to delete user", err)
		return
	}
	util.CallSuccessOK(c, "user delete succcessfully!", nil)
}

// FetchSingleUser function to get single user
func FetchSingleUser(c *gin.Context) {
	var users model.User
	usersID := c.Param("id")
	err := config.DB.Model(&model.User{}).Where("ID = ?", usersID).Find(&users).Error

	if err != nil {
		util.CallErrorNotFound(c, "no user found", nil)
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
	util.CallSuccessOK(c, "success fetch single data", user)
}

// Login function to handle login user
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user := &model.User{}
	if username == "" || password == "" {
		util.CallErrorNotFound(c, "please provide username and password", nil)
		return
	}

	err := config.DB.Model(&user).Where("username = ?", username).First(&user).Error
	if err != nil {
		util.CallErrorNotFound(c, "wrong username or password", err)
		return
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	fmt.Println(user.Password)
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		util.CallErrorNotFound(c, "wrong password or password doesn't match", errf)
		return
	}

	tk := &Token{
		Username: user.Username,
	}
	//Create JWT token
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		util.CallServerError(c, "error create token", err)
		c.Abort()
	}
	util.CallSuccessOK(c, "logged in", tokenString)
}

// Auth function authorization to handle authorized
func Auth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
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
