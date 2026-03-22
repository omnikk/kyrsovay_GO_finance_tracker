package repository

import (
	"finance-tracker/internal/models"
	"time"
)

func CreateTransaction(t *models.Transaction) error {
	return DB.Create(t).Error
}

func GetTransactionsByUser(userID uint, txType string, from, to *time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := DB.Preload("Category").Where("user_id = ?", userID)

	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if from != nil {
		query = query.Where("date >= ?", from)
	}
	if to != nil {
		query = query.Where("date <= ?", to)
	}

	err := query.Order("date desc").Find(&transactions).Error
	return transactions, err
}

func GetTransactionByID(id, userID uint) (*models.Transaction, error) {
	var t models.Transaction
	err := DB.Preload("Category").Where("id = ? AND user_id = ?", id, userID).First(&t).Error
	return &t, err
}

func UpdateTransaction(t *models.Transaction) error {
	return DB.Save(t).Error
}

func DeleteTransaction(id, userID uint) error {
	return DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Transaction{}).Error
}

func GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := DB.Find(&categories).Error
	return categories, err
}

func GetBudgetsByUser(userID uint, month, year int) ([]models.Budget, error) {
	var budgets []models.Budget
	err := DB.Preload("Category").Where("user_id = ? AND month = ? AND year = ?", userID, month, year).Find(&budgets).Error
	return budgets, err
}

func CreateOrUpdateBudget(b *models.Budget) error {
	var existing models.Budget
	err := DB.Where("user_id = ? AND category_id = ? AND month = ? AND year = ?",
		b.UserID, b.CategoryID, b.Month, b.Year).First(&existing).Error
	if err == nil {
		existing.Limit = b.Limit
		return DB.Save(&existing).Error
	}
	return DB.Create(b).Error
}
