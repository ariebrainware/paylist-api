package model

import (
	"github.com/jinzhu/gorm"
)

// User represents a user in the system with their associated details.
// It includes fields for email, name, username, password, and balance.
// The Username field is marked as the primary key and cannot be null.
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
