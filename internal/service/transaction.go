package service

import (
	"finance-tracker/internal/models"
	"finance-tracker/internal/repository"
	"time"
)

func CreateTransaction(userID, categoryID uint, amount float64, txType, description string, date time.Time) (*models.Transaction, error) {
	t := &models.Transaction{
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      amount,
		Type:        txType,
		Description: description,
		Date:        date,
	}
	if err := repository.CreateTransaction(t); err != nil {
		return nil, err
	}
	// Загружаем с категорией
	return repository.GetTransactionByID(t.ID, userID)
}

func GetTransactions(userID uint, txType string, from, to *time.Time) ([]models.Transaction, error) {
	return repository.GetTransactionsByUser(userID, txType, from, to)
}

func UpdateTransaction(id, userID, categoryID uint, amount float64, txType, description string, date time.Time) (*models.Transaction, error) {
	t, err := repository.GetTransactionByID(id, userID)
	if err != nil {
		return nil, err
	}
	t.CategoryID = categoryID
	t.Amount = amount
	t.Type = txType
	t.Description = description
	t.Date = date
	err = repository.UpdateTransaction(t)
	return t, err
}

func DeleteTransaction(id, userID uint) error {
	return repository.DeleteTransaction(id, userID)
}
