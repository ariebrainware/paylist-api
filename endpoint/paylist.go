package endpoint

import (
	"fmt"
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

	err := config.DB.Save(&paylist).Error
	if err != nil {
		util.CallSuccessOK(c, "paylist item created successfully!", paylist.ID)
	}
	util.CallServerError(c, "fail to create paylist", err)
}

//FetchAllPaylist Fetch All Paylist
func FetchAllPaylist(c *gin.Context) {
	var paylist []model.Paylist
	config.DB.Find(&paylist)

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
	err := config.DB.Model(&model.Paylist{}).Where("ID = ?", paylistID).Find(&paylist).Error

	if err != nil {
		util.CallErrorNotFound(c, "no paylist found!", err)
		return
	}
	util.CallSuccessOK(c, "success fetch single paylist", paylist)
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
		util.CallErrorNotFound(c, "paylist not found, make sure to specify the ID", nil)
		return
	}

	err := config.DB.Model(&paylist).Update(&updatedPaylist).Error
	if err != nil {
		util.CallServerError(c, "failed to update the paylist", err)
	}
	util.CallSuccessOK(c, "paylist successfully updated!", paylist)
}

// DeletePaylist remove a paylist
func DeletePaylist(c *gin.Context) {
	var paylist model.Paylist
	paylistID := c.Param("id")

	config.DB.First(&paylist, paylistID)

	if paylist.ID == 0 {
		util.CallErrorNotFound(c, "no paylist found!", nil)
		return
	}
	err := config.DB.Delete(&paylist).Error
	if err != nil {
		util.CallServerError(c, "failed to Delete Paylist", err)
		return
	}
	util.CallSuccessOK(c, "paylist deleted successfully!", nil)
}
