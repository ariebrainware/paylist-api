package model

import (
	"github.com/jinzhu/gorm"
	//"github.com/ariebrainware/paylist-api/model"
)

type Paylist struct {
	gorm.Model
	// ID     int    `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
	Username string
	Completed int `json:"completed"`
}
