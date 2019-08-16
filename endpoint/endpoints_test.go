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
// fun TestFetchSingleUser Functional Testing for FetchSingleUser
func TestFetchSingleUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("GET", "/users/16", nil)
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

// fun TestFetchAllUser Functional Testing for Fetch All User
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

// fun TestCreateUser Functional Testing for CreateUser
func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := bytes.NewBuffer([]byte("{\"email\":\"offler11@gmail.com\",\"name\":\"offler d\",\"username\":\"offler\",\"password\":\"123\"}"))
	
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

// fun TestDeleteUser Functional Testing for DeleteUser
func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("DELETE", "/users/29", nil)
	if err != nil {
		t.Errorf("failed with error code %d", err)
	}
	router := gin.Default()
	router.DELETE("/users/:id", DeleteUser)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Errorf("Failed delete user with error code %d", resp.Code)
	}
}
// fun TestUpdateUser Functional Testing for UpdateUser
func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := bytes.NewBuffer([]byte("{\"username\":\"offler\",\"password\":\"123\",\"email\":\"offler11@gmail.com\",\"name\":\"offler d\"}"))
	req, err := http.NewRequest("PUT", "/users/31", body)
	//req.Header.Set("Content-Type", "application/json")
	
	if err != nil {
		t.Errorf("failed update with error code %d", err)
	}
	
	router := gin.Default()
	resp := httptest.NewRecorder()
	router.PUT("/users/:id", UpdateUser)
	router.ServeHTTP(resp, req)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// fun TestFetchSinglePaylist Functional Testing for FetchSinglePaylist
func TestFetchSinglePaylist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("GET", "/paylist/2", nil)
    if err != nil {
        fmt.Println(err)
	}
    router := gin.Default()
    router.GET("paylist/:id", FetchSinglePaylist)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)
	//assert.Equal(t, resp.Code, 200)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// fun TestFetchAllPaylist Functional Testing for FetchAllPaylist
func TestFetchAllPaylist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("GET", "/paylist", nil)
    if err != nil {
        fmt.Println(err)
	}
    router := gin.Default()
    router.GET("/paylist", FetchAllPaylist)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)
	//assert.Equal(t, resp.Code, 200)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// fun TestCreate Functional Testing for Create User
func TestCreatePaylist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := bytes.NewBuffer([]byte(`{"name":"powerbank","amount":500000}`))
	
	req, err := http.NewRequest("POST", "/paylist", body)
	
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Create paylist failed with error %d.", err)
	}
	router := gin.Default()
	router.POST("/paylist", CreatePaylist)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	
	if resp.Code != 201 {
		t.Errorf("/users/signup failed with error code %d.", resp.Code)
  	}
}

// fun TestDeletePaylist Functional Testing for Delete Paylist
func TestDeletePaylist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("DELETE", "/paylist/5", nil)
	if err != nil {
		t.Errorf("failed with error code %d", err)
	}
	router := gin.Default()
	router.DELETE("/paylist/:id", DeletePaylist)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != 200 {
		t.Errorf("Failed delete user with error code %d", resp.Code)
	}
}

// fun TestUpdatePaylist Functional Testing for UpdatePaylist
func TestUpdatePaylist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := bytes.NewBuffer([]byte("{\"name\":\"ayam\",\"amount\": 20000}"))
	req, err := http.NewRequest("PUT", "/paylist/2", body)
	//req.Header.Set("Content-Type", "application/json")
	
	if err != nil {
		t.Errorf("failed update with error code %d", err)
	}
	
	router := gin.Default()
	resp := httptest.NewRecorder()
	router.PUT("/paylist/:id", UpdatePaylist)
	router.ServeHTTP(resp, req)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}


func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
}