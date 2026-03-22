package handlers

import (
	"encoding/csv"
	"finance-tracker/internal/models"
	"finance-tracker/internal/repository"
	"finance-tracker/internal/service"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func toTitle(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	upper := []rune(strings.ToUpper(string(runes[0])))
	runes[0] = upper[0]
	return string(runes)
}

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
	c.Writer.Write([]byte("\xef\xbb\xbf"))
	c.Writer.WriteString("ID,Тип,Категория,Подкатегория,Сумма,Описание,Дата\n")

	for _, t := range transactions {
		txType := "Расход"
		if t.Type == "income" {
			txType = "Доход"
		}

		categoryName := ""
		subcategoryName := ""

		if t.Category.ParentID != nil {
			subcategoryName = t.Category.Name
			var parent models.Category
			if err := repository.DB.First(&parent, *t.Category.ParentID).Error; err == nil {
				categoryName = parent.Name
			}
		} else {
			categoryName = t.Category.Name
		}

		line := fmt.Sprintf("%d,%s,%s,%s,%.2f,%s,%s\n",
			t.ID,
			txType,
			categoryName,
			subcategoryName,
			t.Amount,
			t.Description,
			t.Date.Format("2006-01-02"),
		)
		c.Writer.WriteString(line)
	}
}

func GetCategoriesTree(c *gin.Context) {
	categories, err := repository.GetCategoriesWithSubs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Type     string `json:"type" binding:"required,oneof=income expense"`
		ParentID *uint  `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat := &models.Category{
		Name:     input.Name,
		Type:     input.Type,
		ParentID: input.ParentID,
		IsCustom: true,
	}
	if err := repository.CreateCategory(cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func DeleteCategory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный id"})
		return
	}
	if err := repository.DeleteCategory(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "категория удалена"})
}

func ImportCSV(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не найден"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка открытия файла"})
		return
	}
	defer f.Close()

	buf := make([]byte, 3)
	f.Read(buf)
	if buf[0] != 0xEF || buf[1] != 0xBB || buf[2] != 0xBF {
		f.Seek(0, 0)
	}

	reader := csv.NewReader(f)
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка чтения заголовка"})
		return
	}

	colIndex := map[string]int{}
	for i, h := range header {
		colIndex[strings.TrimSpace(strings.ToLower(h))] = i
	}

	allCats, _ := repository.GetCategoriesWithSubs()
	catByName := map[string]models.Category{}
	for _, cat := range allCats {
		catByName[strings.ToLower(cat.Name)] = cat
	}

	imported := 0
	skipped := 0
	var errs []string

	rows, _ := reader.ReadAll()
	for i, row := range rows {
		if len(row) == 0 {
			continue
		}

		getCol := func(names ...string) string {
			for _, name := range names {
				if idx, ok := colIndex[name]; ok && idx < len(row) {
					return strings.TrimSpace(row[idx])
				}
			}
			return ""
		}

		rawType := strings.ToLower(getCol("тип", "type"))
		txType := ""
		if strings.Contains(rawType, "доход") || rawType == "income" {
			txType = "income"
		} else if strings.Contains(rawType, "расход") || rawType == "expense" {
			txType = "expense"
		} else {
			skipped++
			errs = append(errs, fmt.Sprintf("Строка %d: неизвестный тип '%s'", i+2, rawType))
			continue
		}

		rawAmount := strings.ReplaceAll(getCol("сумма", "amount"), ",", ".")
		amount, err := strconv.ParseFloat(rawAmount, 64)
		if err != nil || amount <= 0 {
			skipped++
			errs = append(errs, fmt.Sprintf("Строка %d: неверная сумма '%s'", i+2, rawAmount))
			continue
		}

		rawDate := getCol("дата", "date")
		var date time.Time
		for _, layout := range []string{"2006-01-02", "02.01.2006", "01/02/2006", "2006/01/02"} {
			if d, err := time.Parse(layout, rawDate); err == nil {
				date = d
				break
			}
		}
		if date.IsZero() {
			date = time.Now()
		}

		var categoryID uint
		subName := strings.ToLower(getCol("подкатегория", "subcategory"))
		catName := strings.ToLower(getCol("категория", "category"))

		if subName != "" {
			if cat, ok := catByName[subName]; ok {
				categoryID = cat.ID
			}
		}
		if categoryID == 0 && catName != "" {
			if cat, ok := catByName[catName]; ok {
				// Если это корневая категория и подкатегория не указана —
				// ищем подкатегорию с таким же именем внутри неё
				if cat.ParentID == nil {
					for _, sub := range allCats {
						if sub.ParentID != nil && *sub.ParentID == cat.ID &&
							strings.ToLower(sub.Name) == strings.ToLower(cat.Name) {
							categoryID = sub.ID
							break
						}
					}
					// Если такой подкатегории нет — создаём её
					if categoryID == 0 {
						newSub := models.Category{
							Name:     cat.Name,
							Type:     txType,
							ParentID: &cat.ID,
							IsCustom: true,
						}
						repository.DB.Create(&newSub)
						categoryID = newSub.ID
						allCats = append(allCats, newSub)
						catByName[strings.ToLower(newSub.Name)] = newSub
					}
				} else {
					categoryID = cat.ID
				}
			}
		}
		if categoryID == 0 {
			for _, cat := range allCats {
				if strings.ToLower(cat.Name) == "прочее" && cat.Type == txType && cat.ParentID != nil {
					categoryID = cat.ID
					break
				}
			}
		}

		if categoryID == 0 {
			var rootCat models.Category
			if catName != "" {
				rootResult := repository.DB.Where("LOWER(name) = ? AND type = ? AND parent_id IS NULL", catName, txType).First(&rootCat)
				if rootResult.Error != nil {
					rootCat = models.Category{
						Name:     toTitle(catName),
						Type:     txType,
						IsCustom: true,
					}
					repository.DB.Create(&rootCat)
					allCats = append(allCats, rootCat)
					catByName[strings.ToLower(rootCat.Name)] = rootCat
				}
			}

			if subName != "" {
				newCat := models.Category{
					Name:     toTitle(subName),
					Type:     txType,
					IsCustom: true,
				}
				if rootCat.ID > 0 {
					newCat.ParentID = &rootCat.ID
				}
				repository.DB.Create(&newCat)
				categoryID = newCat.ID
				allCats = append(allCats, newCat)
				catByName[strings.ToLower(newCat.Name)] = newCat
			} else if rootCat.ID > 0 {
				categoryID = rootCat.ID
			}
		}

		if categoryID == 0 {
			skipped++
			errs = append(errs, fmt.Sprintf("Строка %d: не удалось определить категорию", i+2))
			continue
		}

		description := getCol("описание", "description")

		t := &models.Transaction{
			UserID:      userID,
			CategoryID:  categoryID,
			Amount:      amount,
			Type:        txType,
			Description: description,
			Date:        date,
		}
		if err := repository.CreateTransaction(t); err == nil {
			imported++
		} else {
			skipped++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"imported": imported,
		"skipped":  skipped,
		"errors":   errs,
	})
}
