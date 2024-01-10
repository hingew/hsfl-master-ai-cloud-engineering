FROM golang:1.21-alpine

WORKDIR /build
COPY ./lib ./lib
COPY ./services/user-service ./services/user-service

WORKDIR /build/services/user-service
RUN go build -o /app/user-service

WORKDIR /app

EXPOSE 3000
CMD ["/app/user-service"]
