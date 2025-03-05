# Stage 1: Download dependencies
FROM golang:1.23.0-bookworm AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Stage 2: Build the application
FROM golang:1.23.0-bookworm AS builder
WORKDIR /app
COPY --from=deps /go/pkg /go/pkg
COPY . .
# Build the binary and include templates/static files
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o notes ./cmd/api

# Stage 3: Runtime image
FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update && \
  apt-get install -y ca-certificates && \
  rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/notes .
COPY --from=builder /app/.env .

EXPOSE 8025
CMD ["./notes"]