package database

import (
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/nzrsh/irr-ca/models"
)

func CreateLec(lec models.Lecture) (*models.Lecture, error) {

	newlec, err := FillTimesFromRowID(lec)
	if err != nil {
		return nil, err
	}

	if err := DB.Create(&newlec).Error; err != nil {
		return nil, errors.New("failed to create lec: " + err.Error())
	}
	log.Infof("Лекция успешно создана: %s", newlec.Date)
	return &newlec, nil
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

	existingLec, err := FillTimesFromRowID(existingLec)
	if err != nil {
		return nil, err
	}

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

func FillTimesFromRowID(lec models.Lecture) (models.Lecture, error) {
	switch lec.RowID {
	// Первая смена
	case "0", "3", "6", "9", "15", "18", "21":
		lec.StartTime = "9:30"
		lec.EndTime = "10:40"

	case "1", "4", "7", "10", "13", "16", "19", "22":
		lec.StartTime = "10:50"
		lec.EndTime = "12:00"

	case "2", "5", "8", "11", "14", "17", "20", "23":
		lec.StartTime = "12:10"
		lec.EndTime = "13:20"

	// Вторая смена
	case "24", "27", "30", "33", "36", "39", "42", "45":
		lec.StartTime = "13:30"
		lec.EndTime = "14:40"

	case "25", "28", "31", "34", "37", "40", "43", "46":
		lec.StartTime = "14:50"
		lec.EndTime = "16:00"

	case "26", "29", "32", "35", "38", "41", "44", "47":
		lec.StartTime = "16:10"
		lec.EndTime = "17:20"

	default:
		return lec, errors.New("неверная позиция строки")
	}

	return lec, nil
}
