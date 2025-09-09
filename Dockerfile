# Стадия сборки
FROM golang:1.22 AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# Финальный образ
FROM alpine:3.20

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/server .

COPY .env .

EXPOSE 8080

CMD ["./server"]