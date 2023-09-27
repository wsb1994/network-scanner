# Use an official Golang runtime as the base image
FROM golang:1.17

# Set the working directory inside the container to /app
WORKDIR /app

# Copy all files from the current directory to the container's /app directory
COPY . .

# Build the Go application
RUN go build -o myapp

# Command to run the application
CMD ["./myapp"]
