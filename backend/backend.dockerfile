FROM alpine:latest

RUN mkdir /app

COPY backendApp /app

COPY . /app

CMD [ "/app/backendApp" ]