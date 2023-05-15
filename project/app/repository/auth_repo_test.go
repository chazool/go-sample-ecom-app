package repository_test

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"sample_app/app/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepository_FindById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
			AddRow(testID, user.Name, user.Email, user.Password, user.Role)
		// Define expected query
		mockDBCon.ExpectQuery(Select).WithArgs(1).WillReturnRows(rows)

		// Call the FindById method
		result, err := userRepo.FindById(int(testID))
		checkUserSuccess(t, user, result, err)

	})

	t.Run("Error: Record Not Found", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(Select).WithArgs(1).WillReturnError(gorm.ErrRecordNotFound)
		// Call method being tested
		_, err := userRepo.FindById(int(testID))
		checkError(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("Error: Other DB Error", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(Select).WithArgs(1).WillReturnError(returnErr)

		// Call method being tested
		_, err := userRepo.FindById((int(testID)))
		checkError(t, err, returnErr)
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		mockDBCon.ExpectQuery(Select).WithArgs(user.Email).WillReturnRows(sqlmock.
			NewRows([]string{"id", "name", "email", "password", "role"}).
			AddRow(testID, user.Name, user.Email, user.Password, user.Role))

		// Call the FindById method
		result, err := userRepo.FindByEmail(user.Email)
		checkUserSuccess(t, user, result, err)
	})

	t.Run("Error: Record Not Found", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(Select).WithArgs(user.Email).WillReturnError(gorm.ErrRecordNotFound)

		// Call method being tested
		_, err := userRepo.FindByEmail(user.Email)
		checkError(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("Error: Other DB Error", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(Select).WithArgs(user.Email).WillReturnError(returnErr)

		// Call method being tested
		_, err := userRepo.FindByEmail(user.Email)
		checkError(t, err, returnErr)
	})
}

func TestUserRepository_Create(t *testing.T) {
	query := regexp.QuoteMeta(fmt.Sprintf(InsertInto, "users"))
	args := []driver.Value{sqlmock.AnyArg(), sqlmock.AnyArg(), nil, user.Name, user.Email, user.Password, user.Role}

	t.Run("Success", func(t *testing.T) {
		returnRows := []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)}
		mockDBCon.ExpectQuery(query).WithArgs(args...).WillReturnRows(returnRows...)

		// Call the Create function with the expected user
		createdUser, err := userRepo.Create(user)
		checkUserSuccess(t, user, createdUser, err)
	})

	t.Run("Error: Other DB Error", func(t *testing.T) {
		returnErr := errors.New("some DB error")
		// Define expected query
		mockDBCon.ExpectQuery(query).WithArgs(args...).WillReturnError(returnErr)

		// Call the Create function with the expected err
		_, err := userRepo.Create(user)
		checkError(t, err, returnErr)
	})
}

func checkUserSuccess(t *testing.T, expected dto.User, result dto.User, err error) {
	assert.NoError(t, err)
	assert.Equal(t, expected.Name, result.Name)
	assert.Equal(t, expected.Email, result.Email)
	assert.Equal(t, expected.Password, result.Password)
	assert.Equal(t, expected.Role, result.Role)
	// Verify that all expectations were met
	assert.NoError(t, mockDBCon.ExpectationsWereMet())
}
