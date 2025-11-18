FROM golang:1.24.6-alpine AS builder

# 1️⃣ Build stage
WORKDIR /app

# Copy root go.mod + go.sum (because build context = repo root)
COPY go.mod go.sum ./
RUN go mod download

# Copy entire project source
COPY . .

# Build the product-service binary
WORKDIR /app/products-service
RUN go build -o main main.go


# 2️⃣ Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy only the built binary
COPY --from=builder /app/products-service/main .

# Copy .env (optional but recommended)
COPY config /app/config
COPY pkg /app/pkg
COPY .env /app/.env


EXPOSE 3000

CMD ["./main"]
