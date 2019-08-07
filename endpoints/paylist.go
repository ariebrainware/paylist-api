package endpoints

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	
	"github.com/ariebrainware/paylist-api/model"
	
)

var db *gorm.DB

func init() {
	var err error

	//connString := "userdb:passworddb/databasename?charset=utf8&parseTime=True&loc=Local"
	connString := "coba:admin123@tcp(localhost:3306)/paylist?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", connString)
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&model.Paylist{})
	db.AutoMigrate(&model.User{})
	fmt.Println("Schema migrated!!")
}

//CreatePaylist
func CreatePaylist(c *gin.Context) {
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	paylist := model.Paylist{
		Name:   c.PostForm("name"),
		Amount: amount,
	}
	fmt.Println(c.PostForm("name"))
	fmt.Println(amount)

	db.Save(&paylist)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"message":    "Paylist item created successfully!",
		"resourceId": paylist.ID,
	})
}

//FetchAllPaylist Fetch All Paylist
func FetchAllPaylist(c *gin.Context) {
	var paylist []model.Paylist
	db.Find(&paylist)

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
	err := db.Model(&model.Paylist{}).Where("ID = ?", paylistID).Find(&paylist).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No paylist found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": paylist})
}

//UpdatePaylists update a paylist
func UpdatePaylist(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	updatedPaylist := model.Paylist{
		Name:   c.PostForm("name"),
		Amount: amount,
	}
	err := db.Model(&model.Paylist{}).Where("ID = ?", id).Update(&updatedPaylist).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusOK,
			"message": "Paylist updated successfully!",
			"errors":  err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Paylist updated successfully!",
		"errors":  err,
	})
}

//DeletePaylist remove a paylist
func DeletePaylist(c *gin.Context) {
	var paylist model.Paylist
	paylistID := c.Param("id")

	db.First(&paylist, paylistID)

	if paylist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No paylist found!"})
		return
	}
	db.Delete(&paylist)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Paylist deleted successfully!"})
}