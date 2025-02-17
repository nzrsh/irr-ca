package database

import (
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/nzrsh/irr-ca/models"
)

func CreateLec(lec models.Lecture) (*models.Lecture, error) {

	if err := DB.Create(&lec).Error; err != nil {
		return nil, errors.New("failed to create lec: " + err.Error())
	}
	log.Infof("Лекция успешно создана: %s", lec.Date)
	return &lec, nil
}

func DeleteLecById(id uint64) error {
	var lec models.Lecture
	if err := DB.First(&lec, id).Error; err != nil {
		return errors.New("lec not found")
	}

	if err := DB.Delete(&lec).Error; err != nil {
		return errors.New("failed to delete lec")
	}
	log.Infof("Лекция с ID %d успешно удалена", id)
	return nil
}

func UpdateLecById(id uint64, lec models.Lecture) (*models.Lecture, error) {
	var existingLec models.Lecture
	if err := DB.First(&existingLec, id).Error; err != nil {
		return nil, errors.New("lec not found")
	}

	existingLec.Date = lec.Date
	existingLec.GroupType = lec.GroupType
	existingLec.StartTime = lec.StartTime
	existingLec.EndTime = lec.EndTime
	existingLec.AbnormalTime = lec.AbnormalTime
	existingLec.Platform = lec.Platform
	existingLec.Corps = lec.Corps
	existingLec.Location = lec.Location
	existingLec.Groups = lec.Groups
	existingLec.Lectors = lec.Lectors
	existingLec.URL = lec.URL
	existingLec.StreamKey = lec.StreamKey
	existingLec.Account = lec.Account
	existingLec.Commentary = lec.Commentary

	if err := DB.Save(&existingLec).Error; err != nil {
		return nil, errors.New("failed to update lec")
	}
	log.Infof("Лекция с ID %d успешно обновлена", id)
	return &existingLec, nil
}

func GetLecsByDay(day string) ([]models.Lecture, error) {
	var lecs []models.Lecture

	if err := DB.Where("date = ?", day).Find(&lecs).Error; err != nil {
		return nil, errors.New("failed to get lecs by day: " + err.Error())
	}
	log.Infof("Лекции на день %s успешно получены", day)
	return lecs, nil
}

func UpdateLecShortURL(lecID uint, shortURL string) error {
	// Находим лекцию по ID
	var lec models.Lecture
	result := DB.First(&lec, lecID)
	if result.Error != nil {
		log.Errorf("Ошибка при поиске лекции с ID %d: %s", lecID, result.Error)
		return result.Error
	}

	// Обновляем поле ShortURL
	lec.ShortURL = shortURL

	// Сохраняем изменения в базе данных
	result = DB.Save(&lec)
	if result.Error != nil {
		log.Errorf("Ошибка при обновлении лекции с ID %d: %s", lecID, result.Error)
		return result.Error
	}

	log.Infof("Сокращенная ссылка для лекции с ID %d успешно обновлена", lecID)
	return nil
}
