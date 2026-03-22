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

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Ошибка получения sql.DB: ", err)
	}
	sqlDB.Exec("PRAGMA encoding = 'UTF-8';")

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

	// Корневые категории
	roots := []models.Category{
		{Name: "Повседневные", Type: "expense"},
		{Name: "Квартира", Type: "expense"},
		{Name: "Крупные", Type: "expense"},
		{Name: "Доходы", Type: "income"},
	}
	DB.Create(&roots)

	// Находим ID корневых
	var everyday, apartment, big, income models.Category
	DB.Where("name = ? AND type = ?", "Повседневные", "expense").First(&everyday)
	DB.Where("name = ? AND type = ?", "Квартира", "expense").First(&apartment)
	DB.Where("name = ? AND type = ?", "Крупные", "expense").First(&big)
	DB.Where("name = ? AND type = ?", "Доходы", "income").First(&income)

	subs := []models.Category{
		// Повседневные
		{Name: "Еда вне дома", Type: "expense", ParentID: &everyday.ID},
		{Name: "Продукты", Type: "expense", ParentID: &everyday.ID},
		{Name: "Бары и рестораны", Type: "expense", ParentID: &everyday.ID},
		{Name: "Транспорт", Type: "expense", ParentID: &everyday.ID},
		{Name: "Алкоголь", Type: "expense", ParentID: &everyday.ID},
		{Name: "Подарки", Type: "expense", ParentID: &everyday.ID},
		{Name: "Здоровье", Type: "expense", ParentID: &everyday.ID},
		{Name: "Одежда", Type: "expense", ParentID: &everyday.ID},
		{Name: "Развлечения", Type: "expense", ParentID: &everyday.ID},
		{Name: "Связь", Type: "expense", ParentID: &everyday.ID},
		{Name: "Лекарства", Type: "expense", ParentID: &everyday.ID},
		{Name: "Прочее", Type: "expense", ParentID: &everyday.ID},
		// Квартира
		{Name: "Ремонт", Type: "expense", ParentID: &apartment.ID},
		{Name: "Ипотека", Type: "expense", ParentID: &apartment.ID},
		{Name: "ЖКХ", Type: "expense", ParentID: &apartment.ID},
		{Name: "Все для дома", Type: "expense", ParentID: &apartment.ID},
		{Name: "Прочее", Type: "expense", ParentID: &apartment.ID},
		// Крупные
		{Name: "Путешествия", Type: "expense", ParentID: &big.ID},
		{Name: "Одежда", Type: "expense", ParentID: &big.ID},
		{Name: "Гаджеты", Type: "expense", ParentID: &big.ID},
		{Name: "Праздники", Type: "expense", ParentID: &big.ID},
		{Name: "Красота и здоровье", Type: "expense", ParentID: &big.ID},
		{Name: "Образование", Type: "expense", ParentID: &big.ID},
		// Доходы
		{Name: "Зарплата", Type: "income", ParentID: &income.ID},
		{Name: "Премия", Type: "income", ParentID: &income.ID},
		{Name: "Фриланс", Type: "income", ParentID: &income.ID},
		{Name: "Кэшбэк", Type: "income", ParentID: &income.ID},
		{Name: "Депозит", Type: "income", ParentID: &income.ID},
		{Name: "Подарки", Type: "income", ParentID: &income.ID},
		{Name: "Больничный", Type: "income", ParentID: &income.ID},
		{Name: "Проценты", Type: "income", ParentID: &income.ID},
		{Name: "Прочее", Type: "income", ParentID: &income.ID},
	}
	DB.Create(&subs)
}
