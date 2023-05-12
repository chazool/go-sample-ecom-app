package services

import (
	"errors"
	"log"
	"sample_app/app/dto"
	"sample_app/app/repository"
)

// error const
var (
	ErrCategoryNotFound         = errors.New("category not found")
	ErrFailToDeleteCategory     = errors.New("failed to delete category")
	ErrFailedToRetrieveCategory = errors.New("failed to retrieve categories")
	ErrFailToCreateCategory     = errors.New("failed to create category record")
)

type CategoryService interface {
	Create(category dto.Category) (dto.Category, error)
	FindById(id int) (dto.Category, error)
	FindAll() ([]dto.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService() CategoryService {
	return &categoryService{
		categoryRepo: repository.NewCategoryRepository(),
	}
}

func (service *categoryService) FindById(id int) (dto.Category, error) {

	log.Printf("Retrieving category with ID %d\n", id)

	// Retrieve category from the database with given ID
	category, err := service.categoryRepo.FindById(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return category, ErrCategoryNotFound
		}
		return category, err
	}

	log.Printf("Retrieved category with ID %d: %+v\n", id, category)

	return category, nil
}

func (service *categoryService) FindAll() ([]dto.Category, error) {

	log.Println("Retrieving all categories")

	// Retrieve all categories from the database
	categories, err := service.categoryRepo.FindAll()
	if err != nil {
		return nil, ErrFailedToRetrieveCategory
	}

	log.Printf("Retrieved %d categories\n", len(categories))

	return categories, nil
}

func (service *categoryService) Create(category dto.Category) (dto.Category, error) {

	log.Printf("Creating category: %+v\n", category)

	// Retrieve category from the database with given ID
	category, err := service.categoryRepo.Create(category)
	if err != nil {
		return category, ErrFailToCreateCategory
	}

	log.Printf("Created category with ID %d\n", category.ID)

	return category, nil
}
