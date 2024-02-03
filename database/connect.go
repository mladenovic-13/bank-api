package database

import (
	"github.com/mladenovic-13/bank-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(connStr string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Transaction{},
		&models.Request{},
	)

	if err != nil {
		return nil, err
	}

	return db, err
}
