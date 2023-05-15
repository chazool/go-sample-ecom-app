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

var productColumns = []string{"id", "name", "description", "price", "category_id", "interactions"}
var productRow = []driver.Value{testID, product.Name, product.Description, product.Price, product.CategoryID, product.Interactions}

func TestProductRepository_Create(t *testing.T) {
	query := regexp.QuoteMeta(fmt.Sprintf(InsertInto, "products"))
	args := []driver.Value{sqlmock.AnyArg(), sqlmock.AnyArg(), nil, product.Name, product.Description, product.Price, product.CategoryID, product.Interactions}

	t.Run("Success", func(t *testing.T) {
		returnRows := []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(1)}
		mockDBCon.ExpectQuery(query).WithArgs(args...).WillReturnRows(returnRows...)

		// Call the Create function with the expected user
		createdProduct, err := productRepo.Create(product)
		checkProductSuccess(t, product, createdProduct, err)
	})

	t.Run("Error: Other DB Error", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(query).WithArgs(args...).WillReturnError(returnErr)

		// Call the Create function with the expected err
		_, err := productRepo.Create(product)
		checkError(t, err, returnErr)
	})
}

func TestCProductRepository_FindById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		//var productColumns = []string{"id", "name", "description", "price", "category_id", "interactions", "weighted_score"}
		//var productRow = []driver.Value{testID, product.Name, product.Description, product.Price, product.CategoryID, product.Interactions, product.WeightedScore}
		rows := sqlmock.NewRows(productColumns).AddRow(productRow...)
		// Define expected query
		mockDBCon.ExpectQuery(Select).WithArgs(1).WillReturnRows(rows)

		// Call the FindById method
		result, err := productRepo.FindById(testID)
		checkProductSuccess(t, product, result, err)

	})

	t.Run("Error: Record Not Found", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(Select).WithArgs(1).WillReturnError(gorm.ErrRecordNotFound)
		// Call method being tested
		_, err := productRepo.FindById(testID)
		checkError(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("Error: Other DB Error", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(Select).WithArgs(1).WillReturnError(returnErr)

		// Call method being tested
		_, err := productRepo.FindById(testID)
		checkError(t, err, returnErr)
	})
}

func TestCProductRepository_FindAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDBCon.ExpectQuery(Select).WillReturnRows(sqlmock.NewRows(productColumns).AddRow(productRow...))

		// Call the FindById method
		result, err := productRepo.FindAll()

		assert.NoError(t, err)
		assert.Equal(t, len([]dto.Product{product}), len(result))

	})

	t.Run("Error: Record Not Found", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(Select).WillReturnError(gorm.ErrRecordNotFound)
		// Call method being tested
		_, err := productRepo.FindAll()
		checkError(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("Error: Other DB Error", func(t *testing.T) {
		// Define expected query
		mockDBCon.ExpectQuery(Select).WillReturnError(returnErr)
		// Call method being tested
		_, err := productRepo.FindAll()
		checkError(t, err, returnErr)
	})
}

func checkProductSuccess(t *testing.T, expected dto.Product, result dto.Product, err error) {
	assert.NoError(t, err)
	assert.Equal(t, expected.Name, result.Name)
	assert.Equal(t, expected.CategoryID, result.CategoryID)
	assert.Equal(t, expected.Price, result.Price)
	assert.Equal(t, expected.Description, result.Description)
	assert.Equal(t, expected.Interactions, result.Interactions)
	assert.Equal(t, expected.WeightedScore, result.WeightedScore)
	// Verify that all expectations were met
	assert.NoError(t, mockDBCon.ExpectationsWereMet())
}
