package repository

import (
	"errors"
	"sample_app/app/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepository_FindById(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Create GORM DB with mock DB
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}
	//defer gormDB.Close()

	// Create mocked user
	mockedUser := dto.User{
		//	ID:   1,

		Name: "John",
	}

	// Create UserRepository with mocked GORM DB
	userRepository := &userRepository{
		db: gormDB,
	}

	// Test case 1: user exists
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(mockedUser.ID, mockedUser.Name))
	user, err := userRepository.FindById(1)
	assert.NoError(t, err)
	assert.Equal(t, mockedUser, user)

	// Test case 2: user does not exist
	mock.ExpectQuery("SELECT").WithArgs(2).WillReturnError(gorm.ErrRecordNotFound)
	_, err = userRepository.FindById(2)
	assert.True(t, errors.Is(err, ErrRecordNotFound))

	// Check all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
