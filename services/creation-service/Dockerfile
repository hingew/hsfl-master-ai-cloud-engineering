FROM golang:1.21-alpine

ENV CREATION_ENDPOINT=

WORKDIR /app
COPY ./lib ./lib
COPY ./services/creation-service ./services/creation-service

WORKDIR /app/services/creation-service
RUN go build -o ./main

EXPOSE 3000
CMD ["/app/services/creation-service/main"]