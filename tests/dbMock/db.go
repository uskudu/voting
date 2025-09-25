package dbMock

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SqliteMock() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
