package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupMiddlewares(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders: "Content-Type, Authorization, Accept, Accept-Language",
	}))
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] [${ip}:${port}] ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Local",
	}))
}

func DeleteFileAfterSend(filename string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Вызываем следующий обработчик (основной обработчик запроса)
		if err := c.Next(); err != nil {
			return err
		}

		// Удаляем файл после завершения обработки запроса
		if err := os.Remove(filename); err != nil {
			log.Errorf("Ошибка удаления файла: %s", err)
		}

		return nil
	}
}
