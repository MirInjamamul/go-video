# Use the official Go Docker image with version 1.19.7
FROM golang:1.19.7
RUN apt-get update && apt-get install ffmpeg -y
# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Set the entry point command to run the built executable
CMD ["./main"]