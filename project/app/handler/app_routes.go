package handler

import "github.com/gofiber/fiber/v2"

// register the routes for the application
func RegisterRoutes(app *fiber.App) {

	categoryHandler := NewCategoryHandler()
	productHandler := NewProductHandler()
	authHandler := NewAuthHandler()

	// register the authentication routes
	auth := app.Group("/auth")
	app.Post("/login", authHandler.login)
	auth.Post("/register", authHandler.register)

	// management routes
	manage := app.Group("/manage")
	manage.Use(authHandler.createAuthenticationMiddleware)

	// manage products
	manageProduct := manage.Group("/product")
	manageProduct.Post("/", productHandler.createProduct)
	manageProduct.Delete("/:id", productHandler.deleteProduct)

	// manage category
	manageCategory := manage.Group("/category")
	manageCategory.Post("/", categoryHandler.createCategory)

	// searching routes
	search := app.Group("/search")
	search.Use(authHandler.searchAuthenticationMiddleware)
	// search product
	searchProduct := search.Group("/product")
	searchProduct.Get("/recommendation", productHandler.getRecommendations)
	searchProduct.Get("/", productHandler.getProducts)
	searchProduct.Get("/:id", productHandler.getProduct)

	// search category
	searCategory := search.Group("/category")
	searCategory.Get("/", categoryHandler.getCategories)
	searCategory.Get("/:id", categoryHandler.getCategory)

}
