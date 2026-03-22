package handlers

import (
	"finance-tracker/internal/models"
	"finance-tracker/internal/repository"
	"finance-tracker/internal/service"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionInput struct {
	CategoryID  uint    `json:"category_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Type        string  `json:"type" binding:"required,oneof=income expense"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func CreateTransaction(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date := time.Now()
	if input.Date != "" {
		parsed, err := time.Parse("2006-01-02", input.Date)
		if err == nil {
			date = parsed
		}
	}

	t, err := service.CreateTransaction(userID, input.CategoryID, input.Amount, input.Type, input.Description, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func GetTransactions(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	txType := c.Query("type")

	var from, to *time.Time
	if f := c.Query("from"); f != "" {
		t, err := time.Parse("2006-01-02", f)
		if err == nil {
			from = &t
		}
	}
	if t := c.Query("to"); t != "" {
		parsed, err := time.Parse("2006-01-02", t)
		if err == nil {
			to = &parsed
		}
	}

	transactions, err := service.GetTransactions(userID, txType, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

func UpdateTransaction(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id"})
		return
	}

	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date := time.Now()
	if input.Date != "" {
		parsed, err := time.Parse("2006-01-02", input.Date)
		if err == nil {
			date = parsed
		}
	}

	t, err := service.UpdateTransaction(uint(id), userID, input.CategoryID, input.Amount, input.Type, input.Description, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func DeleteTransaction(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id"})
		return
	}

	if err := service.DeleteTransaction(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "транзакция удалена"})
}

func GetCategories(c *gin.Context) {
	categories, err := repository.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func SetBudget(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input struct {
		CategoryID uint    `json:"category_id" binding:"required"`
		Limit      float64 `json:"limit" binding:"required,gt=0"`
		Month      int     `json:"month" binding:"required"`
		Year       int     `json:"year" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	budget := &models.Budget{
		UserID:     userID,
		CategoryID: input.CategoryID,
		Limit:      input.Limit,
		Month:      input.Month,
		Year:       input.Year,
	}

	if err := repository.CreateOrUpdateBudget(budget); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "бюджет установлен"})
}

func ExportCSV(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	transactions, err := service.GetTransactions(userID, "", nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=transactions.csv")

	c.Writer.Write([]byte("\xef\xbb\xbf")) // BOM для корректного открытия в Excel
	c.Writer.WriteString("ID,Тип,Категория,Сумма,Описание,Дата\n")

	for _, t := range transactions {
		line := fmt.Sprintf("%d,%s,%s,%.2f,%s,%s\n",
			t.ID,
			t.Type,
			t.Category.Name,
			t.Amount,
			t.Description,
			t.Date.Format("2006-01-02"),
		)
		c.Writer.WriteString(line)
	}
}
