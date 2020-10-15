package endpoint

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/ariebrainware/paylist-api/config"
	"github.com/ariebrainware/paylist-api/model"
	"github.com/ariebrainware/paylist-api/util"
)

//User stuct for parse token
type User struct {
	Username string
	jwt.StandardClaims
}

//FetchAllPaylist Fetch All Paylist
func FetchAllPaylist(c *gin.Context) {
	var paylist []model.Paylist
	tk := User{}
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprint(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallServerError(c, "fail to parse the token, make sure token and signature is valid", err)
		return
	}
	username := tk.Username
	errf := config.DB.Model(&paylist).Where("username = ?", username).Find(&paylist).Error
	if errf != nil {
		util.CallErrorNotFound(c, "no paylist found!", errf)
		return
	}
	util.CallSuccessOK(c, "fetched all paylist", paylist)
}

//FetchSinglePaylist fetch a single paylist
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
		util.CallServerError(c, "fail to parse the token, make sure token is valid", err)
		return
	}
	username := tk.Username
	errf := config.DB.Model(&model.Paylist{}).Where("ID = ? and username = ?", paylistID, username).Find(&paylist).Error
	if errf != nil {
		util.CallErrorNotFound(c, "no paylist found!", errf)
		return
	}
	util.CallSuccessOK(c, "success fetch single paylist", paylist)
}

// CreateUserPaylist function to create new paylist
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
		util.CallServerError(c, "fail to parse the token, make sure token is valid", err)
		return
	}
	username := tk.Username
	// Decrease user balance
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	dueDate, _ := time.Parse("2006-01-02", c.PostForm("due_date"))
	err = config.DB.Model(&users).Where("username  = ?", username).First(&users).Error
	if err != nil {
		util.CallErrorNotFound(c, "can't select balance", err)
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
		util.CallServerError(c, "fail to create paylist", err)
		return
	}
	util.CallSuccessOK(c, "paylist item created successfully!", paylist.ID)

}

//UpdateUserPaylist handle status completed in paylists
func UpdateUserPaylist(c *gin.Context) {
	paylist := model.Paylist{}
	user := model.User{}
	tk := User{}

	// Parse the token payload
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprint(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallUserError(c, "fail to parse the token, make sure the token is valid", err)
	}

	// Check User balance
	username := tk.Username
	id, _ := strconv.Atoi(c.Param("id"))
	if err = config.DB.Model(&user).Select("balance").Where("username = ?", username).Find(&user).Error; err != nil {
		util.CallErrorNotFound(c, "user not found", err)
		return
	}
	if err = config.DB.Model(&paylist).Where("ID = ? AND username = ?", id, username).First(&paylist).Error; err != nil || err == gorm.ErrRecordNotFound {
		util.CallErrorNotFound(c, "no paylist found", nil)
		return
	}
	firstAmount := paylist.Amount
	amount, _ := strconv.Atoi(c.PostForm("amount"))
	updatedPaylist := model.Paylist{
		Name:    c.PostForm("name"),
		Amount:  amount,
		DueDate: c.PostForm("due_date"),
	}
	if tk.Username != paylist.Username {
		util.CallServerError(c, "not authorized", nil)
		return
	}
	// Update paylist
	if err = config.DB.Model(&paylist).Where("username = ?", username).Update(&updatedPaylist).Error; err != nil {
		fmt.Println(err)
		util.CallServerError(c, "something error when try to update paylist", err)
		return
	}

	// Update user balance
	err = config.DB.Model(&user).Where("username = ?", username).Update("balance", (firstAmount-amount)+user.Balance).Error
	if err != nil {
		fmt.Println(err)
		util.CallServerError(c, "something error when try to update user balance", err)
		return
	}
	util.CallSuccessOK(c, "paylist successfully updated!", paylist)
}

