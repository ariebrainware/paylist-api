package model

import (
	"time"
)

//Logging model for table loggings
type Income struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `sql:"column:username" json:"username"`
	Income int    `sql:"column:income" json:"income"`
	CreatedAt time.Time
	UpdatedAt  time.Time
	DeletedAt *time.Time `sql:"index"`
}
