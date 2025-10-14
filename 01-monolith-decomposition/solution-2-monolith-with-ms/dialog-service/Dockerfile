# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем go.mod/go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарь с cgo
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# --- Финальный stage ---
FROM debian:trixie-slim

WORKDIR /app

# Устанавливаем sqlite для рантайма
RUN apt-get update && apt-get install -y sqlite3 && rm -rf /var/lib/apt/lists/*

# Копируем бинарь
COPY --from=builder /app/main .

# Создаём папку для базы
RUN mkdir -p /app/data

EXPOSE 8081

CMD ["./main"]
