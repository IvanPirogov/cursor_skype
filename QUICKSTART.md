# Быстрый запуск мессенджера

## Предварительные требования
- Go 1.21+
- Docker и Docker Compose
- PostgreSQL (или использовать Docker)

## Шаги для запуска

### 1. Настройка окружения
```bash
# Скопировать файл переменных окружения
cp .env.example .env

# Создать папку для загрузок
mkdir -p uploads
```

### 2. Запуск базы данных (Docker)
```bash
# Запустить PostgreSQL и Redis
docker-compose up -d postgres redis

# Или запустить все сервисы
docker-compose up -d
```

### 3. Установка зависимостей
```bash
# Скачать Go модули
go mod download
```

### 4. Запуск приложения
```bash
# Собрать приложение
go build -o bin/server cmd/server/main.go

# Запустить сервер
./bin/server

# Или запустить напрямую
go run cmd/server/main.go
```

## Проверка работы

### Health Check
```bash
curl http://localhost:8080/health
```

### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

### Вход в систему
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

## WebSocket подключение

Для подключения к WebSocket используйте:
```
ws://localhost:8080/ws?token=<your-jwt-token>
```

## Использование Makefile

```bash
# Показать все доступные команды
make help

# Запустить Docker контейнеры
make docker-up

# Собрать приложение
make build

# Запустить приложение
make run

# Остановить Docker контейнеры
make docker-down
```

## Структура API

- **Auth**: `/api/v1/auth/`
- **Users**: `/api/v1/users/`
- **Chats**: `/api/v1/chats/`
- **Messages**: `/api/v1/messages/`
- **Contacts**: `/api/v1/contacts/`
- **Calls**: `/api/v1/calls/`
- **WebSocket**: `/ws`

## Решение проблем

### Ошибки подключения к БД
1. Убедитесь что PostgreSQL запущен
2. Проверьте настройки в .env файле
3. Убедитесь что база данных `messenger` существует

### Ошибки сборки
1. Убедитесь что используете Go 1.21+
2. Выполните `go mod tidy`
3. Проверьте что все зависимости установлены

### Порты заняты
Если порт 8080 занят, измените `SERVER_PORT` в .env файле.

## Следующие шаги

1. Реализуйте недостающие обработчики API
2. Добавьте валидацию входных данных
3. Реализуйте загрузку файлов
4. Добавьте тесты
5. Настройте CI/CD