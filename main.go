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
	v1 := router.Group("/v1/paylist/")
	v1.GET("/paylist", ep.Auth, ep.SignOut, ep.FetchAllPaylist)
	v1.GET("/paylist/:id", ep.Auth, ep.SignOut, ep.FetchSinglePaylist)
	v1.POST("/paylist", ep.Auth, ep.SignOut, ep.CreateUserPaylist)
	v1.PUT("/status/:id", ep.Auth, ep.SignOut, ep.UpdateUserPaylistStatus)
	v1.PUT("/paylist/:id", ep.Auth, ep.SignOut, ep.UpdateUserPaylist)
	v1.DELETE("/paylist/:id", ep.Auth, ep.DeleteUserPaylist)

	v1.GET("/user/:id", ep.SignOut, ep.Auth, ep.FetchSingleUser)
	v1.GET("/users", ep.SignOut, ep.Auth, ep.FetchAllUser)
	v1.GET("/users/signout", ep.SignOut, ep.Logout)
	v1.POST("/user/signin", ep.Login)
	v1.POST("/user/signup", ep.CreateUser)
	v1.POST("/user/refresh-token", ep.RefreshToken)
	v1.POST("/addsaldo", ep.Auth, ep.SignOut, ep.AddBalance)
	v1.PUT("/user/:id", ep.SignOut, ep.Auth, ep.UpdateUser)
	v1.DELETE("/user/:id", ep.SignOut, ep.Auth, ep.DeleteUser)
	router.Run(":8000")
}
