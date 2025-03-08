# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o pdf64 ./cmd/main.go

# Final stage
FROM alpine:3.21

# Install runtime dependencies
RUN apk add --no-cache imagemagick ghostscript qpdf

# Create a non-root user
RUN addgroup -S app && adduser -S app -G app

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/pdf64 .

# Use non-root user
USER app

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080

# Run the application
CMD ["./pdf64"]
