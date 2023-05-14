package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var mockDBCon sqlmock.Sqlmock

func GetMockDB() sqlmock.Sqlmock {
	return mockDBCon
}

func MockDatabase() sqlmock.Sqlmock {
	// Create a new mock database connection
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})

	if err != nil {
		return nil
	}

	SetDBConnection(gormDB)
	mockDBCon = mock
	return mockDBCon
}
