package repository

import (
	"finance-tracker/internal/models"
)

func CreateUser(user *models.User) error {
	return DB.Create(user).Error
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := DB.First(&user, id).Error
	return &user, err
}
