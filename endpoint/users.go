package endpoint

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go" //Used to sign and verify JWT tokens
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	Password  string     `json:"password"`
	Balance   int        `json:"balance"`
}

// CreateUser function to sign up
func CreateUser(c *gin.Context) {
	users := model.User{
		Email:    c.PostForm("email"),
		Name:     c.PostForm("name"),
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	//check field can't empty
	if users.Username == "" {
		util.CallUserError(c, util.APIErrorParams{Msg: "username can't be null"})
		return
	}
	if users.Name == "" {
		util.CallUserError(c, util.APIErrorParams{Msg: "name can't be null"})
		return
	}
	if users.Password == "" {
		util.CallUserError(c, util.APIErrorParams{Msg: "password can't be null"})
		return
	}
	if users.Email == "" {
		util.CallUserError(c, util.APIErrorParams{Msg: "email can't be null"})
		return
	}

	//check username exist or not
	var exists model.User
	if err := config.DB.Where("username = ?", users.Username).First(&exists).Error; err == nil {
		util.CallServerError(c, util.APIErrorParams{Msg: "username already exist", Err: err})
		return
	}
	//Password Encryption
	password, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
	if err != nil {
		util.CallServerError(c, util.APIErrorParams{Msg: "password encryption failed", Err: err})
		return
	}

	users.Password = string(password)
	err = config.DB.Save(&users).Error
	if err != nil {
		util.CallServerError(c, util.APIErrorParams{Msg: "Failed Create User!", Err: err})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{Msg: "User created Successfully!", Data: users.ID})
}

// FetchAllUser function to get list of users
func FetchAllUser(c *gin.Context) {
	var users []model.User
	var user []user1
	tk := User{}
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprint(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, util.APIErrorParams{Msg: "fail to parse the token, make sure token is valid", Err: err})
		return
	}
	username := tk.Username
	config.DB.Model(&users).Where("username = ? ", username).Find(&users)

	if len(users) <= 0 {
		util.CallErrorNotFound(c, util.APIErrorParams{Msg: "No User Found!"})
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
	util.CallSuccessOK(c, util.APISuccessParams{Msg: "Fetch All Users Data", Data: user})
}

// UpdateUser function to update user information
func UpdateUser(c *gin.Context) {
	var users model.User
	ID := c.Param("id")
	tk := User{}
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprint(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, util.APIErrorParams{Msg: "fail to parse the token, make sure token is valid", Err: err})
		return
	}
	username := tk.Username
	balance, _ := strconv.Atoi(c.PostForm("balance"))
	if balance == 0 || balance < 0 {
		util.CallUserError(c, util.APIErrorParams{Msg: "please specify the amount of balance, it can't be negative or zero"})
		return
	}
	user := model.User{
		Email:    c.PostForm("email"),
		Name:     c.PostForm("name"),
		Username: c.PostForm("username"),
		Balance:  balance,
	}
	config.DB.First(&users, ID)
	if users.ID == 0 {
		util.CallErrorNotFound(c, util.APIErrorParams{Msg: "user not found, make sure to specify the id"})
		return
	}
	err = config.DB.Model(&users).Where("username = ? and ID = ?", username, ID).Update(&user).Error
	if err != nil {
		util.CallServerError(c, util.APIErrorParams{Msg: "Failed to update user", Err: err})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{Msg: "User successfully updated!", Data: ID})
}

func EditPassword(c *gin.Context) {
	var users model.User
	ID := c.Param("id")
	tk := User{}
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprint(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, util.APIErrorParams{Msg: "fail to parse the token, make sure token is valid", Err: err})
		return
	}
	username := tk.Username

	OldPassword := c.PostForm("OldPassword")
	NewPassword := c.PostForm("NewPassword")
	config.DB.First(&users, ID)
	if users.ID == 0 {
		util.CallErrorNotFound(c, util.APIErrorParams{Msg: "user not found, make sure to specify the id"})
		return
	}

	errf := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(OldPassword))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		util.CallErrorNotFound(c, util.APIErrorParams{Msg: "password doesn't match", Err: errf})
		return
	}

	//Password Encryption
	password, err := bcrypt.GenerateFromPassword([]byte(NewPassword), bcrypt.DefaultCost)
	if err != nil {
		util.CallServerError(c, util.APIErrorParams{Msg: "password encryption failed", Err: err})
		return
	}
	users.Password = string(password)
	err = config.DB.Model(&users).Where("username = ? and ID = ?", username, ID).Update("password", users.Password).Error
	if err != nil {
		util.CallServerError(c, util.APIErrorParams{Msg: "Failed to update user", Err: err})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{Msg: "Password successfully updated!", Data: ID})
}

