package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ariebrainware/paylist-api/config"
	"github.com/ariebrainware/paylist-api/model"
	"github.com/gin-gonic/gin"
)

var conf config.Config

// CreatePaylist function to create new paylist
func CreatePaylist(c *gin.Context) {
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	paylist := model.Paylist{
		Name:   c.PostForm("name"),
		Amount: amount,
	}
	fmt.Println(c.PostForm("name"))
	fmt.Println(amount)

	config.DB.Save(&paylist)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"message":    "Paylist item created successfully!",
		"resourceId": paylist.ID,
	})
}

//FetchAllPaylist Fetch All Paylist
func FetchAllPaylist(c *gin.Context) {
	var paylist []model.Paylist
	config.DB.Find(&paylist)

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
	err := config.DB.Model(&model.Paylist{}).Where("ID = ?", paylistID).Find(&paylist).Error

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
	config.DB.First(&paylist, id)

	if paylist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No ID found!"})
		return
	}

	config.DB.Model(&paylist).Update(&updatedPaylist)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Paylist updated successfully!"})
}

// DeletePaylist remove a paylist
func DeletePaylist(c *gin.Context) {
	var paylist model.Paylist
	paylistID := c.Param("id")

	config.DB.First(&paylist, paylistID)

	if paylist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No paylist found!"})
		return
	}
	config.DB.Delete(&paylist)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Paylist deleted successfully!"})
}
