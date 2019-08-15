package endpoint

import (
	"testing"
	"bytes"
	"net/http"
	"net/http/httptest"
	"fmt"
	"github.com/gin-gonic/gin"
	//"gopkg.in/go-playground/assert.v1"
	
)

//test Func Login

func TestFetchSingleUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("GET", "/users/15", nil)
    if err != nil {
        fmt.Println(err)
	}
    router := gin.Default()
    router.GET("/users/:id", FetchSingleUser)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)
	//assert.Equal(t, resp.Code, 200)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestFetchAllUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("GET", "/users", nil)
    if err != nil {
        fmt.Println(err)
	}
    router := gin.Default()
    router.GET("/users", FetchAllUser)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)
	//assert.Equal(t, resp.Code, 200)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := bytes.NewBuffer([]byte("{\"ID\":\"18\",\"name\":\"offler d\",\"username\":\"offler\",\"email\":\"offler11@gmail.com\",\"password\":\"123\"}"))
	
	req, err := http.NewRequest("POST", "/users/signup", body)
	
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Create user failed with error %d.", err)
	}
	router := gin.Default()
	router.POST("/users/signup", CreateUser)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	
	if resp.Code != 201 {
		t.Errorf("/users/signup failed with error code %d.", resp.Code)
  	}
}