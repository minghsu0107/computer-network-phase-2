FROM golang:1.14.8-alpine3.11 AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
RUN mkdir -p /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /app/bin /app/*.html ./
RUN apk add --no-cache bash
RUN chmod +x bin

ENTRYPOINT ["./bin"]