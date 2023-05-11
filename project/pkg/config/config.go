package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbCon *gorm.DB

func GetDBConnection() *gorm.DB {
	return dbCon
}

// create a new connection to the database
func InitDBConnection(dst ...interface{}) error {
	dsn := "host=localhost user=postgres password=postgres dbname=sample_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	// migrate the models to the database
	db.AutoMigrate(dst...)

	dbCon = db
	return nil
}
