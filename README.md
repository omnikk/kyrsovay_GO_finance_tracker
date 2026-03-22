# Finance Tracker

Веб-сервис для учёта и анализа личных финансов, разработанный на языке Go.

## Стек технологий

- **Go** — основной язык разработки
- **Gin** — веб-фреймворк
- **GORM** — ORM для работы с базой данных
- **SQLite** — база данных
- **JWT** — аутентификация
- **HTML/CSS/JS** — фронтенд

## Функциональность

- Регистрация и авторизация пользователей (JWT)
- Добавление, редактирование, удаление транзакций
- Категории доходов и расходов
- Аналитика: общий баланс, доходы, расходы
- Расходы по категориям с визуализацией
- Фильтрация транзакций по типу и датам
- Установка бюджетов по категориям
- Экспорт транзакций в CSV

## Структура проекта
```
finance-tracker/
├── cmd/
│   └── main.go                  # Точка входа, маршруты
├── internal/
│   ├── handlers/                # HTTP-обработчики
│   │   ├── auth.go
│   │   ├── transactions.go
│   │   └── analytics.go
│   ├── models/                  # Модели данных
│   │   └── models.go
│   ├── repository/              # Работа с БД
│   │   ├── db.go
│   │   ├── user.go
│   │   └── transaction.go
│   ├── service/                 # Бизнес-логика
│   │   ├── auth.go
│   │   ├── transaction.go
│   │   └── analytics.go
│   └── middleware/              # JWT middleware
│       └── jwt.go
├── static/
│   └── index.html               # Фронтенд
├── migrations/
│   └── schema.sql
├── go.mod
└── go.sum
```

## Запуск

### Требования
- Go 1.21+

### Установка и запуск
```bash
# Клонировать репозиторий
git clone https://github.com/omnikk/kyrsovay_GO_finance_tracker
cd kyrsovay_GO_finance_tracker

# Установить зависимости
go mod download

# Запустить сервер
go run cmd/main.go
```

Открыть в браузере: http://localhost:8080

## API Endpoints

### Аутентификация
| Метод | Путь | Описание |
|-------|------|----------|
| POST | /api/auth/register | Регистрация |
| POST | /api/auth/login | Вход |

### Транзакции (требуют JWT)
| Метод | Путь | Описание |
|-------|------|----------|
| GET | /api/transactions | Список транзакций |
| POST | /api/transactions | Создать транзакцию |
| PUT | /api/transactions/:id | Обновить транзакцию |
| DELETE | /api/transactions/:id | Удалить транзакцию |
| GET | /api/transactions/export | Экспорт в CSV |

### Аналитика (требуют JWT)
| Метод | Путь | Описание |
|-------|------|----------|
| GET | /api/analytics/summary | Сводка по финансам |
| GET | /api/analytics/budgets | Статус бюджетов |
| GET | /api/categories | Список категорий |
| POST | /api/budgets | Установить бюджет |

## Примеры запросов

### Регистрация
```json
POST /api/auth/register
{
  "email": "user@mail.ru",
  "password": "123456",
  "name": "Иван"
}
```

### Добавить транзакцию
```json
POST /api/transactions
Authorization: Bearer <token>
{
  "category_id": 1,
  "amount": 50000,
  "type": "income",
  "description": "Зарплата",
  "date": "2026-03-01"
}
```