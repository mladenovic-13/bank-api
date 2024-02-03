package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mladenovic-13/bank-api/models"
	"github.com/mladenovic-13/bank-api/utils"
)

func (ctx *MiddlewareContext) WithAuth(c *fiber.Ctx) error {
	token := getToken(c.Get("Authorization"))

	if token == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims, err := utils.ValidateJWT(token)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user := new(models.User)

	res := ctx.DB.First(user, "id = ?", claims.SessionUser.ID.String())

	if res.Error != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("user", claims.SessionUser)

	return c.Next()
}

func getToken(cookie string) string {
	strings := strings.Split(cookie, " ")

	if len(strings) != 2 {
		return ""
	}

	if strings[0] != "Bearer" {
		return ""
	}

	return strings[1]
}
