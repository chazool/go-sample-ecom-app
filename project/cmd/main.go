package main

import (
	"log"
	"sample_app/app/handler"
	"sample_app/pkg/config/db"

	"github.com/gofiber/fiber/v2"
)

func init() {
	err := db.InitDBConnection()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	start()
}

// initialize the fiber application
func start() {
	app := fiber.New()

	// register the routes
	handler.RegisterRoutes(app)

	// start the server
	err := app.Listen(":8181")
	if err != nil {
		log.Fatal(err)
	}

}
