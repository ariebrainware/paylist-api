package model

import (
	"time"
	//"github.com/jinzhu/gorm"
)

type Logging struct {
	Username string `json:"username"`
	Token string `json:"token"`
	UserStatus bool
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}