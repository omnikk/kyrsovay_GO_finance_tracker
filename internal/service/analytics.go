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
