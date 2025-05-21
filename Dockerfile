# Stage 1: Build the Go binary
FROM golang:1.24.2 AS builder
WORKDIR /app

# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/main.go

# Stage 2: Create a minimal runtime image
FROM ubuntu:latest

# Install ca-certificates (for HTTPS support)
RUN apt-get update && apt-get install -y \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m -d /home/appuser -s /bin/bash appuser
USER appuser

# Copy the built binary from the builder stage
WORKDIR /home/appuser
COPY --from=builder /app/server ./server
COPY .env .
COPY seed_categories.sh .

RUN chmod +x seed_categories.sh

RUN useradd -m -d /home/appuser -s /bin/bash appuser
USER appuser

# Expose application port
EXPOSE 8080

# Default command
ENTRYPOINT ["sh", "-c", "./seed_categories.sh && exec ./server"]
