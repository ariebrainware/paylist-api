package endpoint

import (
	"fmt"
	"net/http"
	"time"
	"strconv"
	"log"
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
	Balance int `json:"balance"`
}

// CreateUser function to sign up
func CreateUser(c *gin.Context) {
	//balance, _ := strconv.Atoi(c.PostForm("balance"))
	users := model.User{
		Email:    c.PostForm("email"),
		Name:     c.PostForm("name"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		//Balance: balance,
	}
	// fmt.Println(c.PostForm("email"))
	// fmt.Println(c.PostForm("name"))
	// fmt.Println(c.PostForm("username"))
	// fmt.Println(c.PostForm("password"))
	//check username exist or not
	if users.Username == "" || users.Name == "" || users.Password == "" || users.Email == ""  {
		util.CallServerError(c, "field can't be null", nil)
		return
	}
	var exists model.User
	if err := config.DB.Where("username = ?", users.Username).First(&exists).Error; err == nil {
		util.CallServerError(c, "username already exist", err)
		return
	}
	//Password Encryption
	password, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
	if err != nil {
		util.CallServerError(c, "password encryption failed", err)
		return
	}
	//util.CallSuccessOK(c, "password encryption success!", password)
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
			Balance: item.Balance,
		})
	}
	util.CallSuccessOK(c, "Fetch All User Data ", user)
}

// UpdateUser function to update user information
func UpdateUser(c *gin.Context) {
	var users model.User
	ID := c.Param("id")
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)
	username := tk.Username
	balance, _ := strconv.Atoi(c.PostForm("balance"))
	updatedUser := model.User{
		Email:    c.PostForm("email"),
		Name:     c.PostForm("name"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Balance: balance,
	}
	config.DB.First(&users, ID)

	if users.ID == 0 {
		util.CallErrorNotFound(c, "Paylist not found, make sure to specify the ID", nil)
		return
	}

	errf := config.DB.Model(&users).Where("username = ?", username).Update(&updatedUser).Error
	if errf != nil {
		util.CallServerError(c, "Failed to update user", errf)
	}
	util.CallSuccessOK(c, "User successfully updated!", ID)
}

// DeleteUser function to handle user deletion
func DeleteUser(c *gin.Context) {
	var users model.User
	usersID := c.Param("id")
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)
	username := tk.Username

	config.DB.First(&users, usersID)
	if users.ID == 0 {
		util.CallErrorNotFound(c, "user not found", nil)
		return
	}
	config.DB.Model(&users).Where("username = ?", username).Delete(&users)
	if tk.Username != users.Username {
		util.CallServerError(c, "not authorized", nil)
		return
	}
	util.CallSuccessOK(c, "user delete succcessfully!", nil)
}

// FetchSingleUser function to get single user
func FetchSingleUser(c *gin.Context) {
	var users model.User
	usersID := c.Param("id")
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)
	username := tk.Username
	errf := config.DB.Model(&model.User{}).Where("ID = ? and username = ?", usersID, username).Find(&users).Error
	if errf != nil {
		util.CallErrorNotFound(c, "no user found", errf)
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
	util.CallSuccessOK(c, "success fetch single data", user)
}

// Login function to handle login user
func Login(c *gin.Context) {
	logging := &model.Logging{}
	username := c.PostForm("username")
	password := c.PostForm("password")

	user := &model.User{}
	if username == "" || password == "" {
		util.CallErrorNotFound(c, "please provide username and password", nil)
		return
	}

	err := config.DB.Model(&user).Where("username = ?", username).First(&user).Error
	if err != nil {
		util.CallErrorNotFound(c, "wrong username", err)
		return
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		util.CallErrorNotFound(c, "wrong password or password doesn't match", errf)
		return
	}

	expirationTime := time.Now().Add(2 * time.Minute)
	tk := &Token{
		Username: user.Username,
		StandardClaims : jwt.StandardClaims{
			ExpiresAt : expirationTime.Unix(),
		},
	}
	//Create JWT token
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		util.CallServerError(c, "error create token", err)
		c.Abort()
	}
	logging.Token = tokenString
	logging.Username = username
	logging.User_status = true
	config.DB.Save(&logging)
	
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
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

//RefreshToken hanlde refreshing expired token
func RefreshToken(c *gin.Context) {
	claim := Token{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	})
	log.Println(token.Valid, claim, err)
	if time.Unix(claim.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	expirationTime := time.Now().Add(12 * time.Hour)
	claim.ExpiresAt = expirationTime.Unix()
	tokenn:= jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tknString, err := tokenn.SignedString([]byte("secret"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tknString,
		Expires: expirationTime,
	})
	util.CallSuccessOK(c, "refresh token success", tknString)
}

//Logout handle logout user
func Logout(c *gin.Context){
	var logging model.Logging
	tokenStr := c.GetHeader("Authorization")
	token := config.DB.Model(&logging).Where("token = ?",tokenStr ).Find(&logging).Error
	if token != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate session token"})
		return
	}
	config.DB.Model(&logging).Where("token = ?", tokenStr).Delete(&logging)
	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}