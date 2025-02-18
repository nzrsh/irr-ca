package handlers

import (
	"github.com/nzrsh/irr-ca/database"
	"github.com/nzrsh/irr-ca/models"

	"github.com/gofiber/fiber/v2"
)

// Проверка прав администратора
func isAdmin(c *fiber.Ctx) bool {
	role := c.Locals("role")
	return role == string("admin")
}

// Создание пользователя (только для admin)
func CreateUser(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Нет доступа"})
	}

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	database.DB.Create(&user)
	return c.JSON(user)
}

// Получение всех пользователей
func GetUsers(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Нет доступа"})
	}

	var users []models.User
	database.DB.Find(&users)
	return c.JSON(users)
}

// Обновление пользователя (только admin)
func UpdateUser(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Нет доступа"})
	}

	id := c.Params("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Пользователь не найден"})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	database.DB.Save(&user)
	return c.JSON(user)
}

// Удаление пользователя (только admin)
func DeleteUser(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Нет доступа"})
	}

	id := c.Params("id")
	database.DB.Delete(&models.User{}, id)
	return c.JSON(fiber.Map{"message": "Пользователь удален"})
}
