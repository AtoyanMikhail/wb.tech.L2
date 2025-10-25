# Calendar HTTP Server

HTTP-сервер для работы с календарем событий. Реализует CRUD операции для управления событиями пользователей.

## Структура проекта

```
2.19/
├── cmd/
│   └── server/
│       └── main.go           # Точка входа приложения
├── internal/
│   ├── config/
│   │   └── config.go         # Конфигурация приложения
│   ├── handler/
│   │   └── event_handler.go  # HTTP обработчики
│   ├── middleware/
│   │   └── logger.go         # Middleware для логирования
│   ├── model/
│   │   └── event.go          # Модели данных
│   └── service/
│       ├── event_service.go      # Бизнес-логика
│       └── event_service_test.go # Unit-тесты
└── go.mod
```

## API Endpoints

### POST /create_event
Создание нового события.

**Request Body (JSON):**
```json
{
  "user_id": 1,
  "date": "2024-01-15",
  "event": "Meeting with team"
}
```

**Request Body (form-data):**
```
user_id=1&date=2024-01-15&event=Meeting with team
```

**Response:**
```json
{
  "result": {
    "id": 1,
    "user_id": 1,
    "date": "2024-01-15T00:00:00Z",
    "event": "Meeting with team"
  }
}
```

### POST /update_event
Обновление существующего события.

**Request Body (JSON):**
```json
{
  "id": 1,
  "user_id": 1,
  "date": "2024-01-16",
  "event": "Updated meeting"
}
```

**Response:**
```json
{
  "result": {
    "id": 1,
    "user_id": 1,
    "date": "2024-01-16T00:00:00Z",
    "event": "Updated meeting"
  }
}
```

### POST /delete_event
Удаление события.

**Request Body (JSON):**
```json
{
  "id": 1
}
```

**Response:**
```json
{
  "result": "event deleted successfully"
}
```

### GET /events_for_day
Получение событий на день.

**Query Parameters:**
- `user_id` - ID пользователя
- `date` - дата в формате YYYY-MM-DD

**Example:**
```
GET /events_for_day?user_id=1&date=2024-01-15
```

**Response:**
```json
{
  "result": [
    {
      "id": 1,
      "user_id": 1,
      "date": "2024-01-15T00:00:00Z",
      "event": "Meeting with team"
    }
  ]
}
```

### GET /events_for_week
Получение событий на неделю (7 дней начиная с указанной даты).

**Query Parameters:**
- `user_id` - ID пользователя
- `date` - начальная дата в формате YYYY-MM-DD

**Example:**
```
GET /events_for_week?user_id=1&date=2024-01-15
```

### GET /events_for_month
Получение событий на месяц.

**Query Parameters:**
- `user_id` - ID пользователя
- `date` - любая дата месяца в формате YYYY-MM-DD

**Example:**
```
GET /events_for_month?user_id=1&date=2024-01-15
```

## HTTP Status Codes

- **200 OK** - успешное выполнение запроса
- **400 Bad Request** - ошибка валидации входных данных
- **503 Service Unavailable** - бизнес-логическая ошибка (например, попытка удалить несуществующее событие)
- **500 Internal Server Error** - внутренняя ошибка сервера

## Error Response Format

```json
{
  "error": "описание ошибки"
}
```

## Запуск приложения

### Настройка окружения

Убедитесь, что Go установлен и доступен в PATH:

### Установка зависимостей
```bash
cd 2.18
go mod tidy
```

### Запуск сервера
```bash
# С портом по умолчанию (8080)
go run cmd/server/main.go

# С пользовательским портом
PORT=3000 go run cmd/server/main.go

# Используя Makefile
make run

# С пользовательским портом через Makefile
PORT=3000 make run
```

### Запуск тестов
```bash
go test ./internal/service/... -v
```

### Проверка race conditions
```bash
go test ./internal/service/... -race
```

### Линтинг и проверка кода
```bash
go vet ./...
go fmt ./...
```

## Примеры использования

### Создание события с JSON
```bash
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "date": "2024-01-15",
    "event": "Team meeting"
  }'
```

### Создание события с form-data
```bash
curl -X POST http://localhost:8080/create_event \
  -d "user_id=1&date=2024-01-15&event=Team meeting"
```

### Обновление события
```bash
curl -X POST http://localhost:8080/update_event \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "user_id": 1,
    "date": "2024-01-16",
    "event": "Updated team meeting"
  }'
```

### Удаление события
```bash
curl -X POST http://localhost:8080/delete_event \
  -H "Content-Type: application/json" \
  -d '{"id": 1}'
```

### Получение событий на день
```bash
curl "http://localhost:8080/events_for_day?user_id=1&date=2024-01-15"
```

### Получение событий на неделю
```bash
curl "http://localhost:8080/events_for_week?user_id=1&date=2024-01-15"
```

### Получение событий на месяц
```bash
curl "http://localhost:8080/events_for_month?user_id=1&date=2024-01-15"
```
