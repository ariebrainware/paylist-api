package main

import (
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
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
	config.LoadConfiguration()
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://0d01b76a59934b6887a0323eed450b0f@o326252.ingest.sentry.io/5553941",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

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
		{Method: "PUT", URL: "/editpassword/:id", Description: "Edit user password"},
		{Method: "PUT", URL: "/status/:id", Description: "Update user-paylist status"},
		{Method: "POST", URL: "/addsaldo", Description: "Add User Saldo"},
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
	router.GET("/income", ep.Auth, ep.FetchAllIncome)
	router.PUT("/income/:id", ep.Auth, ep.UpdateIncome)
	router.DELETE("/income/:id", ep.Auth, ep.DeleteIncome)
	err = router.Run(fmt.Sprintf(":%d", config.Conf.Port))
	if err != nil {
		sentry.CaptureException(err)
		fmt.Println(err)
		return
	}
}
