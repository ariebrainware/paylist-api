package model

import (
	"time"
)

//Logging model for table loggings
type Incomes struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `sql:"column:username" json:"username"`
	Balance   int    `sql:"column:balance" json:"balance"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (Incomes) TableName() string {
	return "incomes"
}
