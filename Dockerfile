# Этап сборки
FROM golang:1.22-alpine AS builder

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache git

# Создаём рабочую директорию
WORKDIR /build

# Копируем go mod файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходники
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o infra-lint ./cmd

# Финальный образ
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем собранный бинарник из builder
COPY --from=builder /build/infra-lint .

# Делаем бинарник исполняемым
RUN chmod +x infra-lint

# Создаём не-root пользователя
RUN adduser -D -u 1000 linter
USER linter

# Точка входа
ENTRYPOINT ["./infra-lint"] 