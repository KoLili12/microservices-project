# Микросервисная архитектура на Go

Проект демонстрирует работу двух микросервисов с использованием Docker и Docker Compose.

## Структура проекта
```
microservices-project/
├── user-service/          # Сервис управления пользователями
│   ├── main.go
│   ├── Dockerfile
│   └── go.mod
├── order-service/         # Сервис управления заказами
│   ├── main.go
│   ├── Dockerfile
│   └── go.mod
├── docker-compose.yml
└── README.md
```

## Сервисы

### User Service (порт 8080)
- `GET /users` - получить всех пользователей
- `POST /users` - создать пользователя
- `GET /users/{id}` - получить пользователя по ID
- `GET /health` - проверка здоровья

### Order Service (порт 8081)
- `GET /orders` - получить все заказы
- `POST /orders` - создать заказ
- `GET /orders/{id}` - получить заказ по ID
- `GET /health` - проверка здоровья

## Запуск проекта

1. Убедитесь, что установлен Docker и Docker Compose
2. Склонируйте репозиторий
3. Перейдите в директорию проекта
4. Запустите сервисы:
```bash
docker-compose up --build
```

## Тестирование API

### Получить всех пользователей
```bash
curl http://localhost:8080/users
```

### Создать пользователя
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Алексей Смирнов","email":"alexey@example.com"}'
```

### Получить все заказы
```bash
curl http://localhost:8081/orders
```

### Создать заказ
```bash
curl -X POST http://localhost:8081/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"product":"Клавиатура","amount":5000}'
```

## Остановка сервисов
```bash
docker-compose down
```