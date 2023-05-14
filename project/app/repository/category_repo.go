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
	FindById(id uint) (dto.Category, error)
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
	log.Printf("Create category started")

	// Create the category record
	if err := r.db.Create(&category).Error; err != nil {
		log.Printf("Error while creating category, err: %v", err)
		return category, err
	}

	log.Printf("Created category: %+v", category)
	log.Printf("Create category completed")
	return category, nil
}

func (r *categoryRepository) FindById(id uint) (dto.Category, error) {
	log.Printf("FindById started for category ID: %d", id)

	var category dto.Category
	result := r.db.First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("Category not found for ID: %d", id)
			return category, ErrRecordNotFound
		}
		log.Printf("Error while finding category for ID: %d, err: %v", id, result.Error)
		return category, result.Error
	}

	log.Printf("Found category for ID %d: %+v", id, category)
	return category, nil
}

func (r *categoryRepository) FindAll() ([]dto.Category, error) {
	log.Println("FindAll started")

	var categories []dto.Category
	if err := r.db.Find(&categories).Error; err != nil {
		log.Printf("Error while finding categories: %v", err)
		return nil, err
	}

	log.Printf("Found %d categories: %+v", len(categories), categories)
	return categories, nil
}
