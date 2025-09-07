# Dockerfile for testing - keeps Go environment for running tests
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Run tests by default
CMD ["go", "test", "./...", "-v"]
