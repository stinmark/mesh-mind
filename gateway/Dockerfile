# STAGE 1: Build binary
FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY gateway/ ./gateway/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gateway-server ./gateway

# STAGE 2: Minimal runtime image
FROM alpine:latest
WORKDIR /app

# Copy binary from builder stage into WORKDIR
COPY --from=builder /gateway-server /app/gateway-server

EXPOSE 8080

CMD ["./gateway-server"]
