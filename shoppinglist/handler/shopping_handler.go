package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

var ShoppingList []Items

type Items struct {
	Name string
	Amount int
}

func GetShoppingItemByName(c *fiber.Ctx) error {
	name := c.Params("name")
	IsValidItem := SearchItem(name)
	if IsValidItem {
		return c.Status(fiber.StatusOK).JSON("Item found and retrieved successfully." + name)
	}
	return c.Status(fiber.StatusNotFound).JSON("Item not found")
}

func AddNewShoppingItem(c *fiber.Ctx) error {
	name := c.Params("name")
	neues_item := Items{Name: name, Amount: 1}
	ShoppingList = append(ShoppingList,neues_item)
	OutputShoppinglist()
	return c.Status(fiber.StatusCreated).JSON("Item successfully created.")
}

func UpdateAmount(c *fiber.Ctx) error {
	name := c.Params("name")
	ItemCounter(name)
	OutputShoppinglist()
	return c.Status(fiber.StatusOK).JSON("Item updated successfully.")
}

func DeleteShoppingItem(c *fiber.Ctx) error {
    name := c.Params("name")
    IsDeleted := DeleteItem(name)
    if IsDeleted{
        return c.Status(fiber.StatusOK).JSON("Item deleted successfully.")
    }
        return c.Status(fiber.StatusNoContent).JSON("Item not found.")
}

func SearchItem(name string) bool {
	for _, item := range ShoppingList {
		if item.Name == name {
			return true
		}
	}
	return false
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
	fmt.Printf("%v\n", ShoppingList)
	fmt.Println(ShoppingList)
}