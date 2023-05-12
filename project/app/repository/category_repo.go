package repository

import (
	"errors"
	"log"
	"sample_app/app/dto"
	"sample_app/pkg/config/db"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category dto.Category) (dto.Category, error)
	FindById(id int) (dto.Category, error)
	FindAll() ([]dto.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{
		db: db.GetDBConnection(),
	}
}

func (r *categoryRepository) Create(category dto.Category) (dto.Category, error) {
	// Create the category record
	if err := r.db.Create(&category).Error; err != nil {
		return category, err
	}
	log.Printf("Created category: %+v", category)
	return category, nil
}

func (r *categoryRepository) FindById(id int) (dto.Category, error) {
	var category dto.Category
	result := r.db.First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return category, ErrRecordNotFound
		}
		return category, result.Error
	}
	log.Printf("Found category by ID %d: %+v", id, category)
	return category, nil
}

func (r *categoryRepository) FindAll() ([]dto.Category, error) {
	var categories []dto.Category

	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	log.Printf("Found %d categories: %+v", len(categories), categories)
	return categories, nil
}
