package endpoint

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go" //Used to sign and verify JWT tokens
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/ariebrainware/paylist-api/config"
	"github.com/ariebrainware/paylist-api/model"
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
	Balance   int        `json:"balance"`
}

// CreateUser function to sign up
func CreateUser(c *gin.Context) {
	//balance, _ := strconv.Atoi(c.PostForm("balance"))
	users := model.User{
		Email:    c.PostForm("email"),
		Name:     c.PostForm("name"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	//check username exist or not
	if users.Username == "" || users.Name == "" || users.Password == "" || users.Email == "" {
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
			Balance:   item.Balance,
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
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, "fail to parse the token, make sure token is valid", err)
		return
	}
	username := tk.Username
	balance, _ := strconv.Atoi(c.PostForm("balance"))
	updatedUser := model.User{
		Email:    c.PostForm("email"),
		Name:     c.PostForm("name"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Balance:  balance,
	}
	config.DB.First(&users, ID)

	if users.ID == 0 {
		util.CallErrorNotFound(c, "user not found, make sure to specify the id", nil)
		return
	}

	if balance == 0 || balance < 0 {
		util.CallUserError(c, "please specify the amount of balance, it can't be negative or zero", nil)
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
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, "fail to parse the token, make sure token is valid", err)
		return
	}
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
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, "fail to parse the token, make sure token is valid", err)
		return
	}
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
		Balance:   users.Balance,
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

	expirationTime := time.Now().Add(1 * time.Minute)
	tk := &Token{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	//Create JWT token
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		util.CallServerError(c, "error create token", err)
		c.Abort()
	}
	config.DB.Model(&logging).Find(&logging)
	if logging.Username == username {
		util.CallServerError(c, "already login", nil)
		c.Abort()
		return
	}
	data := &model.Logging{
		Token:      tokenString,
		Username:   username,
		UserStatus: true,
	}
	if err = config.DB.Model(&logging).Save(&data).Error; err !=nil {
		util.CallServerError(c,"fail to save logging data",err)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	util.CallSuccessOK(c, "logged in", tokenString)
}

// Auth function authorization to handle authorized
func Auth(c *gin.Context) {
	logging := &model.Logging{}
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected SigningMethod :%v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	config.DB.Model(&logging).Where("token = ? ", tokenString).Find(&logging)
	if logging.Token == "" {
		util.CallServerError(c, "you have to sign in first", nil)
		c.Abort()
		return
	}
	if token != nil && err == nil {
		fmt.Println("token verified")
		return
	}
	result := gin.H{
		"message": "not authorized",
	}
	c.JSON(http.StatusUnauthorized, result)
	c.Abort()
}

//RefreshToken hanlde refreshing expired token
func RefreshToken(c *gin.Context) {
	claim := Token{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, "fail to parse the token, make sure token is valid", err)
		return
	}
	if time.Unix(claim.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	expirationTime := time.Now().Add(12 * time.Hour)
	claim.ExpiresAt = expirationTime.Unix()
	tokenn := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
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
func Logout(c *gin.Context) {
	logging := &model.Logging{}
	tokenStr := c.GetHeader("Authorization")
	erf := config.DB.Model(&logging).Where("token = ?", tokenStr).Update("userStatus", false).Error
	if erf != nil {
		fmt.Println(erf)
	}
	err := config.DB.Model(&logging).Where("token = ?", tokenStr).Delete(&logging).Error
	if err != nil {
		fmt.Println(err)
		util.CallServerError(c,"fail when try to delete the logging", err)
	}
	util.CallSuccessOK(c, "logged out", logging.UserStatus)
}

//SignOut for check token expired
func SignOut(c *gin.Context) {
	claim := Token{}
	tokenString := c.Request.Header.Get("Authorization")
	token, _ := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if token != nil && time.Unix(claim.ExpiresAt, 0).Sub(time.Now()) < 30*time.Second {
		util.CallSuccessOK(c,"token invalid and expired", tokenString)
	}
	if token != nil && time.Unix(claim.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		util.CallSuccessOK(c,"token valid and not expired", tokenString)
	}
}
