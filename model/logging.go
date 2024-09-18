package model

import (
	"time"
)

// Logging represents the logging information for a user.
// It includes the username, token, user status, creation time, and deletion time.
//
// Fields:
// - Username: The username of the user.
// - Token: The authentication token of the user.
// - UserStatus: The status of the user (active/inactive).
// - CreatedAt: The timestamp when the logging record was created.
// - DeletedAt: The timestamp when the logging record was deleted (if applicable).
type Logging struct {
	Username   string `sql:"column:username" json:"username"`
	Token      string `sql:"column:token" json:"token"`
	UserStatus bool   `sql:"column:user_status"`
	CreatedAt  time.Time
	DeletedAt  *time.Time
}

func (Logging) TableName() string {
	return "loggings"
}
