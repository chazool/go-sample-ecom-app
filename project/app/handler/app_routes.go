package handler

import "github.com/gofiber/fiber/v2"

// register the routes for the application
func RegisterRoutes(app *fiber.App) {

	categoryHandler := NewCategoryHandler()
	productHandler := NewProductHandler()
	authHandler := NewAuthHandler()

	// register the authentication routes
	auth := app.Group("/auth")
	auth.Post("/login", authHandler.login)
	auth.Post("/register", authHandler.register)

	// management routes
	manage := app.Group("/manage")
	manage.Use(authHandler.authenticationMiddleware) // require authentication for product routes

	// manage products
	manageProduct := manage.Group("/product")
	manageProduct.Post("/", productHandler.createProduct)
	manageProduct.Delete("/:id", productHandler.deleteProduct)

	// manage category
	manageCategory := manage.Group("/category")
	manageCategory.Post("/", categoryHandler.createCategory)

	// searching routes
	search := app.Group("/search")

	// search product
	searchProduct := search.Group("/product")
	searchProduct.Get("/", productHandler.getProducts)
	searchProduct.Get("/:id", productHandler.getProduct)

	// search category
	searCategory := search.Group("/category")
	searCategory.Get("/", categoryHandler.getCategories)
	searCategory.Get("/:id", categoryHandler.getCategory)

	// register the tracking routes
	app.Post("/track", track)

	// register the recommendation routes
	//app.Get("/recommendation/:id", getRecommendations)
}
