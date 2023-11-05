FROM golang:1.21-alpine

ENV USER_ENDPOINT=

WORKDIR /app
COPY ./lib ./lib
COPY ./services/user-service ./services/user-service

WORKDIR /app/services/user-service
RUN go build -o ./main

EXPOSE 3000
CMD ["/app/services/user-service/main"]