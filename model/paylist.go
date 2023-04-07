package model

import (
	"github.com/jinzhu/gorm"
)

// Paylist is a model for paylist table
type Paylist struct {
	gorm.Model
	Name      string `json:"name"`
	Amount    int    `json:"amount"`
	Username  string `json:"username"`
	DueDate   string `gorm:"default:null"`
	Completed bool   `json:"completed"`
}

func (Paylist) TableName() string {
	return "paylists"
}