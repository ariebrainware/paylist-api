package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	// ID       int    `json:"user_id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}
