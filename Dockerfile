# Stage 1: Build the Go application
FROM golang:1.20 as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache for dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application files
COPY . .

# Compile the Go application for Linux, disabling CGO for static linking
RUN CGO_ENABLED=0 GOOS=linux go build -o env-output-setter main.go

# Stage 2: Create a smaller runtime image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/env-output-setter .

# Command to run the executable
ENTRYPOINT ["./env-output-setter"]
