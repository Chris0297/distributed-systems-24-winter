package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var ShoppingList []Items

type Items struct {
	name string
	amount int
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
	TestItem := Items{} 
	name := c.Params("name")
	TestItem.name = name
	TestItem.amount = 1
	ShoppingList = append(ShoppingList,TestItem)
	OutputShoppinglist()
	return c.Status(fiber.StatusCreated).JSON("Item successfully created.")
}

func SearchItem(name string) bool {
	for _, item := range ShoppingList {
		if item.name == name {
			return true
		}
	}
	return false
}

func UpdateAmount(c *fiber.Ctx) error {
	name := c.Params("name")
	ItemCounter(name)
	OutputShoppinglist()
	return c.Status(fiber.StatusOK).JSON("Item updated successfully.")
}


func ItemCounter(name string) {
	for i := range ShoppingList {
		if ShoppingList[i].name == name {
			ShoppingList[i].amount++
		}
	}
}

func OutputShoppinglist(){
	fmt.Printf("%v\n", ShoppingList)
	fmt.Println(ShoppingList)
}