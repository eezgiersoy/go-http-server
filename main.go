package main

import (
	"awesomeProject/db"
	"awesomeProject/routes"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	db.ConnectDb()
	app := fiber.New()

	setupRoutes(app)

	app.Get("/api", welcome)

	log.Fatal(app.Listen(":3000"))
}

func welcome(ctx *fiber.Ctx) error {
	return ctx.SendString("Welcome to my awesome API")
}

func setupRoutes(app *fiber.App) {
	// welcome endpoint
	app.Get("/api", welcome)

	// User endpoints
	app.Post("/api/users", routes.CreateUser)
	app.Get("/api/users", routes.GetUsers)
	app.Get("/api/users/:id", routes.GetUser)
	app.Patch("/api/users/:id", routes.UpdateUser)
	app.Delete("/api/users/:id", routes.DeleteUser)

	//Product endpoints
	app.Post("/api/products", routes.CreateProduct)
	app.Get("/api/products", routes.GetProducts)
	app.Get("/api/products/:id", routes.GetProduct)
	app.Patch("/api/products/:id", routes.UpdateProduct)
	app.Delete("/api/products/:id", routes.DeleteProduct)

	//Order endpoints
	app.Post("/api/orders", routes.CreateOrder)
	app.Get("/api/orders", routes.GetOrders)
	app.Get("/api/orders/:id", routes.GetOrder)
}
