package handler

import (
	"github.com/gofiber/fiber/v2"
)

func GetShoppingItemByName(c *fiber.Ctx) error {
    //name := c.Params("name")
	
    return c.Status(fiber.StatusOK).JSON("Item found and retrieved successfully.")
}