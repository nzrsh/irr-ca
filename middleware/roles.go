package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func RoleMiddleware(c *fiber.Ctx) error {

	role := c.Locals("role")
	log.Warn(role)
	if role == "viewer" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Нет доступа"})
	} else {
		return c.Next()
	}

}
