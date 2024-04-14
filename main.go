package main

import (
	"be-project/pkg/config"
	"be-project/pkg/database"
	"be-project/pkg/logger"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
)

var log logger.Logger = logger.WithPrefix("main")

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	err = database.Init()
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return c.Status(code).JSON(e)
		},
	})

	app.Listen(":3000")

	log.Wrap("start server!").Info()
}
