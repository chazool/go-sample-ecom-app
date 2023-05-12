package handler

import (
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sample_app/app/dto"
	"sample_app/app/services"
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

	err = h.productService.Delete(id)
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

	// Get product ID from request params
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		// Log the error
		log.Printf("Error parsing product ID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
		})
	}

	product, err := h.productService.FindById(id)
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
