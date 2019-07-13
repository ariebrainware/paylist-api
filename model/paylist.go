package model

import (
	"github.com/jinzhu/gorm"
)

type Paylist struct {
	gorm.Model
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}
