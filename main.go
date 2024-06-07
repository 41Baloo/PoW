package main

import (
	"flag"
	"fmt"
	"uam/server"

	"github.com/gofiber/fiber/v2"
)

var (
	serverPort string
)

func main() {

	flag.IntVar(&server.Difficulty, "difficulty", 5000000, "Pow bruteforce range")
	flag.IntVar(&server.TimeValid, "timeValid", 240, "Amount of seconds a challenge is valid for")
	flag.IntVar(&server.RetriesAllowed, "retries", 10, "How many times a challenge can be requested")
	flag.IntVar(&server.DynamicSaltLength, "saltLength", 30, "The length of the dynamic salt")
	flag.StringVar(&serverPort, "port", ":80", "Public seed to bruteforce with")

	flag.Parse()

	fmt.Printf("[+] Difficulty: %d\n", server.Difficulty)
	fmt.Printf("[+] TimeValid: %d\n", server.TimeValid)
	fmt.Printf("[+] Retries: %d\n", server.RetriesAllowed)
	fmt.Printf("[+] SaltLength: %d\n", server.DynamicSaltLength)

	config := fiber.Config{
		DisableDefaultContentType: true,
	}
	httpServer := fiber.New(config)
	httpServer.Use(func(c *fiber.Ctx) error {
		server.Middleware(c)
		return nil
	})

	go server.ClearCache()

	if err := httpServer.Listen(serverPort); err != nil {
		panic(err)
	}
}
