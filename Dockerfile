# Use an official Go runtime as a parent image
FROM golang:1.23

# Set the working directory
WORKDIR /app

# Copy the Go application files
COPY . .

# Compile the Go application
RUN go build -o /env-output-setter main.go

# Command to run the executable
ENTRYPOINT ["/env-output-setter"]