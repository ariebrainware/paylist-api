package endpoint

import (
	"testing"
	"bytes"
	"net/http"
	"net/http/httptest"
	"fmt"
	"github.com/gin-gonic/gin"	
)
// fun TestFetchSingleUser Functional Testing for FetchSingleUser
func TestFetchSingleUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("GET", "/users/31", nil)
    if err != nil {
        fmt.Println(err)
	}
    router := gin.Default()
    router.GET("/users/:id", FetchSingleUser)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("router returned wrong status code: got %v want %v",
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
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("router returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// fun TestCreateUser Functional Testing for CreateUser
func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := bytes.NewBuffer([]byte("{\"username\":\"offler\",\"password\":\"123\",\"email\":\"offler11@gmail.com\",\"name\":\"offler d\"}"))
	req, err := http.NewRequest("POST", "/users/signup", body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	router := gin.Default()
	resp := httptest.NewRecorder()
	router.POST("/users/signup", CreateUser)
	router.ServeHTTP(resp, req)
	if resp.Code != 201 {
		t.Errorf("router returned wrong status code: got %d", resp.Code)
	}
}

// fun TestDeleteUser Functional Testing for DeleteUser
func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("DELETE", "/users/33", nil)
	if err != nil {
		fmt.Println(err)
	}
	router := gin.Default()
	router.DELETE("/users/:id", DeleteUser)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("router returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
// fun TestUpdateUser Functional Testing for UpdateUser
func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := bytes.NewBuffer([]byte("{\"username\":\"offler\",\"password\":\"123\",\"email\":\"offler11@gmail.com\",\"name\":\"offler d\"}"))
	req, err := http.NewRequest("PUT", "/users/31", body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	router := gin.Default()
	resp := httptest.NewRecorder()
	router.PUT("/users/:id", UpdateUser)
	router.ServeHTTP(resp, req)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("router returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// fun TestFetchSinglePaylist Functional Testing for FetchSinglePaylist
func TestFetchSinglePaylist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("GET", "/paylist/1", nil)
    if err != nil {
        fmt.Println(err)
	}
    router := gin.Default()
    router.GET("paylist/:id", FetchSinglePaylist)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)
	//assert.Equal(t, resp.Code, 200)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("router returned wrong status code: got %v want %v",
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
		t.Errorf("router returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// fun TestCreate Functional Testing for Create User
func TestCreatePaylist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var body = []byte(`{"name":"powerbank","amount":"500000"}`)
	req, err := http.NewRequest("POST", "/paylist", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/paylist", CreatePaylist)
	router.ServeHTTP(resp, req)
	if resp.Code != 201 {
		t.Errorf("router returned wrong status code: got %d",
			resp.Code)
	}
}

// fun TestDeletePaylist Functional Testing for Delete Paylist
func TestDeletePaylist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("DELETE", "/paylist/18", nil)
	if err != nil {
		fmt.Println(err)
	}
	router := gin.Default()
	router.DELETE("/paylist/:id", DeletePaylist)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("router returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// fun TestUpdatePaylist Functional Testing for UpdatePaylist
func TestUpdatePaylist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := bytes.NewBuffer([]byte("{\"name\":\"ayam\",\"amount\": 20000}"))
	req, err := http.NewRequest("PUT", "/paylist/3", body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	router := gin.Default()
	router.PUT("/paylist/:id", UpdatePaylist)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

//func TestLogin Functional Testing for Login User
func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := bytes.NewBuffer([]byte("{\"username\":\"jenkins\",\"password\":\"jenkins123\"}"))
	req, err := http.NewRequest("POST", "/users/signin", body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	router := gin.Default()
	resp := httptest.NewRecorder()
	router.POST("/users/signin", Login)
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("router returned wrong status code: got %d", resp.Code)
	}
}
