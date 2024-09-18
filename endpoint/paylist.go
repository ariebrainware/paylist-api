package endpoint

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/ariebrainware/paylist-api/config"
	"github.com/ariebrainware/paylist-api/model"
	"github.com/ariebrainware/paylist-api/util"
)

// User stuct for parse token
type User struct {
	Username string
	jwt.StandardClaims
}

// FetchAllPaylist handles the request to fetch all paylists for a user.
// It retrieves the user's token from the request context, parses it to get the username,
// and then queries the database for paylists associated with that username.
// If the token is invalid or there are no paylists found, appropriate error responses are returned.
// On success, it returns the list of paylists.
//
// @param c *gin.Context - The Gin context for the request.
//
// @response 200 - Successfully fetched all paylists.
// @response 400 - Failed to parse the token or no paylists found.
func FetchAllPaylist(c *gin.Context) {
	var paylist []model.Paylist
	tk, err := parseToken(c)
	if err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "fail to parse the token, make sure token and signature is valid",
			Err: err,
		})
		return
	}
	username := tk.Username
	errf := config.DB.Model(&paylist).Where("username = ?", username).Find(&paylist).Error
	if errf != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "no paylist found!",
			Err: errf,
		})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "fetched all paylist",
		Data: paylist,
	})
}

// FetchSinglePaylist handles the request to fetch a single paylist by its ID.
// It expects an Authorization token in the request header and validates it.
// If the token is valid, it retrieves the paylist associated with the given ID
// and the username extracted from the token. If the paylist is found, it returns
// the paylist data with a success message. If any error occurs, it returns an
// appropriate error message.
//
// @Summary Fetch a single paylist
// @Description Fetch a single paylist by its ID and the username from the token
// @Tags paylist
// @Accept json
// @Produce json
// @Param id path string true "Paylist ID"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} model.Paylist "success fetch single paylist"
// @Failure 400 {object} util.APIErrorParams "fail to parse the token, make sure token is valid"
// @Failure 404 {object} util.APIErrorParams "no paylist found!"
// @Router /paylist/{id} [get]
func FetchSinglePaylist(c *gin.Context) {
	var paylist model.Paylist
	paylistID := c.Param("id")
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
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
	errf := config.DB.Model(&model.Paylist{}).Where("ID = ? and username = ?", paylistID, username).Find(&paylist).Error
	if errf != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "no paylist found!",
			Err: errf,
		})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "success fetch single paylist",
		Data: paylist,
	})
}

// CreateUserPaylist handles the creation of a new paylist item for a user.
// It performs the following steps:
// 1. Parses the JWT token from the Authorization header to extract the username.
// 2. Decreases the user's balance by the amount specified in the request.
// 3. Creates a new paylist item with the provided details (name, amount, due date).
// 4. Saves the new paylist item to the database.
// 5. Returns a success response with the ID of the created paylist item.
//
// Parameters:
// - c: *gin.Context - The Gin context, which contains the request and response objects.
//
// The function expects the following form data in the request:
// - name: string - The name of the paylist item.
// - amount: string - The amount to be deducted from the user's balance.
// - due_date: string - The due date of the paylist item in the format "YYYY-MM-DD".
//
// The function returns a JSON response with the status of the operation.
func CreateUserPaylist(c *gin.Context) {
	users := model.User{}
	tk := User{}

	// Parse the payload from token
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
	// Decrease user balance
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	dueDate, _ := time.Parse("2006-01-02", c.PostForm("due_date"))
	err = config.DB.Model(&users).Where("username  = ?", username).First(&users).Error
	if err != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "can't select balance",
			Err: err,
		})
		return
	}
	finalAmount := users.Balance - amount
	config.DB.Model(&users).Where("username = ?", username).Update("balance", finalAmount)
	paylist := model.Paylist{
		Name:      c.PostForm("name"),
		Amount:    amount,
		Username:  username,
		DueDate:   dueDate.Format("2006-01-02"),
		Completed: false,
	}

	// Save paylist
	err = config.DB.Save(&paylist).Error
	if err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "fail to create paylist",
			Err: err,
		})
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "paylist item created successfully!",
		Data: paylist.ID,
	})

}

