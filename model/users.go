package model

import "github.com/jinzhu/gorm"

// User is a model for user table
type User struct {
	gorm.Model
	// ID       int    `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}
