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

## 📁 Архитектура

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
- `README.md` - Документация

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

### Быстрый тест всех endpoints

Используйте автоматический скрипт для базовой проверки:

```
./test-api.sh
```


### Ручное тестирование

#### Шаг 1: Проверка health check

```
curl http://localhost:8080/health | jq
```

Ожидается: {"status":"ok"}


#### Шаг 2: Список заданий

```
curl http://localhost:8080/api/v1/tasks | jq
```

Ожидается: массив из 5 заданий


#### Шаг 3: Регистрация пользователя

```
curl -X POST http://localhost:8080/api/v1/auth/register
-H "Content-Type: application/json"
-d '{"username":"alice"}' | jq
```

Сохраните полученный user_id (например, 1)


#### Шаг 4: Генерация JWT токена

**Важно**: Docker использует секрет `your-secret-key-change-in-production`

##### Вариант 1: С помощью Go (рекомендуется)

```
cat > gen_token.go << 'EOF'
package main

import (
"fmt"
"time"
"github.com/golang-jwt/jwt/v5"
)

func main() {
claims := jwt.MapClaims{
"user_id": 1, // Ваш user_id
"exp": time.Now().Add(24 * time.Hour).Unix(),
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, _ := token.SignedString([]byte("your-secret-key-change-in-production"))
fmt.Println(tokenString)
}
EOF

go run gen_token.go
```


##### Вариант 2: С помощью jwt.io

1. Откройте https://jwt.io
2. В **PAYLOAD** вставьте:

```
{
"user_id": 1,
"exp": 1767825600
}
```

3. В **VERIFY SIGNATURE** замените `your-256-bit-secret` на:

```
your-secret-key-change-in-production
```

4. Скопируйте токен из поля **Encoded**

##### Установка токена

```
export TOKEN="ВАШ_JWT_ТОКЕН"
echo $TOKEN # Проверьте что установлен
```


#### Шаг 5: Получить статус пользователя

```
curl -H "Authorization: Bearer $TOKEN"
http://localhost:8080/api/v1/users/1/status | jq
```

Ожидается: {"id":1,"username":"alice","balance":0,...}


#### Шаг 6: Выполнить задание

```
curl -X POST -H "Authorization: Bearer $TOKEN"
-H "Content-Type: application/json"
-d '{"task_id":1}'
http://localhost:8080/api/v1/users/1/task/complete | jq
```

Проверьте обновленный баланс

```
curl -H "Authorization: Bearer $TOKEN"
http://localhost:8080/api/v1/users/1/status | jq
```


#### Шаг 7: Leaderboard

```
curl -H "Authorization: Bearer $TOKEN"
"http://localhost:8080/api/v1/users/leaderboard?limit=10" | jq
```


#### Шаг 8: Реферальная программа

Создайте второго пользователя

```
curl -X POST http://localhost:8080/api/v1/auth/register
-H "Content-Type: application/json"
-d '{"username":"bob"}' | jq
```

Сохраните user_id (например, 2)
Сгенерируйте токен для пользователя 2
Измените user_id: 2 в gen_token.go и запустите снова

```
go run gen_token.go
export TOKEN2="НОВЫЙ_ТОКЕН"
```

Установите реферера (user_id 1 приглашает user_id 2)
```
curl -X POST -H "Authorization: Bearer $TOKEN2"
-H "Content-Type: application/json"
-d '{"referrer_id":1}'
http://localhost:8080/api/v1/users/2/referrer | jq
```

Проверьте балансы обоих пользователей
```
curl -H "Authorization: Bearer $TOKEN"
http://localhost:8080/api/v1/users/1/status | jq # +100 бонус
```
```
curl -H "Authorization: Bearer $TOKEN2"
http://localhost:8080/api/v1/users/2/status | jq # +50 бонус
```


### С Postman

#### Импорт коллекции

**Вариант 1: Локальный файл (рекомендуется)**
1. Откройте Postman
2. **File** → **Import**
3. Вкладка **File** → **Choose Files**
4. Выберите `postman/User-Management-API.postman_collection.json`
5. Нажмите **Import**

**Вариант 2: Из GitHub**
1. В Postman: **File** → **Import**
2. Вкладка **Link**
3. Вставьте URL:
```
https://raw.githubusercontent.com/abdullinmm/user-management-api/main/postman/User-Management-API.postman_collection.json
```
4. Нажмите **Continue** → **Import**

