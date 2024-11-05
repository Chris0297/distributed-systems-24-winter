package main

import (
	"log"
	"os"
	"shoppinglist/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/swaggo/fiber-swagger" // http://localhost:3000/swagger/index.html
    _ "shoppinglist/docs" // Hier sicherstellen, dass der Pfad korrekt ist
)
// @title Shopping List API
// @version 1.0
// @description This is the API for managing a shopping list.
// @host localhost:8080
// @BasePath /
func main() {
	//Setting port via enviorment variable at the docker run command -e PORT=XXXX
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT enviorment variable is not set")
	}
	apiShoppingWithName := "/api/shopping/:name"
	
	app := fiber.New(fiber.Config{
		Immutable: true,
		
	})
	// Enable Cors to avoid Cors Policy Problems
	app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,PUT,DELETE",
        AllowHeaders: "Content-Type, Authorization",
    }))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Get(apiShoppingWithName, handler.GetShoppingItemByName)
	app.Get("/api/shopping", handler.GetAllItems)
	app.Post(apiShoppingWithName, handler.AddNewShoppingItem)
	app.Put(apiShoppingWithName, handler.UpdateItem)
	app.Delete(apiShoppingWithName, handler.DeleteShoppingItem)
	if err := app.Listen(":"+port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// if err := app.Listen(":3000"); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
	
	
}
