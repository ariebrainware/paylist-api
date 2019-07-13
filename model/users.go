package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	// ID       int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
