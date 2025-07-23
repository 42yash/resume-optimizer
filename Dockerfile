# Multi-stage build for smaller final image
FROM node:18-alpine AS repomix-installer

# Install repomix globally
RUN npm install -g repomix

# Main application stage
FROM golang:1.24-alpine AS builder

# Install system dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage - minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates nodejs npm git

# Install repomix in the final image
RUN npm install -g repomix

# Create non-root user for security
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy templates and other necessary files
COPY --from=builder /app/templates ./templates/
COPY --from=builder /app/prompt.md ./
COPY --from=builder /app/repo-prompt.md ./

# Create static directory if it doesn't exist
RUN mkdir -p static

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# Run the application
CMD ["./main"]
