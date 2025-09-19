package db

import (
	"log"
	"os"
	"voting/internal/poll"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() (*gorm.DB, error) {
	dsn := os.Getenv("DB")
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("couldnt connect to db: %v", err)
	}

	if err := db.AutoMigrate(&poll.Poll{}, &poll.Option{}, &poll.Vote{}); err != nil {
		log.Fatalf("couldnt migrate: %v", err)
	}
	return db, nil
}
