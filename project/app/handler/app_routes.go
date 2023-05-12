package handler

import "github.com/gofiber/fiber/v2"

// register the routes for the application
func RegisterRoutes(app *fiber.App) {

	productHandler := NewProductHandler()
	authHandler := NewAuthHandler()

	// register the authentication routes
	auth := app.Group("/auth")
	auth.Post("/login", authHandler.login)
	auth.Post("/register", authHandler.register)

	// register the product routes
	product := app.Group("/product")
	product.Use(authHandler.authenticationMiddleware) // require authentication for product routes
	product.Post("/", productHandler.createProduct)
	product.Get("/", productHandler.getProducts)
	product.Get("/:id", productHandler.getProduct)
	product.Delete("/:id", productHandler.deleteProduct)

	// register the tracking routes
	app.Post("/track", track)

	// register the recommendation routes
	//app.Get("/recommendation/:id", getRecommendations)
}
