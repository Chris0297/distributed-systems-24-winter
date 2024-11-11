package handler

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)



var ShoppingList []ShoppingItem
var dbContext *sql.DB

type ShoppingItem struct {
	Name string `json:"name"`
	Amount int  `json:"amount"`
}

func InitDbContext(db *sql.DB){
	dbContext = db
}

// Suchfunktion, um ein Item mit dem angegebenen Namen zu finden
func SearchItem(name string) (bool, map[string]interface{}) {
	// SQL-Abfrage, um zu prüfen, ob das Item existiert
	query := "SELECT shopping_item, shopping_amount FROM shopping_items WHERE shopping_item = $1"
	
	// Vorbereitung der Rückgabevariablen
	var shoppingItem string
	var shoppingAmount int
	
	// Abfrage ausführen
	err := dbContext.QueryRow(query, name).Scan(&shoppingItem, &shoppingAmount)
	
	// Falls kein Eintrag gefunden wurde, wird `sql.ErrNoRows` zurückgegeben
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		log.Printf("Datenbankfehler: %v", err)
		return false, nil
	}
	
	// Wenn das Item gefunden wurde, die Details in einer Map speichern
	foundItem := map[string]interface{}{
		"shopping_item":  shoppingItem,
		"shopping_amount": shoppingAmount,
	}
	return true, foundItem
}

// Handler-Funktion, um das Item nach Namen abzurufen
func GetShoppingItemByName(c *fiber.Ctx) error {
	name := c.Params("name")
	IsValidItem, foundItem := SearchItem(name)
	if IsValidItem {
		return c.Status(fiber.StatusOK).JSON(foundItem)
	}
	return c.Status(fiber.StatusNotFound).Send([]byte("Item nicht gefunden"))
}






func GetAllItems(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(ShoppingList)
}

func AddNewShoppingItem(c *fiber.Ctx) error {
	var newItem ShoppingItem

	// Body des Requests parsen
	if err := c.BodyParser(&newItem); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Ungültiger Request-Body")
	}

	// Überprüfen, ob der Artikelname gesetzt ist
	if newItem.Name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Artikelname ist erforderlich")
	}

	// Zuerst prüfen, ob der Artikel bereits in der Datenbank existiert
	var existingAmount int
	err := dbContext.QueryRow("SELECT shopping_amount FROM shopping_items WHERE shopping_item = $1", newItem.Name).Scan(&existingAmount)

	// Wenn der Artikel bereits existiert, erhöhen wir den Amount
	if err == nil { // Artikel existiert
		// Artikel existiert, Amount erhöhen
		newItem.Amount = existingAmount + newItem.Amount // Hinzufügen der neuen Menge zum bestehenden Betrag

		// Artikel aktualisieren
		_, err := dbContext.Exec("UPDATE shopping_items SET shopping_amount = $1 WHERE shopping_item = $2", newItem.Amount, newItem.Name)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Fehler beim Aktualisieren des Artikels: %v", err))
		}

		// Erfolgreiche Antwort mit dem aktualisierten Artikel als JSON
		return c.Status(fiber.StatusOK).JSON(newItem)
	}

	// Wenn der Artikel nicht existiert, fügen wir ihn hinzu
	if err == sql.ErrNoRows {
		// Artikel hinzufügen
		query := `INSERT INTO shopping_items (shopping_item, shopping_amount) VALUES ($1, $2) RETURNING shopping_item, shopping_amount`
		err := dbContext.QueryRow(query, newItem.Name, newItem.Amount).Scan(&newItem.Name, &newItem.Amount)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Fehler beim Einfügen des Artikels: %v", err))
		}

		// Erfolgreiche Antwort mit dem eingefügten Artikel als JSON
		return c.Status(fiber.StatusCreated).JSON(newItem)
	}

	// Falls ein anderer Fehler auftritt
	return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Fehler beim Überprüfen des Artikels: %v", err))
}

// func UpdateItem(c *fiber.Ctx) error {
// 	name := c.Params("name")
// 	is_valid_item, item_found := SearchItem(name)
// 	if is_valid_item {
// 		ItemCounter(item_found)
// 		OutputShoppinglist()
// 		return c.Status(fiber.StatusOK).JSON(item_found)
// 	}

// 	return c.Status(fiber.StatusNotFound).Send([]byte{})
// }


// @Summary Get a shopping item by name
// @Description Get details of a specific shopping item
// @ID get-item-by-name
// @Produce json
// @Param name path string true "Name of the shopping item"
// @Success 200 {object} ShoppingItem
// @Router /item/{name} [get]
func DeleteShoppingItem(c *fiber.Ctx) error {
    name := c.Params("name")
    IsDeleted := DeleteItem(name)
    if IsDeleted{
        return c.Status(fiber.StatusOK).Send([]byte{})
    }
        return c.Status(fiber.StatusNoContent).Send([]byte{})
}

func DeleteItem(name string) bool  {
    for i, item := range ShoppingList {
        if item.Name == name {
            ShoppingList = append(ShoppingList[:i], ShoppingList[i+1:]... )
            return true
        }
    }
    return false
}


func ItemCounter(item *ShoppingItem)  {
	item.Amount++
}

func OutputShoppinglist(){
	fmt.Println(ShoppingList)
}