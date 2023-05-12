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
	ErrFailToInteractsProduct  = errors.New("failed to Interact product record")
)

type ProductService interface {
	Create(product dto.Product) (dto.Product, error)
	Delete(id int) error
	findById(id int) (dto.Product, error)
	FindByIdWithInteract(id, userID int) (dto.Product, error)
	FindAll() ([]dto.Product, error)
}

type productService struct {
	productRepo        repository.ProductRepository
	interactionService InteractionService
}

func NewProductService() ProductService {
	return &productService{
		productRepo:        repository.NewProductRepository(),
		interactionService: NewInteractionService(),
	}
}

func (service *productService) Delete(id int) error {

	log.Printf("Deleting product with ID %d\n", id)

	// Find product with given ID
	product, err := service.findById(id)
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

func (service *productService) findById(id int) (dto.Product, error) {

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

func (service *productService) FindByIdWithInteract(id, userID int) (dto.Product, error) {

	log.Printf("Retrieving product with ID %d\n", id)

	// Retrieve product from the database with given ID
	product, err := service.findById(id)
	if err != nil {
		return product, err
	}

	err = service.UpdateInteractions(product.ID, product.Interactions)
	if err != nil {
		return product, err
	}

	interaction := dto.Interaction{
		UserID:     uint(userID),
		ProductID:  product.ID,
		CategoryID: product.CategoryID,
	}

	_, err = service.interactionService.Create(interaction)
	if err != nil {
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

func (service *productService) UpdateInteractions(productID uint, currentInteractions uint) error {
	err := service.productRepo.UpdateInteractions(productID, currentInteractions+1)
	if err != nil {
		return ErrFailToInteractsProduct
	}
	return nil
}
