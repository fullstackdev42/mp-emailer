# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git nodejs npm

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the web assets
RUN cd web && npm install && npm run build

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mp-emailer

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and required files from builder
COPY --from=builder /app/mp-emailer .
COPY --from=builder /app/web/templates ./web/templates
COPY --from=builder /app/web/public ./web/public

# Create non-root user
RUN adduser -D appuser
USER appuser

# Expose port (default to 8080 if not set)
ENV APP_PORT=8080
EXPOSE ${APP_PORT}

# Run the application
CMD ["./mp-emailer"] 