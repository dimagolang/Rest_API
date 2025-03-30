FROM golang:1.23-alpine AS builder

WORKDIR /app


RUN go mod download
COPY go.mod go.sum ./


RUN apk add --no-cache git
RUN go mod tidy
RUN go mod download

COPY . .

FROM alpine:latest

RUN go build -o main.

EXPOSE 8080
CMD ["/app/main"]