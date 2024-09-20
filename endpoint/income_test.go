package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ariebrainware/paylist-api/config"
	"github.com/ariebrainware/paylist-api/model"
	"github.com/ariebrainware/paylist-api/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllIncome(t *testing.T) {
	// Setup Gin router and context
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/income", FetchAllIncome)

	// Mock data
	config.Conf.JWTSignature = "test_signature"
	mockToken := "Bearer test_token"
	mockUsername := "test_user"
	mockIncome := []model.Income{
		{ID: 1, Username: mockUsername, Income: 1000},
		{ID: 2, Username: mockUsername, Income: 2000},
	}

	// Mock database
	config.DB = util.MockDB()
	defer config.DB.Close()
	config.DB.Create(&mockIncome)

	// Test cases
	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid token with income records",
			token:          mockToken,
			expectedStatus: http.StatusOK,
			expectedBody:   `"Msg":"Fetch All Income Data"`,
		},
		{
			name:           "Invalid token",
			token:          "Bearer invalid_token",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"Msg":"fail to parse the token, make sure token is valid"`,
		},
		{
			name:           "No income records found",
			token:          "Bearer no_income_token",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"Msg":"No Income Found!"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req, _ := http.NewRequest(http.MethodGet, "/income", nil)
			req.Header.Set("Authorization", tt.token)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check the response body
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}
