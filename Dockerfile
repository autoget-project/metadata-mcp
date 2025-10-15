# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for creating a static binary
# -o /app/metadata-mcp specifies the output path and name of the binary
# ./cmd/main.go specifies the entry point of the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/metadata-mcp ./cmd/main.go

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS support
RUN apk add --no-cache ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/metadata-mcp .

# Command to run the application
CMD ["./metadata-mcp"]
