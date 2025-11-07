# User Management API

![Go](https://img.shields.io/badge/Go-1.24-blue?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue?logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Compose-blue?logo=docker)
![Chi](https://img.shields.io/badge/Chi-v5-blue)
![JWT](https://img.shields.io/badge/JWT-v5-green)
![License](https://img.shields.io/badge/license-MIT-green)

Production-ready REST API для управления пользователями с системой выполнения заданий, реферальной программой и начислением поинтов. Реализовано с использованием Clean Architecture на Go.

##  Возможности

-  **Управление пользователями** - Регистрация с уникальными username
-  **JWT аутентификация** - Защита endpoints с Bearer токенами  
-  **Система заданий** - 5 типов заданий с разными наградами
-  **Реферальная программа** - Бонусы для реферера и реферала
-  **Баланс поинтов** - Отслеживание баланса в реальном времени
-  **История транзакций** - Полный аудит всех начислений
-  **Leaderboard** - Топ пользователей по балансу
-  **Clean Architecture** - Разделение на слои (domain, usecase, repository, handler)
-  **PostgreSQL** - Надежное хранение данных с pgx v5
-  **Docker** - Контейнеризация с docker-compose
-  **CI/CD** - GitHub Actions для тестирования и линтинга

##  Содержание

- [Архитектура](#архитектура)
- [Технологии](#технологии)
- [Быстрый старт](#быстрый-старт)
- [API Endpoints](#api-endpoints)
- [Аутентификация](#аутентификация)
- [База данных](#база-данных)
- [Тестирование](#тестирование)
- [Разработка](#разработка)

##  Архитектура

Проект следует принципам **Clean Architecture**:

- `cmd/api/` - Точка входа приложения
- `internal/` - Внутренние пакеты
  - `domain/` - Бизнес-логика
    - `entities/` - Сущности (User, Task, Balance)
    - `interfaces/` - Интерфейсы репозиториев
  - `usecase/` - Use cases (бизнес-правила)
  - `repository/postgresql/` - PostgreSQL драйвер
  - `handler/http/` - REST API handlers
  - `middleware/` - HTTP middleware (JWT auth)
  - `pkg/jwt/` - JWT управление
  - `config/` - Конфигурация
- `migrations/` - SQL миграции
- `postman/` - Postman коллекция
- `docker-compose.yml` - Docker оркестрация


##  Технологии

- **Язык**: Go 1.24
- **База данных**: PostgreSQL 16
- **Driver**: pgx v5 (нативный PostgreSQL драйвер)
- **Router**: Chi v5
- **JWT**: golang-jwt/jwt v5
- **Контейнеризация**: Docker & Docker Compose
- **CI/CD**: GitHub Actions

##  Быстрый старт

### Требования

- Docker & Docker Compose
- Go 1.24+ (для локальной разработки)
- Make (опционально)

### Запуск

Клонировать репозиторий

```
git clone https://github.com/abdullinmm/user-management-api.git
cd user-management-api
```

Запустить с Docker Compose

```
make docker-up
```

или

```
docker-compose up -d
```

Проверить health check

```
curl http://localhost:8080/health
```


**Приложение доступно на** `http://localhost:8080`

##  API Endpoints

### Публичные endpoints

- GET /health # Health check
- GET /api/v1/tasks # Список активных заданий
- POST /api/v1/auth/register # Регистрация пользователя


### Защищенные endpoints (требуется JWT)

- GET /api/v1/users/{id}/status # Статус пользователя
- GET /api/v1/users/leaderboard # Топ пользователей
- POST /api/v1/users/{id}/task/complete # Выполнить задание
- POST /api/v1/users/{id}/referrer # Установить реферера


### Примеры запросов

#### 1. Регистрация пользователя

```
curl -X POST http://localhost:8080/api/v1/auth/register
-H "Content-Type: application/json"
-d '{"username":"alice"}'
```


Ответ:

```
{
"id": 1,
"username": "alice"
}
```

#### 2. Список заданий

```
curl http://localhost:8080/api/v1/tasks
```


Ответ:

```
[
{
"id": 1,
"code": "TASK_REGISTER",
"title": "Complete Registration",
"reward_points": 100,
"is_active": true
}
]
```

#### 3. Получить статус пользователя (требуется JWT)

```
curl -H "Authorization: Bearer YOUR_JWT_TOKEN"
http://localhost:8080/api/v1/users/1/status
```

Ответ:

```
{
"id": 1,
"username": "alice",
"referrer_id": null,
"balance": 100,
"created_at": "2025-11-07T10:00:00Z"
}
```


#### 4. Выполнить задание

```
curl -X POST -H "Authorization: Bearer YOUR_JWT_TOKEN"
-H "Content-Type: application/json"
-d '{"task_id":1}'
http://localhost:8080/api/v1/users/1/task/complete

```


#### 5. Leaderboard

```
curl -H "Authorization: Bearer YOUR_JWT_TOKEN"
"http://localhost:8080/api/v1/users/leaderboard?limit=10"
```


##  Аутентификация

Все защищенные endpoints требуют JWT токен в заголовке:

```
Authorization: Bearer <JWT_TOKEN>
```


### Генерация JWT токена

Для тестирования используйте [jwt.io](https://jwt.io):

**Payload:**

```
{
"user_id": 1,
"exp": 1767825600
}
```


**Secret:** `dev_secret`

**Важно**: В production используйте надежный секрет и короткое время жизни токенов!

##  База данных

### Схема

#### Users

```
id BIGSERIAL PRIMARY KEY
username VARCHAR(255) UNIQUE NOT NULL
referrer_id BIGINT REFERENCES users(id)
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```


#### Tasks

```
id BIGSERIAL PRIMARY KEY
code VARCHAR(100) UNIQUE NOT NULL
title VARCHAR(255) NOT NULL
reward_points BIGINT NOT NULL
is_active BOOLEAN DEFAULT true
```


#### Balances

```
user_id BIGINT PRIMARY KEY REFERENCES users(id)
points BIGINT NOT NULL DEFAULT 0
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```


#### Transactions

```
id BIGSERIAL PRIMARY KEY
user_id BIGINT REFERENCES users(id)
delta BIGINT NOT NULL
reason VARCHAR(255) NOT NULL
reference_type VARCHAR(50)
reference_id BIGINT
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```


#### User Tasks

```
user_id BIGINT REFERENCES users(id)
task_id BIGINT REFERENCES tasks(id)
completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
PRIMARY KEY (user_id, task_id)
```


### Миграции

Применить миграции (автоматически при docker-compose up)

```
psql -U postgres -d user_management -f migrations/001_initial_schema.up.sql
```

Откатить миграции

```
psql -U postgres -d user_management -f migrations/001_initial_schema.down.sql
```


##  Тестирование

### С Postman

1. Импортируйте коллекцию: `postman/User-Management-API.postman_collection.json`
2. Установите переменную `baseUrl` = `http://localhost:8080`
3. Зарегистрируйте пользователя через **Auth > Register User**
4. Сгенерируйте JWT токен на [jwt.io](https://jwt.io)
5. Установите переменную `token` со значением JWT
6. Тестируйте защищенные endpoints

### С bash скриптом

```
./test-api.sh
```


### С curl

Health check

```
curl http://localhost:8080/health
```

Список заданий

```
curl http://localhost:8080/api/v1/tasks | jq
```

Регистрация

```
curl -X POST http://localhost:8080/api/v1/auth/register
-H "Content-Type: application/json"
-d '{"username":"testuser"}' | jq
```

Установите TOKEN после генерации

```
export TOKEN="YOUR_JWT_TOKEN"
```

Статус пользователя

```
curl -H "Authorization: Bearer $TOKEN"
http://localhost:8080/api/v1/users/1/status | jq
```

Выполнить задание

```
curl -X POST -H "Authorization: Bearer $TOKEN"
-H "Content-Type: application/json"
-d '{"task_id":1}'
http://localhost:8080/api/v1/users/1/task/complete | jq
```

Leaderboard

```
curl -H "Authorization: Bearer $TOKEN"
"http://localhost:8080/api/v1/users/leaderboard?limit=5" | jq
```


##  Разработка

### Make команды

- make help # Показать все команды
- make build # Собрать приложение
- make run # Запустить локально
- make test # Запустить тесты
- make lint # Запустить линтер
- make docker-build # Собрать Docker образ
- make docker-up # Запустить сервисы
- make docker-down # Остановить сервисы
- make docker-logs # Показать логи


### Локальная разработка

Установить зависимости

```
go mod download
```

Запустить PostgreSQL

```
docker run -d -p 5432:5432
-e POSTGRES_PASSWORD=postgres
-e POSTGRES_DB=user_management
postgres:16-alpine
```


Применить миграции

```
psql -U postgres -d user_management -f migrations/001_initial_schema.up.sql
```

Запустить приложение

```
go run ./cmd/api
```


### Переменные окружения

```
DATABASE_URL="postgresql://postgres:postgres@localhost:5432/user_management?sslmode=disable"
HTTP_PORT="8080"
JWT_SECRET="dev_secret"
REFERRAL_BONUS="100" # Бонус для реферера
REFEREE_BONUS="50" # Бонус для нового пользователя
```


##  Обработка ошибок

API возвращает понятные сообщения об ошибках:

```
{
"error": "user not found"
}
```

undefined

```
{
"error": "task already completed"
}
```

undefined

```
{
"error": "invalid token"
}
```


HTTP коды ответов:
- `200` - Успех
- `201` - Создано
- `400` - Неверный запрос
- `401` - Не авторизован
- `403` - Доступ запрещен
- `404` - Не найдено
- `500` - Внутренняя ошибка сервера

##  Лицензия

MIT License

##  Автор

**Marsel Abdullin**

- GitHub: [@abdullinmm](https://github.com/abdullinmm)
