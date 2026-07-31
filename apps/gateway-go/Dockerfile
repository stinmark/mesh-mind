# ---------------------------------------------------
# Stage 1: Build binary
# ---------------------------------------------------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy module files first to leverage Docker layer caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy the rest of the gateway source code
COPY . .

# Build a statically linked Go binary (CGO_ENABLED=0)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gateway ./cmd/server/main.go

# ---------------------------------------------------
# Stage 2: Final Minimal Runtime Image
# ---------------------------------------------------
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates in case your Go service makes HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/gateway .

ENV PORT=8080 \
  INFERENCE_SERVICE_URL="http://inference:8000"

EXPOSE 8080

CMD ["./gateway"]
