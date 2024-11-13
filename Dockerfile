# Stage 1: Build the Go application
FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o loykrathong-api ./cmd/main.go

# Stage 2: Create a smaller image to run the Go application
FROM alpine:latest
ENV ENV_PATH="/app/config/.env.dev"

WORKDIR /app
COPY --from=builder /app/loykrathong-api .
COPY config/.env.dev config/.env.dev

EXPOSE 8080
CMD ["./loykrathong-api"]