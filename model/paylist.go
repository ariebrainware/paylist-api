package model

import (
	"github.com/jinzhu/gorm"
)

// Paylist is a model for paylist table
type Paylist struct {
	gorm.Model
	// ID     int    `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}
