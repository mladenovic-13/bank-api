package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mladenovic-13/bank-api/utils"
)

func (ctx *MiddlewareContext) WithAdmin(c *fiber.Ctx) error {
	user := utils.GetSessionUser(c)

	log.Println("ID: ", user.ID)
	log.Println("Username: ", user.Username)
	log.Println("Is Admin: ", user.IsAdmin)

	if user.Username == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if !user.IsAdmin {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}
