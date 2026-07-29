FROM golang:1.26.4-alpine

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["dlv", "debug", "./cmd/api", "--headless", "--continue", "--listen=:40000", "--api-version=2", "--accept-multiclient", "--log"]
