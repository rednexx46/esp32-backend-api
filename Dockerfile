# Stage 1 – Build
FROM golang:1.24-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/main.go

# Run tests
RUN go test ./tests/... -v

# Build binary
RUN go build -o esp32-backend-api ./cmd/main.go

# Stage 2 – Runtime
FROM alpine:latest

# Install ca-certificates for HTTPS and create app user
RUN apk --no-cache add ca-certificates && adduser -D -g '' appuser

WORKDIR /home/appuser

# Copy the compiled binary and other necessary files from the builder stage
COPY --from=builder /app/esp32-backend-api .
COPY .env .
COPY docs ./docs/

# Set permissions
USER appuser

# Expose application port
EXPOSE 8080

# Start the server
ENTRYPOINT ["./esp32-backend-api"]