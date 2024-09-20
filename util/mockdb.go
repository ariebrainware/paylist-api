package util

// Add necessary imports

import (
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// MockDB creates a mock database connection for testing
func MockDB() *gorm.DB {

	db, err := gorm.Open("sqlite3", ":memory:")

	if err != nil {

		panic("failed to connect database")

	}

	db.LogMode(false)

	return db

}
