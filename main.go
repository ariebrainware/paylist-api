package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/ariebrainware/paylist-api/config"
	ep "github.com/ariebrainware/paylist-api/endpoint"
	"github.com/ariebrainware/paylist-api/util"
)

type endpoint struct {
	Method      string
	URL         string
	Description string
}

func main() {
	config.Conf()
	router := gin.Default()
	router.Use(cors.Default())
	listEndpoint := []endpoint{
		//Paylist Endpoint
		{Method: "GET", URL: "/paylist", Description: "Get/Fetch All Paylist Data"},
		{Method: "GET", URL: "/paylist/:id", Description: "Get/Fetch Single Paylist Data by ID"},
		{Method: "POST", URL: "/paylist", Description: "Create/Insert User-Paylist Data"},
		{Method: "PUT", URL: "/paylist/:id", Description: "Edit/Update Paylist Data by ID"},
		{Method: "DELETE", URL: "/paylist/:id", Description: "Delete User-Paylist Data by ID"},
		//User Endpoint
		{Method: "GET", URL: "/users", Description: "Get/Fetch All User Data"},
		{Method: "GET", URL: "/users/:id", Description: "Get/Fetch Single User Data by ID"},
		{Method: "POST", URL: "/users", Description: "Sign Up"},
		{Method: "POST", URL: "/users/signin", Description: "Sign In"},
		{Method: "PUT", URL: "/users/:id", Description: "Edit/Update User Data by ID"},
		{Method: "DELETE", URL: "/users/:id", Description: "Delete User Data by ID"},
		{Method: "PUT", URL: "/user-paylist/:id", Description: "Update User-Paylist by ID"},
		{Method: "GET", URL: "/user/signout", Description: "Sign Out / Logout"},
		{Method: "POST", URL: "/users/refresh-token", Description: "Refresh Expired Token"},
	}

	router.GET("/", func(c *gin.Context) {
		util.CallSuccessOK(c, "Paylist-API available endpoint", listEndpoint)
	})
	
	router.GET("/paylist", ep.Auth, ep.FetchAllPaylist)
	router.GET("/paylist/:id", ep.Auth, ep.FetchSinglePaylist)
	router.POST("/paylist", ep.Auth, ep.CreateUserPaylist)
	router.PUT("/status/:id", ep.Auth, ep.UpdateUserPaylistStatus)
	router.PUT("/paylist/:id", ep.Auth, ep.UpdateUserPaylist)
	router.DELETE("/paylist/:id", ep.Auth, ep.DeleteUserPaylist)

	router.GET("/user/:id", ep.Auth, ep.FetchSingleUser)
	router.GET("/users", ep.Auth, ep.FetchAllUser)
	router.GET("/users/signout", ep.Logout)
	router.POST("/user/signin", ep.Login)
	router.POST("/user/signup", ep.CreateUser)
	router.POST("/addsaldo", ep.Auth, ep.AddBalance)
	router.PUT("/user/:id", ep.Auth, ep.UpdateUser)
	router.PUT("/editpassword/:id", ep.Auth, ep.EditPassword)
	router.DELETE("/user/:id", ep.Auth, ep.DeleteUser)
	router.Run(":8000")
}
