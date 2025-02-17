package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/nzrsh/irr-ca/database"
	"github.com/nzrsh/irr-ca/middleware"
	"github.com/xuri/excelize/v2"
)

func MakeBackupToExcel(c *fiber.Ctx) error {
	date := c.Params("date")

	// Проверяем формат даты
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Errorf("Неверный формат даты: %s", date)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса. Используйте формат YYYY-MM-DD.",
		})
	}

	// Создаем новый Excel-файл
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Errorf("Error closing Excel file: %s", err)
		}
	}()

	// Создать новый лист с именем, равным дате
	sheetName := date
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Errorf("Ошибка создания листа: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Внутренняя ошибка сервера.",
		})
	}

	// Получаем данные из базы данных
	lectures, err := database.GetLecsByDay(sheetName)
	if err != nil || len(lectures) == 0 {
		log.Errorf("Ошибка чтения данных: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Внутренняя ошибка сервера.",
		})
	}

	// Заголовки столбцов
	headers := []string{
		"ID", "Current Date", "Group Type", "Start Time", "End Time", "Abnormal Time",
		"Platform", "Corps", "Location", "Groups", "Lectors", "URL", "Stream Key", "Account", "Commentary",
	}
	for col, header := range headers {
		cell, err := excelize.CoordinatesToCellName(col+1, 1) // Столбцы начинаются с 1, строка 1
		if err != nil {
			log.Errorf("Ошибка преобразования координат: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Внутренняя ошибка сервера.",
			})
		}
		f.SetCellValue(sheetName, cell, header)
	}

	// Записываем данные
	for row, lecture := range lectures {
		// Записываем значения в ячейки
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row+2), lecture.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row+2), lecture.Date)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row+2), lecture.GroupType)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row+2), lecture.StartTime)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row+2), lecture.EndTime)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row+2), lecture.AbnormalTime)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row+2), lecture.Platform)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row+2), lecture.Corps)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row+2), lecture.Location)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row+2), lecture.Groups)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row+2), lecture.Lectors)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row+2), lecture.URL)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row+2), lecture.StreamKey)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", row+2), lecture.Account)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", row+2), lecture.Commentary)
	}

	// Устанавливаем активный лист
	f.SetActiveSheet(index)

	timeNow := time.Now().Format("2006-01-02_15-04-05")

	// Сохраняем файл
	filename := fmt.Sprintf("./backups/Backup_%s_by_%s.xlsx", date, timeNow)
	if err := f.SaveAs(filename); err != nil {
		log.Errorf("Ошибка сохранения файла: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Внутренняя ошибка сервера.",
		})
	}

	// Устанавливаем заголовок для скачивания файла
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=Backup_%s.xlsx", date))

	// Отправляем файл клиенту
	if err := c.SendFile(filename); err != nil {
		log.Errorf("Ошибка отправки файла: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка отправки файла.",
		})
	}

	// Удаляем файл после успешной отправки
	defer func() {
		if err := os.Remove(filename); err != nil {
			log.Errorf("Ошибка удаления файла: %s", err)
		}
	}()

	return middleware.DeleteFileAfterSend(filename)(c)
}
