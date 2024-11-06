# Stage 1: Build
FROM golang:latest AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache for dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application files
COPY . .

# Compile the Go application
RUN go build -o /env-output-setter main.go

# Stage 2: Create a smaller final image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled Go binary from the builder image
COPY --from=builder /env-output-setter .

# Command to run the executable
ENTRYPOINT ["./env-output-setter"]