package main

import (
	"fmt"
	"os"

	"github.com/DeepAung/deep-art/config"
	"github.com/DeepAung/deep-art/pkg/databases"
	"github.com/gofiber/fiber/v2"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env.prod"
	}
	return os.Args[1]
}

func main() {
	app := fiber.New()

	cfg := config.LoadConfig(envPath())
	fmt.Println(cfg.App())
	fmt.Println(cfg.Db())
	fmt.Println(cfg.Jwt())

	db := databases.ConnectDb(cfg.Db())
	defer db.Close()

	fmt.Println(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello")
	})

	app.Listen(":3000")
}
