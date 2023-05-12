package services

import (
	"errors"
	"log"
	"sample_app/app/dto"
	"sample_app/app/repository"
	"sort"
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
	Delete(id uint) error
	FindByIdWithInteract(id, userID uint) (dto.Product, error)
	FindAll() ([]dto.Product, error)
	GetRecommendations(userID uint) ([]dto.Product, error)
	findById(id uint) (dto.Product, error)
	calculateProductWeights(interactions []dto.Interaction) (map[uint]float64, error)
	getTopProducts(productWeights map[uint]float64, returnCount int) ([]dto.Product, error)
}

type productService struct {
	productRepo        repository.ProductRepository
	interactionService InteractionService
	categoryService    CategoryService
}

func NewProductService() ProductService {
	return &productService{
		productRepo:        repository.NewProductRepository(),
		categoryService:    NewCategoryService(),
		interactionService: NewInteractionService(),
	}
}

func (service *productService) Delete(id uint) error {

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

func (service *productService) findById(id uint) (dto.Product, error) {

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

func (service *productService) FindByIdWithInteract(id, userID uint) (dto.Product, error) {

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

func (service *productService) GetRecommendations(userID uint) ([]dto.Product, error) {

	// Get user's recent interactions
	recentInteractions, err := service.interactionService.GetRecentInteractions(userID, 10)
	if err != nil {
		return nil, err
	}

	// Calculate weighted scores for each product
	productWeights, err := service.calculateProductWeights(recentInteractions)
	if err != nil {
		return nil, err
	}

	// Sort products by weighted score
	products, err := service.getTopProducts(productWeights, 5)
	if err != nil {
		return nil, err
	}

	// Return top 5 products
	return products, nil
}

func (service *productService) calculateProductWeights(interactions []dto.Interaction) (map[uint]float64, error) {

	productWeights := make(map[uint]float64)
	for _, interaction := range interactions {

		product, err := service.findById(interaction.ProductID)
		if err != nil {
			continue
		}

		category, err := service.categoryService.FindById(product.ID)
		if err != nil {
			continue
		}
		totalInteractions := len(interactions)
		inCategoryInteractions := 0
		for _, interaction2 := range interactions {
			if interaction2.CategoryID == category.ID {
				inCategoryInteractions++
			}
		}
		if totalInteractions > 0 {
			multiplier := float64(inCategoryInteractions) / float64(totalInteractions)
			weightedScore := float64(product.Interactions) * multiplier
			productWeights[product.ID] = weightedScore
		}
	}
	return productWeights, nil
}

func (service *productService) getTopProducts(productWeights map[uint]float64, returnCount int) ([]dto.Product, error) {

	var products []dto.Product
	for productID, score := range productWeights {

		product, err := service.findById(productID)
		if err != nil {
			continue
		}

		product.WeightedScore = score
		products = append(products, product)
	}
	sort.Slice(products, func(i, j int) bool {
		return products[i].WeightedScore > products[j].WeightedScore
	})
	if len(products) > returnCount {
		products = products[:returnCount]
	}
	return products, nil
}
