package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	ep "github.com/ariebrainware/paylist-api/endpoint"
)

type endpoint struct {
	Method      string
	URL         string
	Description string
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	listEndpoint := []endpoint{
		//Paylist Endpoint
		{Method: "GET", URL: "/paylist", Description: "Get/Fetch All Paylist Data"},
		{Method: "GET", URL: "/paylist/:id", Description: "Get/Fetch Single Paylist Data by ID"},
		{Method: "POST", URL: "/paylist", Description: "Create Paylist Data"},
		{Method: "PUT", URL: "/paylist/:id", Description: "Edit/Update Paylist Data by ID"},
		{Method: "DELETE", URL: "/paylist/:id", Description: "Delete Paylist Data by ID"},
		//User Endpoint
		{Method: "GET", URL: "/users", Description: "Get/Fetch All User Data"},
		{Method: "GET", URL: "/users/:id", Description: "Get/Fetch Single User Data by ID"},
		{Method: "POST", URL: "/users", Description: "Sign Up"},
		{Method: "POST", URL: "/users/signin", Description: "Sign In"},
		{Method: "PUT", URL: "/users/:id", Description: "Edit/Update User Data by ID"},
		{Method: "DELETE", URL: "/users/:id", Description: "Delete User Data by ID"},
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  http.StatusOK,
			"message": "Paylist-API available endpoint",
			"data":    listEndpoint,
		})
	})
	v1 := router.Group("/v1/paylist/")
	v1.GET("/paylist", ep.Auth, ep.FetchAllPaylist)
	v1.GET("/paylist/:id", ep.Auth, ep.FetchSinglePaylist)
	v1.POST("/paylist", ep.Auth, ep.CreatePaylist)
	v1.PUT("/paylist/:id", ep.Auth, ep.UpdatePaylist)
	v1.DELETE("/paylist/:id", ep.Auth, ep.DeletePaylist)
	v1.GET("/users/:id", ep.Auth, ep.FetchSingleUser)
	v1.GET("/users", ep.Auth, ep.FetchAllUser)
	v1.POST("/users/signin", ep.Login)
	v1.POST("/users/signup", ep.CreateUser)
	v1.PUT("/users/:id",ep.Auth, ep.UpdateUser)
	v1.DELETE("/users/:id",ep.Auth, ep.DeleteUser)
	router.Run(":3002")
}
