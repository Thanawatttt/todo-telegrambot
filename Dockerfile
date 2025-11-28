# --- Build Stage ---
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# copy go mod
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# build binary
RUN go build -o telegram-todo-bot .

# --- Run Stage ---
FROM alpine:latest

WORKDIR /app

# copy binary จาก build stage
COPY --from=builder /app/telegram-todo-bot .

# default command
CMD ["./telegram-todo-bot"]
