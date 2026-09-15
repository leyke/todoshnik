FROM golang:1.26.4-alpine AS builder

ARG APP_VERSION=dev
ARG APP_COMMIT_HASH=unknown

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build \
    -ldflags="-X main.AppVersion=${APP_VERSION} -X main.CommitHash=${APP_COMMIT_HASH}" \
    -o api ./cmd/api

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/api .

CMD ["./api"]