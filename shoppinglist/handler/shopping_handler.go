package handler

import (
	"database/sql"
	"fmt"

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

func GetShoppingItemByName(c *fiber.Ctx) error {
	name := c.Params("name")
	IsValidItem, found_item := SearchItem(name)
	if IsValidItem {
		return c.Status(fiber.StatusOK).JSON(found_item)
	}
	return c.Status(fiber.StatusNotFound).Send([]byte{})
}

func GetAllItems(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(ShoppingList)
}

func AddNewShoppingItem(c *fiber.Ctx) error {
	var newItem ShoppingItem
	if err := c.BodyParser(&newItem); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Ungültiger Request-Body")
	}
	
	if newItem.Name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Artikelname ist erforderlich")
	}
	query := `INSERT INTO shopping_items (shopping_item, shopping_amount) VALUES ($1, $2) RETURNING shopping_item, shopping_amount`
	
	err := dbContext.QueryRow(query, newItem.Name, newItem.Amount).Scan(&newItem.Name, &newItem.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Fehler beim Einfügen des Artikels: %v", err))
	}

	// Erfolgreiche Antwort mit dem eingefügten Artikel als JSON
	return c.Status(fiber.StatusCreated).JSON(newItem)
}

func UpdateItem(c *fiber.Ctx) error {
	name := c.Params("name")
	is_valid_item, item_found := SearchItem(name)
	if is_valid_item {
		ItemCounter(item_found)
		OutputShoppinglist()
		return c.Status(fiber.StatusOK).JSON(item_found)
	}

	return c.Status(fiber.StatusNotFound).Send([]byte{})
}


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

func SearchItem(name string) (bool, *ShoppingItem) {
	for i := range ShoppingList { 
		if ShoppingList[i].Name == name {
			return true, &ShoppingList[i] 
		}
	}
	return false, nil
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