FROM golang:1.18-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o eventApp ./cmd/api

RUN chmod +x /app/eventApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/eventApp /app

CMD [ "/app/eventApp" ]
