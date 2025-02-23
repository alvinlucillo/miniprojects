# Use a minimal Go base image
FROM golang:1.22 as builder

WORKDIR /app

# Copy go.mod and go.sum first, then download dependencies (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api

# Use a minimal base image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .

# Expose API port
EXPOSE 8080

# Run the application
CMD ["./main"]
