package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/nzrsh/irr-ca/database"
	"github.com/nzrsh/irr-ca/models"
	"github.com/nzrsh/irr-ca/utils"
)

func CreateConf(c *fiber.Ctx) error {
	var conf models.Conf

	if err := c.BodyParser(&conf); err != nil {
		log.Errorf("Ошибка парсинга конфы: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса"})
	}

	// Валидация структуры
	if err := utils.Validate.Struct(conf); err != nil {
		log.Warnf("Ошибка валидации на создание конфы: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат запроса", "details": err.Error()})
	}

	// Создание записи в базе данных
	conference, err := database.CreateConf(conf)
	if err != nil {
		log.Errorf("Ошибка создания конфы: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	// Асинхронное обновление сокращённой ссылки, если URL предоставлен
	if conf.URL != "" {
		go func(confID uint, originalURL string) {
			shortURL, err := utils.GetShortUrl(originalURL)
			if err != nil {
				log.Errorf("Ошибка при сокращении ссылки: %s", err)
				return
			}

			// Обновление записи в базе данных с сокращённой ссылкой
			if err := database.UpdateConfShortURL(confID, shortURL); err != nil {
				log.Errorf("Ошибка обновления конференции с сокращённой ссылкой: %s", err)
			}
		}(conference.ID, conf.URL)
	}

	return c.JSON(conference)
}

// GetConfs обрабатывает HTTP-запрос для получения конференций с пагинацией
func GetConfs(c *fiber.Ctx) error {
	// Извлечение параметров пагинации
	page, err := strconv.Atoi(c.Query("page", "1")) // По умолчанию: 1
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10")) // По умолчанию: 10
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Вызов функции из database
	conferences, totalRecords, err := database.GetConferencesWithPagination(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch conferences",
		})
	}

	// Возвращаем результат с метаданными пагинации
	return c.JSON(fiber.Map{
		"data": conferences,
		"pagination": fiber.Map{
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (int(totalRecords) + pageSize - 1) / pageSize,
			"totalItems": totalRecords,
		},
	})
}

// UpdateConf обновляет конференцию
func UpdateConf(c *fiber.Ctx) error {
	// Получаем ID из параметров запроса
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	// Парсим JSON с обновленными данными
	var conf models.Conf
	if err := c.BodyParser(&conf); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Получаем текущую конференцию из базы данных
	var existingConf models.Conf
	if err := database.DB.First(&existingConf, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Конференция не найдена",
		})
	}

	// Обновляем конференцию в базе данных (без сокращённой ссылки)
	updatedConf, err := database.UpdateConfById(uint(id), conf)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Если URL изменился, запускаем асинхронное обновление сокращённой ссылки
	if conf.URL != "" && conf.URL != existingConf.URL {
		go func(confID uint, originalURL string) {
			shortURL, err := utils.GetShortUrl(originalURL)
			if err != nil {
				log.Errorf("Ошибка при сокращении ссылки: %s", err)
				return
			}

			// Обновление записи в базе данных с сокращённой ссылкой
			if err := database.UpdateConfShortURL(confID, shortURL); err != nil {
				log.Errorf("Ошибка обновления конференции с сокращённой ссылкой: %s", err)
			}
		}(updatedConf.ID, conf.URL)
	}

	// Возвращаем обновлённую конференцию
	return c.JSON(updatedConf)
}

// DeleteConf удаляет конференцию по ID
func DeleteConf(c *fiber.Ctx) error {
	// Получаем ID из параметров запроса
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID",
		})
	}

	// Удаляем конференцию
	if err := database.DeleteConfById(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Возвращаем успешный статус
	return c.SendStatus(fiber.StatusNoContent)
}
