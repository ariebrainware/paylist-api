package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ariebrainware/paylist-api/configdb"
	"github.com/ariebrainware/paylist-api/model"
)
var conf configdb.Config
// CreatePaylist function to create new paylist
func CreatePaylist(c *gin.Context) {
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	paylist := model.Paylist{
		Name:   c.PostForm("name"),
		Amount: amount,
	}
	fmt.Println(c.PostForm("name"))
	fmt.Println(amount)

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
	configdb.DB.Find(&paylist)

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
			"message": "No ID found!"})
			return
		}
	   
	configdb.DB.Model(&paylist).Update(&updatedPaylist)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		 "message": "Paylist updated successfully!"})
}

// DeletePaylist remove a paylist
func DeletePaylist(c *gin.Context) {
	var paylist model.Paylist
	paylistID := c.Param("id")

	configdb.DB.First(&paylist, paylistID)

	if paylist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No paylist found!"})
		return
	}
	configdb.DB.Delete(&paylist)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Paylist deleted successfully!"})
}
