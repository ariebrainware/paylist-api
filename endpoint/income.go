package endpoint

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go" //Used to sign and verify JWT tokens
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/ariebrainware/paylist-api/config"
	"github.com/ariebrainware/paylist-api/model"
	"github.com/ariebrainware/paylist-api/util"
)

// Token is a struct for token model

type Income struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Username  string
	Income    int
}

// FetchAllIncome handles the HTTP request to fetch all income records for a user.
// It extracts the JWT token from the Authorization header, parses it to get the username,
// and retrieves the income records associated with that username from the database.
// If no income records are found, it returns a 404 error. Otherwise, it returns the income records.
//
// @Summary Fetch all income records
// @Description Fetch all income records for the authenticated user
// @Tags income
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []Income "List of income records"
// @Failure 400 {object} util.APIErrorParams "Invalid request"
// @Failure 401 {object} util.APIErrorParams "Unauthorized"
// @Failure 404 {object} util.APIErrorParams "No income records found"
// @Failure 500 {object} util.APIErrorParams "Internal server error"
// @Router /income [get]
func FetchAllIncome(c *gin.Context) {
	var income []model.Income
	var inc []Income
	tk := User{}
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprint(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, util.APIErrorParams{
			Msg: "fail to parse the token, make sure token is valid",
			Err: err,
		})
		return
	}
	username := tk.Username
	config.DB.Model(&income).Where("username = ? ", username).Find(&income)

	if len(income) <= 0 {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "No Income Found!",
			Err: nil,
		})
		return
	}
	for _, item := range income {
		inc = append(inc, Income{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			DeletedAt: item.DeletedAt,
			Username:  item.Username,
			Income:    item.Income,
		})
	}
	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "Fetch All Income Data",
		Data: inc,
	})
}

// UpdateIncome handles the update of an income record for a user.
//
// This function performs the following steps:
// 1. Parses the JWT token from the "Authorization" header to authenticate the user.
// 2. Checks the user's balance to ensure they have sufficient funds.
// 3. Finds the income record based on the provided ID and username.
// 4. Validates that the authenticated user is authorized to update the income record.
// 5. Updates the income record in the database.
// 6. Updates the user's balance in the database.
// 7. Returns a success response if the update is successful, or an error response if any step fails.
//
// Parameters:
// - c: *gin.Context - The Gin context, which provides request and response handling.
//
// The function expects the following in the request:
// - Authorization header containing a valid JWT token.
// - URL parameter "id" representing the income record ID.
// - Form data "income" representing the new income amount.
//
// Responses:
// - 200 OK with a success message and the updated income data if the update is successful.
// - 400 Bad Request with an error message if the token is invalid or any required data is missing.
// - 500 Internal Server Error with an error message if any database operation fails.
func UpdateIncome(c *gin.Context) {
	income := model.Income{}
	user := model.User{}
	tk := User{}

	// Parse the token payload
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprint(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallUserError(c, util.APIErrorParams{
			Msg: "fail to parse the token, make sure the token is valid",
			Err: err,
		})
		return
	}

	// Check User balance
	username := tk.Username
	id, _ := strconv.Atoi(c.Param("id"))
	if !checkUserBalance(c, username, &user) {
		return
	}
	if !findIncome(c, id, username, &income) {
		return
	}

	firstIncome := income.Income
	inc, _ := strconv.Atoi(c.PostForm("income"))
	updatedIncome := model.Income{
		Income: inc,
	}
	if tk.Username != income.Username {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "not authorized",
			Err: nil,
		})
		return
	}
	// Update paylist
	if err = config.DB.Model(&income).Where("username = ?", username).Update(&updatedIncome).Error; err != nil {
		fmt.Println(err)
		util.CallServerError(c, util.APIErrorParams{
			Msg: "something error when try to update paylist",
			Err: err,
		})
		return
	}

	// Update user balance
	err = config.DB.Model(&user).Where("username = ?", username).Update("balance", (inc+user.Balance)-firstIncome).Error
	if err != nil {
		fmt.Println(err)
		util.CallServerError(c, util.APIErrorParams{
			Msg: "something error when try to update user balance",
			Err: err,
		})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "income successfully updated!",
		Data: income,
	})
}

func checkUserBalance(c *gin.Context, username string, user *model.User) bool {
	if err := config.DB.Model(user).Select("balance").Where("username = ?", username).Find(user).Error; err != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "user not found",
			Err: err,
		})
		return false
	}
	return true
}

func findIncome(c *gin.Context, id int, username string, income *model.Income) bool {
	if err := config.DB.Model(income).Where("ID = ? AND username = ?", id, username).First(income).Error; err != nil || err == gorm.ErrRecordNotFound {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "no income found",
			Err: nil,
		})
		return false
	}
	return true
}

// DeleteIncome handles the deletion of an income record for a user.
// It expects an "id" parameter in the URL and an "Authorization" header with a valid JWT token.
//
// The function performs the following steps:
// 1. Parses the JWT token from the "Authorization" header to extract the username.
// 2. Checks if the income record with the given ID exists in the database.
// 3. Verifies that the income record belongs to the user extracted from the token.
// 4. Updates the user's balance by subtracting the income amount.
// 5. Deletes the income record from the database.
//
// If any step fails, an appropriate error response is sent to the client.
//
// Parameters:
// - c: *gin.Context - The context for the request, which includes parameters, headers, and other request data.
//
// Responses:
// - 200 OK: If the income record is successfully deleted.
// - 400 Bad Request: If the token is invalid or missing.
// - 404 Not Found: If the income record does not exist or does not belong to the user.
// - 500 Internal Server Error: If there is an error updating the balance or deleting the income record.
func DeleteIncome(c *gin.Context) {
	var users model.User
	var inc model.Income
	ID := c.Param("id")
	tk := User{}
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprintf(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, util.APIErrorParams{
			Msg: "fail to parse the token, make sure token is valid",
			Err: err,
		})
		return
	}
	username := tk.Username
	if err = config.DB.Where("ID = ?", ID).Find(&inc).Error; err != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "no income data found!",
			Err: err,
		})
		c.Abort()
		return
	}
	config.DB.Model(&inc).Select("income").Where("username = ? and ID = ? ", username, ID).Find(&inc)
	if inc.ID == 0 {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "can't select income",
			Err: nil,
		})
		return
	}
	config.DB.Model(&users).Select("balance").Where("username = ?", username).Find(&users)
	if inc.ID == 0 {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "can't select income",
			Err: nil,
		})
		return
	}
	income := inc.Income
	fmt.Println("inc", inc.Income)
	err = config.DB.Model(&users).Where("username = ?", username).Update("balance", users.Balance-income).Error
	if err != nil {
		fmt.Println(err)
		util.CallServerError(c, util.APIErrorParams{
			Msg: "failed update income",
			Err: err,
		})
		return
	}
	config.DB.Model(&inc).Where("username = ?", username).Delete(&inc)
	if tk.Username != inc.Username {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "not authorized",
			Err: nil,
		})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "user delete successfully!",
		Data: nil,
	})
}
