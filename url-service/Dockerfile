# Stage 1: Build the Go application
FROM golang:1.22.4-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the necessary files to the container
COPY go.mod go.sum ./

# Build the Go application
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/service ./cmd/api/v1/main.go

FROM alpine:latest as runner

WORKDIR /root/

COPY --from=builder /app/service .
COPY cmd/api/v1/.env .env

EXPOSE 8080

CMD ["./service"]
