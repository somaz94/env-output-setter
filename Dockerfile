# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o /env-output-setter ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /env-output-setter /env-output-setter

ENTRYPOINT ["/env-output-setter"]