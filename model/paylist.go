package model

import (
	"github.com/jinzhu/gorm"
)

// Paylist represents a payment list item with details such as name, amount,
// username, due date, and completion status. It embeds gorm.Model to include
// fields like ID, CreatedAt, UpdatedAt, and DeletedAt.
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
