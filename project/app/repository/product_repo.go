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
	log.Printf("Deleting product: %+v", product)
	result := r.db.Delete(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *productRepository) FindById(id uint) (dto.Product, error) {
	var product dto.Product
	result := r.db.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return product, ErrRecordNotFound
		}
		return product, result.Error
	}
	log.Printf("Found product by ID %d: %+v", id, product)
	return product, nil
}

func (r *productRepository) FindAll() ([]dto.Product, error) {
	var products []dto.Product

	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	log.Printf("Found %d products: %+v", len(products), products)
	return products, nil
}

func (r *productRepository) Create(product dto.Product) (dto.Product, error) {
	// Create the product record
	if err := r.db.Create(&product).Error; err != nil {
		return product, err
	}
	log.Printf("Created product: %+v", product)
	return product, nil
}

func (r *productRepository) UpdateInteractions(productID, interactions uint) error {
	result := r.db.Model(&dto.Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
		"interactions": interactions,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
