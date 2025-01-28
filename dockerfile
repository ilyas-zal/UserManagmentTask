# Этап сборки
FROM golang:alpine AS builder

# Устанавливаем рабочий каталог
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod tidy
RUN go mod download

# Копируем остальные файлы
COPY . .

# Собираем приложение
RUN go build -o main cmd/UserManagmentTask/main.go

# Финальный образ
FROM alpine:latest

# Устанавливаем рабочий каталог
WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/main .

# Указываем команду для запуска приложения
CMD ["./main"]
