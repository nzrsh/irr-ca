package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/nzrsh/irr-ca/database"
	"github.com/nzrsh/irr-ca/models"
	"github.com/nzrsh/irr-ca/utils"
)

func CreateLec(c *fiber.Ctx) error {
	var lec models.Lecture

	if err := c.BodyParser(&lec); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса", "details": err.Error()})
	}

	// Валидация структуры
	if err := utils.Validate.Struct(lec); err != nil { // Исправлено: передаем conf, а не CreateConference
		log.Warnf("Ошибка валидации на создание лекции: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса", "details": err.Error()}) // Добавлено: возвращаем детали ошибки
	}

	// Создание записи в базе данных
	lecture, err := database.CreateLec(lec)
	if err != nil {
		log.Errorf("Ошибка создания лекции: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	return c.JSON(lecture)
}

func GetLecsByCurrentDate(c *fiber.Ctx) error {
	// Получаем дату из параметров запроса
	currentDate := c.Params("date")

	// Проверяем формат даты
	_, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		log.Errorf("Неверный формат даты: %s", currentDate)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса. Используйте формат YYYY-MM-DD.",
		})
	}

	// Получаем данные из базы данных
	lecs, err := database.GetLecsByDay(currentDate)
	if err != nil {
		log.Errorf("Ошибка при получении данных из базы данных: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Внутренняя ошибка сервера.",
		})
	}

	if len(lecs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Лекции не найдены.",
		})
	}
	return c.JSON(lecs)
}

func UpdateLec(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var lec models.Lecture
	if err := c.BodyParser(&lec); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	updatedLec, err := database.UpdateLecById(id, lec)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(updatedLec)
}

func DeleteLecById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := database.DeleteLecById(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
