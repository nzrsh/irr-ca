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
	if err := utils.Validate.Struct(lec); err != nil {
		log.Warnf("Ошибка валидации на создание лекции: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса", "details": err.Error()})
	}

	// Проверка, что позиция в группе не занята
	var existingLec models.Lecture
	if err := database.DB.Where("date = ? AND group_type = ? AND position_in_group = ?", lec.Date, lec.GroupType, lec.PositionInGroup).First(&existingLec).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Позиция в группе уже занята"})
	}

	// Создание записи в базе данных
	lecture, err := database.CreateLec(lec)
	if err != nil {
		log.Errorf("Ошибка создания лекции: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	// Сокращение ссылки (если URL указан)
	if lec.URL != "" {
		go func(lecID uint, originalURL string) {
			shortURL, err := utils.GetShortUrl(originalURL)
			if err != nil {
				log.Errorf("Ошибка при сокращении ссылки: %s", err)
				return
			}

			// Обновление записи с сокращённой ссылкой
			if err := database.UpdateLecShortURL(lecID, shortURL); err != nil {
				log.Errorf("Ошибка обновления лекции с сокращённой ссылкой: %s", err)
			}
		}(lecture.ID, lec.URL)
	}

	return c.JSON(lecture)
}

func GetLecsByCurrentDate(c *fiber.Ctx) error {
	currentDate := c.Params("date")

	// Проверка формата даты
	_, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		log.Errorf("Неверный формат даты: %s", currentDate)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса. Используйте формат YYYY-MM-DD.",
		})
	}

	// Получение данных из базы данных
	lecs, err := database.GetLecsByDay(currentDate)
	if err != nil {
		log.Errorf("Ошибка при получении данных из базы данных: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Внутренняя ошибка сервера.",
		})
	}

	// Формирование структуры таблицы
	table := make([]*models.Lecture, 48) // 48 строк
	for _, lec := range lecs {
		index := (lec.GroupType-1)*3 + (lec.PositionInGroup - 1)
		table[index] = &lec
	}

	return c.JSON(table)
}

func UpdateLec(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат ID",
		})
	}

	var lec models.Lecture
	if err := c.BodyParser(&lec); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	// Проверка, что позиция в группе не занята
	var existingLec models.Lecture
	if err := database.DB.Where("date = ? AND group_type = ? AND position_in_group = ? AND id != ?", lec.Date, lec.GroupType, lec.PositionInGroup, id).First(&existingLec).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Позиция в группе уже занята"})
	}

	// Обновление лекции
	updatedLec, err := database.UpdateLecById(id, lec)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Сокращение ссылки (если URL изменился)
	if lec.URL != "" && lec.URL != updatedLec.URL {
		go func(lecID uint, originalURL string) {
			shortURL, err := utils.GetShortUrl(originalURL)
			if err != nil {
				log.Errorf("Ошибка при сокращении ссылки: %s", err)
				return
			}

			if err := database.UpdateLecShortURL(lecID, shortURL); err != nil {
				log.Errorf("Ошибка обновления лекции с сокращённой ссылкой: %s", err)
			}
		}(updatedLec.ID, lec.URL)
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
