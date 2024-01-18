FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o robot ./cmd

CMD ["./robot"]