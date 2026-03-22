package service

import (
	"finance-tracker/internal/repository"
	"time"
)

type Summary struct {
	TotalIncome  float64            `json:"total_income"`
	TotalExpense float64            `json:"total_expense"`
	Balance      float64            `json:"balance"`
	ByCategory   map[string]float64 `json:"by_category"`
}

type BudgetStatus struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Limit        float64 `json:"limit"`
	Spent        float64 `json:"spent"`
	Remaining    float64 `json:"remaining"`
	Exceeded     bool    `json:"exceeded"`
}

func GetSummary(userID uint, from, to *time.Time) (*Summary, error) {
	transactions, err := repository.GetTransactionsByUser(userID, "", from, to)
	if err != nil {
		return nil, err
	}

	summary := &Summary{
		ByCategory: make(map[string]float64),
	}

	for _, t := range transactions {
		if t.Type == "income" {
			summary.TotalIncome += t.Amount
		} else {
			summary.TotalExpense += t.Amount
			summary.ByCategory[t.Category.Name] += t.Amount
		}
	}
	summary.Balance = summary.TotalIncome - summary.TotalExpense
	return summary, nil
}

func GetBudgetStatus(userID uint, month, year int) ([]BudgetStatus, error) {
	budgets, err := repository.GetBudgetsByUser(userID, month, year)
	if err != nil {
		return nil, err
	}

	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)

	transactions, err := repository.GetTransactionsByUser(userID, "expense", &from, &to)
	if err != nil {
		return nil, err
	}

	spent := make(map[uint]float64)
	for _, t := range transactions {
		spent[t.CategoryID] += t.Amount
	}

	var result []BudgetStatus
	for _, b := range budgets {
		s := spent[b.CategoryID]
		result = append(result, BudgetStatus{
			CategoryID:   b.CategoryID,
			CategoryName: b.Category.Name,
			Limit:        b.Limit,
			Spent:        s,
			Remaining:    b.Limit - s,
			Exceeded:     s > b.Limit,
		})
	}
	return result, nil
}

type MonthlyData struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

func GetMonthlyStats(userID uint, year int) ([]MonthlyData, error) {
	from := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	transactions, err := repository.GetTransactionsByUser(userID, "", &from, &to)
	if err != nil {
		return nil, err
	}

	months := map[int]*MonthlyData{}
	for i := 1; i <= 12; i++ {
		months[i] = &MonthlyData{
			Month: time.Month(i).String(),
		}
	}

	for _, t := range transactions {
		m := int(t.Date.Month())
		if t.Type == "income" {
			months[m].Income += t.Amount
		} else {
			months[m].Expense += t.Amount
		}
	}

	result := make([]MonthlyData, 12)
	for i := 1; i <= 12; i++ {
		result[i-1] = *months[i]
	}
	return result, nil
}
