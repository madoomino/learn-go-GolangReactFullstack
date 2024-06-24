package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

func main() {
	app := fiber.New()
	todos := []Todo{}

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"todos": &todos})
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Name == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"msg": "Empty todo name"})
		}

		todo.Id = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(http.StatusCreated).JSON(todo)
	})

	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos[i].Completed = !todos[i].Completed
				return c.Status(http.StatusOK).JSON(todos[i].Completed)
			}
		}
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"msg": "Todo not found."})
	})

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(http.StatusOK).JSON("todo deleted")
			}
		}
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"msg": "Todo not found."})
	})

	app.Get("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				return c.Status(http.StatusOK).JSON(todos[i])
			}
		}
		return c.Status(http.StatusNotFound).JSON("todo not found.")
	})

	log.Fatal(app.Listen(":3000"))
}
