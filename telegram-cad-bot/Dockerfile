# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o cadbot cmd/bot/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/cadbot .

# Create data directory
RUN mkdir -p /data

# Set the database path to the volume mount
ENV DB_PATH=/data/cadbot.db

# Run the bot
CMD ["./cadbot", "-db", "/data/cadbot.db"]
