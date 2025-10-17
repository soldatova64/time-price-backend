FROM golang:1.24.1 as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
# Копируем бинарник из стадии builder
COPY --from=builder /app/app .
# Копируем папку с миграциями
COPY --from=builder /app/migrations ./migrations
# Делаем бинарник исполняемым
RUN chmod +x app
# Копируем .env файл (если нужен)
COPY .env.example .env
EXPOSE 8080
CMD ["./app"]
