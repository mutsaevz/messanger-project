# Chat & Messages API

REST API for managing chats and messages.

Реализованы чаты, сообщения, миграции БД, Docker-окружение и тесты.

---

## 🧱 Tech Stack

<p>
  <img alt="Go (net/http)" src="https://img.shields.io/badge/-Go_(net/http)-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="PostgreSQL" src="https://img.shields.io/badge/-PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white" />
  <img alt="GORM" src="https://img.shields.io/badge/-GORM-2D3748?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="Goose (Migrations)" src="https://img.shields.io/badge/-Goose_(Migrations)-4B5563?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="Docker" src="https://img.shields.io/badge/-Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
  <img alt="Docker Compose" src="https://img.shields.io/badge/-Docker_Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
  <img alt="Testify" src="https://img.shields.io/badge/-Testify-25A162?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="httptest" src="https://img.shields.io/badge/-httptest-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
</p>

---

## 📦 Project Structure

```
.
├── cmd/app              # application entrypoint
├── internal
│   ├── config           # database configuration
│   ├── dto              # request / response DTOs
│   ├── http
│   │   ├── handler      # HTTP handlers
│   │   └── router.go    # routes
│   ├── models           # GORM models
│   ├── repository       # DB access layer
│   └── service          # business logic
├── migrations           # goose SQL migrations
├── Dockerfile
├── docker-compose.yaml
└── Makefile
```
---

## 🚀 Run Project

### Requirements
- Docker
- Docker Compose

### Start application

```bash
make up
```
или напрямую:
```
docker compose up --build
```
После запуска:
	•	API доступен по адресу: http://localhost:8080
	•	PostgreSQL запускается в Docker
	•	Миграции применяются автоматически при старте

⸻

🗄 Database Migrations

Миграции выполняются с помощью goose при запуске docker-compose.

SQL-файлы находятся в директории:

migrations/

Каскадное удаление сообщений реализовано на уровне БД (ON DELETE CASCADE).

⸻

📡 API Endpoints

Create chat

### `POST /chats`

Request:
```
{
  "title": "My chat"
}
```
Response:
```
{
  "id": 1,
  "title": "My chat",
  "created_at": "2026-01-29T20:32:44Z"
}
```

⸻

Send message

### `POST /chats/{id}/messages`

Request:
```
{
  "text": "Hello"
}
```
Response:
```
{
  "id": 1,
  "chat_id": 1,
  "text": "Hello",
  "created_at": "2026-01-29T20:33:01Z"
}
```

⸻

Get chat with messages

### `GET /chats/{id}?limit=20`

Response:
```
{
  "id": 1,
  "title": "My chat",
  "created_at": "2026-01-29T20:32:44Z",
  "messages": [
    {
      "id": 2,
      "chat_id": 1,
      "text": "Hello",
      "created_at": "2026-01-29T20:33:01Z"
    }
  ]
}
```
Сообщения возвращаются:
	•	отсортированными по created_at
	•	ограниченными параметром limit (по умолчанию 20, максимум 100)

⸻

Delete chat

### `DELETE /chats/{id}`

Response:
```
204 No Content
```
Все сообщения чата удаляются каскадно.

⸻

❌ Error Response Format

Все ошибки возвращаются в едином формате:
```
{
  "error": "chat not found"
}
```
HTTP-коды используются корректно (400, 404, 500).

⸻

🧪 Tests

В проекте присутствуют unit и handler-тесты.

Запуск тестов:

make test

Используются:
	•	testing
	•	httptest
	•	testify

⸻

✅ Notes
	•	Валидация входных данных выполняется на уровне DTO
	•	Все даты возвращаются в формате ISO 8601 (RFC3339, UTC)
	•	Бизнес-логика отделена от HTTP и DB слоёв
	•	Проект легко расширяется и тестируется

⸻

