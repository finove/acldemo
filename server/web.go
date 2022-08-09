package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Run() {
	var err error
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	if err = app.Listen(":3000"); err != nil {
		log.Fatal().Err(err).Msg("run web server")
	}
}
