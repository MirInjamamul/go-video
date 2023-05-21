# Use the official Golang base image
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port the application listens on
EXPOSE 8090

# Add a volume mount for the uploads directory
VOLUME ["/app/uploads"]

# Run the application
CMD ["./main"]