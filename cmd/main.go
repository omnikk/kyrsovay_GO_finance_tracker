package main

import (
	"finance-tracker/internal/handlers"
	"finance-tracker/internal/middleware"
	"finance-tracker/internal/repository"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	repository.InitDB()

	r := gin.Default()
	r.StaticFile("/", "./static/index.html")
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		api.GET("/categories", handlers.GetCategoriesTree)
		api.POST("/categories", handlers.CreateCategory)
		api.DELETE("/categories/:id", handlers.DeleteCategory)

		api.POST("/transactions", handlers.CreateTransaction)
		api.GET("/transactions/export", handlers.ExportCSV)
		api.GET("/transactions", handlers.GetTransactions)
		api.PUT("/transactions/:id", handlers.UpdateTransaction)
		api.DELETE("/transactions/:id", handlers.DeleteTransaction)

		api.POST("/budgets", handlers.SetBudget)
		api.POST("/transactions/import", handlers.ImportCSV)
		api.GET("/analytics/summary", handlers.GetSummary)
		api.GET("/analytics/budgets", handlers.GetBudgetStatus)
		api.GET("/analytics/monthly", handlers.GetMonthlyStats)
	}

	log.Println("Сервер запущен на http://localhost:8080")
	r.Run(":8080")
}
