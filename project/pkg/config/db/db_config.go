package db

import (
	"sample_app/app/dto"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbCon *gorm.DB

func GetDBConnection() *gorm.DB {
	return dbCon
}

func SetDBConnection(db *gorm.DB) {
	dbCon = db
}

// create a new connection to the database
func InitDBConnection() error {
	dsn := "host=localhost user=postgres password=postgres dbname=sample_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	// migrate the models to the database
	db.AutoMigrate(&dto.Category{}, &dto.Product{}, &dto.User{}, &dto.Interaction{})

	// create default users
	defaultUser := []dto.User{
		{
			Name:     "admin",
			Email:    "admin",
			Password: "$2a$10$blcTyN6aOhNn/VNRDYWage/prlcjBnpCnKg4HOcnA65AlTYZ.JoX2",
			Role:     dto.RoleAdmin,
		},
		{
			Name:     "system",
			Email:    "system",
			Password: "$2a$10$LA7LOUHsetl/t7dSHreCxOntnGCrEkEzMR.MCP.bapLuCmASWN3Ji",
			Role:     dto.RoleSystem,
		},
	}

	db.Create(&defaultUser)

	dbCon = db
	return nil
}
