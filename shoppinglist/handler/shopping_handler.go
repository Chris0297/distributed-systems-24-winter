package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

var ShoppingList []ShoppingItem

type ShoppingItem struct {
	Name string `json:"name"`
	Amount int  `json:"amount"`
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
	name := c.Params("name")
	neues_item := ShoppingItem{Name: name, Amount: 1}
	ShoppingList = append(ShoppingList,neues_item)
	OutputShoppinglist()
	return c.Status(fiber.StatusCreated).JSON(neues_item)
}

func UpdateItem(c *fiber.Ctx) error {
	name := c.Params("name")
	is_valid_item, item_found := SearchItem(name)
	if is_valid_item {
		ItemCounter(item_found.Name)
		OutputShoppinglist()
		return c.Status(fiber.StatusOK).JSON(item_found)
	}
	
	return c.Status(fiber.StatusNotFound).Send([]byte{})
}

func DeleteShoppingItem(c *fiber.Ctx) error {
    name := c.Params("name")
    IsDeleted := DeleteItem(name)
    if IsDeleted{
        return c.Status(fiber.StatusOK).Send([]byte{})
    }
        return c.Status(fiber.StatusNoContent).Send([]byte{})
}

func SearchItem(name string) (bool, *ShoppingItem) {
	for _, item := range ShoppingList {
		if item.Name == name {
			return true, &item
		}
	}
	return false,nil 
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


func ItemCounter(name string) {
	for i := range ShoppingList {
		if ShoppingList[i].Name == name {
			ShoppingList[i].Amount++
		}
	}
}

func OutputShoppinglist(){
	fmt.Println(ShoppingList)
}