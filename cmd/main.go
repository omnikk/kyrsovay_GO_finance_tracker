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
		api.GET("/categories", handlers.GetCategories)

		api.POST("/transactions", handlers.CreateTransaction)
		api.GET("/transactions", handlers.GetTransactions)
		api.PUT("/transactions/:id", handlers.UpdateTransaction)
		api.DELETE("/transactions/:id", handlers.DeleteTransaction)

		api.POST("/budgets", handlers.SetBudget)
		api.GET("/analytics/summary", handlers.GetSummary)
		api.GET("/analytics/budgets", handlers.GetBudgetStatus)
	}

	log.Println("Сервер запущен на http://localhost:8080")
	r.Run(":8080")
}
