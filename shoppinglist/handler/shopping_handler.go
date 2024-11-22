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
	// SQL-Abfrage, um die ID, den Namen und die Anzahl des Items zu prüfen
	query := "SELECT id, shopping_item, shopping_amount FROM shopping_items WHERE shopping_item = $1"
	
	// Vorbereitung der Rückgabevariablen
	var id int
	var shoppingName string
	var shoppingAmount int
	
	// Abfrage ausführen
	err := dbContext.QueryRow(query, name).Scan(&id, &shoppingName, &shoppingAmount)
	
	// Falls kein Eintrag gefunden wurde, wird `sql.ErrNoRows` zurückgegeben
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		log.Printf("Datenbankfehler: %v", err)
		return false, nil
	}
	
	// Wenn das Item gefunden wurde, die Details in einer Map speichern
	foundItem := map[string]interface{}{
		"id":              id,
		"name":   shoppingName,
		"amount": shoppingAmount,
	}
	return true, foundItem
}

// @Summary Get a shopping item by name
// @Description Get details of a specific shopping item
// @ID get-item-by-name
// @Produce json
// @Param name path string true "Name of the shopping item"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Item not found"
// @Router /api/shopping/{name} [get]
func GetShoppingItemByName(c *fiber.Ctx) error {
	name := c.Params("name")
	IsValidItem, foundItem := SearchItem(name)
	if IsValidItem {
		return c.Status(fiber.StatusOK).JSON(foundItem)
	}
	return c.Status(fiber.StatusNotFound).Send([]byte("Item nicht gefunden"))
}
// @Summary Get all shopping items
// @Description Get a list of all shopping items
// @ID get-all-items
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {string} string "Error fetching items"
// @Router /api/shopping [get]]
func GetAllItems(c *fiber.Ctx) error {
	// SQL-Abfrage, um alle Items abzurufen
	query := "SELECT id, shopping_item, shopping_amount FROM shopping_items"
	rows, err := dbContext.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Fehler beim Abrufen der Items")
	}
	defer rows.Close()

	// Slice, um die Items zu speichern
	var shoppingList []map[string]interface{}

	// Iteration über die Ergebnisse der Abfrage
	for rows.Next() {
		var id int
		var shoppingItem string
		var shoppingAmount int

		// Daten aus der aktuellen Zeile abrufen
		if err := rows.Scan(&id, &shoppingItem, &shoppingAmount); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Fehler beim Verarbeiten der Daten")
		}

		// Item als Map speichern und zur Liste hinzufügen
		item := map[string]interface{}{
			"id":              id,
			"name":   shoppingItem,
			"amount": shoppingAmount,
		}
		shoppingList = append(shoppingList, item)
	}

	// Fehler bei der Zeilenverarbeitung prüfen
	if err = rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Fehler beim Verarbeiten der Daten")
	}
	if len(shoppingList) == 0 {
		return c.Status(fiber.StatusOK).JSON([]string{})
	} 
	// Liste aller Items als JSON zurückgeben
	return c.Status(fiber.StatusOK).JSON(shoppingList)
}

// @Summary Add a new shopping item
// @Description Add a new item to the shopping list
// @ID add-new-item
// @Accept json
// @Produce json
// @Param body body ShoppingItem true "New Shopping Item"
// @Success 201 {object} ShoppingItem
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Error adding item"
// @Router /api/shopping [post]
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

// @Summary Update a shopping item by name
// @Description Update the amount of an existing shopping item
// @ID update-item
// @Accept json
// @Produce json
// @Param name path string true "Name of the shopping item"
// @Param body body ShoppingItem true "Updated Shopping Item"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Item not found"
// @Failure 500 {string} string "Error updating item"
// @Router /api/shopping/{name} [put]
func UpdateItem(c *fiber.Ctx) error {
	name := c.Params("name")
	var newItem ShoppingItem

	// Body des Requests parsen
	if err := c.BodyParser(&newItem); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Ungültiger Request-Body")
	}

	// Überprüfen, ob das Item existiert
	isValidItem, itemFound := SearchItem(name)
	if isValidItem {
		// Auf die ID des Items zugreifen
		id := itemFound["id"].(int)

		// Hier kannst du die Logik einfügen, um das Item in der Datenbank mit der neuen Menge zu aktualisieren
		// Zum Beispiel könnte die shopping_amount aktualisiert werden:
		query := "UPDATE shopping_items SET shopping_amount = $1 WHERE id = $2"
		_, err := dbContext.Exec(query, newItem.Amount, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Fehler beim Aktualisieren des Items")
		}

		// Aktualisierte Informationen zurückgeben
		itemFound["shopping_amount"] = newItem.Amount
		return c.Status(fiber.StatusOK).JSON(itemFound)
	}

	return c.Status(fiber.StatusNotFound).Send([]byte("Item nicht gefunden"))
}

// @Summary Get a shopping item by name
// @Description Get details of a specific shopping item
// @ID get-item-by-name
// @Produce json
// @Param name path string true "Name of the shopping item"
// @Success 200 {object} ShoppingItem
// @Router /api/shopping/{name} [delete]
func DeleteShoppingItem(c *fiber.Ctx) error {
	name := c.Params("name")
	
	// Überprüfen, ob das Item existiert
	isValidItem, itemFound := SearchItem(name)
	if isValidItem {
		// Auf die ID des Items zugreifen
		id := itemFound["id"].(int)

		// SQL-Delete-Query, um das Item zu löschen
		query := "DELETE FROM shopping_items WHERE id = $1"
		_, err := dbContext.Exec(query, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Fehler beim Löschen des Items")
		}

		// Erfolgreiche Löschbestätigung
		return c.Status(fiber.StatusOK).SendString("Item erfolgreich gelöscht")
	}

	return c.Status(fiber.StatusNotFound).Send([]byte("Item nicht gefunden"))
}
// @Summary Hello World
// @Description Simple hello world endpoint
// @ID hello-world
// @Success 200 {string} string "Hello World"
// @Router /hello [get]
func SayHello(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Hello World")
}