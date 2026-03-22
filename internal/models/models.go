package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
	Type string `json:"type" gorm:"not null"` // "income" или "expense"
}

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	CategoryID  uint      `json:"category_id"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // "income" или "expense"
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

type Budget struct {
	ID         uint     `json:"id" gorm:"primaryKey"`
	UserID     uint     `json:"user_id" gorm:"not null"`
	CategoryID uint     `json:"category_id"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
	Limit      float64  `json:"limit" gorm:"not null"`
	Month      int      `json:"month"`
	Year       int      `json:"year"`
}
