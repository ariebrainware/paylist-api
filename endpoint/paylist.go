package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ariebrainware/paylist-api/config"
	"github.com/ariebrainware/paylist-api/model"
	"github.com/ariebrainware/paylist-api/util"
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

	err := config.DB.Save(&paylist).Error
	if err != nil {
		util.CallSuccessOK(c, "Paylist item created successfully!", paylist.ID)
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":     http.StatusInternalServerError,
		"message":    "Fail to create paylist",
		"resourceId": paylist.ID,
	})
}

//FetchAllPaylist Fetch All Paylist
func FetchAllPaylist(c *gin.Context) {
	var paylist []model.Paylist
	config.DB.Find(&paylist)

	if len(paylist) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "No paylist found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": paylist,
		})	
}

//FetchSinglePaylist fetch a single paylist
func FetchSinglePaylist(c *gin.Context) {
	var paylist model.Paylist
	paylistID := c.Param("id")
	err := config.DB.Model(&model.Paylist{}).Where("ID = ?", paylistID).Find(&paylist).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No paylist found!",
			"error":   err,
		})
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

	err := config.DB.Model(&paylist).Update(&updatedPaylist).Error
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{
			"status":  http.StatusNotImplemented,
			"message": "Failed update paylist!",
			"error":   err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Paylist updated successfully!",
		"data":    paylist,
	})
}

// DeletePaylist remove a paylist
func DeletePaylist(c *gin.Context) {
	var paylist model.Paylist
	paylistID := c.Param("id")

	config.DB.First(&paylist, paylistID)

	if paylist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No paylist found!"})
		return
	}
	err := config.DB.Delete(&paylist).Error
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Paylist deleted successfully!"})
}
