package db

import (
	"log"
	"os"

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
	// todo migration
	//if err := db.AutoMigrate(&calculationService.Calculation{}); err != nil {
	//	log.Fatalf("couldnt migrate: %v", err)
	//}
	return db, nil
}
