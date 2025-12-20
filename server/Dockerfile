# Stage 1: Build the Go application
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files and download dependencies first
# This allows Docker to cache the dependencies layer
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for a static binary
# GOOS=linux is to ensure the binary is built for linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/backend ./cmd/api/main.go

# Stage 2: Create a minimal final image
FROM alpine:latest

# It's good practice to run as a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/bin/backend .

# Expose the port the app runs on
EXPOSE 8000

# Set the command to run the application
CMD ["./backend"]
