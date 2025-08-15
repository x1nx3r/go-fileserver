# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o cdn-server ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/cdn-server .

# Expose the port your app listens on
EXPOSE 8891

# Run the binary
ENTRYPOINT ["./cdn-server"]
