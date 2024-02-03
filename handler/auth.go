package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mladenovic-13/bank-api/models"
	"github.com/mladenovic-13/bank-api/utils"
)

type SigninReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary User sign-in
// @Description Handles user sign-in.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body SigninReq true "Body"
// @Success 200 {object} models.User
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /auth/signin [post]
func (ctx *HandlerContext) HandleSignin(c *fiber.Ctx) error {
	credentials := new(SigninReq)

	if err := c.BodyParser(credentials); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user := new(models.User)

	ctx.DB.First(user, "username = ?", credentials.Username)

	if user.Username == credentials.Username {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "User already exists"})
	}

	hashedPassword, err := utils.HashPassword(credentials.Password)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	newUser := models.NewUser(
		credentials.Username,
		credentials.Email,
		hashedPassword,
	)

	result := ctx.DB.Create(newUser)

	if result.Error != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(newUser)
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginRes struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// @Summary User log-in
// @Description Handles user log-in.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param body body LoginReq true "Body"
// @Success 200 {object} LoginRes
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /auth/login [post]
func (ctx *HandlerContext) HandleLogin(c *fiber.Ctx) error {
	credentials := new(LoginReq)

	if err := c.BodyParser(credentials); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user := new(models.User)

	result := ctx.DB.First(user, "username = ?", credentials.Username)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": "User does not exit"})
	}

	if utils.CheckPasswordHash(credentials.Password, user.Password) {
		token, err := utils.CreateJWT(user)

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(
			LoginRes{
				Token: token,
				User:  *user,
			},
		)
	}

	return c.
		Status(fiber.StatusBadRequest).
		JSON(fiber.Map{"error": "Wrong password"})
}

// @Summary User log-out
// @Description Handles user log-out.
// @Tags Authentication
// @Security Bearer Token
// @Success 200
// @Failure 401 {object} Error
// @Router /auth/logout [post]
func (ctx *HandlerContext) HandleLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

	return c.SendStatus(fiber.StatusOK)
}
