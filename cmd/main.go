package main

import (
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

	app.Listen(":3000")
}
