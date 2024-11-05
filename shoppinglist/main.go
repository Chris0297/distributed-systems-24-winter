package main

import (
	"shoppinglist/handler"
	"log"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})
	app.Get("/api/shopping/:name", handler.GetShoppingItemByName)
	app.Get("/api/shopping", handler.GetAllItems)
	app.Post("/api/shopping/:name", handler.AddNewShoppingItem)
	app.Put("/api/shopping/:name", handler.UpdateAmount)
	app.Delete("api/shopping/:name", handler.DeleteShoppingItem)
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	
}
