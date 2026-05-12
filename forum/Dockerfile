FROM golang:1.23-alpine AS builder

# Install gcc and musl-dev for CGO (required by sqlite3)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy dependency files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the binary with CGO enabled (sqlite3 requires it)
RUN CGO_ENABLED=1 GOOS=linux go build -o forum ./cmd/main.go

# Stage 2 — Run
FROM alpine:latest

WORKDIR /app

# Install sqlite libs needed at runtime
RUN apk add --no-cache sqlite-libs

# Copy binary from builder
COPY --from=builder /app/forum .

# Copy web assets and pages
COPY web/ web/


# Create data directory for forum.db
RUN mkdir -p internals/db/data

EXPOSE 8080

CMD ["./forum"]
