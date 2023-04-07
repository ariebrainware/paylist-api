package model

import (
	"github.com/jinzhu/gorm"
)

// User is a model for user table
type User struct {
	gorm.Model
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `gorm:"primary_key;not null"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
}

func (User) TableName() string {
	return "users"
}