# --- Stage 1: Build Stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy root dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY gateway/ ./gateway/
COPY internal/ ./internal/

# Compile gateway/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gateway-server ./gateway

# --- Stage 2: Minimal Runtime Stage ---
FROM alpine:latest

WORKDIR /app
COPY --from=builder /gateway-server .

EXPOSE 8080
CMD ["/gateway-server"]
