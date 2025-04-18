
# This dockerfile isn't used anymore

FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
