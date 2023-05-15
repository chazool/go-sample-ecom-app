package repository_test

import (
	"errors"
	"sample_app/app/dto"
	"sample_app/app/repository"
	"sample_app/pkg/config/db"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	Select     = "SELECT"
	InsertInto = `INSERT INTO "%s"`
)

var returnErr = errors.New("some DB error")

var (
	mockDBCon    sqlmock.Sqlmock
	userRepo     repository.UserRepository
	categoryRepo repository.CategoryRepository
	productRepo  repository.ProductRepository
)

var (
	testID uint = 1
	user        = dto.User{
		Name:     "John",
		Email:    "john@example.com",
		Password: "secret",
		Role:     "admin",
	}
	category = dto.Category{
		Name: "Test Category",
	}
	product = dto.Product{
		CategoryID:   testID,
		Name:         "Test Product",
		Price:        100,
		Description:  "Test Description",
		Interactions: 0,
	}
)

func init() {
	mockDBCon = db.MockDatabase()
	userRepo = repository.NewUserRepository()
	categoryRepo = repository.NewCategoryRepository()
	productRepo = repository.NewProductRepository()
}

func checkError(t *testing.T, resultErr error, expError error) {
	assert.EqualError(t, resultErr, expError.Error())
}
