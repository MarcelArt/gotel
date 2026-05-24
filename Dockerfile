# --- Build Stage ---
FROM golang:1.26.3-alpine AS builder

# Set the working directory inside the container
WORKDIR /src

# Install git, ca-certificates, tzdata, and the swag CLI tool matching the project dependency
RUN apk add --no-cache git ca-certificates tzdata && \
    go install github.com/swaggo/swag/cmd/swag@v1.16.6

# Copy go.mod and go.sum files to cache dependency installation
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Generate Swagger documentation before compiling
RUN swag init --parseDependency --parseInternal

# Build a static binary with optimization flags (disable CGO, strip debugging symbols)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/gotel main.go

# --- Final Runtime Stage ---
FROM alpine:3.21

# Install runtime dependencies (ca-certificates for SSL/TLS, tzdata for timezone info)
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/gotel /app/gotel

# Copy assets and documentation that are required at runtime
COPY --from=builder /src/docs /app/docs
COPY --from=builder /src/public /app/public

# Expose the application port (7001 by default, matching .env setting)
EXPOSE 7001

# Set the binary as the entrypoint
ENTRYPOINT ["/app/gotel"]

# Serve the application by default (can be overridden to "migrate" to run migrations)
CMD ["serve"]
