package handler

import (
	"github.com/gofiber/fiber/v2"
)

// handle product creation
// func createProduct(c *fiber.Ctx) error {
// 	var dbcon *gorm.DB = db.GetDBConnection()

// 	// Parse request body
// 	var product dto.Product
// 	if err := c.BodyParser(&product); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "invalid request",
// 		})
// 	}

// 	// Create the product record
// 	if err := dbcon.Create(&product).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "failed to create product record",
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "product created successfully",
// 		"product": product,
// 	})
// }

// handle product retrieval
// func getProducts1(c *fiber.Ctx) error {
// 	var dbcon *gorm.DB = db.GetDBConnection()

// 	var products []dto.Product

// 	// Retrieve all products from the database
// 	if err := dbcon.Find(&products).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "failed to retrieve products",
// 		})
// 	}

// 	return c.JSON(products)
// }

// handle individual product retrieval
// func getProduct(c *fiber.Ctx) error {
// 	// TODO: implement individual product retrieval
// 	return nil
// }

// func getProduct(c *fiber.Ctx) error {

// 	var dbcon *gorm.DB = db.GetDBConnection()

// 	// Get product ID from request params
// 	id := c.Params("id")

// 	// Retrieve the product with the given ID from the database
// 	var product dto.Product
// 	if err := dbcon.Where("id = ?", id).First(&product).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 				"message": "product not found",
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "failed to retrieve product",
// 		})
// 	}

// 	// // Retrieve the top 5 related products using the recommendation algorithm
// 	// var relatedProducts []dto.Product
// 	// // TODO: implement recommendation algorithm to retrieve related products

// 	// // Add the related products to the product object
// 	// product.RelatedProducts = relatedProducts

// 	return c.JSON(product)
// }

// handle tracking of user behavior
func track(c *fiber.Ctx) error {
	// TODO: implement tracking of user behavior
	return nil
}

// register the tracking routes
//app.Post("/track", track)
//handle recommendation retrieval
// func getRecommendations(c *fiber.Ctx) error {
// 	// TODO: implement recommendation retrieval
// 	return nil
// }

/*
func getRecommendations(c *fiber.Ctx) error {
	// Get user ID from context
	userID := c.Locals("userID").(uint)

	// Get user's recent interactions
	var recentInteractions []dto.Interaction
	if err := db.Where("user_id = ?", userID).Order("created_at desc").Limit(10).Find(&recentInteractions).Error; err != nil {
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
		if err := db.First(&category, product.ID).Error; err != nil {
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

	// Sort products by weighted score
	var products []dto.Product
	for productID, price := range productWeights {
		var product dto.Product
		if err := db.First(&product, productID).Error; err != nil {
			continue
		}
		product.Price = price
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

*/
