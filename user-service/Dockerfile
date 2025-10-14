# Используем официальный образ Go для сборки
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o user-service .

# Используем минимальный образ для запуска
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированное приложение из builder
COPY --from=builder /app/user-service .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./user-service"]