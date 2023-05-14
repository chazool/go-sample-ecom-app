package repository

import (
	"errors"
	"log"
	"sample_app/app/dto"
	"sample_app/pkg/config/db"

	"gorm.io/gorm"
)

var ErrRecordNotFound = errors.New("record not found")
var ErrFailedToFindProduct = errors.New("failed to find product")

type ProductRepository interface {
	Create(product dto.Product) (dto.Product, error)
	Delete(product dto.Product) error
	FindById(id uint) (dto.Product, error)
	FindAll() ([]dto.Product, error)
	UpdateInteractions(productID, interactions uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository() ProductRepository {
	return &productRepository{
		db: db.GetDBConnection(),
	}
}

func (r *productRepository) Delete(product dto.Product) error {
	log.Printf("Deleting product started for product: %+v", product)
	defer log.Printf("Deleting product ended for product: %+v", product)

	result := r.db.Delete(&product)
	if result.Error != nil {
		log.Printf("Error while deleting product: %+v, err: %v", product, result.Error)
		return result.Error
	}

	log.Printf("Deleted product: %+v", product)
	return nil
}

func (r *productRepository) FindById(id uint) (dto.Product, error) {
	log.Printf("Finding product by ID: %d", id)
	var product dto.Product
	result := r.db.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("Product not found with ID: %d", id)
			return product, ErrRecordNotFound
		}
		log.Printf("Error while finding product with ID: %d, err: %v", id, result.Error)
		return product, result.Error
	}
	log.Printf("Found product by ID %d: %+v", id, product)
	return product, nil
}

func (r *productRepository) FindAll() ([]dto.Product, error) {
	log.Printf("FindAll started")
	defer log.Printf("FindAll completed")

	var products []dto.Product
	if err := r.db.Find(&products).Error; err != nil {
		log.Printf("Error while finding all products: %v", err)
		return nil, err
	}

	log.Printf("Found %d products", len(products))
	return products, nil
}

func (r *productRepository) Create(product dto.Product) (dto.Product, error) {
	log.Printf("Create product started for product: %+v", product)

	// Create the product record
	if err := r.db.Create(&product).Error; err != nil {
		log.Printf("Error while creating product: %+v, err: %v", product, err)
		return product, err
	}

	log.Printf("Created product: %+v", product)
	return product, nil
}

func (r *productRepository) UpdateInteractions(productID, interactions uint) error {
	log.Printf("UpdateInteractions started for productID: %d, interactions: %d", productID, interactions)

	result := r.db.Model(&dto.Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
		"interactions": interactions,
	})
	if result.Error != nil {
		log.Printf("Error while updating interactions for productID: %d, err: %v", productID, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Printf("No rows affected while updating interactions for productID: %d", productID)
		return gorm.ErrRecordNotFound
	}

	log.Printf("Updated interactions for productID: %d, new interactions: %d", productID, interactions)
	return nil
}