// UpdateUserPaylist updates the paylist of a user based on the provided context.
// It performs the following steps:
// 1. Parses the token from the context to get the username.
// 2. Retrieves the user and paylist based on the username and paylist ID.
// 3. Validates if the user is authorized to update the paylist.
// 4. Updates the paylist with the new data provided in the request.
// 5. Updates the user's balance based on the changes in the paylist amount.
// 6. Returns appropriate error messages if any step fails, or a success message if the update is successful.
//
// Parameters:
// - c: *gin.Context - The context of the request, containing parameters and form data.
//
// Possible Errors:
// - If the token is invalid or cannot be parsed.
// - If the user or paylist is not found.
// - If the user is not authorized to update the paylist.
// - If there is an error updating the paylist or user balance.
//
// Success Response:
// - Returns a success message indicating the paylist was successfully updated.
func UpdateUserPaylist(c *gin.Context) {
	tk, err := parseToken(c)
	if err != nil {
		util.CallUserError(c, util.APIErrorParams{
			Msg: "fail to parse the token, make sure the token is valid",
			Err: err,
		})
		return
	}

	username := tk.Username
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := findUserByUsername(username)
	if err != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "user not found",
			Err: err,
		})
		return
	}

	paylist, err := findPaylistByIDAndUsername(id, username)
	if err != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "no paylist found",
			Err: err,
		})
		return
	}

	if tk.Username != paylist.Username {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "not authorized",
			Err: nil,
		})
		return
	}

	firstAmount := paylist.Amount
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	updatedPaylist := model.Paylist{
		Name:    c.PostForm("name"),
		Amount:  amount,
		DueDate: c.PostForm("due_date"),
	}

	if err := updatePaylist(username, updatedPaylist); err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "something error when try to update paylist",
			Err: err,
		})
		return
	}

	if err := updateUserBalance(username, firstAmount, amount, user.Balance); err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "something error when try to update user balance",
			Err: err,
		})
		return
	}

	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "paylist successfully updated!",
		Data: paylist,
	})
}

func parseToken(c *gin.Context) (*User, error) {
	tk := User{}
	tokenString := c.GetHeader("Authorization")
	_, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprint(config.Conf.JWTSignature)), nil
	})
	if err != nil {
		return nil, err
	}
	return &tk, nil
}

func findUserByUsername(username string) (*model.User, error) {
	user := model.User{}
	err := config.DB.Model(&user).Select("balance").Where("username = ?", username).Find(&user).Error
	return &user, err
}

func findPaylistByIDAndUsername(id int, username string) (*model.Paylist, error) {
	paylist := model.Paylist{}
	err := config.DB.Model(&paylist).Where("ID = ? AND username = ?", id, username).First(&paylist).Error
	return &paylist, err
}

func updatePaylist(username string, updatedPaylist model.Paylist) error {
	return config.DB.Model(&model.Paylist{}).Where("username = ?", username).Update(&updatedPaylist).Error
}

func updateUserBalance(username string, firstAmount, amount, balance int) error {
	return config.DB.Model(&model.User{}).Where("username = ?", username).Update("balance", (firstAmount-amount)+balance).Error
}

// UpdateUserPaylistStatus updates the status of a user's paylist.
//
// This function performs the following steps:
// 1. Parses the token from the request context.
// 2. Retrieves the username from the token.
// 3. Converts the "id" parameter from the request to an integer.
// 4. Finds the paylist by ID and username.
// 5. Finds the user by username.
// 6. Updates the paylist status based on the user's balance.
// 7. Returns a success response if the update is successful.
//
// If any step fails, an appropriate error response is returned.
//
// Parameters:
// - c: The Gin context containing the request and response objects.
//
// Responses:
// - 200 OK: If the paylist status is successfully updated.
// - 400 Bad Request: If the token parsing fails.
// - 404 Not Found: If the paylist or user is not found.
// - 500 Internal Server Error: If updating the paylist status fails.
func UpdateUserPaylistStatus(c *gin.Context) {
	tk, err := parseToken(c)
	if err != nil {
		util.CallUserError(c, util.APIErrorParams{
			Msg: "fail to parse the token, make sure the token is valid",
			Err: err,
		})
		return
	}

	username := tk.Username
	id, _ := strconv.Atoi(c.Param("id"))

	paylist, err := findPaylistByIDAndUsername(id, username)
	if err != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "no paylist found",
			Err: err,
		})
		return
	}

	user, err := findUserByUsername(username)
	if err != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "user not found",
			Err: err,
		})
		return
	}

	if err := updatePaylistStatus(paylist, user.Balance); err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "fail to update paylist status",
			Err: err,
		})
		return
	}

	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "successfully update user paylist",
		Data: paylist.Completed,
	})
}

