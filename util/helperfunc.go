package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Contains function is to check item whether is exist or not in a list and will return bool
func Contains(d string, dl []string) bool {
	for _, v := range dl {
		if v == d {
			return true
		}
	}
	return false
}

// CallErrorNotFound is for return API response not found
func CallErrorNotFound(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"error":   err,
		"msg":     msg,
		"data":    map[string]interface{}{},
	})
}

// CallUserError is for return error from user side
func CallUserError(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   err,
		"msg":     msg,
		"data":    map[string]interface{}{},
	})
}

// CallServerError is for return API response server error
func CallServerError(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   err,
		"msg":     msg,
		"data":    map[string]interface{}{},
	})
}

// CallSuccessOK is for return API response with status code 200, you need to specify msg, and data as function parameter
func CallSuccessOK(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"error":   "",
		"msg":     msg,
		"data":    data,
	})
}
