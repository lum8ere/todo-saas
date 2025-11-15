# build stage
FROM golang:1.24-alpine AS builder

WORKDIR /src

# копируем go.work и весь backend (быстро будет, если включить build cache позже)
COPY go.work ./
COPY backend ./backend

WORKDIR /src/backend/services/task-service

RUN go build -o /bin/task-service ./cmd/task-service

# runtime stage
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/task-service /app/task-service

EXPOSE 8082

CMD ["/app/task-service"]
