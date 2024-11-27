FROM alpine:latest

RUN mkdir /app

COPY eventApp /app

CMD [ "/app/eventApp"]