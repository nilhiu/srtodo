package main

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DBName = "srtodo.db"

var (
	db       *gorm.DB
	validate = validator.New()
)

func main() {
	db = dbSetup()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":       "Speedrunning roadman.sh Todo List API",
			"Description": "seeing how quickly I can learn to develop APIs in Go, having done no front- or back-end development in more than a year.",
		})
	})
	app.Post("/register", UserRegistrationHandler)
	app.Post("/login", UserLoginHandler)
	app.Get("/todos", TodosGetHandler)
	app.Post("/todos", TodoCreateHandler)
	app.Put("/todos/:id", TodoUpdateHandler)
	app.Delete("/todos/:id", TodoDeleteHandler)

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}

// Sets up the sqlite database for use.
func dbSetup() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(DBName), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&Todo{}, &User{}); err != nil {
		panic(err)
	}

	return db
}

// Parses and validates type structs, which are parsed from JSON.
func parseAndValidate(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		return errors.New("request couldn't be parsed")
	}

	return validate.Struct(out)
}
