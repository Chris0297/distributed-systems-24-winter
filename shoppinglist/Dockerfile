# Start from the official Golang image
FROM golang:1.23-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the Go application files
COPY . .

# Build the Go app
RUN go mod download
RUN go build -o main .

# Run the app
CMD ["./main"]
