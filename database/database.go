package database

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/nzrsh/irr-ca/config"
	"github.com/nzrsh/irr-ca/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open(config.DBName), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	DB = db

	// Авто-миграция моделей
	err = DB.AutoMigrate(&models.User{}, &models.Conf{}, &models.Lecture{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	// Создание администратора, если он не существует
	createAdminUser()
}

// createAdminUser создает администратора, если его нет
func createAdminUser() {
	var existingUser models.User
	result := DB.Where("username = ?", config.AdminUsername).First(&existingUser)

	if result.Error == gorm.ErrRecordNotFound {
		// Хешируем пароль

		admin := models.User{
			Username: config.AdminUsername,
			Password: config.AdminPassword,
			Role:     "admin",
			Corps:    "Первый",
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Fatalf("Ошибка создания администратора:", err)
		}

		log.Infof("Администратор успешно создан:", config.AdminUsername)
	} else if result.Error != nil {
		log.Fatal("Ошибка поиска администратора:", result.Error)
	}
}
