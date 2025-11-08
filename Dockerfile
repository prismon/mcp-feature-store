# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build both binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o synthesis-mcp ./cmd/synthesis-mcp
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o synthesis-api ./cmd/synthesis-api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/synthesis-mcp .
COPY --from=builder /build/synthesis-api .

# Copy configuration files
COPY --from=builder /build/db/migrations ./db/migrations

# Create non-root user
RUN addgroup -g 1000 synthesis && \
    adduser -D -u 1000 -G synthesis synthesis && \
    chown -R synthesis:synthesis /app

USER synthesis

# Expose ports
# 8080 - REST API
# 8081 - MCP HTTP
# 8082 - WebSocket
EXPOSE 8080 8081 8082

# Default to running the MCP server
# Can be overridden with docker run command
CMD ["./synthesis-mcp"]
