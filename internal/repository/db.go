package repository

import (
	"finance-tracker/internal/models"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("finance.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД: ", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Transaction{},
		&models.Budget{},
	)
	if err != nil {
		log.Fatal("Ошибка миграции: ", err)
	}

	seedCategories()
	log.Println("БД инициализирована успешно")
}

func seedCategories() {
	var count int64
	DB.Model(&models.Category{}).Count(&count)
	if count > 0 {
		return
	}

	categories := []models.Category{
		{Name: "Зарплата", Type: "income"},
		{Name: "Фриланс", Type: "income"},
		{Name: "Инвестиции", Type: "income"},
		{Name: "Еда", Type: "expense"},
		{Name: "Транспорт", Type: "expense"},
		{Name: "Жильё", Type: "expense"},
		{Name: "Развлечения", Type: "expense"},
		{Name: "Здоровье", Type: "expense"},
		{Name: "Одежда", Type: "expense"},
		{Name: "Прочее", Type: "expense"},
	}
	DB.Create(&categories)
}
