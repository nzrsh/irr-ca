package database

import (
	"errors"

	"github.com/nzrsh/irr-ca/models"
)

func CreateLec(lec models.Lecture) (*models.Lecture, error) {

	if err := DB.Create(&lec).Error; err != nil {
		return nil, errors.New("failed to create lec: " + err.Error())
	}

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

	return &existingLec, nil
}

func GetLecsByDay(day string) ([]models.Lecture, error) {
	var lecs []models.Lecture

	if err := DB.Where("date = ?", day).Find(&lecs).Error; err != nil {
		return nil, errors.New("failed to get lecs by day: " + err.Error())
	}

	return lecs, nil
}
