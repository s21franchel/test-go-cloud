# Сборка приложения
FROM golang:1.22.3-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /loadbalancer ./cmd/loadbalancer/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app
COPY --from=builder /loadbalancer /app/loadbalancer
COPY configs/ ./configs/

EXPOSE 8080

CMD ["/app/loadbalancer", "-config", "/app/configs/config.yaml"]