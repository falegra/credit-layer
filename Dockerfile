# Stage 1: Build
FROM golang:1.26-alpine AS builder

# Install certificates for HTTPS calls (needed by some Go packages)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy dependency files first to leverage layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a static binary
# CGO_ENABLED=0  → no C dependencies, fully static binary
# -ldflags="-w -s" → strip debug info and symbol table (smaller, less info exposed)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/main.go

# Stage 2: Runtime
FROM alpine:3.21

# Install CA certificates (needed for TLS connections to PostgreSQL)
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/server .

# Copy migrations — goose needs them at runtime
COPY --from=builder /app/db/migrations ./db/migrations

# Set ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8080

CMD ["./server"]
