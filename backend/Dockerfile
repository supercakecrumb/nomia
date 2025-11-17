# Build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

# Copy spec-examples for fixture mode (from build context root)
COPY spec-examples/ ../spec-examples/

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./server"]