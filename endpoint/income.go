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
		util.CallServerError(c, "fail to parse the token, make sure token is valid", err)
		return
	}
	username := tk.Username
	config.DB.Model(&income).Where("username = ? ", username).Find(&income)

	if len(income) <= 0 {
		util.CallErrorNotFound(c, "No Income Found!", nil)
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
	util.CallSuccessOK(c, "Fetch All Income Data ", inc)
}

//Update Func to handle edit income when user wrong input the income
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
		util.CallUserError(c, "fail to parse the token, make sure the token is valid", err)
	}

	// Check User balance
	username := tk.Username
	id, _ := strconv.Atoi(c.Param("id"))
	if err = config.DB.Model(&user).Select("balance").Where("username = ?", username).Find(&user).Error; err != nil {
		util.CallErrorNotFound(c, "user not found", err)
		return
	}
	if err = config.DB.Model(&income).Where("ID = ? AND username = ?", id, username).First(&income).Error; err != nil || err == gorm.ErrRecordNotFound {
		util.CallErrorNotFound(c, "no income found", nil)
		return
	}
	fmt.Println("inc", income.Income)
	firstIncome := income.Income
	inc, _ := strconv.Atoi(c.PostForm("income"))
	updatedIncome := model.Income{
		Income: inc,
	}
	if tk.Username != income.Username {
		util.CallServerError(c, "not authorized", nil)
		return
	}
	// Update paylist
	if err = config.DB.Model(&income).Where("username = ?", username).Update(&updatedIncome).Error; err != nil {
		fmt.Println(err)
		util.CallServerError(c, "something error when try to update paylist", err)
		return
	}

	// Update user balance
	err = config.DB.Model(&user).Where("username = ?", username).Update("balance", (inc+user.Balance)-firstIncome).Error
	if err != nil {
		fmt.Println(err)
		util.CallServerError(c, "something error when try to update user balance", err)
		return
	}
	util.CallSuccessOK(c, "income successfully updated!", income)
}

//func DeleteIncome handle delete income
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
		util.CallServerError(c, "fail to parse the token, make sure token is valid", err)
		return
	}
	username := tk.Username
	if err = config.DB.Where("ID = ?", ID).Find(&inc).Error; err != nil {
		util.CallErrorNotFound(c, "no income data found!", err)
		c.Abort()
		return
	}
	config.DB.Model(&inc).Select("income").Where("username = ? and ID = ? ", username, ID).Find(&inc)
	if inc.ID == 0 {
		util.CallErrorNotFound(c, "can't select income", nil)
		return
	}
	config.DB.Model(&users).Select("balance").Where("username = ?", username).Find(&users)
	if inc.ID == 0 {
		util.CallErrorNotFound(c, "can't select income", nil)
		return
	}
	income := inc.Income
	fmt.Println("inc", inc.Income)
	err = config.DB.Model(&users).Where("username = ?", username).Update("balance", users.Balance-income).Error
	if err != nil {
		fmt.Println(err)
		util.CallServerError(c, "failed update income", nil)
		return
	}
	config.DB.Model(&inc).Where("username = ?", username).Delete(&inc)
	if tk.Username != inc.Username {
		util.CallServerError(c, "not authorized", nil)
		return
	}
	util.CallSuccessOK(c, "user delete successfully!", nil)
}
