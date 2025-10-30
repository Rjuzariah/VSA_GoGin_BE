# Start from the official Go image
FROM golang:1.24

# Install build dependencies including SQLite
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev

# Set the working directory inside the container
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose port 8080
EXPOSE 8081

# Command to run the application
CMD ["./main"]