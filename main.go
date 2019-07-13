package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/ariebrainware/paylist-api/model"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "rob0ne:@_L0c4lDB@/paylist?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&model.Paylist{})
	db.AutoMigrate(&model.User{})
	fmt.Println("Schema migrated!!")
}

// createPaylist
func createPaylist(c *gin.Context) {
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

// fetchAllPaylist fetch all paylist
func fetchAllPaylist(c *gin.Context) {
	var paylist []model.Paylist

	db.Find(&paylist)

	if len(paylist) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No paylist found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": paylist})

}

// fetchSinglePaylist fetch a single paylist
func fetchSinglePaylist(c *gin.Context) {
	var paylist model.Paylist
	paylistID := c.Param("id")
	db.First(&model.Paylist{}, paylistID)

	if paylist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No paylist found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": paylist})
}

// updatePaylistss update a paylist
func updatePaylist(c *gin.Context) {
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

// delete Paylist remove a paylist
func deletePaylist(c *gin.Context) {
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

func main() {
	router := gin.Default()
	// CORS module
	corsConfig := cors.DefaultConfig()
	// corsConfig.AllowWildcard = true
	corsConfig.AllowOrigins = []string{"*", "127.0.0.1", "localhosts"}
	router.Use(cors.New(corsConfig))
	v1 := router.Group("/v1/paylist")
	{
		v1.GET("/", fetchAllPaylist)
		v1.GET("/:id", fetchSinglePaylist)
		v1.POST("/", createPaylist)
		v1.PUT("/:id", updatePaylist)
		v1.DELETE("/:id", deletePaylist)
	}
	router.Run()
}
