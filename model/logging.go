package model

import (
	"time"
	//"github.com/jinzhu/gorm"
)

type Logging struct {
	Username string `json:"username"`
	Token string `json:"token"`
	User_status bool
	Created_At time.Time
	Deleted_At *time.Time `sql:"index"`
}