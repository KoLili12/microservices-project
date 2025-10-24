# Микросервисная архитектура на Go с NGINX и Docker Swarm

## Архитектура
```
Client → NGINX (Port 80) → User Service (Port 8080)
                         → Order Service (Port 8081)
```

## Компоненты

### NGINX
- Reverse Proxy
- API Gateway
- Балансировка нагрузки

### User Service
- Управление пользователями
- REST API

### Order Service
- Управление заказами
- Межсервисное взаимодействие с User Service

## Запуск с Docker Compose
```bash
docker-compose up --build
```

Доступ через NGINX:
- http://localhost/api/users
- http://localhost/api/orders

## Запуск в Docker Swarm

### Инициализация Swarm
```bash
docker swarm init
```

### Сборка образов
```bash
docker-compose build
```

### Развёртывание Stack
```bash
docker stack deploy -c docker-stack.yml microservices
```

### Проверка сервисов
```bash
docker stack services microservices
docker stack ps microservices
```

### Масштабирование
```bash
docker service scale microservices_user-service=5
```

### Удаление Stack
```bash
docker stack rm microservices
```

## Тестирование API

### Через NGINX
```bash
# Получить пользователей
curl http://localhost/api/users

# Создать пользователя
curl -X POST http://localhost/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Иван","email":"ivan@test.com"}'

# Получить заказы
curl http://localhost/api/orders

# Создать заказ
curl -X POST http://localhost/api/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"product":"Товар","amount":1000}'
```

## Остановка

### Docker Compose
```bash
docker-compose down
```

### Docker Swarm
```bash
docker stack rm microservices
docker swarm leave --force
```