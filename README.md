

PayCore — небольшой backend проект на Go для работы с электронными кошельками.

В проекте можно создавать пользователей и кошельки, пополнять и снимать деньги, переводить средства между кошельками и смотреть историю транзакций.

Проект сделан как практический pet-project, чтобы на практике поработать с Go, REST API, PostgreSQL, Redis и Docker.

## Стек

- Go
- Gin
- PostgreSQL
- GORM
- Redis
- Docker
- REST API
- Git

## Что умеет проект

- создание пользователя
- создание кошелька
- получение кошелька
- пополнение кошелька
- снятие средств
- перевод между кошельками
- получение одной транзакции
- получение истории транзакций
- кэширование кошельков через Redis
- удаление кэша после изменения баланса
- PostgreSQL transactions для операций с балансом
## Структура проекта

```text
PayCore/
├── cmd/
│   └── api/
│       └── main.go
│
├── config/
│   ├── database.go
│   └── redis.go
│
├── internal/
│   ├── user/
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   │
│   ├── wallet/
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   │
│   └── transaction/
│       ├── handler.go
│       ├── model.go
│       ├── repository.go
│       └── service.go
│
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
````

## Архитектура

В проекте используется простое разделение на слои:

```text
Handler
   ↓
Service
   ↓
Repository
   ↓
PostgreSQL

```

## Запуск проекта

### 1. Клонировать проект

```bash
git clone <repository-url>
cd PayCore
```

### 2. Создать `.env`

Создай файл `.env` в корне проекта:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5433
DB_USER=paycore
DB_PASSWORD=paycore
DB_NAME=paycore
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
```

Для Docker API эти значения уже задаются в `docker-compose.yml`.

### 3. Запустить проект

```bash
docker compose up -d --build
```

После этого будут запущены:

```text
PayCore API
PostgreSQL
Redis
```

Проверить контейнеры:

```bash
docker ps
```

### 4. Проверить API

```bash
curl http://localhost:8080/health
```

Ожидаемый ответ:

```json
{
  "app": "paycore",
  "status": "ok"
}
```

### 5. Остановить проект

```bash
docker compose down
```

## Локальный запуск без Docker

Если PostgreSQL и Redis уже запущены локально:

```bash
go run ./cmd/api
```

API будет доступен на:

```text
http://localhost:8080
```

## Пример работы

Обычный сценарий работы:

```text
Создание пользователя
        ↓
Создание кошелька
        ↓
Пополнение
        ↓
Получение баланса
        ↓
Снятие средств
        ↓
Создание второго кошелька
        ↓
Перевод средств
        ↓
История транзакций
```

## Цель проекта

Основная цель проекта — практика backend-разработки на Go и работа с PostgreSQL, Redis, REST API и Docker.

Проект небольшой и специально оставлен без лишней сложности, чтобы сосредоточиться на основных backend-задачах:
работе с API, бизнес-логике, базой данных, транзакциями и кэшированием.

