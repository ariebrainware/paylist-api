package model

import "time"

// User is a model for user table
type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Username  string     `gorm:"primary_key;not null"`
	Password  string     `json:"password"`
	Balance   int        `json:"balance"`
}
