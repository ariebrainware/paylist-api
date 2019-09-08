package model

import (
	"github.com/jinzhu/gorm"
	//"github.com/ariebrainware/paylist-api/model"
)

// Paylist is a model for paylist table
type Paylist struct {
	gorm.Model
	// ID     int    `json:"id"`
	Name      string `json:"name"`
	Amount    int    `json:"amount"`
	Username  string
	Completed bool   `json:"completed"`
	User      []User `gorm:"foreignkey:username"`
}