func updatePaylistStatus(paylist *model.Paylist, balance int) error {
	if balance >= 0 && !paylist.Completed {
		paylist.Completed = true
	} else if balance < 0 && !paylist.Completed {
		paylist.Completed = false
	}
	return config.DB.Model(&paylist).Update(paylist).Error
}

func DeleteUserPaylist(c *gin.Context) {
	paylistID, _ := strconv.Atoi(c.Param("id"))
	if paylistID == 0 {
		util.CallUserError(c, util.APIErrorParams{
			Msg: "please specify paylist id",
			Err: nil,
		})
		return
	}

	tk, err := parseToken(c)
	if err != nil {
		util.CallUserError(c, util.APIErrorParams{
			Msg: "fail to parse the token, make sure the token is valid",
			Err: err,
		})
		return
	}

	paylist, err := findPaylistByID(paylistID)
	if err != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "no paylist found!",
			Err: err,
		})
		return
	}

	if tk.Username != paylist.Username {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "user not authorized",
			Err: nil,
		})
		return
	}

	user, err := findUserByUsername(tk.Username)
	if err != nil {
		util.CallErrorNotFound(c, util.APIErrorParams{
			Msg: "can't select balance",
			Err: err,
		})
		return
	}

	if err := updateUserBalanceOnDelete(paylist, user); err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "fail to update the user balance",
			Err: err,
		})
		return
	}

	if err := deletePaylist(paylistID, tk.Username); err != nil {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "paylist fail to delete",
			Err: err,
		})
		return
	}

	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "paylist successfully deleted!",
		Data: nil,
	})
}

func findPaylistByID(paylistID int) (*model.Paylist, error) {
	paylist := &model.Paylist{}
	err := config.DB.Where("ID = ?", paylistID).First(&paylist).Error
	return paylist, err
}

func updateUserBalanceOnDelete(paylist *model.Paylist, user *model.User) error {
	if !paylist.Completed {
		newBalance := paylist.Amount + user.Balance
		return config.DB.Table("users").Where("username = ?", user.Username).Update("balance", newBalance).Error
	}
	return nil
}

func deletePaylist(paylistID int, username string) error {
	return config.DB.Where("ID = ? and username = ?", paylistID, username).Delete(&model.Paylist{}).Error
}

// FilterPaylist filters the paylist based on the provided "created_at" parameter and the username extracted from the JWT token.
// It expects the "Authorization" header to contain a valid JWT token.
//
// @param c *gin.Context - The Gin context which provides request and response handling.
//
// The function performs the following steps:
// 1. Retrieves the "created_at" parameter from the request context.
// 2. If the "created_at" parameter is not provided, it returns an error response indicating the missing parameter.
// 3. Extracts the JWT token from the "Authorization" header and parses it to extract the username.
// 4. Queries the database to filter the paylist records based on the username and the month of the "created_at" parameter.
// 5. If the query is successful, it returns the filtered paylist in the response.
// 6. If any error occurs during the process, it returns an appropriate error response.
func FilterPaylist(c *gin.Context) {
	paylist := model.Paylist{}
	//param,_ := strconv.ParseInt(c.Param("month"),10,64)
	//month := paylist.CreatedAt.Month()
	param := c.Param("created_at")
	if param != "" {
		util.CallServerError(c, util.APIErrorParams{
			Msg: "please specify the filter parameter",
			Err: nil,
		})
		return
	}
	tk := User{}

	// Extract token payload
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
	}
	username := tk.Username

	err = config.DB.Model(&paylist).Where("username = ? and MONTH(created_at) = ?", username, param).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	util.CallSuccessOK(c, util.APISuccessParams{
		Msg:  "success paylist",
		Data: paylist,
	})
}
