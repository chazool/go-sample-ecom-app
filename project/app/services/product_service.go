package services

import (
	"errors"
	"log"
	"sample_app/app/dto"
	"sample_app/app/repository"
)

// error const
var (
	ErrProductNotFound         = errors.New("product not found")
	ErrFailToDeleteProduct     = errors.New("failed to delete product")
	ErrFailedToRetrieveProduct = errors.New("failed to retrieve products")
	ErrFailToCreateProduct     = errors.New("failed to create product record")
)

type ProductService interface {
	Create(product dto.Product) (dto.Product, error)
	Delete(id int) error
	FindById(id int) (dto.Product, error)
	FindAll() ([]dto.Product, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService() ProductService {
	return &productService{
		productRepo: repository.NewProductRepository(),
	}
}

func (service *productService) Delete(id int) error {

	log.Printf("Deleting product with ID %d\n", id)

	// Find product with given ID
	product, err := service.FindById(id)
	if err != nil {
		return err
	}

	// Delete the product
	err = service.productRepo.Delete(product)
	if err != nil {
		return ErrFailToDeleteProduct
	}

	log.Printf("Deleted product with ID %d\n", id)

	return nil
}

func (service *productService) FindById(id int) (dto.Product, error) {

	log.Printf("Retrieving product with ID %d\n", id)

	// Retrieve product from the database with given ID
	product, err := service.productRepo.FindById(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return product, ErrProductNotFound
		}
		return product, err
	}

	log.Printf("Retrieved product with ID %d: %+v\n", id, product)

	return product, nil
}

func (service *productService) FindAll() ([]dto.Product, error) {

	log.Println("Retrieving all products")

	// Retrieve all products from the database
	products, err := service.productRepo.FindAll()
	if err != nil {
		return nil, ErrFailedToRetrieveProduct
	}

	log.Printf("Retrieved %d products\n", len(products))

	return products, nil
}

func (service *productService) Create(product dto.Product) (dto.Product, error) {

	log.Printf("Creating product: %+v\n", product)

	// Retrieve product from the database with given ID
	product, err := service.productRepo.Create(product)
	if err != nil {
		return product, ErrFailToCreateProduct
	}

	log.Printf("Created product with ID %d\n", product.ID)

	return product, nil
}
