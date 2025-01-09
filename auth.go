package main

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Very secret key (but very insecure)
const secretKey = "very_secret_key12345678"

// Checks if the user is authenticated and returns their userID. If the user
// isn't authenticated or the JWT token is missing, `0` gets returned with
// an error.
func checkAuthForUserID(c *fiber.Ctx) (uint, error) {
	tokenStr := c.Get(fiber.HeaderAuthorization)
	token, err := verifyJWT(tokenStr)
	if err != nil {
		return 0, errors.New("unauthorized")
	}

	email, _ := token.Claims.GetSubject()
	usr := User{}
	db.Where("email = ?", email).Take(&usr)

	return usr.ID, nil
}

// Creates a JWT token using `email` as the token's subject.
func createJWT(email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
	})
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return ""
	}

	return tokenStr
}

// Verifies that the JWT token is valid.
func verifyJWT(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
