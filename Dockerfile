# образ go
FROM golang:1.25-alpine

# установка git
RUN apk add --no-cache git

# переход в рабочую директорию контейнера
WORKDIR /app

# копируем и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# копируем весь код проекта в контейнер
COPY . .

# собираем бинарник внутри контейнера
RUN go build -o app ./cmd

# порт
EXPOSE 8080

# окружение по умолчанию
ENV PORT=8080
ENV STORAGE_TYPE=memory

# запуск приложения
CMD ["./app"]
