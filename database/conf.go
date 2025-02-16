package database

import (
	"errors"

	"github.com/nzrsh/irr-ca/models"
)

func CreateConf(conf models.Conf) (*models.Conf, error) {
	if err := DB.Create(&conf).Error; err != nil {
		return nil, errors.New("failed to create conference: " + err.Error())
	}

	return &conf, nil
}

// GetConferencesWithPagination выполняет запрос к базе данных с пагинацией
func GetConferencesWithPagination(page, pageSize int) ([]models.Conf, int64, error) {
	var conferences []models.Conf
	var totalRecords int64

	query := DB.Model(&models.Conf{})

	// Подсчет общего количества записей
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Пагинация
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Выполнение запроса
	if err := query.Find(&conferences).Error; err != nil {
		return nil, 0, err
	}

	return conferences, totalRecords, nil
}

// UpdateConfById обновляет конференцию по ID
func UpdateConfById(id uint, conf models.Conf) (*models.Conf, error) {
	// Находим конференцию по ID
	var existingConf models.Conf
	if err := DB.First(&existingConf, id).Error; err != nil {
		return nil, errors.New("conference not found")
	}

	// Обновляем поля конференции
	existingConf.EventName = conf.EventName
	existingConf.FullName = conf.FullName
	existingConf.Email = conf.Email
	existingConf.Phone = conf.Phone
	existingConf.StartDate = conf.StartDate
	existingConf.StartTime = conf.StartTime
	existingConf.EndDate = conf.EndDate
	existingConf.EndTime = conf.EndTime
	existingConf.Corps = conf.Corps
	existingConf.Location = conf.Location
	existingConf.Platform = conf.Platform
	existingConf.Devices = conf.Devices
	existingConf.Status = conf.Status
	existingConf.URL = conf.URL
	existingConf.User = conf.User
	existingConf.Commentary = conf.Commentary

	// Сохраняем изменения
	if err := DB.Save(&existingConf).Error; err != nil {
		return nil, errors.New("failed to update conference")
	}

	return &existingConf, nil
}

// DeleteConfById удаляет конференцию по ID
func DeleteConfById(id uint) error {
	// Находим конференцию по ID
	var conf models.Conf
	if err := DB.First(&conf, id).Error; err != nil {
		return errors.New("конференция не найдена")
	}

	// Удаляем конференцию
	if err := DB.Delete(&conf).Error; err != nil {
		return errors.New("не удалось удалить конференцию")
	}

	return nil
}
