package main

import (
	"net/http"


	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	
	//"github.com/ariebrainware/paylist-api/model"
	"github.com/ariebrainware/paylist-api/endpoints"	
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
	// v1 := router.Group("/v1/paylist/")
	router.GET("/paylist", endpoints.FetchAllPaylist)
	router.GET("/paylist/:id", endpoints.FetchSinglePaylist)
	router.POST("/paylist", endpoints.CreatePaylist)
	router.PUT("/paylist/:id", endpoints.UpdatePaylist)
	router.DELETE("/paylist/:id", endpoints.DeletePaylist)
	router.Run(":3002")
}
