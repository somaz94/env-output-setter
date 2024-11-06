# # Use an official Go runtime as a parent image
# FROM golang:1.20

# # Set the working directory
# WORKDIR /app

# # Copy the Go application files
# COPY . .

# # Compile the Go application
# RUN go build -o /env-output-setter main.go

# # Command to run the executable
# ENTRYPOINT ["/env-output-setter"]

# Stage 1: Build the Go application
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire Go application source to the container
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