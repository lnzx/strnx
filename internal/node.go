package internal

import (
	"github.com/gofiber/fiber/v2"
)

func GetNodes(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"wallets": "1,2,3,4",
	})
}
