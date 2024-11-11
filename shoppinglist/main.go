package main

import (
"database/sql"
"fmt"
"log"
"os"
_ "shoppinglist/docs" // Hier sicherstellen, dass der Pfad korrekt ist
"shoppinglist/handler"

"github.com/gofiber/fiber/v2"
"github.com/gofiber/fiber/v2/middleware/cors"
"github.com/swaggo/fiber-swagger" // http://localhost:3000/swagger/index.html
_"github.com/lib/pq"
)
var db *sql.DB
// @title Shopping List API
// @version 1.0
// @description This is the API for managing a shopping list.
// @host localhost:8080
// @BasePath /
func main() {
var err error
//Setting port via enviorment variable at the docker run command -e PORT=XXXX
port := os.Getenv("PORT")


Database_Username :=  os.Getenv("DB_USER")
Database_Name := os.Getenv("DB_NAME")
Database_Password :=  os.Getenv("DB_PASSWORD")
Database_Port := 5432;

Database_Host := "db";

psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
Database_Host, Database_Port, Database_Username, Database_Password, Database_Name)
    db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal("Error opening database: ", err)
    }
    
    handler.InitDbContext(db)
    defer db.Close()

    // Überprüfen, ob die Verbindung funktioniert
    err = db.Ping()
    if err != nil {
        log.Fatal("Error connecting to the database: ", err)
    }
    fmt.Println("Successfully connected to the database!")


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
app.Get(apiShoppingWithName, handler.GetShoppingItemByName) //CHECK
app.Get("/api/shopping", handler.GetAllItems)
app.Post("/api/shopping", handler.AddNewShoppingItem)
app.Put(apiShoppingWithName, handler.UpdateItem)
app.Delete(apiShoppingWithName, handler.DeleteShoppingItem)
if err := app.Listen(":"+port); err != nil {
log.Fatalf("Failed to start server: %v", err)
}
// if err := app.Listen(":3000"); err != nil {
// log.Fatalf("Failed to start server: %v", err)
// }


}