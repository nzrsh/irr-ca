package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nzrsh/irr-ca/config"
	"github.com/nzrsh/irr-ca/database"
	"github.com/nzrsh/irr-ca/middleware"
	"github.com/nzrsh/irr-ca/routes"
)

func main() {
	config.InitConfig()
	database.InitDB()

	app := fiber.New()

	middleware.SetupMiddlewares(app)
	routes.SetupRoutes(app)

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal("Сервер не смог запустится", err)
	}
}
