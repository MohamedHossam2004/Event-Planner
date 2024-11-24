FROM golang:1.18-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o authApp ./cmd/api

RUN chmod  +x /app/authApp

FROM alpine:latest

RUN mkdir /app

COPY authApp /app

CMD [ "/app/authApp" ]
