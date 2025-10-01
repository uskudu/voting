package db

import (
	"log"
	"os"
	"voting/internal/poll"
	"voting/internal/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() (*gorm.DB, error) {
	dsn := os.Getenv("DB")
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("couldnt connect to DB: %v", err)
	}

	if err := DB.AutoMigrate(&user.User{}, &poll.Poll{}, &poll.Option{}); err != nil {
		log.Fatalf("couldnt migrate: %v", err)
	}
	return DB, nil
}
