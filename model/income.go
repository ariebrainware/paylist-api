package model

import (
	"time"
)

// Income represents the income data model.
// It includes fields for ID, Username, Income, CreatedAt, UpdatedAt, and DeletedAt.
// The ID field is the primary key.
// The Username field stores the username associated with the income.
// The Income field stores the income amount.
// The CreatedAt field stores the timestamp when the record was created.
// The UpdatedAt field stores the timestamp when the record was last updated.
// The DeletedAt field stores the timestamp when the record was deleted, if applicable.
type Income struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `sql:"column:username" json:"username"`
	Income    int    `sql:"column:income" json:"income"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (Income) TableName() string {
	return "incomes"
}
