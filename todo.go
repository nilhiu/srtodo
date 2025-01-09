package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type (
	// GORM model for the todos.
	Todo struct {
		gorm.Model
		Title       string
		Description string
		UserID      uint
	}

	// Used for parsing JSON for todo creation.
	TodoCreate struct {
		Title       string `json:"title"       validate:"required,min=1"`
		Description string `json:"description" validate:"required"`
	}

	// Used for encoding todos in JSON.
	TodoResponse struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
)

func TodosGetHandler(c *fiber.Ctx) error {
	uID, err := checkAuthForUserID(c)
	if err != nil {
		return err
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	var todos []Todo
	db.Limit(limit).Offset((page-1)*limit).Where("user_id = ?", uID).Find(&todos)

	respTodos := make([]TodoResponse, 0, len(todos))
	for _, todo := range todos {
		respTodos = append(respTodos, TodoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":  respTodos,
		"page":  page,
		"limit": limit,
		"total": len(todos),
	})
}

func TodoCreateHandler(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMEApplicationJSON)

	uID, err := checkAuthForUserID(c)
	if err != nil {
		return err
	}

	create := TodoCreate{}
	if err := parseAndValidate(c, &create); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	todo := Todo{
		Title:       create.Title,
		Description: create.Description,
		UserID:      uID,
	}
	db.Create(&todo)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":          todo.ID,
		"title":       todo.Title,
		"description": todo.Description,
	})
}

func TodoUpdateHandler(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMEApplicationJSON)

	uID, err := checkAuthForUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	tID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "API endpoint wasn't given the id parameter",
		})
	}

	var todo Todo
	db.Where("id = ?", tID).Take(&todo)

	if uID != todo.UserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	create := TodoCreate{}
	if err := parseAndValidate(c, &create); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	todo.Title = create.Title
	todo.Description = create.Description
	db.Save(&todo)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":          todo.ID,
		"title":       todo.Title,
		"description": todo.Description,
	})
}

func TodoDeleteHandler(c *fiber.Ctx) error {
	uID, err := checkAuthForUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	tID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "API endpoint wasn't given the id parameter",
		})
	}

	var todo Todo
	db.Where("id = ?", tID).Take(&todo)

	if uID != todo.UserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	db.Delete(&todo)

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"id": todo.ID,
	})
}
