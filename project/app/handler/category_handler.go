package handler

import (
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sample_app/app/dto"
	"sample_app/app/services"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler() CategoryHandler {
	return CategoryHandler{
		categoryService: services.NewCategoryService(),
	}
}

// handle category retrieval
func (h *CategoryHandler) getCategories(c *fiber.Ctx) error {

	// Log the incoming request
	log.Println("Received get all categories request")

	categories, err := h.categoryService.FindAll()
	if err != nil {
		// Log the error
		log.Printf("Error retrieving categories: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Log the success message
	log.Printf("Found %v categories", len(categories))

	return c.JSON(dto.Response{
		Data: categories,
	})
}

func (h *CategoryHandler) getCategory(c *fiber.Ctx) error {

	// Log the incoming request
	log.Printf("Received get category request for ID: %v", c.Params("id"))

	// Get category ID from request params
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		// Log the error
		log.Printf("Error parsing category ID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid category ID",
		})
	}

	category, err := h.categoryService.FindById(uint(id))
	if err != nil {
		// Log the error
		log.Printf("Error retrieving category: %v", err)
		if errors.Is(err, services.ErrCategoryNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": services.ErrCategoryNotFound.Error(),
			})
		}
		// Log the error
		log.Printf("Error retrieving category: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Log the success message
	log.Printf("Found category with ID %v", id)

	return c.JSON(dto.Response{
		Data: category,
	})
}

// handle category creation
func (h *CategoryHandler) createCategory(c *fiber.Ctx) error {

	// Log the incoming request
	log.Println("Received create category request")

	// Parse request body
	var category dto.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	// Log the received category
	log.Printf("Received category: %+v\n", category)

	// Create the category record
	createdCategory, err := h.categoryService.Create(category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	// Log the created category
	log.Printf("Created category: %+v\n", createdCategory)

	return c.JSON(dto.Response{
		Message: "Category created successfully",
		Data:    createdCategory,
	})
}
