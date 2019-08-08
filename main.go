package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	ep "github.com/ariebrainware/paylist-api/endpoint"
)

type endpoint struct {
	Method string
	URL    string
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	listEndpoint := []endpoint{
		{Method: "GET", URL: "/paylist"},
		{Method: "POST", URL: "/paylist"},
		{Method: "PUT", URL: "/paylist/:id"},
		{Method: "DELETE", URL: "/paylist/:id"},
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  http.StatusOK,
			"message": "Paylist-API available endpoint",
			"data":    listEndpoint,
		})
	})
	v1 := router.Group("/v1/paylist/")
	v1.GET("/paylist", ep.FetchAllPaylist)
	v1.GET("/paylist/:id", ep.FetchSinglePaylist)
	v1.POST("/paylist", ep.CreatePaylist)
	v1.PUT("/paylist/:id", ep.UpdatePaylist)
	v1.DELETE("/paylist/:id", ep.DeletePaylist)
	v1.GET("/users/:id", ep.FetchSingleUser)
	v1.GET("/users", ep.FetchUser)
	v1.POST("/users/signin", ep.Login)
	v1.POST("/users/signup", ep.CreateUser)
	v1.PUT("/users/:id", ep.UpdateUser)
	v1.DELETE("/users/:id", ep.DeleteUser)
	router.Run(":3002")
}
