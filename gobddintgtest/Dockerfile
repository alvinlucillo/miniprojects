# Use an official Golang image as a build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the API binary with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api/main.go
# Use a minimal image for final deployment
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/api .

# Expose the API port
EXPOSE 8080

# Command to run the API
CMD ["./api"]
