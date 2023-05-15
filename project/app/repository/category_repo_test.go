package repository_test

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"sample_app/app/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCategoryRepository_Create(t *testing.T) {
	query := regexp.QuoteMeta(fmt.Sprintf(InsertInto, "categories"))
	args := []driver.Value{sqlmock.AnyArg(), sqlmock.AnyArg(), nil, category.Name}

	t.Run("Success", func(t *testing.T) {
		returnRows := []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)}
		mockDBCon.ExpectQuery(query).WithArgs(args...).WillReturnRows(returnRows...)

		// Call the Create function with the expected user
		createdCategory, err := categoryRepo.Create(category)
		checkCategorySuccess(t, category, createdCategory, err)
	})

	t.Run("Error: Other DB Error", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(query).WithArgs(args...).WillReturnError(returnErr)

		// Call the Create function with the expected err
		_, err := categoryRepo.Create(category)
		checkError(t, err, returnErr)
	})
}

func TestCategoryRepository_FindById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(testID, category.Name)
		// Define expected query
		mockDBCon.ExpectQuery(Select).WithArgs(1).WillReturnRows(rows)

		// Call the FindById method
		result, err := categoryRepo.FindById(testID)
		checkCategorySuccess(t, category, result, err)

	})

	t.Run("Error: Record Not Found", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(Select).WithArgs(1).WillReturnError(gorm.ErrRecordNotFound)
		// Call method being tested
		_, err := categoryRepo.FindById(testID)
		checkError(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("Error: Other DB Error", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(Select).WithArgs(1).WillReturnError(returnErr)

		// Call method being tested
		_, err := categoryRepo.FindById(testID)
		checkError(t, err, returnErr)
	})
}

func checkCategorySuccess(t *testing.T, expected dto.Category, result dto.Category, err error) {
	assert.NoError(t, err)
	assert.Equal(t, expected.Name, result.Name)
	// Verify that all expectations were met
	assert.NoError(t, mockDBCon.ExpectationsWereMet())
}
