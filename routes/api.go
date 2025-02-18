package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nzrsh/irr-ca/handlers"

	"github.com/nzrsh/irr-ca/middleware"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/login", handlers.Login)
	app.Post("/conf", handlers.CreateConf)

	api := app.Group("/api", middleware.AuthMiddleware)
	{
		// Протестировано
		users := api.Group("/users")
		users.Post("/", handlers.CreateUser)
		users.Get("/", handlers.GetUsers)
		users.Put("/:id", handlers.UpdateUser)
		users.Delete("/:id", handlers.DeleteUser)

		// Протестировано
		confs := api.Group("/confs")
		confs.Get("/", handlers.GetConfs)
		confs.Put("/:id", middleware.RoleMiddleware, handlers.UpdateConf)    // Middleware перед обработчиком
		confs.Delete("/:id", middleware.RoleMiddleware, handlers.DeleteConf) // Middleware перед обработчиком

		lecs := api.Group("/lecs")
		lecs.Post("/", handlers.CreateLec)
		lecs.Get("/:date", handlers.GetLecsByCurrentDate)
		lecs.Put("/:id", handlers.UpdateLec)
		lecs.Delete("/:id", handlers.DeleteLecById)

		//backup := api.Group("/backup")
		//backup.Get("/:date", handlers.MakeBackupToExcel)
	}
}
