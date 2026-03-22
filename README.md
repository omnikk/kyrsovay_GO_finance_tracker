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


## Чистая установка (с нуля на Windows)

### Шаг 1 — Установить Go
1. Зайди на https://go.dev/dl/
2. Скачай установщик `go1.xx.windows-amd64.msi`
3. Запусти, нажимай Next до конца
4. Проверь в терминале:
```powershell
go version
```

### Шаг 2 — Установить Git
1. Зайди на https://git-scm.com/download/win
2. Скачай и установи, все настройки по умолчанию
3. Проверь:
```powershell
git --version
```

### Шаг 3 — Установить GCC (нужен для SQLite)
1. Зайди на https://github.com/niXman/mingw-builds-binaries/releases
2. Скачай архив `x86_64-win32-seh` последней версии
3. Распакуй в `C:\mingw64`
4. Добавь `C:\mingw64\bin` в PATH:
   - Поиск Windows → "Переменные среды" → Path → Изменить → Добавить новую
5. Перезапусти терминал и проверь:
```powershell
gcc --version
```

### Шаг 4 — Клонировать и запустить проект
```powershell
git clone https://github.com/omnikk/kyrsovay_GO_finance_tracker
cd kyrsovay_GO_finance_tracker
go mod download
go run cmd/main.go
```

### Шаг 5 — Открыть в браузере
Перейди по адресу: http://localhost:8080

### Возможная проблема
Если при `go run` появится ошибка `cgo: C compiler "gcc" not found` — GCC не добавлен в PATH.
Перепроверь шаг 3 и обязательно перезапусти терминал после добавления переменной.