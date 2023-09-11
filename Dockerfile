# First stage of the build
FROM golang:1.21 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies.
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source files
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/eth-proxy

# Second stage of the build
FROM alpine:latest

# Copy the binary from the builder stage to the final stage.
COPY --from=builder /app/eth-proxy /app/eth-proxy

# Expose the application's port
EXPOSE 8080

# Command to run the application
CMD ["/app/eth-proxy"]