// UpdateUserPaylistStatus is a function to mark the paylist completed or not
func UpdateUserPaylistStatus(c *gin.Context) {
	paylist := model.Paylist{}
	user := model.User{}
	tk := User{}

	// Parse the token payload and validate the username is own the paylist
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprint(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallUserError(c, "fail to parse the token, make sure the token is valid", err)
	}
	username := tk.Username

	// Check User balance
	id, _ := strconv.Atoi(c.Param("id"))
	if err := config.DB.Model(&paylist).Where("ID = ? AND username = ?", id, username).First(&paylist).Error; err != nil || err == gorm.ErrRecordNotFound {
		util.CallErrorNotFound(c, "no paylist found", nil)
		return
	}
	if err := config.DB.Model(&user).Select("balance").Where("username = ?", username).First(&user).Error; err != nil {
		util.CallErrorNotFound(c, "user not found", err)
		return
	}

	// Update the paylist status
	if user.Balance >= 0 && paylist.Completed == false {
		paylist.Completed = true
		config.DB.Model(&paylist).Where("ID = ? and username = ?", id, username).Update(&paylist)
	} else if user.Balance < 0 && paylist.Completed == false {
		paylist.Completed = false
		config.DB.Model(&paylist).Where("ID = ? and username = ?", id, username).Update(&paylist)
	}
	util.CallSuccessOK(c, "successfully update user paylist", paylist.Completed)
}

//DeleteUserPaylist handle deleted user paylist
func DeleteUserPaylist(c *gin.Context) {
	paylistID, _ := strconv.Atoi(c.Param("id"))
	paylist := &model.Paylist{}
	user := &model.User{}
	tk := User{}

	if paylistID == 0 {
		util.CallUserError(c, "please specify paylist id", nil)
		return
	}
	// Parse the token payload and validate the username is own the paylist
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.ParseWithClaims(tokenString, &tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(fmt.Sprint(config.Conf.JWTSignature)), nil
	})
	if err != nil || token == nil {
		fmt.Println(err, token)
		util.CallUserError(c, "fail to parse the token, make sure the token is valid", err)
	}
	username := tk.Username
	if err = config.DB.Where("ID = ?", paylistID).Find(&paylist).Error; err != nil {
		util.CallErrorNotFound(c, "no paylist found!", err)
		c.Abort()
		return
	}

	config.DB.Model(&paylist).Where("username = ?", username).First(&paylist)
	if tk.Username != paylist.Username {
		util.CallServerError(c, "user not authorized", nil)
		c.Abort()
		return
	}

	if err = config.DB.Model(&paylist).Select("amount, completed").Where("ID = ? and username = ?", paylistID, username).First(&paylist).Error; err != nil {
		util.CallErrorNotFound(c, "can't select amount", err)
		return
	}

	if err := config.DB.Model(&user).Select("balance").Where("username = ?", username).First(&user).Error; err != nil {
		util.CallErrorNotFound(c, "can't select balance", err)
		return
	}
	if paylist.Completed == false {
		b := paylist.Amount + user.Balance
		if err = config.DB.Table("users").Where("username = ?", username).Update("balance", b).Error; err != nil {
			fmt.Println(&user.Balance)
			util.CallServerError(c, "fail to update the user balance", err)
			return
		}
	}
	if err := config.DB.Model(&paylist).Where("ID = ? and username = ?", paylistID, username).Delete(&paylist).Error; err != nil {
		util.CallServerError(c, "paylist fail to delete", err)
	}
	util.CallSuccessOK(c, "paylist successfully deleted!", nil)
}

// FilterPaylist is a function to filter paylist based on month
func FilterPaylist(c *gin.Context) {
	paylist := model.Paylist{}
	//param,_ := strconv.ParseInt(c.Param("month"),10,64)
	//month := paylist.CreatedAt.Month()
	param := c.Param("created_at")
	if param != "" {
		util.CallServerError(c, "please specify the filter parameter", nil)
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
		util.CallUserError(c, "fail to parse the token, make sure the token is valid", err)
	}
	username := tk.Username

	err = config.DB.Model(&paylist).Where("username = ? and MONTH(created_at) = ?", username, param).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	util.CallSuccessOK(c, "success paylist", paylist)
}
