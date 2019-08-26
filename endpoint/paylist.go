package endpoint

import (
	"fmt"
	"net/http"
	"strconv"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/ariebrainware/paylist-api/configdb"
	"github.com/ariebrainware/paylist-api/model"
	jwt "github.com/dgrijalva/jwt-go"
)
type User struct {
	Username string
	jwt.StandardClaims
}

var conf configdb.Config
// CreatePaylist function to create new paylist
func CreatePaylist(c *gin.Context) {
	//var users model.User
	//username := c.Param("username")
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	paylist := model.Paylist{
		Name:   c.PostForm("name"),
		Amount: amount,
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
	configdb.DB.Save(&paylist)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"message":    "Paylist item created successfully!",
		"resourceId": paylist.ID,
	})
}

//FetchAllPaylist Fetch All Paylist
func FetchAllPaylist(c *gin.Context) {
	var paylist []model.Paylist
	//var user model.User
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	})
	username := tk.Username
	log.Println(token.Valid, tk, err)
	configdb.DB.Model(&paylist).Where("username = ?", username).Find(&paylist)

	if len(paylist) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No paylist found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": paylist})
}

//FetchSinglePaylist fetch a single paylist
func FetchSinglePaylist(c *gin.Context) {
	var paylist model.Paylist
	paylistID := c.Param("id")
	err := configdb.DB.Model(&model.Paylist{}).Where("ID = ?", paylistID).Find(&paylist).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No paylist found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": paylist})
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

	configdb.DB.First(&paylist, id)
	if paylist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound, 
			"message": "No paylist ID found!"})
			return
		}
	username := tk.Username
	configdb.DB.Model(&paylist).Where("username = ?", username).Update(&updatedPaylist)
	if tk.Username != paylist.Username {
		c.JSON(http.StatusNotFound, gin.H{"message": "not authorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "Paylist updated successfully!"})
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
	configdb.DB.Where("ID = ?", PaylistID).Find(&paylist)
	if paylist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No paylist found!"})
		return
	}
	username := tk.Username
	configdb.DB.Model(&paylist).Where("username = ?", username).Delete(&paylist)
	if tk.Username != paylist.Username {
		c.JSON(http.StatusNotFound, gin.H{"message": "not authorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Paylist deleted successfully!"})
}

func Coba(c *gin.Context) {
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
	configdb.DB.Table("users").Select("balance").Where("username  = ?", username).First(&users)
	fmt.Println(users.Balance)
	if users.Balance > amount {
		sisa = users.Balance - amount
	 } else {
		fmt.Println("Saldo Anda Kurang!")
	 }
	users.Balance = sisa
	configdb.DB.Model(&users).Where("username = ?", username).Update(&users)
	fmt.Println("sisa", sisa)
}
//func 
func Coba2(c *gin.Context){
	var paylist model.Paylist
	var user model.User
	var balance int
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error){
		return []byte("secret"), nil
	}) 
	log.Println(token.Valid, tk, err)
	username := tk.Username
	configdb.DB.Table("paylists").Select("amount").Where("username  = ?", username).Find(&paylist)
	configdb.DB.Table("users").Select("balance").Where("username  = ?", username).Find(&user)
	fmt.Println(paylist.Amount)
	fmt.Println(user.Balance)
	balance = paylist.Amount + user.Balance
	user.Balance = balance 
	fmt.Println(user.Balance)
	configdb.DB.Model(&user).Where("username = ?", username).Update(&user)	
}