// AddBalance is a function to add user balance or income
func AddBalance(c *gin.Context) {
	var users model.User
	var inc model.Income
	tk := User{}
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprintf(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, util.APIErrorParams{Msg: "fail to parse the token, make sure token is valid", Err: err})
		return
	}
	username := tk.Username
	if err = config.DB.Model(&users).Where("username = ?", username).Find(&users).Error; err != nil || err == gorm.ErrRecordNotFound {
		util.CallErrorNotFound(c, util.APIErrorParams{Msg: "no user found"})
		return
	}

	firstBalance := users.Balance
	balance, _ := strconv.Atoi(c.PostForm("balance"))
	if balance == 0 || balance < 0 {
		util.CallUserError(c, util.APIErrorParams{Msg: "please specify the amount of balance, it can't be negative or zero"})
		return
	}
	err = config.DB.Model(&users).Where("username = ?", username).Update("balance", balance+firstBalance).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	data := model.Income{
		Username: username,
		Income:   balance,
	}
	err = config.DB.Model(&inc).Save(&data).Error
	if err != nil {
		fmt.Println("error", err)
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{Msg: "successfully add balance"})
}

// DeleteUser function to handle user deletion
func DeleteUser(c *gin.Context) {
	var users model.User
	usersID := c.Param("id")
	tk := User{}
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprintf(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, util.APIErrorParams{Msg: "fail to parse the token, make sure token is valid", Err: err})
		return
	}
	username := tk.Username
	config.DB.First(&users, usersID)
	if users.ID == 0 {
		util.CallErrorNotFound(c, util.APIErrorParams{Msg: "user not found"})
		return
	}
	config.DB.Model(&users).Where("username = ?", username).Delete(&users)
	if tk.Username != users.Username {
		util.CallServerError(c, util.APIErrorParams{Msg: "not authorized"})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{Msg: "user delete successfully!"})
}

// FetchSingleUser function to get single user
func FetchSingleUser(c *gin.Context) {
	var users model.User
	usersID := c.Param("id")
	tk := User{}
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprintf(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, util.APIErrorParams{Msg: "fail to parse the token, make sure token is valid", Err: err})
		return
	}
	username := tk.Username

	errf := config.DB.Model(&model.User{}).Where("ID = ? and username = ?", usersID, username).Find(&users).Error
	if errf != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{Msg: "no user found", Err: errf})
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
	util.CallSuccessOK(c, util.APISuccessParams{Msg: "success fetch single data", Data: user})
}

// Login function to handle login user
func Login(c *gin.Context) {
	logging := &model.Logging{}
	username := c.PostForm("username")
	password := c.PostForm("password")

	user := &model.User{}
	if username == "" || password == "" {
		util.CallErrorNotFound(c, util.APIErrorParams{Msg: "please provide username and password"})
		return
	}
	err := config.DB.Model(&user).Where("username = ?", username).First(&user).Error
	if err != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{Msg: "wrong username", Err: err})
		return
	}
	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		util.CallErrorNotFound(c, util.APIErrorParams{Msg: "wrong password or password doesn't match", Err: errf})
		return
	}

	expirationTime := time.Now().Add(720 * time.Minute)
	tk := &Token{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	//Create JWT token
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(fmt.Sprintf(config.Conf.JWTSignature)))
	if err != nil {
		util.CallServerError(c, util.APIErrorParams{Msg: "error create token", Err: err})
		c.Abort()
	}
	data := &model.Logging{
		Token:      tokenString,
		Username:   username,
		UserStatus: true,
	}
	config.DB.Model(&logging).Find(&logging)
	if logging.Username == username {
		util.CallUserFound(c, util.APISuccessParams{Msg: "already login"})
		c.Abort()
		return
	}
	if err = config.DB.Model(&logging).Save(&data).Error; err != nil {
		util.CallServerError(c, util.APIErrorParams{Msg: "fail to save logging data", Err: err})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	util.CallSuccessOK(c, util.APISuccessParams{Msg: "logged in", Data: tokenString})
}

// Auth function authorization to handle authorized
func Auth(c *gin.Context) {
	claim := Token{}
	logging := &model.Logging{}
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected SigningMethod :%v", token.Header["alg"])
		}
		return []byte(fmt.Sprintf(config.Conf.JWTSignature)), nil
	})
	config.DB.Model(&logging).Where("token = ? ", tokenString).Find(&logging)
	if logging.Token == "" {
		util.CallServerError(c, util.APIErrorParams{Msg: "you have to sign in first"})
		c.Abort()
	} else if token != nil && time.Unix(claim.ExpiresAt, 0).Sub(time.Now()) < 30*time.Second {
		util.CallUserError(c, util.APIErrorParams{Msg: "token expired", Err: err})
		err = config.DB.Model(&logging).Where("token = ?", tokenString).Delete(&logging).Error
		if err != nil {
			fmt.Println(err)
			util.CallServerError(c, util.APIErrorParams{Msg: "fail when try to delete the logging", Err: err})
		}
		c.Abort()
		return
	}
}

// Logout handle logout user
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
		util.CallServerError(c, util.APIErrorParams{Msg: "fail when try to delete the logging", Err: err})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{Msg: "logged out", Data: logging.UserStatus})
	return
}
