package handler

import (
	"errors"
	"log"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"sample_app/app/dto"
	"sample_app/app/services"
	"sample_app/pkg/config/db"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler() ProductHandler {
	return ProductHandler{
		productService: services.NewProductService(),
	}
}

// handle product deletion
func (h *ProductHandler) deleteProduct(c *fiber.Ctx) error {
	// Log the incoming request
	log.Printf("Received delete request for product ID: %v", c.Params("id"))

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		// Log the error
		log.Printf("Error parsing product ID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
		})
	}

	err = h.productService.Delete(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			// Log the error
			log.Printf("Product not found: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": services.ErrProductNotFound.Error(),
			})
		}
		// Log the error
		log.Printf("Error deleting product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Log the success message
	log.Printf("Product with ID %v deleted successfully", id)

	return c.JSON(dto.Response{
		Message: "Product deleted successfully",
	})
}

// handle product retrieval
func (h *ProductHandler) getProducts(c *fiber.Ctx) error {

	// Log the incoming request
	log.Println("Received get all products request")

	products, err := h.productService.FindAll()
	if err != nil {
		// Log the error
		log.Printf("Error retrieving products: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Log the success message
	log.Printf("Found %v products", len(products))

	return c.JSON(dto.Response{
		Data: products,
	})
}

func (h *ProductHandler) getProduct(c *fiber.Ctx) error {

	// Log the incoming request
	log.Printf("Received get product request for ID: %v", c.Params("id"))

	// Get the authenticated user from the context
	user := c.Locals("user").(*dto.User)

	// Get product ID from request params
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		// Log the error
		log.Printf("Error parsing product ID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
		})
	}

	product, err := h.productService.FindByIdWithInteract(uint(id), user.ID)
	if err != nil {
		// Log the error
		log.Printf("Error retrieving product: %v", err)
		if errors.Is(err, services.ErrProductNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": services.ErrProductNotFound.Error(),
			})
		}
		// Log the error
		log.Printf("Error retrieving product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Log the success message
	log.Printf("Found product with ID %v", id)

	return c.JSON(dto.Response{
		Data: product,
	})
}

// handle product creation
func (h *ProductHandler) createProduct(c *fiber.Ctx) error {

	// Log the incoming request
	log.Println("Received create product request")

	// Parse request body
	var product dto.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	// Log the received product
	log.Printf("Received product: %+v\n", product)

	// Create the product record
	createdProduct, err := h.productService.Create(product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	// Log the created product
	log.Printf("Created product: %+v\n", createdProduct)

	return c.JSON(dto.Response{
		Message: "Product created successfully",
		Data:    createdProduct,
	})
}

func (h *ProductHandler) getRecommendations(c *fiber.Ctx) error {

	var db *gorm.DB = db.GetDBConnection()
	// Get user ID from context
	user := c.Locals("user").(*dto.User)

	// Get user's recent interactions
	var recentInteractions []dto.Interaction
	if err := db.Where("user_id = ?", user.ID).Order("created_at desc").Limit(10).Find(&recentInteractions).Error; err != nil {
		return err
	}

	// Group interactions by product ID
	interactionsByProduct := make(map[uint][]dto.Interaction)
	for _, interaction := range recentInteractions {
		interactionsByProduct[interaction.ProductID] = append(interactionsByProduct[interaction.ProductID], interaction)
	}

	// Calculate weighted scores for each product
	productWeights := make(map[uint]float64)
	for productID, interactions := range interactionsByProduct {
		var totalInteractions int
		var inCategoryInteractions int
		var categoryID uint

		// Get product details
		var product dto.Product
		if err := db.First(&product, productID).Error; err != nil {
			continue
		}

		// Get category details
		var category dto.Category
		if err := db.First(&category, product.CategoryID).Error; err != nil {
			continue
		}
		categoryID = category.ID

		// Calculate number of interactions and interactions within category
		for _, interaction := range interactions {
			totalInteractions++
			if interaction.CategoryID == categoryID {
				inCategoryInteractions++
			}
		}

		// Calculate weighted score based on number of interactions and interactions within category
		if totalInteractions > 0 {
			multiplier := float64(inCategoryInteractions) / float64(totalInteractions)
			weightedScore := float64(product.Interactions) * multiplier
			productWeights[productID] = weightedScore
		}
	}

	// Sort products by Price
	var products []dto.Product
	for productID, score := range productWeights {
		var product dto.Product
		if err := db.First(&product, productID).Error; err != nil {
			continue
		}
		product.WeightedScore = score
		products = append(products, product)
	}
	sort.Slice(products, func(i, j int) bool {
		return products[i].Price > products[j].Price
	})

	// Return top 5 products
	if len(products) > 5 {
		products = products[:5]
	}

	// Return response
	return c.JSON(products)
}
