# Build stage
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Run stage
FROM alpine:3.19
RUN apk add --no-cache sqlite-libs ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/pkg/db/migrations ./pkg/db/migrations
RUN mkdir -p uploads/chat uploads/posts
EXPOSE 8080
CMD ["./server"]
