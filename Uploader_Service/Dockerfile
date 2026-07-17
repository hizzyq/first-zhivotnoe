# Этап 1: Сборка приложения
FROM golang:1.26-alpine AS builder
WORKDIR /app
# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download
# Копируем исходный код и компилируем бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/uploader/main.go

# Этап 2: Финальный легковесный контейнер
FROM alpine:3.19
WORKDIR /app
# Копируем только скомпилированный файл из предыдущего этапа
COPY --from=builder /app/main .
COPY config ./config
EXPOSE 8080
CMD ["./main"]
