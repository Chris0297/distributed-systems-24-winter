package main

import (
	"fmt"
	"shoppinglist/handler"
	"log"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/api/shopping/:name", handler.GetShoppingItemByName)
	fmt.Println("hello")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	
}
