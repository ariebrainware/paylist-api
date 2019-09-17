package model

import "time"

// Paylist is a model for paylist table
type Paylist struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
	Name      string     `json:"name"`
	Amount    int        `json:"amount"`
	Username  string     `json:"username"`
	Completed bool       `json:"completed"`
}
