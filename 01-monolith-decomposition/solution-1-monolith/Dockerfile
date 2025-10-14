# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with cgo
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Final stage
FROM debian:trixie-slim

# Install sqlite3 runtime
RUN apt-get update && apt-get install -y sqlite3 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Create a directory for the database
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./main"]
