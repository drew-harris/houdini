package main

import (
	"fmt"
	"os"

	"github.com/drewharris/houdini/browser"
	"github.com/go-rod/rod"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	b := rod.New().MustConnect()
	app := fiber.New()

	app.Get("/in", func(c *fiber.Ctx) error {
		// Detect if true or false
		go browser.ClockIn(b)
		return c.JSON(fiber.Map{"message": "Started"})
	})

	app.Get("/out", func(c *fiber.Ctx) error {
		// Detect if true or false
		go browser.ClockOut(b)
		return c.JSON(fiber.Map{"message": "Started"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := app.Listen(fmt.Sprintf(":%s", port))

	if err != nil {
		panic("Error starting server")
	}
}
