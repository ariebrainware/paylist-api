package endpoint

import (
	"fmt"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/ariebrainware/paylist-api/config"
	"github.com/ariebrainware/paylist-api/model"
	"github.com/ariebrainware/paylist-api/util"
)

//User stuct for parse token
type User struct {
	Username string
	jwt.StandardClaims
}

//FetchAllPaylist Fetch All Paylist
func FetchAllPaylist(c *gin.Context) {
	var paylist []model.Paylist
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	fmt.Println(token.Valid, tk, err)
	username := tk.Username
	errf := config.DB.Model(&paylist).Where("username = ?", username).Find(&paylist).Error
	if errf != nil {
		util.CallErrorNotFound(c, "no paylist found!", errf)
		return
	}
	util.CallSuccessOK(c, "fetched all paylist", paylist)
}

//FetchSinglePaylist fetch a single paylist
func FetchSinglePaylist(c *gin.Context) {
	var paylist model.Paylist
	paylistID := c.Param("id")
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	fmt.Println(token.Valid, tk, err)
	username := tk.Username
	errf := config.DB.Model(&model.Paylist{}).Where("ID = ? and username = ?", paylistID, username).Find(&paylist).Error
	if errf != nil {
		util.CallErrorNotFound(c, "no paylist found!", errf)
		return
	}
	util.CallSuccessOK(c, "success fetch single paylist", paylist)
}

// CreateUserPaylist function to create new paylist
func CreateUserPaylist(c *gin.Context) {
	users := model.User{}
	tk := User{}

	// Parse the payload from token
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, "fail to parse the token, make sure token is valid", err)
		return
	}
	username := tk.Username

	// Decrease user balance
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	err = config.DB.Model(&users).Where("username  = ?", username).First(&users).Error
	if err != nil {
		util.CallErrorNotFound(c, "can't select balance", err)
		return
	}
	finalAmount := users.Balance - amount
	config.DB.Model(&users).Where("username = ?", username).Update("balance", finalAmount)
	paylist := model.Paylist{
		Name:      c.PostForm("name"),
		Amount:    amount,
		Username:  username,
		Completed: false,
	}

	// Save paylist
	err = config.DB.Model(&paylist).Save(&paylist).Error
	if err != nil {
		util.CallServerError(c, "fail to create paylist", err)
		return
	}
	util.CallSuccessOK(c, "paylist item created successfully!", paylist.ID)
}

//FetchAllPaylist Fetch All Paylist
func FetchAllPaylist(c *gin.Context) {
	var paylist []model.Paylist
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)
	username := tk.Username
	errf := config.DB.Model(&paylist).Where("username = ?", username).Find(&paylist).Error
	if errf != nil {
		util.CallErrorNotFound(c, "no paylist found!", errf)
		return
	}
	util.CallSuccessOK(c, "fetched all paylist", paylist)
}

//FetchSinglePaylist fetch a single paylist
func FetchSinglePaylist(c *gin.Context) {
	var paylist model.Paylist
	paylistID := c.Param("id")
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)
	username := tk.Username
	errf := config.DB.Model(&model.Paylist{}).Where("ID = ? and username = ?", paylistID, username).Find(&paylist).Error
	if errf != nil {
		util.CallErrorNotFound(c, "no paylist found!", errf)
		return
	}
	util.CallSuccessOK(c, "success fetch single paylist", paylist)
}

// UpdatePaylist update a paylist
func UpdatePaylist(c *gin.Context) {
	var paylist model.Paylist
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)
	username := tk.Username
	id, _ := strconv.Atoi(c.Param("id"))
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	updatedPaylist := model.Paylist{
		Name:   c.PostForm("name"),
		Amount: amount,
	}
	config.DB.First(&paylist, id)

	if paylist.ID == 0 {
		util.CallErrorNotFound(c, "no paylist found", nil)
		return
	}
	config.DB.Model(&paylist).Where("username = ?", username).Update(&updatedPaylist)
	if tk.Username != paylist.Username {
		util.CallServerError(c, "not authorized", nil)
		return
	}
	util.CallSuccessOK(c, "paylist successfully updated!", paylist)
}

// DeletePaylist remove a paylist
func DeletePaylist(c *gin.Context) {
	var paylist model.Paylist
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)
	username := tk.Username
	PaylistID := c.Param("id")
	config.DB.Where("ID = ?", PaylistID).Find(&paylist)
	if paylist.ID == 0 {
		util.CallErrorNotFound(c, "no paylist found!", nil)
		c.Abort()
		return
	}
	config.DB.Model(&paylist).Where("username = ?", username).Find(&paylist)
	if tk.Username != paylist.Username {
		util.CallServerError(c, "user not authorized", nil)
		c.Abort()
		return
	}
}

//CreateUserPaylist handle add user paylist
func CreateUserPaylist(c *gin.Context) {
	//var paylist model.Paylist
	var users model.User
	var sisa int
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)
	username := tk.Username
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	eror := config.DB.Table("users").Select("balance").Where("username  = ?", username).First(&users).Error
	if eror != nil {
		util.CallErrorNotFound(c, "can't select balance", eror)
		return
	}
	sisa = users.Balance - amount
	users.Balance = sisa
	config.DB.Model(&users).Where("username = ?", username).Update(&users)
	fmt.Println("sisa", sisa)
}

//DeleteUserPaylist handle deleted user paylist
func DeleteUserPaylist(c *gin.Context) {
	var paylist model.Paylist
	var user model.User
	var balance int
	id := c.Param("id")
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)
	username := tk.Username
	errf := config.DB.Table("paylists").Select("amount, completed").Where("ID = ? and username = ?", id, username).Find(&paylist).Error
	if errf != nil {
		util.CallErrorNotFound(c, "can't select amount", errf)
		return
	}
	fmt.Println(paylist.Amount)
	erf := config.DB.Table("users").Select("balance").Where("username = ?", username).Find(&user).Error
	if erf != nil {
		util.CallErrorNotFound(c, "can't select balance", erf)
		return
	}
	fmt.Println(user.Balance)
	if paylist.Completed == false {
		balance = paylist.Amount + user.Balance
		user.Balance = balance
		fmt.Println(user.Balance)
		config.DB.Model(&user).Where("username = ?", username).Update(&user)
	}
	eror := config.DB.Model(&paylist).Where("ID = ? and username = ?", id, username).Delete(&paylist).Error
	if eror == nil {
		util.CallSuccessOK(c, "paylist successfully deleted!", nil)
	}
}