#### Настройка переменных

1. Кликните на коллекцию **User Management API** (в левой панели)
2. Перейдите на вкладку **Variables**
3. Установите переменные:
   - `baseUrl` = `http://localhost:8080`
   - `token` = (оставьте пустым, заполним после регистрации)
4. Нажмите **Save** (Ctrl+S)

#### Пошаговое тестирование

##### 1️⃣ Публичные endpoints (не требуют токен)

**Health Check**
- Откройте запрос **Health Check**
- Нажмите **Send**
- Ожидается: `{"status":"ok"}`

**Get All Tasks**
- Раскройте папку **Tasks**
- Откройте **List Active Tasks**
- Нажмите **Send**
- Ожидается: JSON массив с 5 заданиями

**Register User**
- Раскройте папку **Auth**
- Откройте **Register User**
- На вкладке **Body** измените username на уникальный (например, `postman_test_001`)
- Нажмите **Send**
- **Сохраните полученный `id`** (например, 9)

##### 2️⃣ Генерация JWT токена

Выполните в терминале:

```
cd ~/git_abdullinmm/user-management-api

cat > gen_token.go << 'EOF'
package main
import ("fmt";"time";"github.com/golang-jwt/jwt/v5")
func main() {
claims := jwt.MapClaims{"user_id": 9, "exp": time.Now().Add(24 * time.Hour).Unix()}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, _ := token.SignedString([]byte("your-secret-key-change-in-production"))
fmt.Println(tokenString)
}
EOF

go run gen_token.go
```


**Добавьте токен в Postman:**
1. Кликните на коллекцию **User Management API**
2. Вкладка **Variables**
3. В строке `token` → колонка **Current value** → вставьте токен
4. Нажмите **Save** (Ctrl+S)

##### 3️⃣ Защищенные endpoints (требуют JWT токен)

**Get User Status**
- Раскройте папку **Users (Protected)**
- Откройте **Get User Status**
- Измените URL на ваш user_id: `/api/v1/users/9/status`
- Нажмите **Send**
- Ожидается: информация о пользователе с `balance: 0`

**Complete Task**
- Откройте **Complete Task**
- Измените URL: `/api/v1/users/9/task/complete`
- На вкладке **Body** проверьте: `{"task_id": 1}`
- Нажмите **Send**
- Ожидается: `{"message":"task completed successfully"}`

**Проверка баланса**
- Снова выполните **Get User Status**
- Ожидается: `balance: 100` (было 0)

**Get Leaderboard**
- Откройте **Get Leaderboard**
- Нажмите **Send**
- Ожидается: список пользователей отсортированных по балансу, ваш пользователь в топе

#### Структура коллекции

- **User Management API/** - Главная коллекция
  - **Auth/** - Аутентификация
    - POST Register User
    - GET Generate Token (Manual) - справочная информация
  - **Tasks/** - Задания
    - GET List Active Tasks
  - **Users (Protected)/** - Защищенные endpoints (требуют JWT)
    - GET Get User Status
    - POST Complete Task
    - POST Set Referrer
    - GET Get Leaderboard
  - **Health Check** - Проверка работы API

#### Troubleshooting

**Ошибка: "Invalid character in header content"**
- Убедитесь что токен добавлен в **переменные коллекции** (Variables)
- Проверьте что на вкладке Authorization тип = **Bearer Token** и значение = `{{token}}`

**Ошибка: "missing auth token" или "invalid token"**
- Сгенерируйте новый токен с правильным user_id
- Убедитесь что секрет = `your-secret-key-change-in-production`
- Проверьте что токен сохранен в переменных и переменная называется `token`

**Ошибка: "access denied"**
- Убедитесь что токен сгенерирован для того же user_id, который используется в URL

### Проверка ошибок

Попытка без токена

```
curl http://localhost:8080/api/v1/users/1/status | jq
```

{"error":"missing auth token"}

Невалидный токен
```
curl -H "Authorization: Bearer invalid_token"
http://localhost:8080/api/v1/users/1/status | jq
```

{"error":"invalid token"}

Повторное выполнение задания
```
curl -X POST -H "Authorization: Bearer $TOKEN"
-H "Content-Type: application/json"
-d '{"task_id":1}'
http://localhost:8080/api/v1/users/1/task/complete | jq
```

{"error":"task already completed"}

undefined

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
