package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/models"
)

type SessionUser struct {
	ID       uuid.UUID
	Username string
	Email    string
	IsAdmin  bool
}

type CustomClaims struct {
	SessionUser
	jwt.RegisteredClaims
}

func CreateJWT(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := &CustomClaims{
		SessionUser: SessionUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			IsAdmin:  user.IsAdmin,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (*CustomClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := new(CustomClaims)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GetSessionUser(c *fiber.Ctx) SessionUser {
	user := c.Locals("user").(SessionUser)
	if user.Username == "" {
		panic("GetSessionUser is not allowed outside protected route")
	}
	return user
}
