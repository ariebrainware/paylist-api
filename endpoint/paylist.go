package endpoint

import (
	"fmt"
	"strconv"
	"log"
	"github.com/gin-gonic/gin"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ariebrainware/paylist-api/config"
	"github.com/ariebrainware/paylist-api/model"
	"github.com/ariebrainware/paylist-api/util"
)

type User struct {
	Username string
	jwt.StandardClaims
}

// CreatePaylist function to create new paylist
func CreatePaylist(c *gin.Context) {
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	paylist := model.Paylist{
		Name:   c.PostForm("name"),
		Amount: amount,
		Completed: completed,
	}
	fmt.Println(c.PostForm("name"))
	fmt.Println(amount)
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)

	paylist.Username = string(tk.Username)
	errf := config.DB.Save(&paylist).Error
	if errf != nil {
		util.CallServerError(c, "fail to create paylist", errf)	
		return
	}
	util.CallSuccessOK(c, "paylist item created successfully!", paylist.ID)
}

//FetchAllPaylist Fetch All Paylist
func FetchAllPaylist(c *gin.Context) {
	var paylist []model.Paylist
	
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)
	username := tk.Username
	config.DB.Model(&paylist).Where("username = ?", username).Find(&paylist)
	if len(paylist) <= 0 {
		util.CallErrorNotFound(c, "no paylist found!", nil)
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
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)
	username := tk.Username
	errf := config.DB.Model(&model.Paylist{}).Where("ID = ? and username = ?",paylistID, username).Find(&paylist).Error
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
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)

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
	username := tk.Username
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
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	})
	log.Println(token.Valid, tk, err)

	PaylistID := c.Param("id")
	config.DB.Where("ID = ?", PaylistID).Find(&paylist)
	if paylist.ID == 0 {
		util.CallErrorNotFound(c, "no paylist found!", nil)
		return
	}
	username := tk.Username
	config.DB.Model(&paylist).Where("username = ?", username).Delete(&paylist)
	if tk.Username != paylist.Username {
		util.CallServerError(c, "not authorized", nil)
		return
	}
	util.CallSuccessOK(c, "paylist successfully deleted!", nil)
}
//CreateUserPaylist handle add user paylist
func CreateUserPaylist(c *gin.Context) {
	//var paylist model.Paylist
	var users model.User
	var sisa int
	
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	}) 
	log.Println(token.Valid, tk, err)
	
	username := tk.Username
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	eror:= config.DB.Table("users").Select("balance").Where("username  = ?", username).First(&users).Error
	if eror != nil{
		util.CallErrorNotFound(c, "can't select balance", eror)
		return
	} 
	sisa = users.Balance - amount
	users.Balance = sisa
	config.DB.Model(&users).Where("username = ?", username).Update(&users)
	fmt.Println("sisa", sisa)
}
//DeleteUserPaylist handle deleted user paylist
func DeleteUserPaylist(c *gin.Context){
	var paylist model.Paylist
	var user model.User
	var balance int
	id := c.Param("id")
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	}) 
	log.Println(token.Valid, tk, err)
	username := tk.Username
	errf := config.DB.Table("paylists").Select("amount, completed").Where("ID = ? and username = ?",id, username).Find(&paylist).Error
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
	if paylist.Completed == 0 {
		balance = paylist.Amount + user.Balance	
		user.Balance = balance 
		fmt.Println(user.Balance)
		config.DB.Model(&user).Where("username = ?", username).Update(&user)
	}
}
//Completed
func Completed(c *gin.Context){
	var paylist model.Paylist
	var user model.User
	id := c.Param("id")
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	}) 
	log.Println(token.Valid, tk, err)
	username := tk.Username
	config.DB.Model(&paylist).Where("ID = ? and username = ?", id, username).Find(&paylist)
	config.DB.Model(&user).Select("balance").Where("username = ?", username).Find(&user)
	
	if user.Balance > 0 && paylist.Completed == 0 {
		paylist.Completed = 1
		config.DB.Model(&paylist).Where("ID = ? and username = ?",id, username).Update(&paylist)
	} else if user.Balance < 0 && paylist.Completed == 0 {
		paylist.Completed = 0
		config.DB.Model(&paylist).Where("ID = ? and username = ?", id, username).Update(&paylist)
	}
}