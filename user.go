package main

import (
	"bytes"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

type (
	// GORM model for the users.
	User struct {
		gorm.Model
		Name     string
		Email    string
		Password []byte
		Todos    []Todo
	}

	// Used for parsing JSON for user registration.
	UserRegistration struct {
		Name     string `json:"name"     validate:"required,min=2"`
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	// Used for parsing JSON for user authentication.
	UserLogin struct {
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
)

// Creates a new user object. The given `password` gets hashed before it's
// placed into the `User` object.
func newUser(name, email, password string) User {
	return User{
		Name:     name,
		Email:    email,
		Password: hashPassword(password),
	}
}

// Used as a handler of the `/register` route. Registeres a user.
func UserRegistrationHandler(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMEApplicationJSON)

	register := UserRegistration{}
	if err := parseAndValidate(c, &register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if isUserRegistered(register.Email) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "user with that email already exists",
		})
	}

	usr := newUser(register.Name, register.Email, register.Password)
	db.Create(&usr)

	return createAndReturnToken(c, usr.Email)
}

// Used as a handler of the `/login` route. Authenticates the user.
func UserLoginHandler(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMEApplicationJSON)

	login := UserLogin{}
	if err := parseAndValidate(c, &login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var usr User
	db.Where("email = ?", login.Email).Take(&usr)

	if !bytes.Equal(usr.Password, hashPassword(login.Password)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	return createAndReturnToken(c, usr.Email)
}

// Creates a JWT token and returns it to the client via JSON. Any handler using
// this function should return immediately after this function, or at least
// be sure that the JSON response isn't getting overwritten.
func createAndReturnToken(c *fiber.Ctx, email string) error {
	token := createJWT(email)
	c.Response().Header.Add(fiber.HeaderAuthorization, token)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

// Hashes the given password with Argon2id.
func hashPassword(password string) []byte {
	return argon2.IDKey(
		[]byte(password),
		[]byte{0xDE, 0xAD, 0xBE, 0xEF},
		1,
		19*1024,
		4,
		32,
	)
}

// Reports if the user is already registered.
func isUserRegistered(email string) bool {
	var usr User
	return db.Where("email = ?", email).Take(&usr).RowsAffected != 0
}
