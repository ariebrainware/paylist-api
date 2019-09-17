package model

import (
	"time"
)

//Logging model for table loggings
type Logging struct {
	Username   string `sql:"column:username" json:"username"`
	Token      string `sql:"column:token" json:"token"`
	UserStatus bool   `sql:"column:user_status"`
	CreatedAt  time.Time
	DeletedAt  *time.Time
}
