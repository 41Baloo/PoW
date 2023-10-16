package main

import (
	"uam/server"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config := fiber.Config{
		DisableDefaultContentType: true,
	}
	httpServer := fiber.New(config)
	httpServer.Use(func(c *fiber.Ctx) error {
		server.Middleware(c)
		return nil
	})

	go server.ClearCache()

	if err := httpServer.Listen(":80"); err != nil {
		panic(err)
	}
}
