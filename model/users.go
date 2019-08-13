package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}
