# Используем официальный образ Go
FROM golang:1.21-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod (go.sum может отсутствовать)
COPY go.mod ./

# Загружаем зависимости (создаст go.sum если нужно)
RUN go mod download && go mod tidy

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Создаем пользователя для безопасности
RUN adduser -D -s /bin/sh appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение
COPY --from=builder /app/main .

# Копируем статические файлы
COPY --from=builder /app/index.html .
COPY --from=builder /app/assets ./assets

# Создаем необходимые папки
RUN mkdir -p data uploads/avatars uploads/achievements

# Устанавливаем права доступа
RUN chown -R appuser:appuser /app

# Переключаемся на пользователя appuser
USER appuser

# Открываем порт
EXPOSE 5000

# Запускаем приложение
CMD ["./main"]
