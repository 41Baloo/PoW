package main

import (
	"uam/server"

	"github.com/gofiber/fiber/v2"
)

func main() {
	httpServer := fiber.New()
	httpServer.Use(func(c *fiber.Ctx) error {
		server.Middleware(c)
		return nil
	})
	if err := httpServer.Listen(":80"); err != nil {
		panic(err)
	}
}
